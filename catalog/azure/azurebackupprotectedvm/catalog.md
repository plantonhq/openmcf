# Azure Backup Protected VM

Registers one virtual machine under a backup policy's protection in a Recovery Services vault. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Protected Item** -- the vault-side registration binding the VM to the policy (ARM derives its name from the VM's group and name)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureRecoveryServicesVault** and **an AzureBackupPolicyVm** -- the protection binds the VM to them.
- **An AzureVirtualMachine** -- the machine to protect, in the vault's region.

### Azure Subscription

- **Creation only registers protection** -- the first backup runs on the policy's schedule, not immediately.
- **Destroying deletes backup data** -- the binding's destroy stops protection AND removes the backups (engine defaults, kept deliberately). Retain-on-destroy is an engine-level switch.
- **Cost follows protection** -- a per-instance protected fee plus backup storage GB, both on the vault's bill.

## Deploy

### Console

Open the deployment store, find **Azure Backup Protected VM**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Protect a VM** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f protected-vm.yaml
```

## After Deploy

The first recovery point appears after the policy's next scheduled run (the item shows IRPending until then -- normal). Restores run from the vault in the portal or CLI (`az backup restore`).
