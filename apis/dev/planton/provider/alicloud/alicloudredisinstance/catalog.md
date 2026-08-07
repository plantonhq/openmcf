# AliCloud Redis Instance

Deploys a managed Redis (KVStore) instance on Alibaba Cloud with configurable engine versions, cluster sharding, multi-zone high availability, read replicas, and encryption at rest. The instance integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to VSwitches for network placement.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **KVStore Instance** -- an `alicloud_kvstore_instance` with the selected engine version, instance class, and network placement in the specified VSwitch
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **A VSwitch** in the target region and availability zone. The Redis instance inherits its VPC and AZ from the VSwitch. Provide the VSwitch ID directly or reference an AliCloudVswitch Cloud Resource via ValueFromRef.
- **A secondary AZ VSwitch** (optional) -- for multi-zone HA, set `secondaryZoneId` to a different AZ so the standby node runs in a separate zone.
- **A KMS key** (optional) -- required when enabling Transparent Data Encryption (`tdeStatus: Enabled`). Once TDE is enabled, it cannot be disabled.

## Deploy

### Console

Open the deployment store, find **AliCloud Redis Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Single** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudRedisInstance
metadata:
  name: app-cache
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  instanceClass: redis.master.small.default
  password: "${REDIS_PASSWORD}"
  vswitchId:
    value: "vsw-abc123"
```

```shell
planton apply -f redis-instance.yaml
```

This creates a Redis 7.0 master-slave instance with VPC password authentication, PostPaid billing, and the default `redis.master.small.default` instance class. Sharding, read replicas, SSL, and TDE are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Redis instance to a VSwitch deployed in the same InfraPipeline:

```yaml
spec:
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: cache-vswitch
      fieldPath: status.outputs.vswitch_id
```

The InfraPipeline resolves the dependency graph, deploys the VSwitch first, then provisions the Redis instance with the resolved VSwitch ID.

## Key Configuration

These are the most important decisions when configuring a Redis instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Instance class and sharding** -- The `instanceClass` determines memory capacity and throughput per node. For workloads that exceed a single node's capacity, set `shardCount` to a value greater than 1 to create a cluster-mode Redis instance with data distributed across shards. Cluster mode is recommended for large-scale caching.

**Read scaling** -- Set `readOnlyCount` (1-9) to add read replicas in the primary zone. Read replicas distribute read traffic across multiple nodes, reducing load on the primary. Available for master-slave and cluster instances.

**VPC authentication mode** -- Set `vpcAuthMode` to `Open` (default) to require password authentication for all connections, including those from within the VPC. Set to `Close` to allow password-free connections from within the same VPC -- useful for trusted internal networks.

**Encryption and protection** -- Enable `tdeStatus: Enabled` with `encryptionKey` for Transparent Data Encryption at rest. Enable `sslEnable: Enable` for encrypted client connections. Set `instanceReleaseProtection: true` to prevent accidental deletion via console or API.

**Backup and maintenance** -- Set `backupPeriod` (e.g., `["Monday", "Wednesday", "Friday"]`) and `backupTime` (e.g., `"02:00Z-03:00Z"`) for automated backup scheduling. Set `maintainStartTime` and `maintainEndTime` for the maintenance window when Alibaba Cloud may apply minor version upgrades.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVswitch** | `vswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | Redis instance ID (e.g., `r-xxxxx`) | Monitoring dashboards, audit references |
| `connection_domain` | Intranet (VPC-internal) connection domain | Application connection strings, DNS CNAME records |
| `private_connection_port` | Private connection port (default: `6379`) | Application connection configuration |
| `private_ip` | Private IP address within the VSwitch | Direct IP-based connections, security group rules |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard single instance** -- A master-slave Redis instance with minimal instance class for development, testing, or lightweight caching workloads. PostPaid billing for cost efficiency. Start from the **Standard Single** preset.

**HA cluster** -- A cluster-mode Redis instance with multiple shards and a secondary AZ for cross-zone failover. Suitable for high-throughput caching and session management. Start from the **HA Cluster** preset.

**Production encrypted** -- A production instance with TDE encryption, SSL for client connections, release protection, and backup policies configured. Suitable for workloads subject to data-at-rest encryption requirements. Start from the **Production Encrypted** preset.

## Works With

- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides the VSwitch for VPC and availability zone placement