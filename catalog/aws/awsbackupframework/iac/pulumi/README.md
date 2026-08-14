# AwsBackupFramework — Pulumi module (Go)

Manages a Backup Audit Manager framework (`backup.Framework`).

Module facts worth knowing before editing:

- **The AWS name is `spec.framework_name`**, never `metadata.name` —
  AWS forbids hyphens in framework names.
- **Controls render from the spec list keyed by name**; the spec
  forbids duplicate names so addresses stay stable.
- **A FAILED deployment is NOT an apply error**: the provider's waiter
  accepts `FAILED` as terminal — the Config-recorder prerequisite
  failure mode is taught in the GUIDE, never masked here.
- **Scope lists render only when non-empty** so both engines send
  identical payloads.

Outputs mirror the Terraform module key-for-key: `framework_arn`.
