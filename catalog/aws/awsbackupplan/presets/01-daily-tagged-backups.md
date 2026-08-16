# Daily Tagged Backups

This preset creates the workhorse plan: a daily 05:00 UTC backup of
everything tagged `backup: true`, kept 35 days, into a vault wired by
reference.

## When to Use

- The first backup plan in an account — opt resources in by tagging
  them, no plan edits needed
- Environments where coverage should follow tags, not inventories

## What You Get

- A daily rule with AWS's default backup windows and a 35-day expiry
  (recovery-point storage is what AWS bills)
- A tag-driven selection assuming the referenced IAM role (it must
  trust `backup.amazonaws.com`)

## Customize

- Add `enableContinuousBackup: true` for point-in-time restore (keep
  `deleteAfterDays` within 35 — the continuous cap)
- Add `lifecycle.coldStorageAfterDays` for long retention (expiry must
  exceed it by at least 90 days — AWS's cold-storage minimum)
- Add `copyActions` targeting another vault/region for DR, or
  `targetLogicallyAirGappedBackupVaultArn` for the ransomware tier
- Cover explicit resources with `resources` ARNs instead of (or on top
  of) tags
