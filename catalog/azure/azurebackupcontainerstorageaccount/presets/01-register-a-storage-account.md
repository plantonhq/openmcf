# Register a Storage Account

This preset registers one storage account with a Recovery Services vault as a backup container -- the one-time prerequisite before any of the account's file shares can be protected. Free, data-free, and deliberately small.

## When to Use

- Always, exactly once per storage-account-and-vault pair -- before the account's first AzureBackupProtectedFileShare
- In charts that provision a storage account and protect its shares in one deploy (the registration node makes the ordering automatic)

## Key Configuration Choices

- **Three references, no dials** -- the vault's resource group, the vault by name, and the storage account by ARM ID; ARM derives the registration's own name
- **Everything is fixed at creation** -- changing any field replaces the registration
- **Same region as the vault** -- Azure Files backup is regional; cross-region registration fails at apply

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The vault's AzureResourceGroup | Your resource group resource's name |
| `<your-recovery-services-vault>` | The AzureRecoveryServicesVault to register with | Its name output wires automatically |
| `<your-storage-account>` | The AzureStorageAccount holding the shares | Its ARM ID output wires automatically |

Protected shares should reference THIS resource's `storage_account_id` output (their default reference does) so the registration always deploys first.
