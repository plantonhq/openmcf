# CloudflareMtlsCertificate Terraform Module

Terraform IaC module for a certificate uploaded to the account-level mTLS certificate store.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareMtlsCertificateSpec
locals.tf     — Empty-string drops (name, private_key)
main.tf       — cloudflare_mtls_certificate
outputs.tf    — certificate_id, expires_on, serial_number
```

## Behavior

Every argument is create-only at the API -- any change plans a replacement, and the certificate ID changes (rotate by replace, re-point consumers, then destroy the old upload). The private key is optional: CA uploads used only to validate clients carry no key. Destroy is a real delete.

## Outputs

| Name | Description |
|------|-------------|
| `certificate_id` | The uploaded certificate's ID -- what consumers reference |
| `expires_on` | Expiry timestamp (RFC3339) |
| `serial_number` | The certificate's serial number |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
