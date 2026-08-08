# AzureFrontDoorRuleSet -- Terraform Module

Creates an `azurerm_cdn_frontdoor_rule_set` on the referenced Front Door
profile plus one `azurerm_cdn_frontdoor_rule` per entry in
`spec.rules`, keyed by the rule's spec-unique name.

## Inputs

- `metadata` -- Planton resource metadata (name, org, env, labels)
- `spec` -- see `variables.tf`; mirrors `spec.proto` exactly. Enum
  fields arrive as the spec enum's FULL value names (e.g.
  `PERMANENT_REDIRECT`, `OVERRIDE_ALWAYS`) and are mapped to ARM's
  casing in `locals.tf`.

## Behavior notes

- **Evaluation order is the rule's own `order` field** -- resource
  creation order carries no meaning. `order` materializes its
  documented 0 default because tfvars drops zero values.
- **Two protocol dialects**: the shared forwarding-protocol enum maps
  to `Http`/`Https` on the redirect action but `HttpOnly`/`HttpsOnly`
  on the route-configuration override -- both maps live in `locals.tf`
  with the reason.
- **Negation is folded into the operator value** -- the provider has no
  `negate_condition` argument; the spec's per-condition boolean becomes
  a `Not` prefix on the operator (`NotEqual`, `NotIPMatch`, ...), and
  the equality-only conditions (method, scheme, HTTP version, device,
  TLS) always emit an explicit `Equal`/`NotEqual` because the provider
  requires an operator on every condition.
- **Address-condition operators default to `IPMatch`** when unspecified
  (the provider's documented default, materialized because tfvars drops
  unset fields).
- **Header DELETE actions send `null` values** -- the provider rejects
  an empty string on Append/Overwrite and any value on Delete; the
  spec's CEL guarantees the right combination arrives.
- **No Azure tags** -- ARM does not support tags on rule sets or rules.

## Outputs

- `rule_set_id` -- what routes reference in `rule_set_ids`
- `rule_set_name`
