# Database Cluster on DigitalOcean

Deploys a managed database cluster on DigitalOcean supporting PostgreSQL, MySQL, Redis, and MongoDB engines with configurable node sizing, replication, optional VPC placement, and custom storage. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for VPC dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Database Cluster** -- a managed database cluster in the specified region with the configured engine, version, node size, and node count (1-3 nodes for primary replication)
- **VPC Network Attachment** -- configured only when `vpc` is provided; places the cluster within a private network for secure internal access
- **Custom Storage** -- configured only when `storageGib` is provided; overrides the default storage allocation for the chosen node size
- **DigitalOcean Tags** -- resource metadata tags (organization, environment, resource kind) applied automatically for tracking

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

**Engine selection** -- The `engine` field accepts `pg` (PostgreSQL), `mysql`, `redis`, or `mongodb`. Each engine has its own set of supported versions. The `engineVersion` field takes a major version string (e.g., `"16"` for PostgreSQL, `"8"` for MySQL, `"7"` for Redis).

**Node count and high availability** -- Set `nodeCount` between 1 and 3. A 3-node cluster provides automatic failover -- when the primary fails, a standby is promoted. Single-node clusters have no redundancy and are suited only for development.

**Node sizing** -- The `sizeSlug` field sets CPU and memory per node (e.g., `"db-s-1vcpu-1gb"` for development, `"db-s-2vcpu-4gb"` for production). Scale up for heavier query workloads or larger datasets.

**VPC placement** -- The `vpc` field places the cluster in a private network so database traffic stays off the public internet. Required for production. When omitted, the cluster uses default networking.

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
| `connection_uri` | Full connection URI including credentials and database name | Application database configuration |
| `host` | Hostname or IP at which the cluster is accessible | Application connection strings |
| `port` | Network port the cluster listens on | Application connection strings |
| `database_user` | Default database user name | Application authentication |
| `database_password` | Default database user password | Application authentication (store as secret) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production PostgreSQL HA** -- 3-node PostgreSQL 16 cluster with VPC isolation for automatic failover and secure private access. Suitable for mission-critical applications. Start from the **PostgreSQL HA** preset.

**Development PostgreSQL** -- Single-node PostgreSQL 16 on the smallest instance, no VPC. Minimal cost for development, CI/CD test databases, and staging. Start from the **PostgreSQL Dev** preset.

**Redis cache** -- Single-node Redis 7 with VPC placement for low-latency caching, session storage, and pub/sub messaging. Start from the **Redis** preset.

## Works With

- [**DigitalOcean VPC**](/cloud-catalog/digital-ocean-vpc) -- provides the private network for database cluster placement