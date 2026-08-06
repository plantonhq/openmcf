variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Managed Disk specification"
  type = object({
    # The Azure region the disk lives in (a disk only attaches to VMs in
    # its own region and zone).
    region = string

    # The resource group the disk lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The disk's name, unique within the resource group. Renaming
    # replaces the disk.
    name = string

    # The storage SKU, as the spec enum's name string (STANDARD_LRS /
    # STANDARD_SSD_LRS / STANDARD_SSD_ZRS / PREMIUM_LRS / PREMIUM_ZRS /
    # PREMIUM_V2_LRS / ULTRA_SSD_LRS).
    storage_account_type = string

    # The disk's origin, as the spec enum's name string (EMPTY / COPY /
    # FROM_IMAGE / IMPORT / IMPORT_SECURE / RESTORE / UPLOAD). Fixed at
    # creation; spec-level validation enforces each option's required
    # source fields.
    create_option = string

    # The size in GiB (required for EMPTY; COPY/FROM_IMAGE inherit the
    # source's size when unset). Can only increase.
    disk_size_gb = optional(number)

    # COPY: the disk/snapshot to clone; RESTORE: the recovery point.
    source_resource_id = optional(string)

    # IMPORT/IMPORT_SECURE: the VHD blob and its storage account.
    source_uri         = optional(string)
    storage_account_id = optional(string)

    # FROM_IMAGE: exactly one of the platform image or the Shared Image
    # Gallery version.
    image_reference_id         = optional(string)
    gallery_image_reference_id = optional(string)

    # UPLOAD: the exact byte size of the VHD to upload (footer included).
    upload_size_bytes = optional(number)

    # OS-carrying disks: which OS (LINUX/WINDOWS) and boot generation
    # (V1/V2), as the spec enums' name strings. Unset for data disks.
    os_type            = optional(string)
    hyper_v_generation = optional(string)

    # The availability zone to pin a zonal disk to (unset for regional or
    # ZRS disks).
    zone = optional(string)

    # Independent performance dials (PREMIUM_V2_LRS/ULTRA_SSD_LRS only;
    # the read-only pair budgets a shared disk's read-only mounts).
    disk_iops_read_write = optional(number)
    disk_mbps_read_write = optional(number)
    disk_iops_read_only  = optional(number)
    disk_mbps_read_only  = optional(number)

    # Premium SSD performance tier (e.g. "P30"); unset uses the size's
    # default tier.
    tier = optional(string)

    # Shared-disk attach limit (2-10); unset for single-attach disks.
    max_shares = optional(number)

    # On-demand bursting for PREMIUM_LRS/PREMIUM_ZRS disks > 512 GiB.
    on_demand_bursting_enabled = optional(bool, false)

    # Logical sector size (512/4096) for PREMIUM_V2_LRS/ULTRA_SSD_LRS.
    logical_sector_size = optional(number)

    # Customer-managed-key encryption: the disk encryption set, as a
    # resolved ARM ID (conflicts with the secure-VM variant).
    disk_encryption_set_id = optional(string)

    # Confidential-VM customer-key encryption set, as a resolved ARM ID
    # (requires security_type CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY).
    secure_vm_disk_encryption_set_id = optional(string)

    # Confidential-VM security profile, as the spec enum's name string.
    security_type = optional(string)

    # Trusted launch (FROM_IMAGE/IMPORT only; conflicts with
    # security_type).
    trusted_launch_enabled = optional(bool, false)

    # Network export posture, as the spec enum's name string (ALLOW_ALL /
    # ALLOW_PRIVATE / DENY_ALL); unset applies Azure's default (AllowAll).
    network_access_policy = optional(string)

    # ALLOW_PRIVATE: the disk-access resource whose private endpoints
    # export traffic uses.
    disk_access_id = optional(string)

    # Whether the export endpoint is publicly reachable (Azure defaults
    # to true).
    public_network_access_enabled = optional(bool, true)

    # Skip fault-domain alignment for very frequent attach/detach cycles.
    optimized_frequent_attach_enabled = optional(bool, false)

    # Raise the baseline performance of an eligible 512 GiB+ disk (fixed
    # at creation).
    performance_plus_enabled = optional(bool, false)

    # Edge Zone pinning for edge-computing workloads (fixed at creation).
    edge_zone = optional(string)

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
