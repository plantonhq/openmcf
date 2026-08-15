package costderivation

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	derivationv1 "github.com/plantonhq/planton/finops/componentcostderivation/v1"
	costprofilev1 "github.com/plantonhq/planton/finops/componentcostprofile/v1"
	pricebookv1 "github.com/plantonhq/planton/finops/pricebook/v1"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/finops/costprofile"
	"github.com/plantonhq/planton/pkg/finops/estimatemodel"
	"github.com/plantonhq/planton/pkg/finops/pricebook"
	"github.com/plantonhq/planton/pkg/specpath"
)

var decimalPattern = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)

// TestCostDerivationConformance holds every cost derivation to its
// contract, offline:
//
//  1. The document parses strictly, names its component (metadata.name
//     equals the filename), the component ships a cost profile, and the
//     component does NOT also ship a hand-authored estimate model -- a
//     component's quantities have exactly one home.
//  2. Cluster-capacity components carry no derivation: capacity
//     footprints are per-preset facts the derivation standard does not
//     yet express, and a monetary derivation for one would be dishonest.
//  3. Every field path anywhere in the document (region binding,
//     conditions, quantity factors, attribute bindings) resolves against
//     the served version's compiled descriptors -- a schema rename that
//     orphans a rule fails CI loudly -- and count_of paths are actually
//     repeated fields.
//  4. Every line rule names a declared meter (a baseline charge or cost
//     driver sku_meter in the component's cost.yaml), at least one
//     quantity factor, exactly one price arm (a slug resolving in the
//     provider's book, or a complete attribute lookup), and its basis
//     prose; every refusal states its reason; conditions carry an op and
//     the comparand exactly when the op compares.
//
// Whether the rules reproduce the hand-verified numbers is the estimate
// generator's replay check -- this gate keeps each derivation internally
// sound.
func TestCostDerivationConformance(t *testing.T) {
	root := repoRoot(t)
	if _, err := os.Stat(filepath.Join(root, "catalog")); err != nil {
		t.Skip("catalog source tree not present (bazel sandbox); runs under go test and the lint.catalog-data lane")
	}

	components, err := Discover(root)
	if err != nil {
		t.Fatalf("discovering cost derivations: %v", err)
	}
	if len(components) == 0 {
		t.Skip("no cost derivations authored yet")
	}

	for _, component := range components {
		component := component
		t.Run(component, func(t *testing.T) {
			derivation, err := Load(root, component)
			if err != nil {
				t.Fatalf("cost derivation: %v", err)
			}
			if derivation.GetKind() != "ComponentCostDerivation" {
				t.Fatalf("kind is %q, want ComponentCostDerivation", derivation.GetKind())
			}
			if derivation.GetMetadata().GetName() != component {
				t.Errorf("metadata.name is %q, want %q (the filename is the component's identity)",
					derivation.GetMetadata().GetName(), component)
			}

			if _, err := os.Stat(estimatemodel.Path(root, component)); err == nil {
				t.Errorf("component also ships an estimate model (%s) -- quantities have exactly one home; delete the model",
					estimatemodel.Path(root, component))
			}

			provider := componentProvider(t, root, component)
			profile, err := costprofile.Load(root, provider, component)
			if err != nil {
				t.Fatalf("the derived component must ship a cost profile: %v", err)
			}
			if profile.GetSpec().GetBillingModel() == costprofilev1.BillingModel_cluster_capacity {
				t.Fatal("cluster-capacity components carry no derivation -- capacity footprints are per-preset facts this standard does not yet express")
			}
			meters := declaredMeters(profile)

			entries, err := pricebook.Entries(root, provider)
			if err != nil {
				t.Fatalf("loading the provider's price book: %v", err)
			}

			specDescriptor := kindSpecDescriptor(t, component)
			spec := derivation.GetSpec()

			if strings.TrimSpace(spec.GetCurrency()) == "" {
				t.Error("currency is empty")
			}
			if spec.GetHoursPerMonth() <= 0 {
				t.Error("hours_per_month must be declared (730 is the industry convention)")
			}

			if region := spec.GetRegion(); region != nil {
				if path := region.GetFromField(); path != "" {
					if err := specpath.Validate(specDescriptor, path); err != nil {
						t.Errorf("region from_field: %v", err)
					}
				}
			}

			for i, refusal := range spec.GetRefusals() {
				if len(refusal.GetWhen()) == 0 {
					t.Errorf("refusal %d has no conditions -- an unconditional refusal means the component should ship no derivation", i)
				}
				checkConditions(t, specDescriptor, refusal.GetWhen())
				if strings.TrimSpace(refusal.GetReason()) == "" {
					t.Errorf("refusal %d has no reason -- the honest sentence is the point", i)
				}
			}

			// A derivation with no lines is legitimate exactly when it
			// derives the honest zero: a zero-commit component whose
			// exclusions name the usage meters. Nothing at all derives
			// nothing.
			if len(spec.GetLines()) == 0 && len(spec.GetRefusals()) == 0 && len(spec.GetExclusions()) == 0 {
				t.Error("derivation declares no line rules, no refusals, and no exclusions -- what does it derive?")
			}
			for _, rule := range spec.GetLines() {
				checkLineRule(t, specDescriptor, rule, meters, entries)
			}

			for i, text := range spec.GetExclusions() {
				if strings.TrimSpace(text.GetText()) == "" {
					t.Errorf("exclusion %d has no text", i)
				}
				checkConditions(t, specDescriptor, text.GetAppliesWhen())
			}
			for i, note := range spec.GetNotes() {
				if strings.TrimSpace(note.GetText()) == "" {
					t.Errorf("note %d has no text", i)
				}
				checkConditions(t, specDescriptor, note.GetAppliesWhen())
			}
		})
	}
}

