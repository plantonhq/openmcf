# Create the managed disk -- the standalone block volume whose data
# outlives any one virtual machine.
#
# Lifecycle notes worth knowing before operating this resource:
# - Name, region, zone, create option (and its source fields), logical
#   sector size, security profile, and performance-plus are the disk's
#   identity -- changing any of them replaces the disk AND ITS DATA.
# - disk_size_gb can only INCREASE. Growing an attached disk may briefly
#   detach it or deallocate the VM, except where Azure supports live
#   resize (and crossing 4 TiB on non-PremiumV2/Ultra SKUs always
#   detaches).
# - Changing tier or the SKU on an attached disk deallocates the VM for
#   the change and restarts it after.
# - The VM-side attachment is NOT here: AzureVirtualMachine's
#   data_disk_attachments owns which VM mounts this disk, at which LUN,
#   with which caching -- so the disk survives VM replacement untouched.
resource "azurerm_managed_disk" "main" {
  name                 = var.spec.name
  location             = var.spec.region
  resource_group_name  = var.spec.resource_group
  storage_account_type = local.storage_account_type
  create_option        = local.create_option

  disk_size_gb = var.spec.disk_size_gb

  # The create option's source fields; spec-level validation enforces the
  # same option-to-source pairings ARM does, so nulls here are safe.
  source_resource_id         = var.spec.source_resource_id
  source_uri                 = var.spec.source_uri
  storage_account_id         = var.spec.storage_account_id
  image_reference_id         = var.spec.image_reference_id
  gallery_image_reference_id = var.spec.gallery_image_reference_id
  upload_size_bytes          = var.spec.upload_size_bytes

  os_type            = local.os_type
  hyper_v_generation = var.spec.hyper_v_generation == "V1" ? "V1" : var.spec.hyper_v_generation == "V2" ? "V2" : null

  zone = var.spec.zone

  # Independent performance dials -- PremiumV2/Ultra only (the read-only
  # pair budgets a shared disk's read-only mounts).
  disk_iops_read_write = var.spec.disk_iops_read_write
  disk_mbps_read_write = var.spec.disk_mbps_read_write
  disk_iops_read_only  = var.spec.disk_iops_read_only
  disk_mbps_read_only  = var.spec.disk_mbps_read_only

  # Premium SSD tier decoupling (null = the size's default tier) and
  # on-demand bursting for >512 GiB premium disks.
  tier                       = var.spec.tier
  on_demand_bursting_enabled = var.spec.on_demand_bursting_enabled

  # The shared-disk seam: >1 lets several VMs attach simultaneously.
  max_shares = var.spec.max_shares

  logical_sector_size = var.spec.logical_sector_size

  # Encryption: customer-managed keys via a disk encryption set, or the
  # confidential-VM customer-key variant (mutually exclusive; spec-level
  # validation enforces the pairing with security_type).
  disk_encryption_set_id           = var.spec.disk_encryption_set_id
  secure_vm_disk_encryption_set_id = var.spec.secure_vm_disk_encryption_set_id
  security_type                    = local.security_type
  trusted_launch_enabled           = var.spec.trusted_launch_enabled

  # Network export posture: who can reach the disk's export endpoint.
  network_access_policy         = local.network_access_policy
  disk_access_id                = var.spec.disk_access_id
  public_network_access_enabled = var.spec.public_network_access_enabled

  optimized_frequent_attach_enabled = var.spec.optimized_frequent_attach_enabled
  performance_plus_enabled          = var.spec.performance_plus_enabled

  edge_zone = var.spec.edge_zone

  tags = local.final_tags
}
