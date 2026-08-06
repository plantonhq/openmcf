# AwsMemorydbAcl — Design Notes

## Provider mapping

Maps 1:1 onto `aws_memorydb_acl` (Terraform) / `memorydb.Acl` (Pulumi).
The full provider surface is modeled: the name and the user-name set.

## Design decisions

- **The ACL name derives from `metadata.name`.** Required + ForceNew in the
  provider (max 40 characters, AWS-enforced at create); clusters attach the
  ACL by this name, which the resource exports as `status.outputs.acl_name`.
- **Membership as references.** `userNames` is a repeated reference to
  `AwsMemorydbUser.status.outputs.user_name` — the membership-as-refs
  pattern: the ACL is the single place access is granted or revoked, and no
  standalone association kind exists (AWS has no such resource; per-edge
  glue kinds are never modeled).
- **No mandatory member.** Unlike ElastiCache user groups (where AWS
  refuses a group without a user named "default"), MemoryDB accepts an
  empty ACL — so membership is optional and the kind registry declares no
  user prerequisite; composed scenarios wire the user chain through the
  e2e-prerequisites annotation instead.
- **"open-access" is not a resource.** The built-in allow-everything ACL
  always exists per account; clusters reference it by literal value. This
  kind never models or mutates it.

## Deliberately skipped (with reasons)

- `name_prefix` — the naming basis is `metadata.name` (no AWS proto uses
  random-suffix naming).

## Update semantics

Membership updates in place — the provider diffs the user set and issues
one UpdateACL, waiting for the ACL to return to active. The name forces
replacement.
