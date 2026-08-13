# Azure Data Protection Backup Vault

Creates a Data Protection backup vault -- the safe that modern Azure Backup data (managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, Data Lake storage) lives in. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Data Protection Backup Vault** -- with its datastore tier, redundancy, soft-delete, immutability, and identity posture
- **Customer-Managed-Key Encryption** (optional) -- backup data encrypted with your own Key Vault key, when the `encryption` block is set

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureResourceGroup** -- the vault is created in it.

### Azure Subscription

- **The vault is free at rest** -- cost starts when instances are protected (per-instance fee plus backup storage GB).
- **Datastore tier and redundancy are fixed at creation** -- changing either replaces the vault; decide both up front.
- **Three settings are one-way doors** -- cross-region restore, `Locked` immutability, and `AlwaysOn` soft delete each replace the vault to walk back. Trial-run immutability as `Unlocked` first.

## Deploy

### Console

Open the deployment store, find **Azure Data Protection Backup Vault**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Backup Vault** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f vault.yaml
```

## After Deploy

The vault's outputs feed the backup kinds: `backup_vault_id` is what backup policies and backup instances reference their vault by, and `system_assigned_identity_principal_id` is what Key Vault grants bind to under customer-managed-key encryption (and what datasource RBAC grants bind to for backup permissions).
