---
title: "Presets"
description: "Ready-to-deploy configuration presets for Redis Cache"
type: "preset-list"
componentSlug: "redis-cache"
componentTitle: "Redis Cache"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-entra"
    rank: "01"
    title: "Standard Cache with Entra Authentication"
    excerpt: "This preset creates the production default: a replicated Standard-tier cache with Microsoft Entra token authentication enabled alongside the access keys -- the on-ramp to a fully keyless cache."
  - slug: "02-premium-enterprise"
    rank: "02"
    title: "Premium Enterprise Cache"
    excerpt: "This preset creates a private, zone-spread, clustered Premium cache with hourly RDB persistence authenticated by managed identity -- the shape a regulated or high-scale workload reaches for."
  - slug: "03-development"
    rank: "03"
    title: "Development Cache"
    excerpt: "This preset creates the smallest, cheapest cache Azure offers -- a single Basic-tier C0 node with an IP allow-list on the public endpoint."
---

# Redis Cache Presets

Ready-to-deploy configuration presets for Redis Cache. Each preset is a complete manifest you can copy, customize, and deploy.
