---
title: "Certificate Manager DNS Authorization"
description: "Certificate Manager DNS Authorization deployment documentation"
icon: "package"
order: 100
componentName: "gcpcertmanagerdnsauthorization"
---

# GCP Certificate Manager DNS Authorization

Deploys one Certificate Manager DNS authorization — the proof of domain
control that lets Google-managed certificates issue BEFORE traffic serves,
and the only validation mode that supports wildcard domains. One
authorization covers a domain and its wildcard.

## What Gets Created

When you deploy a GcpCertManagerDnsAuthorization resource, Planton provisions:

- **DNS Authorization** — a `google_certificate_manager_dns_authorization`
  for the domain
- **Certificate Manager API enablement** — `certificatemanager.googleapis.com`
  is enabled on the target project (never disabled on destroy)

The authorization exports its CNAME validation record; serving that record
from the domain's zone (a `GcpDnsRecord`) completes validation.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- Control of the domain's DNS zone (to serve the validation record)

## Quick Start

Create a file `dns-authorization.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpCertManagerDnsAuthorization
metadata:
  name: example-com-auth
spec:
  projectId: my-gcp-project-123
  domain: example.com
```

Deploy:

```shell
planton apply -f dns-authorization.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `domain` | `string` | The domain being authorized — covers the domain and its wildcard. Bare domain only (no `*.` prefix, no trailing dot). Immutable. | Required, valid bare domain |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project. Can reference a GcpProject resource. |
| `authorizationName` | `string` | `metadata.name` | Authorization name in GCP (1-64 chars). |
| `description` | `string` | — | Human-readable description. |
| `location` | `string` | `global` | Certificate Manager location; regional authorizations pair with regional certificates. |
| `type` | `string` | location-dependent | `FIXED_RECORD` (classic DNS-01 record) or `PER_PROJECT_RECORD` (shareable across projects). Immutable. |
| `labels` | `map<string,string>` | — | User labels merged beneath platform attribution labels. |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `authorization_id` | `string` | Fully-qualified resource ID — the value a certificate's `dnsAuthorizations` list consumes |
| `authorization_name` | `string` | Authorization name in GCP |
| `domain` | `string` | The authorized domain |
| `dns_record_name` | `string` | Fully-qualified name of the validation record to serve |
| `dns_record_type` | `string` | Validation record type (CNAME) |
| `dns_record_data` | `string` | Validation record data — the CNAME target |

## Related Components

- [GcpCertManagerCert](/docs/catalog/gcp/cert-manager-certificate) — references this authorization
- [GcpDnsRecord](/docs/catalog/gcp/dns-record) — serves the validation record
- [GcpDnsZone](/docs/catalog/gcp/dns-zone) — the zone hosting the record
