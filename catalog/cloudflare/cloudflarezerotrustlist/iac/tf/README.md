# CloudflareZeroTrustList Terraform Module

Terraform IaC module for a reusable Zero Trust list (domains, IPs, URLs, emails, and kin) that Gateway policies reference by ID.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareZeroTrustListSpec
locals.tf     — items set shaping
main.tf       — cloudflare_zero_trust_list
outputs.tf    — list_id
```

## Behavior

`type` is immutable at Cloudflare (RequiresReplace). Items are a set -- order is not preserved. URL-type values are normalized by the API and can produce a perpetual plan diff at v5.23.0.

## Outputs

| Name | Description |
|------|-------------|
| `list_id` | UUID referenced by Gateway policies and posture rules |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
