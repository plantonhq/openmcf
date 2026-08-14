// GCP half of the fetcher: refreshes entries carrying a gcp_selector from
// the Cloud Billing Catalog API (cloudbilling.googleapis.com). Google
// assigns every price a stable SKU ID, so selection is direct: the fetcher
// lists the selector's service once (cached across entries), finds the SKU
// by ID, verifies it still serves the entry's region, and converts the
// exact units+nanos price to the decimal string -- floats never touch the
// money.
//
// Credential contract: the API requires authentication. The fetcher uses
// GCP_BILLING_API_KEY when set (a plain API key, passed as ?key=);
// otherwise it mints an application-default-credentials access token via
// `gcloud auth application-default print-access-token`. Both are
// read-only, free-tier credentials; the hard error below says exactly how
// to provision either.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	pricebookv1 "github.com/plantonhq/planton/finops/pricebook/v1"
)

// cloudBillingBaseURL is the Cloud Billing Catalog API endpoint.
const cloudBillingBaseURL = "https://cloudbilling.googleapis.com/v1"

// gcpCatalog authenticates to and caches from the Cloud Billing Catalog
// API. SKU lists are fetched once per service and reused across entries.
type gcpCatalog struct {
	client *http.Client
	apiKey string
	token  string
	skus   map[string][]gcpSku
}

// gcpSku is the slice of one Catalog API SKU the fetcher reads.
type gcpSku struct {
	SkuID          string   `json:"skuId"`
	Description    string   `json:"description"`
	ServiceRegions []string `json:"serviceRegions"`
	PricingInfo    []struct {
		PricingExpression struct {
			UsageUnit   string `json:"usageUnit"`
			TieredRates []struct {
				StartUsageAmount json.Number `json:"startUsageAmount"`
				UnitPrice        struct {
					CurrencyCode string          `json:"currencyCode"`
					Units        json.RawMessage `json:"units"`
					Nanos        int64           `json:"nanos"`
				} `json:"unitPrice"`
			} `json:"tieredRates"`
		} `json:"pricingExpression"`
	} `json:"pricingInfo"`
}

// gcpSkuPage is one page of the Catalog API's skus.list response.
type gcpSkuPage struct {
	Skus          []gcpSku `json:"skus"`
	NextPageToken string   `json:"nextPageToken"`
}

// newGcpCatalog resolves the credential contract lazily-never: the first
// GCP entry needs it, so a book with GCP selectors and no credential fails
// before any network call, with instructions.
func newGcpCatalog(client *http.Client) (*gcpCatalog, error) {
	catalog := &gcpCatalog{client: client, skus: map[string][]gcpSku{}}
	if key := os.Getenv("GCP_BILLING_API_KEY"); key != "" {
		catalog.apiKey = key
		return catalog, nil
	}
	var out bytes.Buffer
	command := exec.Command("gcloud", "auth", "application-default", "print-access-token")
	command.Stdout = &out
	if err := command.Run(); err != nil {
		return nil, errors.New("the Cloud Billing Catalog API needs a credential: set GCP_BILLING_API_KEY " +
			"(an API key from any GCP project, no billing account required) or run " +
			"`gcloud auth application-default login` so the fetcher can mint an access token")
	}
	catalog.token = strings.TrimSpace(out.String())
	return catalog, nil
}

// serviceSkus lists (and caches) every SKU of one Cloud Billing service.
func (c *gcpCatalog) serviceSkus(serviceID string) ([]gcpSku, error) {
	if skus, ok := c.skus[serviceID]; ok {
		return skus, nil
	}
	var skus []gcpSku
	pageToken := ""
	for {
		query := url.Values{}
		query.Set("pageSize", "5000")
		if pageToken != "" {
			query.Set("pageToken", pageToken)
		}
		if c.apiKey != "" {
			query.Set("key", c.apiKey)
		}
		pageURL := fmt.Sprintf("%s/services/%s/skus?%s", cloudBillingBaseURL, serviceID, query.Encode())

		request, err := http.NewRequest(http.MethodGet, pageURL, nil)
		if err != nil {
			return nil, err
		}
		if c.token != "" {
			request.Header.Set("Authorization", "Bearer "+c.token)
		}
		response, err := c.client.Do(request)
		if err != nil {
			return nil, err
		}
		page := &gcpSkuPage{}
		decodeErr := func() error {
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(response.Body)
				return errors.Errorf("GET %s: %s: %s", pageURL, response.Status, strings.TrimSpace(string(body)))
			}
			decoder := json.NewDecoder(response.Body)
			decoder.UseNumber()
			return decoder.Decode(page)
		}()
		if decodeErr != nil {
			return nil, errors.Wrapf(decodeErr, "listing SKUs of service %s", serviceID)
		}
		skus = append(skus, page.Skus...)
		if page.NextPageToken == "" {
			break
		}
		pageToken = page.NextPageToken
	}
	c.skus[serviceID] = skus
	return skus, nil
}

