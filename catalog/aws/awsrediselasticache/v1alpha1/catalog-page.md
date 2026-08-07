# AWS Redis ElastiCache

Deploys an AWS ElastiCache replication group running Redis or Valkey, supporting both non-clustered mode (single primary with up to 5 read replicas) and clustered mode (data partitioned across multiple shards with optional replicas per shard). The AWS identifier is `metadata.name`. The component manages optional or bring-your-own subnet and parameter groups, encryption, authentication, restore sources, logging, and snapshot configuration.

## What Gets Created

When you deploy an AwsRedisElasticache resource, Planton provisions:

- **ElastiCache Replication Group** — an `aws_elasticache_replication_group` running Redis or Valkey with the specified topology, node type, and engine version. Named from `metadata.name`.
- **Subnet Group** — created when `subnetIds` are provided, or attach an existing group via `subnetGroupName`
- **Custom Parameter Group** — created when `parameters` and `parameterGroupFamily` are provided, or attach an existing group via `parameterGroupName`

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **VPC subnets** for in-VPC deployments — provide at least two subnets in different Availability Zones for multi-AZ
- **A security group** allowing inbound traffic on the Redis port (default 6379)
- **A KMS key** if using customer-managed encryption at rest
- **An ACM certificate or TLS-capable client** if enabling transit encryption

## Quick Start

Create a file `redis.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedisElasticache
metadata:
  name: my-redis
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsRedisElasticache.my-redis
spec:
  region: us-west-2
  engine: redis
  engineVersion: "7.1"
  description: Development Redis cache
  nodeType: cache.t3.micro
  numCacheClusters: 1
  subnetIds:
    - subnet-0a1b2c3d4e5f00001
    - subnet-0a1b2c3d4e5f00002
  securityGroupIds:
    - sg-0a1b2c3d4e5f00001
```

Deploy:

```shell
planton apply -f redis.yaml
```

This creates a single-node Redis 7.1 cluster (non-clustered mode) in the specified subnets.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the ElastiCache cluster will be created. Example: `us-west-2`, `eu-west-1`. | Required |
| `engine` | `string` | Cache engine. Values: `redis`, `valkey`. Leave empty when joining a global datastore. | Required unless `globalReplicationGroupId` is set |
| `description` | `string` | Human-readable description for the replication group. | Required by AWS |
| `nodeType` | `string` | ElastiCache node type determining CPU, memory, and network capacity. Examples: `cache.t3.micro`, `cache.r7g.large`, `cache.r6gd.xlarge`. | Required unless `globalReplicationGroupId` is set |
| `numCacheClusters` | `int` | Total node count (primary + replicas) for non-clustered mode. Mutually exclusive with `numNodeGroups`. | 1–6 when set |
| `numNodeGroups` | `int` | Shard count for clustered mode. Mutually exclusive with `numCacheClusters`. Forbidden when joining a global datastore. | Must be > 0 when set |

