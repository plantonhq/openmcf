# AzureFrontDoorFirewallPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureFrontDoorFirewallPolicySpec** defines a Web Application
Firewall (WAF) policy for Azure Front Door
(Microsoft.Network/frontDoorWebApplicationFirewallPolicies) -- the
rule set Front Door enforces on HTTP traffic at Microsoft's edge,
before requests ever reach an origin.

A policy has three layers, evaluated in order:

1. **Custom rules** (`custom_rules`) -- your own match and rate-limit
   rules, evaluated first by ascending priority. Use them for IP/geo
   allowlists, header-based exceptions, per-client rate limiting, and
   (on PREMIUM) JavaScript-challenge or CAPTCHA bot gates.
2. **Managed rules** (`managed_rules`, PREMIUM only) -- Microsoft's
   curated rule sets (Microsoft_DefaultRuleSet, Bot Manager), tuned
   with scoped exclusions and per-rule overrides.
3. **Policy settings** -- the enforcement mode (Prevention vs
   Detection), body inspection, block-response customization, and
   access-log scrubbing.

**The policy is global and resource-group-scoped** -- it is NOT
nested under a Front Door profile (Azure deploys it once across all
edge locations; no region field exists). It also does nothing on its
own: attaching it to a profile's domains through an
AzureFrontDoorSecurityPolicy (which references this policy's
firewall_policy_id output) is what turns enforcement on. One policy
is commonly shared by many profiles' security policies.

**This is the FRONT DOOR policy type.** The regional WAF policy that
Application Gateways attach is a different ARM resource with a
different rule vocabulary (AzureWebApplicationFirewallPolicy).

**SKU pairing**: the policy carries its own STANDARD/PREMIUM sku and
Azure only allows it to be associated with Front Door profiles of the
SAME sku -- decide the tier once, for the profile and its policies
together. PREMIUM unlocks managed rules and the JS-challenge/CAPTCHA
actions.

**ForceNew fields**: `policy_name`, `sku` -- and Azure additionally
rejects a PREMIUM -> STANDARD sku change outright (upgrades recreate;
downgrades are not supported at all).

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorFirewallPolicy
metadata:
  name: test-front-door-firewall-policy
