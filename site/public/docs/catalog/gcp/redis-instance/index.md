---
title: "Redis Instance"
description: "Redis Instance deployment documentation"
icon: "package"
order: 100
componentName: "gcpredisinstance"
---

# GCP Redis Instance

Deploys a Google Cloud Memorystore for Redis instance with configurable tier, replication, persistence, AUTH, transit encryption, and optional CMEK. Supports both standalone (BASIC) and highly-available (STANDARD_HA) configurations with automatic failover and read replicas.

## What Gets Created

When you deploy a GcpRedisInstance resource, Planton provisions:

- **Memorystore Redis Instance** — a fully managed Redis instance in the specified project and region, tagged with organization, environment, and resource labels
- **Primary Endpoint** — a host and port for read/write operations, available immediately after creation
- **Read Replica Endpoint** — created only when tier is `STANDARD_HA` with `readReplicasMode` set to `READ_REPLICAS_ENABLED`, provides a separate endpoint for read-only traffic
- **Maintenance Policy** — configured when a maintenance window is specified, schedules a weekly 1-hour window for GCP-managed updates
- **Persistence (RDB Snapshots)** — configured when persistence is enabled, periodically writes data to disk for durability across restarts

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** where the Redis instance will be created
- **A VPC network** if connecting to a non-default network (referenced via `authorizedNetwork`)
- **A Cloud KMS key** if using customer-managed encryption at rest (CMEK)

## Quick Start

Create a file `redis.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRedisInstance
metadata:
  name: my-redis
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpRedisInstance.my-redis
spec:
  projectId:
    value: my-gcp-project
  instanceName: my-redis
  region: us-central1
  tier: BASIC
  memorySizeGb: 1
```

Deploy:

```shell
planton apply -f redis.yaml
```

This creates a standalone 1 GiB Redis instance in `us-central1` using the default VPC network and the latest supported Redis version.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `projectId` | `StringValueOrRef` | GCP project where the Redis instance will be created. Can reference a GcpProject resource via `valueFrom`. | Required |
| `instanceName` | `string` | Name of the Redis instance. Becomes the GCP resource name. Immutable after creation. | Lowercase letters, numbers, hyphens; 2–40 characters; must start with a letter and end with a letter or number |
| `region` | `string` | GCP region for the instance (e.g., `us-central1`). | Required |
| `tier` | `string` | Service tier. `BASIC` for standalone, `STANDARD_HA` for primary + replica with automatic failover. | `BASIC` or `STANDARD_HA` |
| `memorySizeGb` | `int` | Memory size in GiB for the Redis instance. | Minimum 1 |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `redisVersion` | `string` | Latest supported | Redis engine version (e.g., `REDIS_7_0`, `REDIS_7_2`, `REDIS_6_X`). |
| `displayName` | `string` | — | Human-readable display name for the instance. |
| `locationId` | `string` | GCP-selected | Zone within the region for the primary node. Immutable after creation. |
| `alternativeLocationId` | `string` | GCP-selected | Zone for the `STANDARD_HA` replica; must differ from `locationId`. Pin both to co-locate the cache with zonal workloads. Immutable after creation. |
| `authorizedNetwork` | `StringValueOrRef` | Default network | VPC network the instance connects to. Immutable after creation. Can reference a GcpVpcNetwork resource via `valueFrom`. |
| `connectMode` | `string` | `DIRECT_PEERING` | How the instance connects to the VPC. `DIRECT_PEERING` or `PRIVATE_SERVICE_ACCESS`. Immutable after creation. |
| `reservedIpRange` | `string` | GCP-selected | For `DIRECT_PEERING`: a `/29` CIDR block (e.g., `10.0.0.0/29`), non-overlapping. For `PRIVATE_SERVICE_ACCESS`: the NAME of an allocated range on the private services access connection. Immutable after creation. |
| `secondaryIpRange` | `string` | — | Additional range for node placement — required when enabling read replicas on an existing instance. A `/28` CIDR or `"auto"` (direct peering), or an allocated range name (private service access). Mutable. |
| `authEnabled` | `bool` | `false` | Enables Redis AUTH. When `true`, GCP generates a random AUTH string exported in stack outputs. |
| `transitEncryptionMode` | `string` | `DISABLED` | TLS encryption for client traffic. `DISABLED` or `SERVER_AUTHENTICATION`. Immutable after creation. |
| `redisConfigs` | `map<string,string>` | `{}` | Redis configuration parameters (e.g., `maxmemory-policy`, `notify-keyspace-events`). |
| `maintenanceWindow.day` | `string` | — | Day of week for the maintenance window (`MONDAY` through `SUNDAY`). |
| `maintenanceWindow.hour` | `int` | — | Hour of day (0–23, UTC) when the maintenance window starts. |
| `maintenanceWindow.minute` | `int` | `0` | Minute of the hour (0–59, UTC) — pins the window start to the exact minute. |
| `maintenanceVersion` | `string` | GCP rollout | Self-service maintenance version. Set to a newer available version to apply a patch on your schedule. |
| `readReplicasMode` | `string` | `READ_REPLICAS_DISABLED` | `READ_REPLICAS_DISABLED` or `READ_REPLICAS_ENABLED`. Requires `STANDARD_HA` tier. |
| `replicaCount` | `int` | `0` | Number of read replicas (1–5). Requires `STANDARD_HA` tier with `readReplicasMode` set to `READ_REPLICAS_ENABLED`. |
| `persistenceConfig.persistenceMode` | `string` | — | `DISABLED` or `RDB`. RDB enables periodic snapshots. Only meaningful for `STANDARD_HA` tier. |
| `persistenceConfig.rdbSnapshotPeriod` | `string` | — | Snapshot frequency when mode is `RDB`. One of `ONE_HOUR`, `SIX_HOURS`, `TWELVE_HOURS`, `TWENTY_FOUR_HOURS`. |
| `persistenceConfig.rdbSnapshotStartTime` | `string` | creation time | RFC3339 UTC anchor all snapshots align to — place snapshot I/O in a low-traffic window. |
| `customerManagedKey` | `StringValueOrRef` | Google-managed | Cloud KMS key for encryption at rest (CMEK). Format: `projects/{p}/locations/{l}/keyRings/{kr}/cryptoKeys/{k}`. Immutable after creation. Can reference a GcpKmsKey resource via `valueFrom`. |