// refreshGcpEntry rewrites one entry's price, source, and date from its
// selector's SKU in the Catalog API.
func refreshGcpEntry(catalog *gcpCatalog, entry *pricebookv1.PriceBookEntry, today string) error {
	selector := entry.GetGcpSelector()
	skus, err := catalog.serviceSkus(selector.GetServiceId())
	if err != nil {
		return err
	}

	var matched *gcpSku
	for i := range skus {
		if skus[i].SkuID == selector.GetSkuId() {
			matched = &skus[i]
			break
		}
	}
	if matched == nil {
		return errors.Errorf("service %s carries no SKU %s", selector.GetServiceId(), selector.GetSkuId())
	}

	region := entry.GetRegion()
	regionServed := false
	for _, served := range matched.ServiceRegions {
		if served == region {
			regionServed = true
			break
		}
	}
	if !regionServed {
		return errors.Errorf("SKU %s (%q) does not serve region %q (serves %v)",
			selector.GetSkuId(), matched.Description, region, matched.ServiceRegions)
	}

	if len(matched.PricingInfo) != 1 {
		return errors.Errorf("SKU %s carries %d pricingInfo blocks, want exactly 1", selector.GetSkuId(), len(matched.PricingInfo))
	}
	rates := matched.PricingInfo[0].PricingExpression.TieredRates
	var picked []int
	for i, rate := range rates {
		if selector.GetStartUsageAmount() != "" &&
			trimDecimal(rate.StartUsageAmount.String()) != trimDecimal(selector.GetStartUsageAmount()) {
			continue
		}
		picked = append(picked, i)
	}
	if len(picked) != 1 {
		return errors.Errorf("SKU %s has %d tiered rates and the selector picks %d -- set start_usage_amount to exactly one",
			selector.GetSkuId(), len(rates), len(picked))
	}
	rate := rates[picked[0]]

	if rate.UnitPrice.CurrencyCode != entry.GetCurrency() {
		return errors.Errorf("SKU %s prices in %s, entry is %s", selector.GetSkuId(), rate.UnitPrice.CurrencyCode, entry.GetCurrency())
	}
	price, err := moneyToDecimal(rate.UnitPrice.Units, rate.UnitPrice.Nanos)
	if err != nil {
		return errors.Wrapf(err, "SKU %s unit price", selector.GetSkuId())
	}

	entry.ListUnitPrice = price
	entry.PriceSource = fmt.Sprintf("%s/services/%s/skus", cloudBillingBaseURL, selector.GetServiceId())
	entry.RetrievedOn = today
	return nil
}

// moneyToDecimal renders a google.type.Money units+nanos pair as the exact
// decimal string ("0"+100000000 -> "0.1"). Units arrives as a JSON string
// (proto3 int64 mapping) or bare number depending on the emitter; both are
// parsed as integers, never floats. Negative prices do not exist in the
// catalog and are rejected.
func moneyToDecimal(unitsRaw json.RawMessage, nanos int64) (string, error) {
	unitsToken := strings.Trim(strings.TrimSpace(string(unitsRaw)), `"`)
	if unitsToken == "" {
		unitsToken = "0"
	}
	units, err := strconv.ParseInt(unitsToken, 10, 64)
	if err != nil {
		return "", errors.Errorf("units %q is not an integer", unitsToken)
	}
	if units < 0 || nanos < 0 {
		return "", errors.Errorf("negative money (units %d, nanos %d)", units, nanos)
	}
	if nanos > 999999999 {
		return "", errors.Errorf("nanos %d out of range", nanos)
	}
	price := trimDecimal(fmt.Sprintf("%d.%09d", units, nanos))
	if !decimalPattern.MatchString(price) {
		return "", errors.Errorf("computed price %q is not a plain decimal", price)
	}
	return price, nil
}
