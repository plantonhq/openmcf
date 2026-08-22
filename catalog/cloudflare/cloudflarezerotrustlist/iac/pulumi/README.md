# CloudflareZeroTrustList Pulumi Module

Pulumi (Go) IaC module for a reusable Zero Trust list (domains, IPs, URLs, emails, and kin) that Gateway policies reference by ID.

## Architecture

```
main.go                    — Entrypoint loading the stack input
module/main.go             — Resources(): provider setup, resource, outputs
module/locals.go           — Locals initialization
module/zero_trust_list.go  — cloudflare.ZeroTrustList
module/outputs.go          — Stack output keys
```

## Behavior

Mirrors the Terraform module's contract exactly: immutable `type`, set-semantics items, `list_id` stack output.

## Outputs

| Name | Description |
|------|-------------|
| `list_id` | UUID referenced by Gateway policies and posture rules |

## Provider Version

Uses `pulumi-cloudflare` SDK v6 (Cloudflare Terraform provider v5 line).
