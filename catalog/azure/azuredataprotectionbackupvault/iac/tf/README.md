# AzureDataProtectionBackupVault Terraform Module

## Overview

Creates a Data Protection backup vault -- the safe that modern Azure Backup data (managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, Data Lake storage) lives in. The vault is free at rest; cost follows the protected instances and their backup storage.

## Resources Created

- `azurerm_data_protection_backup_vault` -- the vault
- `azurerm_data_protection_backup_vault_customer_managed_key` -- created only when `spec.encryption` is set (the provider's sibling resource that rewrites the vault's own security settings)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureDataProtectionBackupVaultSpec fields; the resource group, Key Vault key, and identity references arrive as resolved literals

## Outputs

- `backup_vault_id` -- the vault's full ARM ID; what backup policies and backup instances reference their vault by
- `backup_vault_name` -- the vault's name
- `system_assigned_identity_principal_id` -- for Key Vault grants under customer-managed-key encryption and datasource RBAC grants, when a system identity is enabled

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **`cross_region_restore_enabled` ships only when true**: the provider ERRORS when the argument is explicitly present -- even as `false` -- on a non-GeoRedundant vault, so an unset spec value reaches the provider as null, never false. Enabling is in-place; disabling replaces the vault (one-way ForceNew).
- **Two more one-way doors**: `immutability: Locked` and `soft_delete: AlwaysOn` are permanent -- leaving either replaces the vault (the provider's ForceNew transitions). `Disabled <-> Unlocked` and `On <-> Off` move freely.
- **CMK encryption can never be removed**: the composed CMK resource's delete is a documented provider no-op -- only deleting the vault removes the encryption. The KEY rotates in place (the one updatable part); versionless key URIs are accepted (`VersionTypeAny`) -- the spec reference's default target, so rotation propagates automatically. Azure unwraps the key with the vault's SYSTEM-assigned identity (the provider hardcodes it).
- **Deletion outlives the API's answer**: the provider polls past Azure's premature "deleted" response until the vault is actually gone (its own workaround for the service bug) -- destroy runs a little longer than the API suggests.
- **`retention_duration_in_days` is the SOFT-DELETE window** (how long deleted backups stay recoverable), not backup retention -- backup retention lives on policies.

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege control-plane actions the deploying principal needs.
