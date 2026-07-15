# Azure Managed Disk

Deploys an Azure Managed Disk — the standalone block storage volume virtual machines attach for data that must outlive any one VM. The disk is a first-class resource in Azure's own model: it has its own lifecycle, SKU, encryption, and network posture, and a VM *attaches* it rather than containing it. The VM-side attachment (which VM, at which LUN, with which caching) lives on `AzureVirtualMachine`'s `dataDiskAttachments`, referencing this disk's `disk_id` output — the disk spec deliberately knows nothing about its consumers.

## What Gets Created

When you deploy an AzureManagedDisk resource, Planton provisions:

- **Managed Disk** — an `azurerm_managed_disk` / `compute.ManagedDisk` in the specified region and resource group, with the configured SKU, origin (create option), size, performance dials, encryption, and network posture
- **Azure Tags** — resource metadata tags applied to the disk for tracking and governance, merged with any user-supplied tags (user tags win on key collision)

Nothing else is created here. The disk does not attach itself to any VM (each `AzureVirtualMachine` declares its attachments via `dataDiskAttachments`), and it does not create disk encryption sets or disk-access resources — those are referenced by ARM ID.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the disk will be created (can reference an AzureResourceGroup resource)
- **For customer-managed-key encryption**: an existing disk encryption set, referenced by ARM ID
- **For a private network posture** (`ALLOW_PRIVATE`): an existing disk-access resource, referenced by ARM ID

## Quick Start

Create a file `disk.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedDisk
metadata:
  name: orders-db-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureManagedDisk.orders-db-data
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: orders-db-data
  storageAccountType: PREMIUM_LRS
  createOption: EMPTY
  diskSizeGb: 256
  zone: "1"
```

Deploy:

```shell
planton apply -f disk.yaml
```

This creates an empty 256 GiB Premium SSD pinned to zone 1. To mount it, reference the disk's `status.outputs.disk_id` from an `AzureVirtualMachine`'s `dataDiskAttachments` (the VM must be in the same region and zone).

## The Create Option Matrix

`createOption` is the disk's origin story, fixed at creation. Each option requires specific source fields, enforced at spec level exactly as ARM enforces them:

| Create Option | Required Source Fields | What It Does |
|---------------|------------------------|--------------|
| `EMPTY` | `diskSizeGb` | A blank volume of the given size |
| `COPY` | `sourceResourceId` | Clones an existing managed disk or snapshot |
| `FROM_IMAGE` | exactly one of `imageReferenceId` / `galleryImageReferenceId` | Copies a platform/marketplace image's disk or a Shared Image Gallery version |
| `IMPORT` | `sourceUri` + `storageAccountId` | Wraps an existing VHD blob |
| `IMPORT_SECURE` | `sourceUri` + `storageAccountId` (and `hyperVGeneration: V2`) | Securely imports a VHD for confidential-VM scenarios |
| `RESTORE` | `sourceResourceId` | Materializes a backup recovery point |
| `UPLOAD` | `uploadSizeBytes` | A direct-upload target for streaming a VHD without a staging storage account |

For `COPY` and `FROM_IMAGE`, `diskSizeGb` may be omitted to inherit the source's size (the actual value surfaces in the `disk_size_gb` output) or set larger to grow at creation.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the disk (e.g., `eastus`). A disk only attaches to VMs in its own region (and zone, when pinned). | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. | Required |
| `name` | `string` | Name of the disk, unique within the resource group. Name it after the data it carries (`orders-db-data`), not the VM it happens to attach to. | Required, 1-80 chars, Azure naming rules |
| `storageAccountType` | `enum` | The storage SKU: `STANDARD_LRS`, `STANDARD_SSD_LRS`, `STANDARD_SSD_ZRS`, `PREMIUM_LRS`, `PREMIUM_ZRS`, `PREMIUM_V2_LRS`, `ULTRA_SSD_LRS`. Updatable in place between compatible SKUs. | Required |
| `createOption` | `enum` | The disk's origin (see the matrix above). Fixed at creation. | Required, with per-option source-field pairings |

### Sizing and Placement

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `diskSizeGb` | `int32` | source's size | Size in GiB, 1–65536 (the upper reaches need `PREMIUM_V2_LRS`/`ULTRA_SSD_LRS`). Required for `EMPTY`. **Size can only ever increase** — growing an attached disk may briefly detach it or deallocate the VM except where Azure supports live resize. |
| `zone` | `string` | none (regional) | Availability zone to pin the disk to (`"1"`, `"2"`, `"3"`). A zonal disk only attaches to VMs in the same zone. Must be unset for ZRS SKUs — a ZRS disk replicates across zones and cannot also be pinned to one. Fixed at creation. |
| `edgeZone` | `string` | none | Azure Edge Zone for edge-computing workloads. Fixed at creation. |

