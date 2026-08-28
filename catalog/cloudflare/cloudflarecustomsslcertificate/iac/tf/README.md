# CloudflareCustomSslCertificate Terraform Module

Terraform IaC module for a bring-your-own TLS certificate uploaded to a Cloudflare zone.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareCustomSslCertificateSpec
locals.tf     — Empty-string drops (policy, custom_csr_id) and geo_restrictions shaping
main.tf       — cloudflare_custom_ssl
outputs.tf    — certificate_id, zone_id, expires_on
```

## Behavior

Certificate, private key, and zone changes force replacement (rotation is destroy-and-create; Cloudflare serves the previous certificate until the replacement deploys). `priority` is deliberately absent -- read-only at provider v5.23.0. Custom certificates are a Business/Enterprise zone feature and must be publicly trusted; both walls are Cloudflare's, at create. Destroy is a real delete.

## Outputs

| Name | Description |
|------|-------------|
| `certificate_id` | The uploaded certificate's ID |
| `zone_id` | The zone the certificate belongs to |
| `expires_on` | Expiry timestamp (RFC3339) |

Deployment status is deliberately not an output: it transitions asynchronously (pending before active), so a point-in-time phase would flip on the first refresh and re-plan forever (the class was measured live on the sibling AOP certificate).

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
