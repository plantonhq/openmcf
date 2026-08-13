// Package main (estimategen) generates the committed per-component cost
// estimates at catalog/_pricing/estimates/ by joining each component's
// estimate model (catalog/_pricing/models/, the authored quantities and
// prose) with its provider's price book (catalog/_pricing/pricebook/, the
// pinned unit prices). The generator computes every line cost and preset
// total exactly (big.Rat, no floats), orders lines largest cost first, and
// emits deterministic YAML: identical inputs produce byte-identical files,
// and the drift test in this package holds committed files to that promise.
//
// The generator is also the cross-artifact coherence gate -- the one
// program that reads the cost profiles, the models, and the price books
// together -- and it fails loudly when they disagree: a price reference
// that resolves to no book entry, a book entry no model references (a dead
// price), a book entry whose pricing_unit contradicts the component's
// cost.yaml declaration for that meter, a currency mismatch, or a
// regionally priced entry used under a different region assumption.
// Disagreement between sources is a build failure, never a published
// number.
package main

import (
	"fmt"
	"math/big"
	"path"
	"regexp"
	"sort"
	"strings"

	estimatemodelv1 "github.com/plantonhq/planton/finops/componentcostestimatemodel/v1"
	pricebookv1 "github.com/plantonhq/planton/finops/pricebook/v1"
	"github.com/plantonhq/planton/pkg/finops/costestimate"
	"github.com/plantonhq/planton/pkg/finops/costprofile"
	"github.com/plantonhq/planton/pkg/finops/estimatemodel"
	"github.com/plantonhq/planton/pkg/finops/pricebook"
	"github.com/plantonhq/planton/pkg/yamlemit"
)

// header opens every generated estimate file. It names the regeneration
// command and both input documents so a reader lands on the authoring
// surfaces, never on the generated output.
const header = `# GENERATED -- DO NOT EDIT. Regenerate with ` + "`make generate-cost-estimates`" + `.
# This file joins the component's estimate model (../models/, the authored
# quantities with their derivations) with the provider's price book
# (../pricebook/, pinned unit prices with source URL and retrieval date):
# every line cost and preset total is computed exactly, lines are ordered
# largest cost first, and the catalog-data conformance gate re-verifies all
# arithmetic. Schema: finops/componentcostestimate/v1.
`

var decimalPattern = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)

// Summary is one whole-tree generation, rendered in memory -- writing is
// the caller's decision so tests can byte-compare without touching disk.
type Summary struct {
	// Files maps repo-relative output paths to rendered content.
	Files map[string]string
	// Problems are cross-artifact coherence violations. Any problem fails
	// both the CLI run and the drift gate; generated content is only
	// trustworthy when this is empty.
	Problems []string
}