spec:
  resourceGroup:
    value: test-rg
  policyName: testedgewaf
  # PREMIUM so the plan exercises the managed rules, the challenge
  # lifetimes, and the JS_CHALLENGE custom-rule action -- the deepest
  # rendering path.
  sku: PREMIUM
  mode: PREVENTION
  requestBodyCheckEnabled: true
  redirectUrl: https://example.com/blocked
  customBlockResponseStatusCode: 429
  # base64 of "<h1>blocked</h1>"
  customBlockResponseBody: PGgxPmJsb2NrZWQ8L2gxPg==
  jsChallengeCookieExpirationInMinutes: 45
  captchaCookieExpirationInMinutes: 60
  customRules:
    # Exercises: rate limiting with explicit window/threshold, IP_MATCH,
    # and priority ordering.
    - name: ratelimitapi
      priority: 10
      ruleType: RATE_LIMIT_RULE
      rateLimitDurationInMinutes: 5
      rateLimitThreshold: 300
      action: BLOCK
      matchConditions:
        - matchVariable: REQUEST_URI
          operator: BEGINS_WITH
          matchValues: ["/api/"]
    # Exercises: the keyed-variable selector, negation, transforms
    # (URL_DECODE included -- its "UrlDecode" wire casing differs from
    # the SDK's Go identifier, so the plan keeps that row exercised),
    # and the Premium-only JS_CHALLENGE action.
    - name: botgate
      priority: 20
      ruleType: MATCH_RULE
      action: JS_CHALLENGE
      matchConditions:
        - matchVariable: REQUEST_HEADER
          selector: User-Agent
          operator: CONTAINS
          matchValues: [curl, python-requests]
          negateCondition: false
          transforms: [LOWERCASE, TRIM, URL_DECODE]
    # Exercises: geo allowlisting through negation and a disabled rule.
    - name: geoblock
      enabled: false
      priority: 30
      ruleType: MATCH_RULE
      action: LOG
      matchConditions:
        - matchVariable: SOCKET_ADDR
          operator: GEO_MATCH
          matchValues: [US, DE, IN]
          negateCondition: true
  managedRules:
    # Exercises: the 2.x anomaly-scoring model with a set-wide
    # exclusion, a group override, and a per-rule LOG override.
    - type: Microsoft_DefaultRuleSet
      version: "2.1"
      action: RULE_SET_BLOCK
      exclusions:
        - matchVariable: EXCLUDE_REQUEST_COOKIE_NAMES
          operator: SELECTOR_EQUALS
          selector: session-token
      overrides:
        - ruleGroupName: SQLI
          exclusions:
            - matchVariable: EXCLUDE_REQUEST_BODY_POST_ARG_NAMES
              operator: SELECTOR_STARTS_WITH
              selector: comment
          rules:
            - ruleId: "942100"
              enabled: true
              action: OVERRIDE_LOG
    # Exercises: the bot manager set with the JSChallenge override
    # dialect.
    - type: Microsoft_BotManagerRuleSet
      version: "1.1"
      action: RULE_SET_LOG
      overrides:
        - ruleGroupName: UnknownBots
          rules:
            - ruleId: Bot300700
              enabled: true
              action: OVERRIDE_JS_CHALLENGE
  logScrubbing:
    scrubbingRules:
      # Keyed rule with an explicit selector.
      - matchVariable: SCRUB_REQUEST_HEADER_NAMES
        operator: SELECTOR_EQUALS
        selector: Authorization
      # Equals-any rule on a non-keyed variable.
      - matchVariable: SCRUB_REQUEST_IP_ADDRESS
        operator: SELECTOR_EQUALS_ANY
  tags:
    team: edge-security
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.policyName` | `string` | yes |  |  |
| `spec.sku` | `enum` |  |  |  |
| `spec.mode` | `enum` |  |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.requestBodyCheckEnabled` | `bool` |  | `true` |  |
| `spec.redirectUrl` | `string` |  |  |  |
| `spec.customBlockResponseStatusCode` | `int32` |  |  |  |
| `spec.customBlockResponseBody` | `string` |  |  |  |
| `spec.jsChallengeCookieExpirationInMinutes` | `int32` |  |  |  |
| `spec.captchaCookieExpirationInMinutes` | `int32` |  |  |  |
| `spec.customRules` | `[]AzureFrontDoorFirewallPolicyCustomRule` |  |  |  |
| `spec.customRules[].name` | `string` | yes |  |  |
| `spec.customRules[].enabled` | `bool` |  | `true` |  |
| `spec.customRules[].priority` | `int32` |  | `1` |  |
| `spec.customRules[].ruleType` | `enum` |  |  |  |
| `spec.customRules[].rateLimitDurationInMinutes` | `int32` |  | `1` |  |
| `spec.customRules[].rateLimitThreshold` | `int32` |  | `10` |  |
| `spec.customRules[].action` | `enum` |  |  |  |
| `spec.customRules[].matchConditions` | `[]AzureFrontDoorFirewallPolicyMatchCondition` |  |  |  |
| `spec.customRules[].matchConditions[].matchVariable` | `enum` |  |  |  |
| `spec.customRules[].matchConditions[].selector` | `string` |  |  |  |
| `spec.customRules[].matchConditions[].operator` | `enum` |  |  |  |
| `spec.customRules[].matchConditions[].matchValues` | `[]string` | yes |  |  |
| `spec.customRules[].matchConditions[].negateCondition` | `bool` |  |  |  |
| `spec.customRules[].matchConditions[].transforms` | `[]enum` |  |  |  |
| `spec.managedRules` | `[]AzureFrontDoorFirewallPolicyManagedRuleSet` |  |  |  |
| `spec.managedRules[].type` | `string` | yes |  |  |
| `spec.managedRules[].version` | `string` | yes |  |  |
| `spec.managedRules[].action` | `enum` |  |  |  |
| `spec.managedRules[].exclusions` | `[]AzureFrontDoorFirewallPolicyManagedRuleExclusion` |  |  |  |
| `spec.managedRules[].exclusions[].matchVariable` | `enum` |  |  |  |
| `spec.managedRules[].exclusions[].operator` | `enum` |  |  |  |
| `spec.managedRules[].exclusions[].selector` | `string` | yes |  |  |
| `spec.managedRules[].overrides` | `[]AzureFrontDoorFirewallPolicyManagedRuleGroupOverride` |  |  |  |
| `spec.managedRules[].overrides[].ruleGroupName` | `string` | yes |  |  |
| `spec.managedRules[].overrides[].exclusions` | `[]AzureFrontDoorFirewallPolicyManagedRuleExclusion` |  |  |  |
| `spec.managedRules[].overrides[].exclusions[].matchVariable` | `enum` |  |  |  |
| `spec.managedRules[].overrides[].exclusions[].operator` | `enum` |  |  |  |
| `spec.managedRules[].overrides[].exclusions[].selector` | `string` | yes |  |  |
| `spec.managedRules[].overrides[].rules` | `[]AzureFrontDoorFirewallPolicyManagedRuleOverride` |  |  |  |
| `spec.managedRules[].overrides[].rules[].ruleId` | `string` | yes |  |  |
| `spec.managedRules[].overrides[].rules[].enabled` | `bool` |  | `false` |  |
| `spec.managedRules[].overrides[].rules[].action` | `enum` |  |  |  |
| `spec.managedRules[].overrides[].rules[].exclusions` | `[]AzureFrontDoorFirewallPolicyManagedRuleExclusion` |  |  |  |
| `spec.managedRules[].overrides[].rules[].exclusions[].matchVariable` | `enum` |  |  |  |
| `spec.managedRules[].overrides[].rules[].exclusions[].operator` | `enum` |  |  |  |
| `spec.managedRules[].overrides[].rules[].exclusions[].selector` | `string` | yes |  |  |
| `spec.logScrubbing` | `AzureFrontDoorFirewallPolicyLogScrubbing` |  |  |  |
| `spec.logScrubbing.enabled` | `bool` |  | `true` |  |
| `spec.logScrubbing.scrubbingRules` | `[]AzureFrontDoorFirewallPolicyScrubbingRule` | yes |  |  |
| `spec.logScrubbing.scrubbingRules[].enabled` | `bool` |  | `true` |  |
| `spec.logScrubbing.scrubbingRules[].matchVariable` | `enum` |  |  |  |
| `spec.logScrubbing.scrubbingRules[].operator` | `enum` |  |  |  |
| `spec.logScrubbing.scrubbingRules[].selector` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

