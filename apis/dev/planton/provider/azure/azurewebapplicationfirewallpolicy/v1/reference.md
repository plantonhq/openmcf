# AzureWebApplicationFirewallPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureWebApplicationFirewallPolicySpec** defines a regional Web
Application Firewall (WAF) policy
(Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies) --
the rule set an Azure Application Gateway enforces on HTTP traffic.

A policy has three layers, evaluated in order:

1. **Custom rules** (`custom_rules`) -- your own match and rate-limit
   rules, evaluated first by ascending priority. Use them for IP/geo
   allowlists, header-based exceptions, and per-client rate limiting.
2. **Managed rules** (`managed_rules`) -- Microsoft's curated rule sets
   (OWASP core rule set, bot manager), tuned with per-rule overrides and
   scoped exclusions.
3. **Policy settings** (`policy_settings`) -- the enforcement mode
   (Prevention vs Detection), body-inspection limits, and log scrubbing.

One policy attaches to Application Gateways at three levels, all by
reference to its `policy_id` output: gateway-wide
(`firewall_policy_id`), per HTTP listener, and per URL path rule -- so a
single org-standard policy governs many gateways while specific routes
carry stricter or looser variants.

This is the APPLICATION GATEWAY policy type. Azure Front Door's WAF
policy is a different ARM resource with a different rule vocabulary.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureWebApplicationFirewallPolicy
metadata:
  name: test-waf
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  policyName: test-waf-policy
  # Exercises both custom-rule types, the action/operator/variable/
  # transform enum mappings, and the rate-limit trio.
  customRules:
    - name: allowOffice
      priority: 5
      ruleType: MATCH_RULE
      action: ALLOW
      matchConditions:
        - matchVariables:
            - variableName: REMOTE_ADDR
          operator: IP_MATCH
          matchValues:
            - "203.0.113.0/24"
    - name: throttleApi
      priority: 20
      ruleType: RATE_LIMIT_RULE
      action: BLOCK
      rateLimitDuration: ONE_MIN
      rateLimitThreshold: 300
      groupRateLimitBy: CLIENT_ADDR
      matchConditions:
        - matchVariables:
            - variableName: REQUEST_URI
          operator: BEGINS_WITH
          matchValues:
            - /api/
    - name: blockBadAgents
      priority: 30
      ruleType: MATCH_RULE
      action: JS_CHALLENGE
      matchConditions:
        - matchVariables:
            - variableName: REQUEST_HEADERS
              selector: User-Agent
          operator: CONTAINS
          matchValues:
            - scanner
          transforms:
            - LOWERCASE
            - TRIM
  managedRules:
    # Exercises the exclusion path with a narrowed rule set.
    exclusions:
      - matchVariable: REQUEST_COOKIE_NAMES
        selectorMatchOperator: SELECTOR_EQUALS
        selector: session-token
        excludedRuleSet:
          ruleGroups:
            - ruleGroupName: REQUEST-942-APPLICATION-ATTACK-SQLI
              excludedRules:
                - "942440"
    managedRuleSets:
      - version: "3.2"
        ruleGroupOverrides:
          - ruleGroupName: REQUEST-920-PROTOCOL-ENFORCEMENT
            rules:
              - id: "920300"
                enabled: false
              - id: "920440"
                enabled: true
                action: OVERRIDE_LOG
      - type: MICROSOFT_BOT_MANAGER_RULE_SET
        version: "1.1"
  # Exercises the mode enum, the body dials, and log scrubbing.
  policySettings:
    mode: PREVENTION
    maxRequestBodySizeInKb: 256
    fileUploadLimitInMb: 200
    jsChallengeCookieExpirationInMinutes: 60
    logScrubbing:
      rules:
        - matchVariable: SCRUB_REQUEST_HEADER_NAMES
          selectorMatchOperator: SELECTOR_EQUALS
          selector: Authorization
        - matchVariable: SCRUB_REQUEST_IP_ADDRESS
          selectorMatchOperator: SELECTOR_EQUALS
  tags:
    team: security
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.policyName` | `string` | yes |  |  |
| `spec.customRules` | `[]AzureWebApplicationFirewallPolicyCustomRule` |  |  |  |
| `spec.customRules[].name` | `string` |  |  |  |
| `spec.customRules[].priority` | `int32` | yes |  |  |
| `spec.customRules[].enabled` | `bool` |  | `true` |  |
| `spec.customRules[].ruleType` | `enum` | yes |  |  |
| `spec.customRules[].action` | `enum` | yes |  |  |
| `spec.customRules[].rateLimitDuration` | `enum` |  |  |  |
| `spec.customRules[].rateLimitThreshold` | `int32` |  |  |  |
| `spec.customRules[].groupRateLimitBy` | `enum` |  |  |  |
| `spec.customRules[].matchConditions` | `[]AzureWebApplicationFirewallPolicyMatchCondition` | yes |  |  |
| `spec.customRules[].matchConditions[].matchVariables` | `[]AzureWebApplicationFirewallPolicyMatchVariable` | yes |  |  |
| `spec.customRules[].matchConditions[].matchVariables[].variableName` | `enum` | yes |  |  |
| `spec.customRules[].matchConditions[].matchVariables[].selector` | `string` |  |  |  |
| `spec.customRules[].matchConditions[].operator` | `enum` | yes |  |  |
| `spec.customRules[].matchConditions[].matchValues` | `[]string` |  |  |  |
| `spec.customRules[].matchConditions[].negationCondition` | `bool` |  |  |  |
| `spec.customRules[].matchConditions[].transforms` | `[]enum` |  |  |  |
| `spec.managedRules` | `AzureWebApplicationFirewallPolicyManagedRules` | yes |  |  |
| `spec.managedRules.exclusions` | `[]AzureWebApplicationFirewallPolicyManagedRulesExclusion` |  |  |  |
| `spec.managedRules.exclusions[].matchVariable` | `enum` | yes |  |  |
| `spec.managedRules.exclusions[].selectorMatchOperator` | `enum` | yes |  |  |
| `spec.managedRules.exclusions[].selector` | `string` | yes |  |  |
| `spec.managedRules.exclusions[].excludedRuleSet` | `AzureWebApplicationFirewallPolicyExcludedRuleSet` |  |  |  |
| `spec.managedRules.exclusions[].excludedRuleSet.type` | `enum` |  |  |  |
| `spec.managedRules.exclusions[].excludedRuleSet.version` | `string` |  | `3.2` |  |
| `spec.managedRules.exclusions[].excludedRuleSet.ruleGroups` | `[]AzureWebApplicationFirewallPolicyExcludedRuleGroup` |  |  |  |
| `spec.managedRules.exclusions[].excludedRuleSet.ruleGroups[].ruleGroupName` | `string` | yes |  |  |
| `spec.managedRules.exclusions[].excludedRuleSet.ruleGroups[].excludedRules` | `[]string` |  |  |  |
| `spec.managedRules.managedRuleSets` | `[]AzureWebApplicationFirewallPolicyManagedRuleSet` | yes |  |  |
| `spec.managedRules.managedRuleSets[].type` | `enum` |  |  |  |
| `spec.managedRules.managedRuleSets[].version` | `string` | yes |  |  |
| `spec.managedRules.managedRuleSets[].ruleGroupOverrides` | `[]AzureWebApplicationFirewallPolicyRuleGroupOverride` |  |  |  |
| `spec.managedRules.managedRuleSets[].ruleGroupOverrides[].ruleGroupName` | `string` | yes |  |  |
| `spec.managedRules.managedRuleSets[].ruleGroupOverrides[].rules` | `[]AzureWebApplicationFirewallPolicyRuleOverride` | yes |  |  |
| `spec.managedRules.managedRuleSets[].ruleGroupOverrides[].rules[].id` | `string` | yes |  |  |
| `spec.managedRules.managedRuleSets[].ruleGroupOverrides[].rules[].enabled` | `bool` |  | `false` |  |
| `spec.managedRules.managedRuleSets[].ruleGroupOverrides[].rules[].action` | `enum` |  |  |  |
| `spec.policySettings` | `AzureWebApplicationFirewallPolicySettings` |  |  |  |
| `spec.policySettings.enabled` | `bool` |  | `true` |  |
| `spec.policySettings.mode` | `enum` |  |  |  |
| `spec.policySettings.requestBodyCheck` | `bool` |  | `true` |  |
| `spec.policySettings.requestBodyEnforcement` | `bool` |  | `true` |  |
| `spec.policySettings.requestBodyInspectLimitInKb` | `int32` |  | `128` |  |
| `spec.policySettings.maxRequestBodySizeInKb` | `int32` |  | `128` |  |
| `spec.policySettings.fileUploadEnforcement` | `bool` |  | `true` |  |
| `spec.policySettings.fileUploadLimitInMb` | `int32` |  | `100` |  |
| `spec.policySettings.jsChallengeCookieExpirationInMinutes` | `int32` |  | `30` |  |
| `spec.policySettings.logScrubbing` | `AzureWebApplicationFirewallPolicyLogScrubbing` |  |  |  |
| `spec.policySettings.logScrubbing.enabled` | `bool` |  | `true` |  |
| `spec.policySettings.logScrubbing.rules` | `[]AzureWebApplicationFirewallPolicyLogScrubbingRule` | yes |  |  |
| `spec.policySettings.logScrubbing.rules[].enabled` | `bool` |  | `true` |  |
| `spec.policySettings.logScrubbing.rules[].matchVariable` | `enum` | yes |  |  |
| `spec.policySettings.logScrubbing.rules[].selectorMatchOperator` | `enum` |  |  |  |
| `spec.policySettings.logScrubbing.rules[].selector` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region the policy lives in. A policy can only attach to
Application Gateways in the same region. Changing it replaces the
policy.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the policy will be created in. Can be a
literal resource-group name or a reference to an AzureResourceGroup's
name output. Changing it replaces the policy.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.policyName

`string` · required

The policy's name, unique within the resource group. Changing it
replaces the policy.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.customRules

`[]AzureWebApplicationFirewallPolicyCustomRule`

Custom rules, evaluated before the managed rule sets in ascending
priority order. The first matching rule's action wins (except LOG,
which records and continues).

- rule: rate_limit_duration and rate_limit_threshold are required for RATE_LIMIT_RULE and must be omitted for MATCH_RULE (group_rate_limit_by likewise)
- rule: a RATE_LIMIT_RULE cannot use the ALLOW action -- exceeding a rate limit is blocked, logged, or challenged

### spec.customRules[].name

`string`

The rule's name, shown in WAF logs. Letters and digits, starting with
a letter.

- rule: custom rule names start with a letter and contain only letters and digits

### spec.customRules[].priority

`int32` · required

Evaluation order among custom rules: LOWER runs first, and each rule
must have a unique priority. 1-100.

- rule: {"required":true,"int32":{"lte":100,"gte":1}}

### spec.customRules[].enabled

`bool` · optional (explicit presence)

Whether the rule is evaluated. Azure's default is true; disable to
stage a rule without deleting it.

- default: `true`

### spec.customRules[].ruleType

`enum` · required

MATCH_RULE acts on every matching request; RATE_LIMIT_RULE acts only
when matching requests exceed rate_limit_threshold within
rate_limit_duration.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_custom_rule_type_unspecified` -- Not specified -- invalid; declare the rule type.
- `MATCH_RULE` -- Acts on every request matching the conditions.
- `RATE_LIMIT_RULE` -- Acts only when matching requests exceed the rate-limit threshold within the window.

