---
title: "Regional Address"
description: "Regional Address deployment documentation"
icon: "package"
order: 100
componentName: "gcpaddress"
---

# GCP Regional Address

Reserves a static IP address at regional scope — either a public IPv4/IPv6 address for Cloud NAT, regional load balancers, and VM instances, or a private IP/range within a VPC subnetwork for GCE endpoints, internal LB VIPs, VPC peering, and IPsec interconnect. The component automatically enables the Compute Engine API on the target project.

## What Gets Created

When you deploy a GcpAddress resource, Planton provisions:

- **Compute Engine API enablement** — activates `compute.googleapis.com` on the target project
- **Regional Address** — a `google_compute_address` resource with the specified name, region, address type, and network configuration

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **An existing VPC network** — required for INTERNAL VPC_PEERING / IPSEC_INTERCONNECT addresses
- **An existing subnetwork** — required for INTERNAL GCE_ENDPOINT / DNS_RESOLVER addresses
- **IAM permissions** — `roles/compute.networkAdmin` or equivalent on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpAddress
metadata:
  name: nat-external-ip
spec:
  projectId:
    value: my-gcp-project-123
  addressName: nat-external-ip
  region: us-central1
```

```shell
planton apply -f regional-address.yaml
```

This reserves a public IPv4 address in `us-central1` for Cloud NAT or a regional load balancer.

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `addressName` | `string` | — (required) | Cloud-side name (RFC1035). Immutable. |
| `region` | `string` | — (required) | GCP region (e.g. `us-central1`). Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the address. |
| `addressType` | `string` | `EXTERNAL` | `EXTERNAL` or `INTERNAL`. Immutable. |
| `ipVersion` | `string` | `IPV4` | `IPV4` or `IPV6`. Immutable. |
| `network` | `StringValueOrRef` | — | VPC for INTERNAL VPC_PEERING / IPSEC_INTERCONNECT. |
| `subnetwork` | `StringValueOrRef` | — | Subnetwork for INTERNAL GCE_ENDPOINT / DNS_RESOLVER. |
| `networkTier` | `string` | `PREMIUM` | `PREMIUM` or `STANDARD` — EXTERNAL only. |
| `prefixLength` | `int32` | — | CIDR prefix (8-29) for peering/interconnect ranges. |
| `purpose` | `string` | — | INTERNAL purpose enum. Immutable. |
| `ipv6EndpointType` | `string` | — | `VM` or `NETLB` for external IPv6. |

## Examples

### External Static IP for Cloud NAT

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpAddress
metadata:
  name: nat-ip
spec:
  projectId:
    value: my-prod-project-123
  addressName: nat-ip
  region: us-central1
  description: Static IP for Cloud NAT
```

### Internal GCE Endpoint in a Subnetwork

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpAddress
metadata:
  name: vm-internal-ip
spec:
  projectId:
    value: my-prod-project-123
  addressName: vm-internal-ip
  region: us-central1
  addressType: INTERNAL
  purpose: GCE_ENDPOINT
  subnetwork:
    valueFrom:
      kind: GcpSubnetwork
      name: app-subnet
      fieldPath: status.outputs.subnetwork_self_link
```

### Internal Shared Load Balancer VIP

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpAddress
metadata:
  name: ilb-vip
spec:
  projectId:
    value: my-prod-project-123
  addressName: ilb-vip
  region: us-central1
  addressType: INTERNAL
  purpose: SHARED_LOADBALANCER_VIP
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `address` | Reserved IP address or range start |
| `self_link` | Self-link URL — the composition handle for NAT, LBs, and VMs |
| `name` | Name of the address resource in GCP |
| `region` | Plain region name from the spec |

## Related Components

- [GcpGlobalAddress](/docs/catalog/gcp/global-address) — global-scope addresses (HTTP(S) LB VIPs, global VPC peering ranges, PSC)
- [GcpVpcNetwork](/docs/catalog/gcp/vpc) — VPC network for INTERNAL addresses
- [GcpSubnetwork](/docs/catalog/gcp/subnetwork) — subnetwork for GCE endpoint addresses
- [GcpProject](/docs/catalog/gcp/project) — GCP project and API enablement
