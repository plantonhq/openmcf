# AwsBackupSettings — Terraform/OpenTofu module

Manages the AWS Backup settings singletons (`aws_backup_global_settings`
account-wide, `aws_backup_region_settings` per region).

Module facts worth knowing before editing:

- **The arms are count-gated on spec presence** — an omitted arm
  leaves that settings object completely untouched.
- **Both deletes are provider no-ops** — destroy changes nothing at
  AWS; the module never fakes a reset.
- **The maps pass through whole**: the provider requires the full map
  and AWS returns every supported key/type on read — partial maps show
  perpetual plan differences (taught on the spec).
- **`resource_type_management_preference` renders only when set** —
  once set at AWS it cannot be cleared back to unset, only flipped per
  type.
- **Outputs come from data sources** (caller identity + region) — the
  settings resources expose no ARNs; their identities ARE the outputs.

Outputs mirror the Pulumi module key-for-key: `account_id`, `region`.
