// Package costestimator evaluates a component's cost derivation against
// one typed manifest. It is the execution engine of the derivation
// standard: conditions select which rules apply to THIS configuration,
// quantity factors multiply out the monthly consumption, and price
// selection resolves each line to exactly one pinned price-book entry --
// or the whole estimate refuses with the reason it cannot know the
// number. It never guesses: an indeterminate rate, an unpriced value, or
// an ambiguous lookup is a refusal, and a partial estimate (some lines
// priced, one silently missing) does not exist as an outcome.
//
// The result is the same shape an authored estimate model preset carries
// (quantity lines referencing book entries by slug), so the estimate
// generator's join, arithmetic, coherence checks, and emission run
// unchanged on evaluated output -- one home for pricing arithmetic,
// whether the quantities were authored by hand or derived from a
// manifest.
package costestimator

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	derivationv1 "github.com/plantonhq/planton/finops/componentcostderivation/v1"
	estimatemodelv1 "github.com/plantonhq/planton/finops/componentcostestimatemodel/v1"
	pricebookv1 "github.com/plantonhq/planton/finops/pricebook/v1"
	"github.com/plantonhq/planton/pkg/finops/pricebook"
	"github.com/plantonhq/planton/pkg/specpath"
)

// Refusal is the honest non-answer: this configuration's committed cost
// cannot be known, and here is why. A refusal is a first-class outcome,
// never an error -- errors mean the derivation itself is malformed.
type Refusal struct {
	Reason string
}

// Evaluate runs a derivation against one manifest and returns exactly one
// of: the evaluated preset-model shape (quantity lines with resolved
// price slugs, configuration-scoped exclusions and notes), a refusal, or
// an error for a malformed derivation (unresolvable path, bad decimal,
// unknown slug -- the conformance gate catches these before CI ever runs
// the engine, so errors here mean the gate was bypassed).
//
// entries is the component's provider price book keyed by slug (the
// pricebook.Entries shape).
func Evaluate(
	manifest proto.Message,
	spec *derivationv1.ComponentCostDerivationSpec,
	entries map[string]*pricebookv1.PriceBookEntry,
) (*estimatemodelv1.PresetEstimateModel, *Refusal, error) {
	specMsg, err := manifestSpec(manifest)
	if err != nil {
		return nil, nil, err
	}

	for _, refusal := range spec.GetRefusals() {
		holds, err := conditionsHold(specMsg, refusal.GetWhen())
		if err != nil {
			return nil, nil, err
		}
		if holds {
			return nil, &Refusal{Reason: refusal.GetReason()}, nil
		}
	}

	region, err := resolveRegion(specMsg, spec.GetRegion())
	if err != nil {
		return nil, nil, err
	}

	lines, refusal, err := evaluateLines(specMsg, spec, entries, region)
	if err != nil || refusal != nil {
		return nil, refusal, err
	}

	exclusions, err := matchingTexts(specMsg, spec.GetExclusions())
	if err != nil {
		return nil, nil, err
	}
	notes, err := matchingTexts(specMsg, spec.GetNotes())
	if err != nil {
		return nil, nil, err
	}

	return &estimatemodelv1.PresetEstimateModel{
		RegionAssumption: region,
		Currency:         spec.GetCurrency(),
		HoursPerMonth:    spec.GetHoursPerMonth(),
		QuantityLines:    lines,
		Exclusions:       exclusions,
		Notes:            strings.Join(notes, " "),
	}, nil, nil
}

