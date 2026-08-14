# AwsSsmPatchBaseline — Pulumi module (Go)

Manages a patch baseline (`ssm.PatchBaseline`) with its folded
patch-group registrations (`ssm.PatchGroup`, one per spec entry) and
the optional default designation (`ssm.DefaultPatchBaseline`), all
parented to the baseline.

Module facts worth knowing before editing:

- **`OperatingSystem` renders only on an explicit choice** (unset =
  WINDOWS, the provider default) and forces replacement.
- **`ApproveAfterDays` passes through presence-typed** — 0 (approve
  on release day) is distinct from unset.
- **The designation wires the folded baseline's own resolved OS**
  (`createdBaseline.OperatingSystem.Elem()`) — the standalone
  resource's OS-mismatch failure mode cannot happen here.
- **Destroying the designation RESTORES AWS's predefined default** for
  the OS (the provider looks it up and re-registers it) — the
  reversible designation class.
- **`BaselineId` on every registration is the folded baseline's own
  ID** — the composition edge is structural, not configurable.

Outputs mirror the Terraform module key-for-key: `baseline_id`,
`baseline_arn`, `operating_system` (the resolved OS — also the
designation's import ID).