### Performance

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `diskIopsReadWrite` | `int32` | SKU baseline | Provisioned read/write IOPS. `PREMIUM_V2_LRS`/`ULTRA_SSD_LRS` only (other SKUs have fixed per-size performance). Updatable in place. |
| `diskMbpsReadWrite` | `int32` | SKU baseline | Provisioned read/write throughput in MBps. `PREMIUM_V2_LRS`/`ULTRA_SSD_LRS` only. Updatable in place. |
| `diskIopsReadOnly` | `int32` | none | IOPS budget shared by VMs mounting a **shared** disk read-only. Requires `maxShares`; `PREMIUM_V2_LRS`/`ULTRA_SSD_LRS` only. Updatable in place. |
| `diskMbpsReadOnly` | `int32` | none | Throughput budget (MBps) for a shared disk's read-only mounts. Requires `maxShares`; `PREMIUM_V2_LRS`/`ULTRA_SSD_LRS` only. Updatable in place. |
| `tier` | `string` | size's default tier | Performance tier for `PREMIUM_LRS`/`PREMIUM_ZRS` disks (e.g., `P30`) — a small disk can buy a bigger tier's IOPS for bursty workloads. Changing it on an attached disk briefly deallocates the VM. |
| `onDemandBurstingEnabled` | `bool` | `false` | On-demand bursting for `PREMIUM_LRS`/`PREMIUM_ZRS` disks larger than 512 GiB — bursts beyond the provisioned tier on demand (billed per burst). |
| `maxShares` | `int32` | none (single-attach) | Maximum simultaneous VM attachments (2–10) — the shared-disk seam clustered workloads build on. The limit depends on SKU and size. |
| `logicalSectorSize` | `int32` | `4096` | Logical sector size in bytes (512 or 4096), `PREMIUM_V2_LRS`/`ULTRA_SSD_LRS` only. Choose 512 only for legacy applications. Fixed at creation. |
| `performancePlusEnabled` | `bool` | `false` | Raises the baseline IOPS/throughput of an eligible disk (512 GiB+, supported create options). Fixed at creation. |
| `optimizedFrequentAttachEnabled` | `bool` | `false` | Skips fault-domain alignment with the VM to optimize for very frequent attach/detach cycles (more than 5 a day). Leaving alignment on is right for virtually all disks. |

### OS-Carrying Disks

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `osType` | `enum` | none (data disk) | `LINUX` or `WINDOWS`, for disks carrying an operating system (imports/copies of an OS disk). |
| `hyperVGeneration` | `enum` | image's default | `V2` (UEFI — required for `IMPORT_SECURE`, trusted launch, and confidential VMs) or `V1` (BIOS, legacy images). Fixed at creation. |

### Security and Encryption

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `diskEncryptionSetId` | `StringValueOrRef` | platform keys | The customer-managed-key disk encryption set, by ARM ID. Conflicts with `secureVmDiskEncryptionSetId`. |
| `secureVmDiskEncryptionSetId` | `StringValueOrRef` | none | The disk encryption set for `CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY` security — required with that security type and only valid then. Fixed at creation. |
| `securityType` | `enum` | none | Confidential-VM security profile for OS disks: `CONFIDENTIAL_VM_VMGUEST_STATE_ONLY_ENCRYPTED_WITH_PLATFORM_KEY`, `CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_PLATFORM_KEY`, or `CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY` (requires `createOption` `FROM_IMAGE` or `IMPORT_SECURE`). Cannot be combined with `trustedLaunchEnabled`. Fixed at creation. |
| `trustedLaunchEnabled` | `bool` | `false` | Trusted launch (secure boot + vTPM) for OS disks. Requires `createOption` `FROM_IMAGE` or `IMPORT`; cannot be combined with `securityType`. Fixed at creation. |

### Network Export Posture

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `networkAccessPolicy` | `enum` | `ALLOW_ALL` | Who can reach the disk's export endpoint: `ALLOW_ALL` (reachable with proper authorization), `ALLOW_PRIVATE` (only through a disk-access resource's private endpoints — set `diskAccessId`), or `DENY_ALL` (network export disabled entirely — the lockdown posture). |
| `diskAccessId` | `string` | none | The disk-access resource whose private endpoints export traffic uses, by ARM ID. Required with `ALLOW_PRIVATE`, and only valid then. |
| `publicNetworkAccessEnabled` | `bool` | `true` | Whether the export endpoint is reachable over the public network at all. `false` pairs with `ALLOW_PRIVATE` for a fully private posture. Updatable in place. |

### Tags

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `tags` | `map<string, string>` | `{}` | Free-form tags, merged over Planton-derived resource tags (user wins on collision). |

