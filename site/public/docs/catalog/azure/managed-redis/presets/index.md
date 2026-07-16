---
title: "Presets"
description: "Ready-to-deploy configuration presets for Managed Redis"
type: "preset-list"
componentSlug: "managed-redis"
componentTitle: "Managed Redis"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-keyless-balanced"
    rank: "01"
    title: "Keyless Balanced Cache"
    excerpt: "This preset creates the Managed Redis default posture done right: a general-purpose Balanced instance with high availability and NO access keys -- every client authenticates with Microsoft Entra..."
  - slug: "02-search-json-enterprise"
    rank: "02"
    title: "Search + JSON Document Store"
    excerpt: "This preset creates a Managed Redis instance with the RediSearch and RedisJSON modules -- a queryable, full-text-searchable JSON document store with Redis latency. Modules exist only on Managed..."
  - slug: "03-geo-replicated"
    rank: "03"
    title: "Geo-Replicated Group Member"
    excerpt: "This preset creates one member of an ACTIVE geo-replication group -- multi-primary Redis where every region accepts writes and Azure merges the datasets with conflict-free semantics. Deploy one per..."
---

# Managed Redis Presets

Ready-to-deploy configuration presets for Managed Redis. Each preset is a complete manifest you can copy, customize, and deploy.
