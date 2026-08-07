# GCP Certificate Manager Certificate

Deploys one Certificate Manager certificate — Google-managed (auto-renewed)
or self-managed (uploaded PEM) — the modern certificate resource external
Application Load Balancers consume via a target HTTPS proxy's
`certificateManagerCertificates` list or a certificate map.

## What Gets Created

When you deploy a GcpCertManagerCert resource, Planton provisions:

- **Certificate** — a `google_certificate_manager_certificate` with either
  the managed or the self-managed arm
- **Certificate Manager API enablement** — `certificatemanager.googleapis.com`
  is enabled on the target project (never disabled on destroy)

DNS authorizations are NOT created here — they are first-class
`GcpCertManagerDnsAuthorization` resources this certificate references.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- For DNS-authorized managed certificates: a **GcpCertManagerDnsAuthorization**
  per distinct domain, with its validation record served from the domain's zone
  (compose a **GcpDnsRecord** from the authorization's outputs)

## Quick Start

Create a file `certificate.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpCertManagerCert
metadata:
  name: web-cert
spec:
  projectId: my-gcp-project-123
  managed:
    domains:
      - app.example.com
    dnsAuthorizations:
      - valueFrom:
          kind: GcpCertManagerDnsAuthorization
          name: app-example-com-auth
          fieldPath: status.outputs.authorization_id
```

Deploy:

```shell
planton apply -f certificate.yaml
```

## Configuration Reference

### Required Fields

Exactly one of `managed` or `selfManaged` must be set:

| Field | Type | Description |
|-------|------|-------------|
| `managed` | `object` | Google-managed arm: `domains[]` (required; wildcards allowed with DNS auth), `dnsAuthorizations[]` (refs to GcpCertManagerDnsAuthorization), `issuanceConfig` (private-PKI config path; mutually exclusive with dnsAuthorizations). Omit both auth fields for load-balancer authorization. |
| `selfManaged` | `object` | Uploaded arm: `pemCertificate` (leaf first, then intermediates) and `pemPrivateKey` (secret — masked in outputs). |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project. Can reference a GcpProject resource. |
| `certName` | `string` | `metadata.name` | Certificate name in GCP (1-64 chars, unique per location). Immutable. |
| `description` | `string` | — | Human-readable description. |
| `location` | `string` | `global` | Certificate Manager location; regional certificates serve regional load balancers only. Immutable. |
| `scope` | `string` | `DEFAULT` | `DEFAULT`, `EDGE_CACHE` (Media CDN), `ALL_REGIONS` (global certs only), or `CLIENT_AUTH` (backend mTLS client cert). Immutable. |
| `labels` | `map<string,string>` | — | User labels merged beneath platform attribution labels. |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `certificate_id` | `string` | Fully-qualified resource ID (`projects/{project}/locations/{location}/certificates/{name}`) |
| `certificate_name` | `string` | Certificate name — the value a target HTTPS proxy's `certificateManagerCertificates` list consumes |
| `san_dnsnames` | `string[]` | Subject Alternative Names in the issued certificate |
| `location` | `string` | The Certificate Manager location |
| `managed_state` | `string` | `PROVISIONING`/`FAILED`/`ACTIVE` for managed certificates; empty for self-managed. Stays `PROVISIONING` until domain validation completes. |

## Related Components

- [GcpCertManagerDnsAuthorization](/docs/catalog/gcp/cert-manager-dns-authorization) — the domain-control proof this certificate references
- [GcpDnsRecord](/docs/catalog/gcp/dns-record) — serves the authorization's validation record
- [GcpTargetHttpsProxy](/docs/catalog/gcp/target-https-proxy) — consumes the certificate
- [GcpManagedSslCertificate](/docs/catalog/gcp/managed-ssl-certificate) / [GcpSslCertificate](/docs/catalog/gcp/ssl-certificate) — the classic compute certificate kinds
