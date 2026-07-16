---
title: "Presets"
description: "Ready-to-deploy configuration presets for Firestore Backup Schedule"
type: "preset-list"
componentSlug: "firestore-backup-schedule"
componentTitle: "Firestore Backup Schedule"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-daily-short-retention"
    rank: "01"
    title: "Daily Short-Retention Backups"
    excerpt: "A daily backup schedule with a one-week retention window — the common companion to PITR for recent restore points without long-term storage cost."
  - slug: "02-weekly-long-retention"
    rank: "02"
    title: "Weekly Long-Retention Archive"
    excerpt: "A weekly backup on Sunday with the maximum 14-week retention — the long-lived archive leg of the daily-plus-weekly pattern."
---

# Firestore Backup Schedule Presets

Ready-to-deploy configuration presets for Firestore Backup Schedule. Each preset is a complete manifest you can copy, customize, and deploy.
