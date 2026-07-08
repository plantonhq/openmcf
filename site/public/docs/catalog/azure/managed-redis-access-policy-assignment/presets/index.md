---
title: "Presets"
description: "Ready-to-deploy configuration presets for Managed Redis Access Policy Assignment"
type: "preset-list"
componentSlug: "managed-redis-access-policy-assignment"
componentTitle: "Managed Redis Access Policy Assignment"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-identity-grant"
    rank: "01"
    title: "Workload Identity Grant"
    excerpt: "This preset grants a user-assigned managed identity data-plane access to a Managed Redis instance -- the secretless workload pattern: the application connects with its object ID as the Redis username..."
  - slug: "02-human-operator"
    rank: "02"
    title: "Human Operator Grant"
    excerpt: "This preset grants a human user -- or an Entra group, covering a whole on-call rotation with one assignment -- data-plane access to a Managed Redis instance. Personal, auditable access with no shared..."
---

# Managed Redis Access Policy Assignment Presets

Ready-to-deploy configuration presets for Managed Redis Access Policy Assignment. Each preset is a complete manifest you can copy, customize, and deploy.
