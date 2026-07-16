---
title: "Presets"
description: "Ready-to-deploy configuration presets for Disk Encryption Set"
type: "preset-list"
componentSlug: "disk-encryption-set"
componentTitle: "Disk Encryption Set"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-system-assigned-rotation"
    rank: "01"
    title: "System-Assigned Set with Auto Key Rotation"
    excerpt: "This preset creates a disk encryption set with a system-assigned identity and automatic key rotation -- the recommended posture. The set follows the key's latest version as it rotates (referencing..."
  - slug: "02-user-assigned-preprovisioned"
    rank: "02"
    title: "User-Assigned Set with Pre-Provisioned Access"
    excerpt: "This preset creates a disk encryption set backed by a user-assigned managed identity. Because the identity exists independently of the set, you grant it Key Vault crypto access BEFORE creating the..."
---

# Disk Encryption Set Presets

Ready-to-deploy configuration presets for Disk Encryption Set. Each preset is a complete manifest you can copy, customize, and deploy.
