---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Share"
type: "preset-list"
componentSlug: "storage-share"
componentTitle: "Storage Share"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-team-file-share"
    rank: "01"
    title: "Team File Share"
    excerpt: "This preset creates a standard SMB file share -- the general-purpose shared drive: Windows mounts it natively, Linux mounts it via cifs, and Azure's default TransactionOptimized tier balances storage..."
  - slug: "02-nfs-premium-share"
    rank: "02"
    title: "NFS Premium Share"
    excerpt: "This preset creates an NFS v4.1 share on a premium FileStorage account -- SSD-backed provisioned performance with real POSIX semantics (hard links, symlinks, chmod) for Linux workloads that SMB..."
  - slug: "03-policy-anchored-access"
    rank: "03"
    title: "Policy-Anchored Access Share"
    excerpt: "This preset creates a share whose external access rides stored access policies (signed identifiers). SAS tokens issued against a policy inherit its window and permissions -- and revoking the policy..."
---

# Storage Share Presets

Ready-to-deploy configuration presets for Storage Share. Each preset is a complete manifest you can copy, customize, and deploy.
