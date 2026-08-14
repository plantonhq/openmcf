---
title: "Protect a VM"
description: "This preset creates the standard protection binding: one VM under one policy in one vault, all wired by reference. All disks are backed up; the protection posture is Azure-managed."
type: "preset"
rank: "01"
presetSlug: "01-protect-a-vm"
componentSlug: "backup-protected-vm"
componentTitle: "Backup Protected VM"
provider: "azure"
icon: "package"
order: 1
---

# Protect a VM

This preset creates the standard protection binding: one VM under one policy in one vault, all wired by reference. All disks are backed up; the protection posture is Azure-managed.

## When to Use

- Every production VM -- the default binding to the environment's backup policy
- Chart-composed stacks where the VM and its protection deploy together, so nothing ships unprotected

## Key Configuration Choices

- **Everything by reference** -- vault name, policy ID, and VM ID all resolve from sibling resources; deploy order is handled by the platform
- **No disk filters** -- every disk is backed up; add `excludeDiskLuns` for rebuildable data disks (caches, natively-backed-up databases) to stop paying twice for the same bytes
- **No `protectionState`** -- Azure manages the posture; the item reads `IRPending` until the policy's first scheduled backup runs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The VAULT's AzureResourceGroup (not necessarily the VM's) | Your resource group resource's name |
| `<your-recovery-services-vault>` | The protecting AzureRecoveryServicesVault | Its name output wires automatically |
| `<your-virtual-machine>` | The AzureVirtualMachine to protect (in the vault's region) | Its vm_id output wires automatically |
| `<your-backup-policy>` | The AzureBackupPolicyVm governing schedule and retention | Its backup_policy_id output wires automatically |

Remember: creation only REGISTERS protection -- the first backup runs at the policy's next scheduled window. For critical machines, trigger one manually (`az backup protection backup-now`).
