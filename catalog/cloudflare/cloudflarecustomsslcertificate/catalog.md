# Cloudflare Custom SSL Certificate

A bring-your-own TLS certificate uploaded to a zone -- Cloudflare presents it to visitors instead of Universal SSL. The certificate must be publicly trusted (self-signed is rejected), and custom certificates are a Business/Enterprise zone feature enforced at create. Rotation is replacement: the certificate ID changes, and Cloudflare serves the old certificate until the new one deploys.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Custom certificate** -- one `cloudflare_custom_ssl` upload on the zone, with the chosen SNI class, bundle method, and private-key geo restrictions

## Prerequisites

- **A Cloudflare zone on a Business plan or above** -- lower plans are rejected at the API
- **A certificate issued by a publicly trusted CA** covering the zone's hostnames, plus its private key
- **A Cloudflare API token** with Zone → SSL and Certificates → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareCustomSslCertificate
metadata:
  name: www-custom-certificate
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  certificate: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  privateKey:
    value: ${secrets-group/prod-tls/www-key}
  type: sni_custom
```

```shell
planton apply -f custom-certificate.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone the certificate is uploaded to. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. |
| `certificate` | string | PEM certificate, publicly trusted, covering the zone's hostnames. | Required. |
| `privateKey` | StringValueOrRef | PEM private key -- provide a managed-secret reference; the API never returns it. | Required, sensitive. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `type` | string | `legacy_custom` | `legacy_custom` (every client, one slot) or `sni_custom` (modern clients, multiple uploads). Changing it replaces the upload. |
| `bundleMethod` | string | `ubiquitous` | `ubiquitous`, `optimal`, or `force` (chain exactly as uploaded). |
| `policy` | string | unset | Geo-policy expression for the private key, e.g. `(country: US) or (region: EU)`. |
| `geoRestrictions.label` | string | unset | Coarse geo label: `us`, `eu`, `highest_security`. |
| `customCsrId` | string | unset | The Cloudflare-generated CSR the certificate was issued from, when the CSR flow was used. |
| `deploy` | string | unset (API: production) | `staging` validates on Cloudflare's staging network first (Business/Enterprise). |

## Examples

### EV certificate with a US private-key restriction

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareCustomSslCertificate
metadata:
  name: compliance-certificate
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: acme-zone
      fieldPath: status.outputs.zone_id
  certificate: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  privateKey:
    value: ${secrets-group/prod-tls/compliance-key}
  type: sni_custom
  bundleMethod: optimal
  geoRestrictions:
    label: us
```

## Destroy Semantics

Destroy is a real delete. The zone falls back to Universal SSL / Advanced certificates -- have one active first, or visitors get handshake failures.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `certificate_id` | string | The uploaded certificate's ID |
| `zone_id` | string | The zone the certificate belongs to |
| `expires_on` | string | Expiry timestamp (RFC3339) |

Deployment status is deliberately not a stack output: deployment is asynchronous (pending before active), so a point-in-time phase would flip on the first refresh and re-plan forever. Read deployment status from the Cloudflare API or dashboard.

## Related Components

- [Cloudflare Certificate Pack](/docs/catalog/cloudflare/cloudflarecertificatepack) -- Cloudflare-managed Advanced certificates that renew themselves
- [Cloudflare Zone TLS Settings](/docs/catalog/cloudflare/cloudflarezonetlssettings) -- the zone's TLS posture
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key
