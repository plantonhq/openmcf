# AzureDataProtectionBackupPolicy Terraform Module

## Overview

Creates a Data Protection backup policy -- the schedule and retention rules for one datasource type, expressed as one spec with six variant blocks. Exactly one variant is set (validated at admission); the module creates exactly one provider resource.

## Resources Created (exactly one, by variant)

- `azurerm_data_protection_backup_policy_blob_storage` -- the `blob_storage` variant
- `azurerm_data_protection_backup_policy_disk` -- the `disk` variant
- `azurerm_data_protection_backup_policy_kubernetes_cluster` -- the `kubernetes_cluster` variant
- `azurerm_data_protection_backup_policy_mysql_flexible_server` -- the `mysql_flexible_server` variant
- `azurerm_data_protection_backup_policy_postgresql_flexible_server` -- the `postgresql_flexible_server` variant
- `azurerm_data_protection_backup_policy_data_lake_storage` -- the `data_lake_storage` variant

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureDataProtectionBackupPolicySpec fields; the vault reference arrives as a resolved literal ARM ID

## Outputs

- `backup_policy_id` -- the policy's full ARM ID (coalesced across the six variant resources); what backup instances bind their policy by
- `backup_policy_name` -- the policy's name, unique on its vault

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Every variant is immutable after create**: the provider ships no update path (near-total ForceNew) -- changing anything replaces the policy, and bound backup instances re-bind to the replacement.
- **One vault reference, three provider addressing styles**: most variants take the vault's ARM ID directly; the `kubernetes_cluster` resource takes vault NAME + resource group (both derived from the ID in `locals.tf` -- ARM IDs are structured); the `data_lake_storage` resource takes the same ID under its own argument name.
- **Data Lake priorities are order-derived**: the provider stamps rule N with priority N+1 (there is no priority argument on that resource) -- list rules in priority order.
- **Stores are pinned per datasource** (the service's own surface today): disk and Kubernetes retain on `OperationalStore`; the flexible servers, Data Lake, and blob's named rules on `VaultStore`; blob's dual default tiers select their stores by which duration field is set. The spec vocabularies mirror these exactly and widen if the service does.
- **Empty-string and empty-list optionals are sent as null** so provider defaults and optional blocks behave exactly as they do in hand-written Terraform.

## Required Permissions

The deploying principal needs `Microsoft.DataProtection/backupVaults/backupPolicies/*` on the vault's resource group (Contributor covers it).
