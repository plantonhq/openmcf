// Package main (the price-book fetcher) refreshes every price-book entry
// that carries a machine selector from the provider's public price API --
// the AWS Price List bulk API today. For each refreshed entry it rewrites
// the unit price, the price_source (pinned to the offer's VERSIONED
// document URL, a citation that stays fetchable after prices change), and
// the retrieval date. Entries without a selector -- and whole books without
// any -- are left untouched: they are refreshed by hand, or by a future
// provider fetcher.
//
// Selection is deterministic or dead: a selector that matches zero or
// several products, or zero or several price dimensions, is a hard error,
// never a guess. Free-tier $0 dimensions sit beside the standard rate in
// the offer files, so tiered SKUs disambiguate with description_contains
// (list-price semantics never net free tiers out).
//
// Run through `make generate-price-book`. Requires network access; CI never
// fetches -- the lint lane validates the committed snapshot. An unchanged
// upstream offer version rewrites a byte-identical book (apart from
// retrieved_on, which honestly records the re-verification date).
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	pricebookv1 "github.com/plantonhq/planton/finops/pricebook/v1"
)

// priceListBaseURL is the AWS Price List bulk API host. The offer files
// are public and unauthenticated; the regional files are small (tens of
// KB) where the all-region files can be enormous.
const priceListBaseURL = "https://pricing.us-east-1.amazonaws.com"

var decimalPattern = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)

// offerIndex is the per-offer version index at
// /offers/v1.0/aws/<offer>/index.json.
type offerIndex struct {
	CurrentVersion string `json:"currentVersion"`
}

// offerFile is the slice of a regional offer document the fetcher reads:
// product attributes for selection, on-demand price dimensions for the
// price.
type offerFile struct {
	Products map[string]struct {
		Attributes map[string]string `json:"attributes"`
	} `json:"products"`
	Terms struct {
		OnDemand map[string]map[string]struct {
			PriceDimensions map[string]struct {
				Description  string            `json:"description"`
				PricePerUnit map[string]string `json:"pricePerUnit"`
			} `json:"priceDimensions"`
		} `json:"OnDemand"`
	} `json:"terms"`
}

// regionalOffer is one downloaded regional offer document plus the
// versioned URL it was read from (the price_source every entry it prices
// will cite).
type regionalOffer struct {
	file         *offerFile
	versionedURL string
}

// globalRegionCode in a selector means the SKU is priced globally. Globally
// priced SKUs (Route 53 hosted zones) appear only in the offer-level price
// document, never in the regional slices, so the fetcher reads the whole
// offer for them.
const globalRegionCode = "global"

// fetchOffer resolves an offer's current version and downloads its
// versioned price document -- the regional slice for regional SKUs, the
// offer-level document for globally priced ones.
func fetchOffer(client *http.Client, offerCode, regionCode string) (*regionalOffer, error) {
	indexURL := fmt.Sprintf("%s/offers/v1.0/aws/%s/index.json", priceListBaseURL, offerCode)
	var index offerIndex
	if err := getJSON(client, indexURL, &index); err != nil {
		return nil, errors.Wrapf(err, "resolving current version of offer %s", offerCode)
	}
	if index.CurrentVersion == "" {
		return nil, errors.Errorf("offer %s version index carries no currentVersion", offerCode)
	}

	versionedURL := fmt.Sprintf("%s/offers/v1.0/aws/%s/%s/%s/index.json",
		priceListBaseURL, offerCode, index.CurrentVersion, regionCode)
	if regionCode == globalRegionCode {
		versionedURL = fmt.Sprintf("%s/offers/v1.0/aws/%s/%s/index.json",
			priceListBaseURL, offerCode, index.CurrentVersion)
	}
	var file offerFile
	if err := getJSON(client, versionedURL, &file); err != nil {
		return nil, errors.Wrapf(err, "downloading offer %s %s for %s", offerCode, index.CurrentVersion, regionCode)
	}
	return &regionalOffer{file: &file, versionedURL: versionedURL}, nil
}

// refreshEntry rewrites one entry's price, source, and date from its
// selector's match in the regional offer document.
func refreshEntry(entry *pricebookv1.PriceBookEntry, offer *regionalOffer, today string) error {
	selector := entry.GetAwsSelector()

	var matched []string
	for sku, product := range offer.file.Products {
		if product.Attributes["usagetype"] != selector.GetUsageType() {
			continue
		}
		if product.Attributes["operation"] != selector.GetOperation() {
			continue
		}
		matched = append(matched, sku)
	}
	if len(matched) != 1 {
		return errors.Errorf("selector (usagetype %q, operation %q) matches %d products, want exactly 1",
			selector.GetUsageType(), selector.GetOperation(), len(matched))
	}

	var prices []string
	var descriptions []string
	for _, term := range offer.file.Terms.OnDemand[matched[0]] {
		for _, dimension := range term.PriceDimensions {
			if selector.GetDescriptionContains() != "" &&
				!strings.Contains(dimension.Description, selector.GetDescriptionContains()) {
				continue
			}
			price, ok := dimension.PricePerUnit[entry.GetCurrency()]
			if !ok {
				return errors.Errorf("price dimension %q carries no %s price", dimension.Description, entry.GetCurrency())
			}
			prices = append(prices, price)
			descriptions = append(descriptions, dimension.Description)
		}
	}
	if len(prices) != 1 {
		return errors.Errorf("selector matches %d price dimensions %q, want exactly 1 -- disambiguate with description_contains",
			len(prices), descriptions)
	}

	price := trimDecimal(prices[0])
	if !decimalPattern.MatchString(price) {
		return errors.Errorf("fetched price %q is not a plain decimal", prices[0])
	}

	entry.ListUnitPrice = price
	entry.PriceSource = offer.versionedURL
	entry.RetrievedOn = today
	return nil
}

// trimDecimal drops the bulk API's zero-padding ("0.0225000000" ->
// "0.0225") without ever rounding.
func trimDecimal(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

func getJSON(client *http.Client, url string, out any) error {
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.Errorf("GET %s: %s", url, response.Status)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func newClient() *http.Client {
	return &http.Client{Timeout: 60 * time.Second}
}
