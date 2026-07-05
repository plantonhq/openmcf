---
title: "Daily Full Backups"
description: "Creates a complete, self-contained backup of the database every day at 02:00 UTC and keeps each one for 31 days. The straightforward baseline every production database should start from."
type: "preset"
rank: "01"
presetSlug: "01-daily-full-backups"
componentSlug: "spanner-backup-schedule"
componentTitle: "Spanner Backup Schedule"
provider: "gcp"
icon: "package"
order: 1
---

# Daily Full Backups

Creates a complete, self-contained backup of the database every day at 02:00 UTC and keeps each one for 31 days. The straightforward baseline every production database should start from.

## When to Use

- The default protection posture for any production Spanner database
- Restore points at day granularity are sufficient
- STANDARD-edition instances (incremental backups require ENTERPRISE or above)

## Key Configuration

- **Instance and database by reference** — composes against the `GcpSpannerInstance` and `GcpSpannerDatabase` resources' name outputs
- **Daily cadence** — `0 2 * * *`, evaluated in UTC; Spanner accepts 12-hour, daily, weekly, or monthly frequencies
- **31-day retention** — `2678400s`; applies to backups created after any change

## Customization Notes

- `metadata.name` doubles as the schedule name when `scheduleName` is omitted
- Backups inherit the database's encryption posture unless `encryptionConfig` overrides it
- Retention accepts up to 366 days (`31622400s`)

## Related Presets

- **02-incremental-enterprise** — cheaper storage via incremental chains on ENTERPRISE instances
- **03-weekly-long-retention** — weekly cadence with long retention for compliance archives
