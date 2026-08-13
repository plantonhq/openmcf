# AzureManagedDisk

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureManagedDiskSpec** defines the configuration for creating an Azure
Managed Disk: the standalone block storage volume virtual machines
attach for data that must outlive any one VM.

The disk is a first-class resource in Azure's own model -- it has its
own lifecycle, SKU, encryption, and network posture, and a VM ATTACHES
it rather than containing it:
- an AzureVirtualMachine's data_disk_attachments references this disk's
  disk_id output (with a LUN and caching mode),
- a shared disk (max_shares) attaches to SEVERAL VMs at once -- the
  clustered-database seam that only a standalone disk can express,
- detaching and re-attaching to a replacement VM is the disk's whole
  point: the data outlives the machine.

The create_option is the disk's origin story -- an empty volume, a copy
of a snapshot or another disk, a platform or gallery image's disk, an
imported VHD, or a direct-upload target -- and the source fields it
requires are enforced by spec-level validation exactly as ARM enforces
them.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureManagedDisk
metadata:
  name: test-disk
  labels:
    environment: production
    team: platform
spec:
  region: eastus

  # The resource group the disk lives in (literal value here; a manifest
  # can also reference an AzureResourceGroup's name output via valueFrom).
  resourceGroup:
    value: test-rg

  name: orders-db-data

  # Premium SSD pinned to zone 1 with a bought-up performance tier.
  storageAccountType: PREMIUM_LRS
  createOption: EMPTY
  diskSizeGb: 512
  zone: "1"
  tier: P30

  # User tags merged over the metadata-derived tags.
  tags:
    cost-center: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.storageAccountType` | `enum` | yes |  |  |
| `spec.createOption` | `enum` | yes |  |  |
| `spec.diskSizeGb` | `int32` |  |  |  |
| `spec.sourceResourceId` | `string` |  |  |  |
| `spec.sourceUri` | `string` |  |  |  |
| `spec.storageAccountId` | `string` |  |  |  |
| `spec.imageReferenceId` | `string` |  |  |  |
| `spec.galleryImageReferenceId` | `string` |  |  |  |
| `spec.uploadSizeBytes` | `int64` |  |  |  |
| `spec.osType` | `enum` |  |  |  |
| `spec.hyperVGeneration` | `enum` |  |  |  |
| `spec.zone` | `string` |  |  |  |
| `spec.diskIopsReadWrite` | `int32` |  |  |  |
| `spec.diskMbpsReadWrite` | `int32` |  |  |  |
| `spec.diskIopsReadOnly` | `int32` |  |  |  |
| `spec.diskMbpsReadOnly` | `int32` |  |  |  |
| `spec.tier` | `string` |  |  |  |
| `spec.maxShares` | `int32` |  |  |  |
| `spec.onDemandBurstingEnabled` | `bool` |  |  |  |
| `spec.logicalSectorSize` | `int32` |  |  |  |
| `spec.diskEncryptionSetId` | `string \| valueFrom` |  |  | AzureDiskEncryptionSet (`status.outputs.disk_encryption_set_id`) |
| `spec.secureVmDiskEncryptionSetId` | `string \| valueFrom` |  |  | AzureDiskEncryptionSet (`status.outputs.disk_encryption_set_id`) |
| `spec.securityType` | `enum` |  |  |  |
| `spec.trustedLaunchEnabled` | `bool` |  |  |  |
| `spec.networkAccessPolicy` | `enum` |  |  |  |
| `spec.diskAccessId` | `string` |  |  |  |
| `spec.publicNetworkAccessEnabled` | `bool` |  | `true` |  |
| `spec.optimizedFrequentAttachEnabled` | `bool` |  |  |  |
| `spec.performancePlusEnabled` | `bool` |  |  |  |
| `spec.edgeZone` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region the disk lives in, e.g. "eastus". A disk can only
attach to VMs in its own region (and zone, when pinned). Changing the
region replaces the disk.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the disk will be created in. Can be a
literal resource-group name or a reference to an AzureResourceGroup's
name output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the disk, unique within the resource group. 1-80
characters (alphanumerics, underscores, periods, and hyphens).
Changing the name replaces the disk -- name it after the data it
carries ("orders-db-data"), not the VM it happens to attach to.

- rule: Managed disk names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.storageAccountType

`enum` · required

The disk's storage SKU -- the fundamental performance/redundancy
choice. STANDARD_LRS (HDD) for cold data; STANDARD_SSD_LRS/_ZRS for
light workloads; PREMIUM_LRS/_ZRS for production databases and OS
disks (fixed per-size performance tiers, credit-based bursting);
PREMIUM_V2_LRS and ULTRA_SSD_LRS decouple capacity from performance
-- size, IOPS, and throughput are dialed independently (the
disk_iops_*/disk_mbps_* fields). The ZRS variants replicate across
availability zones and cannot also be zone-pinned. Updatable in place
between compatible SKUs (the attached VM size must support the
target).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_managed_disk_storage_account_type_unspecified` -- Not specified -- invalid; the SKU is the disk's fundamental choice and must be explicit.
- `STANDARD_LRS` -- HDD, locally redundant: cold data, dev/test, infrequent access.
- `STANDARD_SSD_LRS` -- SSD, locally redundant: light production, web servers, small databases.
- `STANDARD_SSD_ZRS` -- SSD, zone redundant: light workloads that must survive a zone outage.
- `PREMIUM_LRS` -- Premium SSD, locally redundant: the production default for OS disks and databases (fixed per-size tiers, credit-based bursting).
- `PREMIUM_ZRS` -- Premium SSD, zone redundant: premium performance that survives a zone outage.
- `PREMIUM_V2_LRS` -- Premium SSD v2: capacity, IOPS, and throughput dialed independently (data disks only; no OS disks).
- `ULTRA_SSD_LRS` -- Ultra Disk: the highest performance envelope, dialed independently; zonal only, data disks only.

