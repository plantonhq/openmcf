# GcpFirestoreDatabase

Provisions a [Google Cloud Firestore](https://cloud.google.com/firestore) database with configurable database type, edition, PITR, CMEK encryption, delete protection, ENTERPRISE data-access modes (classic Firestore API, MongoDB-compatible API, realtime updates), resource-manager tags, and teardown policy.

## What It Does

A Firestore database is the top-level container for collections, documents, and indexes in Google Cloud Firestore. This component creates and manages the database itself, including its type (Native or Datastore mode), edition, encryption configuration, and point-in-time recovery settings.

Each GCP project can have multiple named databases plus one special `(default)` database. The default database is what client libraries connect to when no database ID is specified.

## When to Use

- You need a managed NoSQL document database for your application
- You want to manage the database lifecycle (creation, encryption, PITR) through Planton
- You need to enforce CMEK encryption or delete protection for compliance
- You want a named database separate from the project's default database

## Key Configuration

### Database Type (`type`)

Choose the database type at creation time:

| Type | When to Use |
|---|---|
| `FIRESTORE_NATIVE` | Modern Firestore with real-time listeners, offline support, and mobile SDKs. Recommended for new applications. |
| `DATASTORE_MODE` | Legacy Datastore-compatible mode. Use for existing Datastore applications or server-side workloads. |

### Database Name (`database_name`)

- Use `(default)` for the project's primary database (one per project)
- Use a custom name (4-63 chars, lowercase letters/digits/hyphens) for additional databases

### Database Edition (`database_edition`)

| Edition | When to Use |
|---|---|
| `STANDARD` | Default. Suitable for most workloads. |
| `ENTERPRISE` | Enhanced SLA, advanced security. Requires `FIRESTORE_NATIVE` type. |

### Point-in-Time Recovery (`point_in_time_recovery_enablement`)

When enabled, retains 7 days of version history for disaster recovery. When disabled (default), retains 1 hour.

### Data-Access Modes (ENTERPRISE edition only)

Three switches control which protocols can read and write an ENTERPRISE database — the shape for a database dedicated to MongoDB drivers, or one that must stay single-protocol:

| Field | Values | Controls |
|---|---|---|
| `firestore_data_access_mode` | `DATA_ACCESS_MODE_ENABLED` / `DATA_ACCESS_MODE_DISABLED` | The classic Firestore API |
| `mongodb_compatible_data_access_mode` | `DATA_ACCESS_MODE_ENABLED` / `DATA_ACCESS_MODE_DISABLED` | MongoDB drivers and tools (pair with `MONGODB_COMPATIBLE_API`-scoped GcpFirestoreIndex indexes) |
| `realtime_updates_mode` | `REALTIME_UPDATES_MODE_ENABLED` / `REALTIME_UPDATES_MODE_DISABLED` | Live query snapshots |

### Teardown Policy (`deletion_policy`)

Defaults to `DELETE` so destroying the resource destroys the database (the raw provider's own default, ABANDON, would leave it running unmanaged). `PREVENT` makes a destroy fail — belt-and-suspenders beside `delete_protection_state`; `ABANDON` unmanages the database while it keeps serving.

### Labels and Tags

Firestore databases do **not** support GCP labels. Resource Manager **tags** (`resource_manager_tags`, `tagKeys/{id}` → `tagValues/{id}`) are supported at create time only — mutating them replaces the database.

## Outputs

| Output | Description |
|---|---|
| `database_id` | Fully qualified path (`projects/{project}/databases/{name}`) |
| `database_name` | Database name (e.g., `(default)` or custom name) |
| `uid` | Server-generated UUID4 |
| `create_time` | Creation timestamp (RFC3339) |
| `earliest_version_time` | Earliest PITR recovery timestamp (RFC3339) |
| `version_retention_period` | Version history window (`3600s` without PITR, `604800s` with) |
| `key_prefix` | Datastore Mode App Engine key prefix (empty otherwise) |
| `update_time` | Last configuration-change timestamp (RFC3339) |

## Relationships

- **Depends on**: GcpProject (project_id), optionally GcpKmsKey (kms_key_name)
- **Referenced by**: GcpFirestoreIndex (composite indexes) and GcpFirestoreBackupSchedule (managed backups) — both many-per-database, composing against the `database_name` output — plus application connection strings and Firebase client libraries

## Deployment

```shell
planton apply -f firestore-database.yaml
```

For copy-paste ready manifests, see e2e/manifest.yaml.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
