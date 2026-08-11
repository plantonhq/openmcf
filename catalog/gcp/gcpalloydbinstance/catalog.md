# GCP AlloyDB Instance

Deploys a standalone AlloyDB instance (`google_alloydb_instance`) attached to an existing cluster. Read pools scale read traffic independently of the bundled primary in GcpAlloydbCluster. PRIMARY and SECONDARY types exist for advanced topologies; presets target READ_POOL.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **AlloyDB Instance** — compute node within an existing cluster (READ_POOL, PRIMARY, or SECONDARY)
- **Read Pool Config** — when `instanceType` is READ_POOL (or blank, which defaults to READ_POOL), configures `nodeCount` for read capacity
- **Machine Configuration** — `cpuCount` XOR `machineType` (mutually exclusive)
- **Network Surface** — optional public IP with authorized CIDR ranges, outbound public IP, and per-instance PSC config
- **Observability** — optional query insights and per-instance database flags

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** — credentials for the target project
- **Planton Runner** — when using Runner-based credential delivery

### GCP Prerequisites

- **GcpAlloydbCluster** — the cluster this instance joins. Reference `status.outputs.cluster_id`.
- **AlloyDB API** (`alloydb.googleapis.com`) enabled in the target project (usually already enabled by the parent cluster).

## Deploy

### Console

Open the deployment store, find **GCP AlloyDB Instance**, and click **Deploy**. Start from the **Read Pool Basic** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpAlloydbInstance
metadata:
  name: orders-read-pool
spec:
  cluster:
    valueFrom:
      kind: GcpAlloydbCluster
      name: orders-db
      fieldPath: status.outputs.cluster_id
  instanceId: orders-read-pool
  readPoolConfig:
    nodeCount: 2
  cpuCount: 4
  availabilityType: REGIONAL
```

### InfraChart

Wire the instance to a cluster deployed in the same InfraPipeline:

```yaml
spec:
  cluster:
    valueFrom:
      kind: GcpAlloydbCluster
      name: orders-db
      fieldPath: status.outputs.cluster_id
```

## Key Configuration

**Instance type** — Blank or `READ_POOL` is the common read-scaling shape. `PRIMARY` and `SECONDARY` exist for advanced topologies where the cluster does not bundle the compute node.

**Read pool sizing** — `readPoolConfig.nodeCount` sets how many read nodes serve traffic. Pair with `availabilityType: REGIONAL` for multi-zone read HA.

**Machine sizing** — Set `cpuCount` (2, 4, 8, 16, 32, 64, 96, 128) or `machineType` explicitly — never both.

**Connection security** — `requireConnectors: true` enforces IAM-based auth via AlloyDB Auth Proxy. `sslMode: ENCRYPTED_ONLY` requires TLS.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpAlloydbCluster** | `cluster` | `status.outputs.cluster_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | Fully qualified instance resource name | Monitoring, connection routing |
| `instance_name` | Short instance name | Display, logging |
| `ip_address` | Private IP of the instance | Application read connection strings |
| `state` | Instance state (`READY`, `CREATING`) | Health checks |

## Common Patterns

**Basic read pool** — single-node ZONAL pool for dev/staging read offload. Start from **Read Pool Basic**.

**HA read pool** — REGIONAL pool with two nodes for production read scaling. Start from **Read Pool HA**.

**Production read pool** — connectors enforced, TLS-only, query insights enabled. Start from **Read Pool Production**.

## Works With

- [**GCP AlloyDB Cluster**](/cloud-catalog/gcp-alloydb-cluster) — parent cluster this instance attaches to
- [**GCP AlloyDB User**](/cloud-catalog/gcp-alloydb-user) — application credentials on the same cluster
- [**GCP Project**](/cloud-catalog/gcp-project) — optional project override when the cluster lives elsewhere
