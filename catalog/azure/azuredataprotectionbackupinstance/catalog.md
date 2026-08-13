# Azure Data Protection Backup Instance

Creates a Data Protection backup instance -- the binding that puts one datasource (a managed disk, a storage account's blobs, an AKS cluster, a MySQL/PostgreSQL flexible server, or a Data Lake storage account) under a backup policy's protection. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Backup instance** -- ONE of the six datasource-specific ARM bindings, per the variant block set in the spec

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureDataProtectionBackupVault** -- the vault that will hold the backups, WITH a system-assigned identity (the grants below bind to it).
- **An AzureDataProtectionBackupPolicy** -- of the SAME datasource type as the variant you deploy (a disk instance needs a disk policy, on the same vault).
- **The datasource component** -- the AzureManagedDisk, AzureStorageAccount, AzureAksCluster, or flexible-server resource being protected.

### Azure Subscription

- **The vault's managed identity needs datasource permissions BEFORE this deploys** -- Azure validates them at create time. Compose AzureRoleAssignment resources referencing the vault's `system_assigned_identity_principal_id` output: disk needs "Disk Backup Reader" on the disk and "Disk Snapshot Contributor" on the snapshot resource group; blob and Data Lake need "Storage Account Backup Contributor" on the storage account; AKS needs the Backup extension and trusted access; the flexible servers need the vault identity's backup roles.
- **The instance is free** -- backup STORAGE is what bills, according to the policy's retention.
- **Destroying the instance stops protection and deletes the backup data** -- with vault soft delete on, the data lingers 14 days as a soft-deleted item and holds the vault's own deletion for that window.

## Deploy

### Console

Open the deployment store, find **Azure Data Protection Backup Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Disk Backup** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f backup-instance.yaml
```

## After Deploy

The instance's `backup_instance_id` output identifies the protection binding. The first backup runs on the policy's schedule -- protection configured does NOT mean a restore point exists yet.
