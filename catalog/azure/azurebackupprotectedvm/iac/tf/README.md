# AzureBackupProtectedVm Terraform Module

## Overview

Registers one virtual machine under a backup policy's protection in a Recovery Services vault. Creation only REGISTERS protection -- the first backup runs on the policy's schedule, not immediately.

## Resources Created

- `azurerm_backup_protected_vm` -- the protected item (`.../vaults/{vault}/backupFabrics/Azure/protectionContainers/.../protectedItems/VM;iaasvmcontainerv2;{vm-rg};{vm-name}` -- ARM derives the name from the VM's group and name)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureBackupProtectedVmSpec fields; the resource group, vault name, VM ID, and policy ID references arrive as resolved literals

## Outputs

- `backup_protected_vm_id` -- the protected item's full ARM ID

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Destroy deletes backup data (deliberate)**: the provider's `features` block stays at defaults, so destroying this resource stops protection AND deletes the backup data. Teams that need data to outlive the binding enable `vm_backup_stop_protection_and_retain_data_on_destroy` (or the suspend variant) in their own engine configuration.
- **Destroy ordering matters for the vault**: the vault refuses to delete while protected items remain -- destroy protections before their vault.
- **`protection_state` is service-managed when unset**: Azure reports transient states (IRPending, ProtectionPaused) that the provider reads back as Protected. `BackupsSuspended` requires the VAULT to be immutable (Unlocked/Locked) -- an apply-time contract Azure checks against the live vault.
- **Disk LUN filters are exclusive**: `exclude_disk_luns` (back up everything except) or `include_disk_luns` (back up OS disk plus only these) -- never both.
- **Soft-deleted ghosts hold the registration**: if the same VM was recently unprotected with data in soft delete, re-creating protection can collide with the ghost; the provider recovers it only when its `recover_soft_deleted_backup_protected_vm` feature is enabled.

## Required Permissions

The deploying principal needs `Microsoft.RecoveryServices/vaults/backupFabrics/*` on the vault's resource group plus read on the source VM (Contributor on both covers it).
