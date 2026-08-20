# AwsBackupVault — Pulumi module (Go)

Manages a backup vault as exactly one of AWS's two vault types
(`backup.Vault` OR `backup.LogicallyAirGappedVault`) plus the standard
arm's satellites (`backup.VaultLockConfiguration`,
`backup.VaultPolicy`, `backup.VaultNotifications`).

Module facts worth knowing before editing:

- **The arms are nil-gated on the spec union** — exactly one vault
  resource renders per instance, and the satellites render only inside
  the standard arm (the provider's readers reject other vault types).
- **Satellites attach by the vault's OWN name** (the created vault's
  Name output), never a foreign vault's.
- **`changeable_for_days` is write-only at AWS** (the compliance-mode
  opt-in): imports never see it and applies re-assert it.
- **`force_destroy` is deploy-side delete behavior** — rendered only
  when true, honored at destroy to drain recovery points.
- **The policy Struct marshals to normalized JSON** (matching the
  Terraform module's `jsonencode`) so both engines send identical
  payloads.

Outputs mirror the Terraform module key-for-key: `vault_arn`,
`vault_name`.
