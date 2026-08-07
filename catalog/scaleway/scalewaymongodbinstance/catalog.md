# Scaleway MongoDB Instance

Deploys a managed MongoDB instance on Scaleway as a composite resource that bundles the database engine and application-level users with role-based access control into a single declarative unit. Configurable node counts (standalone or 3-node replica set), Private Network integration for production isolation, automated snapshot scheduling, and block storage volumes. Supports ValueFromRef for Private Network dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **MongoDB Instance** -- a managed MongoDB database engine with the configured version, node type, node count (standalone or 3-node replica set), storage volume, and optional Private Network attachment and snapshot scheduling
- **Database Users** -- created only when the `users` list is populated; each user receives role-based access (read, read_write, db_admin) scoped to specific databases or all databases
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway Private Network** in the target region for production deployments. MongoDB instances receive private-only endpoints when attached to a Private Network. Provide the Private Network UUID directly or reference a ScalewayPrivateNetwork Cloud Resource via ValueFromRef.
- **MongoDB is currently available only in fr-par** (Paris). Additional regions may become available in the future.

## Deploy

### Console

Open the deployment store, find **Scaleway MongoDB Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production Replica Set** preset in the [Presets](#presets) tab for a 3-node replica set with Private Network connectivity and snapshot scheduling.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayMongodbInstance
metadata:
  name: app-mongodb
  org: acme-corp
  env: prod
spec:
  region: fr-par
  version: "7.0.12"
  nodeType: MGDB-PLAY2-NANO
  nodeNumber: 1
  adminUser: admin
  adminPassword: changeme123
```

```shell
planton apply -f scaleway-mongodb-instance.yaml
```

This creates a standalone MongoDB 7.0 instance with default 5 GB block storage and a public endpoint. No Private Network, snapshot schedule, or additional users are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the instance to a Private Network deployed in the same InfraPipeline:

```yaml
spec:
  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id
```

The InfraPipeline resolves the dependency graph, deploys the Private Network first, then provisions the MongoDB instance with the resolved Private Network ID.

## Key Configuration

These are the most important decisions when configuring a MongoDB instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Node count** -- The `nodeNumber` field selects standalone (1) or replica set (3) mode. A 3-node replica set provides automatic failover via primary election. There is no 2-node mode -- MongoDB requires an odd number of voting members. Changing between 1 and 3 may destroy and recreate the instance.

**Node type** -- Choose between shared vCPU types (`MGDB-PLAY2-NANO` through `MGDB-PRO2-L`) for cost-optimized workloads and dedicated vCPU types (`MGDB-POP2-2C-8G` through `MGDB-POP2-64C-256G`) for production-grade performance. Node type can be changed after creation.

**Storage** -- The `volumeType` field selects between `sbs_5k` (5,000 IOPS) and `sbs_15k` (15,000 IOPS) block storage. Volume type cannot be changed after creation. The `volumeSizeInGb` field must be a multiple of 5 and can only be increased, never decreased.

**Networking** -- When `privateNetworkId` is set and `enablePublicNetwork` is false (default), the instance is private-only. MongoDB on Scaleway has no IP-based ACL rules -- the public endpoint is accessible from any IP address when enabled.

**Snapshot schedule** -- Enable `enableSnapshotSchedule` with `snapshotScheduleFrequencyHours` and `snapshotScheduleRetentionDays` for automated point-in-time recovery. Not enabled by default.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayPrivateNetwork** (optional) | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | Unique identifier of the MongoDB instance | Monitoring dashboards, snapshot tooling, management automation |
| `public_dns_record` | Public endpoint DNS hostname | Application connection strings when public endpoint is enabled |
| `public_port` | Public endpoint TCP port | Application connection configuration |
| `private_dns_records` | Private Network endpoint DNS hostnames | Application connection strings from Private Network resources |
| `private_ips` | Private Network endpoint IPv4 addresses | Direct private connectivity from applications |
| `private_port` | Private Network endpoint TCP port | Application connection configuration |
| `tls_certificate` | TLS CA certificate in PEM format | Client TLS verification via `tlsCAFile` connection option |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Development standalone** -- A single-node MGDB-PLAY2-NANO instance with public endpoint access and no snapshots. Minimal cost for development, testing, and prototyping. Start from the **Dev Standalone** preset.

**Production replica set** -- A 3-node MGDB-POP2-2C-8G replica set with Private Network connectivity, 12-hour snapshot schedule, and 14-day snapshot retention. Automatic failover for production workloads. Start from the **Production Replica Set** preset.

## Works With

- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides network isolation for private-only MongoDB endpoints