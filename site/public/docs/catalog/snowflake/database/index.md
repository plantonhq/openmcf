---
title: "Database"
description: "Database deployment documentation"
icon: "package"
order: 100
componentName: "snowflakedatabase"
---

# Snowflake Database

Deploys a Snowflake database with configurable Time Travel retention, transient mode, Iceberg table defaults, and task execution parameters. All database-level settings are declared in the spec and applied idempotently. Integrates with Planton's Snowflake Provider Connection for credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Snowflake Database** -- a database with the specified name and all configured parameters including Time Travel retention, default collation, Iceberg catalog settings, logging levels, and user task execution behavior

## Before You Deploy

### Planton Setup

- **Snowflake Provider Connection** -- an active connection in the Connect module with Snowflake account credentials (account name, username, password). Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credential authentication.

### Snowflake Account

- **A Snowflake account** with permissions to create databases. The user configured in the Provider Connection must have the `CREATE DATABASE` privilege.
- **An external volume** if configuring Iceberg table defaults -- create the external volume in Snowflake first, then reference it via the `externalVolume` field.
- **Database naming** -- the `name` field is the Snowflake identifier and must be unique within the account. Snowflake convention uses uppercase identifiers. Avoid characters: `|`, `.`, `(`, `)`, `"`.
- **A warehouse** (optional) -- required if configuring user task execution with managed warehouses via the `userTask` sub-object.

## Deploy

### Console

Open the deployment store, find **Snowflake Database**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production** preset in the [Presets](#presets) tab to pre-populate a standard production configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: snowflake.planton.dev/v1
kind: SnowflakeDatabase
metadata:
  name: analytics-db
  org: acme-corp
  env: prod
spec:
  name: ANALYTICS
  comment: Production analytics database
  dataRetentionTimeInDays: 30
  dropPublicSchemaOnCreation: true
```

```shell
planton apply -f snowflake-database.yaml
```

This creates a persistent Snowflake database named `ANALYTICS` with 30-day Time Travel retention and the public schema dropped on creation. No Iceberg, logging, or task settings are configured. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Snowflake database. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Time Travel retention** -- The `dataRetentionTimeInDays` field controls how long CLONE and UNDROP operations remain available. Set to 30 for production databases that need a wide recovery window. Use 1 for development databases to minimize storage costs. Set `maxDataExtensionTimeInDays` to prevent streams from going stale on tables with infrequent changes.

**Transient vs. persistent** -- Set `isTransient` to `true` for databases that do not need Fail-safe protection. Transient databases skip the 7-day Fail-safe period after Time Travel expires, reducing storage costs. Use for development, staging, and ephemeral workloads.

**Iceberg table defaults** -- Set `catalog` and `externalVolume` to configure the database for Apache Iceberg tables. Choose `storageSerializationPolicy: OPTIMIZED` for best Snowflake performance, or `COMPATIBLE` when third-party compute engines also read the data files.

**Logging and tracing** -- The `logLevel` field controls event table ingestion severity (TRACE through OFF). Set to WARN or ERROR in production to keep event tables manageable. Enable `enableConsoleOutput` and set `traceLevel` to ON_EVENT during development for stored procedure debugging.

**User task execution** -- The `userTask` sub-object controls managed warehouse sizing, minimum trigger intervals, and execution timeouts. Set `suspendTaskAfterNumFailures` to automatically suspend repeatedly failing tasks, and `taskAutoRetryAttempts` to allow automatic retries.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Provider-assigned unique ID for the database resource | API operations, resource tracking |
| `name` | Fully-qualified name of the created database | Schema creation, cross-database queries |
| `owner` | Owner role of the database | Role-based access control verification |
| `created_on` | Timestamp when the database was created | Audit logs, lifecycle tracking |
| `is_transient` | Whether the database is transient | Environment classification, cost reporting |
| `data_retention_time_in_days` | Configured data retention time in days | Compliance verification, backup policy validation |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production data warehouse** -- Persistent database with 30-day Time Travel retention, warning-level logging, and the public schema dropped to enforce explicit schema design. Start from the **Production** preset.

**Development database** -- Transient database with 1-day retention, debug logging, and console output enabled for stored procedure troubleshooting. Minimal storage costs with no Fail-safe period. Start from the **Development** preset.

**Iceberg analytics** -- Database configured with Snowflake as the Iceberg catalog, an external volume for data storage, and optimized serialization policy. Enables open table format interoperability with Spark, Trino, and other compute engines. Start from the **Iceberg Analytics** preset.

## Works With

This component operates independently and does not reference other components.