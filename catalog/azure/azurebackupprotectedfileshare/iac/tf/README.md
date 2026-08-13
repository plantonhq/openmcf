# AzureBackupProtectedFileShare Terraform Module

## Overview

Puts one Azure Files share under a backup policy's protection in a Recovery Services vault. Creating the binding only REGISTERS protection -- the first backup runs on the policy's schedule, not immediately.

## Resources Created

- `azurerm_backup_protected_file_share` -- the protected item (`.../vaults/{vault}/backupFabrics/Azure/protectionContainers/StorageContainer;storage;{sa-rg};{sa-name}/protectedItems/AzureFileShare;{system-name}`)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureBackupProtectedFileShareSpec fields; all five references arrive as resolved literals

## Outputs

- `backup_protected_file_share_id` -- the protected item's full ARM ID (Azure names it by the share's SYSTEM name, not its friendly name)

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **The storage account must already be registered with the vault** (AzureBackupContainerStorageAccount). The provider runs an Inquire pass to discover protectable shares inside the registered container and fails loudly -- `fileshare not found in protectable or protected fileshares` -- when the account is not registered. The spec's default reference wires `source_storage_account_id` THROUGH the registration so it always deploys first.
- **`backup_policy_id` is the only updatable field** -- everything else replaces the protection (a new protected item; the old share's backup data follows the vault's soft-delete rules).
- **Create and delete run up to 80 minutes** (the provider's own timeout class) -- the Inquire discovery and protection configuration are long-running ARM operations.
- **Destroying deletes the backup data** (engines' default destroy behavior; vault soft delete -- always on -- may hold it 14 days, which also delays unregistering the container).
- **No tags**: ARM protected items carry no tags -- deliberately no tag map in this module.

## Required Permissions

The deploying principal needs `Microsoft.RecoveryServices/vaults/backupFabrics/*` on the vault's resource group plus read on the storage account (Contributor covers both).
