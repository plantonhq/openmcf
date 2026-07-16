---
title: "SQL Database"
description: "SQL Database deployment documentation"
icon: "package"
order: 100
componentName: "azuremssqldatabase"
---

# Azure SQL Database

Creates an Azure SQL Database on an AzureMssqlServer logical server. The database is the unit of compute and billing in Azure SQL -- it carries its own SKU (DTU, vCore, serverless, or Hyperscale), storage ceiling, availability posture, backup retention, and encryption, or joins an AzureMssqlElasticPool to share pooled compute.

## What Gets Created

When you deploy an AzureMssqlDatabase resource, Planton provisions:

- **SQL Database** -- an `azurerm_mssql_database` on the referenced server with your chosen SKU (or pool membership), storage, availability, lifecycle mode, encryption, and backup policies. The resource internally orchestrates ARM's TDE, short/long-term retention, and threat-detection sub-APIs.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureMssqlServer** to create the database on (referenced through `serverId`)
- **For pool membership**: an `AzureMssqlElasticPool` on the same server
- **For the database-scoped CMK**: a user-assigned identity with wrap/unwrap access on the Key Vault key

## Quick Start

Create a file `database.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlDatabase
metadata:
  name: orders-db
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureMssqlDatabase.orders-db
spec:
  serverId:
    valueFrom:
      kind: AzureMssqlServer
      name: my-sql
      fieldPath: status.outputs.server_id
  databaseName: orders
  skuName: GP_S_Gen5_1
  autoPauseDelayInMinutes: 60
```

Deploy:

```shell
planton apply -f database.yaml
```

This creates a serverless database that pauses after an hour idle. Connect with the server's `fqdn` output and `Database=orders`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `serverId` | `StringValueOrRef` | The parent logical server. Defaults to referencing an `AzureMssqlServer`'s `server_id` output. Fixed at creation. | Required |
| `databaseName` | `string` | Unique within the server. Changing it replaces the database. | Required, ≤128 chars, no `<>*%&:\/?`, no trailing `.`/space |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `skuName` | `string` | Azure picks (serverless) | DTU (`Basic`, `S0`-`S12`, `P1`-`P15`), vCore (`GP_Gen5_2`, `BC_Gen5_4`, `HS_Gen5_2`), serverless (`GP_S_Gen5_1`), `DW`/`DS`, `ElasticPool`, `Free`. Hyperscale transitions replace the database. |
| `elasticPoolId` | `StringValueOrRef` | -- | The pool to join (requires `skuName: ElasticPool`; forbids a maintenance window). References `AzureMssqlElasticPool`. |
| `maxSizeGb` | `double` | SKU default | 0.1-4096; fractional sizes are legal. Hyperscale ignores it. |
| `collation` | `string` | `SQL_Latin1_General_CP1_CI_AS` | Fixed at creation. |
| `licenseType` | `enum` | LicenseIncluded | `BASE_PRICE` (Hybrid Benefit) / `LICENSE_INCLUDED`. vCore only. |
| `autoPauseDelayInMinutes` | `int32` | -- | Serverless only: 60-10080, or -1 to disable auto-pause. |
| `minCapacity` | `double` | -- | Serverless only: the always-warm vCore floor (0.25-40). |
| `readReplicaCount` | `int32` | -- | Hyperscale only: readable HA replicas (0-4). |
| `readScale` | `bool` | `false` | Premium/BC: read-intent connections go to a secondary. |
| `zoneRedundant` | `bool` | `false` | Spread replicas across availability zones. |
| `ledgerEnabled` | `bool` | `false` | Tamper-evident tables. Fixed at creation. |
| `enclaveType` | `enum` | -- | `VBS` / `DEFAULT_ENCLAVE`. Changing it replaces the database. |
| `maintenanceConfigurationName` | `string` | `SQL_Default` | Region-specific windows (`SQL_EastUS_DB_1`). Unset on pooled databases. |
| `createMode` | `enum` | `DEFAULT` | `COPY`, `SECONDARY`, `ONLINE_SECONDARY`, `POINT_IN_TIME_RESTORE`, `RECOVERY`, `RESTORE`, `RESTORE_LONG_TERM_RETENTION_BACKUP`. Fixed at creation. |
| `creationSourceDatabaseId` | `StringValueOrRef` | -- | The source for copy/secondary/PITR modes; references another database's `database_id`. |
| `secondaryType` | `enum` | `GEO` | `GEO` / `NAMED` (secondary modes only). |
| `restorePointInTime` | `string` | -- | RFC-3339 restore instant (PITR only). |
| `storageAccountType` | `enum` | `GEO_REDUNDANT` | Backup redundancy: `GEO_REDUNDANT`, `GEO_ZONE_REDUNDANT`, `ZONE_REDUNDANT`, `LOCALLY_REDUNDANT`. |
| `userAssignedIdentityIds` | `list` | `[]` | Identities for the database-scoped CMK. |
| `transparentDataEncryptionEnabled` | `bool` | `true` | Hyperscale cannot disable it. |
| `transparentDataEncryptionKeyVaultKeyId` | `StringValueOrRef` | -- | Database-scoped CMK (VERSIONED `AzureKeyVaultKey.key_id`). Requires an attached identity. |
| `transparentDataEncryptionKeyAutomaticRotationEnabled` | `bool` | `false` | Re-encrypt automatically as the key rotates. |
| `import` | `object` | -- | A bacpac import on creation (storage URI + sensitive credentials). Fresh databases only. |
| `shortTermRetentionPolicy` | `object` | 7 days / 12h | The PITR window (1-35 days) + differential cadence (12/24h). |
| `longTermRetentionPolicy` | `object` | -- | Weekly/monthly/yearly horizons as ISO-8601 durations + `weekOfYear`. |
| `threatDetectionPolicy` | `object` | inherits server | Database-scoped Defender policy. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags. |

## Examples

### Pooled Database

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlDatabase
metadata:
  name: tenant-42
spec:
  serverId:
    valueFrom:
      name: my-sql
  databaseName: tenant-42
  skuName: ElasticPool
  elasticPoolId:
    valueFrom:
      kind: AzureMssqlElasticPool
      name: tenant-pool
      fieldPath: status.outputs.elastic_pool_id
```

### Geo Secondary for Disaster Recovery

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlDatabase
metadata:
  name: orders-dr
spec:
  serverId:
    valueFrom:
      name: my-sql-westus
  databaseName: orders
  skuName: GP_Gen5_2
  createMode: SECONDARY
  secondaryType: GEO
  creationSourceDatabaseId:
    valueFrom:
      kind: AzureMssqlDatabase
      name: orders-db
      fieldPath: status.outputs.database_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `database_id` | `string` | The database's ARM ID -- referenced by copy/secondary/restore databases' `creationSourceDatabaseId` |
| `database_name` | `string` | The `Database=` segment of connection strings against the server's fqdn |

## Related Components

- [AzureMssqlServer](/docs/catalog/azure/sql-server) — the parent logical server
- [AzureMssqlElasticPool](/docs/catalog/azure/sql-elastic-pool) — pooled compute the database can join
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — the database-scoped CMK unwrap identity
- [AzureKeyVaultKey](/docs/catalog/azure/key-vault-key) — the database-scoped TDE key