## Examples

### Premium Data Disk with VM Attachment

The production default: an empty zonal Premium SSD, attached by the VM that mounts it. The disk knows nothing about the VM — replacing the VM never touches the data:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedDisk
metadata:
  name: orders-db-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureManagedDisk.orders-db-data
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: orders-db-data
  storageAccountType: PREMIUM_LRS
  createOption: EMPTY
  diskSizeGb: 512
  zone: "1"
---
apiVersion: azure.planton.dev/v1
kind: AzureVirtualMachine
metadata:
  name: orders-db
spec:
  # ... VM configuration (same region and zone as the disk) ...
  dataDiskAttachments:
    - managedDiskId:
        valueFrom:
          name: orders-db-data
      lun: 0
      caching: READ_ONLY
```

### Premium SSD v2 with Dialed Performance

Capacity, IOPS, and throughput provisioned independently — a small disk with big performance, impossible on the classic per-size tiers:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedDisk
metadata:
  name: ledger-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureManagedDisk.ledger-data
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: ledger-data
  storageAccountType: PREMIUM_V2_LRS
  createOption: EMPTY
  diskSizeGb: 128
  zone: "1"
  diskIopsReadWrite: 8000
  diskMbpsReadWrite: 300
```

### Shared Disk for a Failover Cluster

One disk, several VMs: `maxShares` enables the clustered-database seam, and the read-only dials budget the standby nodes' mounts:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedDisk
metadata:
  name: cluster-quorum
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureManagedDisk.cluster-quorum
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: cluster-quorum
  storageAccountType: PREMIUM_V2_LRS
  createOption: EMPTY
  diskSizeGb: 256
  zone: "1"
  maxShares: 3
  diskIopsReadWrite: 6000
  diskMbpsReadWrite: 250
  diskIopsReadOnly: 3000
  diskMbpsReadOnly: 125
```

### Snapshot Clone

A full, independent copy of a snapshot — the restore and environment-duplication workhorse. Size is inherited from the source; the SKU is free to differ:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedDisk
metadata:
  name: orders-db-restored
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: staging.AzureManagedDisk.orders-db-restored
spec:
  region: eastus
  resourceGroup:
    value: staging-rg
  name: orders-db-restored
  storageAccountType: PREMIUM_LRS
  createOption: COPY
  sourceResourceId: /subscriptions/xxx/resourceGroups/prod-rg/providers/Microsoft.Compute/snapshots/orders-db-nightly
  zone: "1"
```

### Customer-Managed Keys with a Locked-Down Network Posture

Encryption through a disk encryption set and network export disabled entirely — the posture for regulated data that never needs SAS-based export:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedDisk
metadata:
  name: pii-vault-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureManagedDisk.pii-vault-data
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: pii-vault-data
  storageAccountType: PREMIUM_LRS
  createOption: EMPTY
  diskSizeGb: 1024
  zone: "2"
  diskEncryptionSetId:
    value: /subscriptions/xxx/resourceGroups/security-rg/providers/Microsoft.Compute/diskEncryptionSets/prod-cmk
  networkAccessPolicy: DENY_ALL
  publicNetworkAccessEnabled: false
  tags:
    data-classification: restricted
```

### Zone-Redundant Data

A ZRS SKU replicates the disk across availability zones — no `zone` field, because a ZRS disk cannot also be pinned to one:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedDisk
metadata:
  name: shared-content
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureManagedDisk.shared-content
spec:
  region: westeurope
  resourceGroup:
    value: prod-rg
  name: shared-content
  storageAccountType: PREMIUM_ZRS
  createOption: EMPTY
  diskSizeGb: 512
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `disk_id` | `string` | Azure Resource Manager ID of the disk — the primary output; `AzureVirtualMachine`'s `dataDiskAttachments` references it to attach the disk |
| `disk_name` | `string` | The disk's name as deployed |
| `disk_size_gb` | `int32` | The disk's **actual** size in GiB — inherited from the source for `COPY`/`FROM_IMAGE` disks that omitted `diskSizeGb`, so downstream capacity planning reads the real value |

## Operational Notes

- **Identity fields replace the disk and its data**: name, region, zone, create option (and its source fields), logical sector size, security profile, and performance-plus are fixed at creation — changing any of them replaces the disk
- **Size only grows**: `diskSizeGb` can never decrease; growing an attached disk may briefly detach it or deallocate the VM except where Azure supports live resize
- **SKU and tier changes on an attached disk** deallocate the VM for the change and restart it after

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) -- provides the resource group for disk placement
- [AzureVirtualMachine](/docs/catalog/azure/azurevirtualmachine) -- attaches the disk via `dataDiskAttachments`, owning the LUN and caching mode for each mount
