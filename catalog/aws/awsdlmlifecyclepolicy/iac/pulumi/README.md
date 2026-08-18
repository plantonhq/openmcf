# AwsDlmLifecyclePolicy — Pulumi module (Go)

Manages one Data Lifecycle Manager policy (`dlm.LifecyclePolicy`) as a default-XOR-custom mode union.

Module facts worth knowing before editing:

- **Two provider arguments are DERIVED here, never spec surface** — the top-level `DefaultPolicy` (from the default arm's `resource_type`) and `PolicyLanguage` (left to the provider's own derivation). A spec field for either could contradict the configured arm.
- **Default mode omits what AWS defaults** — the provider diff-suppresses the create/retain interval defaults, so those dials are sent only when the spec sets them.
- **The bridge flattens two MaxItems-1 lists to scalars** — `ResourceLocations` and `CreateRule.Times` are single strings in the Pulumi SDK while the spec (mirroring the provider schema) carries lists; this module maps `[0]` accordingly.
- **`Schedule.CopyTags` is ForceNew** — changing it replaces the whole schedule; everything else in a schedule updates in place.
- **Event-based policies react to snapshots shared INTO this account** — `EventSource` + `Action` replace `TargetTags` + `Schedules`; the spec's CELs enforce the provider's ConflictsWith walls before any preview.

Outputs mirror the Terraform module key-for-key: `policy_id` (import ID), `policy_arn`.
