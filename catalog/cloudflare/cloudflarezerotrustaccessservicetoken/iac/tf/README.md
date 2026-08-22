# CloudflareZeroTrustAccessServiceToken Terraform Module

Terraform IaC module for an Access service token -- a machine credential (client-ID / client-secret) for non-human clients.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareZeroTrustAccessServiceTokenSpec
locals.tf     — Scope + rotation-pair shaping
main.tf       — cloudflare_zero_trust_access_service_token
outputs.tf    — service_token_id, client_id, client_secret (sensitive), expires_at
```

## Behavior

Set exactly one of `account_id` or `zone_id`. `client_secret_version` and `previous_client_secret_expires_at` travel as a pair. `client_secret` is marked sensitive and is returned only at create and rotation.

## Outputs

| Name | Description |
|------|-------------|
| `service_token_id` | UUID used for import and policy `service_token` rules |
| `client_id` | Presented in `CF-Access-Client-ID` |
| `client_secret` | Sensitive. Returned only at create and rotation |
| `expires_at` | RFC3339 expiry |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
