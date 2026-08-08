# AzureFrontDoorFirewallPolicy -- Terraform Module

Creates an `azurerm_cdn_frontdoor_firewall_policy` (azurerm ~> 5.0) in
the referenced resource group. Credentials arrive as ARM_* environment
variables (service principal or keyless OIDC); the provider block stays
empty.

## Behavior notes

- **Global resource** -- ARM fixes the location; the provider sends no
  region.
- **Challenge lifetimes are Premium-conditional**: on Premium the
  module always sends the JS-challenge and CAPTCHA lifetimes (spec
  value or the documented default of 30) because Azure enables both
  policies there unconditionally; on Standard `null` is sent (Azure
  rejects the fields -- the spec's CELs keep values out).
- **Provider defaults pass through**: null lets enabled (true),
  request-body check (true), custom-rule priority (1), the rate-limit
  pair (1/10), and the override rule's enabled (FALSE -- listing a
  rule is the disable gesture) fall to azurerm's defaults.
- **Enum vocabulary maps live in `locals.tf`**: the spec's prefixed
  value names (`RULE_SET_*`, `OVERRIDE_*`, `SELECTOR_*`, `EXCLUDE_*`,
  `SCRUB_*`) translate to ARM's bare wire values.

## Outputs

- `firewall_policy_id` -- what AzureFrontDoorSecurityPolicy references
  in `firewall_policy_id`
- `firewall_policy_name`

## Validate offline

```shell
planton tofu plan --manifest ../../e2e/manifest.yaml --module-dir .
```
