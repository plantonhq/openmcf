# AwsBackupPlan — Pulumi module (Go)

Manages a backup plan (`backup.Plan`) with its folded selections
(`backup.Selection`, one per spec entry, parented to the plan).

Module facts worth knowing before editing:

- **Rules render from the spec list keyed by name**; the spec forbids
  duplicate names so addresses stay stable.
- **Provider-default knobs render only on an explicit choice**
  (timezone Etc/UTC, start window 60, completion window 180) so the
  module never fights them.
- **Lifecycle day counts pass through presence-typed** — the provider
  cannot transmit an explicit zero, and the spec never pretends
  otherwise.
- **`PlanId` on every selection is the folded plan's own ID** — the
  composition edge is structural, not configurable.
- **Selection deletion order is the provider's problem by design**:
  AWS refuses plan deletion while selections exist, and the provider
  retries while they drain.

Outputs mirror the Terraform module key-for-key: `plan_id`,
`plan_arn`, `plan_version`, `selection_ids` (keyed by selection name).
