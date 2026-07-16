# GCP SSL Certificate (Self-Managed)

Creates a self-managed Compute Engine SSL certificate — you bring the PEM chain and private key, and the load balancer presents them to clients. Reference its `self_link` from a target HTTPS (or SSL) proxy exactly like a Google-managed certificate.

## What Gets Created

A single SSL certificate uploaded from your PEM material — global when `region` is empty, regional when set.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **IAM permissions** — any role carrying `compute.sslCertificates.*` on the target project
- **A PEM certificate chain and unencrypted private key** from your CA

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSslCertificate
metadata:
  name: prod-app-cert
spec:
  projectId:
    value: my-gcp-project-123
  certificate: |
    -----BEGIN CERTIFICATE-----
    ...leaf, then intermediates...
    -----END CERTIFICATE-----
  privateKey: |
    -----BEGIN PRIVATE KEY-----
    ...matching unencrypted key...
    -----END PRIVATE KEY-----
```

```shell
planton apply -f cert.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `certificate` | `string` | — (required) | PEM chain: leaf first, then intermediates (max 5 certs). Immutable. |
| `privateKey` | `string` (sensitive) | — (required) | Matching unencrypted PEM key. Write-only in GCP. Immutable. |
| `region` | `string` | `""` (global) | Region for a regional certificate; empty means global. Immutable. |
| `certificateName` | `string` | `metadata.name` | Cloud-side name (RFC1035); shared namespace with managed certificates. Immutable. |
| `description` | `string` | `""` | What this certificate secures. Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the certificate. Immutable. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value a target HTTPS proxy references in `ssl_certificates` |
| `certificate_name` | Name of the SSL certificate in GCP |
| `certificate_id` | Server-assigned numeric ID of the certificate |
| `expire_time` | Expiry in RFC3339, parsed from the uploaded chain — plan rotation off this |
| `region` | Region of a regional certificate; empty for global |

## Related Components

- [GcpTargetHttpsProxy](/docs/catalog/gcp/gcptargethttpsproxy) — presents this certificate to clients
- [GcpManagedSslCertificate](/docs/catalog/gcp/gcpmanagedsslcertificate) — the Google-managed alternative
- [GcpSslPolicy](/docs/catalog/gcp/gcpsslpolicy) — hardens TLS versions and ciphers on the same proxy
