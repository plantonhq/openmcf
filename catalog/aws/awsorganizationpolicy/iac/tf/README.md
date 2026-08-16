# AwsOrganizationPolicy — Terraform/OpenTofu module

Manages one Organizations policy (`aws_organizations_policy`) with its
folded attachments (`aws_organizations_policy_attachment`, for_each
keyed by resolved target).

Module facts worth knowing before editing:

- **`name` renders `spec.policy_name`** — the explicit name field
  (policy names allow spaces `metadata.name` cannot carry).
- **`type` renders only on an explicit choice** (unset =
  SERVICE_CONTROL_POLICY, the provider default) and forces
  replacement.
- **`content` is `jsonencode(var.spec.content)`** — the spec carries
  the document structured (the catalog's uniform policy-document
  idiom); the provider suppresses JSON-equivalent diffs.
- **Attachment targets arrive resolved** — the platform resolves each
  value-or-reference before the module runs; the resolved target IS
  the for_each key and half of the `{target_id}:{policy_id}` import
  composite.
- **`skip_destroy` is deliberately not rendered** on either resource —
  destroy means detach and delete (the recorded apply-behavior
  exclusion).

Outputs mirror the Pulumi module key-for-key: `policy_id` (the import
ID) and `arn`.
