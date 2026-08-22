# Cloudflare mTLS Certificate

A certificate uploaded to the account-level mTLS store: the trust material that Authenticated Origin Pulls rows, zone TLS CA associations, and Workers mTLS bindings reference. Self-signed CAs are the normal case -- your infrastructure validates these, not the public. Every field is create-only: any change replaces the upload and the certificate ID changes.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **mTLS certificate** -- one `cloudflare_mtls_certificate` in the account store, CA or leaf per the `ca` flag

## Prerequisites

- **A Cloudflare account** (free plans included -- the store has no plan gate)
- **The PEM certificate or CA chain** to upload; a private key only when Cloudflare must present the certificate itself
- **A Cloudflare API token** with Account → SSL and Certificates → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareMtlsCertificate
metadata:
  name: origin-pull-ca
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: origin-pull-ca
  ca: true
  certificates: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
```

```shell
planton apply -f mtls-ca.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account the certificate is uploaded to. | Required, 32-hex. |
| `ca` | bool | CA certificate (true) or leaf certificate (false). | Required -- must be stated explicitly. |
| `certificates` | string | The certificate (or CA chain) in PEM form. | Required. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | unset | Display name in the dashboard. |
| `privateKey` | StringValueOrRef | unset | Only for leaf certificates Cloudflare presents itself. Sensitive -- provide a managed-secret reference; the API never returns it. |

## Examples

### CA for Authenticated Origin Pulls

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareMtlsCertificate
metadata:
  name: origin-pull-ca
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: origin-pull-ca
  ca: true
  certificates: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
```

### Leaf certificate with its key

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareMtlsCertificate
metadata:
  name: worker-client-leaf
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  ca: false
  certificates: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  privateKey:
    value: ${secrets-group/prod-mtls/worker-client-key}
```

## Destroy Semantics

Destroy is a real delete. Re-point consumers (CA associations, Workers bindings) at replacement material BEFORE destroying, or they lose their trust anchor.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `certificate_id` | string | The uploaded certificate's ID -- what consumers reference |
| `expires_on` | string | Expiry timestamp (RFC3339) |
| `serial_number` | string | The certificate's serial number |

## Related Components

- [Cloudflare Zone TLS Settings](/docs/catalog/cloudflare/cloudflarezonetlssettings) -- CA hostname associations reference this kind
- [Cloudflare Authenticated Origin Pulls](/docs/catalog/cloudflare/cloudflareauthenticatedoriginpulls) -- the zone-level enablement this trust material serves
- [Cloudflare Worker](/docs/catalog/cloudflare/cloudflareworker) -- mTLS bindings reference uploads from this store
