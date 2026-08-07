---
title: "Cloud SQL Database"
description: "Cloud SQL Database deployment documentation"
icon: "package"
order: 100
componentName: "gcpcloudsqldatabase"
---

# GCP Cloud SQL Database

Creates a database inside an existing Google Cloud SQL instance. Databases are composable satellites of the instance — one instance hosts many databases, each owned by its own application, each created, reviewed, and deleted as a first-class Cloud Resource. Pairs naturally with GCP Cloud SQL User for per-application credentials.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cloud SQL Database** -- a `google_sql_database` on the referenced instance, with the specified name and (optionally) an engine-specific character set and collation

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### GCP Project

- **A Cloud SQL instance** -- the [GcpCloudSql](/cloud-catalog/gcp-cloud-sql) instance that hosts the database. Reference it via ValueFromRef so the pipeline deploys the instance first.

## Deploy

### Console

Open the deployment store, find **GCP Cloud SQL Database**, and click **Create**. The wizard walks two decisions: which instance hosts the database, then the database's name and collation. The [Presets](#presets) tab offers **PostgreSQL App Database** and **MySQL utf8mb4 Database** starting points.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSqlDatabase
metadata:
  name: orders
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  instance:
    valueFrom:
      kind: GcpCloudSql
      name: orders-db-prod
      fieldPath: status.outputs.instance_name
  databaseName: orders
  charset: UTF8
  collation: en_US.UTF8
```

```shell
planton apply -f database.yaml
```

## Key Configuration

**Instance** -- Immutable. A database never moves between instances; relocating data is an export/import. Reference the GcpCloudSql resource rather than typing the name.

**Database name** -- Immutable, max 128 characters. It is what applications put in their connection strings — name by the owning application.

**Charset and collation** -- Engine-interpreted: MySQL accepts any supported pair (`utf8mb4` + `utf8mb4_0900_ai_ci` is the modern default); PostgreSQL wants `UTF8` with an OS locale collation (`en_US.UTF8`); SQL Server ignores charset and uses its own collation names. Empty keeps the engine default.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpCloudSql** | `instance` | `status.outputs.instance_name` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `database_name` | The created database's name | Application connection strings, service configuration |
| `self_link` | GCP resource self link | Audit log filters |

## Works With

- [**GCP Cloud SQL**](/cloud-catalog/gcp-cloud-sql) -- the instance that hosts this database
- [**GCP Cloud SQL User**](/cloud-catalog/gcp-cloud-sql-user) -- per-application credentials for connecting to this database
