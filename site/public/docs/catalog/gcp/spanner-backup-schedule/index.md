---
title: "Spanner Backup Schedule"
description: "Spanner Backup Schedule deployment documentation"
icon: "package"
order: 100
componentName: "gcpspannerbackupschedule"
---

# GCP Spanner Backup Schedule

Deploys a cron-driven backup schedule on a Cloud Spanner database, with full or incremental backup chains, retention up to 366 days, and optional customer-managed encryption for the backup copies.

## What Gets Created

When you deploy a GcpSpannerBackupSchedule resource, Planton provisions:

- **Backup Schedule** — a `google_spanner_backup_schedule` on the specified database with the chosen cadence and retention; the Spanner API is enabled automatically
- **Backup Kind** — a full-backup or incremental-backup specification (immutable per schedule)
- **Backup Encryption** — when `encryptionConfig` is provided, an explicit encryption posture for the backups (otherwise they inherit the database's)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing Spanner instance and database** (deploy via GcpSpannerInstance and GcpSpannerDatabase first)
- **ENTERPRISE or ENTERPRISE_PLUS edition** on the instance when using `backupType: INCREMENTAL`
- **A KMS key** with the Spanner service agent granted `roles/cloudkms.cryptoKeyEncrypterDecrypter` when using CUSTOMER_MANAGED_ENCRYPTION

## Quick Start

Create a file `backup-schedule.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSpannerBackupSchedule
metadata:
  name: orders-daily-backups
spec:
  instance:
    valueFrom:
      kind: GcpSpannerInstance
      name: prod-spanner
      fieldPath: status.outputs.instance_name
  database:
    valueFrom:
      kind: GcpSpannerDatabase
      name: orders-db
      fieldPath: status.outputs.database_name
  cron: "0 2 * * *"
  retentionDuration: 2678400s
```

Deploy:

```shell
planton apply -f backup-schedule.yaml
```

This backs the database up daily at 02:00 UTC, keeping each full backup for 31 days, in the provider's default project.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `instance` | `StringValueOrRef` | Spanner instance hosting the database. References a GcpSpannerInstance's `instance_name` output. Immutable. | Required |
| `database` | `StringValueOrRef` | Database to back up. References a GcpSpannerDatabase's `database_name` output. Immutable. | Required |
| `cron` | `string` | Crontab cadence, evaluated in UTC. Allowed frequencies: 12-hour, daily, weekly, monthly. Mutable. | Required |
| `retentionDuration` | `string` | Per-backup retention as a seconds duration (`86400s`); up to 366 days (`31622400s`). Mutable. | Required, pattern `^[0-9]+(\.[0-9]{1,9})?s$` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project owning the instance. Can reference a GcpProject. |
| `scheduleName` | `string` | `metadata.name` | Schedule name within the database. Immutable. Pattern `^[a-z][-a-z0-9]*[a-z0-9]$`, up to 63 chars. |
| `backupType` | `string` | `FULL` | `FULL` or `INCREMENTAL` (incremental requires ENTERPRISE+ edition). Immutable. |
| `encryptionConfig.encryptionType` | `string` | `USE_DATABASE_ENCRYPTION` | `USE_DATABASE_ENCRYPTION`, `GOOGLE_DEFAULT_ENCRYPTION`, or `CUSTOMER_MANAGED_ENCRYPTION`. Mutable. |
| `encryptionConfig.kmsKeyName` | `StringValueOrRef` | — | Regional backup CMEK key (references a GcpKmsKey `key_id`). CMEK only; exactly one key shape. |
| `encryptionConfig.kmsKeyNames[]` | `list<StringValueOrRef>` | — | Multi-region backup CMEK: one key per region of the instance configuration. CMEK only. |

## Examples

### Incremental Backups on an Enterprise Instance

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSpannerBackupSchedule
metadata:
  name: orders-incremental-backups
spec:
  instance:
    valueFrom:
      kind: GcpSpannerInstance
      name: prod-spanner
      fieldPath: status.outputs.instance_name
  database:
    valueFrom:
      kind: GcpSpannerDatabase
      name: orders-db
      fieldPath: status.outputs.database_name
  cron: "0 2/12 * * *"
  retentionDuration: 1209600s
  backupType: INCREMENTAL
```

### CMEK-Encrypted Weekly Archive

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSpannerBackupSchedule
metadata:
  name: orders-weekly-archive
spec:
  instance:
    valueFrom:
      kind: GcpSpannerInstance
      name: prod-spanner
      fieldPath: status.outputs.instance_name
  database:
    valueFrom:
      kind: GcpSpannerDatabase
      name: orders-db
      fieldPath: status.outputs.database_name
  cron: "0 4 * * 0"
  retentionDuration: 31622400s
  encryptionConfig:
    encryptionType: CUSTOMER_MANAGED_ENCRYPTION
    kmsKeyName:
      valueFrom:
        kind: GcpKmsKey
        name: spanner-backup-cmek
        fieldPath: status.outputs.key_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `schedule_id` | `string` | Fully qualified schedule ID (`projects/{project}/instances/{instance}/databases/{database}/backupSchedules/{name}`) |
| `schedule_name` | `string` | Short schedule name within the database |

## Related Components

- [GcpSpannerInstance](/docs/catalog/gcp/spanner-instance) — the instance whose `default_backup_schedule_type` and `force_destroy` interact with schedules and their backups
- [GcpSpannerDatabase](/docs/catalog/gcp/spanner-database) — the database this schedule protects
- [GcpKmsKey](/docs/catalog/gcp/kms-key) — customer-managed encryption keys for backup CMEK
- [GcpProject](/docs/catalog/gcp/project) — the project the instance lives in
