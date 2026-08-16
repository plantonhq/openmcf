# AwsBackupRestoreTestingPlan — Terraform/OpenTofu module

Manages a restore testing plan (`aws_backup_restore_testing_plan`)
with its folded selections (`aws_backup_restore_testing_selection`,
for_each keyed by selection name).

Module facts worth knowing before editing:

- **The AWS names are `spec.plan_name` / `selections[].name`**, never
  `metadata.name` — AWS forbids hyphens and periods here.
- **`restore_testing_plan_name` on every selection is the folded
  plan's own name** — the composition edge is structural, not
  configurable.
- **The conditions block renders only when the spec carries it** —
  AWS's present-but-empty read artifact is normalized by the provider,
  and the module never sends an empty block.
- **One-way knobs pass through as-is** (timezone, windows, exclude
  vaults keep an AWS-side value once set) — the spec comments carry
  the truth; the module never fakes a reset.

Outputs mirror the Pulumi module key-for-key:
`restore_testing_plan_arn`.
