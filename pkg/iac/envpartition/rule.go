package envpartition

import (
	_ "embed"
	"strings"

	"github.com/pkg/errors"
	rulev1 "github.com/plantonhq/planton/apis/dev/planton/iac/environmentpartitionrule/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// defaultRuleYAML is the untaught state, shipped as a document of the rule's
// own kind so the taught rule is just an edited copy -- untaught and taught
// differ only in data.
//
//go:embed default-rule.yaml
var defaultRuleYAML []byte

// Rule is a compiled EnvironmentPartitionRule, validated and normalized for
// the engine. Compile once, apply to any number of resources.
type Rule struct {
	// authoritativeTagKeys in declared order (declared order is the
	// deterministic scan order for a resource's tags; the keys are
	// expected to agree, and disagreement is flagged, never resolved by
	// order).
	authoritativeTagKeys []string
	// environments in declared order -- canonical names.
	environments []string
	// tokenToEnv maps each lowercased declared token to its canonical
	// environment name.
	tokenToEnv map[string]string
	// overrides keyed by scan coordinates.
	overrides map[resourceKey]string
	// containmentInheritance is tier 4's switch (on unless disabled).
	containmentInheritance bool
	// fallbackEnvironment, "" when the rule declares none.
	fallbackEnvironment string
}

type resourceKey struct {
	typeName   string
	identifier string
}

// DefaultRule returns the compiled untaught rule. The embedded document is
// part of the build; failing to compile it is a programming error, so this
// panics rather than making every caller handle an impossible error.
func DefaultRule() *Rule {
	rule, err := ParseRuleYAML(defaultRuleYAML)
	if err != nil {
		panic("envpartition: embedded default rule is invalid: " + err.Error())
	}
	return rule
}

// ParseRuleYAML parses and compiles an EnvironmentPartitionRule document.
func ParseRuleYAML(yamlBytes []byte) (*Rule, error) {
	doc := &rulev1.EnvironmentPartitionRule{}
	if err := protobufyaml.LoadYamlBytes(yamlBytes, doc); err != nil {
		return nil, errors.Wrap(err, "parsing EnvironmentPartitionRule yaml")
	}
	return CompileRule(doc)
}

// CompileRule validates and compiles a rule document.
func CompileRule(doc *rulev1.EnvironmentPartitionRule) (*Rule, error) {
	if doc.GetKind() != "EnvironmentPartitionRule" {
		return nil, errors.Errorf("kind is %q, want EnvironmentPartitionRule", doc.GetKind())
	}
	spec := doc.GetSpec()
	rule := &Rule{
		tokenToEnv:             map[string]string{},
		overrides:              map[resourceKey]string{},
		containmentInheritance: !spec.GetDisableContainmentInheritance(),
		fallbackEnvironment:    spec.GetFallbackEnvironment(),
	}
	for _, key := range spec.GetAuthoritativeTagKeys() {
		if key == "" {
			return nil, errors.New("authoritative_tag_keys contains an empty key")
		}
		rule.authoritativeTagKeys = append(rule.authoritativeTagKeys, key)
	}
	for i, env := range spec.GetEnvironments() {
		if env.GetName() == "" {
			return nil, errors.Errorf("environments[%d]: name is empty", i)
		}
		rule.environments = append(rule.environments, env.GetName())
		for _, token := range env.GetTokens() {
			normalized := strings.ToLower(token)
			if normalized == "" {
				return nil, errors.Errorf("environment %q declares an empty token", env.GetName())
			}
			if owner, taken := rule.tokenToEnv[normalized]; taken && owner != env.GetName() {
				return nil, errors.Errorf("token %q is declared by both %q and %q -- a token must identify exactly one environment", normalized, owner, env.GetName())
			}
			rule.tokenToEnv[normalized] = env.GetName()
		}
	}
	for i, override := range spec.GetOverrides() {
		if override.GetIdentifier() == "" || override.GetEnvironment() == "" {
			return nil, errors.Errorf("overrides[%d]: identifier and environment are both required", i)
		}
		key := resourceKey{typeName: override.GetTypeName(), identifier: override.GetIdentifier()}
		if existing, dup := rule.overrides[key]; dup && existing != override.GetEnvironment() {
			return nil, errors.Errorf("overrides[%d]: %s %q is pinned to both %q and %q", i, key.typeName, key.identifier, existing, override.GetEnvironment())
		}
		rule.overrides[key] = override.GetEnvironment()
	}
	return rule, nil
}

// normalizeEnvValue maps an authoritative tag's value to a canonical
// environment name: through the declared vocabulary when the value matches
// a token ("prd" -> prod), literally otherwise (an explicit tag is the
// user's word; "ops" names an "ops" environment even if the rule never
// declared it).
func (r *Rule) normalizeEnvValue(value string) string {
	if canonical, declared := r.tokenToEnv[strings.ToLower(value)]; declared {
		return canonical
	}
	return value
}
