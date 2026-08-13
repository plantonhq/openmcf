---
title: "Service Networking Connection"
description: "Service Networking Connection deployment documentation"
icon: "package"
order: 100
componentName: "gcpservicenetworkingconnection"
---

# GCP Service Networking Connection

Establishes **Private Services Access (PSA)** — the VPC peering join that makes AlloyDB, Cloud SQL, Memorystore, and Filestore private IP connectivity real. Without this connection, managed services cannot allocate private IPs from your network. Compose it after reserving a `VPC_PEERING` range on a GcpGlobalAddress and before deploying any private-IP database.

Peers a VPC with a Google managed-services producer (default `servicenetworking.googleapis.com`) so those services allocate private IPs from ranges you reserve inside the network.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Service Networking API enablement** on the target project
- **Private services access connection** — a VPC peering between your network and the producer, backed by one or more INTERNAL `VPC_PEERING` address ranges (`GcpGlobalAddress` resources referenced by name)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** with credentials for the target project
- **Planton Runner** when using Runner-based credential delivery

### GCP Prerequisites

- **GcpVpcNetwork** — the consumer VPC to peer
- **GcpGlobalAddress** — at least one INTERNAL address with purpose `VPC_PEERING` (a /16 is the common default)

## Deploy

### Console

Open the deployment store, find **GCP Service Networking Connection**, and click **Deploy**. Start from the **Private Services Access** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpServiceNetworkingConnection
metadata:
  name: my-vpc-psa
spec:
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  service: servicenetworking.googleapis.com
  reservedPeeringRanges:
    - valueFrom:
        kind: GcpGlobalAddress
        name: prod-vpc-psa-range
        fieldPath: status.outputs.name
```

```shell
planton apply -f psa-connection.yaml
```

### InfraChart

When deploying as part of a multi-resource environment, wire the connection to a VPC and reserved range deployed in the same InfraPipeline:

```yaml
spec:
  network:
    valueFrom:
      kind: GcpVpcNetwork
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  reservedPeeringRanges:
    - valueFrom:
        kind: GcpGlobalAddress
        name: prod-vpc-psa-range
        fieldPath: status.outputs.name
```

The InfraPipeline deploys the VPC and global address first, then establishes the PSA peering that downstream AlloyDB and Cloud SQL resources require.

## Key Configuration

**Network + service** — One connection per (network, service) pair. Both are create-time.

**Reserved peering ranges** — Named `GcpGlobalAddress` resources (by NAME, not CIDR). Append ranges later when the producer runs out of space.

**Adopt on creation conflict** — Recovery lever when a connection already exists outside Planton.

## Outputs and Dependencies

### Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpVpcNetwork** | `network` | `status.outputs.network_self_link` |
| **GcpGlobalAddress** | `reservedPeeringRanges` | `status.outputs.name` |

### Provides

| Output | Description |
|--------|-------------|
| `peering` | VPC peering name on the network (e.g. servicenetworking-googleapis-com) |
| `network` | Self-link of the peered network |

## Presets

**Private Services Access** — The canonical single-range PSA setup for managed services private IP. Start from the **Private Services Access** preset.

**Multi Range Growth** — Two reserved ranges from the start — the shape for networks expecting heavy private-IP growth. Start from the **Multi Range Growth** preset.

## Works With

- [**GCP Global Address**](/cloud-catalog/gcp-global-address) — reserves VPC_PEERING space upstream
- [**GCP VPC Network**](/cloud-catalog/gcp-vpc-network) — the peered consumer network
- [**GCP AlloyDB Cluster**](/cloud-catalog/gcp-alloydb-cluster) — private IP requires PSA
- [**GCP Cloud SQL**](/cloud-catalog/gcp-cloud-sql) — private IP requires PSA
