# AwsDlmLifecyclePolicy — Terraform/OpenTofu module

Manages one Data Lifecycle Manager policy (`aws_dlm_lifecycle_policy`) as a default-XOR-custom mode union.

Module facts worth knowing before editing:

- **Two provider arguments are DERIVED here, never spec surface** — the top-level `default_policy` (from the default arm's `resource_type`) and `policy_language` (SIMPLIFIED for default mode, STANDARD for custom — left to the provider's own derivation). A spec field for either could contradict the configured arm.
- **Default mode omits what AWS defaults** — the provider diff-suppresses `create_interval` 0→1 and `retain_interval` 0→7, so this module sends those dials only when the spec sets them.
- **`schedule.copy_tags` is ForceNew** — changing it replaces the whole schedule; everything else in a schedule updates in place.
- **The policy has no volume/snapshot edges** — it targets by tags at fire time; `execution_role_arn` (an AwsIamRole reference) is its only dependency.
- **Event-based policies react to snapshots shared INTO this account** — `event_source` + `action` replace `target_tags` + `schedule`; the spec's CELs enforce the provider's ConflictsWith walls before AWS sees the plan.

Outputs mirror the Pulumi module key-for-key: `policy_id` (import ID), `policy_arn`.
