---
title: "Presets"
description: "Ready-to-deploy configuration presets for Redis Cache Access Policy Assignment"
type: "preset-list"
componentSlug: "redis-cache-access-policy-assignment"
componentTitle: "Redis Cache Access Policy Assignment"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-identity-data-reader"
    rank: "01"
    title: "Workload Identity Data Reader Grant"
    excerpt: "This preset grants the built-in read-only policy to a user-assigned managed identity -- the standard first grant on the road to a keyless cache."
  - slug: "02-custom-policy-grant"
    rank: "02"
    title: "Custom Policy Grant"
    excerpt: "This preset completes the three-kind composition: the cache defines the boundary, an AzureRedisCacheAccessPolicy defines WHAT is allowed, and this grant defines WHO gets it."
  - slug: "03-human-operator-grant"
    rank: "03"
    title: "Human Operator Grant"
    excerpt: "This preset grants full data access (including admin commands) to a human user or an Entra group -- the break-glass and on-call path that replaces sharing the access keys."
---

# Redis Cache Access Policy Assignment Presets

Ready-to-deploy configuration presets for Redis Cache Access Policy Assignment. Each preset is a complete manifest you can copy, customize, and deploy.