### spec.createOption

`enum` · required

How the disk is created -- its origin story. Fixed at creation:
EMPTY provisions a blank volume (set disk_size_gb); COPY clones an
existing disk or snapshot (set source_resource_id); FROM_IMAGE copies
a platform/marketplace image's disk (set image_reference_id) or a
Shared Image Gallery version (set gallery_image_reference_id);
IMPORT/IMPORT_SECURE wrap an existing VHD blob (set source_uri +
storage_account_id); RESTORE materializes a backup recovery point
(set source_resource_id); UPLOAD creates a direct-upload target (set
upload_size_bytes).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_managed_disk_create_option_unspecified` -- Not specified -- invalid; every disk has an origin.
- `EMPTY` -- A blank volume of disk_size_gb.
- `COPY` -- A clone of an existing managed disk or snapshot (source_resource_id).
- `FROM_IMAGE` -- A copy of a platform/marketplace image's disk (image_reference_id) or a Shared Image Gallery version (gallery_image_reference_id).
- `IMPORT` -- Wraps an existing VHD blob (source_uri + storage_account_id).
- `IMPORT_SECURE` -- Securely imports a VHD for confidential-VM scenarios (source_uri + storage_account_id; hyper_v_generation must be V2).
- `RESTORE` -- Materializes a backup recovery point (source_resource_id).
- `UPLOAD` -- A direct-upload target for streaming a VHD without a staging storage account (upload_size_bytes).

### spec.diskSizeGb

`int32` · optional (explicit presence)

The disk's size in GiB (1-65536; the upper reaches need
PREMIUM_V2_LRS/ULTRA_SSD_LRS). Required for EMPTY; for COPY and
FROM_IMAGE it may be omitted to inherit the source's size, or set
larger to grow at creation. Size can only ever INCREASE -- growing an
attached disk may briefly detach or deallocate the VM except where
Azure supports live resize.

- rule: {"int32":{"lte":65536,"gte":1}}

### spec.sourceResourceId

`string`

For COPY: the ARM ID of the managed disk or snapshot to clone. For
RESTORE: the recovery point to materialize. Fixed at creation.

### spec.sourceUri

`string`

For IMPORT/IMPORT_SECURE: the URI of the VHD blob to wrap. Fixed at
creation.

### spec.storageAccountId

`string`

For IMPORT/IMPORT_SECURE: the ARM ID of the storage account holding
source_uri. Fixed at creation.

### spec.imageReferenceId

`string`

For FROM_IMAGE: the ARM ID of the platform/marketplace image to copy.
Exactly one of image_reference_id / gallery_image_reference_id.
Fixed at creation.

### spec.galleryImageReferenceId

`string`

For FROM_IMAGE: the ARM ID of the Shared Image Gallery version to
copy. Exactly one of image_reference_id / gallery_image_reference_id.
Fixed at creation.

### spec.uploadSizeBytes

`int64`

For UPLOAD: the upload target's size in BYTES, which must equal the
source VHD's byte size exactly (footer included). Fixed at creation.

### spec.osType

`enum`

For disks carrying an operating system (IMPORT/IMPORT_SECURE/COPY of
an OS disk): which OS it holds. Leave unspecified for data disks.

Allowed values (use exactly as shown):

- `azure_managed_disk_os_type_unspecified` -- Not specified: a data disk.
- `LINUX` -- The disk carries Linux.
- `WINDOWS` -- The disk carries Windows.

### spec.hyperVGeneration

`enum`

The Hyper-V boot generation for OS-carrying disks (V2 = UEFI, the
modern default for trusted launch and confidential VMs; V1 = BIOS).
IMPORT_SECURE requires V2. Leave unspecified for data disks. Fixed at
creation.

Allowed values (use exactly as shown):

- `azure_managed_disk_hyper_v_generation_unspecified` -- Not specified: data disk, or the image's own generation.
- `V1` -- BIOS boot -- legacy images only.
- `V2` -- UEFI boot -- the modern generation; required for IMPORT_SECURE, trusted launch, and confidential VMs.

### spec.zone

`string`

The availability zone to pin the disk to ("1", "2", or "3"). A zonal
disk only attaches to VMs in the same zone. Leave unset for regional
(non-zonal) or ZRS disks -- a ZRS SKU must not also be zone-pinned.
Fixed at creation.

- rule: {"string":{"in":["","1","2","3"]}}

### spec.diskIopsReadWrite

`int32` · optional (explicit presence)

Provisioned read/write IOPS, for PREMIUM_V2_LRS and ULTRA_SSD_LRS
only (other SKUs have fixed per-size performance). One I/O operation
moves 4k-256k bytes. Updatable in place.

- rule: {"int32":{"gte":1}}

### spec.diskMbpsReadWrite

`int32` · optional (explicit presence)

Provisioned read/write throughput in MBps, for PREMIUM_V2_LRS and
ULTRA_SSD_LRS only. Updatable in place.

- rule: {"int32":{"gte":1}}

### spec.diskIopsReadOnly

`int32` · optional (explicit presence)

IOPS budget shared by VMs mounting a SHARED disk read-only (requires
max_shares; PREMIUM_V2_LRS/ULTRA_SSD_LRS only). Updatable in place.

- rule: {"int32":{"gte":1}}

### spec.diskMbpsReadOnly

`int32` · optional (explicit presence)

Throughput budget (MBps) shared by VMs mounting a SHARED disk
read-only (requires max_shares; PREMIUM_V2_LRS/ULTRA_SSD_LRS only).
Updatable in place.

- rule: {"int32":{"gte":1}}

### spec.tier

`string`

The performance tier for PREMIUM_LRS/PREMIUM_ZRS disks (e.g. "P30"),
decoupling performance from size -- a small disk can buy a bigger
tier's IOPS for bursty workloads, and pre-provisioned tiers avoid
resize-for-performance. Leave unset for the size's default tier.
Changing it on an attached disk briefly deallocates the VM.

### spec.maxShares

`int32` · optional (explicit presence)

The maximum number of VMs that can attach this disk simultaneously
(2-10) -- the shared-disk seam clustered workloads (failover
databases, scale-out file systems) build on. Leave unset for a
normal single-attach disk. The limit depends on SKU and size.

- rule: {"int32":{"lte":10,"gte":2}}

### spec.onDemandBurstingEnabled

`bool`

Whether on-demand bursting is enabled, for PREMIUM_LRS/PREMIUM_ZRS
disks larger than 512 GiB: the disk bursts beyond its provisioned
tier on demand (billed per burst) instead of relying on the
credit-based bursting smaller premium disks get automatically.

### spec.logicalSectorSize

`int32` · optional (explicit presence)

The disk's logical sector size in bytes (512 or 4096), for
PREMIUM_V2_LRS/ULTRA_SSD_LRS only. Azure's default is 4096; choose
512 only for legacy applications that require it. Fixed at creation.

- rule: {"int32":{"in":[512,4096]}}

### spec.diskEncryptionSetId

`string | valueFrom`

The customer-managed-key disk encryption set encrypting this disk.
A disk encryption set by ARM ID, or a reference to an
AzureDiskEncryptionSet's output. Conflicts with
secure_vm_disk_encryption_set_id.

- references: AzureDiskEncryptionSet (`status.outputs.disk_encryption_set_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDiskEncryptionSet, name: <that resource's name>, fieldPath: status.outputs.disk_encryption_set_id}} -- a bare string does not parse

