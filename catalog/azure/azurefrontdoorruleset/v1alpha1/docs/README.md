# AzureFrontDoorRuleSet -- Design Research

## Scope

The rule set is Front Door's edge delivery policy: an ordered list of
rules, each pairing match conditions with actions, that routes attach
by ARM id. It is a first-class kind because one policy is shared across
many routes and has its own lifecycle; the RULES fold inside it because
they form one ordered document -- nothing references an individual
rule, evaluation order is a property of the set, and a rule has no life
outside it.

Source of truth: `azurerm_cdn_frontdoor_rule_set` +
`azurerm_cdn_frontdoor_rule` (terraform-provider-azurerm v4.80,
`internal/services/cdn/`), parity-verified against pulumi-azure v6
(`cdn.FrontdoorRuleSet` / `cdn.FrontdoorRule`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `profile_id` | `cdn_frontdoor_profile_id` | FK to AzureFrontDoorProfile; ForceNew |
| `rule_set_name` | rule set `name` | ForceNew; 1-60 letters/digits, NO hyphens (the provider's own regex -- stricter than other Front Door names) |
| `rules[]` | one `azurerm_cdn_frontdoor_rule` each | keyed by rule `name` (the ARM child identity; spec CEL enforces uniqueness); rule name ForceNew, everything else updates in place |
| `rules[].order` | `order` | plain int >= 0; the provider enforces nothing else (no uniqueness) |
| `rules[].behavior_on_match` | `behavior_on_match` | Continue (default) / Stop |
| `rules[].conditions.*` | `conditions` block | all 19 provider condition types modeled; see contracts below |
| `rules[].actions.*` | `actions` block | all 5 provider action types; redirect/rewrite/override modeled as SINGULAR fields so the provider's only-once contracts are structural |

## Shape decisions

- **Conditions/actions as typed lists in one message, not ARM's
  heterogeneous array** -- conditions AND together and actions all
  apply, so cross-type ordering carries no semantics; the typed shape is
  what both engines' wire formats speak, and it keeps every condition's
  own vocabulary independently validatable.
- **One shared operator enum (13 values)** with per-condition `in`
  subsets: the standard 10 comparisons everywhere, `WILDCARD` only on
  url_path, `GEO_MATCH`/`IP_MATCH` only on the address conditions.
  Unspecified is legal ONLY on the address conditions, where it means
  IP_MATCH (the provider default, materialized in both modules).
- **Equal-only conditions carry no operator field** (request_method,
  request_scheme, http_version, is_device, ssl_protocol): the provider
  accepts exactly one operator for them, and a one-value knob is a
  constant, not a choice. Both modules let the provider default apply.
- **Closed string vocabularies stay strings with `in`-list CELs**
  (methods, HTTP versions, TLS versions, server ports, device classes,
  schemes) -- they are ARM's own wire values (`"2.0"`, `"TLSv1.2"`,
  `"443"`), following the route kind's content-type-allowlist pattern.
- **One shared forwarding-protocol enum, two wire dialects**: the
  redirect action speaks `Http`/`Https`/`MatchRequest` while the
  route-configuration override speaks `HttpOnly`/`HttpsOnly`/
  `MatchRequest`. Same semantics; each module maps its own dialect
  (documented at both map sites).
- **`cache_behavior` is required on the override**: the provider
  rejects an override without it, and ARM round-trips an absent cache
  configuration as Disabled -- every override makes an explicit cache
  decision.

## Validation contracts (front-loaded as CELs)

Mirrored from the provider's expand-time validators (there is no
CustomizeDiff on either resource):

- <= 10 conditions and <= 5 actions per rule; at least one action
- redirect XOR rewrite; redirect/rewrite/override at most once
  (structural -- singular fields)
- operator `ANY` <-> empty match_values (per condition); the two
  conditions whose match_values the provider REQUIRES (request_body,
  url_file_extension) exclude ANY from their operator vocabulary --
  the conjunction of the provider's two rules makes ANY unusable there
- GEO_MATCH values are two-letter uppercase country codes
- header DELETE carries no value; APPEND/OVERWRITE require one
- the override's ~10 cache/forwarding pairings: origin <-> protocol
  both-or-neither, DISABLED forbids the caching fields, caching
  behaviors require query_string_caching_behavior, HONOR_ORIGIN
  forbids cache_duration while OVERRIDE_* require it, *_SPECIFIED
  behaviors require query_string_parameters (and the other two forbid
  them), cache_duration in d.HH:MM:SS (days 1-365, never "0.")
- rule names unique within the set (fold-derived: each rule is an ARM
  child resource keyed by name -- duplicates would silently target one
  resource)

Deliberately NOT validated (the provider enforces neither): rule-order
uniqueness; CIDR validity/duplicates/overlap on the address conditions
(the provider checks these apply-time with `net.ParseCIDR` -- regex
cannot honestly replicate CIDR parsing, and inventing a partial rule
would reject valid input or admit invalid).

## Apply-time contracts left to Azure (documented, not CELs)

- CIDR validity, duplicate, and overlap checks on
  remote_address/socket_address (provider expand-time, exact parsing)
- Premium-only surface: the RegEx operator (and Wildcard in some
  regions), ssl_protocol/socket_address/server_port/client_port
  conditions, and the route-configuration override are rejected by ARM
  on Standard profiles -- the provider has NO client-side SKU gate, so
  the spec does not invent one (cross-resource: the SKU lives on the
  profile)

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `rule_set_id` | rule set resource id | `AzureFrontDoorRoute.rule_set_ids` |
| `rule_set_name` | `name` | human/portal orientation |

No per-rule ids on purpose: nothing references a rule -- routes attach
the whole set.

## Recorded skips (with reasons)

- **Rule-order uniqueness / gap validation** -- the provider enforces
  only `order >= 0`; anything stricter would be an invented rule.
- **Per-rule ARM-id outputs** -- no consumer exists; the set id is the
  only composition seam.

## Lifecycle notes

- The rule set has NO update in the provider -- both of its fields are
  ForceNew; renaming the set replaces every rule under it.
- Rules update in place (order, behavior, conditions, actions); only a
  rule's name and its parent are ForceNew.
- The provider polls rule deletion to true absence (the API returns
  success early); both engines inherit the behavior from the provider
  layer.
