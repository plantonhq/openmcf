package module

import (
	"strings"

	azurewebapplicationfirewallpolicyv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurewebapplicationfirewallpolicy/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureWebApplicationFirewallPolicy *azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicy
	ResourceGroupName                 string
	AzureTags                         map[string]string
}

// The spec's enums arrive as proto enum values; each map below carries the
// complete vocabulary for its enum, translated to ARM's strings. A missing
// entry would silently drop the setting, so the maps are exhaustive by
// construction.

var ruleTypeStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyCustomRuleType]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyCustomRuleType_MATCH_RULE:      "MatchRule",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyCustomRuleType_RATE_LIMIT_RULE: "RateLimitRule",
}

var customRuleActionStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyCustomRuleAction]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyCustomRuleAction_ALLOW:        "Allow",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyCustomRuleAction_BLOCK:        "Block",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyCustomRuleAction_LOG:          "Log",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyCustomRuleAction_JS_CHALLENGE: "JSChallenge",
}

var rateLimitDurationStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRateLimitDuration]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRateLimitDuration_ONE_MIN:   "OneMin",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRateLimitDuration_FIVE_MINS: "FiveMins",
}

var groupRateLimitByStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyGroupRateLimitBy]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyGroupRateLimitBy_CLIENT_ADDR:             "ClientAddr",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyGroupRateLimitBy_CLIENT_ADDR_XFF_HEADER:  "ClientAddrXFFHeader",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyGroupRateLimitBy_GEO_LOCATION:            "GeoLocation",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyGroupRateLimitBy_GEO_LOCATION_XFF_HEADER: "GeoLocationXFFHeader",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyGroupRateLimitBy_NONE:                    "None",
}

var matchVariableNameStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName_REMOTE_ADDR:     "RemoteAddr",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName_REQUEST_METHOD:  "RequestMethod",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName_QUERY_STRING:    "QueryString",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName_POST_ARGS:       "PostArgs",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName_REQUEST_URI:     "RequestUri",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName_REQUEST_HEADERS: "RequestHeaders",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName_REQUEST_BODY:    "RequestBody",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchVariableName_REQUEST_COOKIES: "RequestCookies",
}

var matchOperatorStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_ANY:                   "Any",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_IP_MATCH:              "IPMatch",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_GEO_MATCH:             "GeoMatch",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_EQUAL:                 "Equal",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_CONTAINS:              "Contains",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_LESS_THAN:             "LessThan",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_GREATER_THAN:          "GreaterThan",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_LESS_THAN_OR_EQUAL:    "LessThanOrEqual",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_GREATER_THAN_OR_EQUAL: "GreaterThanOrEqual",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_BEGINS_WITH:           "BeginsWith",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_ENDS_WITH:             "EndsWith",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMatchOperator_REGEX:                 "Regex",
}

var transformStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyTransform]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyTransform_HTML_ENTITY_DECODE: "HtmlEntityDecode",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyTransform_LOWERCASE:          "Lowercase",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyTransform_REMOVE_NULLS:       "RemoveNulls",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyTransform_TRIM:               "Trim",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyTransform_URL_DECODE:         "UrlDecode",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyTransform_URL_ENCODE:         "UrlEncode",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyTransform_UPPERCASE:          "Uppercase",
}

var managedRuleSetTypeStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyManagedRuleSetType]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyManagedRuleSetType_OWASP:                          "OWASP",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyManagedRuleSetType_MICROSOFT_BOT_MANAGER_RULE_SET: "Microsoft_BotManagerRuleSet",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyManagedRuleSetType_MICROSOFT_DEFAULT_RULE_SET:     "Microsoft_DefaultRuleSet",
}

var ruleOverrideActionStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRuleOverrideAction]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRuleOverrideAction_OVERRIDE_ALLOW:           "Allow",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRuleOverrideAction_OVERRIDE_ANOMALY_SCORING: "AnomalyScoring",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRuleOverrideAction_OVERRIDE_BLOCK:           "Block",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRuleOverrideAction_OVERRIDE_JS_CHALLENGE:    "JSChallenge",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyRuleOverrideAction_OVERRIDE_LOG:             "Log",
}

var exclusionMatchVariableStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_ARG_KEYS:      "RequestArgKeys",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_ARG_NAMES:     "RequestArgNames",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_ARG_VALUES:    "RequestArgValues",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_COOKIE_KEYS:   "RequestCookieKeys",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_COOKIE_NAMES:  "RequestCookieNames",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_COOKIE_VALUES: "RequestCookieValues",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_HEADER_KEYS:   "RequestHeaderKeys",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_HEADER_NAMES:  "RequestHeaderNames",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyExclusionMatchVariable_REQUEST_HEADER_VALUES: "RequestHeaderValues",
}

var selectorMatchOperatorStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicySelectorMatchOperator]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicySelectorMatchOperator_SELECTOR_EQUALS:      "Equals",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicySelectorMatchOperator_SELECTOR_CONTAINS:    "Contains",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicySelectorMatchOperator_SELECTOR_STARTS_WITH: "StartsWith",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicySelectorMatchOperator_SELECTOR_ENDS_WITH:   "EndsWith",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicySelectorMatchOperator_SELECTOR_EQUALS_ANY:  "EqualsAny",
}

var scrubbingMatchVariableStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyScrubbingMatchVariable]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_ARG_NAMES:      "RequestArgNames",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_COOKIE_NAMES:   "RequestCookieNames",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_HEADER_NAMES:   "RequestHeaderNames",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_IP_ADDRESS:     "RequestIPAddress",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_JSON_ARG_NAMES: "RequestJSONArgNames",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyScrubbingMatchVariable_SCRUB_REQUEST_POST_ARG_NAMES: "RequestPostArgNames",
}

var modeStrings = map[azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMode]string{
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMode_PREVENTION: "Prevention",
	azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyMode_DETECTION:  "Detection",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurewebapplicationfirewallpolicyv1.AzureWebApplicationFirewallPolicyStackInput) *Locals {
	locals := &Locals{}

	locals.AzureWebApplicationFirewallPolicy = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureWebApplicationFirewallPolicy.String()),
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

	// The user's spec tags merge over the metadata-derived tags -- user
	// tags deliberately win so an org's governance conventions can
	// override the derived values where they collide.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
