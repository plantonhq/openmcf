# AzureBackupPolicyVm Terraform Module

## Overview

Creates an Azure Backup policy for IaaS virtual machines -- WHEN VMs under it are backed up (the schedule) and HOW LONG each backup is kept (layered daily/weekly/monthly/yearly retention). The policy is a free configuration object, an ARM child of its Recovery Services vault.

## Resources Created

- `azurerm_backup_policy_vm` -- the policy (`.../vaults/{vault}/backupPolicies/{name}`)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureBackupPolicyVmSpec fields; the resource group and vault-name references arrive as resolved literals

## Outputs

- `backup_policy_id` -- the policy's full ARM ID (what protected VMs bind to)
- `backup_policy_name` -- unique on its vault

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Frequency decides the retention layers** (spec CEL mirrors the provider's CustomizeDiff): Hourly/Daily schedules require `retention_daily`; Weekly requires `retention_weekly` and forbids `retention_daily`. Hourly needs `policy_type` V2.
- **Azure rejects 2-6 days of daily retention** at create (a service rule the provider surfaces only at apply) -- the spec front-loads it: count is 1, or 7+.
- **`policy_type` is ForceNew** -- changing the policy generation replaces the policy. V2 unlocks Hourly schedules, `consistency_type: OnlyCrashConsistent`, and instant-restore retention beyond 5 days.
- **No tags**: ARM backup policies carry no tags -- deliberately no tag map in this module, unlike the vault sibling.
- **Retention forms are exclusive**: monthly/yearly pick backups by week-of-month (weeks+weekdays) OR by month days (days/include_last_days), never both (spec CEL mirrors the provider's ConflictsWith/AtLeastOneOf lattice).

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege control-plane actions the deploying principal needs.
