---
title: "Presets"
description: "Ready-to-deploy configuration presets for Redis Linked Server"
type: "preset-list"
componentSlug: "redis-linked-server"
componentTitle: "Redis Linked Server"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-geo-dr-link"
    rank: "01"
    title: "Geo-Replication DR Link"
    excerpt: "This preset links two Premium caches into the standard disaster-recovery pair: the primary serves reads and writes while continuously replicating to a warm secondary in another region."
  - slug: "02-relink-after-failover"
    rank: "02"
    title: "Re-Link After Failover"
    excerpt: "This preset closes the DR loop: after a regional failover promoted the secondary (by deleting the original link), it re-establishes replication in the opposite direction once the failed region..."
  - slug: "03-cross-manifest-link"
    rank: "03"
    title: "Cross-Manifest Geo Link"
    excerpt: "This preset links to a secondary cache that is NOT managed in the same manifest set -- another team's cache, or one provisioned outside Planton -- by passing its ARM id and region as literal values."
---

# Redis Linked Server Presets

Ready-to-deploy configuration presets for Redis Linked Server. Each preset is a complete manifest you can copy, customize, and deploy.
