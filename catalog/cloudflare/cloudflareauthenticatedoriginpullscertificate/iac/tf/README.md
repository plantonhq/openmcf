# CloudflareAuthenticatedOriginPullsCertificate Terraform Module

Terraform IaC module for an Authenticated Origin Pulls client-certificate upload (zone-wide or hostname-scoped).

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareAuthenticatedOriginPullsCertificateSpec
locals.tf     — Scope resolution (zone vs hostname)
main.tf       — cloudflare_authenticated_origin_pulls_certificate (zone scope)
                XOR cloudflare_authenticated_origin_pulls_hostname_certificate (hostname scope)
outputs.tf    — certificate_id, zone_id, expires_on, status
```

## Behavior

The `scope` selects which provider resource is created -- exactly one of the two. Rotation is replacement on both surfaces; never rotate only the private key (the zone-scoped provider resource silently ignores key-only changes at v5.23.0 -- its Update is empty and the key does not force replacement). Destroy is a real delete settling asynchronously.

## Outputs

| Name | Description |
|------|-------------|
| `certificate_id` | The uploaded certificate's ID -- what associations reference |
| `zone_id` | The zone the certificate belongs to |
| `expires_on` | Expiry timestamp (RFC3339) |
| `status` | Deployment status (asynchronous) |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
