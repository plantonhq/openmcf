---
title: "Cloud SQL Database"
description: "Cloud SQL Database deployment documentation"
icon: "package"
order: 100
componentName: "gcpcloudsqldatabase"
---

# GCP Cloud SQL Database

Creates a logical database inside a Cloud SQL instance. Databases are first-class composable nodes: one per application on a shared instance, each with its own lifecycle — create and drop them freely without touching the instance.

## What Gets Created

When you deploy a GcpCloudSqlDatabase resource, Planton provisions:

- **Cloud SQL database** — a `google_sql_database` inside the referenced instance

## Prerequisites

- **An existing Cloud SQL instance** — a [GcpCloudSql](/docs/catalog/gcp/cloud-sql) resource (or a literal instance name)
- **GCP credentials** with Cloud SQL admin permissions on the project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSqlDatabase
metadata:
  name: orders-database
spec:
  instance:
    valueFrom:
      kind: GcpCloudSql
      name: my-postgres
      fieldPath: status.outputs.instance_name
  databaseName: orders
```

```shell
planton apply -f database.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `instance` | `StringValueOrRef` | — (required) | The hosting instance (ref → GcpCloudSql). Immutable. |
| `databaseName` | `string` | — (required) | Database name inside the instance. Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the instance. |
| `charset` | `string` | engine default | MySQL: `utf8mb4` recommended. PostgreSQL: `UTF8` only. |
| `collation` | `string` | engine default | Engine-specific collation name. |

## Examples

### MySQL Database with Modern UTF-8

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSqlDatabase
metadata:
  name: app-database
spec:
  instance:
    valueFrom:
      kind: GcpCloudSql
      name: my-mysql
      fieldPath: status.outputs.instance_name
  databaseName: app
  charset: utf8mb4
  collation: utf8mb4_0900_ai_ci
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `database_name` | Name of the database inside the instance |
| `self_link` | GCP resource self link |

## Related Components

- [GcpCloudSql](/docs/catalog/gcp/cloud-sql) — the instance this database lives on
- [GcpCloudSqlUser](/docs/catalog/gcp/cloud-sql-user) — pair each database with its application user