// evaluateLines applies every line rule, merges same-meter-same-price
// lines, and drops zero quantities. One unpriceable applicable line
// refuses the WHOLE estimate: a total missing a line it owes is a lie.
func evaluateLines(
	specMsg protoreflect.Message,
	spec *derivationv1.ComponentCostDerivationSpec,
	entries map[string]*pricebookv1.PriceBookEntry,
	region string,
) ([]*estimatemodelv1.QuantityLine, *Refusal, error) {
	type merged struct {
		line     *estimatemodelv1.QuantityLine
		quantity *big.Rat
		bases    []string
		order    int
	}
	mergedByKey := map[string]*merged{}

	for _, rule := range spec.GetLines() {
		holds, err := conditionsHold(specMsg, rule.GetAppliesWhen())
		if err != nil {
			return nil, nil, err
		}
		if !holds {
			continue
		}

		quantity, err := evaluateQuantity(specMsg, spec, rule)
		if err != nil {
			return nil, nil, err
		}
		if quantity.Sign() == 0 {
			continue
		}

		slug, refusal, err := resolvePrice(specMsg, rule, entries, region, spec.GetCurrency())
		if err != nil || refusal != nil {
			return nil, refusal, err
		}

		key := rule.GetSkuMeter() + "\x00" + slug
		if existing, ok := mergedByKey[key]; ok {
			existing.quantity.Add(existing.quantity, quantity)
			existing.bases = append(existing.bases, rule.GetBasis())
			continue
		}
		mergedByKey[key] = &merged{
			line: &estimatemodelv1.QuantityLine{
				SkuMeter: rule.GetSkuMeter(),
				Price:    slug,
			},
			quantity: quantity,
			bases:    []string{rule.GetBasis()},
			order:    len(mergedByKey),
		}
	}

	ordered := make([]*merged, 0, len(mergedByKey))
	for _, m := range mergedByKey {
		ordered = append(ordered, m)
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].order < ordered[j].order })

	var lines []*estimatemodelv1.QuantityLine
	for _, m := range ordered {
		m.line.PricingQuantity = plainDecimal(m.quantity)
		m.line.QuantityBasis = strings.Join(m.bases, "; ")
		lines = append(lines, m.line)
	}
	return lines, nil, nil
}

// resolvePrice returns the slug of the one book entry pricing a rule --
// static, or selected by the manifest's own values.
func resolvePrice(
	specMsg protoreflect.Message,
	rule *derivationv1.LineRule,
	entries map[string]*pricebookv1.PriceBookEntry,
	region string,
	currency string,
) (string, *Refusal, error) {
	if slug := rule.GetPriceSlug(); slug != "" {
		if _, ok := entries[slug]; !ok {
			return "", nil, fmt.Errorf("line %q: price_slug %q resolves to no price book entry", rule.GetSkuMeter(), slug)
		}
		return slug, nil, nil
	}

	lookup := rule.GetPriceLookup()
	if lookup == nil {
		return "", nil, fmt.Errorf("line %q carries neither price_slug nor price_lookup", rule.GetSkuMeter())
	}

	want := map[string]string{}
	for _, binding := range lookup.GetAttributes() {
		value, refusal, err := bindingValue(specMsg, rule, binding)
		if err != nil || refusal != nil {
			return "", refusal, err
		}
		want[binding.GetKey()] = value
	}

	var matches []string
	for slug, entry := range entries {
		if len(entry.GetAttributes()) == 0 {
			continue
		}
		if entry.GetServiceName() != lookup.GetServiceName() ||
			entry.GetPricingUnit() != lookup.GetPricingUnit() ||
			entry.GetCurrency() != currency {
			continue
		}
		if entry.GetRegion() != pricebook.GlobalRegion && entry.GetRegion() != region {
			continue
		}
		if !attributesEqual(entry.GetAttributes(), want) {
			continue
		}
		matches = append(matches, slug)
	}
	sort.Strings(matches)

	switch len(matches) {
	case 1:
		return matches[0], nil, nil
	case 0:
		return "", &Refusal{Reason: fmt.Sprintf(
			"no pinned price for %s (%s, %s in %s, %s) -- the price book does not yet cover this configuration, and a number without a verified price would be a guess",
			rule.GetSkuMeter(), lookup.GetServiceName(), describeAttributes(want), region, currency)}, nil
	default:
		return "", &Refusal{Reason: fmt.Sprintf(
			"several pinned prices match %s (%s) -- the lookup is ambiguous and picking one would be a guess",
			rule.GetSkuMeter(), describeAttributes(want))}, nil
	}
}

