# AzureDataProtectionBackupInstance Terraform Module

## Overview

Creates a Data Protection backup instance -- the binding that puts one datasource under a backup policy's protection. Exactly one variant block is set in the spec; each variant creates its own provider resource, so ONE resource exists per deployment.

## Resources Created

Exactly one of, per the spec's variant:

- `azurerm_data_protection_backup_instance_blob_storage`
- `azurerm_data_protection_backup_instance_disk`
- `azurerm_data_protection_backup_instance_kubernetes_cluster`
- `azurerm_data_protection_backup_instance_mysql_flexible_server`
- `azurerm_data_protection_backup_instance_postgresql_flexible_server`
- `azurerm_data_protection_backup_instance_data_lake_storage`

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureDataProtectionBackupInstanceSpec fields; the vault, policy, and datasource references arrive as resolved literals

## Outputs

- `backup_instance_id` -- the instance's full ARM ID (coalesced across the six variant resources)
- `backup_instance_name` -- the instance's name, unique on its vault

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Grants precede the instance**: the vault's managed identity must hold the datasource roles BEFORE create (disk: "Disk Backup Reader" on the disk + "Disk Snapshot Contributor" on the snapshot resource group; blob/Data Lake: "Storage Account Backup Contributor" on the account; Kubernetes: the AKS Backup extension + trusted access; MySQL/PostgreSQL: the vault identity's backup roles). An authorization-class create failure means missing or still-propagating role assignments, not a module defect.
- **Nearly everything is ForceNew**; only `backup_policy_id` updates in place -- and the kubernetes_cluster variant has NO update path at all (even the policy binding replaces it).
- **Blob's container list is one-way**: it can change but never be removed entirely once set (the provider ForceNews on clearing it); the module sends null when empty so operational-only protection holds exactly.
- **Destroy deletes the backup data**: with vault soft delete on, the data lingers as a soft-deleted item for 14 days and holds the vault's own deletion for that window.
- **Backup instances carry NO tags** -- the provider has no tags argument on any of the six resources.

## Required Permissions

The deploying principal needs `Microsoft.DataProtection/backupVaults/backupInstances/*` on the vault's resource group (Contributor covers it) -- distinct from the VAULT IDENTITY's datasource grants above, which must also exist.
