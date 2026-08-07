---
title: "RDS Instance"
description: "RDS Instance deployment documentation"
icon: "package"
order: 100
componentName: "alicloudrdsinstance"
---

# AliCloud RDS Instance

Deploys a managed relational database instance on Alibaba Cloud with bundled databases, accounts, and fine-grained account privileges. Supports MySQL, PostgreSQL, SQL Server, MariaDB, and PPAS engines through a single component type. The instance integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to VSwitches for network placement.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **RDS Instance** -- an `alicloud_db_instance` with the selected engine, instance class, storage size, and HA category, placed in the specified VSwitch
- **Databases** -- one `alicloud_db_database` per entry in the `databases` list, with engine-appropriate default character sets (e.g., `utf8mb4` for MySQL, `UTF8` for PostgreSQL)
- **Accounts** -- one `alicloud_rds_account` per entry in the `accounts` list, created after all databases are provisioned
- **Account Privileges** -- one `alicloud_db_account_privilege` per privilege entry, granting specific access levels (ReadOnly, ReadWrite, DDLOnly, DMLOnly, DBOwner) on named databases
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **A VSwitch** in the target region and availability zone. The RDS instance inherits its VPC and AZ from the VSwitch. Provide the VSwitch ID directly or reference an AliCloudVswitch Cloud Resource via ValueFromRef.
- **A standby AZ VSwitch** (optional) -- for HighAvailability or Finance category deployments, set `zoneIdSlaveA` to a different AZ for cross-zone failover.
- **A KMS key** (optional) -- for disk encryption (`encryptionKey`) or Transparent Data Encryption (`tdeStatus`). Once TDE is enabled, it cannot be disabled.

## Deploy

### Console

Open the deployment store, find **AliCloud RDS Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **MySQL Basic** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudRdsInstance
metadata:
  name: app-database
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  engine: MySQL
  engineVersion: "8.0"
  instanceType: rds.mysql.s2.large
  instanceStorage: 50
  vswitchId:
    value: "vsw-abc123"
  databases:
    - name: myapp
  accounts:
    - accountName: app_user
      accountPassword: "${DB_PASSWORD}"
      privileges:
        - databaseNames: [myapp]
          privilege: ReadWrite
```

```shell
planton apply -f rds-instance.yaml
```

This creates a MySQL 8.0 HighAvailability instance with 50 GB cloud storage, one database, one account with ReadWrite access, and Postpaid billing. SSL, TDE, and deletion protection are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the RDS instance to a VSwitch deployed in the same InfraPipeline:

```yaml
spec:
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: db-vswitch
      fieldPath: status.outputs.vswitch_id
```

The InfraPipeline resolves the dependency graph, deploys the VSwitch first, then provisions the RDS instance with the resolved VSwitch ID.

## Key Configuration

These are the most important decisions when configuring an RDS instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Engine choice** -- Set `engine` to one of MySQL, PostgreSQL, SQLServer, MariaDB, or PPAS, and `engineVersion` to a supported version string (e.g., `"8.0"` for MySQL, `"16.0"` for PostgreSQL). The engine determines valid instance types, character sets, and parameter names. The engine choice is immutable after creation.

**HA category** -- Set `category` to control the deployment architecture. `HighAvailability` (default) provisions a primary + standby pair with automatic failover. `Basic` is a single node suitable for dev/test. `AlwaysOn` enables SQL Server Availability Groups. `Finance` deploys a three-node enterprise cluster. Set `zoneIdSlaveA` to a different AZ for cross-zone high availability.

**Bundled databases and accounts** -- The `databases`, `accounts`, and `accounts[].privileges` fields let you declare the full database topology in a single manifest. Accounts are created after databases, and privileges are granted after both exist. Use `accountType: Super` for administrative access or `accountType: Normal` with explicit privilege grants for application accounts.

**Encryption** -- Two independent encryption mechanisms: `encryptionKey` encrypts the underlying storage disk with a customer-managed KMS key, while `tdeStatus: Enabled` turns on Transparent Data Encryption at the engine level. TDE is irreversible once enabled. Both can be used together.

**Billing and protection** -- Defaults to `Postpaid` (pay-as-you-go). Set `instanceChargeType: Prepaid` with `period` for subscription billing. Enable `deletionProtection` to prevent accidental deletion via console or API.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVswitch** | `vswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | RDS instance ID (e.g., `rm-xxxxx`) | Monitoring dashboards, audit references |
| `connection_string` | Intranet (VPC-internal) connection endpoint | Application connection strings, DNS CNAME records |
| `port` | Database service port (e.g., `3306` for MySQL, `5432` for PostgreSQL) | Application connection configuration |
| `database_ids` | Map of database names to their IDs | Resource tracking, audit references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**MySQL for development** -- A Basic-category MySQL 8.0 instance with minimal storage, one database, and one ReadWrite account. Postpaid billing for cost efficiency. Start from the **MySQL Basic** preset.

**PostgreSQL with high availability** -- A HighAvailability PostgreSQL instance with cross-AZ standby, increased storage, and SSL enabled. Start from the **PostgreSQL HA** preset.

**MySQL for production** -- A HighAvailability MySQL 8.0 instance with cloud ESSD storage, deletion protection, maintenance window, and security IP whitelist configured. Start from the **MySQL Production** preset.

## Works With

- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides the VSwitch for VPC and availability zone placement