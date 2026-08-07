# Database on Civo

Deploys a managed PostgreSQL or MySQL database instance on Civo Cloud with configurable engine version, instance sizing, read replicas, VPC networking, and firewall protection. Integrates with Planton's Provider Connections for Civo credential management and ValueFromRef for wiring VPC and firewall dependencies.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo Database Instance** -- a managed database in the specified region running the chosen engine (PostgreSQL or MySQL) with the configured version, instance size, and node count
- **Read Replicas** -- created only when `replicas` is greater than 0, adding replica nodes to the primary for read scaling and automatic failover
- **Network Attachment** -- the database is placed on the specified VPC network for private connectivity
- **Firewall Binding** -- created only when `firewallIds` is provided, restricting database access to allowed source CIDRs

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A Civo VPC network** in the target region. Provide the network ID directly or reference a CivoVpc Cloud Resource via ValueFromRef.
- **A Civo Firewall** (recommended for production) restricting database port access to the application tier. Provide the firewall ID directly or reference a CivoFirewall Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **Database on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **PostgreSQL Production** preset in the [Presets](#presets) tab for a production-grade database with read replicas.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoDatabase
metadata:
  name: app-db
  org: acme-corp
  env: prod
spec:
  dbInstanceName: app-db
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.medium
  replicas: 2
  networkId:
    value: "abc12345-6789-def0-1234-567890abcdef"
```

```shell
planton apply -f civo-database.yaml
```

This creates a managed PostgreSQL 16 database with 2 read replicas (3 nodes total) on the specified VPC network in Civo's London region. No firewall is attached. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the database to a VPC and firewall deployed in the same InfraPipeline:

```yaml
spec:
  networkId:
    valueFrom:
      kind: CivoVpc
      name: app-network
      fieldPath: status.outputs.network_id
  firewallIds:
    - valueFrom:
        kind: CivoFirewall
        name: db-firewall
        fieldPath: status.outputs.firewall_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC and firewall first, then provisions the database with the resolved values.

## Key Configuration

These are the most important decisions when configuring a Civo database. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Engine selection** -- The `engine` field accepts `postgres` or `mysql`. Choose based on your application's ORM and query requirements. PostgreSQL is the default choice for most workloads; MySQL is available for applications with existing MySQL dependencies.

**Engine version** -- The `engineVersion` field accepts major version strings (e.g., `"16"` for PostgreSQL, `"8.0"` for MySQL). Use the latest stable major version for new projects. Changing the version after deployment requires a migration.

**Replicas** -- The `replicas` field sets the number of read replica nodes (0 to 4). A value of 2 creates a 3-node cluster (1 primary + 2 replicas) with automatic failover and read scaling. Use 0 for development to minimize cost.

**Instance size** -- The `sizeSlug` field sets the instance type (e.g., `g3.db.small` for development, `g3.db.medium` for production). This determines CPU, memory, and base storage for all nodes in the cluster.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CivoVpc** | `networkId` | `status.outputs.network_id` |
| **CivoFirewall** (optional) | `firewallIds` | `status.outputs.firewall_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `database_id` | Unique identifier of the database instance on Civo | Civo API operations, monitoring dashboards |
| `host` | Hostname or IP of the primary database endpoint | Application connection strings, DNS configuration |
| `port` | Network port the database is listening on | Application connection strings |
| `username` | Default database user name | Application database configuration |
| `password_secret_ref` | Reference to the secret containing the default user password | Application secret injection |
| `replica_endpoints` | Hostnames of read replica nodes (empty when no replicas) | Read-heavy query routing, connection pooler configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production PostgreSQL** -- PostgreSQL 16 with 2 read replicas, medium instance sizing, VPC networking, and firewall protection. Covers most production workloads requiring high availability and read scaling. Start from the **PostgreSQL Production** preset.

**Development MySQL** -- MySQL 8.0 single node with small instance sizing and no firewall. Minimal cost for development and testing. Start from the **MySQL Development** preset.

## Works With

- [**Civo VPC**](/cloud-catalog/civo-vpc) -- provides the VPC network for database connectivity
- [**Civo Firewall**](/cloud-catalog/civo-firewall) -- controls inbound access to the database instance