// checkLineRule verifies one rule's internal soundness.
func checkLineRule(
	t *testing.T,
	specDescriptor protoreflect.MessageDescriptor,
	rule *derivationv1.LineRule,
	meters map[string]bool,
	entries map[string]*pricebookv1.PriceBookEntry,
) {
	t.Helper()
	meter := strings.TrimSpace(rule.GetSkuMeter())
	if meter == "" {
		t.Error("line rule has no sku_meter")
	} else if !meters[meter] {
		t.Errorf("sku_meter %q is not declared by the component's cost.yaml (baseline charges and cost drivers) -- a derivation cannot price an undeclared meter", meter)
	}
	checkConditions(t, specDescriptor, rule.GetAppliesWhen())

	if len(rule.GetQuantity()) == 0 {
		t.Errorf("line %q has no quantity factors", meter)
	}
	for _, factor := range rule.GetQuantity() {
		switch f := factor.GetFactor().(type) {
		case *derivationv1.QuantityFactor_Constant:
			if !decimalPattern.MatchString(f.Constant) {
				t.Errorf("line %q: constant %q is not a plain decimal string", meter, f.Constant)
			}
		case *derivationv1.QuantityFactor_FieldValue:
			terminal, err := specpath.Terminal(specDescriptor, f.FieldValue.GetFieldPath())
			if err != nil {
				t.Errorf("line %q: field_value: %v", meter, err)
			} else if terminal.IsList() || terminal.IsMap() {
				t.Errorf("line %q: field_value %q is not a scalar -- use count_of for repeated fields", meter, f.FieldValue.GetFieldPath())
			}
			if d := f.FieldValue.GetDefaultWhenUnset(); d != "" && !decimalPattern.MatchString(d) {
				t.Errorf("line %q: default_when_unset %q is not a plain decimal string", meter, d)
			}
		case *derivationv1.QuantityFactor_CountOf:
			terminal, err := specpath.Terminal(specDescriptor, f.CountOf)
			if err != nil {
				t.Errorf("line %q: count_of: %v", meter, err)
			} else if !terminal.IsList() {
				t.Errorf("line %q: count_of %q is not a repeated field", meter, f.CountOf)
			}
		case *derivationv1.QuantityFactor_HoursInMonth:
			if !f.HoursInMonth {
				t.Errorf("line %q: hours_in_month must be true when used", meter)
			}
		default:
			t.Errorf("line %q has a quantity factor with no arm set", meter)
		}
	}

	switch price := rule.GetPrice().(type) {
	case *derivationv1.LineRule_PriceSlug:
		if _, ok := entries[price.PriceSlug]; !ok {
			t.Errorf("line %q: price_slug %q resolves to no price book entry", meter, price.PriceSlug)
		}
	case *derivationv1.LineRule_PriceLookup:
		lookup := price.PriceLookup
		if strings.TrimSpace(lookup.GetServiceName()) == "" || strings.TrimSpace(lookup.GetPricingUnit()) == "" {
			t.Errorf("line %q: price_lookup must name service_name and pricing_unit", meter)
		}
		if len(lookup.GetAttributes()) == 0 {
			t.Errorf("line %q: price_lookup binds no attributes -- with nothing selecting the price, use price_slug", meter)
		}
		for _, binding := range lookup.GetAttributes() {
			if strings.TrimSpace(binding.GetKey()) == "" {
				t.Errorf("line %q: attribute binding has no key", meter)
			}
			switch value := binding.GetValue().(type) {
			case *derivationv1.AttributeBinding_Constant:
				if strings.TrimSpace(value.Constant) == "" {
					t.Errorf("line %q: attribute %q has an empty constant", meter, binding.GetKey())
				}
			case *derivationv1.AttributeBinding_FromField:
				if err := specpath.Validate(specDescriptor, value.FromField); err != nil {
					t.Errorf("line %q: attribute %q: %v", meter, binding.GetKey(), err)
				}
			default:
				t.Errorf("line %q: attribute %q carries neither constant nor from_field", meter, binding.GetKey())
			}
		}
	default:
		t.Errorf("line %q carries neither price_slug nor price_lookup", meter)
	}

	if strings.TrimSpace(rule.GetBasis()) == "" {
		t.Errorf("line %q has no basis -- the audit trail from manifest values to the quantity", meter)
	}
}