// bindingValue resolves one attribute binding to its string value.
func bindingValue(
	specMsg protoreflect.Message,
	rule *derivationv1.LineRule,
	binding *derivationv1.AttributeBinding,
) (string, *Refusal, error) {
	if constant := binding.GetConstant(); constant != "" {
		return constant, nil, nil
	}
	path := binding.GetFromField()
	if path == "" {
		return "", nil, fmt.Errorf("line %q: attribute %q carries neither constant nor from_field", rule.GetSkuMeter(), binding.GetKey())
	}
	resolved, err := specpath.Resolve(specMsg, path)
	if err != nil {
		return "", nil, fmt.Errorf("line %q: attribute %q: %w", rule.GetSkuMeter(), binding.GetKey(), err)
	}
	if resolved.Field.IsList() {
		count := 0
		if resolved.Value.IsValid() {
			count = resolved.Value.List().Len()
		}
		if count > 1 {
			return "", &Refusal{Reason: fmt.Sprintf(
				"%q carries %d values -- the provider picks among them at its own discretion, so the committed rate is indeterminate",
				path, count)}, nil
		}
		value := ""
		if count == 1 {
			var err error
			value, err = renderScalar(resolved.Field, resolved.Value.List().Get(0))
			if err != nil {
				return "", nil, fmt.Errorf("line %q: attribute %q: %w", rule.GetSkuMeter(), binding.GetKey(), err)
			}
		}
		return bindingOrAssumption(path, value, binding)
	}
	value, err := effectiveString(resolved)
	if err != nil {
		return "", nil, fmt.Errorf("line %q: attribute %q: %w", rule.GetSkuMeter(), binding.GetKey(), err)
	}
	return bindingOrAssumption(path, value, binding)
}

// bindingOrAssumption applies the binding's declared assumption to an
// unset value -- a stated assumption (the basis prose carries it) is
// honest; an unset value with no assumption is a refusal, never a guess.
func bindingOrAssumption(path, value string, binding *derivationv1.AttributeBinding) (string, *Refusal, error) {
	if value != "" {
		return value, nil, nil
	}
	if assumption := binding.GetAssumptionWhenUnset(); assumption != "" {
		return assumption, nil, nil
	}
	return "", &Refusal{Reason: fmt.Sprintf("%q is unset -- nothing selects the price", path)}, nil
}

// evaluateQuantity multiplies out a rule's factors.
func evaluateQuantity(
	specMsg protoreflect.Message,
	spec *derivationv1.ComponentCostDerivationSpec,
	rule *derivationv1.LineRule,
) (*big.Rat, error) {
	factors := rule.GetQuantity()
	if len(factors) == 0 {
		return nil, fmt.Errorf("line %q has no quantity factors", rule.GetSkuMeter())
	}
	product := big.NewRat(1, 1)
	for _, factor := range factors {
		value, err := factorValue(specMsg, spec, rule, factor)
		if err != nil {
			return nil, err
		}
		product.Mul(product, value)
	}
	return product, nil
}

