---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual Machine"
type: "preset-list"
componentSlug: "virtual-machine"
componentTitle: "Virtual Machine"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-ubuntu-ssh"
    rank: "01"
    title: "Ubuntu Server with SSH Keys"
    excerpt: "This preset creates a zonal Ubuntu 24.04 LTS VM authenticated by SSH keys only, attached to a referenced `AzureNetworkInterface`, with managed boot diagnostics. It is the canonical Linux production..."
  - slug: "02-windows-server"
    rank: "02"
    title: "Windows Server with Trusted Launch"
    excerpt: "This preset creates a zonal Windows Server 2022 VM with trusted launch (secure boot + vTPM), password authentication sourced from a secret reference, and Azure Hybrid Benefit licensing. The Windows..."
  - slug: "03-spot-data-disk"
    rank: "03"
    title: "Spot Worker with a Persistent Data Disk"
    excerpt: "This preset creates a spot (evictable, deeply discounted) Linux worker whose working data lives on a referenced `AzureManagedDisk` that survives eviction, with a system-assigned identity for..."
---

# Virtual Machine Presets

Ready-to-deploy configuration presets for Virtual Machine. Each preset is a complete manifest you can copy, customize, and deploy.
