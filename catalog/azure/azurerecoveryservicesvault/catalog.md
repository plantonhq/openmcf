# Azure Recovery Services Vault

Creates a Recovery Services vault -- the safe that classic Azure Backup data (VM and file-share backups) and Site Recovery configuration live in. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Recovery Services Vault** -- with its redundancy, immutability, encryption, identity, and monitoring posture
- **Resource Guard Association** (optional) -- multi-user authorization on privileged vault operations, when `resourceGuardId` is set

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureResourceGroup** -- the vault is created in it.

### Azure Subscription

- **The vault is free at rest** -- cost starts when items are protected (per-instance fee plus backup storage GB).
- **Redundancy locks once items are protected** -- decide `storageModeType` before the first protection, not after.
- **Deletion is guarded** -- a vault with protected items refuses to delete; stop and delete protections first.

## Deploy

### Console

Open the deployment store, find **Azure Recovery Services Vault**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Backup Vault** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f vault.yaml
```

## After Deploy

The vault's outputs feed the backup kinds: `recovery_services_vault_name` is what backup policies and protected items address their vault by, and `system_assigned_identity_principal_id` is what Key Vault grants bind to under customer-managed-key encryption.
