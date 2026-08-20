# AwsSsmMaintenanceWindow — Pulumi module (Go)

Manages a maintenance window (`ssm.MaintenanceWindow`) with its folded
target registrations (`ssm.MaintenanceWindowTarget`) and tasks
(`ssm.MaintenanceWindowTask`), one per spec entry, parented to the
window.

Module facts worth knowing before editing:

- **`Enabled` passes through as a tri-state** — nil stays unset so the
  provider default (true) rules; only an explicit false creates the
  window paused.
- **Rate controls render only on targeted tasks** — AWS rejects
  `MaxConcurrency`/`MaxErrors` on untargeted tasks, so the module sets
  them only when selectors exist.
- **The invocation union renders exactly one arm** (the spec's CELs
  guarantee one, matching TaskType) — see `taskInvocation`.
- **`WindowId` on every registration is the folded window's own ID** —
  the composition edge is structural, not configurable.

Outputs mirror the Terraform module key-for-key: `window_id`,
`target_ids` (keyed by target name), `task_ids` (keyed by task name —
the `WindowTaskId`, not the synthetic composite).
