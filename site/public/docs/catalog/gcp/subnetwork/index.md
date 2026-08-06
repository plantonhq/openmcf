---
title: "Subnetwork"
description: "Subnetwork deployment documentation"
icon: "package"
order: 100
componentName: "gcpsubnetwork"
---

# GCP Subnetwork

Creates a subnetwork in a custom-mode VPC — the regional address space workloads live in: a primary IPv4 range, secondary ranges for alias IPs (GKE pods and services), optional IPv6 dual-stack, special-purpose roles (proxy-only subnets for regional load balancers, Private Service Connect), and VPC Flow Logs.

## What Gets Created

When you deploy a GcpSubnetwork resource, Planton provisions:

- **Compute Engine API enablement** — ensures `compute.googleapis.com` is active in the target project (never disabled on destroy)
- **Subnetwork** — a `google_compute_subnetwork` in the specified region and VPC with all configured ranges, purpose, IPv6, and logging

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing custom-mode VPC** — referenced via `vpcSelfLink`
- **IAM permissions** — `roles/compute.networkAdmin` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpSubnetwork
metadata:
  name: app-subnet
spec:
  projectId:
    value: my-gcp-project-123
  vpcSelfLink:
    valueFrom:
      kind: GcpVpcNetwork
      name: my-vpc
      fieldPath: status.outputs.network_self_link
  subnetworkName: app-subnet
  region: us-central1
  ipCidrRange: 10.10.0.0/20
  privateIpGoogleAccess: true
```

```shell
planton apply -f subnetwork.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `vpcSelfLink` | `StringValueOrRef` | — | Required. The parent VPC. Can reference a GcpVpcNetwork. Immutable. |
| `subnetworkName` | `string` | — | Required. Name in GCP (RFC1035). Immutable. |
| `region` | `string` | — | Required. Region of the subnet. Immutable. |
| `ipCidrRange` | `string` | — | Primary IPv4 CIDR. Required except for IPv6-only subnets. Expandable in place; never shrinkable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the subnet. Immutable. |
| `purpose` | `string` | `PRIVATE` | `PRIVATE`, `REGIONAL_MANAGED_PROXY` (proxy-only), `GLOBAL_MANAGED_PROXY`, `PRIVATE_SERVICE_CONNECT`, `PEER_MIGRATION`, `PRIVATE_NAT`. |
| `role` | `string` | — | `ACTIVE` / `BACKUP` on proxy-only subnets. |
| `secondaryIpRanges` | `list` | `[]` | Named alias-IP ranges — how GKE gets pod/service IPs. |
| `privateIpGoogleAccess` | `bool` | `false` | Internal access to Google APIs for VMs without external IPs. |
| `privateIpv6GoogleAccess` | `string` | GCP default | IPv6 counterpart of the above. |
| `stackType` | `string` | `IPV4_ONLY` | `IPV4_ONLY`, `IPV4_IPV6`, `IPV6_ONLY`. |
| `ipv6AccessType` | `string` | — | `EXTERNAL` or `INTERNAL`; required for IPv6-carrying stack types. Immutable. |
| `allowSubnetCidrRoutesOverlap` | `bool` | `false` | Permit deliberate CIDR overlap with routes outside the VPC. |
| `sendSecondaryIpRangeIfEmpty` | `bool` | `false` | Safety latch: whether an empty range list removes existing secondary ranges. |
| `logConfig` | object | off | VPC Flow Logs: aggregation interval, sampling, metadata scope, CEL filter. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `subnetwork_self_link` | Self-link — the value GKE clusters, instances, and other consumers reference |
| `subnetwork_name` | Name in GCP (referenced by by-name consumers like Cloud Run Direct VPC egress) |
| `region` / `ip_cidr_range` | Placement and primary range |
| `secondary_ranges` | Names + CIDRs of secondary ranges (GKE selects them by name) |
| `gateway_address` | IPv4 address of the subnet's default gateway |
| `subnetwork_id` | Server-assigned numeric ID |
| `internal_ipv6_prefix` / `external_ipv6_prefix` | Allocated IPv6 prefixes (empty without IPv6) |

## Related Components

- [GcpVpcNetwork](/docs/catalog/gcp/vpc) — provides the parent VPC network
- [GcpGkeCluster](/docs/catalog/gcp/gke-cluster) — consumes the subnet and secondary ranges for node, pod, and service networking
- [GcpRouterNat](/docs/catalog/gcp/router-nat) — outbound internet for private-only subnets
- [GcpProject](/docs/catalog/gcp/project) — manages the GCP project that hosts the subnet
- [GcpComputeInstance](/docs/catalog/gcp/compute-instance) — VMs deployed into this subnet
