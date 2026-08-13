# AzureBackupPolicyFileShare Terraform Module

## Overview

Creates an Azure Backup policy for Azure Files shares -- WHEN shares under it are backed up (the schedule) and HOW LONG each backup is kept (layered daily/weekly/monthly/yearly retention). The policy is a free configuration object, an ARM child of its Recovery Services vault.

## Resources Created

- `azurerm_backup_policy_file_share` -- the policy (`.../vaults/{vault}/backupPolicies/{name}`)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureBackupPolicyFileShareSpec fields; the resource group and vault-name references arrive as resolved literals

## Outputs

- `backup_policy_id` -- the policy's full ARM ID (what protected file shares bind to)
- `backup_policy_name` -- unique on its vault

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Daily or Hourly only** -- file-share policies have no Weekly schedule (the provider's own surface). The schedule SHAPE follows the frequency (spec CEL): Daily names its `time`; Hourly configures the `hourly` window (start, interval 4/6/8/12, duration 4-24) instead.
- **`retention_daily` is always required** (the provider's own contract) -- the base layer both frequencies retain into; counts run 1-200 (weekly 1-200, monthly 1-120, yearly 1-10 -- shorter than VM retention).
- **`backup_tier` picks where backups live**: `snapshot` (default) keeps share snapshots in the storage account only; `vault-standard` additionally copies backups into the vault. `snapshot_retention_in_days` applies to vault-standard only and must be LESS THAN `retention_daily.count` (spec CEL mirrors the provider's CustomizeDiff).
- **No tags**: ARM backup policies carry no tags -- deliberately no tag map in this module, unlike the vault sibling.
- **Retention forms are exclusive**: monthly/yearly pick backups by week-of-month (weeks+weekdays) OR by month days (days/include_last_days), never both (spec CEL mirrors the provider's ConflictsWith/AtLeastOneOf lattice).

## Required Permissions

The deploying principal needs `Microsoft.RecoveryServices/vaults/backupPolicies/*` on the vault's resource group (Contributor covers it).