Exactly one of `numCacheClusters` or `numNodeGroups` must be provided to select the topology mode — except a global-datastore secondary, where topology is inherited.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `engineVersion` | `string` | Provider default | Engine version. Examples: `7.1`, `7.0`, `6.2` for Redis; `7.2`, `8.0` for Valkey. Leave empty when joining a global datastore. |
| `port` | `int` | `6379` | Port for client connections. **ForceNew** — changing this destroys and recreates the cluster. Range: 1–65535. |
| `preferredCacheClusterAzs` | `string[]` | `[]` | AZ placement per cache cluster (non-clustered). Length must match `numCacheClusters`. Mutually exclusive with `nodeGroupConfigurations`. |
| `replicasPerNodeGroup` | `int` | `0` | Read replicas per shard. Only valid when `numNodeGroups` is set. Range: 0–5. |
| `nodeGroupConfigurations` | `object[]` | `[]` | Per-shard AZ placement, replica count, and slot assignment. Requires clustered mode. |
| `durability` | `string` | — | Write-acknowledgement mode: `default`, `async`, `sync`, `disabled`. Valkey 9.0+ clustered only. **ForceNew**. |
| `globalReplicationGroupId` | `string` | — | Join an existing global replication group as secondary. **ForceNew**. |
| `automaticFailoverEnabled` | `bool` | `false` | Promote a replica to primary on failure. Requires `numCacheClusters >= 2` or clustered mode. |
| `multiAzEnabled` | `bool` | `false` | Spread replicas across Availability Zones. Requires `automaticFailoverEnabled` to be `true`. |
| `subnetIds` | `StringValueOrRef[]` | `[]` | Subnet IDs for the ElastiCache subnet group. Mutually exclusive with `subnetGroupName`. |
| `subnetGroupName` | `string` | — | Existing ElastiCache subnet group (bring-your-own). **ForceNew**. |
| `securityGroupIds` | `StringValueOrRef[]` | `[]` | VPC security groups attached to cluster nodes. Can reference `AwsSecurityGroup` via `valueFrom`. |
| `networkType` | `string` | `ipv4` | IP addressing: `ipv4`, `ipv6`, `dual_stack`. **ForceNew**. |
| `ipDiscovery` | `string` | — | DNS discovery address family: `ipv4` or `ipv6`. Meaningful with dual-stack. |
| `atRestEncryptionEnabled` | `bool` | `false` | Encrypt data on disk and in snapshots. **ForceNew** — changing this destroys and recreates the cluster. Recommended: `true`. |
| `transitEncryptionEnabled` | `bool` | `false` | Encrypt all client and replication traffic with TLS. Recommended: `true`. |
| `transitEncryptionMode` | `string` | — | TLS enforcement mode. Values: `preferred` (allows non-TLS), `required` (TLS only). Requires `transitEncryptionEnabled`. |
| `kmsKeyId` | `StringValueOrRef` | — | Customer-managed KMS key ARN for at-rest encryption. **ForceNew**. Can reference `AwsKmsKey` via `valueFrom`. |
| `authToken` | `StringValueOrRef` | — | Redis AUTH password (16–128 printable chars). Requires `transitEncryptionEnabled`. Mutually exclusive with `userGroupIds`. |
| `authTokenUpdateStrategy` | `string` | — | Token change rollout: `ROTATE`, `SET`, or `DELETE`. Requires `authToken`. |
| `userGroupIds` | `StringValueOrRef[]` | `[]` | Redis ACL user groups via `AwsElasticacheUserGroup`. Mutually exclusive with `authToken`. |
| `snapshotArns` | `string[]` | `[]` | S3 ARNs of RDB files to seed from. Mutually exclusive with `snapshotName`. **ForceNew**, create-time only. |
| `snapshotName` | `string` | — | Existing ElastiCache snapshot to restore from. Mutually exclusive with `snapshotArns`. **ForceNew**, create-time only. |
| `maintenanceWindow` | `string` | AWS default | Weekly maintenance window in UTC. Format: `ddd:hh24:mi-ddd:hh24:mi`. Example: `sun:05:00-sun:06:00`. |
| `snapshotRetentionLimit` | `int` | `0` | Days to retain automatic snapshots. `0` disables snapshots. Range: 0–35. |
| `snapshotWindow` | `string` | AWS default | Daily snapshot window in UTC. Format: `hh24:mi-hh24:mi`. Example: `03:00-04:00`. |
| `finalSnapshotIdentifier` | `string` | — | Name for the final snapshot taken on deletion. If omitted, no final snapshot is created. |
| `applyImmediately` | `bool` | `false` | Apply changes immediately instead of during the next maintenance window. May cause brief downtime. |
| `parameterGroupFamily` | `string` | — | Parameter group family. Required when `parameters` is provided. Examples: `redis7`, `redis6.x`, `valkey7`. |
| `parameters` | `object[]` | `[]` | Custom cache parameters applied via a managed parameter group. Mutually exclusive with `parameterGroupName`. |
| `parameters[].name` | `string` | — | Parameter name (e.g., `maxmemory-policy`, `timeout`). Required. |
| `parameters[].value` | `string` | — | Parameter value (e.g., `volatile-lru`, `300`). Required. |
| `parameterGroupName` | `string` | — | Existing parameter group (bring-your-own). Cluster mode requires `.cluster.on` family. |
| `logDeliveryConfigurations` | `object[]` | `[]` | Log delivery configs. At most 2 entries — one per log type. |
| `logDeliveryConfigurations[].destinationType` | `string` | — | Destination type. Values: `cloudwatch-logs`, `kinesis-firehose`. Required. |
| `logDeliveryConfigurations[].destination` | `StringValueOrRef` | — | Destination identifier (log group name or delivery stream name). Required. |
| `logDeliveryConfigurations[].logFormat` | `string` | — | Serialization format. Values: `text`, `json`. Required. |
| `logDeliveryConfigurations[].logType` | `string` | — | Log type. Values: `slow-log`, `engine-log`. Required. |
| `notificationTopicArn` | `StringValueOrRef` | — | SNS topic ARN for cluster event notifications. Can reference `AwsSnsTopic` via `valueFrom`. |
| `autoMinorVersionUpgrade` | `bool` | `false` | Automatically apply minor engine version upgrades during maintenance windows. |
| `dataTieringEnabled` | `bool` | `false` | Move less-frequently-accessed data to SSD. Only on `r6gd` node types. **ForceNew**. |
| `clusterMode` | `string` | — | Migration setting: `enabled`, `compatible`, `disabled`. Online path from non-clustered to clustered. |