// factorValue evaluates one quantity factor.
func factorValue(
	specMsg protoreflect.Message,
	spec *derivationv1.ComponentCostDerivationSpec,
	rule *derivationv1.LineRule,
	factor *derivationv1.QuantityFactor,
) (*big.Rat, error) {
	switch f := factor.GetFactor().(type) {
	case *derivationv1.QuantityFactor_Constant:
		rat, ok := parseDecimal(f.Constant)
		if !ok {
			return nil, fmt.Errorf("line %q: constant %q is not a plain decimal string", rule.GetSkuMeter(), f.Constant)
		}
		return rat, nil

	case *derivationv1.QuantityFactor_FieldValue:
		resolved, err := specpath.Resolve(specMsg, f.FieldValue.GetFieldPath())
		if err != nil {
			return nil, fmt.Errorf("line %q: %w", rule.GetSkuMeter(), err)
		}
		if resolved.Field.IsList() || resolved.Field.IsMap() {
			return nil, fmt.Errorf("line %q: field_value %q is not a scalar -- use count_of for repeated fields",
				rule.GetSkuMeter(), f.FieldValue.GetFieldPath())
		}
		text, err := effectiveString(resolved)
		if err != nil {
			return nil, fmt.Errorf("line %q: %w", rule.GetSkuMeter(), err)
		}
		rat, ok := parseDecimal(text)
		if !ok && text != "" {
			return nil, fmt.Errorf("line %q: field %q reads %q, not a number", rule.GetSkuMeter(), f.FieldValue.GetFieldPath(), text)
		}
		if (text == "" || rat.Sign() == 0) && f.FieldValue.GetDefaultWhenUnset() != "" {
			rat, ok = parseDecimal(f.FieldValue.GetDefaultWhenUnset())
			if !ok {
				return nil, fmt.Errorf("line %q: default_when_unset %q is not a plain decimal string",
					rule.GetSkuMeter(), f.FieldValue.GetDefaultWhenUnset())
			}
			return rat, nil
		}
		if text == "" {
			return big.NewRat(0, 1), nil
		}
		return rat, nil

	case *derivationv1.QuantityFactor_CountOf:
		resolved, err := specpath.Resolve(specMsg, f.CountOf)
		if err != nil {
			return nil, fmt.Errorf("line %q: %w", rule.GetSkuMeter(), err)
		}
		if !resolved.Field.IsList() {
			return nil, fmt.Errorf("line %q: count_of %q is not a repeated field", rule.GetSkuMeter(), f.CountOf)
		}
		count := 0
		if resolved.Value.IsValid() {
			count = resolved.Value.List().Len()
		}
		return big.NewRat(int64(count), 1), nil

	case *derivationv1.QuantityFactor_HoursInMonth:
		if !f.HoursInMonth {
			return nil, fmt.Errorf("line %q: hours_in_month must be true when used", rule.GetSkuMeter())
		}
		if spec.GetHoursPerMonth() == 0 {
			return nil, fmt.Errorf("line %q uses hours_in_month but the spec declares no hours_per_month", rule.GetSkuMeter())
		}
		return big.NewRat(int64(spec.GetHoursPerMonth()), 1), nil

	default:
		return nil, fmt.Errorf("line %q has a quantity factor with no arm set", rule.GetSkuMeter())
	}
}

// conditionsHold evaluates a condition set; ALL must hold, and an empty
// set holds vacuously.
func conditionsHold(specMsg protoreflect.Message, conditions []*derivationv1.Condition) (bool, error) {
	for _, condition := range conditions {
		holds, err := conditionHolds(specMsg, condition)
		if err != nil {
			return false, err
		}
		if !holds {
			return false, nil
		}
	}
	return true, nil
}

