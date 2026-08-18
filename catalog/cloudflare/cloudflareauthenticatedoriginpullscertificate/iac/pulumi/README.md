# CloudflareAuthenticatedOriginPullsCertificate Pulumi Module

Pulumi (Go) IaC module for an Authenticated Origin Pulls client-certificate upload (zone-wide or hostname-scoped).

## Architecture

```
main.go                   — Entrypoint loading the stack input
module/main.go            — Resources(): provider setup, resource, outputs
module/locals.go          — Locals initialization
module/certificate.go     — cloudflare.AuthenticatedOriginPullsCertificate
                            XOR cloudflare.AuthenticatedOriginPullsHostnameCertificate
module/outputs.go         — Stack output keys
```

## Behavior

Mirrors the Terraform module's contract exactly: the `scope` selects which resource is created (exactly one), rotation is replacement (never key-only -- the zone surface silently ignores it at v5.23.0), and the `certificate_id` / `zone_id` / `expires_on` / `status` stack outputs.

## Outputs

| Name | Description |
|------|-------------|
| `certificate_id` | The uploaded certificate's ID -- what associations reference |
| `zone_id` | The zone the certificate belongs to |
| `expires_on` | Expiry timestamp (RFC3339) |
| `status` | Deployment status (asynchronous) |

## Provider Version

Uses `pulumi-cloudflare` SDK v6 (Cloudflare Terraform provider v5 line).
