---
title: "Presets"
description: "Ready-to-deploy configuration presets for Backup Policy (File Share)"
type: "preset-list"
componentSlug: "backup-policy-file-share"
componentTitle: "Backup Policy (File Share)"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-daily-snapshot-policy"
    rank: "01"
    title: "Daily Snapshot Policy"
    excerpt: "This preset creates the everyday file-share policy: one nightly snapshot backup with the classic grandfather-father-son retention ladder -- 30 dailies, 12 weekly Sundays, 12 first-Sunday monthlies, 5..."
  - slug: "02-hourly-vault-standard-policy"
    rank: "02"
    title: "Hourly Vault-Standard Policy"
    excerpt: "This preset creates the durable low-RPO policy: backups every 4 hours inside a business-hours window, copied INTO the vault (vault-standard tier) so they survive storage-account deletion or..."
---

# Backup Policy (File Share) Presets

Ready-to-deploy configuration presets for Backup Policy (File Share). Each preset is a complete manifest you can copy, customize, and deploy.
