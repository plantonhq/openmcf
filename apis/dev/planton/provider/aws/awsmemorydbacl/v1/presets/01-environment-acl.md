# Environment ACL

This preset creates a MemoryDB Access Control List holding one
application's user, ready to grow one user per application. The ACL is the
single attachment point between identities and clusters: clusters attach
the ACL, users join it — access changes never touch the cluster.

## When to Use

- Any production MemoryDB deployment — the alternative ("open-access") is
  unauthenticated and belongs to development only
- One ACL per environment that trusts the same set of application
  identities, shared by every cluster in that environment

## Key Configuration Choices

- **`metadata.name` as the AWS ACL name** — create-time immutable; the
  cluster's `aclName` references the exported `status.outputs.acl_name`
- **Membership as references** — each entry references an AwsMemorydbUser's
  exported user name, so the resource graph shows exactly which identities
  reach which clusters

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<acl-name>` | The ACL name (e.g. `payments-env-acl`); becomes `metadata.name`, max 40 chars | Your naming convention |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |
| `<app-user-name>` | The AwsMemorydbUser resource this ACL grants access to | Your user manifests |

## Common Additions

- More `userNames` entries — one per application
- A read-only analytics user beside the read/write service users

## Related Presets

- **AwsMemorydbUser 01-password-auth / 02-iam-auth** — the identities this
  ACL groups
