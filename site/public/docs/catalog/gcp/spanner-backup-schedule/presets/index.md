---
title: "Presets"
description: "Ready-to-deploy configuration presets for Spanner Backup Schedule"
type: "preset-list"
componentSlug: "spanner-backup-schedule"
componentTitle: "Spanner Backup Schedule"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-daily-full-backups"
    rank: "01"
    title: "Daily Full Backups"
    excerpt: "Creates a complete, self-contained backup of the database every day at 02:00 UTC and keeps each one for 31 days. The straightforward baseline every production database should start from."
  - slug: "02-incremental-enterprise"
    rank: "02"
    title: "Incremental Backups (Enterprise)"
    excerpt: "Backs the database up every 12 hours using incremental chains — each backup stores only the changes since the previous one, at a fraction of full-backup storage cost with the same restore semantics."
  - slug: "03-weekly-long-retention"
    rank: "03"
    title: "Weekly Long-Retention Archive"
    excerpt: "Creates a full backup every Sunday and keeps each one for 366 days — the maximum Spanner allows — encrypted with an explicit customer-managed key. The compliance-archive pattern that runs alongside a..."
---

# Spanner Backup Schedule Presets

Ready-to-deploy configuration presets for Spanner Backup Schedule. Each preset is a complete manifest you can copy, customize, and deploy.
