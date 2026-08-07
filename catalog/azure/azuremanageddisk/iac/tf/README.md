# AzureManagedDisk Terraform Module

## Overview

This Terraform module provisions an Azure Managed Disk using the `azurerm`
provider (`~> 5.0`). It creates a single `azurerm_managed_disk` -- the
standalone block storage volume whose data outlives any one virtual
machine -- with the configured SKU, origin (create option), size,
performance dials, encryption, and network posture.

The VM-side attachment is deliberately NOT here: AzureVirtualMachine's
`data_disk_attachments` owns which VM mounts this disk, at which LUN,
with which caching -- so the disk survives VM replacement untouched, and
a shared disk (`max_shares`) can attach to several VMs at once.

Spec-level validation enforces the same rules ARM does -- each create
option's required source fields, the SKU gates on the performance dials,
the encryption pairings, and the network-posture pairing -- so the module
maps fields without re-validating them.

Lifecycle notes worth knowing before operating this resource:

- Name, region, zone, create option (and its source fields), logical
  sector size, security profile, and performance-plus are the disk's
  identity -- changing any of them replaces the disk **and its data**.
- `disk_size_gb` can only INCREASE. Growing an attached disk may briefly
  detach it or deallocate the VM, except where Azure supports live
  resize (and crossing 4 TiB on non-PremiumV2/Ultra SKUs always
  detaches).
- Changing `tier` or the SKU on an attached disk deallocates the VM for
  the change and restarts it after.
- The performance dials, network posture, and tags update in place.

## Resources Created

- `azurerm_managed_disk.main` -- the disk itself

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Managed disk specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region; a disk only attaches to VMs in its own region (and zone) |
| `resource_group` | yes | Resource group name |
| `name` | yes | Disk name, unique within the resource group; renaming replaces the disk |
| `storage_account_type` | yes | The SKU, as the spec enum's name string (`STANDARD_LRS` ... `ULTRA_SSD_LRS`) |
| `create_option` | yes | The disk's origin (`EMPTY` / `COPY` / `FROM_IMAGE` / `IMPORT` / `IMPORT_SECURE` / `RESTORE` / `UPLOAD`); fixed at creation |
| `disk_size_gb` | for `EMPTY` | Size in GiB; `COPY`/`FROM_IMAGE` inherit the source's size when unset. Can only increase |
| `source_resource_id` | per option | `COPY`: the disk/snapshot to clone; `RESTORE`: the recovery point |
| `source_uri` + `storage_account_id` | per option | `IMPORT`/`IMPORT_SECURE`: the VHD blob and its storage account |
| `image_reference_id` / `gallery_image_reference_id` | per option | `FROM_IMAGE`: exactly one of the platform image or gallery version |
| `upload_size_bytes` | per option | `UPLOAD`: the exact byte size of the VHD (footer included) |
| `os_type` / `hyper_v_generation` | no | OS-carrying disks only (`LINUX`/`WINDOWS`, `V1`/`V2`); unset for data disks |
| `zone` | no | Availability zone to pin a zonal disk to; unset for regional or ZRS disks |
| `disk_iops_read_write` / `disk_mbps_read_write` | no | Independent performance dials (`PREMIUM_V2_LRS`/`ULTRA_SSD_LRS` only) |
| `disk_iops_read_only` / `disk_mbps_read_only` | no | Read-only budgets for a shared disk's mounts (require `max_shares`) |
| `tier` | no | Premium SSD performance tier (e.g. `P30`); unset uses the size's default |
| `max_shares` | no | Shared-disk attach limit (2-10); unset for single-attach |
| `on_demand_bursting_enabled` | no | On-demand bursting for `PREMIUM_LRS`/`PREMIUM_ZRS` disks > 512 GiB |
| `logical_sector_size` | no | 512 or 4096 (`PREMIUM_V2_LRS`/`ULTRA_SSD_LRS` only); fixed at creation |
| `disk_encryption_set_id` | no | Customer-managed-key encryption set, as a resolved ARM ID |
| `secure_vm_disk_encryption_set_id` | no | Confidential-VM customer-key encryption set (pairs with the matching `security_type`) |
| `security_type` | no | Confidential-VM security profile, as the spec enum's name string |
| `trusted_launch_enabled` | no | Trusted launch (`FROM_IMAGE`/`IMPORT` only; conflicts with `security_type`) |
| `network_access_policy` | no | Export posture (`ALLOW_ALL` / `ALLOW_PRIVATE` / `DENY_ALL`); unset applies Azure's default |
| `disk_access_id` | no | The disk-access resource for `ALLOW_PRIVATE`, as an ARM ID |
| `public_network_access_enabled` | no | Whether the export endpoint is publicly reachable (Azure defaults to true) |
| `optimized_frequent_attach_enabled` | no | Skip fault-domain alignment for very frequent attach/detach cycles |
| `performance_plus_enabled` | no | Raise the baseline performance of an eligible 512 GiB+ disk; fixed at creation |
| `edge_zone` | no | Azure Edge Zone pinning; fixed at creation |
| `tags` | no | User tags, merged over metadata-derived tags (user wins) |

The module maps the spec enums' name strings to ARM values internally
(e.g. `PREMIUM_V2_LRS` becomes `PremiumV2_LRS`, `FROM_IMAGE` becomes
`FromImage`); only explicit choices are ever sent, so an unspecified
optional enum and Azure's default deploy identically.

## Outputs

| Output | Description |
|--------|-------------|
| `disk_id` | Full ARM ID of the disk -- the join key AzureVirtualMachine's `data_disk_attachments` references |
| `disk_name` | The disk's name as deployed |
| `disk_size_gb` | The disk's ACTUAL size in GiB (inherited from the source when the spec omitted it) |

## Usage

```hcl
module "managed_disk" {
  source = "./iac/tf"

  metadata = { name = "orders-db-data", org = "mycompany", env = "production" }

  spec = {
    region               = "eastus"
    resource_group       = "prod-rg"
    name                 = "orders-db-data"
    storage_account_type = "PREMIUM_LRS"
    create_option        = "EMPTY"
    disk_size_gb         = 512
    zone                 = "1"
  }
}
```

## Required Permissions

The deploying credential needs `Microsoft.Compute/disks/write` on the
resource group -- held via Contributor or Owner. Customer-managed-key
encryption additionally requires read access on the referenced disk
encryption set (`Microsoft.Compute/diskEncryptionSets/read`).
