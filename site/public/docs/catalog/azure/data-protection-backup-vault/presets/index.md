---
title: "Presets"
description: "Ready-to-deploy configuration presets for Data Protection Backup Vault"
type: "preset-list"
componentSlug: "data-protection-backup-vault"
componentTitle: "Data Protection Backup Vault"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-backup-vault"
    rank: "01"
    title: "Standard Backup Vault"
    excerpt: "This preset creates the everyday production vault: the standard vault-store tier on geo-redundant backup storage with cross-region restore, soft delete at its default posture, Microsoft-managed..."
  - slug: "02-immutable-cmk-vault"
    rank: "02"
    title: "Immutable CMK Vault"
    excerpt: "This preset creates the compliance-grade vault: trial-run immutability, a 30-day soft-delete window, and backup data encrypted with your own Key Vault key. For organizations whose backup posture must..."
---

# Data Protection Backup Vault Presets

Ready-to-deploy configuration presets for Data Protection Backup Vault. Each preset is a complete manifest you can copy, customize, and deploy.
