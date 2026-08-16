# AwsBackupSettings — Pulumi module (Go)

Manages the AWS Backup settings singletons (`backup.GlobalSettings`
account-wide, `backup.RegionSettings` per region).

Module facts worth knowing before editing:

- **The arms are nil-gated on spec presence** — an omitted arm leaves
  that settings object completely untouched.
- **Both deletes are provider no-ops** — destroy changes nothing at
  AWS; the module never fakes a reset.
- **The maps pass through whole**: the provider requires the full map
  and AWS returns every supported key/type on read — partial maps show
  perpetual preview differences (taught on the spec).
- **`ResourceTypeManagementPreference` renders only when set** — once
  set at AWS it cannot be cleared back to unset, only flipped per
  type.
- **Outputs come from `GetCallerIdentity` + spec.region** — the
  settings resources expose no ARNs; their identities ARE the outputs.

Outputs mirror the Terraform module key-for-key: `account_id`,
`region`.
