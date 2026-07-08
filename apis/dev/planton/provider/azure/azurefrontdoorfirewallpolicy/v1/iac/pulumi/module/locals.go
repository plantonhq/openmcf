package module

import (
	"strings"

	azurefrontdoorfirewallpolicyv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurefrontdoorfirewallpolicy/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureFrontDoorFirewallPolicy *azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicy
	ResourceGroupName            string
	AzureTags                    map[string]string
	// SkuName is ARM's tier value, materialized from the spec enum with
	// the documented STANDARD default (stack inputs never carry proto
	// defaults).
	SkuName string
	// IsPremium gates the Premium-only policy settings (the
	// JS-challenge/CAPTCHA expirations, which Azure ALWAYS enables on
	// Premium and REJECTS on Standard).
	IsPremium bool
}

// skuStrings maps the spec's sku enum to ARM's tier values.
var skuStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySku]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySku_STANDARD: "Standard_AzureFrontDoor",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySku_PREMIUM:  "Premium_AzureFrontDoor",
}

// modeStrings maps the enforcement-mode enum to ARM's values.
var modeStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMode]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMode_DETECTION:  "Detection",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMode_PREVENTION: "Prevention",
}

// customRuleTypeStrings maps the custom-rule type enum to ARM's values.
var customRuleTypeStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleType]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleType_MATCH_RULE:      "MatchRule",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleType_RATE_LIMIT_RULE: "RateLimitRule",
}

// customRuleActionStrings maps the custom-rule action enum to ARM's
// values (JSChallenge/CAPTCHA are Premium-only, spec-enforced).
var customRuleActionStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleAction]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleAction_ALLOW:        "Allow",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleAction_BLOCK:        "Block",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleAction_LOG:          "Log",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleAction_REDIRECT:     "Redirect",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleAction_JS_CHALLENGE: "JSChallenge",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyCustomRuleAction_CAPTCHA:      "CAPTCHA",
}

// matchVariableStrings maps the match-variable enum to ARM's values.
var matchVariableStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_COOKIES:        "Cookies",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_POST_ARGS:      "PostArgs",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_QUERY_STRING:   "QueryString",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_REMOTE_ADDR:    "RemoteAddr",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_REQUEST_BODY:   "RequestBody",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_REQUEST_HEADER: "RequestHeader",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_REQUEST_METHOD: "RequestMethod",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_REQUEST_URI:    "RequestUri",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyMatchVariable_SOCKET_ADDR:    "SocketAddr",
}

// operatorStrings maps the condition-operator enum to ARM's exact
// (case-sensitive) operator names.
var operatorStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_ANY:                   "Any",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_BEGINS_WITH:           "BeginsWith",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_CONTAINS:              "Contains",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_ENDS_WITH:             "EndsWith",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_EQUAL:                 "Equal",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_GEO_MATCH:             "GeoMatch",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_GREATER_THAN:          "GreaterThan",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_GREATER_THAN_OR_EQUAL: "GreaterThanOrEqual",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_IP_MATCH:              "IPMatch",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_LESS_THAN:             "LessThan",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_LESS_THAN_OR_EQUAL:    "LessThanOrEqual",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyOperator_REG_EX:                "RegEx",
}

// transformStrings maps the transform enum to ARM's values.
var transformStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyTransform]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyTransform_LOWERCASE:    "Lowercase",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyTransform_REMOVE_NULLS: "RemoveNulls",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyTransform_TRIM:         "Trim",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyTransform_UPPERCASE:    "Uppercase",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyTransform_URL_DECODE:   "URLDecode",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyTransform_URL_ENCODE:   "URLEncode",
}

// managedRuleSetActionStrings maps the rule-set action enum to ARM's
// values (the RULE_SET_ prefix exists only to keep the proto enum
// names collision-free within the kind).
var managedRuleSetActionStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleSetAction]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleSetAction_RULE_SET_BLOCK:    "Block",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleSetAction_RULE_SET_LOG:      "Log",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleSetAction_RULE_SET_REDIRECT: "Redirect",
}

