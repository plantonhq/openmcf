# Daily Short-Retention Backups

A daily backup schedule with a one-week retention window — the common
companion to PITR for recent restore points without long-term storage cost.

## When to use

Production databases where you want managed backups beyond PITR's 7-day
version window, but do not need months of archive retention.

## What to customize

- `retention` — the per-backup lifetime (`604800s` = 7 days; up to
  `8467200s` = 14 weeks).
- `database` — reference your `GcpFirestoreDatabase` resource.

## Composes with

`GcpFirestoreDatabase` upstream (reference its `database_name` output).
Pair with `02-weekly-long-retention` on the same database for the
daily-plus-weekly pattern Firestore supports.
