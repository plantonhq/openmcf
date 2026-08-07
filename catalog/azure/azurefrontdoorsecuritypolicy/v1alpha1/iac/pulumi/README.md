# AzureFrontDoorSecurityPolicy -- Pulumi Module

Creates a `cdn.FrontdoorSecurityPolicy` (pulumi-azure classic v6) on
the referenced Front Door profile, through the shared Azure provider
builder (static client secret, keyless web identity, or ambient chain).

## Behavior notes

- **The wrapper nesting is rebuilt here**: the spec's flat four fields
  become the provider's securityPolicies -> firewall -> association
  shape (a one-choice ARM union).
- **`patternsToMatch` is the constant `"/*"`** -- the service accepts
  no other pattern. Bridge dialect: pulumi flattens azurerm's one-item
  list to a singular string (the Terraform module sends `["/*"]`);
  same ARM payload.
- **Domain references arrive resolved**: endpoint and custom-domain
  ARM ids are literals by the time the module runs; both types pass
  through unchanged.
- **No Azure tags** -- ARM does not support tags on security policies.

## Outputs

- `security_policy_id` -- operational addressing (nothing composes on
  the association itself)
- `security_policy_name`

## Build

```shell
make build
```
