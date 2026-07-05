---
title: "Redis User Group with Default User"
description: "This preset creates a Redis RBAC user group whose membership references a locked-down \"default\" user plus one or more application users. The group is the attachment unit — replication groups and..."
type: "preset"
rank: "01"
presetSlug: "01-redis-group"
componentSlug: "elasticache-user-group"
componentTitle: "ElastiCache User Group"
provider: "aws"
icon: "package"
order: 1
---

# Redis User Group with Default User

This preset creates a Redis RBAC user group whose membership references a
locked-down "default" user plus one or more application users. The group
is the attachment unit — replication groups and serverless caches reference
the group's id, not individual users.

## When to Use

- Attaching RBAC to a Redis replication group (`user_group_ids`) or
  serverless cache (`user_group_id`)
- Granting or revoking an application's cache access by editing group
  membership — the cache itself never changes
- Any Redis cluster with RBAC enabled

## Key Configuration Choices

- **`metadata.name` as the AWS user group id** — create-time immutable;
  caches reference this id in their `user_group_ids` / `user_group_id` fields
- **Default user by reference** — include an **AwsElasticacheUser** whose
  `userName` is `default` (preset **03-disabled-default**); AWS refuses
  group creation without one
- **`engine: redis`** — must match every member user's engine and every
  cache the group attaches to
- **Membership updates in place** — adding or removing a user id is an
  in-place edit on the group

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-group-id>` | AWS user group id (e.g. `orders-rbac`); becomes `metadata.name` | Your naming convention |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |
| `<default-user-resource-name>` | Name of the locked-down default AwsElasticacheUser | Preset 03-disabled-default manifest |
| `<app-user-resource-name>` | Name of an application AwsElasticacheUser | Preset 01-password-auth or 02-iam-auth manifest |

## Common Additions

- Additional `userIds` entries per application service
- A second group per environment (dev/staging/prod) with different membership

## Related Presets

- **02-valkey-group** — same pattern for Valkey engines
- **03-disabled-default** (user) — the mandatory default user this group references