// conditionHolds evaluates one condition. Comparisons read the field's
// EFFECTIVE value (an unset scalar reads as its zero value); presence
// checks read explicit presence. Catalog placeholders ("<aws-region>")
// read as unset either way.
func conditionHolds(specMsg protoreflect.Message, condition *derivationv1.Condition) (bool, error) {
	resolved, err := specpath.Resolve(specMsg, condition.GetFieldPath())
	if err != nil {
		return false, fmt.Errorf("condition on %q: %w", condition.GetFieldPath(), err)
	}

	switch condition.GetOp() {
	case derivationv1.Condition_is_set, derivationv1.Condition_is_unset:
		present := resolved.Present
		if resolved.Field.IsList() {
			present = resolved.Value.IsValid() && resolved.Value.List().Len() > 0
		} else if present && !resolved.Field.IsMap() && resolved.Field.Kind() != protoreflect.MessageKind {
			text, err := renderScalar(resolved.Field, resolved.Value)
			if err == nil && text == "" {
				present = false
			}
		} else if present && resolved.Field.Kind() == protoreflect.MessageKind && !resolved.Field.IsMap() {
			if text, unwrapped := unwrapValue(resolved.Field, resolved.Value); unwrapped && text == "" {
				present = false
			}
		}
		if condition.GetOp() == derivationv1.Condition_is_set {
			return present, nil
		}
		return !present, nil

	case derivationv1.Condition_equals, derivationv1.Condition_not_equals:
		text, err := effectiveString(resolved)
		if err != nil {
			return false, fmt.Errorf("condition on %q: %w", condition.GetFieldPath(), err)
		}
		if condition.GetOp() == derivationv1.Condition_equals {
			return text == condition.GetValue(), nil
		}
		return text != condition.GetValue(), nil

	default:
		return false, fmt.Errorf("condition on %q has no op", condition.GetFieldPath())
	}
}

// resolveRegion resolves the pricing region: the manifest's own field
// when it answers, the declared assumption otherwise.
func resolveRegion(specMsg protoreflect.Message, binding *derivationv1.RegionBinding) (string, error) {
	if binding == nil {
		return "", nil
	}
	if path := binding.GetFromField(); path != "" {
		resolved, err := specpath.Resolve(specMsg, path)
		if err != nil {
			return "", fmt.Errorf("region from_field: %w", err)
		}
		text, err := effectiveString(resolved)
		if err != nil {
			return "", fmt.Errorf("region from_field: %w", err)
		}
		if text != "" {
			if binding.GetZoneToRegion() {
				text = zoneToRegion(text)
			}
			return text, nil
		}
	}
	return binding.GetAssumption(), nil
}

// zoneToRegion resolves a GCP-style zonal location to its region: a
// single-letter suffix after the last hyphen is the zone letter
// ("us-central1-a" -> "us-central1"); anything else passes through.
func zoneToRegion(location string) string {
	if i := strings.LastIndex(location, "-"); i > 0 && len(location)-i == 2 {
		suffix := location[i+1]
		if suffix >= 'a' && suffix <= 'z' {
			return location[:i]
		}
	}
	return location
}

// matchingTexts returns the texts whose conditions hold, in authored order.
func matchingTexts(specMsg protoreflect.Message, texts []*derivationv1.ConditionalText) ([]string, error) {
	var matched []string
	for _, text := range texts {
		holds, err := conditionsHold(specMsg, text.GetAppliesWhen())
		if err != nil {
			return nil, err
		}
		if holds {
			matched = append(matched, text.GetText())
		}
	}
	return matched, nil
}

// manifestSpec returns the manifest's spec message.
func manifestSpec(manifest proto.Message) (protoreflect.Message, error) {
	reflected := manifest.ProtoReflect()
	specField := reflected.Descriptor().Fields().ByName("spec")
	if specField == nil || specField.Kind() != protoreflect.MessageKind {
		return nil, fmt.Errorf("%s has no spec message field", reflected.Descriptor().FullName())
	}
	return reflected.Get(specField).Message(), nil
}

// effectiveString renders a resolved field's effective value as text.
// Wrapper messages (a message whose scalar `value` field carries the
// literal beside a reference arm) render their literal; a placeholder or
// genuinely empty value renders "".
func effectiveString(resolved specpath.Resolved) (string, error) {
	if resolved.Field.IsMap() {
		return "", fmt.Errorf("field %q is a map -- maps have no effective value", resolved.Field.Name())
	}
	if resolved.Field.Kind() == protoreflect.MessageKind {
		if !resolved.Present || !resolved.Value.IsValid() {
			return "", nil
		}
		if text, unwrapped := unwrapValue(resolved.Field, resolved.Value); unwrapped {
			return text, nil
		}
		return "", fmt.Errorf("field %q is a message with no scalar value field -- it has no effective value", resolved.Field.Name())
	}
	value := resolved.Value
	if !value.IsValid() {
		value = resolved.Field.Default()
	}
	return renderScalar(resolved.Field, value)
}

