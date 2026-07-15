# ---------------------------------------------------------------------------
# AWS FSx for OpenZFS file system
# ---------------------------------------------------------------------------
# One aws_fsx_openzfs_file_system resource carries the whole spec. ForceNew
# attributes (replace the file system when changed): deployment_type,
# storage_type, subnet_ids, security_group_ids, kms_key_id, backup_id,
# preferred_subnet_id, endpoint_ip_address_range, and — the subtle one —
# root_volume_configuration.copy_tags_to_snapshots. Storage capacity grows in
# place (never shrinks); throughput_capacity scales up or down in place.
# ---------------------------------------------------------------------------

resource "aws_fsx_openzfs_file_system" "this" {
  # Core shape. deployment_type is always sent (spec default SINGLE_AZ_2);
  # storage_capacity is presence-honest — null means the capacity comes from
  # a backup restore or the elastic INTELLIGENT_TIERING class and the
  # argument must be omitted entirely.
  deployment_type     = var.spec.deployment_type
  storage_capacity    = var.spec.storage_capacity_gib
  storage_type        = var.spec.storage_type
  throughput_capacity = var.spec.throughput_capacity

  # Provisioned SSD read cache for the INTELLIGENT_TIERING storage class.
  dynamic "read_cache_configuration" {
    for_each = var.spec.read_cache_configuration != null ? [var.spec.read_cache_configuration] : []
    content {
      sizing_mode = read_cache_configuration.value.sizing_mode
      size        = read_cache_configuration.value.size_gib
    }
  }

  # Networking (ForceNew). Empty security_group_ids lets AWS attach the VPC
  # default SG; the MULTI_AZ_1 arms are null for the single-AZ types.
  subnet_ids                = var.spec.subnet_ids
  security_group_ids        = local.security_group_ids
  preferred_subnet_id       = local.preferred_subnet_id
  endpoint_ip_address_range = local.endpoint_ip_address_range
  route_table_ids           = local.route_table_ids

  # Encryption at rest: null falls back to the AWS-managed FSx key.
  kms_key_id = local.kms_key_id

  # Restore-from-backup create shape (capacity and settings come from the
  # backup; the spec forbids provisioning capacity alongside it).
  backup_id = local.backup_id

  # Disk IOPS: AUTOMATIC scales 3 IOPS/GiB with storage; USER_PROVISIONED
  # pins an exact number.
  dynamic "disk_iops_configuration" {
    for_each = var.spec.disk_iops_configuration != null ? [var.spec.disk_iops_configuration] : []
    content {
      mode = disk_iops_configuration.value.mode
      iops = disk_iops_configuration.value.iops
    }
  }

  # Root volume — the file system's default NFS mount target. Only emitted
  # when configured so an omitted block stays on AWS defaults.
  # copy_tags_to_snapshots is ForceNew: flipping it replaces the file system.
  dynamic "root_volume_configuration" {
    for_each = var.spec.root_volume_configuration != null ? [var.spec.root_volume_configuration] : []
    content {
      data_compression_type  = root_volume_configuration.value.data_compression_type
      read_only              = root_volume_configuration.value.read_only
      record_size_kib        = root_volume_configuration.value.record_size_kib
      copy_tags_to_snapshots = root_volume_configuration.value.copy_tags_to_snapshots

      dynamic "nfs_exports" {
        for_each = root_volume_configuration.value.nfs_exports != null ? [root_volume_configuration.value.nfs_exports] : []
        content {
          dynamic "client_configurations" {
            for_each = nfs_exports.value.client_configurations
            content {
              clients = client_configurations.value.clients
              options = client_configurations.value.options
            }
          }
        }
      }

      dynamic "user_and_group_quotas" {
        for_each = root_volume_configuration.value.user_and_group_quotas
        content {
          id                         = user_and_group_quotas.value.id
          storage_capacity_quota_gib = user_and_group_quotas.value.storage_capacity_quota_gib
          type                       = user_and_group_quotas.value.type
        }
      }
    }
  }

  # Automatic backups. Zero is a real value ("no automatic backups") — the
  # resolved default flows through as-is.
  automatic_backup_retention_days   = var.spec.automatic_backup_retention_days
  daily_automatic_backup_start_time = local.daily_automatic_backup_start_time
  copy_tags_to_backups              = var.spec.copy_tags_to_backups
  copy_tags_to_volumes              = var.spec.copy_tags_to_volumes

  # The final-backup decision is presence-honest: an explicit false ("take a
  # final backup on delete") must reach AWS rather than being coerced away.
  skip_final_backup = var.spec.skip_final_backup
  final_backup_tags = local.final_backup_tags

  # Cascading-delete opt-in: without DELETE_CHILD_VOLUMES_AND_SNAPSHOTS,
  # deletion fails while child volumes or snapshots exist.
  delete_options = local.delete_options

  # Maintenance window ("d:HH:MM", UTC).
  weekly_maintenance_start_time = local.weekly_maintenance_start_time

  tags = local.aws_tags
}
