---
title: "Protect a File Share"
description: "This preset binds one Azure Files share to a file-share backup policy -- the last link in the backup chain (vault → registration → policy → protection). Creation registers protection; the first..."
type: "preset"
rank: "01"
presetSlug: "01-protect-a-file-share"
componentSlug: "backup-protected-file-share"
componentTitle: "Backup Protected File Share"
provider: "azure"
icon: "package"
order: 1
---

# Protect a File Share

This preset binds one Azure Files share to a file-share backup policy -- the last link in the backup chain (vault → registration → policy → protection). Creation registers protection; the first backup lands at the policy's next scheduled run.

## When to Use

- Every share whose contents should survive deletion, corruption, or ransomware
- App charts that provision a share and its protection together -- restore-ready from day one

## Key Configuration Choices

- **The storage account wires THROUGH its registration** (`AzureBackupContainerStorageAccount`'s echoed `storage_account_id` output) -- the reference guarantees the registration deploys first; a direct account reference can race the discovery pass
- **`backupPolicyId` is the only updatable field** -- vault, account, and share name are identity; changing them replaces the protection
- **No dials beyond the five references** -- schedule and retention live on the policy, shared by every share bound to it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The VAULT's AzureResourceGroup | Your resource group resource's name |
| `<your-recovery-services-vault>` | The protecting AzureRecoveryServicesVault | Its name output wires automatically |
| `<your-backup-container-registration>` | The account's AzureBackupContainerStorageAccount registration | Its echoed account-ID output wires automatically |
| `<your-storage-share>` | The AzureStorageShare to protect | Its name output wires automatically |
| `<your-file-share-backup-policy>` | The AzureBackupPolicyFileShare to bind | Its ARM-ID output wires automatically |

Destroying the protection deletes the backup data (vault soft delete may hold it 14 days) -- teardown runs protections first, then the registration, then the vault.
