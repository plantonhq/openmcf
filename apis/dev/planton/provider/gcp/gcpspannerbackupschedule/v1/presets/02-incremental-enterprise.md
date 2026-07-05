# Incremental Backups (Enterprise)

Backs the database up every 12 hours using incremental chains — each backup stores only the changes since the previous one, at a fraction of full-backup storage cost with the same restore semantics.

## When to Use

- Twice-daily restore points without paying for two full copies a day
- Large databases where daily full backups dominate storage cost
- Instances on ENTERPRISE or ENTERPRISE_PLUS edition (incremental backups are an edition feature)

## Key Configuration

- **12-hour cadence** — `0 2/12 * * *` runs at 02:00 and 14:00 UTC
- **`backupType: INCREMENTAL`** — immutable; switching a schedule between FULL and INCREMENTAL recreates it
- **14-day retention** — a chain's restore window; tune to your recovery objectives

## Customization Notes

- The instance's edition is validated by GCP at create time — a STANDARD instance rejects incremental schedules
- Pair with a weekly FULL schedule (see 03) so restores never depend on one long chain
- Backups inherit the database's encryption posture unless `encryptionConfig` overrides it

## Related Presets

- **01-daily-full-backups** — the simpler baseline for STANDARD instances
- **03-weekly-long-retention** — the complementary full-backup archive
