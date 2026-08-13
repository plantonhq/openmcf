# Hourly Vault-Standard Policy

This preset creates the durable low-RPO policy: backups every 4 hours inside a business-hours window, copied INTO the vault (vault-standard tier) so they survive storage-account deletion or compromise, with five days of local snapshots alongside for fast restores.

## When to Use

- Shares whose loss would matter to an auditor -- vaulted copies are what "backup" usually means in compliance terms
- Low-RPO workloads that cannot afford losing a day of changes
- Protection against ransomware or account-level compromise (local snapshots share the account's fate; vaulted copies do not)

## Key Configuration Choices

- **`backupTier: vault-standard`** -- backups additionally copy into the Recovery Services vault
- **`hourly` window** -- interval 4, start 06:00, duration 12: three backups across business hours; hourly schedules have no `time` field (the window replaces it)
- **`snapshotRetentionInDays: 5`** -- local snapshots for fast operational restores; must be STRICTLY LESS than `retentionDaily.count` (the provider rejects equality)
- **Month-days retention grammar** -- monthlies kept on the 1st and the last day; the alternative week-of-month grammar (`weeks`+`weekdays`) is mutually exclusive with it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The vault's AzureResourceGroup | Your resource group resource's name |
| `<your-recovery-services-vault>` | The AzureRecoveryServicesVault the policy lives in | Its name output wires automatically |

Hourly backups multiply kept points quickly -- watch the share's 200-snapshot budget when extending the ladder.