### spec.customRules[].action

`enum` · required

What happens when the rule matches. LOG records the match and
continues evaluation; the others stop it. Rate-limit rules cannot
ALLOW (Azure rejects it) -- exceeding a rate limit is only ever
blocked, logged, or challenged.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_custom_rule_action_unspecified` -- Not specified -- invalid; declare the action.
- `ALLOW` -- Let the request through, skipping the remaining custom AND managed rules -- the allowlist gesture.
- `BLOCK` -- Reject the request.
- `LOG` -- Record the match and continue evaluating.
- `JS_CHALLENGE` -- Serve a JavaScript challenge; only real browsers that solve it proceed (bot mitigation without CAPTCHAs).

### spec.customRules[].rateLimitDuration

`enum`

For RATE_LIMIT_RULE: the sliding window the threshold counts over.

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_rate_limit_duration_unspecified` -- Not specified: only valid on MATCH_RULE rules.
- `ONE_MIN` -- A one-minute sliding window.
- `FIVE_MINS` -- A five-minute sliding window.

### spec.customRules[].rateLimitThreshold

`int32` · optional (explicit presence)

For RATE_LIMIT_RULE: the number of matching requests allowed per
group (see group_rate_limit_by) within the window before the action
fires.

- rule: {"int32":{"gte":1}}