// SortedPaths returns the output paths in deterministic order.
func (s *Summary) SortedPaths() []string {
	paths := make([]string, 0, len(s.Files))
	for p := range s.Files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// Generate renders every component's estimate document from its model and
// its provider's price book. It reads the repo tree and writes nothing.
func Generate(repoRoot string) (*Summary, error) {
	summary := &Summary{Files: map[string]string{}}

	components, err := estimatemodel.Discover(repoRoot)
	if err != nil {
		return nil, err
	}

	profiles, err := costprofile.Discover(repoRoot)
	if err != nil {
		return nil, err
	}
	providerOf := map[string]string{}
	for provider, providerComponents := range profiles {
		for _, component := range providerComponents {
			providerOf[component] = provider
		}
	}

	// referenced tracks which book entries some model actually cites, per
	// provider -- the dead-price check reads it after all components render.
	referenced := map[string]map[string]bool{}

	for _, component := range components {
		provider, ok := providerOf[component]
		if !ok {
			summary.problemf("%s: estimate model exists but the component ships no cost.yaml -- a model without a cost anatomy is unreviewable", component)
			continue
		}
		model, err := estimatemodel.Load(repoRoot, component)
		if err != nil {
			return nil, err
		}
		profile, err := costprofile.Load(repoRoot, provider, component)
		if err != nil {
			return nil, err
		}
		entries, err := pricebook.Entries(repoRoot, provider)
		if err != nil {
			return nil, err
		}
		if referenced[provider] == nil {
			referenced[provider] = map[string]bool{}
		}

		meterUnits := map[string]map[string]bool{}
		meterService := map[string]string{}
		for _, charge := range profile.GetSpec().GetBaselineCharges() {
			addMeterUnit(meterUnits, charge.GetSkuMeter(), charge.GetPricingUnit())
			meterService[charge.GetSkuMeter()] = charge.GetServiceName()
		}
		for _, driver := range profile.GetSpec().GetCostDrivers() {
			addMeterUnit(meterUnits, driver.GetSkuMeter(), driver.GetPricingUnit())
		}

		doc, ok := renderComponent(summary, component, model, entries, meterUnits, meterService, referenced[provider])
		if ok {
			summary.Files[path.Join(costestimate.Dir, component+".yaml")] = doc
		}
	}

	// The dead-price sweep: an entry no model references is dead weight --
	// either an orphaned pin (its component's model moved off it) or a slug
	// typo waiting to bite. The book stays exactly as large as its readers.
	providers, err := pricebook.Discover(repoRoot)
	if err != nil {
		return nil, err
	}
	for _, provider := range providers {
		book, err := pricebook.Load(repoRoot, provider)
		if err != nil {
			return nil, err
		}
		for _, entry := range book.GetSpec().GetEntries() {
			if !referenced[provider][entry.GetName()] {
				summary.problemf("%s price book: entry %q is referenced by no estimate model -- a dead price; remove it or fix the model that should cite it", provider, entry.GetName())
			}
		}
	}

	return summary, nil
}

// renderComponent renders one component's estimate document. It returns
// ok=false when a coherence problem makes the document untrustworthy.
func renderComponent(
	summary *Summary,
	component string,
	model *estimatemodelv1.ComponentCostEstimateModel,
	entries map[string]*pricebookv1.PriceBookEntry,
	meterUnits map[string]map[string]bool,
	meterService map[string]string,
	referenced map[string]bool,
) (string, bool) {
	var b strings.Builder
	b.WriteString(header)
	b.WriteString("apiVersion: finops.planton.dev/v1\n")
	b.WriteString("kind: ComponentCostEstimate\n")
	b.WriteString("metadata:\n")
	yamlemit.WriteKV(&b, 2, "name", component, false)
	b.WriteString("spec:\n")
	b.WriteString("  presets:\n")

	ok := true
	for _, preset := range model.GetSpec().GetPresets() {
		yamlemit.WriteKV(&b, 4, "- preset", preset.GetPreset(), false)
		if preset.GetRegionAssumption() != "" {
			yamlemit.WriteKV(&b, 6, "region_assumption", preset.GetRegionAssumption(), false)
		}
		if preset.GetCurrency() != "" {
			yamlemit.WriteKV(&b, 6, "currency", preset.GetCurrency(), false)
		}
		if preset.GetHoursPerMonth() != 0 {
			fmt.Fprintf(&b, "      hours_per_month: %d\n", preset.GetHoursPerMonth())
		}

		if footprint := preset.GetCapacityFootprint(); footprint != nil {
			b.WriteString("      capacity_footprint:\n")
			for _, kv := range []struct{ key, value string }{
				{"cpu_requests", footprint.GetCpuRequests()},
				{"memory_requests", footprint.GetMemoryRequests()},
				{"cpu_limits", footprint.GetCpuLimits()},
				{"memory_limits", footprint.GetMemoryLimits()},
				{"persistent_storage", footprint.GetPersistentStorage()},
				{"basis", footprint.GetBasis()},
			} {
				if kv.value != "" {
					yamlemit.WriteKV(&b, 8, kv.key, kv.value, false)
				}
			}
		} else {
			lines, total, linesOK := priceLines(summary, component, preset, entries, meterUnits, meterService, referenced)
			ok = ok && linesOK
			if len(lines) > 0 {
				b.WriteString("      line_items:\n")
				for _, line := range lines {
					yamlemit.WriteKV(&b, 8, "- service_name", line.serviceName, false)
					yamlemit.WriteKV(&b, 10, "sku_meter", line.skuMeter, false)
					yamlemit.WriteKV(&b, 10, "pricing_unit", line.pricingUnit, false)
					yamlemit.WriteKV(&b, 10, "pricing_quantity", line.pricingQuantity, true)
					yamlemit.WriteKV(&b, 10, "quantity_basis", line.quantityBasis, false)
					yamlemit.WriteKV(&b, 10, "list_unit_price", line.listUnitPrice, true)
					// A URL with embedded spaces (the Azure Retail Prices
					// filter style) is plain-scalar-legal but reads as
					// several tokens; quote it so it reads as one.
					yamlemit.WriteKV(&b, 10, "price_source", line.priceSource, strings.Contains(line.priceSource, " "))
					yamlemit.WriteKV(&b, 10, "retrieved_on", line.retrievedOn, true)
					yamlemit.WriteKV(&b, 10, "list_cost", line.listCost, true)
				}
			}
			yamlemit.WriteKV(&b, 6, "total_list_cost", total, true)
		}

		if len(preset.GetExclusions()) > 0 {
			b.WriteString("      exclusions:\n")
			for _, exclusion := range preset.GetExclusions() {
				yamlemit.WriteListItem(&b, 8, exclusion)
			}
		}
		if preset.GetNotes() != "" {
			yamlemit.WriteKV(&b, 6, "notes", preset.GetNotes(), false)
		}
	}

	return b.String(), ok
}

// renderedLine is one priced estimate line, carrying its exact cost for
// the largest-first ordering.
type renderedLine struct {
	serviceName     string
	skuMeter        string
	pricingUnit     string
	pricingQuantity string
	quantityBasis   string
	listUnitPrice   string
	priceSource     string
	retrievedOn     string
	listCost        string
	cost            *big.Rat
}

// priceLines joins a preset's quantity lines with the price book, checks
// every cross-artifact agreement, and returns the priced lines ordered
// largest cost first plus the exact total.
func priceLines(
	summary *Summary,
	component string,
	preset *estimatemodelv1.PresetEstimateModel,
	entries map[string]*pricebookv1.PriceBookEntry,
	meterUnits map[string]map[string]bool,
	meterService map[string]string,
	referenced map[string]bool,
) ([]renderedLine, string, bool) {
	ok := true
	var lines []renderedLine
	total := new(big.Rat)

	for _, quantity := range preset.GetQuantityLines() {
		where := fmt.Sprintf("%s preset %s line %q", component, preset.GetPreset(), quantity.GetSkuMeter())

		entry, found := entries[quantity.GetPrice()]
		if !found {
			summary.problemf("%s: price reference %q resolves to no price book entry", where, quantity.GetPrice())
			ok = false
			continue
		}
		referenced[entry.GetName()] = true

		// The three-source agreements. The cost profile owns the meter's
		// unit vocabulary; the model owns the estimate's currency and
		// region assumption; the entry must fit both or the join is wrong.
		if units := meterUnits[quantity.GetSkuMeter()]; len(units) > 0 && !units[entry.GetPricingUnit()] {
			summary.problemf("%s: price entry %q is per %q, but the component's cost.yaml declares this meter in %s",
				where, entry.GetName(), entry.GetPricingUnit(), joinSorted(units))
			ok = false
		}
		if service := meterService[quantity.GetSkuMeter()]; service != "" && service != entry.GetServiceName() {
			summary.problemf("%s: price entry %q bills service %q, but the component's cost.yaml declares this meter under %q",
				where, entry.GetName(), entry.GetServiceName(), service)
			ok = false
		}
		if entry.GetCurrency() != preset.GetCurrency() {
			summary.problemf("%s: price entry %q is in %s, but the preset estimates in %s",
				where, entry.GetName(), entry.GetCurrency(), preset.GetCurrency())
			ok = false
		}
		if entry.GetRegion() != pricebook.GlobalRegion && entry.GetRegion() != preset.GetRegionAssumption() {
			summary.problemf("%s: price entry %q is priced for %s, but the preset assumes %s",
				where, entry.GetName(), entry.GetRegion(), preset.GetRegionAssumption())
			ok = false
		}

		quantityRat, qOK := parseDecimal(quantity.GetPricingQuantity())
		priceRat, pOK := parseDecimal(entry.GetListUnitPrice())
		if !qOK || !pOK {
			summary.problemf("%s: quantity %q or unit price %q is not a plain decimal string",
				where, quantity.GetPricingQuantity(), entry.GetListUnitPrice())
			ok = false
			continue
		}
		cost := new(big.Rat).Mul(quantityRat, priceRat)
		total.Add(total, cost)

		lines = append(lines, renderedLine{
			serviceName:     entry.GetServiceName(),
			skuMeter:        quantity.GetSkuMeter(),
			pricingUnit:     entry.GetPricingUnit(),
			pricingQuantity: quantity.GetPricingQuantity(),
			quantityBasis:   quantity.GetQuantityBasis(),
			listUnitPrice:   entry.GetListUnitPrice(),
			priceSource:     entry.GetPriceSource(),
			retrievedOn:     entry.GetRetrievedOn(),
			listCost:        exactDecimal(cost, 4),
			cost:            cost,
		})
	}

	sort.SliceStable(lines, func(i, j int) bool { return lines[i].cost.Cmp(lines[j].cost) > 0 })

	if total.Sign() == 0 {
		return lines, "0.00", ok
	}
	return lines, exactDecimal(total, 4), ok
}

// addMeterUnit records a cost profile's declared pricing unit for a meter.
// Drivers may declare no unit (they switch between pricing models); only
// concrete declarations constrain the price book.
func addMeterUnit(meterUnits map[string]map[string]bool, meter, unit string) {
	if unit == "" {
		return
	}
	if meterUnits[meter] == nil {
		meterUnits[meter] = map[string]bool{}
	}
	meterUnits[meter][unit] = true
}

// parseDecimal parses a plain decimal string into an exact rational.
func parseDecimal(value string) (*big.Rat, bool) {
	if !decimalPattern.MatchString(value) {
		return nil, false
	}
	rat, ok := new(big.Rat).SetString(value)
	return rat, ok
}

// exactDecimal renders a rational as an exact decimal string with at least
// minFrac fraction digits. Products and sums of finite decimals are always
// finite decimals, so the denominator is a power of 2s and 5s and the
// rendering never rounds.
func exactDecimal(r *big.Rat, minFrac int) string {
	denom := new(big.Int).Set(r.Denom())
	scale := 0
	for _, factor := range []int64{2, 5} {
		f := big.NewInt(factor)
		count := 0
		mod := new(big.Int)
		for {
			quo, m := new(big.Int).QuoRem(denom, f, mod)
			if m.Sign() != 0 {
				break
			}
			denom.Set(quo)
			count++
		}
		if count > scale {
			scale = count
		}
	}
	if scale < minFrac {
		scale = minFrac
	}
	return r.FloatString(scale)
}

// joinSorted renders a unit set for error messages, deterministically.
func joinSorted(set map[string]bool) string {
	var values []string
	for value := range set {
		values = append(values, fmt.Sprintf("%q", value))
	}
	sort.Strings(values)
	return strings.Join(values, ", ")
}

func (s *Summary) problemf(format string, args ...any) {
	s.Problems = append(s.Problems, fmt.Sprintf(format, args...))
}
