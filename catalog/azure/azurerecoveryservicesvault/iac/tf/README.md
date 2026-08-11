# AzureRecoveryServicesVault Terraform Module

## Overview

Creates a Recovery Services vault -- the safe that classic Azure Backup data (VM and file-share backups) and Site Recovery configuration live in. The vault is free at rest; cost follows the protected items and their backup storage.

## Resources Created

- `azurerm_recovery_services_vault` -- the vault
- `azurerm_recovery_services_vault_resource_guard_association` -- created only when `spec.resource_guard_id` is set (Multi-User Authorization; ARM pins the association's name to the literal `VaultProxy`)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureRecoveryServicesVaultSpec fields; the resource group, Key Vault key, identity, and Resource Guard references arrive as resolved literals

## Outputs

- `recovery_services_vault_id` -- the vault's full ARM ID
- `recovery_services_vault_name` -- what backup policies and protected items address their vault by
- `system_assigned_identity_principal_id` -- for Key Vault grants under customer-managed-key encryption, when a system identity is enabled
- `resource_guard_association_id` -- the association's ARM ID when the guard arm is configured; empty otherwise

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Destroy semantics (deliberate)**: the provider's `features` block stays at defaults, so deleting the vault FAILS while protected items remain inside it (`purge_protected_items_from_vault_on_destroy` off). Stop and delete protections first, then the vault -- the honest, least-surprising posture for backup data.
- **Immutability is a one-way street at the end**: transitions run Disabled <-> Unlocked -> Locked; leaving Locked replaces the vault. The provider stages a requested Locked through Unlocked automatically (ARM refuses Locked at create).
- **Cross-region restore is one-way ForceNew**: enabling updates in place; disabling replaces the vault. It also requires geo-redundant storage (spec CEL mirrors the provider's apply-time check).
- **CMK encryption ratchets**: once enabled it can never be disabled, `infrastructure_encryption_enabled` can never change, and the `sku` freezes (the provider's own update guards). Versionless key URIs are accepted (`VersionTypeAny`) -- the spec reference's default target, so key rotation propagates automatically.
- **Identity never downgrades**: once set, removing the identity or switching to the opposite single flavor is rejected by the provider (Azure's CMK guidance).
- **Monitoring**: all five alert switches default ON both provider- and service-side. This engine honors every switch, including the three v5-new ones the Pulumi engine cannot express (its PARITY-EXCEPTION).

## Required Permissions

The deploying principal needs `Microsoft.RecoveryServices/vaults/*` on the resource group (Contributor covers it). Customer-managed-key encryption additionally needs wrap/unwrap access on the Key Vault key for the vault's identity.
