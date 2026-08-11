# ---------------------------------------------------------------------------
# AWS FSx for Windows File Server
# ---------------------------------------------------------------------------
# One aws_fsx_windows_file_system resource carries the whole spec. ForceNew
# attributes (replace the file system when changed): deployment_type,
# storage_type, subnet_ids, preferred_subnet_id, security_group_ids,
# kms_key_id, backup_id, active_directory_id, and copy_tags_to_backups.
# Storage capacity grows in place (never shrinks); throughput_capacity scales
# up or down in place; aliases and the audit configuration update in place.
# ---------------------------------------------------------------------------

resource "aws_fsx_windows_file_system" "this" {
  # Core shape. deployment_type is always sent (spec default SINGLE_AZ_2 —
  # the provider default is the legacy SINGLE_AZ_1); storage_capacity is
  # presence-honest — null means the capacity comes from a backup restore and
  # the argument must be omitted entirely.
  deployment_type     = var.spec.deployment_type
  storage_capacity    = var.spec.storage_capacity_gib
  storage_type        = var.spec.storage_type
  throughput_capacity = var.spec.throughput_capacity

  # Networking (ForceNew). The preferred subnet pins the active file server
  # of a MULTI_AZ_1 pair; empty security_group_ids lets AWS attach the VPC
  # default SG.
  subnet_ids          = var.spec.subnet_ids
  preferred_subnet_id = local.preferred_subnet_id
  security_group_ids  = local.security_group_ids

  # Encryption at rest: null falls back to the AWS-managed FSx key.
  kms_key_id = local.kms_key_id

  # Restore-from-backup create shape (capacity and settings come from the
  # backup; the spec forbids provisioning capacity alongside it).
  backup_id = local.backup_id

  # Active Directory — mandatory for Windows File Server; exactly one of the
  # two arms is present (CEL-enforced at validation).
  active_directory_id = local.active_directory_id

  # Self-managed AD: the domain join uses EITHER the Secrets Manager
  # service-account secret (recommended — credentials never touch the
  # manifest) OR direct username/password; empty strings become null so the
  # provider sees exactly one arm.
  dynamic "self_managed_active_directory" {
    for_each = var.spec.self_managed_active_directory != null ? [var.spec.self_managed_active_directory] : []
    content {
      domain_name                            = self_managed_active_directory.value.domain_name
      dns_ips                                = self_managed_active_directory.value.dns_ips
      username                               = self_managed_active_directory.value.username != "" ? self_managed_active_directory.value.username : null
      password                               = self_managed_active_directory.value.password != "" ? self_managed_active_directory.value.password : null
      domain_join_service_account_secret     = self_managed_active_directory.value.domain_join_service_account_secret_arn != "" ? self_managed_active_directory.value.domain_join_service_account_secret_arn : null
      file_system_administrators_group       = self_managed_active_directory.value.file_system_administrators_group != "" ? self_managed_active_directory.value.file_system_administrators_group : null
      organizational_unit_distinguished_name = self_managed_active_directory.value.organizational_unit_distinguished_name != "" ? self_managed_active_directory.value.organizational_unit_distinguished_name : null
    }
  }

  # DNS aliases (in-place; each needs a CNAME to the file system's DNS name).
  aliases = local.aliases

  # Windows audit logging to CloudWatch Logs.
  dynamic "audit_log_configuration" {
    for_each = var.spec.audit_log_configuration != null ? [var.spec.audit_log_configuration] : []
    content {
      file_access_audit_log_level       = audit_log_configuration.value.file_access_audit_log_level
      file_share_access_audit_log_level = audit_log_configuration.value.file_share_access_audit_log_level
      audit_log_destination             = audit_log_configuration.value.audit_log_destination != "" ? audit_log_configuration.value.audit_log_destination : null
    }
  }

  # SSD IOPS: AUTOMATIC scales 3 IOPS/GiB with storage; USER_PROVISIONED pins
  # an exact number (0-350000).
  dynamic "disk_iops_configuration" {
    for_each = var.spec.disk_iops_configuration != null ? [var.spec.disk_iops_configuration] : []
    content {
      mode = disk_iops_configuration.value.mode
      iops = disk_iops_configuration.value.iops
    }
  }

  # Automatic backups. Zero is a real value ("no automatic backups") — the
  # resolved default (7) flows through as-is.
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
