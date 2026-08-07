# AzureWebApplicationFirewallPolicy -- Design Research

## The Resource

A regional Web Application Firewall policy
(`Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies`) is
the rule set an Azure Application Gateway enforces on HTTP traffic. The
component maps onto `azurerm_web_application_firewall_policy` (azurerm
v4.x, `internal/services/network/web_application_firewall_policy_resource.go`),
parity-verified against pulumi-azure v6 (`waf.Policy`).

Azure Front Door's WAF policy is a DIFFERENT ARM resource
(`cdn_frontdoor_firewall_policy`) with its own rule vocabulary; it belongs
to the Front Door family, never here.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `policy_name` | Required, ForceNew |
| `location` | `region` | Required, ForceNew; must match the attaching gateways' region |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `custom_rules` | `custom_rules` | Full surface: match + rate-limit types, priorities 1-100 (ARM's real range; azurerm has no range check), enabled staging, JSChallenge action, XFF/geo rate grouping |
| `custom_rules[].match_conditions` | same | Variable enum (8), operator enum (12), transforms enum (7), negation; ANY-operator ⊕ match_values CEL mirrors ARM |
| `managed_rules` | `managed_rules` | Required (Max1 in azurerm → a message field) |
| `managed_rules[].managed_rule_set` | `managed_rule_sets[]` | Type enum (OWASP/BotManager/DefaultRuleSet, unspecified = OWASP like azurerm's default), version vocabulary CEL (9 values) |
| `rule_group_override.rule[]` | `rule_group_overrides[].rules[]` | id + enabled (default FALSE -- listing = disabling, azurerm's contract) + action enum (OVERRIDE_* prefix avoids colliding with custom-rule actions) |
| `managed_rules[].exclusion` | `exclusions[]` | 9-value collection enum, 5-value selector operator enum (SELECTOR_* prefix), required selector, nested excluded_rule_set with its NARROWER version set (no 2.2.9/3.0/3.1) |
| `policy_settings` | `policy_settings` | All 9 dials + log_scrubbing (6-value scrub-variable enum, Equals/EqualsAny operator subset via CEL) |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `http_listener_ids` / `path_based_rule_ids` (computed) | -- (skipped) | Reverse-link conveniences; the resource graph already models attachment from the gateway side, and echoing it here would create a second, staleable source of truth |

## Decomposition Decisions

- **The policy is a first-class kind, not a gateway sub-message.** It has
  an independent lifecycle (tuned continuously while gateways stay
  untouched), is FK-referenced at three levels (gateway, listener, path
  rule), and is shared across gateways -- all three split criteria.
- **Everything inside the policy folds.** Custom rules, managed-rule
  tuning, and settings have no life outside their policy and are ARM
  sub-properties, not resources.

## Recorded Skips (with reasons)

- **`rule_type: Invalid`** -- azurerm passes ARM's enum through verbatim,
  which includes a nonsense "Invalid" member; modeling it would let users
  declare a rule that can never work.
- **`http_listener_ids` / `path_based_rule_ids` outputs** -- see the field
  map: reverse links duplicated from the gateway's spec would go stale the
  moment a gateway detaches. The resource graph is the source of truth for
  attachment.
- **`rule_group_name` as a string, not an enum** -- the 43-value
  vocabulary mixes hyphens, underscores, and casing that cannot survive
  proto enum naming (`REQUEST-942-...`, `crs_41_...`, `Known-CVEs`), and
  it grows with every rule-set release. The field comment teaches the
  format per rule set; azurerm validates the value at plan time.

## Design Decisions

- **Custom-rule priority bounded 1-100** -- ARM's real constraint for
  custom WAF rules; azurerm's schema carries no range check, so the spec
  supplies the guard the provider forgot.
- **Rate-limit trio CEL-paired to RATE_LIMIT_RULE** and **ALLOW rejected
  on rate-limit rules** -- both are ARM's runtime contract, surfaced at
  validation time instead of after a deploy.
- **Enum prefixes where vocabularies collide**: `OVERRIDE_*` (override
  actions vs custom-rule actions), `SELECTOR_*` (selector operators vs
  match operators), `SCRUB_*` (scrubbing variables vs exclusion
  variables) -- proto enum values share a namespace per package, and the
  three unprefixed vocabularies are the ones users type most.
- **The exclusion `excluded_rule_set.version` has its own narrower CEL**
  (1.0/1.1/2.1/2.2/3.2) -- azurerm validates the two version fields
  against different sets; collapsing them would accept versions ARM
  rejects.
- **`file_upload_enforcement` is forwarded only on explicit presence** on
  both engines -- it is only honored with OWASP 3.2, and materializing a
  default would error on older rule sets.

## Operational Behavior Worth Knowing

- **Everything except name/region/resource-group updates in place** --
  tuning never recreates the policy or detaches gateways.
- **Detection first, then Prevention** is the standard rollout: run the
  policy in DETECTION against real traffic, review matches, add
  exclusions/overrides, then flip the mode.
- **A policy and its gateways must share a region**; attaching a policy
  from another region fails at the gateway.
- **Deleting a policy still referenced by a gateway fails** -- detach (or
  destroy the gateway) first; the resource graph's dependency ordering
  handles this in composed environments.
- **JS challenge requires WAF_v2** (as does the policy attachment itself);
  the challenge cookie lifetime is a policy-wide dial.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `policy_id` output is consumed by `AzureApplicationGateway`:
  - `firewall_policy_id` (gateway-wide)
  - `http_listeners[].firewall_policy_id` (per listener)
  - `url_path_map[].path_rules[].firewall_policy_id` (per route)
