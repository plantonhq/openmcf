# AzureFrontDoorFirewallPolicy -- Design Research

## Scope

The Front Door WAF policy is the edge firewall rule set: policy
settings (mode, body inspection, block response, challenge lifetimes),
custom match/rate-limit rules, Microsoft's managed rule sets, and log
scrubbing. It is a first-class kind because one policy is shared across
many profiles (via security policies) and carries its own lifecycle.

It is a DIFFERENT ARM type than the regional WAF policy
(`Microsoft.Network/frontDoorWebApplicationFirewallPolicies`, global,
vs `ApplicationGatewayWebApplicationFirewallPolicies`, regional) with a
different rule vocabulary -- the two kinds deliberately do not share
shapes.

Source of truth: `azurerm_cdn_frontdoor_firewall_policy`
(terraform-provider-azurerm v4.80, `internal/services/cdn/`, backed by
the frontdoor 2025-03-01 ARM API), parity-verified against pulumi-azure
v6 (`cdn.FrontdoorFirewallPolicy`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `resource_group` | `resource_group_name` | FK to AzureResourceGroup; ForceNew |
| `policy_name` | `name` | ForceNew; 1-128, begins with a letter, letters/digits only (the provider's own regex -- NO hyphens) |
| `sku` | `sku_name` | ForceNew; Standard_AzureFrontDoor / Premium_AzureFrontDoor; unspecified deploys STANDARD; PREMIUM->STANDARD rejected outright (CustomizeDiff) |
| `mode` | `mode` | Detection / Prevention; REQUIRED (no provider default, unlike the regional WAF) |
| `enabled` | `enabled` | default true |
| `request_body_check_enabled` | `request_body_check_enabled` | default true |
| `redirect_url` | `redirect_url` | http/https scheme validated; ARM requires it when any rule action is Redirect (apply-time, cross-field) |
| `custom_block_response_status_code` | same | closed 15-value set (200/403/405/406/429/990-999) |
| `custom_block_response_body` | same | base64-validated |
| `js_challenge_cookie_expiration_in_minutes` | same | 5-1440; PREMIUM only; Azure always enables the policy there (default 30) |
| `captcha_cookie_expiration_in_minutes` | same | 5-1440; PREMIUM only; same always-on semantics |
| `custom_rules[]` | `custom_rule` | <= 100; priority default 1, enabled default true, rate-limit window/threshold defaults 1/10 (always present in the provider; ARM ignores them on MatchRule) |
| `custom_rules[].match_conditions[]` | `match_condition` | <= 10 per rule; match_values REQUIRED (1-600 x 1-256 chars); selector only for the keyed variables |
| `managed_rules[]` | `managed_rule` | <= 100; PREMIUM only (create-path validator); type/version open strings with pairing checks |
| `managed_rules[].exclusions[]` / `overrides[]` | `exclusion` / `override` | exclusions <= 100 at all three scopes (set/group/rule); overrides <= 100, rules <= 1000 |
| `log_scrubbing` | `log_scrubbing` | enabled default true; <= 100 rules; 7 match variables (a superset of the profile's 3) |
| `tags` | `tags` | user tags merged over Planton-derived tags |

## Shape decisions

- **No `region` field** -- the resource is global; the provider
  hardcodes location `Global`. (The regional WAF kind has a region;
  copying it here would be inventing a knob ARM does not have.)
- **`match_values` is genuinely required** (min 1) -- this provider
  marks it Required even for the `ANY` operator, unlike the regional
  WAF's "ANY means empty values" contract. Modeled provider-exactly;
  no invented ANY rule.
- **Managed-rule `type`/`version` are documented strings, not enums** --
  the provider deliberately validates them as non-empty strings so
  Microsoft can ship new set types and versions server-side; only the
  two statically-known version pairings are CEL-enforced. (The
  regional WAF's closed type enum was faithful to ITS provider
  validator; here it would be invention.)
- **Prefixed enum value names** (`RULE_SET_*`, `OVERRIDE_*`,
  `SELECTOR_*`, `EXCLUDE_*`, `SCRUB_*`) exist only to keep proto enum
  names collision-free within the kind; both modules map them to ARM's
  bare wire values.
- **The challenge lifetimes are modeled WITHOUT platform defaults** --
  on PREMIUM, Azure always enables both policies and both engines pin
  the documented default of 30 when the spec is silent; on STANDARD the
  fields are rejected (spec CEL), so a platform-materialized default
  would break Standard deployments.
- **One shared selector-operator enum** serves managed-rule exclusions
  (all five values) and log scrubbing (Equals/EqualsAny only,
  CEL-restricted) -- same ARM vocabulary, one map per module.

## Validation contracts (front-loaded as CELs)

Mirrored from the provider's CustomizeDiff and create/expand
validators:

- **The three Standard-sku gates** (CustomizeDiff + create path):
  managed rules, the JS-challenge/CAPTCHA lifetimes, and the
  JS_CHALLENGE/CAPTCHA custom-rule actions are all PREMIUM-only.
- **DefaultRuleSet version pairing** (expand): the legacy type only at
  1.0/preview-0.1; Microsoft_DefaultRuleSet only at 1.1+.
- **Override action version gates** (expand): AnomalyScoring only on
  2.0+ sets; 2.0+ sets allow only AnomalyScoring/Log; JSChallenge only
  on bot-manager set types (matched case-insensitively, as the
  provider does).
- **Scrubbing-rule contracts** (expand): RequestIPAddress/RequestUri
  accept only EqualsAny; Equals requires a selector; EqualsAny forbids
  one.
- redirect_url scheme, block status code set, block body base64
  (schema validators).

Not statically checkable (documented, enforced by Azure at deploy
time): the sku-must-match-the-profile pairing (cross-resource, checked
at association), redirect_url required when any rule redirects
(cross-field against the rule list at ARM), and selector requiredness
per match variable (keyed variables need one; ARM rejects it
elsewhere).

## Recorded skips

- **ARM's `groupBy` rate-limit grouping** (GeoLocation/SocketAddr/None
  per custom rule, frontdoor API 2025-03-01): NOT in the azurerm
  schema, so neither engine can express it. Re-enable trigger: the
  field landing in `azurerm_cdn_frontdoor_firewall_policy`.
- **`frontend_endpoint_ids` computed attribute**: a read-back of
  classic Front Door frontend links -- not configuration and not a
  composition seam for Standard/Premium profiles (security policies
  own the association).
- **The provider's PREMIUM default-reset CustomizeDiff** (drifting the
  challenge lifetimes back to 30 when removed from config) is realized
  behaviorally: both engines always send value-or-30 on Premium, so
  the reset semantics hold declaratively.

## Preview-status fields

Azure flags the JS_CHALLENGE/CAPTCHA custom-rule actions and log
scrubbing as preview features; the spec comments carry that status.
