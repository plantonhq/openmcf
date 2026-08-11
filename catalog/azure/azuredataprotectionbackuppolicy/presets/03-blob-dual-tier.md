# Blob Dual-Tier Backup

This preset creates the defense-in-depth blob policy: continuous point-in-time restore inside the storage account (30 days) PLUS daily vaulted copies that survive account deletion (90 days, with the first of each month kept a year). Blob storage is the only datasource with two retention tiers -- this preset uses both.

## When to Use

- Blob data whose protection must survive the storage account itself (deletion, compromise, ransomware)
- Anywhere in-account point-in-time restore alone is not enough
- The template for tuning tier durations independently to their different price points

## Key Configuration Choices

- **`operationalDefaultRetentionDuration: P30D`** -- continuous in-account restore, no schedule needed (the tier is always-on)
- **`vaultDefaultRetentionDuration: P90D` + daily intervals** -- the vault tier is SCHEDULED; the intervals are required with it (the spec enforces the provider's own pairing)
- **`retentionRules[monthly]` on `VaultStore`** -- named rules exist only on the vault tier (the operational tier cannot carry them)
- **Drop either tier by removing its duration** -- operational-only needs just the one field; vault-only drops the operational duration

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-backup-vault>` | The AzureDataProtectionBackupVault the policy lives in | Your backup vault resource's name |

The policy is free -- cost follows the protected storage accounts and both tiers' backup storage.
