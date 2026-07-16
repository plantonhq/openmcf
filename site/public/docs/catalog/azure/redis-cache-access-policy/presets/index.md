---
title: "Presets"
description: "Ready-to-deploy configuration presets for Redis Cache Access Policy"
type: "preset-list"
componentSlug: "redis-cache-access-policy"
componentTitle: "Redis Cache Access Policy"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-read-only-prefix"
    rank: "01"
    title: "Read-Only Policy on a Key Prefix"
    excerpt: "This preset defines the most common custom policy: read-only access scoped to one application's key prefix -- finer than the built-in \"Data Reader\", which reads every key on the cache."
  - slug: "02-app-writer-no-admin"
    rank: "02"
    title: "Application Writer Without Admin Commands"
    excerpt: "This preset defines the standard application-workload policy: full data read/write on the app's key prefix, with Redis's dangerous command category carved out."
  - slug: "03-queue-worker-commands"
    rank: "03"
    title: "Queue Worker with Single-Command Grants"
    excerpt: "This preset shows the finest grain the ACL syntax offers: individual command grants. The identity can push, pop, and measure one queue key -- and do literally nothing else."
---

# Redis Cache Access Policy Presets

Ready-to-deploy configuration presets for Redis Cache Access Policy. Each preset is a complete manifest you can copy, customize, and deploy.
