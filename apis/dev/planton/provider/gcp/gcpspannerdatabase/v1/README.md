# GcpSpannerDatabase

Provisions a [Google Cloud Spanner](https://cloud.google.com/spanner) database within an existing Spanner instance.

## What It Does

A Spanner database is a collection of tables, views, indexes, and other schema objects that live inside a Spanner instance. This component creates and manages the database itself — its SQL dialect, initial schema (via DDL), encryption posture, point-in-time recovery window, and deletion guards. The Spanner API is enabled automatically in the target project.

A Spanner database is **not** an instance. The instance ([GcpSpannerInstance](../gcpspannerinstance/v1)) provides compute capacity and geographic configuration; the database stores your data and schema. Multiple databases share a single instance, and backup schedules ([GcpSpannerBackupSchedule](../gcpspannerbackupschedule/v1)) attach to a database by reference.

## When to Use

- You have a GcpSpannerInstance and need a database on it
- You want the database lifecycle (creation, initial schema, encryption) managed through Planton
- You need CMEK encryption or deletion locks for compliance

## Key Configuration

### SQL Dialect (`database_dialect`)

Chosen at creation time; permanent:

| Dialect | When to Use |
|---|---|
| `GOOGLE_STANDARD_SQL` | Full Spanner feature set including interleaved tables and STRUCT types. Default. |
| `POSTGRESQL` | Teams standardized on PostgreSQL syntax and tooling. Some Spanner-specific features are unavailable. |

### Initial DDL (`ddl`)

DDL statements create tables, indexes, and views atomically with the database — if any statement fails, the database is not created.

**Append-only lifecycle:** new statements added later are applied via UpdateDDL; modifying or removing an existing statement forces database recreation. For ongoing schema management, use a migration tool (Liquibase, Flyway).

### Encryption (`encryption_config`)

Omit for Google-managed encryption. For CMEK, set exactly one shape — both immutable:

- `kms_key_name` — one key, for **regional** instance configurations (key must live in the config's location); references a `GcpKmsKey`'s `key_id` output
- `kms_key_names` — one key **per region** of a **multi-region** instance configuration

### Point-in-Time Recovery (`version_retention_period`)

How far back data can be read at a previous timestamp. Range: 1 hour to 7 days (default 1h). Longer windows consume more storage. Mutable.

### Two Deletion Guards

| Guard | Enforced by | Effect |
|---|---|---|
| `deletion_protection` (default **true**) | IaC engines | A destroy plan fails before touching GCP. Set false explicitly before intentional teardown. |
| `enable_drop_protection` (default false) | GCP API | NO interface (console, gcloud, API, IaC) can delete the database — and the parent instance cannot be deleted either. Compliance-grade lock. |

### Labels

Spanner databases do **not** support GCP labels. Labels are managed at the instance level via GcpSpannerInstance.

## Outputs

| Output | Description |
|---|---|
| `database_id` | Fully qualified path (`projects/{project}/instances/{instance}/databases/{name}`) — the IAM/API handle |
| `database_name` | Short name (referenced by GcpSpannerBackupSchedule) |
| `state` | CREATING or READY |

## Relationships

- **Depends on**: GcpSpannerInstance (`instance`), GcpProject (`project_id`, ambient fallback), optionally GcpKmsKey (CMEK)
- **Referenced by**: GcpSpannerBackupSchedule (`database`), application connection strings, IAM bindings

## Deployment

```shell
planton apply -f spanner-database.yaml
```

For copy-paste ready manifests, see the [presets](presets/).
