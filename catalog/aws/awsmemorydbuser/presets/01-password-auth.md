# Password-Authenticated Application User

This preset creates one MemoryDB ACL identity that authenticates with a
password in the AUTH command. Each application gets its own user with an
access string scoping exactly which keys and commands it may touch — the
standard production pattern for MemoryDB's ACL-based authentication.

## When to Use

- Applications connecting to a MemoryDB cluster through a custom ACL
  (instead of the development-only "open-access" ACL)
- Per-service credential isolation — rotating one application's password
  never disturbs the others
- Workloads that cannot use IAM-signed tokens (running outside AWS, or
  connecting to a cluster without TLS)

## Key Configuration Choices

- **`metadata.name` as the AWS user name** — create-time immutable and the
  single identity clients present in `AUTH <name> <password>`; ACLs
  reference it in their membership lists
- **`accessString`** — Redis ACL syntax: key patterns (`~app:*`) plus
  command categories (`+@read`, `-@dangerous`)
- **One or two passwords** — two entries enable zero-downtime rotation
  (add the new password, roll clients, remove the old)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-name>` | The AUTH identity (e.g. `orders-service`); becomes `metadata.name`, max 40 chars | Your naming convention |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |
| `<key-prefix>` | Key namespace this user may touch (e.g. `orders`) | Your data model |
| `$secret/memorydb-user-password` | A managed-secret reference to the password material (16–128 characters) — the field is sensitive, so a `$secret/<slug>` reference belongs here, never plaintext | Your org's managed secrets |

## Common Additions

- A second password in `passwords` for rotation
- Tighter access strings (`+@read` only, `-@dangerous`, narrower key
  patterns)
- Membership in an environment ACL (AwsMemorydbAcl) that clusters attach

## Related Presets

- **02-iam-auth** — IAM-signed tokens instead of long-lived passwords
