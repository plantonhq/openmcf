# AzureBackupContainerStorageAccount Terraform Module

## Overview

Registers a storage account with a Recovery Services vault as a backup container -- the one-time prerequisite that lets the account's file shares be protected. One registration per storage-account-and-vault pair; each share then gets its own AzureBackupProtectedFileShare binding. Registration is free and moves no data.

## Resources Created

- `azurerm_backup_container_storage_account` -- the registration (`.../vaults/{vault}/backupFabrics/Azure/protectionContainers/StorageContainer;storage;{sa-rg};{sa-name}`)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureBackupContainerStorageAccountSpec fields; the resource group, vault-name, and storage-account references arrive as resolved literals

## Outputs

- `backup_container_id` -- the registration's full ARM ID
- `storage_account_id` -- the registered account's ARM ID, echoed so protected file shares reference the REGISTRATION for their `source_storage_account_id` (the reference carries both the value and the deploy-order edge -- the provider docs' own wiring pattern)

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Everything is ForceNew** -- ARM has no update on protection containers; changing any field replaces the registration.
- **The account must live in the vault's region** (Azure Files backup is regional) -- Azure rejects cross-region registration at apply time.
- **Azure Backup locks the storage account while registered** (an ARM resource lock protecting the backups' source); unregistering removes it.
- **Unregister refuses while shares are protected** -- destroy the account's AzureBackupProtectedFileShare bindings first. Vault soft delete (always on) may hold deleted protections for 14 days, which can delay the unregister -- the GUIDE carries the teardown recipe.

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege control-plane actions the deploying principal needs.
