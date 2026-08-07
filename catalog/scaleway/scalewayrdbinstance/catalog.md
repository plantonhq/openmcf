# Scaleway RDB Instance

Deploys a managed PostgreSQL or MySQL database on Scaleway as a composite resource that bundles the database engine, logical databases, users with per-database privileges, and network ACL rules into a single declarative unit. Configurable high availability with automatic failover, Private Network integration, encryption at rest, automated backups, and flexible storage options. Supports ValueFromRef for Private Network dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **RDB Instance** -- a managed database engine (PostgreSQL or MySQL) with the configured version, node type, HA mode, storage volume, backup schedule, encryption settings, and optional Private Network attachment
- **Logical Databases** -- created only when the `databases` list is populated; each entry is a separate database namespace for application data
- **Database Users** -- created only when the `users` list is populated; each user has a name, password, and optional admin flag
- **User Privileges** -- created only when users have `privileges` entries; each privilege links a user to a database with a specific permission level (readonly, readwrite, all, none)
- **Network ACL** -- created only when `aclRules` is populated; a single resource that atomically replaces all IP-based access rules on the public endpoint
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway Private Network** in the target region for production deployments. RDB instances receive private endpoints when attached to a Private Network. Provide the Private Network UUID directly or reference a ScalewayPrivateNetwork Cloud Resource via ValueFromRef.
- **A supported engine version** -- PostgreSQL 14, 15, or 16 and MySQL 8 are supported. The engine cannot be changed after creation.

## Deploy

### Console

Open the deployment store, find **Scaleway RDB Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production PostgreSQL HA** preset in the [Presets](#presets) tab for a high-availability PostgreSQL instance with Private Network, encryption, and frequent backups.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayRdbInstance
metadata:
  name: app-db
  org: acme-corp
  env: prod
spec:
  region: fr-par
  engine: PostgreSQL-16
  nodeType: DB-DEV-S
  adminUser: admin
  adminPassword: changeme123
```

```shell
planton apply -f scaleway-rdb-instance.yaml
```

This creates a single-node PostgreSQL 16 instance on local SSD storage with default backup settings. No HA, Private Network, databases, users, or ACL rules are configured.

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

The InfraPipeline resolves the dependency graph, deploys the Private Network first, then provisions the RDB instance with the resolved Private Network ID.

## Key Configuration

These are the most important decisions when configuring an RDB instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Engine** -- The `engine` field selects the database engine and major version (e.g., `PostgreSQL-16`, `MySQL-8`). This cannot be changed after creation. Changing the engine requires creating a new instance and migrating data.

**High availability** -- Set `isHaCluster` to true for a synchronous standby replica with automatic failover. Doubles the cost but provides near-zero RPO. Recommended for production; leave false for development.

**Storage** -- The `volumeType` field selects between `lssd` (local SSD, lowest latency), `bssd` (block SSD, 5K IOPS), and `sbs_15k` (block SSD, 15K IOPS). Volume type cannot be changed after creation. The `volumeSizeInGb` field can only be increased, never decreased.

**Networking and ACL** -- When `privateNetworkId` is set, the instance receives a private endpoint. ACL rules control the public endpoint only. For production, combine Private Network connectivity with restrictive ACL rules to minimize attack surface.

**Backup** -- Automated backups are enabled by default. Adjust `backupScheduleFrequencyHours` (1-24) and `backupScheduleRetentionDays` (1-365) based on your RPO requirements. Set `disableBackup` to true only for disposable development instances.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayPrivateNetwork** (optional) | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | Unique identifier of the RDB instance | Read replicas, monitoring tools, management integrations |
| `endpoint_ip` | Public endpoint IPv4 address | Admin access, external integrations (subject to ACL rules) |
| `endpoint_port` | Public endpoint TCP port | Application connection configuration |
| `private_endpoint_ip` | Private Network endpoint IPv4 address | Application connection strings from Private Network resources |
| `private_endpoint_port` | Private Network endpoint TCP port | Application connection configuration |
| `certificate` | TLS CA certificate in PEM format | Client TLS verification via `sslrootcert` (PostgreSQL) or `ssl-ca` (MySQL) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Development PostgreSQL** -- A single-node DB-DEV-S instance running PostgreSQL 16 on local SSD with default daily backups. Minimal cost for development and testing. Start from the **Dev PostgreSQL** preset.

**Production PostgreSQL HA** -- A high-availability DB-GP-XS instance with synchronous standby, Private Network, encryption at rest, 6-hour backup frequency, and 30-day retention. The standard production configuration. Start from the **Production PostgreSQL HA** preset.

**MySQL web application** -- A DB-GP-XS instance running MySQL 8 with Private Network connectivity. Sized for typical web application databases (CMS, e-commerce, SaaS). Start from the **MySQL Web App** preset.

## Works With

- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides network isolation for private database endpoints