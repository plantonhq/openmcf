---
title: "Presets"
description: "Ready-to-deploy configuration presets for ElastiCache User Group"
type: "preset-list"
componentSlug: "elasticache-user-group"
componentTitle: "ElastiCache User Group"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-redis-group"
    rank: "01"
    title: "Redis User Group with Default User"
    excerpt: "This preset creates a Redis RBAC user group whose membership references a locked-down \"default\" user plus one or more application users. The group is the attachment unit — replication groups and..."
  - slug: "02-valkey-group"
    rank: "02"
    title: "Valkey User Group"
    excerpt: "This preset creates a Valkey RBAC user group — the same membership and attachment model as the Redis preset, with `engine: valkey` so AWS accepts only Valkey-engine users and attaches the group only..."
---

# ElastiCache User Group Presets

Ready-to-deploy configuration presets for ElastiCache User Group. Each preset is a complete manifest you can copy, customize, and deploy.
