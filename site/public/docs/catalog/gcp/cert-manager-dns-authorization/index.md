---
title: "Cert Manager DNS Authorization"
description: "Cert Manager DNS Authorization deployment documentation"
icon: "package"
order: 100
componentName: "gcpcertmanagerdnsauthorization"
---

# GCP Cert Manager DNS Authorization

Creates one Certificate Manager DNS authorization — the proof of domain control a Google-managed certificate needs before it can be issued for a domain that is not yet serving traffic. One authorization covers a single domain AND its wildcard: authorizing `example.com` issues certificates for both `example.com` and `*.example.com`. It is the only validation mode that supports wildcards, and it decouples certificate issuance from traffic serving — the key to zero-downtime TLS migration.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Certificate Manager DNS Authorization** -- a `google_certificate_manager_dns_authorization` for the domain, exporting the DNS validation record (a CNAME) that must exist in the domain's zone

The validation record itself is NOT created here — serve it with a [GcpDnsRecord](/cloud-catalog/gcp-dns-record) wired to this kind's outputs, and validation completes automatically.

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.

### GCP Project

- **Certificate Manager API** (`certificatemanager.googleapis.com`) enabled in the target project.
- **An authoritative DNS zone** for the domain — typically a [GcpDnsZone](/cloud-catalog/gcp-dns-zone) — where the exported validation CNAME will be served.

## Deploy

### Console

Open the deployment store, find **GCP Cert Manager DNS Authorization**, and click **Create**. The wizard walks two decisions: the authorization's envelope (project, name, location), then the domain being authorized and its record scope. The [Presets](#presets) tab offers **Standard domain** and **Shared per-project** starting points.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCertManagerDnsAuthorization
metadata:
  name: orders-example-auth
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  domain: example.com
```

```shell
planton apply -f authorization.yaml
```

## Key Configuration

**Domain** -- Immutable, bare domain only (no `*.` prefix, no trailing dot). The authorization covers the domain AND its wildcard — never create a separate authorization for `*.example.com`. Subdomain wildcards need their own authorization (`*.example.com` does not cover `a.b.example.com`).

**Validation record scope** (`type`) -- Immutable. `FIXED_RECORD` is the classic DNS-01 style record, one per (domain, authorization); `PER_PROJECT_RECORD` scopes the record to (domain, project) so multiple Certificate Manager resources across projects share one record. Empty lets GCP pick the location-appropriate default.

**Location** -- Empty means `global`. Regional authorizations pair with regional certificates only — keep it aligned with the certificate that will consume this authorization.

**Authorization name** -- 1-64 chars starting with a letter; empty defaults to `metadata.name`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `authorization_id` | Fully-qualified resource ID (`projects/*/locations/*/dnsAuthorizations/*`) | A GcpCertManagerCert's `dnsAuthorizations` list — the family's composition key |
| `authorization_name` | The authorization's GCP name | Auditing |
| `domain` | The authorized domain | Cross-checking coverage |
| `dns_record_name` | The validation record's fully-qualified name | A GcpDnsRecord's record name |
| `dns_record_type` | The validation record's type (CNAME) | A GcpDnsRecord's record type |
| `dns_record_data` | The value the CNAME must point at | A GcpDnsRecord's record values |

## Works With

- [**GCP Cert Manager Cert**](/cloud-catalog/gcp-cert-manager-cert) -- the Google-managed certificate that references this authorization by ID
- [**GCP DNS Record**](/cloud-catalog/gcp-dns-record) -- serves the exported validation CNAME; wire its name/type/values to this kind's outputs
- [**GCP DNS Zone**](/cloud-catalog/gcp-dns-zone) -- the authoritative zone the record lives in