## Examples

### High-Availability Instance with AUTH

A `STANDARD_HA` instance with Redis AUTH enabled for production workloads:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRedisInstance
metadata:
  name: prod-cache
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpRedisInstance.prod-cache
spec:
  projectId:
    value: my-gcp-project
  instanceName: prod-cache
  region: us-central1
  tier: STANDARD_HA
  memorySizeGb: 5
  redisVersion: REDIS_7_2
  authEnabled: true
  maintenanceWindow:
    day: SUNDAY
    hour: 4
    minute: 30
```

### Read Replicas with Persistence

A `STANDARD_HA` instance with read replicas and RDB snapshots for high-throughput, durable workloads:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRedisInstance
metadata:
  name: analytics-redis
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpRedisInstance.analytics-redis
spec:
  projectId:
    value: my-gcp-project
  instanceName: analytics-redis
  region: europe-west1
  tier: STANDARD_HA
  memorySizeGb: 16
  redisVersion: REDIS_7_0
  readReplicasMode: READ_REPLICAS_ENABLED
  replicaCount: 3
  persistenceConfig:
    persistenceMode: RDB
    rdbSnapshotPeriod: SIX_HOURS
    rdbSnapshotStartTime: "2026-01-01T03:00:00Z"
  redisConfigs:
    maxmemory-policy: allkeys-lru
```

### Private Network with TLS and CMEK

A locked-down instance using Private Service Access, transit encryption, and customer-managed encryption keys:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRedisInstance
metadata:
  name: secure-redis
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpRedisInstance.secure-redis
spec:
  projectId:
    value: my-gcp-project
  instanceName: secure-redis
  region: us-east1
  tier: STANDARD_HA
  memorySizeGb: 8
  authorizedNetwork:
    valueFrom:
      kind: GcpVpcNetwork
      name: my-vpc
      fieldPath: status.outputs.network_self_link
  connectMode: PRIVATE_SERVICE_ACCESS
  reservedIpRange: managed-services-range # name of the GcpGlobalAddress VPC_PEERING allocation
  authEnabled: true
  transitEncryptionMode: SERVER_AUTHENTICATION
  customerManagedKey:
    valueFrom:
      kind: GcpKmsKey
      name: redis-cmek
      fieldPath: status.outputs.key_id
  maintenanceWindow:
    day: WEDNESDAY
    hour: 2
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `host` | `string` | Hostname or IP address of the primary Redis endpoint |
| `port` | `int` | Port number of the primary Redis endpoint (typically 6379) |
| `read_endpoint` | `string` | Hostname or IP address of the read replica endpoint. Only populated when `STANDARD_HA` tier with read replicas enabled. |
| `read_endpoint_port` | `int` | Port number of the read replica endpoint. Only populated when `STANDARD_HA` tier with read replicas enabled. |
| `current_location_id` | `string` | Zone where the Redis primary is currently running. May change after a failover event. |
| `auth_string` | `string` | Redis AUTH string for client authentication. Only populated when `authEnabled` is `true`. Treat as a secret. |
| `server_ca_certs` | `string[]` | PEM-encoded CA certificates protecting the server endpoint. Only populated when transit encryption is on — clients install these as trust anchors. |
| `persistence_iam_identity` | `string` | Service identity (`serviceAccount:<email>`) used by import/export — grant it access to GCS buckets used for RDB import/export. |
| `effective_reserved_ip_range` | `string` | The CIDR range actually in use, whether explicitly reserved or GCP-selected. |
| `instance_name` | `string` | Name of the Redis instance in GCP — the handle API callers and automation address the instance by. |
| `region` | `string` | Plain region name hosting the instance (e.g., `us-central1`). |

## Related Components

- [GcpVpcNetwork](/docs/catalog/gcp/vpc) — provides the VPC network for instance connectivity
- [GcpGlobalAddress](/docs/catalog/gcp/global-address) — reserves the VPC_PEERING range named by `reservedIpRange` in PRIVATE_SERVICE_ACCESS mode
- [GcpServiceNetworkingConnection](/docs/catalog/gcp/service-networking-connection) — the private services access peering PRIVATE_SERVICE_ACCESS instances ride on
- [GcpKmsKey](/docs/catalog/gcp/kms-key) — provides the Cloud KMS key for customer-managed encryption at rest
- [GcpKmsKeyRing](/docs/catalog/gcp/kms-key-ring) — manages the key ring containing the KMS key
- [GcpFirewallRule](/docs/catalog/gcp/firewall-rule) — controls network access to the VPC where the instance resides
