---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Data Lake Gen2 Filesystem"
type: "preset-list"
componentSlug: "storage-data-lake-gen2-filesystem"
componentTitle: "Storage Data Lake Gen2 Filesystem"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-data-lake-zone"
    rank: "01"
    title: "Data Lake Zone Filesystem"
    excerpt: "This preset creates one zone of a medallion-style data lake (raw, curated, gold) as its own filesystem -- the grain that gives each zone its own POSIX posture and its own RBAC scope."
  - slug: "02-team-scoped-workspace"
    rank: "02"
    title: "Team-Scoped Workspace Filesystem"
    excerpt: "This preset creates a self-service analytics workspace owned by an Entra group -- the team writes freely inside its filesystem, and membership changes happen in Entra, never in storage config."
  - slug: "03-regulated-zone-cmk"
    rank: "03"
    title: "Regulated Zone with Customer-Managed Key"
    excerpt: "This preset creates a compliance zone whose data encrypts under YOUR Key Vault key (via an encryption scope), with a deny-by-default root ACL -- key custody and least-privilege for just the regulated..."
---

# Storage Data Lake Gen2 Filesystem Presets

Ready-to-deploy configuration presets for Storage Data Lake Gen2 Filesystem. Each preset is a complete manifest you can copy, customize, and deploy.
