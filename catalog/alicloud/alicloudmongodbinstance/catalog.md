# AliCloud MongoDB Instance

Deploys a managed MongoDB replica-set instance on Alibaba Cloud with configurable replication factors, multi-zone high availability across three AZs, read-only replicas, and dual encryption options. The instance integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to VSwitches for network placement.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **MongoDB Instance** -- an `alicloud_mongodb_instance` in replica-set mode with the selected engine version, instance class, storage size, and replication factor, placed in the specified VSwitch
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **A VSwitch** in the target region and availability zone. The MongoDB instance inherits its VPC and AZ from the VSwitch. Provide the VSwitch ID directly or reference an AliCloudVswitch Cloud Resource via ValueFromRef.
- **VSwitches in two additional AZs** (optional) -- for three-zone high availability, set `secondaryZoneId` and `hiddenZoneId` to different AZs so each replica-set member runs in a separate zone.
- **A KMS key** (optional) -- for Transparent Data Encryption (`tdeStatus`) or cloud disk encryption (`encrypted`). These two encryption approaches are mutually exclusive.

## Deploy

### Console

Open the deployment store, find **AliCloud MongoDB Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Development** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudMongodbInstance
metadata:
  name: app-mongodb
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  engineVersion: "7.0"
  dbInstanceClass: dds.mongo.mid
  dbInstanceStorage: 20
  accountPassword: "${MONGODB_PASSWORD}"
  vswitchId:
    value: "vsw-abc123"
```

```shell
planton apply -f mongodb-instance.yaml
```

This creates a MongoDB 7.0 replica-set instance with 3 nodes (default replication factor), WiredTiger storage engine, 20 GB storage, and PostPaid billing. Multi-zone HA, encryption, and read-only replicas are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the MongoDB instance to a VSwitch deployed in the same InfraPipeline:

```yaml
spec:
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: db-vswitch
      fieldPath: status.outputs.vswitch_id
```

The InfraPipeline resolves the dependency graph, deploys the VSwitch first, then provisions the MongoDB instance with the resolved VSwitch ID.

## Key Configuration

These are the most important decisions when configuring a MongoDB instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Replication factor** -- Set `replicationFactor` to 1, 3, 5, or 7 to control the number of replica-set nodes. Higher values increase read capacity and fault tolerance. The default of 3 provides a primary, secondary, and hidden node. Add `readonlyReplicas` (0-5) for additional read scaling without affecting the replication factor.

**Multi-zone high availability** -- For production deployments, set `zoneId`, `secondaryZoneId`, and `hiddenZoneId` to three different availability zones. This distributes the primary, secondary, and hidden nodes across AZs, so the replica set survives a full zone failure.

**Encryption** -- Two mutually exclusive encryption approaches: `tdeStatus: enabled` with `encryptionKey` enables Transparent Data Encryption at the engine level (irreversible once enabled), while `encrypted: true` with `cloudDiskEncryptionKey` encrypts the underlying cloud disk at the infrastructure layer. Choose TDE for compliance scenarios requiring engine-level encryption; choose cloud disk encryption for simpler infrastructure-level protection.

**Storage engine and type** -- `storageEngine` defaults to WiredTiger (recommended for most workloads). Set `storageType` to control the disk tier: `cloud_essd1` through `cloud_essd3` offer increasing IOPS levels, `cloud_auto` auto-scales IOPS, and `local_ssd` provides lowest latency with local NVMe disks. Set `provisionedIops` for explicit IOPS allocation on cloud storage types.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVswitch** | `vswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | MongoDB instance ID (e.g., `dds-xxxxx`) | Monitoring dashboards, audit references |
| `replica_set_name` | Replica set name for MongoDB connection strings | Application connection strings (e.g., `?replicaSet=<name>` for automatic failover) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Development instance** -- A small replica-set instance with minimal storage and PostPaid billing for cost-efficient development and testing. Start from the **Development** preset.

**Production with high availability** -- A three-zone HA deployment with higher instance class, increased storage, and backup policies configured for production workloads. Start from the **Production HA** preset.

**Encrypted for compliance** -- A production instance with TDE or cloud disk encryption enabled, deletion protection, and SSL for client connections. Suitable for workloads subject to data-at-rest encryption requirements. Start from the **Encrypted Compliance** preset.

## Works With

- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides the VSwitch for VPC and availability zone placement