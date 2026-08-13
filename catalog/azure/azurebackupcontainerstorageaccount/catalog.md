# Azure Backup Container (Storage Account)

Registers a storage account with a Recovery Services vault as a backup container -- the one-time prerequisite that lets the account's file shares be protected. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Backup Container Registration** -- an ARM child of the vault (`.../vaults/{vault}/backupFabrics/Azure/protectionContainers/StorageContainer;storage;{sa-rg};{sa-name}`); ARM derives its name from the storage account's group and name

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureRecoveryServicesVault** -- the registration lives in it (referenced by name).
- **An AzureStorageAccount** -- the account to register (referenced by ARM ID).

### Azure Subscription

- **Registration is free** and moves no data -- cost starts when shares are protected.
- **The account must live in the vault's region** (Azure Files backup is regional).
- **Azure Backup places a resource lock on the account while registered** -- expected, and removed at unregister.
- **Everything is fixed at creation** -- changing any field replaces the registration.

## Deploy

### Console

Open the deployment store, find **Azure Backup Container (Storage Account)**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Register a Storage Account** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f registration.yaml
```

## After Deploy

The registration's `storage_account_id` output (the account ID, echoed) is what AzureBackupProtectedFileShare resources reference for their `sourceStorageAccountId` -- wiring through the registration guarantees it deploys first.
