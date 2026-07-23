---
title: "Presets"
description: "Ready-to-deploy configuration presets for CronJob"
type: "preset-list"
componentSlug: "cronjob"
componentTitle: "CronJob"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-nightly-backup"
    rank: "01"
    title: "Nightly Backup CronJob"
    excerpt: "This preset runs a backup every night at 03:00 in an explicit time zone, never overlaps with a still-running previous backup, keeps a week of successful run history, and pulls the database credential..."
  - slug: "02-frequent-sync"
    rank: "02"
    title: "Frequent Sync CronJob"
    excerpt: "This preset runs a synchronization task every 15 minutes with semantics tuned for high-frequency, latest-state-wins work: a stale run is replaced rather than protected, and a run that cannot start..."
  - slug: "03-monthly-report"
    rank: "03"
    title: "Monthly Report CronJob (Indexed)"
    excerpt: "This preset generates a partitioned report on the first of every month: each scheduled run stamps out an Indexed Job with six numbered completions (one per report section, region, or data shard),..."
---

# CronJob Presets

Ready-to-deploy configuration presets for CronJob. Each preset is a complete manifest you can copy, customize, and deploy.
