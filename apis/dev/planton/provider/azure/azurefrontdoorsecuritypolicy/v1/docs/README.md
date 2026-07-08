# AzureFrontDoorSecurityPolicy -- Design Research

## Scope

The security policy is the WAF-to-domain association -- an ARM child of
a Front Door profile (`Microsoft.Cdn/profiles/{p}/securityPolicies/{n}`)
that binds one Front Door WAF policy to a set of the profile's
hostnames. It is a first-class kind because it is the enforcement seam
(a WAF without one does nothing) and because one profile carries many
associations with independent lifecycles.

Source of truth: `azurerm_cdn_frontdoor_security_policy`
(terraform-provider-azurerm v4.80, `internal/services/cdn/`, CDN API
2024-02-01), parity-verified against pulumi-azure v6
(`cdn.FrontdoorSecurityPolicy`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `profile_id` | `cdn_frontdoor_profile_id` | FK to AzureFrontDoorProfile; ForceNew |
| `security_policy_name` | `name` | ForceNew; begins/ends alphanumeric, letters/digits/hyphens (the provider's own regex; no length bound beyond the regex) |
| `firewall_policy_id` | `security_policies.firewall.cdn_frontdoor_firewall_policy_id` | FK to AzureFrontDoorFirewallPolicy; ForceNew |
| `domain_ids[]` | `...association.domain[].cdn_frontdoor_domain_id` | 1-500; each an endpoint OR custom-domain ARM id (the provider's dual-ID validator accepts both); updatable in place |
| -- | `...association.patterns_to_match` | a CONSTANT, not a field (see below) |

## Shape decisions

- **The provider's triple nesting is flattened.** azurerm wraps the
  parameters in `security_policies` (1 item) -> `firewall` (1 item) ->
  `association` (1 item) -- TF ergonomics for a one-choice ARM union
  (`SecurityPolicyType` has exactly one value, WebApplicationFirewall).
  The spec models the four real inputs flat; each module rebuilds the
  wrapper shape it needs.
- **One association, not a list.** ARM's model carries an associations
  ARRAY, but the provider caps it at one -- and because every
  association must use the same `/*` pattern, multiple associations
  add nothing today. If Azure ever ships real per-path scoping, a
  repeated shape lands with it.
- **`patterns_to_match` is a constant.** The service accepts exactly
  `/*` (the provider validates in-list with one value). A one-value
  knob is a constant, not configuration: both modules send it
  unconditionally. Re-enable trigger: Azure accepting a second
  pattern. Engine dialect: the pulumi bridge flattens the one-item
  list to a singular string -- TF sends `["/*"]`, Pulumi `"/*"`; same
  ARM payload, documented at both sites.
- **`default_kind` points at the endpoint.** The domain list accepts
  two kinds; the endpoint is the default (every profile has one; the
  default-domain association is the first one every deployment makes)
  with the custom-domain alternative documented in the field comment
  and preset.

## Validation contracts

- Name regex, min/max domain counts (1-500) -- schema-level.
- **The 100/500 domain cap rides the PROFILE's sku** -- the provider
  reads the profile at apply time to pick the cap; cross-resource, so
  it stays apply-time (documented in the spec).
- **The WAF-policy-sku-must-match-the-profile-sku pairing** -- enforced
  by ARM at association time; cross-resource, documented in both this
  kind and the WAF kind.

## Recorded skips

- **The per-domain `active` computed flag** (`domain[].active`): a
  read-back status of whether the association is live; nothing
  composes on it and it carries no configuration. Operational
  observability belongs to diagnostics, not stack outputs.

## E2E note

The live proof associates the WAF with the ENDPOINT's default domain:
associating a validated custom domain requires a delegated DNS zone
the test subscription does not own. The endpoint-scoped association is
a real association (the provider's domain validator accepts endpoint
IDs as first-class), not a test shortcut.
