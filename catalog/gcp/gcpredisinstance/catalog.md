# GCP Redis Instance

Deploys a fully managed Memorystore for Redis instance with configurable tier (BASIC or STANDARD_HA), VPC networking, Redis AUTH, TLS encryption, RDB persistence, read replicas, customer-managed encryption keys (CMEK), and maintenance windows. The instance connects to a VPC via direct peering or Private Service Access and supports Redis versions 6.x through 7.2. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects, VPCs, and KMS keys.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Memorystore Redis Instance** -- a managed Redis instance in the specified GCP project and region, configured with the chosen tier, memory size, and Redis version
- **VPC Network Attachment** -- created only when `authorizedNetwork` is specified; connects the instance to the given VPC via direct peering (default) or Private Service Access
- **Redis AUTH** -- created only when `authEnabled` is true; GCP generates and auto-rotates an AUTH string exported in stack outputs
- **TLS Encryption** -- created only when `transitEncryptionMode` is set to SERVER_AUTHENTICATION; enables client-to-server TLS verification
- **RDB Persistence** -- created only when `persistenceConfig` is specified with mode RDB; configures periodic snapshots at the specified interval
- **Read Replicas** -- created only when `readReplicasMode` is READ_REPLICAS_ENABLED with STANDARD_HA tier; provisions 1-5 read replicas with a dedicated read endpoint
- **CMEK Encryption** -- created only when `customerManagedKey` is specified; encrypts data at rest with the provided Cloud KMS key
- **Maintenance Policy** -- created only when `maintenanceWindow` is specified; schedules a weekly 1-hour maintenance window on the chosen day and hour
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the Redis instance will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **A VPC network** (if using private connectivity) for the instance to attach to. The VPC can use direct peering (default) or Private Service Access. Provide the network self-link directly or reference a GcpVpcNetwork Cloud Resource via ValueFromRef.
- **Redis API** (`redis.googleapis.com`) enabled in the target project.

## Deploy

### Console

Open the deployment store, find **GCP Redis Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic Cache** preset in the [Presets](#presets) tab to pre-populate a minimal development instance.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpRedisInstance
metadata:
  name: app-cache
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  instanceName: app-cache
  region: us-central1
  tier: STANDARD_HA
  memorySizeGb: 5
```

```shell
planton apply -f redis-instance.yaml
```

This creates a 5 GB STANDARD_HA Redis instance with automatic failover, no AUTH, no TLS, and no persistence. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Redis instance to a GCP project and VPC deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-project
      fieldPath: status.outputs.project_id
  authorizedNetwork:
    valueFrom:
      kind: GcpVpcNetwork
      name: production-vpc
      fieldPath: status.outputs.network_self_link
```

The InfraPipeline resolves the dependency graph, deploys the project and VPC first, then provisions the Redis instance with private VPC connectivity.

## Key Configuration

These are the most important decisions when configuring a Redis instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Tier selection** -- BASIC provides a single-node instance with no replication and no SLA -- suitable for caching in development or non-critical workloads. STANDARD_HA provides a primary with automatic failover to a replica in a different zone, backed by a 99.9% availability SLA. The tier is immutable after creation.

**Authentication and encryption** -- Enable `authEnabled` to require clients to present a GCP-generated AUTH string (exported in outputs). Set `transitEncryptionMode: SERVER_AUTHENTICATION` for TLS. Both add security layers beyond VPC network controls. Transit encryption mode is immutable after creation.

**Read replicas** -- Available only with STANDARD_HA. Set `readReplicasMode: READ_REPLICAS_ENABLED` and `replicaCount` (1-5) to expose a dedicated read endpoint for scaling read-heavy workloads. Applications direct write traffic to the primary endpoint and read traffic to the read endpoint.

**Persistence** -- Configure `persistenceConfig` with `persistenceMode: RDB` and a snapshot period (ONE_HOUR, SIX_HOURS, TWELVE_HOURS, or TWENTY_FOUR_HOURS) for durability across restarts and failovers. Only meaningful with STANDARD_HA tier.

**Instance sizing** -- `memorySizeGb` sets the total memory available for data (minimum 1 GB). For STANDARD_HA, the replica consumes additional memory that is managed by GCP and not counted against your allocation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpVpcNetwork** (optional) | `authorizedNetwork` | `status.outputs.network_self_link` |
| **GcpKmsKey** (optional) | `customerManagedKey` | `status.outputs.key_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `host` | Primary endpoint hostname or IP address | Application connection strings, Kubernetes ConfigMaps |
| `port` | Primary endpoint port (typically 6379) | Application connection configuration |
| `read_endpoint` | Read replica endpoint (only when read replicas enabled) | Read-only application connections for load distribution |
| `read_endpoint_port` | Read replica endpoint port (only when read replicas enabled) | Read-only connection configuration |
| `current_location_id` | Zone where the primary is currently running | Monitoring, failover tracking |
| `auth_string` | Redis AUTH string (only when auth enabled) | Application secrets, Kubernetes Secrets |
| `server_ca_certs` | PEM CA certificates (only with TLS SERVER_AUTHENTICATION) | Client trust anchors for the TLS handshake |
| `persistence_iam_identity` | Service-account identity used by import/export operations | Granting Cloud Storage access for RDB import/export |
| `effective_reserved_ip_range` | The CIDR range actually in use (set or auto-selected) | Planning non-overlapping address space |
| `instance_name` | The instance's GCP resource name | API calls and automation addressing the instance |
| `region` | Region hosting the instance | Constructing the instance's locations path |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Basic cache** -- A BASIC tier instance with 1 GB memory for development, testing, or lightweight caching. No authentication, no TLS, no VPC attachment. Optimized for cost and simplicity. Start from the **Basic Cache** preset.

**HA production** -- A STANDARD_HA instance with 5 GB memory, Redis AUTH, TLS encryption, RDB persistence (12-hour snapshots), and a scheduled maintenance window. Suitable for production caching and session storage with 99.9% availability SLA. Start from the **HA Production** preset.

**HA with read replicas** -- A STANDARD_HA instance with 10 GB memory, three read replicas, AUTH, TLS, 6-hour RDB snapshots, and customer-managed encryption (CMEK). Designed for read-heavy workloads like leaderboards, dashboards, and counters that benefit from a dedicated read endpoint. Start from the **HA Read Replicas** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the Redis instance is created
- [**GCP VPC Network**](/cloud-catalog/gcp-vpc-network) -- provides the VPC network for private connectivity via direct peering or Private Service Access
- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- provides the Cloud KMS key for customer-managed encryption at rest