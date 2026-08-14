---
title: "Presets"
description: "Ready-to-deploy configuration presets for Data Protection Backup Policy"
type: "preset-list"
componentSlug: "data-protection-backup-policy"
componentTitle: "Data Protection Backup Policy"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-daily-disk-backup"
    rank: "01"
    title: "Daily Disk Backup"
    excerpt: "This preset creates the everyday managed-disk policy: daily incremental snapshots at 02:00 UTC, a week of default retention, and the first backup of each week kept for 90 days. The right starting..."
  - slug: "02-aks-cluster-backup"
    rank: "02"
    title: "AKS Cluster Backup"
    excerpt: "This preset creates the Kubernetes cluster policy: backups every four hours with two weeks of default retention, and the first backup of each day kept for eight weeks. The modern-backup capability..."
  - slug: "03-blob-dual-tier"
    rank: "03"
    title: "Blob Dual-Tier Backup"
    excerpt: "This preset creates the defense-in-depth blob policy: continuous point-in-time restore inside the storage account (30 days) PLUS daily vaulted copies that survive account deletion (90 days, with the..."
---

# Data Protection Backup Policy Presets

Ready-to-deploy configuration presets for Data Protection Backup Policy. Each preset is a complete manifest you can copy, customize, and deploy.
