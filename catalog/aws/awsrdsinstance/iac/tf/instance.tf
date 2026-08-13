# The instance composes onto its neighbors instead of embedding them:
# subnets, security groups, KMS keys, and the monitoring role attach by
# reference, and database ingress rules live on the referenced
# AwsSecurityGroup nodes -- this module never creates or mutates a
# resource that deserves to be its own node.
#
# Create-only in AWS: the engine, username, db_name, character sets,
# timezone, availability-zone pin, storage encryption + KMS key, and the
# restore sources. Everything else updates in place (immediately or at
# the next maintenance window, per apply_immediately). Growing storage
# applies in place; shrinking requires a new instance.
resource "aws_db_instance" "this" {
  identifier = local.instance_identifier

  # Empty on a replica or restore: AWS derives the engine from the
  # source (the CEL contract requires it otherwise).
  engine = var.spec.engine != "" ? var.spec.engine : null
  # Empty pins nothing: AWS picks the engine's current default version,
  # so an unpinned manifest never goes stale.
  engine_version = var.spec.engine_version != "" ? var.spec.engine_version : null

  instance_class = var.spec.instance_class

  # Storage: allocated is inherited from the source on replicas and
  # restores (0 -> null); the autoscaling ceiling is the cheap insurance
  # against disk-full outages.
  allocated_storage     = var.spec.allocated_storage_gb != 0 ? var.spec.allocated_storage_gb : null
  max_allocated_storage = var.spec.max_allocated_storage_gb != 0 ? var.spec.max_allocated_storage_gb : null
  storage_type          = var.spec.storage_type != "" ? var.spec.storage_type : null
  iops                  = var.spec.iops != 0 ? var.spec.iops : null
  storage_throughput    = var.spec.storage_throughput != 0 ? var.spec.storage_throughput : null
  # A dedicated EBS volume for logs -- steadier I/O for audit-heavy or
  # WAL-heavy workloads. Forwarded only when true.
  dedicated_log_volume = var.spec.dedicated_log_volume ? true : null

  # Storage encryption is a one-way door: it can only be chosen at
  # create time, which is why the spec recommends it on by default.
  storage_encrypted = var.spec.storage_encrypted
  kms_key_id        = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  db_name  = var.spec.db_name != "" ? var.spec.db_name : null
  username = var.spec.username != "" ? var.spec.username : null

  # The three-way password contract (CEL enforces exactly one strategy):
  # AWS-managed secret (recommended -- no secret in manifest or state)
  # or a directly supplied password. manage_master_user_password must be
  # forwarded ONLY when true: an explicit false conflicts with password
  # in the provider's ConflictsWith machinery.
  manage_master_user_password   = var.spec.manage_master_user_password ? true : null
  master_user_secret_kms_key_id = var.spec.master_user_secret_kms_key_id != "" ? var.spec.master_user_secret_kms_key_id : null
  password                      = var.spec.password != "" ? var.spec.password : null

  # Networking: the subnet group managed here (or referenced), the VPC
  # default SG when no groups are given (AWS's own default).
  db_subnet_group_name   = local.manage_subnet_group ? aws_db_subnet_group.this[0].name : (var.spec.db_subnet_group_name != "" ? var.spec.db_subnet_group_name : null)
  vpc_security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
  network_type           = var.spec.network_type != "" ? var.spec.network_type : null
  port                   = var.spec.port != 0 ? var.spec.port : null

  # Availability: a synchronous standby with automatic failover
  # (multi_az), or a single-AZ instance optionally pinned to one zone --
  # the two are mutually exclusive (CEL).
  multi_az            = var.spec.multi_az
  availability_zone   = var.spec.availability_zone != "" ? var.spec.availability_zone : null
  publicly_accessible = var.spec.publicly_accessible

  # Read replica: engine, storage, and credentials are inherited from
  # the source (CEL keeps them empty here). replica_mode is Oracle's
  # queryable-vs-mounted choice.
  replicate_source_db = var.spec.replicate_source_db != "" ? var.spec.replicate_source_db : null
  replica_mode        = var.spec.replica_mode != "" ? var.spec.replica_mode : null

  # Create-time restore sources (mutually exclusive with each other and
  # with replicate_source_db, CEL-enforced). s3_import is the MySQL
  # migration on-ramp: restore a Percona XtraBackup straight from S3.
  snapshot_identifier = var.spec.snapshot_identifier != "" ? var.spec.snapshot_identifier : null

  dynamic "s3_import" {
    for_each = var.spec.s3_import != null ? [var.spec.s3_import] : []
    content {
      bucket_name           = s3_import.value.bucket_name
      bucket_prefix         = s3_import.value.bucket_prefix != "" ? s3_import.value.bucket_prefix : null
      ingestion_role        = s3_import.value.ingestion_role
      source_engine         = s3_import.value.source_engine
      source_engine_version = s3_import.value.source_engine_version
    }
  }

  # One-way storage file-system upgrade, consulted on replicas and
  # snapshot restores whose source still runs the older configuration.
  upgrade_storage_config = var.spec.upgrade_storage_config ? true : null

  dynamic "restore_to_point_in_time" {
    for_each = var.spec.restore_to_point_in_time != null ? [var.spec.restore_to_point_in_time] : []
    content {
      source_db_instance_identifier            = restore_to_point_in_time.value.source_db_instance_identifier != "" ? restore_to_point_in_time.value.source_db_instance_identifier : null
      source_dbi_resource_id                   = restore_to_point_in_time.value.source_dbi_resource_id != "" ? restore_to_point_in_time.value.source_dbi_resource_id : null
      source_db_instance_automated_backups_arn = restore_to_point_in_time.value.source_db_instance_automated_backups_arn != "" ? restore_to_point_in_time.value.source_db_instance_automated_backups_arn : null
      restore_time                             = restore_to_point_in_time.value.restore_time != "" ? restore_to_point_in_time.value.restore_time : null
      use_latest_restorable_time               = restore_to_point_in_time.value.use_latest_restorable_time ? true : null
    }
  }

  # Backups: 0 disables automated backups (and point-in-time recovery);
  # null keeps the AWS default (which differs for replicas), so 0 is
  # forwarded explicitly only through the retention field's semantics.
  backup_retention_period = var.spec.backup_retention_period != 0 ? var.spec.backup_retention_period : null
  backup_window           = var.spec.backup_window != "" ? var.spec.backup_window : null
  maintenance_window      = var.spec.maintenance_window != "" ? var.spec.maintenance_window : null
  copy_tags_to_snapshot   = var.spec.copy_tags_to_snapshot
  # Tri-state: null keeps the AWS default (true). An explicit false
  # retains automated backups after deletion -- the last line of defense
  # against a mistaken teardown.
  delete_automated_backups = var.spec.delete_automated_backups

  # Deletion safety: the CEL contract requires a final-snapshot name
  # unless skipping is explicit, so a delete can never fail late on a
  # missing snapshot identifier.
  skip_final_snapshot       = var.spec.skip_final_snapshot
  final_snapshot_identifier = var.spec.final_snapshot_identifier != "" ? var.spec.final_snapshot_identifier : null
  deletion_protection       = var.spec.deletion_protection

  iam_database_authentication_enabled = var.spec.iam_database_authentication_enabled

  enabled_cloudwatch_logs_exports = length(var.spec.enabled_cloudwatch_logs_exports) > 0 ? var.spec.enabled_cloudwatch_logs_exports : null

  # Observability: Performance Insights (per-query telemetry, free at
  # 7-day retention), Enhanced Monitoring (OS-level metrics through the
  # referenced role), and Database Insights.
  performance_insights_enabled          = var.spec.performance_insights_enabled ? true : null
  performance_insights_kms_key_id       = var.spec.performance_insights_kms_key_id != "" ? var.spec.performance_insights_kms_key_id : null
  performance_insights_retention_period = var.spec.performance_insights_retention_period != 0 ? var.spec.performance_insights_retention_period : null
  monitoring_interval                   = var.spec.monitoring_interval != 0 ? var.spec.monitoring_interval : null
  monitoring_role_arn                   = var.spec.monitoring_role_arn != "" ? var.spec.monitoring_role_arn : null
  database_insights_mode                = var.spec.database_insights_mode != "" ? var.spec.database_insights_mode : null

  # Engine-configuration attachments: the managed inline group, an
  # existing referenced group, or the engine default -- parameter and
  # option groups are configuration lists, not composable nodes.
  parameter_group_name = local.effective_parameter_group
  option_group_name    = local.effective_option_group

  # Active Directory join -- either the AWS-managed directory shape or
  # the self-managed shape (never both; CEL-enforced). Nested ternaries,
  # never `&&`: HCL's `&&` does not short-circuit, so a pruned
  # active_directory block would error on the inner attribute access.
  domain                 = local.active_directory == null ? null : (local.active_directory.domain != "" ? local.active_directory.domain : null)
  domain_iam_role_name   = local.active_directory == null ? null : (local.active_directory.domain_iam_role_name != "" ? local.active_directory.domain_iam_role_name : null)
  domain_fqdn            = local.active_directory == null ? null : (local.active_directory.domain_fqdn != "" ? local.active_directory.domain_fqdn : null)
  domain_ou              = local.active_directory == null ? null : (local.active_directory.domain_ou != "" ? local.active_directory.domain_ou : null)
  domain_auth_secret_arn = local.active_directory == null ? null : (local.active_directory.domain_auth_secret_arn != "" ? local.active_directory.domain_auth_secret_arn : null)
  domain_dns_ips         = local.active_directory == null ? null : (length(local.active_directory.domain_dns_ips) > 0 ? local.active_directory.domain_dns_ips : null)

  license_model            = var.spec.license_model != "" ? var.spec.license_model : null
  character_set_name       = var.spec.character_set_name != "" ? var.spec.character_set_name : null
  nchar_character_set_name = var.spec.nchar_character_set_name != "" ? var.spec.nchar_character_set_name : null
  timezone                 = var.spec.timezone != "" ? var.spec.timezone : null

  # Extended support posture -- empty keeps the AWS default (paid
  # extended support engages automatically when the version leaves
  # standard support).
  engine_lifecycle_support = var.spec.engine_lifecycle_support != "" ? var.spec.engine_lifecycle_support : null

  ca_cert_identifier = var.spec.ca_cert_identifier != "" ? var.spec.ca_cert_identifier : null

  # RDS Blue/Green Deployments: a synchronized green copy takes the
  # change and switches over in under a minute -- near-zero-downtime
  # engine upgrades and parameter changes.
  dynamic "blue_green_update" {
    for_each = var.spec.blue_green_update_enabled ? [true] : []
    content {
      enabled = true
    }
  }

  # Tri-state: null keeps the AWS default (true).
  auto_minor_version_upgrade  = var.spec.auto_minor_version_upgrade
  allow_major_version_upgrade = var.spec.allow_major_version_upgrade
  apply_immediately           = var.spec.apply_immediately

  tags = local.aws_tags
}
