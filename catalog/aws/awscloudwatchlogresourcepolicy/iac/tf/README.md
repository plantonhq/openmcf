# AwsCloudwatchLogResourcePolicy — Terraform/OpenTofu module

Manages one CloudWatch Logs resource policy (`aws_cloudwatch_log_resource_policy`).

Module facts worth knowing before editing:

- **Exactly one scope** (spec-guaranteed): `policy_name` = a named account-wide policy; `resource_arn` = a policy pinned to one log group. Both are identity — changing either replaces the policy.
- **Revision-guarded updates** — the provider passes AWS's revision ID from state, so concurrent out-of-band edits fail loudly instead of being overwritten. Resource-scoped deletes REQUIRE the tracked revision.
- **The document is a Struct, JSON-encoded here** and diffed semantically as an IAM document.
- **No tags** — resource policies are untaggable at AWS (mirrored in the Pulumi module).

Outputs mirror the Pulumi module key-for-key: `policy_id` (name or ARN — also the import ID; an ARN-shaped import selects the resource scope), `policy_scope`, `revision_id`.
