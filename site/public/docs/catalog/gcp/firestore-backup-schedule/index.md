---
title: "Firestore Backup Schedule"
description: "Firestore Backup Schedule deployment documentation"
icon: "package"
order: 100
componentName: "gcpfirestorebackupschedule"
---

# GCP Firestore Backup Schedule

Deploys a managed backup schedule on a Cloud Firestore database — periodic backups with retention up to 14 weeks, distinct from point-in-time recovery.

## What Gets Created

A Firestore backup schedule with daily or weekly recurrence and configurable retention. The Firestore API is enabled automatically.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing Firestore database** — referenced via `database` (a `GcpFirestoreDatabase` resource or a literal name)
- **IAM permissions** — Firestore Admin access (e.g. `roles/datastore.owner`)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpFirestoreBackupSchedule
metadata:
  name: orders-daily-backups
spec:
  database:
    valueFrom:
      kind: GcpFirestoreDatabase
      name: prod-firestore
      fieldPath: status.outputs.database_name
  retention: 604800s
  daily: true
```

```shell
planton apply -f backup-schedule.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `database` | `StringValueOrRef` | — (required) | Database to back up. Immutable. |
| `retention` | `string` | — (required) | Per-backup retention (`604800s` = 7 days). Mutable. |
| `daily` | `bool` | — | Daily recurrence (exactly one of daily or weekly). Immutable. |
| `weeklyRecurrence.day` | `string` | — | Weekly day (`MONDAY`–`SUNDAY`). Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project owning the database. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `schedule_id` | Server-assigned schedule ID (last path segment) |
| `database` | The database the schedule protects |

## Related Components

- [GcpFirestoreDatabase](/docs/catalog/gcp/firestore-database) — the database this schedule protects
- [GcpProject](/docs/catalog/gcp/project) — provides the GCP project
