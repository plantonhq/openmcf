---
title: "IAM-Authenticated User"
description: "This preset creates a Redis/Valkey RBAC identity that authenticates with a short-lived IAM-signed token instead of a long-lived password. No secret material lives in the manifest — clients sign..."
type: "preset"
rank: "02"
presetSlug: "02-iam-auth"
componentSlug: "elasticache-user"
componentTitle: "ElastiCache User"
provider: "aws"
icon: "package"
order: 2
---

# IAM-Authenticated User

This preset creates a Redis/Valkey RBAC identity that authenticates with
a short-lived IAM-signed token instead of a long-lived password. No secret
material lives in the manifest — clients sign tokens with their AWS IAM
identity at connection time.

## When to Use

- Workloads running in AWS (ECS, EKS, Lambda) that already have an IAM role
- Caches with **transit encryption (TLS) enabled** — AWS requires TLS for
  IAM authentication
- Eliminating password rotation toil entirely

## Key Configuration Choices

- **`userName` must equal `metadata.name`** — AWS enforces this coupling
  for IAM auth; the preset uses the same placeholder for both
- **`authenticationMode.type: iam`** — no `passwords` field; CEL rejects
  credential material on non-password types
- **Cache-side requirement** — the attached replication group or serverless
  cache must have transit encryption enabled, and the client's IAM principal
  needs `elasticache:Connect` on both the user ARN and the cache ARN

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-id>` | AWS user id and AUTH user name (must match for IAM auth) | Your naming convention |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |

## Common Additions

- A companion **03-disabled-default** user in the same user group
- Tighter `accessString` scopes (`~app:* +@read` instead of `~* +@read +@write`)

## Related Presets

- **01-password-auth** — password-based auth for workloads outside AWS
- **03-disabled-default** — the mandatory locked-down "default" user