// exclusionMatchVariableStrings maps the exclusion-collection enum to
// ARM's values.
var exclusionMatchVariableStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyExclusionMatchVariable]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyExclusionMatchVariable_EXCLUDE_QUERY_STRING_ARG_NAMES:      "QueryStringArgNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyExclusionMatchVariable_EXCLUDE_REQUEST_BODY_JSON_ARG_NAMES: "RequestBodyJsonArgNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyExclusionMatchVariable_EXCLUDE_REQUEST_BODY_POST_ARG_NAMES: "RequestBodyPostArgNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyExclusionMatchVariable_EXCLUDE_REQUEST_COOKIE_NAMES:        "RequestCookieNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyExclusionMatchVariable_EXCLUDE_REQUEST_HEADER_NAMES:        "RequestHeaderNames",
}

// selectorOperatorStrings maps the selector-operator enum to ARM's
// values (shared by managed-rule exclusions and log scrubbing).
var selectorOperatorStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySelectorOperator]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySelectorOperator_SELECTOR_CONTAINS:    "Contains",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySelectorOperator_SELECTOR_ENDS_WITH:   "EndsWith",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySelectorOperator_SELECTOR_EQUALS:      "Equals",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySelectorOperator_SELECTOR_EQUALS_ANY:  "EqualsAny",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicySelectorOperator_SELECTOR_STARTS_WITH: "StartsWith",
}

// managedRuleOverrideActionStrings maps the per-rule override action
// enum to ARM's values (the OVERRIDE_ prefix exists only to keep the
// proto enum names collision-free within the kind).
var managedRuleOverrideActionStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleOverrideAction]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleOverrideAction_OVERRIDE_ALLOW:           "Allow",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleOverrideAction_OVERRIDE_ANOMALY_SCORING: "AnomalyScoring",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleOverrideAction_OVERRIDE_BLOCK:           "Block",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleOverrideAction_OVERRIDE_CAPTCHA:         "CAPTCHA",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleOverrideAction_OVERRIDE_JS_CHALLENGE:    "JSChallenge",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleOverrideAction_OVERRIDE_REDIRECT:        "Redirect",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyManagedRuleOverrideAction_OVERRIDE_LOG:             "Log",
}

// scrubbingMatchVariableStrings maps the scrubbing-variable enum to
// ARM's values (a superset of the profile's access-log vocabulary --
// the WAF logs carry more request parts).
var scrubbingMatchVariableStrings = map[azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyScrubbingMatchVariable]string{
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyScrubbingMatchVariable_SCRUB_QUERY_STRING_ARG_NAMES:      "QueryStringArgNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_BODY_JSON_ARG_NAMES: "RequestBodyJsonArgNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_BODY_POST_ARG_NAMES: "RequestBodyPostArgNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_COOKIE_NAMES:        "RequestCookieNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_HEADER_NAMES:        "RequestHeaderNames",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_IP_ADDRESS:          "RequestIPAddress",
	azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_URI:                 "RequestUri",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurefrontdoorfirewallpolicyv1.AzureFrontDoorFirewallPolicyStackInput) *Locals {
	locals := &Locals{}

	locals.AzureFrontDoorFirewallPolicy = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// Materialize the tier default: unspecified deploys STANDARD (the
	// spec's documented default -- stack inputs never carry proto
	// defaults).
	locals.SkuName = skuStrings[target.Spec.Sku]
	if locals.SkuName == "" {
		locals.SkuName = "Standard_AzureFrontDoor"
	}
	locals.IsPremium = locals.SkuName == "Premium_AzureFrontDoor"

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		// PARITY-EXCEPTION: resource_kind here is the lowered
		// CloudResourceKind enum string and resource_id is omitted when
		// metadata.id is empty, while the Terraform module emits the
		// family-wide snake-case literal and falls back to metadata.name.
		// Output-neutral (tags never feed stack outputs); aligning the two
		// shapes is a family-wide convention change, not a per-kind fix.
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureFrontDoorFirewallPolicy.String()),
	}

	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}

	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
