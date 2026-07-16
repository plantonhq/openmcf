---
title: "Presets"
description: "Ready-to-deploy configuration presets for Managed Disk"
type: "preset-list"
componentSlug: "managed-disk"
componentTitle: "Managed Disk"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-premium-data-disk"
    rank: "01"
    title: "Premium Data Disk"
    excerpt: "This preset creates an empty zonal Premium SSD data disk -- the production default for database volumes and any data that must outlive its virtual machine. The VM attaches it through its..."
  - slug: "02-premium-v2-dialed-performance"
    rank: "02"
    title: "Premium SSD v2 with Dialed Performance"
    excerpt: "This preset creates a Premium SSD v2 data disk whose capacity, IOPS, and throughput are provisioned independently -- a small disk with big performance, impossible on the classic per-size tiers. It is..."
  - slug: "03-snapshot-clone"
    rank: "03"
    title: "Snapshot Clone"
    excerpt: "This preset creates a disk by cloning an existing snapshot (or another managed disk) -- the restore, environment-duplication, and troubleshooting workhorse. The clone is a full, independent disk:..."
---

# Managed Disk Presets

Ready-to-deploy configuration presets for Managed Disk. Each preset is a complete manifest you can copy, customize, and deploy.
