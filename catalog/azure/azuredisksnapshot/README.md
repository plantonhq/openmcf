# Overview

The **AzureDiskSnapshot** component deploys a managed disk snapshot -- a point-in-time copy of a disk used for backup, cloning, and as the source of gallery image versions (AzureComputeGalleryImage's `versions[].os_disk_snapshot_id`).

## Purpose

- **The cheapest point-in-time copy**: an INCREMENTAL snapshot stores only the delta since the previous snapshot of the same disk, on standard storage -- the right default for backup chains.
- **The image pipeline's handoff artifact**: snapshot a prepared OS disk, then publish it as a gallery image version -- the chart-native golden-image chain.
- **A restore and clone source**: new managed disks create from a snapshot; disaster-recovery workflows copy snapshots across regions.

## Key Features

- Full azurerm v5 surface: both creation modes (Copy from a disk/snapshot, Import from a VHD blob), incremental mode, the network-access posture (policy, disk access, public access), and legacy ADE encryption settings.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup, `source_resource_id` to AzureManagedDisk, `storage_account_id` to AzureStorageAccount, and the encryption vault references to AzureKeyVault; the `snapshot_id` output is exactly what a gallery image version's source references.
- The provider's own looseness is preserved honestly: the schema does not tie `create_option` to its source fields (Azure validates the pairing at create), and the spec transcribes that -- with the working pairs documented on the fields.

## Use Cases

- **Golden image publishing**: snapshot the prepared OS disk, feed it to AzureComputeGalleryImage as a version source.
- **Pre-change backups**: an incremental snapshot before risky maintenance, deleted after the change settles.
- **Disk cloning**: create dev/test disks from a production snapshot.

## Future Enhancements

- Disk-access resources (private-endpoint export) are a P2 catalog kind; the `disk_access_id` reference gains its typed default when that kind lands.