### spec.secureVmDiskEncryptionSetId

`string | valueFrom`

For CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY security: the
disk encryption set encrypting the confidential VM's guest state.
A disk encryption set by ARM ID, or a reference to an
AzureDiskEncryptionSet's output. Conflicts with
disk_encryption_set_id. Fixed at creation.

- references: AzureDiskEncryptionSet (`status.outputs.disk_encryption_set_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDiskEncryptionSet, name: <that resource's name>, fieldPath: status.outputs.disk_encryption_set_id}} -- a bare string does not parse

### spec.securityType

`enum`

The confidential-VM security profile for OS disks of confidential
VMs. Leave unspecified for everything else. When set to
CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY, create_option must
be FROM_IMAGE or IMPORT_SECURE and secure_vm_disk_encryption_set_id
is required. Cannot be combined with trusted_launch_enabled. Fixed
at creation.

Allowed values (use exactly as shown):

- `azure_managed_disk_security_type_unspecified` -- Not specified: not a confidential-VM disk.
- `CONFIDENTIAL_VM_VMGUEST_STATE_ONLY_ENCRYPTED_WITH_PLATFORM_KEY` -- Only the VM guest state is encrypted, with a platform-managed key.
- `CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_PLATFORM_KEY` -- The full disk is encrypted with a platform-managed key.
- `CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY` -- The full disk is encrypted with a customer-managed key (set secure_vm_disk_encryption_set_id; create_option FROM_IMAGE or IMPORT_SECURE).

### spec.trustedLaunchEnabled

`bool`

Whether trusted launch (secure boot + vTPM support) is enabled for
this OS disk. Only valid when create_option is FROM_IMAGE or IMPORT
and the image supports it; cannot be combined with security_type.
Fixed at creation.

### spec.networkAccessPolicy

`enum`

The disk's network access posture. Unspecified applies Azure's
default (ALLOW_ALL: the disk's export endpoint is reachable over the
network with proper authorization). ALLOW_PRIVATE restricts export
to the disk-access resource's private endpoints (set disk_access_id);
DENY_ALL disables network export entirely -- the lockdown posture
for disks that never need SAS-based export.

Allowed values (use exactly as shown):

- `azure_managed_disk_network_access_policy_unspecified` -- Not specified: Azure's default (AllowAll).
- `ALLOW_ALL` -- The disk's export endpoint is reachable with proper authorization.
- `ALLOW_PRIVATE` -- Export only through the disk-access resource's private endpoints (set disk_access_id).
- `DENY_ALL` -- Network export disabled entirely -- the lockdown posture.

### spec.diskAccessId

`string`

For ALLOW_PRIVATE: the disk-access resource whose private endpoints
export traffic uses, by ARM ID. Plain ARM ID: disk-access resources
are not modeled as a Planton kind.
Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Compute/diskAccesses/{name}

### spec.publicNetworkAccessEnabled

`bool` · optional (explicit presence)

Whether the disk's export endpoint is reachable over the public
network at all. Azure's default is true; false pairs with
ALLOW_PRIVATE for a fully private posture. Updatable in place.

- default: `true`

### spec.optimizedFrequentAttachEnabled

`bool`

Whether the disk skips fault-domain alignment with its VM to
optimize for very frequent attach/detach cycles (more than 5 a day).
Azure's default is false; leaving alignment on is right for
virtually all disks.

### spec.performancePlusEnabled

`bool`

Whether performance-plus is enabled, raising the baseline IOPS/
throughput of an eligible disk (512 GiB+, deployed from a supported
create option). Fixed at creation.

### spec.edgeZone

`string`

The Azure Edge Zone the disk is deployed in, for edge-computing
workloads. Leave unset for regular regional deployment. Fixed at
creation.

### spec.tags

`map<string, string>`

Free-form tags applied to the disk, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Tags are Azure's governance surface -- Azure
Policy enforces them and Microsoft Cost Management groups by them.
Updatable in place.

## Validation Rules

- `disk_empty_requires_size`: create_option EMPTY requires disk_size_gb
- `disk_copy_restore_requires_source_resource`: create_option COPY and RESTORE require source_resource_id (the disk/snapshot to clone or the recovery point to restore)
- `disk_import_requires_source_uri_and_account`: create_option IMPORT and IMPORT_SECURE require source_uri and storage_account_id
- `disk_from_image_requires_exactly_one_image`: create_option FROM_IMAGE requires exactly one of image_reference_id (platform/marketplace) or gallery_image_reference_id (Shared Image Gallery)
- `disk_image_sources_only_for_from_image`: image_reference_id and gallery_image_reference_id are only valid with create_option FROM_IMAGE
- `disk_upload_requires_byte_size`: create_option UPLOAD requires upload_size_bytes (the exact byte size of the VHD to upload, footer included)
- `disk_performance_dials_need_flexible_sku`: disk_iops_*/disk_mbps_*/logical_sector_size are only settable on PREMIUM_V2_LRS and ULTRA_SSD_LRS (other SKUs have fixed per-size performance)
- `disk_read_only_dials_need_shared_disk`: disk_iops_read_only/disk_mbps_read_only require max_shares (they budget the read-only mounts of a shared disk)
- `disk_bursting_and_tier_need_premium`: tier and on_demand_bursting_enabled apply to PREMIUM_LRS/PREMIUM_ZRS disks only
- `disk_encryption_sets_conflict`: disk_encryption_set_id and secure_vm_disk_encryption_set_id are mutually exclusive
- `disk_secure_vm_set_pairs_with_cvm_customer_key`: secure_vm_disk_encryption_set_id is required when security_type is CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY, and only valid then
- `disk_cvm_customer_key_create_options`: security_type CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY requires create_option FROM_IMAGE or IMPORT_SECURE
- `disk_trusted_launch_conflicts_security_type`: trusted_launch_enabled cannot be combined with security_type (a disk is trusted-launch or confidential, not both)
- `disk_trusted_launch_create_options`: trusted_launch_enabled requires create_option FROM_IMAGE or IMPORT
- `disk_access_pairs_with_allow_private`: disk_access_id is required when network_access_policy is ALLOW_PRIVATE, and only valid then
- `disk_zrs_cannot_be_zonal`: a ZRS disk replicates across zones and cannot also be pinned to one (leave zone unset)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureManagedDisk, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.disk_id` | `string` | The Azure Resource Manager ID of the disk. This is the primary output: AzureVirtualMachine's data_disk_attachments references it to attach the disk. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Compute/disks/{name} |
| `status.outputs.disk_name` | `string` | The name of the disk. |
| `status.outputs.disk_size_gb` | `int32` | The disk's ACTUAL size in GiB -- inherited from the source for COPY/FROM_IMAGE disks that omitted disk_size_gb, so downstream capacity planning reads the real value. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.diskEncryptionSetId` | AzureDiskEncryptionSet | `status.outputs.disk_encryption_set_id` |
| `spec.secureVmDiskEncryptionSetId` | AzureDiskEncryptionSet | `status.outputs.disk_encryption_set_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureDataProtectionBackupInstance | `spec.disk.diskId` | `status.outputs.disk_id` |
| AzureDiskSnapshot | `spec.sourceResourceId` | `status.outputs.disk_id` |
| AzureVirtualMachine | `spec.osManagedDiskId` | `status.outputs.disk_id` |
| AzureVirtualMachine | `spec.dataDiskAttachments[].managedDiskId` | `status.outputs.disk_id` |

## See Also

- [Overview](../README.md)
