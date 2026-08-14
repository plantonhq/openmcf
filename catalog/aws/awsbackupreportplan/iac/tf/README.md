# AwsBackupReportPlan — Terraform/OpenTofu module

Manages a Backup Audit Manager report plan (`aws_backup_report_plan`).

Module facts worth knowing before editing:

- **The AWS name is `spec.report_plan_name`**, never `metadata.name` —
  AWS forbids hyphens in report plan names.
- **`report_template` is ForceNew inside the nested block** — a
  template change replaces the whole report plan; the module passes it
  through without softening the semantics.
- **`number_of_frameworks` renders only when positive**, mirroring the
  provider's own send condition (AWS computes it otherwise).
- **Optional lists render null when empty** so both engines send
  identical payloads.

Outputs mirror the Pulumi module key-for-key: `report_plan_arn`.
