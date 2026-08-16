# AwsBackupReportPlan — Pulumi module (Go)

Manages a Backup Audit Manager report plan (`backup.ReportPlan`).

Module facts worth knowing before editing:

- **The AWS name is `spec.report_plan_name`**, never `metadata.name` —
  AWS forbids hyphens in report plan names.
- **`ReportTemplate` is ForceNew inside the nested block** — a
  template change replaces the whole report plan; the module passes it
  through without softening the semantics.
- **`NumberOfFrameworks` renders only when positive**, mirroring the
  provider's own send condition (AWS computes it otherwise).
- **Optional lists render only when non-empty** so both engines send
  identical payloads.

Outputs mirror the Terraform module key-for-key: `report_plan_arn`.
