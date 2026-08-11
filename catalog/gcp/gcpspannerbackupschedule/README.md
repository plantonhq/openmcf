# GcpSpannerBackupSchedule

Provisions a [Google Cloud Spanner](https://cloud.google.com/spanner) backup schedule on an existing Spanner database.

## What It Does

A backup schedule creates backups of one database on a cron cadence and retains each backup for a configured duration. This component manages the schedule itself — its cadence, retention, backup kind (full or incremental), and the encryption posture of the backups it creates. The Spanner API is enabled automatically in the target project.

Schedules are first-class, many-per-database resources: a production database commonly carries a daily incremental schedule AND a weekly full schedule side by side. The schedule is owned by the database — deleting the database deletes its schedules — but the BACKUPS themselves survive until their retention expires, which is what makes destroying the instance fail unless `GcpSpannerInstance.force_destroy` is set.

## When to Use

- Any production Spanner database (point-in-time recovery via `version_retention_period` covers at most 7 days; backups cover months)
- Compliance archives with retention measured in months (up to 366 days)
- Cost-efficient high-frequency restore points via incremental chains on ENTERPRISE instances

## Key Configuration

### Cadence (`cron`)

A crontab expression evaluated in **UTC**. Spanner accepts a bounded set of frequencies — every 12 hours, daily, weekly, or monthly:

| Expression | Meaning |
|---|---|
| `0 2/12 * * *` | every 12 hours, at 02:00 and 14:00 UTC |
| `0 2 * * *` | daily at 02:00 UTC |
| `0 2 * * 0` | weekly on Sunday at 02:00 UTC |
| `0 2 8 * *` | monthly on the 8th at 02:00 UTC |

Mutable — cadence changes apply in place.

### Retention (`retention_duration`)

How long each backup lives, as a seconds duration string (`86400s` = 1 day, `31622400s` = 366 days, the maximum). Mutable — applies to backups created AFTER the change; existing backups keep the retention they were created with.

### Backup Kind (`backup_type`)

Immutable; chosen at creation:

| Kind | When to Use |
|---|---|
| `FULL` (default) | Every backup is a complete, self-contained copy. Works on every edition. |
| `INCREMENTAL` | Backups form chains storing only changes since the previous one — significantly cheaper storage, same restore semantics. Requires the instance to be ENTERPRISE or ENTERPRISE_PLUS edition. |

### Backup Encryption (`encryption_config`)

Omit for `USE_DATABASE_ENCRYPTION` (backups inherit the database's posture — CMEK databases get CMEK backups). Or set explicitly:

- `GOOGLE_DEFAULT_ENCRYPTION` — Google-managed keys regardless of the database's posture
- `CUSTOMER_MANAGED_ENCRYPTION` — explicit CMEK: exactly one of `kms_key_name` (regional instance configurations) or `kms_key_names` (one key per region of a multi-region configuration); both reference a `GcpKmsKey`'s `key_id` output

### Destroy Behavior (`deletion_policy`)

`DELETE` (default) removes the schedule — no further backups are taken; `PREVENT` fails the destroy, protecting a cadence a recovery objective depends on; `ABANDON` drops it from management but keeps it running in GCP. None of these touches backups already taken — those live until their retention expires.

## Outputs

| Output | Description |
|---|---|
| `schedule_id` | Fully qualified path (`projects/{project}/instances/{instance}/databases/{database}/backupSchedules/{name}`) — the API handle |
| `schedule_name` | Short schedule name within the database |

## Relationships

- **Depends on**: GcpSpannerInstance (`instance`), GcpSpannerDatabase (`database`), GcpProject (`project_id`, ambient fallback), optionally GcpKmsKey (backup CMEK)
- **Interacts with**: `GcpSpannerInstance.default_backup_schedule_type` (AUTOMATIC attaches GCP-managed default schedules to new databases; explicit schedules like this one give full control instead) and `GcpSpannerInstance.force_destroy` (surviving backups block instance destruction while false)

## Deployment

```shell
planton apply -f backup-schedule.yaml
```

For copy-paste ready manifests, see the [presets](presets/).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
