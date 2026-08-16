# AwsSsmMaintenanceWindow — Terraform/OpenTofu module

Manages a maintenance window (`aws_ssm_maintenance_window`) with its
folded target registrations (`aws_ssm_maintenance_window_target`) and
tasks (`aws_ssm_maintenance_window_task`), both for_each keyed by
name.

Module facts worth knowing before editing:

- **`enabled` passes through as a tri-state** — unset stays null so
  the provider default (true) rules; only an explicit false creates
  the window paused.
- **Rate controls render only on targeted tasks** — AWS rejects
  `max_concurrency`/`max_errors` on untargeted tasks, so the module
  gates them on `length(targets) > 0`.
- **The invocation union renders exactly one arm** (the spec's CELs
  guarantee one, matching task_type) — each arm is a null-guarded
  dynamic block.
- **`window_id` on every registration is the folded window's own ID**
  — the composition edge is structural, not configurable.

Outputs mirror the Pulumi module key-for-key: `window_id`,
`target_ids` (keyed by target name), `task_ids` (keyed by task name —
the `window_task_id`, not the synthetic composite).
