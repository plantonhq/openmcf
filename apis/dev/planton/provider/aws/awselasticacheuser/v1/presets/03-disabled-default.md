# Disabled Default User

This preset creates the mandatory "default" user every ElastiCache user
group must contain. With the access string switched "off", unauthenticated
clients are rejected outright — the standard production posture alongside
per-application password or IAM users.

## When to Use

- Every Redis/Valkey user group — AWS refuses to create a group without a
  member whose user NAME is exactly "default"
- Locking down anonymous access so clients MUST authenticate as a named user
- Pairing with **01-password-auth** or **02-iam-auth** application users in
  the same user group

## Key Configuration Choices

- **`userName: default`** — not optional; AWS enforces this membership rule
  at group-create time (CEL cannot prove name-data across resources, so the
  constraint surfaces at deploy time)
- **`accessString: "off ~* +@all"`** — the user exists but grants nothing;
  clients cannot connect without authenticating as someone else
- **`no-password-required`** — no credential material; this user is never
  meant to be used for authentication

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<default-user-id>` | AWS user id for the default user (e.g. `rbac-default-user`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |

## Common Additions

- One or more application users (presets 01 or 02) referenced in an
  **AwsElasticacheUserGroup** `userIds` list

## Related Presets

- **01-redis-group** (user group) — wires this default user plus application
  users into a group that attaches to a cache
