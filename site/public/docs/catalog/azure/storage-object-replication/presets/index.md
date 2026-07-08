---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Object Replication"
type: "preset-list"
componentSlug: "storage-object-replication"
componentTitle: "Storage Object Replication"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-cross-region-dr"
    rank: "01"
    title: "Cross-Region DR Replication"
    excerpt: "This preset continuously replicates a container to an account in another region, backfilling everything that already exists -- the blob-level disaster-recovery posture where reads can fail over to..."
  - slug: "02-prefix-scoped-distribution"
    rank: "02"
    title: "Prefix-Scoped Content Distribution"
    excerpt: "This preset replicates only one namespace prefix of a container to a regional account, forward-looking from creation -- the read-local distribution pattern that keeps consumers close to their data..."
---

# Storage Object Replication Presets

Ready-to-deploy configuration presets for Storage Object Replication. Each preset is a complete manifest you can copy, customize, and deploy.
