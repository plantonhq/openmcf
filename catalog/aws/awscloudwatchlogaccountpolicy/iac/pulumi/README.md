# AwsCloudwatchLogAccountPolicy — Pulumi module (Go)

Manages one CloudWatch Logs account-level policy (`cloudwatch.LogAccountPolicy`).

Module facts worth knowing before editing:

- **Identity is the (policy_name, policy_type) pair** — changing either replaces the policy; the provider imports the pair as `policy_name:policy_type`.
- **The document is a Struct, JSON-encoded here** — each policy type carries its own document schema, validated server-side at Put time.
- **`Scope` is pinned to `ALL`** — the only value the provider's enum carries at the pin; a recorded exclusion, never a spec field.
- **`SelectionCriteria` replaces on change** — only the document updates in place.
- **No tags** — account policies are untaggable at AWS (mirrored in the Terraform module).

Outputs mirror the Terraform module key-for-key: `policy_name`, `policy_type`.