// unwrapValue reads a wrapper message's scalar `value` field (the
// value-or-reference shape catalog specs use). Returns unwrapped=false
// when the message carries no scalar value field.
func unwrapValue(field protoreflect.FieldDescriptor, value protoreflect.Value) (string, bool) {
	inner := field.Message().Fields().ByName("value")
	if inner == nil || inner.Kind() == protoreflect.MessageKind || inner.IsList() || inner.IsMap() {
		return "", false
	}
	text, err := renderScalar(inner, value.Message().Get(inner))
	if err != nil {
		return "", false
	}
	return text, true
}

// renderScalar renders one scalar value the way conditions and attribute
// bindings compare it: enums by name, bools as "true"/"false", numbers
// as plain decimals, strings as themselves -- with catalog placeholders
// ("<aws-region>") reading as unset.
func renderScalar(field protoreflect.FieldDescriptor, value protoreflect.Value) (string, error) {
	switch field.Kind() {
	case protoreflect.StringKind:
		text := value.String()
		if isPlaceholder(text) {
			return "", nil
		}
		return text, nil
	case protoreflect.BoolKind:
		return strconv.FormatBool(value.Bool()), nil
	case protoreflect.EnumKind:
		enumValue := field.Enum().Values().ByNumber(value.Enum())
		if enumValue == nil {
			return "", fmt.Errorf("field %q carries unknown enum number %d", field.Name(), value.Enum())
		}
		return string(enumValue.Name()), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return strconv.FormatInt(value.Int(), 10), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return strconv.FormatUint(value.Uint(), 10), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("field %q has kind %s -- no effective value rendering", field.Name(), field.Kind())
	}
}

// isPlaceholder recognizes the catalog's template placeholders
// ("<aws-region>", "<node-group-name>") so they never masquerade as real
// configuration values.
func isPlaceholder(text string) bool {
	return len(text) > 2 && strings.HasPrefix(text, "<") && strings.HasSuffix(text, ">")
}

// describeAttributes renders bound attribute values for refusal
// sentences, deterministically ordered.
func describeAttributes(attributes map[string]string) string {
	keys := make([]string, 0, len(attributes))
	for key := range attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+" "+attributes[key])
	}
	return strings.Join(parts, ", ")
}

// attributesEqual compares an entry's attribute identity against the
// bound values -- exact map equality, so an entry carrying extra
// attributes never matches a narrower lookup by accident.
func attributesEqual(entry map[string]string, want map[string]string) bool {
	if len(entry) != len(want) {
		return false
	}
	for key, value := range want {
		if entry[key] != value {
			return false
		}
	}
	return true
}

// parseDecimal parses a plain decimal string into an exact rational.
func parseDecimal(text string) (*big.Rat, bool) {
	if text == "" {
		return nil, false
	}
	for _, r := range text {
		if (r < '0' || r > '9') && r != '.' {
			return nil, false
		}
	}
	rat, ok := new(big.Rat).SetString(text)
	return rat, ok
}

// plainDecimal renders an exact rational as the minimal plain decimal
// string ("730", "1460", "0.5") -- quantities read as the numbers a
// person would write, never with manufactured trailing zeros.
func plainDecimal(r *big.Rat) string {
	if r.IsInt() {
		return r.Num().String()
	}
	// Products and sums of finite decimals are finite decimals: the
	// denominator is a power of 2s and 5s, so FloatString at that scale
	// never rounds; trailing zeros are then trimmed.
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
	text := r.FloatString(scale)
	text = strings.TrimRight(text, "0")
	return strings.TrimSuffix(text, ".")
}
