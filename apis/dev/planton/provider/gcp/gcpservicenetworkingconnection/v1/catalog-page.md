# GCP Service Networking Connection

Creates a private services access connection — the VPC peering between your network and a service producer's network. This is the prerequisite that lets Cloud SQL, AlloyDB, Memorystore (PRIVATE_SERVICE_ACCESS mode), and Filestore hand out private IPs from ranges reserved inside your VPC.

## What Gets Created

One connection per (network, service) pair: a VPC peering on the network (e.g. `servicenetworking-googleapis-com`) from which the producer carves service subnets out of your reserved `VPC_PEERING` address ranges.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A VPC network** — referenced via `network` (a `GcpVpc` resource or a literal self-link)
- **At least one reserved range** — a `GcpGlobalAddress` with `addressType: INTERNAL` and `purpose: VPC_PEERING`
- **IAM permissions** — `servicenetworking.services.addPeering` (e.g. `roles/servicenetworking.networksAdmin`)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpServiceNetworkingConnection
metadata:
  name: prod-vpc-psa
spec:
  network:
    valueFrom:
      kind: GcpVpc
      name: prod-vpc
      fieldPath: status.outputs.network_self_link
  reservedPeeringRanges:
    - valueFrom:
        kind: GcpGlobalAddress
        name: prod-vpc-psa-range
        fieldPath: status.outputs.name
```

```shell
planton apply -f psa-connection.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `network` | `StringValueOrRef` | — (required) | VPC network to peer — name or self-link. Can reference a GcpVpc. Immutable. |
| `reservedPeeringRanges` | `StringValueOrRef[]` | — (min 1) | Names of INTERNAL `VPC_PEERING` global address ranges. Mutable — append to grow capacity. |
| `service` | `string` | `servicenetworking.googleapis.com` | The service producer to peer with. Immutable. |
| `updateOnCreationFail` | `bool` | `false` | Adopt a pre-existing connection for the same pair instead of failing. |
| `projectId` | `StringValueOrRef` | provider default | Project used for in-module API enablement. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `peering` | Name of the VPC peering created on the network (e.g. `servicenetworking-googleapis-com`) |
| `network` | The peered VPC network as the connection resolved it |

## Related Components

- [GcpVpc](/docs/catalog/gcp/gcpvpc) — the network being peered
- [GcpGlobalAddress](/docs/catalog/gcp/gcpglobaladdress) — reserves the `VPC_PEERING` ranges handed to the producer
- [GcpCloudSql](/docs/catalog/gcp/gcpcloudsql) — uses private IP through this connection
