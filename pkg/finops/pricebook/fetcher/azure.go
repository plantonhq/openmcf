// Azure half of the fetcher: refreshes entries carrying an azure_selector
// from the Azure Retail Prices API (prices.azure.com, public and
// unauthenticated). The fetcher composes an OData $filter from the
// selector, follows pagination, and requires the surviving rows to agree
// on exactly one price -- the API repeats a meter once per display region,
// so several rows carrying one identical price are that one price, while
// rows disagreeing on price or unit are a hard error, never a guess.
//
// Money discipline: the API types retailPrice as a JSON float, so every
// number is decoded as json.Number and kept as its raw decimal token --
// float64 never touches a price.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	pricebookv1 "github.com/plantonhq/planton/finops/pricebook/v1"
)

// retailPricesBaseURL is the Azure Retail Prices API endpoint.
const retailPricesBaseURL = "https://prices.azure.com/api/retail/prices"

// retailPricesAPIVersion is pinned deliberately: from this version on the
// API matches filter values CASE-SENSITIVELY, which is what makes a
// selector deterministic instead of best-effort.
const retailPricesAPIVersion = "2023-01-01-preview"

// retailPriceItem is the slice of one Retail Prices API row the fetcher
// reads. Prices decode as json.Number (raw decimal tokens, never float64).
type retailPriceItem struct {
	CurrencyCode     string      `json:"currencyCode"`
	RetailPrice      json.Number `json:"retailPrice"`
	TierMinimumUnits json.Number `json:"tierMinimumUnits"`
	ArmRegionName    string      `json:"armRegionName"`
	MeterID          string      `json:"meterId"`
	MeterName        string      `json:"meterName"`
	ProductName      string      `json:"productName"`
	SkuName          string      `json:"skuName"`
	ArmSkuName       string      `json:"armSkuName"`
	ServiceName      string      `json:"serviceName"`
	UnitOfMeasure    string      `json:"unitOfMeasure"`
	Type             string      `json:"type"`
}

// retailPricesPage is one page of the Retail Prices API response.
type retailPricesPage struct {
	Items        []retailPriceItem `json:"Items"`
	NextPageLink string            `json:"NextPageLink"`
}

// azureFilter composes the OData $filter the selector describes. A meter_id
// anchors alone (plus priceType); otherwise every populated name field
// becomes an equality predicate, in a stable order so the recorded
// price_source is reproducible.
func azureFilter(selector *pricebookv1.AzurePriceSelector) string {
	priceType := selector.GetPriceType()
	if priceType == "" {
		priceType = "Consumption"
	}
	var parts []string
	if selector.GetMeterId() != "" {
		parts = append(parts, fmt.Sprintf("meterId eq '%s'", selector.GetMeterId()))
	} else {
		for _, predicate := range []struct{ field, value string }{
			{"serviceName", selector.GetServiceName()},
			{"armRegionName", selector.GetArmRegionName()},
			{"armSkuName", selector.GetArmSkuName()},
			{"meterName", selector.GetMeterName()},
			{"productName", selector.GetProductName()},
			{"skuName", selector.GetSkuName()},
		} {
			if predicate.value != "" {
				parts = append(parts, fmt.Sprintf("%s eq '%s'", predicate.field, predicate.value))
			}
		}
	}
	parts = append(parts, fmt.Sprintf("priceType eq '%s'", priceType))
	return strings.Join(parts, " and ")
}

// azureSourceURL is the refetchable citation recorded as price_source: the
// exact query, version-pinned, in readable (unencoded) form -- the same
// convention the book's hand-verified Azure entries already use.
func azureSourceURL(filter string) string {
	return fmt.Sprintf("%s?api-version=%s&currencyCode='USD'&$filter=%s",
		retailPricesBaseURL, retailPricesAPIVersion, filter)
}

// fetchAzurePrices runs one filter against the Retail Prices API and
// returns every row across all pages.
func fetchAzurePrices(client *http.Client, filter string) ([]retailPriceItem, error) {
	query := url.Values{}
	query.Set("api-version", retailPricesAPIVersion)
	query.Set("currencyCode", "'USD'")
	query.Set("$filter", filter)
	next := retailPricesBaseURL + "?" + query.Encode()

	var items []retailPriceItem
	for next != "" {
		page, err := getRetailPricesPage(client, next)
		if err != nil {
			return nil, err
		}
		items = append(items, page.Items...)
		next = page.NextPageLink
	}
	return items, nil
}

func getRetailPricesPage(client *http.Client, pageURL string) (*retailPricesPage, error) {
	response, err := client.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, errors.Errorf("GET %s: %s: %s", pageURL, response.Status, strings.TrimSpace(string(body)))
	}
	decoder := json.NewDecoder(response.Body)
	decoder.UseNumber()
	page := &retailPricesPage{}
	if err := decoder.Decode(page); err != nil {
		return nil, errors.Wrapf(err, "decoding %s", pageURL)
	}
	return page, nil
}

// refreshAzureEntry rewrites one entry's price, source, and date from its
// selector's matches in the Retail Prices API.
func refreshAzureEntry(client *http.Client, entry *pricebookv1.PriceBookEntry, today string) error {
	selector := entry.GetAzureSelector()
	filter := azureFilter(selector)

	items, err := fetchAzurePrices(client, filter)
	if err != nil {
		return err
	}

	// Tier discipline: a tiered meter must name the tier it prices; an
	// untiered meter must not need to.
	var survivors []retailPriceItem
	tiers := map[string]bool{}
	for _, item := range items {
		tier := trimDecimal(item.TierMinimumUnits.String())
		tiers[tier] = true
		if selector.GetTierMinimumUnits() != "" && tier != trimDecimal(selector.GetTierMinimumUnits()) {
			continue
		}
		survivors = append(survivors, item)
	}
	if selector.GetTierMinimumUnits() == "" && len(tiers) > 1 {
		return errors.Errorf("filter %q matches %d tiers %v -- disambiguate with tier_minimum_units", filter, len(tiers), keys(tiers))
	}
	if len(survivors) == 0 {
		return errors.Errorf("filter %q matches no prices", filter)
	}

	// Agreement discipline: the API repeats one meter across display
	// regions; identical rows are one price, disagreeing rows are a bug in
	// the selector.
	price := ""
	unit := ""
	for _, item := range survivors {
		if item.CurrencyCode != entry.GetCurrency() {
			return errors.Errorf("filter %q matched a %s price, entry is %s", filter, item.CurrencyCode, entry.GetCurrency())
		}
		itemPrice := trimDecimal(item.RetailPrice.String())
		if !decimalPattern.MatchString(itemPrice) {
			return errors.Errorf("fetched price %q is not a plain decimal", item.RetailPrice.String())
		}
		if price == "" {
			price, unit = itemPrice, item.UnitOfMeasure
			continue
		}
		if itemPrice != price || item.UnitOfMeasure != unit {
			return errors.Errorf("filter %q matches disagreeing prices (%s per %q vs %s per %q) -- tighten the selector",
				filter, price, unit, itemPrice, item.UnitOfMeasure)
		}
	}

	entry.ListUnitPrice = price
	entry.PriceSource = azureSourceURL(filter)
	entry.RetrievedOn = today
	return nil
}

func keys(set map[string]bool) []string {
	var out []string
	for k := range set {
		out = append(out, k)
	}
	return out
}
