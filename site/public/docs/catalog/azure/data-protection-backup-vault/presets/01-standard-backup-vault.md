---
title: "Standard Backup Vault"
description: "This preset creates the everyday production vault: the standard vault-store tier on geo-redundant backup storage with cross-region restore, soft delete at its default posture, Microsoft-managed..."
type: "preset"
rank: "01"
presetSlug: "01-standard-backup-vault"
componentSlug: "data-protection-backup-vault"
componentTitle: "Data Protection Backup Vault"
provider: "azure"
icon: "package"
order: 1
---

# Standard Backup Vault

This preset creates the everyday production vault: the standard vault-store tier on geo-redundant backup storage with cross-region restore, soft delete at its default posture, Microsoft-managed encryption. The right starting point for a region's disk, blob, AKS, and database backups.

## When to Use

- The first Data Protection vault a region's environment needs -- one home for its modern backup protections
- Standard production posture without customer-managed-key requirements
- Anywhere the paired-region restore capability is worth its modest storage premium

## Key Configuration Choices

- **`datastoreType: VaultStore`** -- the standard tier; fixed at creation (operational-tier and archive-tier vaults serve narrower datasources)
- **`redundancy: GeoRedundant`** -- a copy in the paired region; fixed at creation
- **`crossRegionRestoreEnabled: true`** -- restore in the paired region during an outage; note disabling later REPLACES the vault (one-way)
- **No `softDelete` / `retentionDurationInDays`** -- the service defaults apply (On, 14 days): deleted backups stay recoverable for two weeks
- **No immutability yet** -- run `Unlocked` once retention settings have settled (see the second preset)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the vault lives in | Your resource group resource's name |

The vault is free at rest -- cost starts when instances are protected.
