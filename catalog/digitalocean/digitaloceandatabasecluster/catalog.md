# Database Cluster on DigitalOcean

Deploys a managed database cluster on DigitalOcean supporting PostgreSQL, MySQL, Redis, MongoDB, Kafka, OpenSearch, and Valkey engines with configurable node sizing, replication, optional VPC placement, custom storage with autoscale, maintenance windows, backup-restore provisioning, and engine-specific tuning. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for VPC dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Database Cluster** -- a managed database cluster in the specified region with the configured engine, version, node size, and node count
- **VPC Network Attachment** -- configured only when `vpc` is provided; places the cluster within a private network for secure internal access
- **Custom Storage** -- configured only when `storageGib` is provided; overrides the default storage allocation for the chosen node size, with optional automatic growth via `storageAutoscale`
- **Maintenance Window** -- configured only when `maintenanceWindow` is provided; pins automatic updates to a weekly slot
- **Engine Tuning** -- `sqlMode` on MySQL clusters and `evictionPolicy` on Redis/Valkey clusters
- **DigitalOcean Tags** -- your `tags` plus resource metadata tags (organization, environment, resource kind) applied automatically for tracking

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A target region** where managed databases are available. Not all regions support all engines -- check the DigitalOcean documentation.
- **A VPC network** (recommended for production) in the target region. Provide the VPC UUID directly or reference a DigitalOceanVpc Cloud Resource via ValueFromRef.
- **A valid node size slug** (e.g., `"db-s-2vcpu-4gb"`) -- check available database sizes via `doctl databases options slugs`.

## Deploy

### Console

Open the deployment store, find **Database Cluster on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **PostgreSQL HA** preset in the [Presets](#presets) tab for a production-ready multi-node configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanDatabaseCluster
metadata:
  name: app-db
  org: acme-corp
  env: prod
spec:
  clusterName: app-db
  engine: pg
  engineVersion: "16"
  region: nyc3
  sizeSlug: db-s-2vcpu-4gb
  nodeCount: 3
```

```shell
planton apply -f do-database.yaml
```

This creates a 3-node PostgreSQL 16 cluster in the NYC3 region with no VPC attachment. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the database cluster to a VPC deployed in the same InfraPipeline:

```yaml
spec:
  vpc:
    valueFrom:
      kind: DigitalOceanVpc
      name: app-network
      fieldPath: status.outputs.vpc_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC first, then provisions the database cluster within it.

## Key Configuration

These are the most important decisions when configuring a database cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Engine selection** -- The `engine` field accepts `pg` (PostgreSQL), `mysql`, `redis`, `mongodb`, `kafka`, `opensearch`, or `valkey`. Each engine has its own set of supported versions. The `engineVersion` field takes a major or major.minor version string (e.g., `"16"` for PostgreSQL, `"8"` for MySQL, `"3.5"` for Kafka).

**Node count and high availability** -- Valid node counts are engine-specific: most engines accept 1-3 nodes (a 3-node cluster provides automatic failover), Kafka requires at least 3, and OpenSearch scales to 15. Single-node clusters have no redundancy and are suited only for development.

**Node sizing** -- The `sizeSlug` field sets CPU and memory per node (e.g., `"db-s-1vcpu-1gb"` for development, `"db-s-2vcpu-4gb"` for production). Scale up for heavier query workloads or larger datasets. Changing the size later resizes the cluster in place.

**VPC placement** -- The `vpc` field places the cluster in a private network so database traffic stays off the public internet. Required for production. When omitted, the cluster uses default networking. Cannot be changed after creation.

**Storage and growth** -- `storageGib` sets custom disk space beyond the size slug's default (increase-only), and `storageAutoscale` grows it automatically when usage crosses a threshold. Autoscale is applied by the Terraform provisioner; the Pulumi bridge does not support it yet and rejects it loudly rather than dropping it.

**Engine tuning** -- `sqlMode` applies only to MySQL and `evictionPolicy` only to Redis/Valkey; the manifest is rejected at validation time if they are paired with any other engine, mirroring DigitalOcean's own rules.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanVpc** (optional) | `vpc` | `status.outputs.vpc_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | Unique database cluster identifier (UUID) | DigitalOcean API operations, monitoring |
| `connection_uri` | Full public connection URI including credentials and database name | Application database configuration |
| `host` | Public hostname at which the cluster is accessible | Application connection strings |
| `port` | Network port the cluster listens on | Application connection strings |
| `database_user` | Default database user name | Application authentication |
| `database_password` | Default database user password | Application authentication (store as secret) |
| `private_host` | Private-network hostname (same-VPC access) | In-VPC application connection strings |
| `private_uri` | Private-network connection URI including credentials | In-VPC application configuration |
| `database_name` | Name of the default database | Application configuration |
| `ui_host` / `ui_port` / `ui_uri` / `ui_database` / `ui_user` / `ui_password` | OpenSearch Dashboards connection details (OpenSearch clusters only) | Dashboards access |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production PostgreSQL HA** -- 3-node PostgreSQL 16 cluster with VPC isolation for automatic failover and secure private access. Suitable for mission-critical applications. Start from the **PostgreSQL HA** preset.

**Development PostgreSQL** -- Single-node PostgreSQL 16 on the smallest instance, no VPC. Minimal cost for development, CI/CD test databases, and staging. Start from the **PostgreSQL Dev** preset.

**Redis cache** -- Single-node Redis 7 with VPC placement and an LRU eviction policy for low-latency caching, session storage, and pub/sub messaging. Start from the **Redis** preset.

**Kafka streaming** -- 3-node Kafka 3.5 inside a VPC, the minimum DigitalOcean Kafka topology, for event streaming between services. Start from the **Kafka** preset.

**OpenSearch analytics** -- Single-node OpenSearch 2 with 100 GiB storage for log analytics and full-text search; OpenSearch Dashboards connection details arrive as the `ui_*` outputs. Start from the **OpenSearch** preset.

## Works With

- [**DigitalOcean VPC**](/cloud-catalog/digital-ocean-vpc) -- provides the private network for database cluster placement