## Examples

### Non-Clustered with Encryption and Failover

A 3-node Redis cluster (1 primary + 2 replicas) with encryption and automatic failover across multiple AZs:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedisElasticache
metadata:
  name: session-cache
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRedisElasticache.session-cache
spec:
  region: us-west-2
  engine: redis
  engineVersion: "7.1"
  description: Session cache with HA
  nodeType: cache.r7g.large
  numCacheClusters: 3
  automaticFailoverEnabled: true
  multiAzEnabled: true
  subnetIds:
    - subnet-private-az1
    - subnet-private-az2
    - subnet-private-az3
  securityGroupIds:
    - sg-redis-prod
  atRestEncryptionEnabled: true
  transitEncryptionEnabled: true
  transitEncryptionMode: required
  snapshotRetentionLimit: 7
  snapshotWindow: "03:00-04:00"
  maintenanceWindow: "sun:05:00-sun:06:00"
```

### Clustered Mode with Custom Parameters

A sharded Redis cluster with 3 shards and 2 replicas per shard, custom parameter overrides, and slow-log delivery to CloudWatch:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedisElasticache
metadata:
  name: analytics-cache
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRedisElasticache.analytics-cache
spec:
  region: us-west-2
  engine: redis
  engineVersion: "7.1"
  description: Sharded analytics cache
  nodeType: cache.r7g.xlarge
  numNodeGroups: 3
  replicasPerNodeGroup: 2
  automaticFailoverEnabled: true
  multiAzEnabled: true
  subnetIds:
    - subnet-private-az1
    - subnet-private-az2
    - subnet-private-az3
  securityGroupIds:
    - sg-redis-analytics
  atRestEncryptionEnabled: true
  transitEncryptionEnabled: true
  parameterGroupFamily: redis7
  parameters:
    - name: maxmemory-policy
      value: volatile-lru
    - name: timeout
      value: "300"
  logDeliveryConfigurations:
    - destinationType: cloudwatch-logs
      destination: /aws/elasticache/analytics-cache/slow-log
      logFormat: json
      logType: slow-log
  applyImmediately: true
```

### Valkey with Data Tiering and Foreign Key References

A Valkey cluster using `r6gd` nodes for data tiering, referencing other Planton-managed resources:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedisElasticache
metadata:
  name: tiered-cache
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsRedisElasticache.tiered-cache
spec:
  region: us-west-2
  engine: valkey
  engineVersion: "7.2"
  description: Valkey cache with data tiering
  nodeType: cache.r6gd.xlarge
  numCacheClusters: 3
  automaticFailoverEnabled: true
  multiAzEnabled: true
  dataTieringEnabled: true
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: main-private-subnet-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: main-private-subnet-b
        fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: AwsSecurityGroup
        name: redis-sg
        field: status.outputs.security_group_id
  atRestEncryptionEnabled: true
  kmsKeyId:
    valueFrom:
      kind: AwsKmsKey
      name: data-key
      field: status.outputs.key_arn
  transitEncryptionEnabled: true
  transitEncryptionMode: required
  notificationTopicArn:
    valueFrom:
      kind: AwsSnsTopic
      name: infra-alerts
      field: status.outputs.topic_arn
  snapshotRetentionLimit: 14
  finalSnapshotIdentifier: tiered-cache-final
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `replication_group_id` | `string` | Identifier of the replication group, used in AWS CLI/API calls |
| `primary_endpoint_address` | `string` | Primary (writer) endpoint DNS name for read-write operations |
| `reader_endpoint_address` | `string` | Reader endpoint DNS name distributing reads across replicas. Empty for single-node deployments. |
| `configuration_endpoint_address` | `string` | Configuration endpoint for Cluster Mode Enabled clients. Empty when Cluster Mode is disabled. |
| `arn` | `string` | Amazon Resource Name of the replication group |
| `port` | `int` | Port on which the cluster accepts connections |
| `subnet_group_name` | `string` | Name of the created subnet group. Only populated when `subnetIds` were provided. |
| `parameter_group_name` | `string` | Name of the created parameter group. Only populated when `parameters` were provided. |

## Related Components

- [AwsVpc](/docs/catalog/aws/awsvpc) — provides subnets for cluster placement
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — controls network-level access to the Redis/Valkey endpoint
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — provides a customer-managed key for at-rest encryption
- [AwsSnsTopic](/docs/catalog/aws/awssnstopic) — receives cluster event notifications
- [AwsElasticacheUser](/docs/catalog/aws/awselasticacheuser) — RBAC identity referenced via user groups
- [AwsElasticacheUserGroup](/docs/catalog/aws/awselasticacheusergroup) — collects users; referenced in `userGroupIds`
