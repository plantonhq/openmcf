---
title: "Blob Backup"
description: "This preset protects a storage account's blob data with vault-tier backups of named containers (plus whatever operational-tier protection the referenced blob policy configures)."
type: "preset"
rank: "03"
presetSlug: "03-blob-backup"
componentSlug: "data-protection-backup-instance"
componentTitle: "Data Protection Backup Instance"
provider: "azure"
icon: "package"
order: 3
---

# Blob Backup

This preset protects a storage account's blob data with vault-tier backups of named containers (plus whatever operational-tier protection the referenced blob policy configures).

## When to Use

- Business-critical blob data (documents, exports, event archives) that needs restore points independent of the storage account itself
- Pairing vault-tier copies with the operational tier's continuous in-account protection for a two-layer posture
- Protecting against account-level incidents -- vaulted copies live in the backup vault, not the account

## Key Configuration Choices

- **Named containers** -- vault-tier backup is per-container; list what matters. ONE-WAY once set: the list can change but never be removed entirely (Azure's own contract)
- **Omit the container list only for operational-only policies** -- the operational tier protects the whole account continuously inside the account
- **Grants precede the instance** -- the vault identity needs "Storage Account Backup Contributor" on the account; compose an AzureRoleAssignment referencing the vault's `system_assigned_identity_principal_id` output

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-backup-vault>` | The AzureDataProtectionBackupVault holding the backups | The vault component's name |
| `<your-blob-policy>` | An AzureDataProtectionBackupPolicy with the `blobStorage` variant, on the same vault | The policy component's name |
| `<your-storage-account>` | The AzureStorageAccount whose blobs are protected | The storage-account component's name |

The instance is free; vaulted backup storage bills per the policy's retention.
