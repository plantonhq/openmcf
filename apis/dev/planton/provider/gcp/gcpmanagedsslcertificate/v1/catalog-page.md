# GCP Managed SSL Certificate

Creates a Google-managed SSL certificate — a global Compute Engine SSL certificate whose private key and issuance are handled entirely by Google. Reference its `self_link` from a target HTTPS proxy to terminate TLS at a global external Application Load Balancer.

## What Gets Created

A single global Google-managed SSL certificate covering the domains you specify.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **IAM permissions** — any role carrying `compute.sslCertificates.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpManagedSslCertificate
metadata:
  name: app-cert
spec:
  projectId:
    value: my-gcp-project-123
  domains:
    - app.example.com
```

```shell
planton apply -f cert.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `domains` | `string[]` | — (required) | FQDNs the certificate is valid for. No wildcards. Immutable. |
| `certificateName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable. |
| `description` | `string` | `""` | What this certificate secures. Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the certificate. Immutable. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value a target HTTPS proxy references in `ssl_certificates` |
| `certificate_name` | Name of the SSL certificate in GCP |
| `certificate_id` | Server-assigned numeric ID of the certificate |
| `expire_time` | Expiry time in RFC3339 format; empty until provisioning completes |

## Related Components

- [GcpUrlMap](/docs/catalog/gcp/gcpurlmap) — routes traffic behind the HTTPS proxy
- [GcpBackendService](/docs/catalog/gcp/gcpbackendservice) — backends the load balancer targets
- [GcpGlobalAddress](/docs/catalog/gcp/gcpglobaladdress) — static IP for DNS A records