The Azure Resource Group the policy is created in. The policy is a
global resource (no region), but every ARM resource belongs to a
resource group for organization, RBAC scoping, and lifecycle
grouping.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.policyName

`string` · required

The policy's name -- unique within the resource group. 1-128
characters; must begin with a letter and may contain only letters
and numbers (Azure allows no hyphens here, unlike most Front Door
names).

**ForceNew**: changing the name replaces the policy, which detaches
it from every security policy that referenced it.

- rule: policy_name must be 1-128 characters, begin with a letter, and contain only letters and numbers (no hyphens)
- rule: {"required":true}

### spec.sku

`enum`

The pricing/capability tier. Unspecified deploys STANDARD -- the
right answer unless you need the managed rule sets or the
JS-challenge/CAPTCHA actions, which are PREMIUM-only. Must match
the sku of every Front Door profile this policy gets associated
with (Azure enforces the pairing at association time).

**ForceNew** -- and Azure additionally refuses a PREMIUM ->
STANDARD downgrade even as a replace, so choose PREMIUM
deliberately.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_sku_unspecified` -- Not specified -- deploys STANDARD, the production default.
- `STANDARD` -- Custom match and rate-limit rules only. Up to 100 domains per security-policy association.
- `PREMIUM` -- Everything in STANDARD plus Microsoft's managed rule sets (Microsoft_DefaultRuleSet, Bot Manager) and the JS-challenge/CAPTCHA custom-rule actions. Up to 500 domains per security-policy association.

### spec.mode

`enum`

The enforcement mode. PREVENTION blocks (or redirects/challenges)
matching requests -- the production posture. DETECTION only logs
matches -- the tuning mode for new policies watching real traffic
before they start blocking it.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_mode_unspecified` -- Not specified -- invalid; choose the enforcement mode explicitly.
- `DETECTION` -- Log matches without acting on them -- the tuning mode for new policies watching real traffic.
- `PREVENTION` -- Act on matches (block, redirect, challenge) -- the production posture.

### spec.enabled

`bool` · optional (explicit presence)

Whether the policy is enforced at all. Azure's default is true;
false parks the policy without detaching it from its security
policies.

- default: `true`

### spec.requestBodyCheckEnabled

`bool` · optional (explicit presence)

Whether request bodies are inspected. Azure's default is true;
turning it off blinds the WAF to body-borne attacks (SQL injection
in POST forms, JSON payload attacks). In DETECTION mode the policy
only logs regardless.

- default: `true`

### spec.redirectUrl

`string`

Where clients are sent when a rule's action is REDIRECT. Must be an
http:// or https:// URL. Azure requires this to be set when any
custom rule or managed-rule override uses the REDIRECT action
(enforced at deploy time -- the rule-to-URL pairing is not
statically checkable).

- rule: redirect_url must be an http:// or https:// URL

### spec.customBlockResponseStatusCode

`int32` · optional (explicit presence)

The HTTP status code returned when a request is blocked. One of
200, 403, 405, 406, 429, or 990-999 (the 99x codes are Azure
WAF-specific). Omit for Azure's default (403).

- rule: custom_block_response_status_code must be one of 200, 403, 405, 406, 429, or 990-999

### spec.customBlockResponseBody

`string`