// checkConditions verifies each condition names a resolvable path, an op,
// and the comparand exactly when the op compares.
func checkConditions(t *testing.T, specDescriptor protoreflect.MessageDescriptor, conditions []*derivationv1.Condition) {
	t.Helper()
	for _, condition := range conditions {
		if err := specpath.Validate(specDescriptor, condition.GetFieldPath()); err != nil {
			t.Errorf("condition: %v", err)
		}
		switch condition.GetOp() {
		case derivationv1.Condition_equals, derivationv1.Condition_not_equals:
			if condition.GetValue() == "" {
				t.Errorf("condition on %q compares against an empty value -- use is_unset for absence", condition.GetFieldPath())
			}
		case derivationv1.Condition_is_set, derivationv1.Condition_is_unset:
			if condition.GetValue() != "" {
				t.Errorf("condition on %q is a presence check but carries value %q", condition.GetFieldPath(), condition.GetValue())
			}
		default:
			t.Errorf("condition on %q has no op", condition.GetFieldPath())
		}
	}
}

// declaredMeters collects the sku_meter vocabulary a component's cost
// profile declares across baseline charges and cost drivers.
func declaredMeters(profile *costprofilev1.ComponentCostProfile) map[string]bool {
	meters := map[string]bool{}
	for _, charge := range profile.GetSpec().GetBaselineCharges() {
		meters[strings.TrimSpace(charge.GetSkuMeter())] = true
	}
	for _, driver := range profile.GetSpec().GetCostDrivers() {
		meters[strings.TrimSpace(driver.GetSkuMeter())] = true
	}
	return meters
}

// componentProvider locates the provider directory a component lives under.
func componentProvider(t *testing.T, repoRoot, component string) string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(repoRoot, "catalog", "*", component))
	if err != nil || len(matches) == 0 {
		t.Fatalf("cost derivation names component %q, which exists nowhere under catalog/", component)
	}
	if len(matches) > 1 {
		t.Fatalf("component %q exists under multiple providers: %v", component, matches)
	}
	return filepath.Base(filepath.Dir(matches[0]))
}

// kindSpecDescriptor resolves a component directory name to its kind's spec
// message descriptor via the kind registry.
func kindSpecDescriptor(t *testing.T, component string) protoreflect.MessageDescriptor {
	t.Helper()
	kind := crkreflect.KindFromString(component)
	apiMessage, err := crkreflect.NewInstance(kind)
	if err != nil {
		t.Fatalf("NewInstance(%s): %v", component, err)
	}
	specField := apiMessage.ProtoReflect().Descriptor().Fields().ByName("spec")
	if specField == nil || specField.Kind() != protoreflect.MessageKind {
		t.Fatalf("%s has no spec message field", apiMessage.ProtoReflect().Descriptor().FullName())
	}
	return specField.Message()
}

// repoRoot walks up from this test file to the directory containing go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found walking up from test file")
		}
		dir = parent
	}
}
