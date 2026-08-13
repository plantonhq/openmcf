# Azure Backup Policy (File Share)

Creates an Azure Backup policy for Azure Files shares -- the schedule and layered retention rules that govern file-share backups in a Recovery Services vault. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **File-Share Backup Policy** -- an ARM child of its vault (`.../vaults/{vault}/backupPolicies/{name}`) carrying the schedule, the retention ladder, and the snapshot/vault-standard tier choice

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureRecoveryServicesVault** -- the policy lives in it (referenced by name).

### Azure Subscription

- **The policy is free** -- cost follows the shares protected under it and their backup storage.
- **Daily or Hourly only** -- file-share policies have no Weekly schedule. Daily schedules name a `time`; Hourly schedules configure an `hourly` window instead. The manifest validation walks you through it.
- **`backupTier: vault-standard`** additionally copies backups into the vault; its `snapshotRetentionInDays` must be less than `retentionDaily.count`.

## Deploy

### Console

Open the deployment store, find **Azure Backup Policy (File Share)**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Daily Snapshot Policy** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f policy.yaml
```

## After Deploy

The policy's `backup_policy_id` output is what AzureBackupProtectedFileShare resources bind their shares to.