### spec.customRules[].groupRateLimitBy

`enum`

For RATE_LIMIT_RULE: how requests are grouped for counting --
per-client address (the norm), per-XFF-header address (behind another
proxy), per-geo, or NONE (one global counter for all matching
traffic).

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_group_rate_limit_by_unspecified` -- Not specified: Azure groups by client address.
- `CLIENT_ADDR` -- Count per client IP (the norm).
- `CLIENT_ADDR_XFF_HEADER` -- Count per address in the X-Forwarded-For header -- for gateways behind another proxy.
- `GEO_LOCATION` -- Count per client geography.
- `GEO_LOCATION_XFF_HEADER` -- Count per geography resolved from the X-Forwarded-For header.
- `NONE` -- One global counter for all matching traffic.

### spec.customRules[].matchConditions

`[]AzureWebApplicationFirewallPolicyMatchCondition` · required

The conditions that select requests. Multiple conditions AND
together; values within one condition OR together.

- rule: {"repeated":{"minItems":"1"}}
- rule: match_values is required for every operator except ANY (and must be omitted for ANY)

### spec.customRules[].matchConditions[].matchVariables

`[]AzureWebApplicationFirewallPolicyMatchVariable` · required

The request parts the condition inspects. Multiple variables OR
together.

- rule: {"repeated":{"minItems":"1"}}

### spec.customRules[].matchConditions[].matchVariables[].variableName

`enum` · required

Which part of the request to inspect.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_match_variable_name_unspecified` -- Not specified -- invalid; declare the variable.
- `REMOTE_ADDR` -- The client's IP address.
- `REQUEST_METHOD` -- The HTTP method.
- `QUERY_STRING` -- The raw query string.
- `POST_ARGS` -- POST form arguments (selector picks one).
- `REQUEST_URI` -- The request URI (path + query).
- `REQUEST_HEADERS` -- Request headers (selector picks one).
- `REQUEST_BODY` -- The request body.
- `REQUEST_COOKIES` -- Request cookies (selector picks one).

