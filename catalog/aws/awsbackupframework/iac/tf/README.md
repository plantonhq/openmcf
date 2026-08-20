# AwsBackupFramework — Terraform/OpenTofu module

Manages a Backup Audit Manager framework (`aws_backup_framework`).

Module facts worth knowing before editing:

- **The AWS name is `spec.framework_name`**, never `metadata.name` —
  AWS forbids hyphens in framework names.
- **Controls key by name** (a for_each map inside the dynamic block)
  so entry order never churns the framework; the spec forbids
  duplicate names.
- **A FAILED deployment is NOT an apply error**: the provider's waiter
  accepts `FAILED` as terminal — the Config-recorder prerequisite
  failure mode is taught in the GUIDE, never masked here.
- **Scope lists render null when empty** so both engines send
  identical payloads.

Outputs mirror the Pulumi module key-for-key: `framework_arn`.
