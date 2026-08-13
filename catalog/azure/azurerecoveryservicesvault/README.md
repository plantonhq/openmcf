# Overview

The **AzureRecoveryServicesVault** component creates a Recovery Services vault -- the safe that classic Azure Backup data (VM and file-share backups) and Site Recovery configuration live in. Backup policies and protected items are children of a vault; one vault typically serves a region's workloads. The vault itself is free at rest -- cost accrues per protected instance and per GB of backup storage.

## Purpose

- **The backup foundation as declarative infrastructure**: redundancy, immutability, encryption, and network posture -- reviewed and versioned like everything else.
- **Ransomware-resistant by option**: the immutability posture stops backup deletion and retention reduction; the Resource Guard association adds multi-user authorization on privileged operations.
- **Bring-your-own-key encryption**: backup data encrypted with a Key Vault key you control, wired as a typed reference that rotates automatically.
- **Typed references end-to-end**: resource group, Key Vault key, user-assigned identities, and Resource Guard all wire by reference -- chart-ready.

## Key Features

- Full azurerm v5 surface: SKU, storage redundancy (geo/local/zone) with cross-region restore, public network toggle, the three-state immutability posture, system/user-assigned identity, customer-managed-key encryption with infrastructure (double) encryption, all five built-in Azure Monitor alert switches, classic VMware replication, and the composed Resource Guard association.
- The service's contracts front-loaded as manifest-time validation: cross-region restore requires geo-redundant storage; encryption requires an identity; exactly one identity unwraps the key.
- Destroy semantics recorded honestly: a vault with protected items refuses to delete (deliberately -- backup data never disappears as a side effect).

## Use Cases

- **The region's backup home**: one vault per region per environment, holding VM and file-share protections.
- **Compliance-grade backups**: immutability plus multi-user authorization for retention that nobody -- including admins -- can quietly reduce.
- **Customer-managed-key posture**: organizations whose data-at-rest policy requires their own key material.

## Future Enhancements

- File-share backup policy and protection kinds complete the classic Backup story as their contracts land in the catalog.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
