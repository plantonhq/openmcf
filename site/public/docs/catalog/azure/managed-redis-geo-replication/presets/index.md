---
title: "Presets"
description: "Ready-to-deploy configuration presets for Managed Redis Geo Replication"
type: "preset-list"
componentSlug: "managed-redis-geo-replication"
componentTitle: "Managed Redis Geo Replication"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-two-region-pair"
    rank: "01"
    title: "Two-Region Active Pair"
    excerpt: "This preset links two Managed Redis instances in different regions into an active geo-replication group: both accept writes, Azure merges the datasets conflict-free, and applications read and write..."
  - slug: "02-global-mesh"
    rank: "02"
    title: "Global Active Mesh"
    excerpt: "This preset links four Managed Redis instances across continents into one active geo-replication group -- a write-anywhere global cache where every region serves local reads and writes and Azure..."
---

# Managed Redis Geo Replication Presets

Ready-to-deploy configuration presets for Managed Redis Geo Replication. Each preset is a complete manifest you can copy, customize, and deploy.
