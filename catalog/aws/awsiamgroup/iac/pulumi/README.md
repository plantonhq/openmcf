# AwsIamGroup — Pulumi module (Go)

Manages one IAM group (`iam.Group`), its declarative membership
(`iam.GroupMembership`), managed-policy attachments
(`iam.GroupPolicyAttachment`), and inline policies
(`iam.GroupPolicy`).

Module facts worth knowing before editing:

- **The name comes from `metadata.name`** and renames update the
  group IN PLACE at AWS (the ARN recomputes; members and policies
  persist).
- **Membership is ONE resource carrying the whole users list** — the
  AUTHORITATIVE group-centric form: out-of-band additions are removed
  on the next apply, and clearing the list removes the resource (and
  every membership). Rendered only when users are declared.
- **Attachments key by the sanitized policy ARN**, never the list
  index — reordering `managed_policy_arns` must be a no-op, not a
  transient detach/re-attach on a live group.
- **Inline policies live and die with the group**; each
  `google.protobuf.Struct` document is JSON-encoded at render.
- **Nothing here is taggable at AWS** — the module deliberately
  carries no tag map (the one absence against the catalog's tag
  convention, mirrored in the Terraform module).

Outputs mirror the Terraform module key-for-key: `group_arn`,
`group_name`, `group_id`.
