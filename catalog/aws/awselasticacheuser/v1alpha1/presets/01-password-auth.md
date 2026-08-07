# Password-Authenticated Application User

This preset creates one Redis/Valkey RBAC identity that authenticates with
a password in the AUTH command. Each application gets its own user with an
access string scoping exactly which keys and commands it may touch — the
standard production pattern for ElastiCache RBAC.

## When to Use

- Applications connecting to a Redis or Valkey replication group or
  serverless cache with RBAC enabled
- Per-service credential isolation — rotating one application's password
  never disturbs the others
- Workloads that cannot use IAM-signed tokens (outside AWS, or caches
  without transit encryption)

## Key Configuration Choices

- **`metadata.name` as the AWS user id** — create-time immutable; this is
  what user groups reference in `userIds`, not `metadata.id`
- **`userName` for AUTH** — what clients send in `AUTH <userName>
  <password>`; several users may share a user name (AWS unions their
  credentials at authentication time)
- **`accessString`** — Redis ACL syntax: key patterns (`~app:*`) plus
  command categories (`+@read`, `-@dangerous`)
- **One or two passwords** — two entries enable zero-downtime rotation
  (add the new password, roll clients, remove the old)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-id>` | AWS user id (e.g. `orders-service`); becomes `metadata.name` | Your naming convention |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |
| `<auth-user-name>` | Name clients present in AUTH (often matches the user id) | Your application identity |
| `<key-prefix>` | Key namespace this user may touch (e.g. `orders`) | Your data model |
| `<strong-password-16-to-128-chars>` | Password material (16–128 printable characters) | Your secret manager |

## Common Additions

- A second password in `passwords` for rotation
- A companion **03-disabled-default** user in the same user group
- Tighter access strings (`+@read` only, `-@dangerous`, narrower key patterns)

## Related Presets

- **02-iam-auth** — IAM-signed tokens instead of long-lived passwords
- **03-disabled-default** — the mandatory locked-down "default" user every
  user group must include
