---
title: "Presets"
description: "Ready-to-deploy configuration presets for Recovery Services Vault"
type: "preset-list"
componentSlug: "recovery-services-vault"
componentTitle: "Recovery Services Vault"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-backup-vault"
    rank: "01"
    title: "Standard Backup Vault"
    excerpt: "This preset creates the everyday production vault: geo-redundant backup storage with cross-region restore, Microsoft-managed encryption, all alert switches at their all-on defaults. The right..."
  - slug: "02-immutable-cmk-vault"
    rank: "02"
    title: "Immutable CMK Vault"
    excerpt: "This preset creates the compliance-grade vault: immutability enforced (reversibly, at `Unlocked`), backup data encrypted with your own Key Vault key through a user-assigned identity, and the public..."
---

# Recovery Services Vault Presets

Ready-to-deploy configuration presets for Recovery Services Vault. Each preset is a complete manifest you can copy, customize, and deploy.
