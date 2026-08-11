# Azure Data Protection Backup Policy

Creates a Data Protection backup policy -- the schedule and retention rules for one datasource type (blob storage, managed disks, AKS clusters, MySQL/PostgreSQL flexible servers, or Data Lake storage), expressed as one component with variant blocks. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions exactly ONE of:

- **Blob Storage Backup Policy** -- operational (continuous, in-account) and/or vault (scheduled) retention tiers
- **Disk Backup Policy** -- scheduled incremental snapshots on the operational tier
- **Kubernetes Cluster Backup Policy** -- scheduled AKS backups on the operational tier
- **MySQL / PostgreSQL Flexible Server Backup Policy** -- scheduled full backups on the vault tier
- **Data Lake Storage Backup Policy** -- scheduled ADLS backups on the vault tier

The variant block you set IS the datasource type.

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureDataProtectionBackupVault** -- the policy is a child of it.

### Azure Subscription

- **The policy is free** -- cost follows the protected instances and their backup storage.
- **Policies are immutable** -- every field is fixed at creation; changing anything replaces the policy (backup instances then re-bind).
- **Schedules are ISO-8601 repeating intervals** -- `R/2024-01-01T00:00:00+00:00/P1D` means daily from that instant; durations like `P7D`/`P4M`/`P10Y` set retention.

## Deploy

### Console

Open the deployment store, find **Azure Data Protection Backup Policy**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Daily Disk Backup** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f policy.yaml
```

## After Deploy

The policy's `backup_policy_id` output is what backup instances bind to when putting a disk, blob container, cluster, or database under this policy's protection.
