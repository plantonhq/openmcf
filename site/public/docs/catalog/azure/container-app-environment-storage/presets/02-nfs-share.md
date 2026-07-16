---
title: "NFS Azure Files Share"
description: "This preset registers an NFS Azure Files share (premium FileStorage account) on a Container App Environment. Workloads mount it by declaring an `NFS_AZURE_FILE` volume whose `storage_name` references..."
type: "preset"
rank: "02"
presetSlug: "02-nfs-share"
componentSlug: "container-app-environment-storage"
componentTitle: "Container App Environment Storage"
provider: "azure"
icon: "package"
order: 2
---

# NFS Azure Files Share

This preset registers an NFS Azure Files share (premium FileStorage account) on a Container App Environment. Workloads mount it by declaring an `NFS_AZURE_FILE` volume whose `storage_name` references this registration. NFS traffic never leaves the VNet, so the environment must be VNet-injected.

## When to Use

- POSIX-semantics workloads (file locking, permissions) that SMB serves poorly
- Higher-throughput shared storage on premium FileStorage accounts
- VNet-injected environments where storage traffic must stay private

## Key Configuration Choices

- **NFS path** (`nfsServerUrl`) -- The share is addressed by the account's file endpoint; no account key travels to the environment (access is network-scoped)
- **VNet requirement** -- Pair only with environments created with `infrastructure_subnet_id`; the storage account's network rules must admit the environment's subnet
- **Read-write access** (`accessMode: READ_WRITE`) -- Use READ_ONLY for shared reference data

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<container-app-environment-id>` | ARM ID of the VNet-injected Container App Environment | `AzureContainerAppEnvironment` status outputs |
| `<azure-files-share-name>` | Name of the NFS share | `AzureStorageShare` status outputs (share_name; protocol NFS) |
| `<storage-account-name>` | The premium FileStorage account name | `AzureStorageAccount` status outputs (storage_account_name) |

## Related Presets

- **01-smb-share** -- Use instead for the common SMB case (works without VNet injection)
