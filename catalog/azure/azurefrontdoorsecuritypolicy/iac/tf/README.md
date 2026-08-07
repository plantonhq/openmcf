# AzureFrontDoorSecurityPolicy -- Terraform Module

Creates an `azurerm_cdn_frontdoor_security_policy` (azurerm ~> 5.0) on
the referenced Front Door profile. Credentials arrive as ARM_*
environment variables (service principal or keyless OIDC); the provider
block stays empty.

## Behavior notes

- **The wrapper nesting is rebuilt here**: the spec's flat four fields
  become the provider's security_policies -> firewall -> association
  blocks (a one-choice ARM union).
- **`patterns_to_match` is the constant `["/*"]`** -- the service
  accepts no other pattern. Engine dialect: the pulumi bridge flattens
  this one-item list to a singular string; same ARM payload.
- **Domain references arrive resolved**: endpoint and custom-domain
  ARM ids are literals by the time the module runs; both types render
  as `domain` blocks unchanged.
- **No Azure tags** -- ARM does not support tags on security policies.

## Outputs

- `security_policy_id` -- operational addressing (nothing composes on
  the association itself)
- `security_policy_name`

## Validate offline

```shell
planton tofu plan --manifest ../../e2e/manifest.yaml --module-dir .
```
