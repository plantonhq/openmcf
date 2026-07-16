# GcpSpannerBackupSchedule — Research and Design Documentation

## 1. Why Schedules Are a First-Class Kind

```
Project
└── Spanner Instance (compute + topology)         GcpSpannerInstance
    └── Database (schema + data + encryption)     GcpSpannerDatabase
        ├── Backup Schedule: daily incremental    GcpSpannerBackupSchedule
        └── Backup Schedule: weekly full archive  GcpSpannerBackupSchedule
```

Backup schedules pass the split test on all three axes: they are the provider's own resource (`google_spanner_backup_schedule`), they are many-per-database (the daily-plus-weekly pattern above is the production norm), and they carry an independent lifecycle (cadence, retention, and encryption change without touching the database). Bundling them into the database spec would force one schedule per database and couple schedule edits to the database resource.

Point-in-time recovery (`version_retention_period` on the database, at most 7 days) and backups are complementary, not alternatives: PITR answers "read the data as of an hour ago"; backups answer "restore the database from last month."

## 2. Terraform Provider Floor

Designed from `google_spanner_backup_schedule` on the released Terraform Google provider 6.x line (`~> 6.0`); the Spanner surface is fully GA (GA and beta schemas identical). Both engines enable `spanner.googleapis.com` before creating the schedule, and both build `schedule_id` from the created resource's resolved attributes so the ambient-project fallback stays honest.

### Field coverage

| Provider surface | Modeled | Notes |
|---|---|---|
| `instance`, `database`, `name`, `project` | ✅ | instance/database by FK to the parents' name outputs; name defaults to `metadata.name`; project ambient fallback |
| `spec.cron_spec.text` | ✅ | surfaced as the flat `cron` field — the wrapper blocks carry no other content |
| `retention_duration` | ✅ | seconds-duration pattern validated pre-deploy; max 366 days |
| `full_backup_spec` / `incremental_backup_spec` | ✅ | folded into one `backup_type` field (default FULL) — the provider expresses one immutable choice as a pair of empty marker blocks; a string field is the honest shape of that choice |
| `encryption_config` | ✅ | all three types; CMEK key XOR keys enforced pre-deploy; both FK → GcpKmsKey `key_id` |

### Recorded skips (evidence-based)

| Feature | Reason |
|---|---|
| `deletion_policy` | Client-side Terraform lever (PREVENT/ABANDON) that conflicts with managed destroy semantics; catalog-wide exclusion. |

## 3. Behavioral Notes

- **Ownership and survival:** the schedule is owned by the database (deleting the database deletes its schedules), but the backups it created live until their retention expires. Surviving backups are what make destroying the parent instance fail unless `GcpSpannerInstance.force_destroy` is set.
- **Cadence bounds:** the API accepts only 12-hour, daily, weekly, and monthly cron frequencies, evaluated in UTC; the crontab syntax is validated server-side (richer expressions are rejected at create).
- **Retention semantics:** changing `retention_duration` affects only backups created after the change.
- **Edition gate:** INCREMENTAL schedules require the instance to be ENTERPRISE or ENTERPRISE_PLUS; GCP rejects them on STANDARD instances at create time.
- **Relationship to AUTOMATIC default schedules:** `GcpSpannerInstance.default_backup_schedule_type = AUTOMATIC` makes GCP attach a default schedule to each new database. Explicit GcpSpannerBackupSchedule resources are the full-control alternative; both can coexist.

## 4. Immutability

ForceNew (recreate on change): `instance`, `database`, `schedule_name`, `backup_type`, `project_id`. Mutable in place: `cron`, `retention_duration`, `encryption_config`.

## 5. Downstream Composition

The schedule is a leaf — nothing references it by FK. Its `schedule_id` output is the handle for API callers and audit tooling. In the spanner application pattern:

```
GcpSpannerInstance (instance_name)
  └── GcpSpannerDatabase (database_name)
        ├── GcpSpannerBackupSchedule (daily operational)
        └── GcpSpannerBackupSchedule (weekly archive)
```
