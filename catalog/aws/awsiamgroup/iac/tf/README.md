# AwsIamGroup — Terraform/OpenTofu module

Manages one IAM group (`aws_iam_group`), its declarative membership
(`aws_iam_group_membership`), managed-policy attachments
(`aws_iam_group_policy_attachment`), and inline policies
(`aws_iam_group_policy`).

Module facts worth knowing before editing:

- **The name comes from `metadata.name`** and renames update the
  group IN PLACE at AWS (the ARN recomputes; members and policies
  persist).
- **Membership is ONE resource carrying the whole users list** — the
  AUTHORITATIVE group-centric form: out-of-band additions are removed
  on the next apply, and clearing the list removes the resource (and
  every membership). Rendered only when users are declared (`count`).
- **Attachments use `for_each = toset(managed_policy_arns)`** —
  keyed by the policy ARN itself, so reordering is a no-op, never a
  transient detach/re-attach on a live group.
- **Inline policies live and die with the group**; the heterogeneous
  Struct documents are JSON-encoded into a homogeneous map(string) in
  locals (for_each cannot iterate heterogeneous objects).
- **Nothing here is taggable at AWS** — the module deliberately
  carries no tag map (the one absence against the catalog's tag
  convention, mirrored in the Pulumi module).

Outputs mirror the Pulumi module key-for-key: `group_arn`,
`group_name`, `group_id`.
