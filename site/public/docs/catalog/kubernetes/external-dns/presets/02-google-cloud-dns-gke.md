---
title: "Google Cloud DNS on GKE (Keyless via Workload Identity)"
description: "This preset installs ExternalDNS on a GKE cluster publishing to a Cloud DNS zone, authenticating keylessly through GKE Workload Identity — no service-account key anywhere. It scopes the instance to..."
type: "preset"
rank: "02"
presetSlug: "02-google-cloud-dns-gke"
componentSlug: "external-dns"
componentTitle: "External DNS"
provider: "kubernetes"
icon: "package"
order: 2
---

# Google Cloud DNS on GKE (Keyless via Workload Identity)

This preset installs ExternalDNS on a GKE cluster publishing to a Cloud DNS
zone, authenticating keylessly through GKE Workload Identity — no
service-account key anywhere. It scopes the instance to one zone, enables
full `sync` reconciliation, and tags ownership with a per-cluster TXT owner
ID. This is the standard production posture on GCP.

## When to Use

- GKE clusters whose Services/Ingresses publish records into a Cloud DNS zone
- Zones dedicated to (or safely shareable with) this cluster
- Production deployments — keyless Workload Identity is the recommended
  authentication

## Key Configuration Choices

- **Keyless authentication** (`workloadIdentity.gke.serviceAccountEmail`) —
  the controller ServiceAccount impersonates a GCP service account with DNS
  permissions; no credential Secrets are created
- **Zone scoping** (`zoneIdFilters` + `domainFilters`) — restricts
  management to one zone even if the GCP identity can see more
- **`policy: sync`** — full reconciliation; the TXT registry limits deletes
  to records tagged with this instance's `txtOwnerId`
- **Sources `service` + `ingress`** — the chart defaults, stated explicitly

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<gcp-project-id>` | GCP project owning the Cloud DNS zones | GCP console or `GcpProject` outputs |
| `<cloud-dns-zone-name>` | Cloud DNS managed-zone name | Cloud DNS console or `GcpDnsZone` outputs |
| `<external-dns-gsa-email>` | GCP service account with `roles/dns.admin`, carrying a `roles/iam.workloadIdentityUser` binding for `<project>.svc.id.goog[external-dns/my-external-dns-cloud-dns]` | IAM console or `GcpServiceAccount` outputs |
| `<cluster-name>` | Unique owner ID for this instance | Your cluster naming |
| `<example.com>` | Domain suffix this instance manages | Your zone's domain |

## Related Presets

- **01-aws-route53-eks-keyless** — the same posture on EKS + Route 53
- **04-cloudflare-any-cluster** — publish to Cloudflare from any cluster (including GKE) with a token
