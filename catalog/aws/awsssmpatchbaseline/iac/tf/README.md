# AwsSsmPatchBaseline — Terraform/OpenTofu module

Manages a patch baseline (`aws_ssm_patch_baseline`) with its folded
patch-group registrations (`aws_ssm_patch_group`, for_each keyed by
group name) and the optional default designation
(`aws_ssm_default_patch_baseline`, count-gated).

Module facts worth knowing before editing:

- **`operating_system` renders only on an explicit choice** (unset =
  WINDOWS, the provider default) and forces replacement.
- **`approve_after_days` passes through presence-typed** — 0 (approve
  on release day) is distinct from unset.
- **The designation wires the folded baseline's own resolved OS**
  (`aws_ssm_patch_baseline.this.operating_system`) — the standalone
  resource's OS-mismatch failure mode cannot happen here.
- **Destroying the designation RESTORES AWS's predefined default** for
  the OS (the provider looks it up and re-registers it) — the
  reversible designation class.
- **`baseline_id` on every registration is the folded baseline's own
  ID** — the composition edge is structural, not configurable.

Outputs mirror the Pulumi module key-for-key: `baseline_id`,
`baseline_arn`, `operating_system` (the resolved OS — also the
designation's import ID).
