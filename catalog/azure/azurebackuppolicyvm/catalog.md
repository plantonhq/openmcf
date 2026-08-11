# Azure Backup Policy (VM)

Creates an Azure Backup policy for IaaS virtual machines -- the schedule and layered retention rules that govern VM backups in a Recovery Services vault. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VM Backup Policy** -- an ARM child of its vault (`.../vaults/{vault}/backupPolicies/{name}`) carrying the schedule, retention ladder, tiering, and instant-restore settings

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureRecoveryServicesVault** -- the policy lives in it (referenced by name).

### Azure Subscription

- **The policy is free** -- cost follows the VMs protected under it and their backup storage.
- **The schedule's frequency decides the legal retention layers** -- Hourly/Daily need `retentionDaily`; Weekly needs `retentionWeekly` instead. The manifest validation walks you through it.
- **Hourly schedules need `policyType: V2`** -- the enhanced policy generation.

## Deploy

### Console

Open the deployment store, find **Azure Backup Policy (VM)**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Daily Backup Policy** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f policy.yaml
```

## After Deploy

The policy's `backup_policy_id` output is what AzureBackupProtectedVm resources bind their VMs to.