The response body returned when a request is blocked, as a
base64-encoded string (e.g. a branded HTML error page). Omit for
Azure's default block page.

- rule: custom_block_response_body must be base64-encoded (encode your HTML/text before setting it)

### spec.jsChallengeCookieExpirationInMinutes

`int32` · optional (explicit presence)

How long a solved JavaScript challenge stays valid before the
client is re-challenged, in minutes (5-1440). **PREMIUM only** --
on PREMIUM, Azure always enables the JS-challenge policy and both
engines send 30 (Azure's default) when this is unset; this field
tunes that lifetime. Leave unset on STANDARD (Azure rejects it, so
the modules never send it there and no platform default applies).
The JavaScript challenge is an Azure PREVIEW feature.

- rule: {"int32":{"lte":1440,"gte":5}}

### spec.captchaCookieExpirationInMinutes

`int32` · optional (explicit presence)

How long a solved CAPTCHA stays valid before the client is
re-challenged, in minutes (5-1440). **PREMIUM only** -- on PREMIUM,
Azure always enables the CAPTCHA policy and both engines send 30
(Azure's default) when this is unset; this field tunes that
lifetime. Leave unset on STANDARD (Azure rejects it, so the modules
never send it there and no platform default applies). CAPTCHA is an
Azure PREVIEW feature.

- rule: {"int32":{"lte":1440,"gte":5}}

### spec.customRules

`[]AzureFrontDoorFirewallPolicyCustomRule`

Your own match and rate-limit rules, evaluated before the managed
rule sets in ascending priority order. Up to 100 rules.

- rule: {"repeated":{"maxItems":"100"}}

### spec.customRules[].name

`string` · required

The rule's name, shown in WAF logs and metrics.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.customRules[].enabled

`bool` · optional (explicit presence)

Whether the rule is evaluated. Azure's default is true; disable to
stage a rule without deleting it.

- default: `true`

### spec.customRules[].priority

`int32` · optional (explicit presence)

Evaluation order among custom rules: LOWER runs first. Azure's
default is 1.

- default: `1`
- rule: {"int32":{"gte":0}}

### spec.customRules[].ruleType

`enum`

MATCH_RULE acts on every matching request; RATE_LIMIT_RULE acts
only when matching requests exceed rate_limit_threshold within
rate_limit_duration_in_minutes.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_custom_rule_type_unspecified` -- Not specified -- invalid; declare the rule type.
- `MATCH_RULE` -- Acts on every request matching the conditions.
- `RATE_LIMIT_RULE` -- Acts only when matching requests exceed the rate-limit threshold within the window.

### spec.customRules[].rateLimitDurationInMinutes

`int32` · optional (explicit presence)

For RATE_LIMIT_RULE: the sliding window the threshold counts over,
in minutes. Azure's default is 1. Ignored by MATCH_RULE rules.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.customRules[].rateLimitThreshold

`int32` · optional (explicit presence)

For RATE_LIMIT_RULE: the number of matching requests allowed per
client within the window before the action fires. Azure's default
is 10. Ignored by MATCH_RULE rules.

- default: `10`
- rule: {"int32":{"gte":1}}

### spec.customRules[].action

`enum`

What happens when the rule matches. LOG records the match and
continues evaluation; ALLOW lets the request through skipping the
remaining rules; BLOCK rejects it; REDIRECT sends the client to
the policy's redirect_url; JS_CHALLENGE and CAPTCHA (PREMIUM only)
gate the request behind a browser challenge.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_custom_rule_action_unspecified` -- Not specified -- invalid; declare the action.
- `ALLOW` -- Let the request through, skipping the remaining custom AND managed rules -- the allowlist gesture.
- `BLOCK` -- Reject the request with the policy's block response.
- `LOG` -- Record the match and continue evaluating.
- `REDIRECT` -- Send the client to the policy's redirect_url.
- `JS_CHALLENGE` -- Serve a JavaScript challenge; only real browsers that solve it proceed (bot mitigation without user friction). PREMIUM only; an Azure PREVIEW feature.
- `CAPTCHA` -- Serve a CAPTCHA; only humans that solve it proceed. PREMIUM only; an Azure PREVIEW feature.

### spec.customRules[].matchConditions

`[]AzureFrontDoorFirewallPolicyMatchCondition`

The conditions that select requests. ALL conditions must match for
the action to apply (conditions AND together; the values WITHIN one
condition OR together). Up to 10 conditions per rule.

- rule: {"repeated":{"maxItems":"10"}}

### spec.customRules[].matchConditions[].matchVariable

`enum`

The request part the condition inspects.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_match_variable_unspecified` -- Not specified -- invalid; declare the variable.
- `COOKIES` -- Request cookies (selector picks one).
- `POST_ARGS` -- POST form arguments (selector picks one).
- `QUERY_STRING` -- The raw query string (selector picks one argument).
- `REMOTE_ADDR` -- The client's IP address.
- `REQUEST_BODY` -- The request body.
- `REQUEST_HEADER` -- Request headers (selector picks one).
- `REQUEST_METHOD` -- The HTTP method.
- `REQUEST_URI` -- The request URI (path + query).
- `SOCKET_ADDR` -- The client's socket address (the direct TCP peer -- differs from REMOTE_ADDR behind proxies that set X-Forwarded-For).

### spec.customRules[].matchConditions[].selector

`string`

For the keyed variables (COOKIES, POST_ARGS, QUERY_STRING,
REQUEST_HEADER): the specific key to inspect (e.g. "User-Agent"
with REQUEST_HEADER). Azure requires it for those variables and
ignores it elsewhere (enforced at deploy time).

### spec.customRules[].matchConditions[].operator

`enum`

How the variable is compared against match_values.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_operator_unspecified` -- Not specified -- invalid; declare the operator.
- `ANY` -- Matches every request carrying the variable.
- `BEGINS_WITH` -- String prefix match.
- `CONTAINS` -- Substring containment.
- `ENDS_WITH` -- String suffix match.
- `EQUAL` -- Exact string equality.
- `GEO_MATCH` -- The client's country matches (two-letter ISO codes in match_values).
- `GREATER_THAN` -- Numeric greater-than.
- `GREATER_THAN_OR_EQUAL` -- Numeric greater-than-or-equal.
- `IP_MATCH` -- The variable is an IP inside one of the match_values CIDRs/addresses.
- `LESS_THAN` -- Numeric less-than.
- `LESS_THAN_OR_EQUAL` -- Numeric less-than-or-equal.
- `REG_EX` -- Regular-expression match.

### spec.customRules[].matchConditions[].matchValues

`[]string` · required

The values compared against (OR semantics). 1-600 values, each
1-256 characters. IP_MATCH takes CIDRs or addresses; GEO_MATCH
takes two-letter ISO country codes.

- rule: {"repeated":{"minItems":"1","maxItems":"600","items":{"string":{"minLen":"1","maxLen":"256"}}}}

### spec.customRules[].matchConditions[].negateCondition

`bool`

Inverts the result: the condition matches when the comparison does
NOT hold. Azure's default is false.

### spec.customRules[].matchConditions[].transforms

`[]enum`

Normalizations applied to the variable before comparison (e.g.
LOWERCASE + URL_DECODE to catch encoding evasions). Up to 5.

- rule: {"repeated":{"maxItems":"5","items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_transform_unspecified` -- Not specified -- invalid; list only real transforms.
- `LOWERCASE` -- Lowercase the value.
- `REMOVE_NULLS` -- Strip NULL bytes.
- `TRIM` -- Trim leading/trailing whitespace.
- `UPPERCASE` -- Uppercase the value.
- `URL_DECODE` -- URL-decode the value (%27 -> ').
- `URL_ENCODE` -- URL-encode the value.

### spec.managedRules

`[]AzureFrontDoorFirewallPolicyManagedRuleSet`

Microsoft's curated rule sets with per-rule tuning. **PREMIUM
only** (Azure rejects managed rules on STANDARD). The known set
types and versions: `Microsoft_DefaultRuleSet` (1.1, 2.0, 2.1 --
the OWASP successor with Microsoft threat-intelligence rules),
`Microsoft_BotManagerRuleSet` (1.0, 1.1 -- good/bad/unknown bot
classification), and the legacy `DefaultRuleSet` (1.0,
preview-0.1). Azure ships new types and versions server-side, so
the type is a string, not a closed list.

- rule: {"repeated":{"maxItems":"100"}}
- rule: the legacy DefaultRuleSet type only supports versions '1.0' and 'preview-0.1' -- for 1.1 and above use the Microsoft_DefaultRuleSet type
- rule: Microsoft_DefaultRuleSet starts at version '1.1' ('1.1', '2.0', '2.1') -- for 1.0 and below use the legacy DefaultRuleSet type
- rule: rule overrides on 2.0+ rule sets must use OVERRIDE_ANOMALY_SCORING or OVERRIDE_LOG (the anomaly-scoring model); OVERRIDE_ANOMALY_SCORING is invalid below 2.0
- rule: the OVERRIDE_JS_CHALLENGE action is only valid on bot-manager rule sets (Microsoft_BotManagerRuleSet)

### spec.managedRules[].type

`string` · required

The rule-set family. Known values: `Microsoft_DefaultRuleSet` (the
OWASP successor with Microsoft threat-intelligence rules),
`Microsoft_BotManagerRuleSet` (bot classification), and the legacy
`DefaultRuleSet`. Azure ships new families server-side, so this is
a string, not a closed list.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules[].version

`string` · required

The rule-set version. `Microsoft_DefaultRuleSet`: "1.1", "2.0",
"2.1" (2.x uses anomaly scoring); `Microsoft_BotManagerRuleSet`:
"1.0", "1.1"; the legacy `DefaultRuleSet`: "1.0", "preview-0.1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules[].action

`enum`

The action for requests that trip the set. With
Microsoft_DefaultRuleSet 2.x (anomaly scoring), this is the action
taken when the accumulated anomaly score crosses the blocking
threshold.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_managed_rule_set_action_unspecified` -- Not specified -- invalid; declare the set's action.
- `RULE_SET_BLOCK` -- Reject requests that trip the set.
- `RULE_SET_LOG` -- Record matches without acting.
- `RULE_SET_REDIRECT` -- Send matching clients to the policy's redirect_url.

### spec.managedRules[].exclusions

`[]AzureFrontDoorFirewallPolicyManagedRuleExclusion`

Set-wide exclusions: request parts (by key) that EVERY rule in the
set skips -- the tool for false positives on specific cookies,
headers, or form fields, without disabling rules. Up to 100.

- rule: {"repeated":{"maxItems":"100"}}

### spec.managedRules[].exclusions[].matchVariable

`enum`

The keyed request collection the exclusion applies to.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_exclusion_match_variable_unspecified` -- Not specified -- invalid; declare the collection.
- `EXCLUDE_QUERY_STRING_ARG_NAMES` -- Query-string argument names.
- `EXCLUDE_REQUEST_BODY_JSON_ARG_NAMES` -- JSON body argument names.
- `EXCLUDE_REQUEST_BODY_POST_ARG_NAMES` -- POST body argument names.
- `EXCLUDE_REQUEST_COOKIE_NAMES` -- Cookie names.
- `EXCLUDE_REQUEST_HEADER_NAMES` -- Header names.

### spec.managedRules[].exclusions[].operator

`enum`

How selector picks keys out of the collection.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_selector_operator_unspecified` -- Not specified -- invalid for exclusions; means SELECTOR_EQUALS for log-scrubbing rules (Azure's default there).
- `SELECTOR_CONTAINS` -- The key contains the selector.
- `SELECTOR_ENDS_WITH` -- The key ends with the selector.
- `SELECTOR_EQUALS` -- The key equals the selector exactly.
- `SELECTOR_EQUALS_ANY` -- Every key in the collection matches (no selector).
- `SELECTOR_STARTS_WITH` -- The key starts with the selector.

### spec.managedRules[].exclusions[].selector

`string` · required

The key (or key fragment, per the operator) to exclude -- e.g. a
session cookie name that trips SQL-injection rules.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules[].overrides

`[]AzureFrontDoorFirewallPolicyManagedRuleGroupOverride`

Per-rule-group tuning: disable individual rules, change their
action, or scope exclusions to a group or rule. Up to 100 groups.

- rule: {"repeated":{"maxItems":"100"}}

### spec.managedRules[].overrides[].ruleGroupName

`string` · required

The rule group, using Azure's exact group names (e.g. "SQLI",
"XSS", "PROTOCOL-ATTACK" for Microsoft_DefaultRuleSet; "BadBots"
for the bot manager).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules[].overrides[].exclusions

`[]AzureFrontDoorFirewallPolicyManagedRuleExclusion`

Group-scoped exclusions: request parts every rule IN THIS GROUP
skips. Up to 100.

- rule: {"repeated":{"maxItems":"100"}}

### spec.managedRules[].overrides[].exclusions[].matchVariable

`enum`

The keyed request collection the exclusion applies to.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_exclusion_match_variable_unspecified` -- Not specified -- invalid; declare the collection.
- `EXCLUDE_QUERY_STRING_ARG_NAMES` -- Query-string argument names.
- `EXCLUDE_REQUEST_BODY_JSON_ARG_NAMES` -- JSON body argument names.
- `EXCLUDE_REQUEST_BODY_POST_ARG_NAMES` -- POST body argument names.
- `EXCLUDE_REQUEST_COOKIE_NAMES` -- Cookie names.
- `EXCLUDE_REQUEST_HEADER_NAMES` -- Header names.

### spec.managedRules[].overrides[].exclusions[].operator

`enum`

How selector picks keys out of the collection.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_selector_operator_unspecified` -- Not specified -- invalid for exclusions; means SELECTOR_EQUALS for log-scrubbing rules (Azure's default there).
- `SELECTOR_CONTAINS` -- The key contains the selector.
- `SELECTOR_ENDS_WITH` -- The key ends with the selector.
- `SELECTOR_EQUALS` -- The key equals the selector exactly.
- `SELECTOR_EQUALS_ANY` -- Every key in the collection matches (no selector).
- `SELECTOR_STARTS_WITH` -- The key starts with the selector.

### spec.managedRules[].overrides[].exclusions[].selector

`string` · required

The key (or key fragment, per the operator) to exclude -- e.g. a
session cookie name that trips SQL-injection rules.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules[].overrides[].rules

`[]AzureFrontDoorFirewallPolicyManagedRuleOverride`

The individual rules to tune within the group. Up to 1000.

- rule: {"repeated":{"maxItems":"1000"}}

### spec.managedRules[].overrides[].rules[].ruleId

`string` · required

The rule ID, from the rule set's documentation -- note each
family's own ID shape: Microsoft_DefaultRuleSet rules are numeric
(e.g. "942100" for the SQL-injection libinjection rule) while
Microsoft_BotManagerRuleSet rules carry a "Bot" prefix (e.g.
"Bot300700" for unknown bots). Azure rejects an ID that does not
exist in the chosen set and version.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules[].overrides[].rules[].enabled

`bool` · optional (explicit presence)

Whether the rule runs. Azure's default here is FALSE -- listing a
rule without enabled=true DISABLES it, which is the common tuning
gesture for false positives.

- default: `false`

### spec.managedRules[].overrides[].rules[].action

`enum`

Replace the rule's action. On 2.0+ rule sets only
OVERRIDE_ANOMALY_SCORING and OVERRIDE_LOG are valid (the
anomaly-scoring model); OVERRIDE_JS_CHALLENGE only on the bot
manager set.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_managed_rule_override_action_unspecified` -- Not specified -- invalid; declare the substituted action.
- `OVERRIDE_ALLOW` -- Let matching requests through.
- `OVERRIDE_ANOMALY_SCORING` -- Contribute to the anomaly score instead of acting directly -- the only rule action (besides OVERRIDE_LOG) on 2.0+ rule sets.
- `OVERRIDE_BLOCK` -- Reject matching requests.
- `OVERRIDE_CAPTCHA` -- Serve a CAPTCHA.
- `OVERRIDE_JS_CHALLENGE` -- Serve a JavaScript challenge (Microsoft_BotManagerRuleSet only).
- `OVERRIDE_LOG` -- Record matches without acting.
- `OVERRIDE_REDIRECT` -- Send matching clients to the policy's redirect_url.

### spec.managedRules[].overrides[].rules[].exclusions

`[]AzureFrontDoorFirewallPolicyManagedRuleExclusion`

Rule-scoped exclusions: request parts THIS RULE skips. Up to 100.

- rule: {"repeated":{"maxItems":"100"}}

### spec.managedRules[].overrides[].rules[].exclusions[].matchVariable

`enum`

The keyed request collection the exclusion applies to.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_exclusion_match_variable_unspecified` -- Not specified -- invalid; declare the collection.
- `EXCLUDE_QUERY_STRING_ARG_NAMES` -- Query-string argument names.
- `EXCLUDE_REQUEST_BODY_JSON_ARG_NAMES` -- JSON body argument names.
- `EXCLUDE_REQUEST_BODY_POST_ARG_NAMES` -- POST body argument names.
- `EXCLUDE_REQUEST_COOKIE_NAMES` -- Cookie names.
- `EXCLUDE_REQUEST_HEADER_NAMES` -- Header names.

### spec.managedRules[].overrides[].rules[].exclusions[].operator

`enum`

How selector picks keys out of the collection.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_selector_operator_unspecified` -- Not specified -- invalid for exclusions; means SELECTOR_EQUALS for log-scrubbing rules (Azure's default there).
- `SELECTOR_CONTAINS` -- The key contains the selector.
- `SELECTOR_ENDS_WITH` -- The key ends with the selector.
- `SELECTOR_EQUALS` -- The key equals the selector exactly.
- `SELECTOR_EQUALS_ANY` -- Every key in the collection matches (no selector).
- `SELECTOR_STARTS_WITH` -- The key starts with the selector.

### spec.managedRules[].overrides[].rules[].exclusions[].selector

`string` · required

The key (or key fragment, per the operator) to exclude -- e.g. a
session cookie name that trips SQL-injection rules.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.logScrubbing

`AzureFrontDoorFirewallPolicyLogScrubbing`

Scrub (mask) sensitive request data out of the WAF's logs before
they are written -- auth headers, PII query arguments, client IPs.
Omit to leave scrubbing off. WAF log scrubbing is an Azure PREVIEW
feature.

### spec.logScrubbing.enabled

`bool` · optional (explicit presence)

Whether scrubbing is active. Azure's default is true (when the
block is present); false stages the rules without applying them.

- default: `true`

### spec.logScrubbing.scrubbingRules

`[]AzureFrontDoorFirewallPolicyScrubbingRule` · required

The request parts to redact. At least one; up to 100.

- rule: {"repeated":{"minItems":"1","maxItems":"100"}}
- rule: SCRUB_REQUEST_IP_ADDRESS and SCRUB_REQUEST_URI are not keyed collections -- they only accept the SELECTOR_EQUALS_ANY operator
- rule: selector is required with SELECTOR_EQUALS (the default) and must be omitted with SELECTOR_EQUALS_ANY

### spec.logScrubbing.scrubbingRules[].enabled

`bool` · optional (explicit presence)

Whether this rule is active. Azure's default is true.

- default: `true`

### spec.logScrubbing.scrubbingRules[].matchVariable

`enum`

The request part to redact.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_scrubbing_match_variable_unspecified` -- Not specified -- invalid; declare the part to redact.
- `SCRUB_QUERY_STRING_ARG_NAMES` -- Query-string argument names.
- `SCRUB_REQUEST_BODY_JSON_ARG_NAMES` -- JSON body argument names.
- `SCRUB_REQUEST_BODY_POST_ARG_NAMES` -- POST body argument names.
- `SCRUB_REQUEST_COOKIE_NAMES` -- Cookie names.
- `SCRUB_REQUEST_HEADER_NAMES` -- Header names.
- `SCRUB_REQUEST_IP_ADDRESS` -- The client IP address (SELECTOR_EQUALS_ANY only).
- `SCRUB_REQUEST_URI` -- The request URI (SELECTOR_EQUALS_ANY only).

### spec.logScrubbing.scrubbingRules[].operator

`enum`

How selector picks keys: SELECTOR_EQUALS redacts one named key
(Azure's default), SELECTOR_EQUALS_ANY redacts every key in the
collection. SCRUB_REQUEST_IP_ADDRESS and SCRUB_REQUEST_URI only
accept SELECTOR_EQUALS_ANY (they are not keyed collections).

- rule: log-scrubbing rules support only SELECTOR_EQUALS and SELECTOR_EQUALS_ANY

Allowed values (use exactly as shown):

- `azure_front_door_firewall_policy_selector_operator_unspecified` -- Not specified -- invalid for exclusions; means SELECTOR_EQUALS for log-scrubbing rules (Azure's default there).
- `SELECTOR_CONTAINS` -- The key contains the selector.
- `SELECTOR_ENDS_WITH` -- The key ends with the selector.
- `SELECTOR_EQUALS` -- The key equals the selector exactly.
- `SELECTOR_EQUALS_ANY` -- Every key in the collection matches (no selector).
- `SELECTOR_STARTS_WITH` -- The key starts with the selector.

### spec.logScrubbing.scrubbingRules[].selector

`string`

The key to redact (e.g. "Authorization" with
SCRUB_REQUEST_HEADER_NAMES). Required with SELECTOR_EQUALS; must be
omitted with SELECTOR_EQUALS_ANY.

### spec.tags

`map<string, string>`

Free-form tags applied to the policy, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins. Tags are Azure's governance
surface -- Azure Policy enforces them and Microsoft Cost Management
groups by them. Updatable in place.

## Validation Rules

- `front_door_firewall_policy_managed_rules_premium_only`: managed_rules are only supported on the PREMIUM sku -- Azure's managed rule sets (Microsoft_DefaultRuleSet, Bot Manager) are a Premium Front Door capability
- `front_door_firewall_policy_js_challenge_premium_only`: js_challenge_cookie_expiration_in_minutes is only supported on the PREMIUM sku -- the JavaScript challenge is a Premium Front Door capability
- `front_door_firewall_policy_captcha_premium_only`: captcha_cookie_expiration_in_minutes is only supported on the PREMIUM sku -- CAPTCHA is a Premium Front Door capability
- `front_door_firewall_policy_challenge_actions_premium_only`: custom rules with the JS_CHALLENGE or CAPTCHA action are only supported on the PREMIUM sku

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFrontDoorFirewallPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.firewall_policy_id` | `string` | The Azure Resource Manager ID of the WAF policy -- what AzureFrontDoorSecurityPolicy's firewall_policy_id references to attach this policy to a profile's domains (the policy enforces nothing until a security policy associates it). Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/{name} |
| `status.outputs.firewall_policy_name` | `string` | The policy's name -- unique within its resource group. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFrontDoorSecurityPolicy | `spec.firewallPolicyId` | `status.outputs.firewall_policy_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
