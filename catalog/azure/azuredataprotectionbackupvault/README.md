# Overview

The **AzureDataProtectionBackupVault** component creates a Data Protection backup vault -- the safe that MODERN Azure Backup data lives in: managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, and Data Lake storage. Backup policies and backup instances are children of a vault. This is the newer generation alongside the classic Recovery Services vault (which serves VM and file-share backup). The vault itself is free at rest -- cost accrues per protected instance and per GB of backup storage.

## Purpose

- **The modern backup foundation as declarative infrastructure**: datastore tier, redundancy, soft delete, immutability, and encryption -- reviewed and versioned like everything else.
- **Ransomware-resistant by option**: the immutability posture stops backup deletion and retention reduction; AlwaysOn soft delete makes the deletion safety net permanent.
- **Bring-your-own-key encryption**: backup data encrypted with a Key Vault key you control, wired as a typed reference that rotates automatically.
- **Typed references end-to-end**: resource group, Key Vault key, and user-assigned identities wire by reference -- chart-ready.

## Key Features

- Full azurerm v5 surface: the datastore tier (vault/operational/archive), storage redundancy (geo/local/zone) with cross-region restore, the soft-delete state machine with its 14-180-day retention window, the three-state immutability posture, system/user-assigned identity, and customer-managed-key encryption composed as an inline block.
- The service's contracts front-loaded as manifest-time validation: cross-region restore requires geo-redundant storage; encryption requires a system-assigned identity (Azure unwraps the key with it -- hardcoded service-side).
- The three one-way doors recorded honestly on their fields: cross-region restore, Locked immutability, and AlwaysOn soft delete each replace the vault to walk back.

## Use Cases

- **The region's modern backup home**: one vault per region per environment, holding disk, blob, AKS, and database protections.
- **Compliance-grade backups**: immutability plus permanent soft delete for retention nobody -- including admins -- can quietly reduce.
- **Customer-managed-key posture**: organizations whose data-at-rest policy requires their own key material.

## Future Enhancements

- Backup instance kinds (the bindings that put specific disks, blobs, and clusters under a policy's protection) complete the modern Backup story as their contracts land in the catalog.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
