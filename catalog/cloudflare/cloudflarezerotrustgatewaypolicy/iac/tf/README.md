# CloudflareZeroTrustGatewayPolicy Terraform Module

Terraform IaC module for one Gateway policy -- a filter over employee DNS/HTTP/network traffic plus the action taken on a match.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareZeroTrustGatewayPolicySpec
locals.tf     — filters wrap, expiration/schedule, rule_settings (always emitted)
main.tf       — cloudflare_zero_trust_gateway_policy
outputs.tf    — policy_id, precedence
```

## Behavior

`filter` is a singular string; locals wrap it to the provider's one-element list. `rule_settings` is always sent -- an empty object when the spec configures nothing -- because the provider drifts when the block is omitted. `enabled` is passed through explicitly so a missing value stays visible (Cloudflare defaults it to false).

Known upstream defect at v5.23.0: `rule_settings.add_headers` and `override_ips` drift on first apply.

## Outputs

| Name | Description |
|------|-------------|
| `policy_id` | UUID of the created policy |
| `precedence` | Evaluation order |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
