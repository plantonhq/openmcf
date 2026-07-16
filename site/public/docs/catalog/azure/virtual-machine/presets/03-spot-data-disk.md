---
title: "Spot Worker with a Persistent Data Disk"
description: "This preset creates a spot (evictable, deeply discounted) Linux worker whose working data lives on a referenced `AzureManagedDisk` that survives eviction, with a system-assigned identity for..."
type: "preset"
rank: "03"
presetSlug: "03-spot-data-disk"
componentSlug: "virtual-machine"
componentTitle: "Virtual Machine"
provider: "azure"
icon: "package"
order: 3
---

# Spot Worker with a Persistent Data Disk

This preset creates a spot (evictable, deeply discounted) Linux worker whose working data lives on a referenced `AzureManagedDisk` that survives eviction, with a system-assigned identity for credential-less access to Azure services and the maximum pre-eviction drain window. It is the cost-optimized shape for interruption-tolerant batch and CI workloads.

## When to Use

- Batch processing, rendering, CI runners, and queue consumers that tolerate interruption
- Any workload where a 60-90% compute discount beats guaranteed availability
- Workers whose in-progress state can be checkpointed to the data disk

## Key Configuration Choices

- **`spot.evictionPolicy: DEALLOCATE`** -- eviction stops compute billing but keeps the VM and its disks; the worker restarts when capacity returns. Use `DELETE` only for fully stateless fleets managed by an orchestrator
- **No `maxBidPrice`** -- the default (-1) pays up to the on-demand price and never gets evicted on price, only on capacity; a cap adds price-based eviction for stricter budgets
- **`dataDiskAttachments` referencing a first-class disk** -- checkpoints and working data live on the `AzureManagedDisk`, which outlives eviction, deletion, and even this VM's replacement
- **`terminationNotification: PT15M`** -- the scheduled event gives the worker its maximum drain window; poll the Instance Metadata Service and checkpoint on the signal
- **`identity: SYSTEM_ASSIGNED`** -- grant the worker queue/storage/registry access with `AzureRoleAssignment` against the VM's principal output; no credentials in the image
- **Cheap OS disk** -- the OS disk is disposable here; `STANDARD_SSD_LRS` is enough (the data disk carries what matters)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the NIC's region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-network-interface-resource-name>` | Planton metadata name of the `AzureNetworkInterface` | Your NIC resource |
| `<your-managed-disk-resource-name>` | Planton metadata name of the `AzureManagedDisk` | Your disk resource |
| `<your-ssh-public-key>` | An OpenSSH public key | `~/.ssh/id_ed25519.pub` or your key management |
