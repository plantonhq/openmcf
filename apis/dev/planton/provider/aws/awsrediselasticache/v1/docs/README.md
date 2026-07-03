# AwsRedisElasticache — Architecture and Design

## Overview

AwsRedisElasticache provisions a managed ElastiCache replication group running Redis or Valkey. The component abstracts subnet groups, parameter groups, encryption, topology selection, restore sources, and authentication into a single declarative resource while preserving production-critical features.

The replication group's AWS identifier is taken from `metadata.name` — create-time immutable, so renaming means replacement.

## Why Redis/Valkey Only

ElastiCache supports three engines: Redis, Valkey, and Memcached. This component covers **Redis and Valkey only** because they share the same Terraform resource (`aws_elasticache_replication_group`) and identical configuration surface. Memcached uses a fundamentally different resource (`aws_elasticache_cluster`) with a different topology model (no replication, no persistence, no encryption at rest, no AUTH). A separate `AwsMemcachedElasticache` component covers Memcached.

## Topology Modes

### Non-Clustered (Cluster Mode Disabled)

Set `numCacheClusters` (1–6). One primary handles all writes; additional nodes are read replicas that receive asynchronous replication. The primary endpoint handles writes; the reader endpoint load-balances reads across replicas.

Use `preferredCacheClusterAzs` to pin each cache cluster to a specific Availability Zone (list length must match `numCacheClusters`; first entry hosts the primary).

**Best for:** workloads under ~113 GB, read-heavy patterns, simple client configuration.

### Clustered (Cluster Mode Enabled)

Set `numNodeGroups` (shard count) and optionally `replicasPerNodeGroup` (0–5). Data is hash-slotted across shards. Each shard has its own primary and replicas. The configuration endpoint enables automatic slot discovery for cluster-aware clients.

Use `nodeGroupConfigurations` for per-shard placement — pin each shard's primary and replicas to specific AZs, override replica count per shard, or assign keyspace slots. Mutually exclusive with `preferredCacheClusterAzs`.

**Best for:** multi-terabyte datasets, write-heavy patterns, horizontal scaling needs.

### Cluster Mode Migration

The `clusterMode` field controls the online migration path from non-clustered to clustered:

- `"compatible"` — runs a non-clustered group in cluster-mode-compatible form so clients can migrate to the cluster protocol before the topology actually shards
- `"enabled"` / `"disabled"` — explicit cluster-mode state

Leave empty to let the topology fields decide (the common case).

### Topology Selection CEL

The spec enforces exactly one topology via CEL: `(numCacheClusters > 0) != (numNodeGroups > 0)`, except when joining a global datastore as a secondary — there the shard layout is inherited from the primary and only `numCacheClusters` may be set.

## Global Datastore Membership

Set `globalReplicationGroupId` to join an **existing** global replication group as a secondary region. The secondary inherits engine, engine version, node type, encryption settings, shard layout, and parameter group from the global primary — leave `engine`, `engineVersion`, `nodeType`, `numNodeGroups`, encryption fields, parameter-group fields, and restore sources empty when set.

Global replication groups themselves are created outside this component; this field is the join path. Creating the global resource is a separate deferred kind.

## Durability (Valkey 9.0+, Cluster Mode Enabled)

The `durability` field controls how writes are acknowledged relative to replica propagation. Values: `"default"`, `"async"`, `"sync"`, `"disabled"`. Requires Cluster Mode Enabled (`numNodeGroups`) and engine Valkey 9.0 or later. `"sync"` trades write latency for zero-data-loss failovers. ForceNew — changing it replaces the group.

## Bundled Resources

Following the AwsRdsCluster pattern, the IaC modules conditionally create supporting resources:

- **Subnet Group**: Created when `subnetIds` are provided. Name is sanitized from the resource metadata ID to meet AWS naming constraints (lowercase, alphanumeric, hyphens). Mutually exclusive with `subnetGroupName` (bring-your-own existing subnet group).
- **Parameter Group**: Created when `parameters` are provided with a `parameterGroupFamily`. Uses `name_prefix` to avoid naming collisions during updates. Mutually exclusive with `parameterGroupName` (bring-your-own existing parameter group). Cluster Mode Enabled requires a family `.cluster.on` group.

## Restore (Create-Time)

Two alternative restore sources seed a new replication group at creation:

- **`snapshotArns`** — S3 ARNs of RDB snapshot files (migration path from self-managed Redis)
- **`snapshotName`** — name of an existing ElastiCache snapshot (clone-from-backup path)

