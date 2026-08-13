# Overview

The **AzureComputeGalleryImage** component deploys a gallery image: one image definition inside an Azure Compute Gallery -- its marketplace-style identity, OS type, security posture, and recommended sizing -- plus the published, region-replicated versions VMs actually deploy.

## Purpose

- **The image's contract and its releases in one place**: the definition declares what the image IS (publisher/offer/SKU, Linux/Windows, Gen1/Gen2, trusted-launch posture); the versions list publishes the actual release artifacts, each built from exactly one source and replicated to chosen regions.
- **Versions are the image owner's own release history**: add an entry to publish, remove it to unpublish -- the manifest is the authoritative list of what is deployable.
- **Consumers pin or float**: VMs deploy from a version's ARM ID (exact release) or the definition's ID (latest non-excluded version).

## Key Features

- Full azurerm v5 surface across BOTH provider resources (`azurerm_shared_image` + the composed `azurerm_shared_image_version`): the four-flag security clique, recommended sizing ranges, disk-type exclusions, purchase plans, per-region replica counts, storage types, and customer-managed-key encryption.
- Chart-ready: `gallery_name` defaults its reference to AzureComputeGallery, a version's `os_disk_snapshot_id` to AzureDiskSnapshot, `storage_account_id` to AzureStorageAccount, and per-region `disk_encryption_set_id` to AzureDiskEncryptionSet; outputs carry the definition's id/name plus a per-version id map.
- The provider's contracts are front-loaded as spec validation rules: at most one security flag, exactly one version source, blob-and-storage-account together, no encryption sets on Shallow replication, ordered sizing ranges, unique version names.

## Use Cases

- **Golden image publishing**: the image build pipeline snapshots a prepared OS disk and publishes it as the next version; consumers on "latest" pick it up on their next deploy.
- **Multi-region image distribution**: one version replicated to every workload region, with per-region replica counts sized to deployment concurrency.
- **Compliance-pinned fleets**: production VMs pin exact version IDs from the `version_ids` output; dev floats on latest.

## Future Enhancements

- Legacy managed images (`azurerm_image`) are deliberately excluded -- superseded by the gallery model this kind implements.
