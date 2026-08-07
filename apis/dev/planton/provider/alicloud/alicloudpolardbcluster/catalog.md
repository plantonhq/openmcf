# AliCloud PolarDB Cluster

Deploys an Alibaba Cloud PolarDB cluster with bundled databases, accounts, and account privileges. PolarDB uses a shared-storage, compute-storage-separated architecture supporting MySQL, PostgreSQL, and Oracle compatibility modes. The cluster integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to VSwitches for network placement.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **PolarDB Cluster** -- an `alicloud_polardb_cluster` with the selected engine type, node class, node count, and edition, placed in the specified VSwitch
- **Databases** -- one `alicloud_polardb_database` per entry in the `databases` list, with engine-appropriate default character sets and optional collation settings for PostgreSQL/Oracle modes
- **Accounts** -- one `alicloud_polardb_account` per entry in the `accounts` list, created after all databases are provisioned
- **Account Privileges** -- one `alicloud_polardb_account_privilege` per privilege entry, granting specific access levels (ReadOnly, ReadWrite, DDLOnly, DMLOnly) on named databases
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **A VSwitch** in the target region and availability zone. The PolarDB cluster inherits its VPC and AZ from the VSwitch. Provide the VSwitch ID directly or reference an AliCloudVswitch Cloud Resource via ValueFromRef.
- **A KMS key** (optional) -- required when enabling Transparent Data Encryption (`tdeStatus: Enabled`). TDE encrypts data at rest at the storage level and cannot be disabled once enabled.

## Deploy

### Console

Open the deployment store, find **AliCloud PolarDB Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **MySQL Dev** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudPolardbCluster
metadata:
  name: app-polardb
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  dbType: MySQL
  dbVersion: "8.0"
  dbNodeClass: polar.mysql.x4.large
  vswitchId:
    value: "vsw-abc123"
  databases:
    - dbName: myapp
  accounts:
    - accountName: app_user
      accountPassword: "${DB_PASSWORD}"
      privileges:
        - dbNames: [myapp]
          accountPrivilege: ReadWrite
```

```shell
planton apply -f polardb-cluster.yaml
```

This creates a PolarDB MySQL 8.0 Enterprise Edition cluster with 2 nodes (1 primary + 1 read replica), one database, one account with ReadWrite access, and PostPaid billing. TDE, deletion lock, and SQL audit are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the PolarDB cluster to a VSwitch deployed in the same InfraPipeline:

```yaml
spec:
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: db-vswitch
      fieldPath: status.outputs.vswitch_id
```

The InfraPipeline resolves the dependency graph, deploys the VSwitch first, then provisions the PolarDB cluster with the resolved VSwitch ID.

## Key Configuration

These are the most important decisions when configuring a PolarDB cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Engine type** -- Set `dbType` to MySQL, PostgreSQL, or Oracle, and `dbVersion` to a supported version string (e.g., `"8.0"` for MySQL, `"14"` for PostgreSQL). The engine type determines valid node classes, character sets, and collation settings. This choice is immutable after creation.

**Edition and storage model** -- Set `creationCategory` to control the cluster architecture. `Normal` (default) is Enterprise Edition with shared distributed storage that auto-scales -- you do not provision storage capacity. `SENormal` is Standard Edition with local ESSD storage where you pre-allocate `storageSpace` in GB. `Basic` is a single-node edition for dev/test. `NormalMultimaster` enables multi-master writes (MySQL only).

**Node count and scaling** -- Set `dbNodeCount` to control the number of nodes (1-16, default 2). Node count includes 1 primary and N-1 read replicas that share the same storage layer. Scale read capacity by adding nodes without data migration.

**Bundled databases and accounts** -- The `databases`, `accounts`, and `accounts[].privileges` fields declare the full database topology in a single manifest. PostgreSQL and Oracle modes additionally support `collate` and `ctype` settings per database. Use `accountType: Super` for administrative access or `accountType: Normal` with explicit privilege grants.

**Encryption and audit** -- Enable `tdeStatus: Enabled` with `encryptionKey` for Transparent Data Encryption. TDE is irreversible once enabled. Enable `collectorStatus: Enable` to turn on SQL audit logging for compliance. Set `backupRetentionPolicyOnClusterDeletion` to `ALL`, `LATEST`, or `NONE` to control backup retention when the cluster is deleted.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVswitch** | `vswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | PolarDB cluster ID (e.g., `pc-xxxxx`) | Monitoring dashboards, audit references |
| `connection_string` | Primary endpoint connection string for read-write operations | Application connection strings, DNS CNAME records |
| `port` | Database service port (e.g., `3306` for MySQL, `5432` for PostgreSQL/Oracle) | Application connection configuration |
| `database_ids` | Map of database names to their IDs | Resource tracking, audit references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**MySQL for development** -- A Basic-edition PolarDB MySQL cluster with minimal node class, one database, and one ReadWrite account. PostPaid billing for cost efficiency. Start from the **MySQL Dev** preset.

**MySQL for production** -- An Enterprise Edition MySQL cluster with larger node class, multiple nodes for read scaling, and production-grade settings. Start from the **MySQL Production** preset.

**PostgreSQL for production** -- An Enterprise Edition PostgreSQL cluster with UTF8 character set, suitable for applications requiring PostgreSQL-specific features like JSONB, window functions, and advanced indexing. Start from the **PostgreSQL Production** preset.

## Works With

- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides the VSwitch for VPC and availability zone placement