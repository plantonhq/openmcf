---
title: "MongoDB Database"
description: "MongoDB Database deployment documentation"
icon: "package"
order: 100
componentName: "atlasmongodb"
---

# MongoDB Database on Atlas

Deploys a Atlas MongoDB cluster with configurable cluster type, replication topology, cloud provider selection, instance sizing, and backup settings. Supports replica set, sharded, and geo-sharded configurations across AWS, Azure, and GCP. Integrates with Planton's Atlas MongoDB Provider Connection for API key management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Atlas MongoDB Advanced Cluster** -- a cluster with the specified type (REPLICASET, SHARDED, or GEOSHARDED), configured with electable and optional read-only nodes, backup settings, and MongoDB version selection in the target cloud provider region
- **Replication Specification** -- region configuration with electable node specs, read-only node specs (when `readOnlyNodes` > 0), auto-scaling settings, and shard count based on the cluster type

## Before You Deploy

### Planton Setup

- **Atlas MongoDB Provider Connection** -- an active connection in the Connect module with Atlas MongoDB API keys (public key and private key). Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API key authentication.

### Atlas MongoDB Account

- **A Atlas MongoDB project** -- provide the `clusterConfig.projectId` of an existing project where the cluster will be created. Projects act as the organizational container for clusters, users, and network settings.
- **Sufficient Atlas permissions** -- the API key must have Project Owner or equivalent role to create and manage clusters in the target project.
- **Cloud provider selection** -- choose `clusterConfig.providerName` as `AWS`, `GCP`, `AZURE`, or `TENANT` (for M2/M5 shared-tier instances). The cluster is provisioned in the selected provider's infrastructure.
- **Instance sizing** -- select `clusterConfig.providerInstanceSizeName` based on workload requirements (e.g., M10 for development, M30 for production, M40+ for high-throughput workloads).
- **MongoDB version** -- decide on `clusterConfig.mongoDbMajorVersion`. Atlas supports 4.4, 5.0, 6.0, and 7.0 for M10+ clusters. If omitted, Atlas deploys 7.0. M0, M2, and M5 instances default to 5.0.

## Deploy

### Console

Open the deployment store, find **MongoDB Database on Atlas**, and click **Deploy**. The creation wizard walks you through environment and connection configuration, and spec fields including cluster type, cloud provider, instance size, and replication topology.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: atlas.planton.dev/v1
kind: AtlasMongodb
metadata:
  name: app-db
  org: acme-corp
  env: prod
spec:
  clusterConfig:
    projectId: "64a1234567890abcdef12345"
    clusterType: REPLICASET
    electableNodes: 3
    priority: 7
    providerName: AWS
    providerInstanceSizeName: M10
    mongoDbMajorVersion: "7.0"
```

```shell
planton apply -f atlas-mongodb.yaml
```

This creates a 3-node replica set on AWS with M10 instances running MongoDB 7.0. No cloud backup, read-only nodes, or auto-scaling is configured. A Stack Job tracks the provisioning in real time.

For a sharded cluster, change `clusterType` to `SHARDED` and increase the instance size to M30 or higher for production throughput.

## Key Configuration

These are the most important decisions when configuring a Atlas MongoDB cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Cluster type** -- The `clusterConfig.clusterType` field determines the deployment topology. REPLICASET provides a single-region replica set for most workloads. SHARDED distributes data across multiple shards for horizontal scaling. GEOSHARDED enables global data distribution with zone-based sharding. You cannot convert between sharded and replica set types after creation.

**Instance sizing** -- The `clusterConfig.providerInstanceSizeName` field sets the instance tier for all data-bearing servers. Use M10 for development, M30 for production, and M40+ for high-throughput workloads. M0, M2, and M5 are shared-tier instances with limited configurability.

**Replication topology** -- Set `clusterConfig.electableNodes` to 3, 5, or 7 for the total number of electable nodes. The `priority` field (1-7) designates the preferred region -- priority 7 identifies where Atlas places the primary node. Add `readOnlyNodes` to offload analytics queries to dedicated read replicas.

**Cloud backup** -- Set `clusterConfig.cloudBackup` to `true` to enable continuous backup with point-in-time recovery. Disabled by default to minimize cost for development clusters.

**Auto-scaling** -- Set `clusterConfig.autoScalingDiskGbEnabled` to `true` to let Atlas automatically increase storage capacity when usage approaches the provisioned limit. Recommended for production workloads with unpredictable data growth.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Provider-assigned unique ID for the Atlas MongoDB cluster | API operations, resource tracking |
| `bootstrap_endpoint` | Standard connection string in SRV format (mongodb+srv://) | Application database connection configuration |
| `crn` | Cluster identifier for resource identification and API operations | Atlas API operations, monitoring integration |
| `rest_endpoint` | Standard connection string in legacy format (mongodb://) | Legacy application connections, direct host access |

## Common Patterns

No presets are available yet. Configure directly using the fields documented in the [API Explorer](#api-explorer) tab.

## Works With

This component operates independently and does not reference other deployment components.