### spec.customRules[].matchConditions[].matchVariables[].selector

`string`

For the keyed variables (headers, cookies, args): the specific key to
inspect (e.g. "User-Agent" with REQUEST_HEADERS). Omit to inspect the
whole collection.

### spec.customRules[].matchConditions[].operator

`enum` · required

How the variable is compared against match_values.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_match_operator_unspecified` -- Not specified -- invalid; declare the operator.
- `ANY` -- Matches every request carrying the variable (no match_values).
- `IP_MATCH` -- The variable is an IP inside one of the match_values CIDRs/addresses (REMOTE_ADDR / XFF).
- `GEO_MATCH` -- The client's country code matches (two-letter ISO codes in match_values).
- `EQUAL` -- Exact string equality.
- `CONTAINS` -- Substring containment.
- `LESS_THAN` -- Numeric less-than.
- `GREATER_THAN` -- Numeric greater-than.
- `LESS_THAN_OR_EQUAL` -- Numeric less-than-or-equal.
- `GREATER_THAN_OR_EQUAL` -- Numeric greater-than-or-equal.
- `BEGINS_WITH` -- String prefix match.
- `ENDS_WITH` -- String suffix match.
- `REGEX` -- Regular-expression match.

### spec.customRules[].matchConditions[].matchValues

`[]string`

The values compared against (OR semantics). Required for every
operator except ANY, which matches all requests carrying the
variable.

### spec.customRules[].matchConditions[].negationCondition

`bool`

Inverts the result: the condition matches when the comparison does
NOT hold. Azure's default is false.

### spec.customRules[].matchConditions[].transforms

`[]enum`

Normalizations applied to the variable before comparison (e.g.
LOWERCASE + URL_DECODE to catch encoding evasions).

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_transform_unspecified` -- Not specified -- invalid; list only real transforms.
- `HTML_ENTITY_DECODE` -- Decode HTML entities (&lt; -> <).
- `LOWERCASE` -- Lowercase the value.
- `REMOVE_NULLS` -- Strip NULL bytes.
- `TRIM` -- Trim leading/trailing whitespace.
- `URL_DECODE` -- URL-decode the value (%27 -> ').
- `URL_ENCODE` -- URL-encode the value.
- `UPPERCASE` -- Uppercase the value.

### spec.managedRules

`AzureWebApplicationFirewallPolicyManagedRules` · required

The managed (Microsoft-curated) rule configuration. Required -- a WAF
policy without a managed rule set is rejected by Azure; the OWASP 3.2
core rule set is the baseline virtually every policy carries.

- rule: {"required":true}

### spec.managedRules.exclusions

`[]AzureWebApplicationFirewallPolicyManagedRulesExclusion`

Scoped exclusions: request parts (by key) that the managed rules skip
-- the tool for false positives on specific cookies, headers, or form
fields, without disabling whole rules for all traffic.

### spec.managedRules.exclusions[].matchVariable

`enum` · required

The keyed request collection the exclusion applies to (header names,
cookie values, query-arg names, ...).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_exclusion_match_variable_unspecified` -- Not specified -- invalid; declare the collection.
- `REQUEST_ARG_KEYS` -- Query/POST argument keys.
- `REQUEST_ARG_NAMES` -- Query/POST argument names.
- `REQUEST_ARG_VALUES` -- Query/POST argument values.
- `REQUEST_COOKIE_KEYS` -- Cookie keys.
- `REQUEST_COOKIE_NAMES` -- Cookie names.
- `REQUEST_COOKIE_VALUES` -- Cookie values.
- `REQUEST_HEADER_KEYS` -- Header keys.
- `REQUEST_HEADER_NAMES` -- Header names.
- `REQUEST_HEADER_VALUES` -- Header values.

### spec.managedRules.exclusions[].selectorMatchOperator

`enum` · required

How selector picks keys out of the collection.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_selector_match_operator_unspecified` -- Not specified: SELECTOR_EQUALS behavior for log scrubbing; invalid for exclusions, which declare an operator explicitly.
- `SELECTOR_EQUALS` -- The key equals the selector exactly.
- `SELECTOR_CONTAINS` -- The key contains the selector.
- `SELECTOR_STARTS_WITH` -- The key starts with the selector.
- `SELECTOR_ENDS_WITH` -- The key ends with the selector.
- `SELECTOR_EQUALS_ANY` -- Every key in the collection matches (no selector).

### spec.managedRules.exclusions[].selector

`string` · required

The key (or key fragment, per the operator) to exclude -- e.g. a
session cookie name that trips SQL-injection rules.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules.exclusions[].excludedRuleSet

`AzureWebApplicationFirewallPolicyExcludedRuleSet`

Narrow the exclusion to specific rule sets, groups, or rules instead
of all managed rules -- the surgical form.

### spec.managedRules.exclusions[].excludedRuleSet.type

`enum`

The rule-set family the exclusion narrows to. Unspecified applies
OWASP (azurerm's default).

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_managed_rule_set_type_unspecified` -- Not specified: OWASP (azurerm's default).
- `OWASP` -- The OWASP core rule set (versions 2.2.9-3.2) -- the standard SQLI/XSS/RCE/LFI protection baseline.
- `MICROSOFT_BOT_MANAGER_RULE_SET` -- Microsoft's bot manager rule set (versions 0.1-1.1): good/bad/unknown bot classification.
- `MICROSOFT_DEFAULT_RULE_SET` -- Microsoft's default rule set (versions 2.1-2.2): the OWASP successor with Microsoft threat-intelligence rules.

### spec.managedRules.exclusions[].excludedRuleSet.version

`string` · optional (explicit presence)

The rule-set version. Unspecified applies "3.2" (azurerm's default).
Note the narrower allowed set than managed_rule_sets[].version --
exclusions cannot target OWASP 3.0/3.1/2.2.9.

- default: `3.2`
- rule: version must be one of: 1.0, 1.1, 2.1, 2.2, 3.2

### spec.managedRules.exclusions[].excludedRuleSet.ruleGroups

`[]AzureWebApplicationFirewallPolicyExcludedRuleGroup`

The rule groups (and optionally individual rules) the exclusion
narrows to.

### spec.managedRules.exclusions[].excludedRuleSet.ruleGroups[].ruleGroupName

`string` · required

The rule group, using Azure's exact group names.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules.exclusions[].excludedRuleSet.ruleGroups[].excludedRules

`[]string`

Specific rule IDs within the group. Omit to exclude the whole group.

### spec.managedRules.managedRuleSets

`[]AzureWebApplicationFirewallPolicyManagedRuleSet` · required

The rule sets to run. Typically one OWASP set, optionally alongside
the Microsoft bot-manager set.

- rule: {"repeated":{"minItems":"1"}}

### spec.managedRules.managedRuleSets[].type

`enum`

The rule-set family. Unspecified applies OWASP (azurerm's default).

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_managed_rule_set_type_unspecified` -- Not specified: OWASP (azurerm's default).
- `OWASP` -- The OWASP core rule set (versions 2.2.9-3.2) -- the standard SQLI/XSS/RCE/LFI protection baseline.
- `MICROSOFT_BOT_MANAGER_RULE_SET` -- Microsoft's bot manager rule set (versions 0.1-1.1): good/bad/unknown bot classification.
- `MICROSOFT_DEFAULT_RULE_SET` -- Microsoft's default rule set (versions 2.1-2.2): the OWASP successor with Microsoft threat-intelligence rules.

### spec.managedRules.managedRuleSets[].version

`string` · required

The rule-set version. OWASP: "3.2", "3.1", "3.0", "2.2.9";
Microsoft_DefaultRuleSet: "2.1", "2.2"; Microsoft_BotManagerRuleSet:
"0.1", "1.0", "1.1". OWASP 3.2 is the current production standard --
several policy dials (file-upload enforcement) only work with it.

- rule: version must be one of: 0.1, 1.0, 1.1, 2.1, 2.2, 2.2.9, 3.0, 3.1, 3.2
- rule: {"required":true}

### spec.managedRules.managedRuleSets[].ruleGroupOverrides

`[]AzureWebApplicationFirewallPolicyRuleGroupOverride`

Per-rule-group tuning: disable individual rules or change their
action without turning off the whole set.

### spec.managedRules.managedRuleSets[].ruleGroupOverrides[].ruleGroupName

`string` · required

The rule group, using Azure's exact group names (e.g.
"REQUEST-942-APPLICATION-ATTACK-SQLI" for OWASP 3.x,
"crs_41_sql_injection_attacks" for 2.2.9, "SQLI" for
Microsoft_DefaultRuleSet, "BadBots" for the bot manager).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules.managedRuleSets[].ruleGroupOverrides[].rules

`[]AzureWebApplicationFirewallPolicyRuleOverride` · required

The individual rules to tune within the group.

- rule: {"repeated":{"minItems":"1"}}

### spec.managedRules.managedRuleSets[].ruleGroupOverrides[].rules[].id

`string` · required

The rule ID, from the rule set's documentation (e.g. "942100" for the
OWASP SQL-injection libinjection rule).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedRules.managedRuleSets[].ruleGroupOverrides[].rules[].enabled

`bool` · optional (explicit presence)

Whether the rule runs. azurerm's default is false -- listing a rule
here without enabled=true DISABLES it, which is the common tuning
gesture.

- default: `false`

### spec.managedRules.managedRuleSets[].ruleGroupOverrides[].rules[].action

`enum`

Replace the rule's action (e.g. downgrade BLOCK to OVERRIDE_LOG while
investigating false positives). Omit to keep the rule set's action.

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_rule_override_action_unspecified` -- Not specified: keep the rule set's own action.
- `OVERRIDE_ALLOW` -- Let matching requests through.
- `OVERRIDE_ANOMALY_SCORING` -- Contribute to the anomaly score instead of acting directly (the OWASP 3.x scoring model).
- `OVERRIDE_BLOCK` -- Reject matching requests.
- `OVERRIDE_JS_CHALLENGE` -- Serve a JavaScript challenge.
- `OVERRIDE_LOG` -- Record matches without acting.

### spec.policySettings

`AzureWebApplicationFirewallPolicySettings`

Enforcement mode and body-inspection dials. Omit for Azure's
defaults: enabled, Prevention mode, request-body inspection on with a
128 KB limit, 100 MB file-upload limit.

### spec.policySettings.enabled

`bool` · optional (explicit presence)

Whether the policy is enforced at all. Azure's default is true;
false parks the policy without detaching it.

- default: `true`

### spec.policySettings.mode

`enum`

PREVENTION blocks matching requests (Azure's default and the
production posture); DETECTION only logs them -- the tuning mode for
new policies watching real traffic.

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_mode_unspecified` -- Not specified: PREVENTION (Azure's default).
- `PREVENTION` -- Block requests that trip the rules -- the production posture.
- `DETECTION` -- Log matches without blocking -- the tuning mode.

### spec.policySettings.requestBodyCheck

`bool` · optional (explicit presence)

Whether request bodies are inspected at all. Azure's default is
true; turning it off blinds the WAF to body-borne attacks.

- default: `true`

### spec.policySettings.requestBodyEnforcement

`bool` · optional (explicit presence)

Whether requests with bodies LARGER than max_request_body_size_in_kb
are blocked (true, Azure's default) or passed through with only the
first bytes inspected (false).

- default: `true`

### spec.policySettings.requestBodyInspectLimitInKb

`int32` · optional (explicit presence)

How many KB of each request body the WAF inspects. Azure's default
is 128. 0 means unlimited inspection.

- default: `128`
- rule: {"int32":{"gte":0}}

### spec.policySettings.maxRequestBodySizeInKb

`int32` · optional (explicit presence)

The maximum request body size in KB before request_body_enforcement
applies. 8-2000; Azure's default is 128.

- default: `128`
- rule: {"int32":{"lte":2000,"gte":8}}

### spec.policySettings.fileUploadEnforcement

`bool` · optional (explicit presence)

Whether uploads larger than file_upload_limit_in_mb are blocked.
Only honored with OWASP 3.2. Azure's default is true.

- default: `true`

### spec.policySettings.fileUploadLimitInMb

`int32` · optional (explicit presence)

The maximum file upload size in MB. 1-4000; Azure's default is 100.

- default: `100`
- rule: {"int32":{"lte":4000,"gte":1}}

### spec.policySettings.jsChallengeCookieExpirationInMinutes

`int32` · optional (explicit presence)

How long a solved JavaScript challenge stays valid before the client
is re-challenged, in minutes. 5-1440; Azure's default is 30. Only
meaningful when some rule uses the JS_CHALLENGE action.

- default: `30`
- rule: {"int32":{"lte":1440,"gte":5}}

### spec.policySettings.logScrubbing

`AzureWebApplicationFirewallPolicyLogScrubbing`

Log scrubbing: redact sensitive request parts (auth headers, PII
arguments) from WAF logs before they land in Log Analytics.

### spec.policySettings.logScrubbing.enabled

`bool` · optional (explicit presence)

Whether scrubbing is active. Azure's default is true (when the block
is present).

- default: `true`

### spec.policySettings.logScrubbing.rules

`[]AzureWebApplicationFirewallPolicyLogScrubbingRule` · required

The request parts to redact from logs.

- rule: {"repeated":{"minItems":"1"}}
- rule: selector is required with SELECTOR_EQUALS (except for SCRUB_REQUEST_IP_ADDRESS) and must be omitted with SELECTOR_EQUALS_ANY

### spec.policySettings.logScrubbing.rules[].enabled

`bool` · optional (explicit presence)

Whether this rule is active. Azure's default is true.

- default: `true`

### spec.policySettings.logScrubbing.rules[].matchVariable

`enum` · required

The request part to redact.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_scrubbing_match_variable_unspecified` -- Not specified -- invalid; declare the part to redact.
- `SCRUB_REQUEST_ARG_NAMES` -- Query/POST argument names.
- `SCRUB_REQUEST_COOKIE_NAMES` -- Cookie names.
- `SCRUB_REQUEST_HEADER_NAMES` -- Header names.
- `SCRUB_REQUEST_IP_ADDRESS` -- The client IP address (no selector).
- `SCRUB_REQUEST_JSON_ARG_NAMES` -- JSON body argument names.
- `SCRUB_REQUEST_POST_ARG_NAMES` -- POST argument names.

### spec.policySettings.logScrubbing.rules[].selectorMatchOperator

`enum`

How selector picks keys: SELECTOR_EQUALS redacts one named key (the
default), SELECTOR_EQUALS_ANY redacts every key in the collection.

- rule: log-scrubbing rules support only SELECTOR_EQUALS and SELECTOR_EQUALS_ANY

Allowed values (use exactly as shown):

- `azure_web_application_firewall_policy_selector_match_operator_unspecified` -- Not specified: SELECTOR_EQUALS behavior for log scrubbing; invalid for exclusions, which declare an operator explicitly.
- `SELECTOR_EQUALS` -- The key equals the selector exactly.
- `SELECTOR_CONTAINS` -- The key contains the selector.
- `SELECTOR_STARTS_WITH` -- The key starts with the selector.
- `SELECTOR_ENDS_WITH` -- The key ends with the selector.
- `SELECTOR_EQUALS_ANY` -- Every key in the collection matches (no selector).

### spec.policySettings.logScrubbing.rules[].selector

`string`

The key to redact (e.g. "Authorization" with
SCRUB_REQUEST_HEADER_NAMES). Required with SELECTOR_EQUALS; must be
omitted with SELECTOR_EQUALS_ANY. REQUEST_IP_ADDRESS takes no
selector.

### spec.tags

`map<string, string>`

Free-form tags applied to the policy, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureWebApplicationFirewallPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.policy_id` | `string` | The Azure Resource Manager ID of the WAF policy. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/applicationGatewayWebApplicationFirewallPolicies/{name} Referenced by AzureApplicationGateway (firewall_policy_id, per-listener firewall_policy_id, and per-path-rule firewall_policy_id). |
| `status.outputs.policy_name` | `string` | The name of the WAF policy. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureApplicationGateway | `spec.httpListeners[].firewallPolicyId` | `status.outputs.policy_id` |
| AzureApplicationGateway | `spec.urlPathMaps[].pathRules[].firewallPolicyId` | `status.outputs.policy_id` |
| AzureApplicationGateway | `spec.firewallPolicyId` | `status.outputs.policy_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
