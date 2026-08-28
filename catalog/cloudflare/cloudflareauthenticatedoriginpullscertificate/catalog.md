# Cloudflare Authenticated Origin Pulls Certificate

The client certificate Cloudflare presents to your origin under Authenticated Origin Pulls. Self-signed is the normal case -- the origin validates it, not the public. `scope: zone` replaces the zone-wide certificate; `scope: hostname` uploads material for per-hostname associations. Deletion settles asynchronously, and a key-only rotation against the zone surface silently no-ops (rotate certificate and key together).

## What Gets Created

When you deploy this resource, the IaC module provisions exactly one of:

- **Zone upload** -- `cloudflare_authenticated_origin_pulls_certificate` when `scope` is `zone` (the default)
- **Hostname upload** -- `cloudflare_authenticated_origin_pulls_hostname_certificate` when `scope` is `hostname`

## Prerequisites

- **A Cloudflare zone** (free plans included)
- **A client certificate and key** -- self-signed pairs are normal; the origin's trust store validates them
- **A Cloudflare API token** with Zone → SSL and Certificates → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareAuthenticatedOriginPullsCertificate
metadata:
  name: app-client-certificate
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  scope: hostname
  certificate: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  privateKey:
    value: ${secrets-group/prod-aop/app-client-key}
```

```shell
planton apply -f aop-certificate.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone the certificate is uploaded to. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. |
| `certificate` | string | The client certificate in PEM form. | Required. |
| `privateKey` | StringValueOrRef | The private key in PEM form -- provide a managed-secret reference; the API never returns it. | Required, sensitive. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `scope` | string | `zone` | `zone` (replaces the zone-wide client certificate) or `hostname` (material for per-hostname associations). |

## Examples

### Zone-wide client certificate

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareAuthenticatedOriginPullsCertificate
metadata:
  name: zone-client-certificate
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: acme-zone
      fieldPath: status.outputs.zone_id
  scope: zone
  certificate: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  privateKey:
    value: ${secrets-group/prod-aop/zone-client-key}
```

## Destroy Semantics

Real delete on both surfaces, settling asynchronously (200 with `pending_deletion`/`deleted` before the record goes). Revert or re-point referencing associations first.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `certificate_id` | string | The uploaded certificate's ID -- what associations reference |
| `zone_id` | string | The zone the certificate belongs to |
| `expires_on` | string | Expiry timestamp (RFC3339) |

Deployment status is deliberately not a stack output: it transitions asynchronously (pending_deployment to active seconds after create), so a point-in-time phase would flip on the first refresh and re-plan forever. Read deployment status from the Cloudflare API or dashboard.

## Related Components

- [Cloudflare Authenticated Origin Pulls](/docs/catalog/cloudflare/cloudflareauthenticatedoriginpulls) -- the enablement whose rows reference this kind
- [Cloudflare mTLS Certificate](/docs/catalog/cloudflare/cloudflaremtlscertificate) -- account-level CA trust material
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key
