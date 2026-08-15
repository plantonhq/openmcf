// Command fetcher entrypoint: walks every price book, refreshes the
// entries that carry a machine selector, and rewrites the refreshed books
// in their canonical form (preserving each book's leading comment header).
// See fetch.go for the selection contract.
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	pricebookv1 "github.com/plantonhq/planton/finops/pricebook/v1"
	"github.com/plantonhq/planton/pkg/finops/pricebook"
	"github.com/plantonhq/planton/pkg/yamlemit"
)

func main() {
	repoRoot, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	providers, err := pricebook.Discover(repoRoot)
	if err != nil {
		fatal(err)
	}

	client := newClient()
	today := time.Now().Format("2006-01-02")
	// One offer document serves every entry that selects from it.
	offers := map[string]*regionalOffer{}
	// The GCP catalog authenticates and caches lazily: books without GCP
	// selectors never demand a credential.
	var catalog *gcpCatalog
	refreshed := 0
	var problems []string

	for _, provider := range providers {
		book, err := pricebook.Load(repoRoot, provider)
		if err != nil {
			fatal(err)
		}

		touched := false
		for _, entry := range book.GetSpec().GetEntries() {
			var refreshErr error
			switch {
			case entry.GetAwsSelector() != nil:
				refreshErr = refreshAwsEntry(client, offers, entry, today)
			case entry.GetAzureSelector() != nil:
				refreshErr = refreshAzureEntry(client, entry, today)
			case entry.GetGcpSelector() != nil:
				if catalog == nil {
					catalog, err = newGcpCatalog(client)
					if err != nil {
						fatal(err)
					}
				}
				refreshErr = refreshGcpEntry(catalog, entry, today)
			default:
				continue
			}
			if refreshErr != nil {
				problems = append(problems, fmt.Sprintf("%s price book entry %q: %v", provider, entry.GetName(), refreshErr))
				continue
			}
			touched = true
			refreshed++
		}

		if !touched {
			continue
		}
		path := pricebook.Path(repoRoot, provider)
		header, err := readHeader(path)
		if err != nil {
			fatal(err)
		}
		if err := os.WriteFile(path, []byte(renderBook(header, book)), 0644); err != nil {
			fatal(err)
		}
	}

	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "%d entr(ies) could not be refreshed -- fix the selectors:\n", len(problems))
		for _, problem := range problems {
			fmt.Fprintf(os.Stderr, "  - %s\n", problem)
		}
		os.Exit(1)
	}
	fmt.Printf("refreshed %d price book entr(ies)\n", refreshed)
}

// readHeader returns a book's leading comment block, which the fetcher
// preserves verbatim -- the header is authored prose, not generated data.
func readHeader(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var header strings.Builder
	for _, line := range strings.SplitAfter(string(content), "\n") {
		if !strings.HasPrefix(line, "#") {
			break
		}
		header.WriteString(line)
	}
	return header.String(), nil
}

// renderBook renders a price book in its canonical form: schema field
// order, decimals and dates quoted, empty fields omitted.
func renderBook(header string, book *pricebookv1.PriceBook) string {
	var b strings.Builder
	b.WriteString(header)
	b.WriteString("apiVersion: finops.planton.dev/v1\n")
	b.WriteString("kind: PriceBook\n")
	b.WriteString("metadata:\n")
	yamlemit.WriteKV(&b, 2, "name", book.GetMetadata().GetName(), false)
	b.WriteString("spec:\n")
	b.WriteString("  entries:\n")
	for _, entry := range book.GetSpec().GetEntries() {
		yamlemit.WriteKV(&b, 4, "- name", entry.GetName(), false)
		yamlemit.WriteKV(&b, 6, "service_name", entry.GetServiceName(), false)
		yamlemit.WriteKV(&b, 6, "pricing_unit", entry.GetPricingUnit(), false)
		yamlemit.WriteKV(&b, 6, "region", entry.GetRegion(), false)
		yamlemit.WriteKV(&b, 6, "currency", entry.GetCurrency(), false)
		yamlemit.WriteKV(&b, 6, "list_unit_price", entry.GetListUnitPrice(), true)
		yamlemit.WriteKV(&b, 6, "price_source", entry.GetPriceSource(), strings.Contains(entry.GetPriceSource(), " "))
		yamlemit.WriteKV(&b, 6, "retrieved_on", entry.GetRetrievedOn(), true)
		// Attributes are hand-authored lookup identity, never fetched --
		// but a refresh rewrites the whole book, so dropping them here
		// would silently destroy every value-keyed price. Sorted keys keep
		// the rendering deterministic.
		if attributes := entry.GetAttributes(); len(attributes) > 0 {
			b.WriteString("      attributes:\n")
			keys := make([]string, 0, len(attributes))
			for key := range attributes {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				yamlemit.WriteKV(&b, 8, key, attributes[key], false)
			}
		}
		if selector := entry.GetAwsSelector(); selector != nil {
			b.WriteString("      aws_selector:\n")
			yamlemit.WriteKV(&b, 8, "offer_code", selector.GetOfferCode(), false)
			yamlemit.WriteKV(&b, 8, "region_code", selector.GetRegionCode(), false)
			yamlemit.WriteKV(&b, 8, "usage_type", selector.GetUsageType(), false)
			if selector.GetOperation() != "" {
				yamlemit.WriteKV(&b, 8, "operation", selector.GetOperation(), false)
			}
			if selector.GetDescriptionContains() != "" {
				yamlemit.WriteKV(&b, 8, "description_contains", selector.GetDescriptionContains(), false)
			}
		}
		if selector := entry.GetAzureSelector(); selector != nil {
			b.WriteString("      azure_selector:\n")
			for _, field := range []struct{ key, value string }{
				{"service_name", selector.GetServiceName()},
				{"arm_region_name", selector.GetArmRegionName()},
				{"arm_sku_name", selector.GetArmSkuName()},
				{"meter_name", selector.GetMeterName()},
				{"product_name", selector.GetProductName()},
				{"sku_name", selector.GetSkuName()},
				{"price_type", selector.GetPriceType()},
			} {
				if field.value != "" {
					yamlemit.WriteKV(&b, 8, field.key, field.value, false)
				}
			}
			if selector.GetTierMinimumUnits() != "" {
				yamlemit.WriteKV(&b, 8, "tier_minimum_units", selector.GetTierMinimumUnits(), true)
			}
			if selector.GetMeterId() != "" {
				yamlemit.WriteKV(&b, 8, "meter_id", selector.GetMeterId(), false)
			}
		}
		if selector := entry.GetGcpSelector(); selector != nil {
			b.WriteString("      gcp_selector:\n")
			yamlemit.WriteKV(&b, 8, "service_id", selector.GetServiceId(), false)
			yamlemit.WriteKV(&b, 8, "sku_id", selector.GetSkuId(), false)
			if selector.GetStartUsageAmount() != "" {
				yamlemit.WriteKV(&b, 8, "start_usage_amount", selector.GetStartUsageAmount(), true)
			}
		}
	}
	return b.String()
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "price-book fetcher:", err)
	os.Exit(1)
}
