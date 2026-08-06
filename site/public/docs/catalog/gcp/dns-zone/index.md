---
title: "DNS Zone"
description: "DNS Zone deployment documentation"
icon: "package"
order: 100
componentName: "gcpdnszone"
---

# GCP DNS Zone

Creates a Google Cloud DNS managed zone (`google_dns_managed_zone`) — the zone shell only. Individual DNS records are composed via the separate [GcpDnsRecord](/docs/catalog/gcp/dns-record) kind.

## What Gets Created

- **Cloud DNS API enablement** on the project (`dns.googleapis.com`)
- **Managed zone** — public, private, forwarding, or peering depending on spec

## Prerequisites

- **GCP credentials** with Cloud DNS permissions
- **A GCP project** where the zone will live
- **For private zones:** a [GcpVpcNetwork](/docs/catalog/gcp/vpc) (and/or [GcpGkeCluster](/docs/catalog/gcp/gke-cluster)) referenced in `privateVisibilityConfig`

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsZone
metadata:
  name: example.com
spec:
  projectId:
    value: my-gcp-project-123
```

Deploy creates a public zone for `example.com.` (dns_name derived from metadata.name).

## Configuration Reference

### Required Fields

| Field | Description |
|-------|-------------|
| `projectId` | GCP project ID (literal or GcpProject ref). Required. |

### Optional Fields

| Field | Default | Description |
|-------|---------|-------------|
| `dnsName` | `{metadata.name}.` | FQDN ending with `.` |
| `description` | `managed-zone for {name}` | Console description |
| `visibility` | `public` | `public` or `private` |
| `privateVisibilityConfig` | — | VPC networks and/or GKE clusters (private standard zones) |
| `dnssecConfig` | `state: off` | DNSSEC for public zones |
| `forwardingConfig` | — | Private forwarding zone targets |
| `peeringConfig` | — | Private peering zone target network |
| `cloudLoggingConfig` | — | Query logging toggle |
| `forceDestroy` | `false` | Delete all records on zone destroy |
| `labels` | — | Additional GCP labels |

### Validation Rules

- `forwardingConfig` and `peeringConfig` require `visibility: private`
- `forwardingConfig` and `peeringConfig` are mutually exclusive
- `privateVisibilityConfig` requires `visibility: private`

## Stack Outputs

| Output | Description |
|--------|-------------|
| `zone_id` | Managed zone ID |
| `zone_name` | Zone name — use as `managed_zone` FK in GcpDnsRecord |
| `nameservers` | Delegation nameservers (public zones) |

## Related Components

- [GcpDnsRecord](/docs/catalog/gcp/dns-record) — individual DNS records
- [GcpVpcNetwork](/docs/catalog/gcp/vpc) — private zone visibility
- [GcpGkeCluster](/docs/catalog/gcp/gke-cluster) — cluster-scoped private visibility
- [GcpProject](/docs/catalog/gcp/project) — owning project
