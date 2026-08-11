# Daily Snapshot Policy

This preset creates the everyday file-share policy: one nightly snapshot backup with the classic grandfather-father-son retention ladder -- 30 dailies, 12 weekly Sundays, 12 first-Sunday monthlies, 5 first-Sunday-of-January yearlies. The right default for most shares.

## When to Use

- The first policy a vault needs -- one schedule most shares share
- Workloads whose recovery point objective tolerates up to 24 hours
- Shares where fast snapshot restores matter more than surviving storage-account loss (for that, use the vault-standard preset)

## Key Configuration Choices

- **Snapshot tier (the default)** -- backups live as share snapshots in the storage account itself: fast restores, no vault copy
- **`time: "23:00"`** -- the nightly backup start; must land on the hour or half past
- **30/12/12/5 ladder** -- 59 kept points against the share's 200-snapshot budget; each kept point is incremental backup storage
- **Week-of-month retention grammar** -- "First Sunday" for monthlies/yearlies; the alternative month-days grammar (`days`/`includeLastDays`) is mutually exclusive with it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The vault's AzureResourceGroup | Your resource group resource's name |
| `<your-recovery-services-vault>` | The AzureRecoveryServicesVault the policy lives in | Its name output wires automatically |

The policy protects nothing by itself: register the share's storage account (AzureBackupContainerStorageAccount), then bind each share (AzureBackupProtectedFileShare).
