---
title: "Presets"
description: "Ready-to-deploy configuration presets for Backup Policy (VM)"
type: "preset-list"
componentSlug: "backup-policy-vm"
componentTitle: "Backup Policy (VM)"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-daily-backup-policy"
    rank: "01"
    title: "Daily Backup Policy"
    excerpt: "This preset creates the everyday production policy: one nightly backup with the classic grandfather-father-son retention ladder -- 30 dailies, 12 weekly Sundays, 12 first-Sunday monthlies, 7..."
  - slug: "02-hourly-enhanced-policy"
    rank: "02"
    title: "Hourly Enhanced Policy"
    excerpt: "This preset creates a low-RPO policy on the V2 (enhanced) generation: a backup every 4 hours inside a 12-hour working-day window, a week of instant-restore snapshots, and age-based archive tiering...."
---

# Backup Policy (VM) Presets

Ready-to-deploy configuration presets for Backup Policy (VM). Each preset is a complete manifest you can copy, customize, and deploy.
