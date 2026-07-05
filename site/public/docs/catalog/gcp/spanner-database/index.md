---
title: "Spanner Database"
description: "Spanner Database deployment documentation"
icon: "package"
order: 100
componentName: "gcpspannerdatabase"
---

# GCP Spanner Database

Deploys a Cloud Spanner database within an existing Spanner instance, with support for GoogleSQL and PostgreSQL dialects, initial DDL schema creation, single- and multi-region CMEK encryption, configurable point-in-time recovery, and two-level deletion guards.

## What Gets Created

When you deploy a GcpSpannerDatabase resource, Planton provisions:

- **Spanner Database** — a `google_spanner_database` in the specified instance with the chosen SQL dialect and version retention period; the Spanner API is enabled automatically
- **Initial Schema** — when `ddl` is provided, statements execute atomically with database creation
- **CMEK Encryption** — when `encryptionConfig` is provided, the database is encrypted with customer-managed KMS keys

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing Spanner instance** (deploy via GcpSpannerInstance first)
- **A KMS key** in the same location as the instance configuration when enabling CMEK (one key per region for multi-region instances); the Spanner service agent needs `roles/cloudkms.cryptoKeyEncrypterDecrypter` on it

## Quick Start

Create a file `spanner-database.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSpannerDatabase
metadata:
  name: orders-db
spec:
  instance:
    valueFrom:
      kind: GcpSpannerInstance
      name: prod-spanner
      fieldPath: status.outputs.instance_name
```

Deploy:

```shell
planton apply -f spanner-database.yaml
```

This creates an empty database named `orders-db` with the default GoogleSQL dialect and 1-hour version retention, in the provider's default project.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `instance` | `StringValueOrRef` | Spanner instance to create the database on. References a GcpSpannerInstance's `instance_name` output via `valueFrom`. Immutable. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project owning the instance. Can reference a GcpProject. |
| `databaseName` | `string` | `metadata.name` | Database name within the instance. Immutable. 2-30 chars, pattern `^[a-z][a-z0-9_\-]*[a-z0-9]$`. |
| `databaseDialect` | `string` | `GOOGLE_STANDARD_SQL` | `GOOGLE_STANDARD_SQL` or `POSTGRESQL`. Immutable and permanent. |
| `versionRetentionPeriod` | `string` | `1h` | Point-in-time recovery window, 1h-7d (e.g. `24h`, `3d`). Mutable. |
| `ddl` | `list<string>` | — | Initial DDL executed atomically with creation. Append-only afterwards; editing an existing entry recreates the database. |
| `enableDropProtection` | `bool` | `false` | GCP API-side lock: blocks deletion through ANY interface and blocks deletion of the parent instance. Mutable. |
| `encryptionConfig.kmsKeyName` | `StringValueOrRef` | — | Regional CMEK key (references a GcpKmsKey `key_id`). Immutable. Exactly one of the two key shapes. |
| `encryptionConfig.kmsKeyNames[]` | `list<StringValueOrRef>` | — | Multi-region CMEK: one key per region of the instance configuration. Immutable. |
| `defaultTimeZone` | `string` | `America/Los_Angeles` | IANA time zone affecting time-zone-dependent SQL functions. |
| `deletionProtection` | `bool` | `true` | IaC-side guard: both engines refuse to destroy the database while true. Set `false` explicitly before intentional teardown. |

## Examples

### PostgreSQL-Dialect Database

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSpannerDatabase
metadata:
  name: analytics-db
spec:
  instance:
    valueFrom:
      kind: GcpSpannerInstance
      name: prod-spanner
      fieldPath: status.outputs.instance_name
  databaseDialect: POSTGRESQL
  versionRetentionPeriod: "7d"
```

### Database with Initial Schema

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSpannerDatabase
metadata:
  name: orders-db
spec:
  instance:
    valueFrom:
      kind: GcpSpannerInstance
      name: prod-spanner
      fieldPath: status.outputs.instance_name
  ddl:
    - CREATE TABLE orders (order_id STRING(36) NOT NULL, created_at TIMESTAMP NOT NULL) PRIMARY KEY (order_id)
    - CREATE INDEX orders_by_created ON orders(created_at)
```

### CMEK-Encrypted with Drop Protection

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSpannerDatabase
metadata:
  name: payments-db
spec:
  instance:
    valueFrom:
      kind: GcpSpannerInstance
      name: prod-spanner
      fieldPath: status.outputs.instance_name
  encryptionConfig:
    kmsKeyName:
      valueFrom:
        kind: GcpKmsKey
        name: spanner-cmek
        fieldPath: status.outputs.key_id
  enableDropProtection: true
  versionRetentionPeriod: "3d"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `database_id` | `string` | Fully qualified database ID (`projects/{project}/instances/{instance}/databases/{name}`) |
| `database_name` | `string` | Short database name, referenced by GcpSpannerBackupSchedule |
| `state` | `string` | Database state: `CREATING` or `READY` |

## Related Components

- [GcpSpannerInstance](/docs/catalog/gcp/spanner-instance) — the instance this database lives on
- [GcpSpannerBackupSchedule](/docs/catalog/gcp/spanner-backup-schedule) — cron-driven full/incremental backups for this database
- [GcpKmsKey](/docs/catalog/gcp/kms-key) — customer-managed encryption keys for CMEK
- [GcpProject](/docs/catalog/gcp/project) — the project the instance lives in
