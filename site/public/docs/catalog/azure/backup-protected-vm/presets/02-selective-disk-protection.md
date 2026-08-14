---
title: "Selective Disk Protection"
description: "This preset protects a database VM's OS disk while EXCLUDING the data disks -- the pattern for machines whose data already has native backup tooling (database dumps, log shipping). You stop paying..."
type: "preset"
rank: "02"
presetSlug: "02-selective-disk-protection"
componentSlug: "backup-protected-vm"
componentTitle: "Backup Protected VM"
provider: "azure"
icon: "package"
order: 2
---

# Selective Disk Protection

This preset protects a database VM's OS disk while EXCLUDING the data disks -- the pattern for machines whose data already has native backup tooling (database dumps, log shipping). You stop paying for the same bytes twice; the OS disk restore still recovers the machine itself.

## When to Use

- Database VMs with native backups on their data disks (SQL dumps, WAL archiving)
- VMs carrying rebuildable data disks -- caches, scratch space, tempdb
- Cost-tuning large VMs where full-disk backup doubles an already-covered dataset

## Key Configuration Choices

- **`excludeDiskLuns: [0, 1]`** -- the data disks at those LUNs are skipped; find a disk's LUN in the VM's data-disk list. The OS disk is ALWAYS backed up regardless
- **The inverse exists** -- `includeDiskLuns` backs up the OS disk plus ONLY the listed disks; the two filters are mutually exclusive (validated at manifest time)
- **Restore implication** -- a restore of this VM recovers the OS disk and any included disks; the excluded disks come back from their native tooling. Rehearse that combined restore before trusting it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The VAULT's AzureResourceGroup | Your resource group resource's name |
| `<your-recovery-services-vault>` | The protecting AzureRecoveryServicesVault | Its name output wires automatically |
| `<your-virtual-machine>` | The AzureVirtualMachine to protect | Its vm_id output wires automatically |
| `<your-backup-policy>` | The AzureBackupPolicyVm governing schedule and retention | Its backup_policy_id output wires automatically |
