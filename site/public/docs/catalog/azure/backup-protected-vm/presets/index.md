---
title: "Presets"
description: "Ready-to-deploy configuration presets for Backup Protected VM"
type: "preset-list"
componentSlug: "backup-protected-vm"
componentTitle: "Backup Protected VM"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-protect-a-vm"
    rank: "01"
    title: "Protect a VM"
    excerpt: "This preset creates the standard protection binding: one VM under one policy in one vault, all wired by reference. All disks are backed up; the protection posture is Azure-managed."
  - slug: "02-selective-disk-protection"
    rank: "02"
    title: "Selective Disk Protection"
    excerpt: "This preset protects a database VM's OS disk while EXCLUDING the data disks -- the pattern for machines whose data already has native backup tooling (database dumps, log shipping). You stop paying..."
---

# Backup Protected VM Presets

Ready-to-deploy configuration presets for Backup Protected VM. Each preset is a complete manifest you can copy, customize, and deploy.
