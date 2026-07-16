# ---------------------------------------------------------------------------
# AWS FSx for Lustre file system
# ---------------------------------------------------------------------------
# One aws_fsx_lustre_file_system resource carries the whole spec. ForceNew
# attributes (replace the file system when changed): deployment_type,
# storage_type, subnet_ids, security_group_ids, kms_key_id, backup_id,
# efa_enabled, drive_cache_type, copy_tags_to_backups, and the legacy S3 link
# arm (import_path / export_path / imported_file_chunk_size). Storage capacity
# grows in place (a decrease — or any growth on SCRATCH_1 — replaces the file
# system); metadata IOPS grow in place and shrink by replacement.
# ---------------------------------------------------------------------------

resource "aws_fsx_lustre_file_system" "this" {
  # Core shape. deployment_type is always sent: the spec default (SCRATCH_2)
  # deliberately diverges from the provider default (legacy SCRATCH_1), so
  # relying on provider defaulting would deploy the wrong generation.
  deployment_type = var.spec.deployment_type
  storage_type    = var.spec.storage_type

  # Presence-honest capacity: null means the capacity comes from a backup
  # restore or the elastic INTELLIGENT_TIERING class, and the argument must
  # be omitted entirely (AWS rejects explicit capacity for both shapes).
  storage_capacity = var.spec.storage_capacity_gib

  # Throughput is one knob per storage generation: per-TiB throughput for
  # provisioned (SSD/HDD) PERSISTENT capacity, absolute throughput for
  # INTELLIGENT_TIERING. The spec's CEL rules keep them mutually exclusive.
  per_unit_storage_throughput = var.spec.per_unit_storage_throughput
  throughput_capacity         = var.spec.throughput_capacity

  data_compression_type    = var.spec.data_compression_type
  file_system_type_version = local.file_system_type_version

  # EFA/GPUDirect must be decided at creation (ForceNew) and pins
  # per_unit_storage_throughput while enabled. Sent only when true so
  # unset stays on AWS's computed default.
  efa_enabled = var.spec.efa_enabled ? true : null

  # HDD read cache decision — required by AWS for HDD storage (CEL-enforced
  # at validation), meaningless elsewhere.
  drive_cache_type = local.drive_cache_type

  # Provisioned SSD read cache for the INTELLIGENT_TIERING storage class.
  dynamic "data_read_cache_configuration" {
    for_each = var.spec.data_read_cache_configuration != null ? [var.spec.data_read_cache_configuration] : []
    content {
      sizing_mode = data_read_cache_configuration.value.sizing_mode
      size        = data_read_cache_configuration.value.size_gib
    }
  }

  # Networking (ForceNew — Lustre is single-AZ, exactly one subnet).
  subnet_ids         = [var.spec.subnet_id]
  security_group_ids = local.security_group_ids

  # Encryption at rest: null falls back to the AWS-managed FSx key.
  kms_key_id = local.kms_key_id

  # Restore-from-backup create shape (capacity and settings come from the
  # backup; the spec forbids provisioning capacity alongside it).
  backup_id = local.backup_id

  # Legacy S3 data repository link (not PERSISTENT_2 — data repository
  # associations are the modern, many-per-filesystem generation).
  import_path              = local.import_path
  export_path              = local.export_path
  auto_import_policy       = local.auto_import_policy
  imported_file_chunk_size = var.spec.imported_file_chunk_size

  # POSIX root squash — maps root clients to an unprivileged UID:GID, with
  # NID-listed administrative hosts exempt. Updatable in place.
  dynamic "root_squash_configuration" {
    for_each = var.spec.root_squash_configuration != null ? [var.spec.root_squash_configuration] : []
    content {
      root_squash    = root_squash_configuration.value.root_squash != "" ? root_squash_configuration.value.root_squash : null
      no_squash_nids = length(root_squash_configuration.value.no_squash_nids) > 0 ? root_squash_configuration.value.no_squash_nids : null
    }
  }

  # CloudWatch logging for data repository events.
  dynamic "log_configuration" {
    for_each = var.spec.log_configuration != null ? [var.spec.log_configuration] : []
    content {
      destination = log_configuration.value.destination != "" ? log_configuration.value.destination : null
      level       = log_configuration.value.level
    }
  }

  # Metadata IOPS configuration (PERSISTENT_2 only; CEL-enforced).
  dynamic "metadata_configuration" {
    for_each = var.spec.metadata_configuration != null ? [var.spec.metadata_configuration] : []
    content {
      mode = metadata_configuration.value.mode
      iops = metadata_configuration.value.iops
    }
  }

  # Automatic backups (PERSISTENT deployments). Zero is a real value ("no
  # automatic backups") — the resolved default flows through as-is.
  automatic_backup_retention_days   = var.spec.automatic_backup_retention_days
  daily_automatic_backup_start_time = local.daily_automatic_backup_start_time
  copy_tags_to_backups              = var.spec.copy_tags_to_backups

  # The final-backup decision is presence-honest: an explicit false ("take a
  # final backup on delete") must reach AWS rather than being coerced away.
  skip_final_backup = var.spec.skip_final_backup
  final_backup_tags = local.final_backup_tags

  # Maintenance window ("d:HH:MM", UTC).
  weekly_maintenance_start_time = local.weekly_maintenance_start_time

  tags = local.aws_tags
}
