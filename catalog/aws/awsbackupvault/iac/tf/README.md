# AwsBackupVault — Terraform/OpenTofu module

Manages a backup vault as exactly one of AWS's two vault types
(`aws_backup_vault` OR `aws_backup_logically_air_gapped_vault`) plus
the standard arm's satellites (`aws_backup_vault_lock_configuration`,
`aws_backup_vault_policy`, `aws_backup_vault_notifications`).

Module facts worth knowing before editing:

- **The arms are count-gated on the spec union** — exactly one vault
  resource renders per instance, and the satellites render only inside
  the standard arm (the provider's readers reject other vault types).
- **Satellites attach by the vault's OWN name**, never a foreign
  vault's — the composition edge is structural, not configurable.
- **`changeable_for_days` is write-only at AWS** (the compliance-mode
  opt-in): imports never see it and applies re-assert it.
- **`force_destroy` is deploy-side delete behavior** — never sent as
  AWS state, only honored at destroy to drain recovery points.
- **The policy is a JSON document** rendered via `jsonencode` from the
  spec's Struct (whitespace/key-order changes never diff).

Outputs mirror the Pulumi module key-for-key: `vault_arn`,
`vault_name`.
