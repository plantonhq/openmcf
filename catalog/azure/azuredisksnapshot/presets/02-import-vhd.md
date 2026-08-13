# Import VHD

This preset creates a snapshot from an uploaded VHD blob -- the on-ramp for images and disks migrated from outside Azure (or exported from another subscription).

## When to Use

- Migrating VM images built outside Azure (on-prem, another cloud) into managed-disk workflows
- Rehydrating a disk exported to blob storage back into a managed snapshot

## Key Configuration Choices

- **`createOption: Import` + `sourceUri` + `storageAccountId`** -- the working triple for blob sources; the storage account reference carries the read grant, and Azure (not the schema) validates the pairing
- **Page-blob VHDs only** -- the classic fixed-VHD format, not VHDX
- **`diskSizeGb`** -- optional; unset inherits the VHD's size, set it only to create a larger snapshot

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<your-storage-account>` | The Planton name of the `AzureStorageAccount` holding the VHD | Planton console (or replace `valueFrom` with `value:` and a literal account ARM ID) |
| `https://myorgvhds.blob.core.windows.net/vhds/uploaded.vhd` | The VHD page blob's URL | Storage account -> Containers -> the blob's URL |
| `eastus` | The Azure region | Your region strategy |

## Related Presets

- **Incremental Disk Backup** -- the Copy-mode default for disks already in Azure
