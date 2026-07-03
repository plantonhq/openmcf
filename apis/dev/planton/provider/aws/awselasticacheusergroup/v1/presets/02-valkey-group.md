# Valkey User Group

This preset creates a Valkey RBAC user group — the same membership and
attachment model as the Redis preset, with `engine: valkey` so AWS accepts
only Valkey-engine users and attaches the group only to Valkey caches.

## When to Use

- RBAC on ElastiCache for Valkey replication groups or serverless caches
- New deployments choosing Valkey over Redis OSS
- Environments standardizing on Valkey while keeping the same RBAC graph
  shape (users → group → cache)

## Key Configuration Choices

- **`engine: valkey`** — create-time immutable; must match every member
  user's `engine` field and the target cache's engine family
- **Same membership rules as Redis** — a user whose `userName` is
  `default` is mandatory; reference it via `status.outputs.user_id`
- **Regional coupling** — every member user and the target cache must live
  in the same region as the group

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-group-id>` | AWS user group id (e.g. `platform-valkey-rbac`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |
| `<default-user-resource-name>` | Valkey-engine default user (`engine: valkey`, `userName: default`) | A user manifest with matching engine |
| `<app-user-resource-name>` | Valkey-engine application user | A user manifest with matching engine |

## Common Additions

- Additional application users per service
- IAM-authenticated users (preset **02-iam-auth** with `engine: valkey`)

## Related Presets

- **01-redis-group** — the Redis-engine equivalent
- **03-disabled-default** (user) — create the default user with
  `engine: valkey` for this group
