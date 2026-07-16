---
title: "Presets"
description: "Ready-to-deploy configuration presets for Container App Environment Storage"
type: "preset-list"
componentSlug: "container-app-environment-storage"
componentTitle: "Container App Environment Storage"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-smb-share"
    rank: "01"
    title: "SMB Azure Files Share (Read-Write)"
    excerpt: "This preset registers a standard SMB Azure Files share on a Container App Environment as read-write working storage. Apps and jobs then mount it by declaring an `AZURE_FILE` volume whose..."
  - slug: "02-nfs-share"
    rank: "02"
    title: "NFS Azure Files Share"
    excerpt: "This preset registers an NFS Azure Files share (premium FileStorage account) on a Container App Environment. Workloads mount it by declaring an `NFS_AZURE_FILE` volume whose `storage_name` references..."
---

# Container App Environment Storage Presets

Ready-to-deploy configuration presets for Container App Environment Storage. Each preset is a complete manifest you can copy, customize, and deploy.