Mutually exclusive with each other and forbidden when joining a global datastore (the secondary receives data from the primary). Both are ForceNew and create-time-only.

## Encryption Architecture

- **At Rest**: Enabled via `atRestEncryptionEnabled`. Uses AWS-managed key by default; `kmsKeyId` overrides with a customer-managed KMS key. Both are ForceNew — changing them destroys and recreates the cluster.
- **In Transit**: Enabled via `transitEncryptionEnabled`. The `transitEncryptionMode` field controls enforcement: `"preferred"` allows both TLS and non-TLS connections (for migration); `"required"` enforces TLS for all connections.

## Networking

- **`networkType`**: IP addressing for the cluster — `"ipv4"` (default), `"ipv6"`, or `"dual_stack"`. ForceNew — changing the network type replaces the cluster. Dual-stack requires subnets with both IPv4 and IPv6 CIDRs.
- **`ipDiscovery`**: Which address family DNS discovery returns to clients — `"ipv4"` or `"ipv6"`. Only meaningful alongside dual-stack `networkType`; updates in place, letting clients migrate address families without replacing the cluster.

## Authentication

Two mutually exclusive methods:

1. **AUTH Token** (`authToken`): A single password (16–128 chars) that all clients must provide. Simple but shared-secret model. Requires transit encryption. `authTokenUpdateStrategy` controls how token changes roll out: `"ROTATE"` (zero-downtime — old and new tokens both work until rotation completes), `"SET"` (immediate replacement), or `"DELETE"` (remove token and turn AUTH off).
2. **User Groups** (`userGroupIds`): Redis ACL with fine-grained user permissions via `AwsElasticacheUser` / `AwsElasticacheUserGroup` — AWS's recommended production authentication model.

### RBAC Composition Path

```
AwsElasticacheUser  →  AwsElasticacheUserGroup  →  AwsRedisElasticache
     (WHO / WHAT)            (membership)              (userGroupIds)
```

Adding an application to a cache is a membership edit on the group — the cache resource never changes. Each `userGroupIds` entry is a `StringValueOrRef` defaulting to `AwsElasticacheUserGroup.status.outputs.user_group_id`.

## Infra Chart Composability

### Inputs (StringValueOrRef)

| Field | Default Reference |
|-------|-------------------|
| `subnetIds` | `AwsSubnet.status.outputs.subnet_id` |
| `securityGroupIds` | `AwsSecurityGroup.status.outputs.security_group_id` |
| `kmsKeyId` | `AwsKmsKey.status.outputs.key_arn` |
| `notificationTopicArn` | `AwsSnsTopic.status.outputs.topic_arn` |
| `userGroupIds` | `AwsElasticacheUserGroup.status.outputs.user_group_id` |

### Outputs (for downstream)

| Output | Downstream Use |
|--------|---------------|
| `primary_endpoint_address` + `port` | Application connection config |
| `reader_endpoint_address` | Read-only connection pool |
| `configuration_endpoint_address` | Cluster-aware client config |
| `arn` | IAM policies, metric dimensions |

### Typical DAG Position

```
Layer 0: AwsVpc
Layer 1: AwsSecurityGroup, AwsKmsKey, AwsElasticacheUser
Layer 2: AwsElasticacheUserGroup
Layer 3: AwsRedisElasticache  ← this component
Layer 4: Application configs, AwsLambda event triggers
```

## Deliberately Omitted (v1)

| Feature | Reason |
|---------|--------|
| Outpost fields (`CacheOutpostArns`) | ElastiCache on Outposts is a separate deployment surface; deferred until real demand appears |
| EC2-Classic `security_group_names` | Legacy networking model; VPC security groups via `securityGroupIds` are the supported path |
| Creating a global replication group | Join-as-secondary via `globalReplicationGroupId` is supported; provisioning the global resource itself is a separate deferred kind |
| ElastiCache Serverless | Fundamentally different resource and config model — see `AwsServerlessElasticache` |
| Memcached | Different TF resource, different topology, no replication — see `AwsMemcachedElasticache` |

## References

- [AWS ElastiCache for Redis User Guide](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/WhatIs.html)
- [ElastiCache Replication Group API](https://docs.aws.amazon.com/AmazonElastiCache/latest/APIReference/API_CreateReplicationGroup.html)
- [Redis Cluster Mode Explained](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Replication.Redis-RedisCluster.html)
- [ElastiCache Parameter Groups](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/ParameterGroups.html)
- [ElastiCache Global Datastore](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Redis-Global-Datastore.html)
