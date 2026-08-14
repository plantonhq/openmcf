# AwsBackupPlan — Terraform/OpenTofu module

Manages a backup plan (`aws_backup_plan`) with its folded selections
(`aws_backup_selection`, for_each keyed by selection name).

Module facts worth knowing before editing:

- **Rules key by name** (a for_each map) so entry order never churns
  the plan; the spec forbids duplicate names.
- **Provider-default knobs render only on an explicit choice**
  (timezone Etc/UTC, start window 60, completion window 180) so the
  module never fights them.
- **Lifecycle day counts pass through presence-typed** — the provider
  cannot transmit an explicit zero, and the spec never pretends
  otherwise.
- **`plan_id` on every selection is the folded plan's own ID** — the
  composition edge is structural, not configurable.
- **Selection deletion order is the provider's problem by design**:
  AWS refuses plan deletion while selections exist, and the provider
  retries while they drain.

Outputs mirror the Pulumi module key-for-key: `plan_id`, `plan_arn`,
`plan_version`, `selection_ids` (keyed by selection name).
