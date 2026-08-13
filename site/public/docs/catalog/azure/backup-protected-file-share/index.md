---
title: "Backup Protected File Share"
description: "Backup Protected File Share deployment documentation"
icon: "package"
order: 100
componentName: "azurebackupprotectedfileshare"
---

# Azure Backup Protected File Share

Puts one Azure Files share under a backup policy's protection in a Recovery Services vault. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Protected File Share** -- an ARM child of the vault (`.../protectionContainers/StorageContainer;storage;{sa-rg};{sa-name}/protectedItems/AzureFileShare;{system-name}`); Azure names it by the share's SYSTEM name

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureBackupContainerStorageAccount** -- the share's storage account REGISTERED with the vault (the default reference wires through it).
- **An AzureBackupPolicyFileShare** -- the schedule and retention the share binds to.
- **The share itself** (AzureStorageShare) -- referenced by name.

### Azure Subscription

- **Creation only registers protection** -- the first backup runs on the policy's schedule, not immediately.
- **Creates and deletes run up to 80 minutes** -- protection configuration is a long-running ARM operation.
- **Destroying deletes the backup data** -- vault soft delete (always on) may hold it 14 days.

## Deploy

### Console

Open the deployment store, find **Azure Backup Protected File Share**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Protect a File Share** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f protected-share.yaml
```

## After Deploy

The share appears under the vault's backup items; the first recovery point lands after the policy's next scheduled run. Re-pointing `backupPolicyId` at a different policy updates in place -- everything else replaces the protection.
