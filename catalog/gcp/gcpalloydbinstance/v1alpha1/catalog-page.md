# GCP AlloyDB Instance

Creates an AlloyDB instance on an existing cluster — typically a READ_POOL for read scaling, with optional PRIMARY/SECONDARY types for advanced topologies.

## What Gets Created

- **AlloyDB API enablement** on the project
- **AlloyDB instance** — `google_alloydb_instance` on the referenced cluster

## Prerequisites

- **An existing [GcpAlloydbCluster](/docs/catalog/gcp/gcpalloydbcluster)** (or literal cluster resource path)
- **GCP credentials** with AlloyDB permissions

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpAlloydbInstance
metadata:
  name: orders-read-pool
spec:
  cluster:
    valueFrom:
      kind: GcpAlloydbCluster
      name: my-orders-cluster
      fieldPath: status.outputs.cluster_id
  instanceId: orders-read-pool
  readPoolConfig:
    nodeCount: 2
```

## Configuration Reference

| Field | Default | Description |
|-------|---------|-------------|
| `cluster` | — (required) | Cluster ref → `cluster_id`. Immutable. |
| `instanceId` | — (required) | Instance ID. Immutable. |
| `instanceType` | `READ_POOL` | `PRIMARY`, `READ_POOL`, or `SECONDARY`. |
| `readPoolConfig.nodeCount` | — | Required for READ_POOL. |
| `cpuCount` / `machineType` | — | XOR machine sizing. |
| `requireConnectors` / `sslMode` | — | Connection security. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `instance_name` | Full instance resource path |
| `ip_address` | Private connection IP |
| `state` | Instance state |

## Related Components

- [GcpAlloydbCluster](/docs/catalog/gcp/gcpalloydbcluster)
- [GcpAlloydbUser](/docs/catalog/gcp/gcpalloydbuser)
