# The cluster composes onto its neighbors instead of embedding them:
# subnets, security groups, KMS keys, and IAM roles attach by reference,
# and database ingress rules live on the referenced AwsSecurityGroup
# nodes -- this module never creates or mutates a resource that deserves
# to be its own node. The compute serving queries is the folded
# instances list (cluster_instances.tf); this resource is the
# shared-storage brain: endpoints, credentials, backups, encryption,
# and engine lifecycle.
#
# Create-only in AWS: the identifier, engine, engine_mode, subnet group,
# availability zones, master username, database name, storage
# encryption + KMS key, restore sources, and source_region. Everything
# else updates in place (immediately or at the next maintenance window,
# per apply_immediately).
resource "aws_rds_cluster" "this" {
  cluster_identifier = local.cluster_identifier

  engine = var.spec.engine
  # Empty pins nothing: AWS picks the engine's current default version,
  # so an unpinned manifest never goes stale.
  engine_version = var.spec.engine_version != "" ? var.spec.engine_version : null
  # Empty and "provisioned" are the same thing to AWS; forwarding null
  # keeps the diff surface minimal. Serverless v2 is provisioned mode +
  # a serverlessv2 block -- only legacy Serverless v1 sets "serverless".
  engine_mode              = var.spec.engine_mode != "" ? var.spec.engine_mode : null
  engine_lifecycle_support = var.spec.engine_lifecycle_support != "" ? var.spec.engine_lifecycle_support : null

  # Networking: the subnet group managed here (or referenced), the VPC
  # default SG when no groups are given (AWS's own default), and
  # AWS-picked AZ spread unless explicitly pinned (the list is
  # create-only -- letting AWS choose is almost always right).
  db_subnet_group_name   = local.manage_subnet_group ? aws_db_subnet_group.this[0].name : (var.spec.db_subnet_group_name != "" ? var.spec.db_subnet_group_name : null)
  vpc_security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
  availability_zones     = length(var.spec.availability_zones) > 0 ? var.spec.availability_zones : null
  network_type           = var.spec.network_type != "" ? var.spec.network_type : null
  port                   = var.spec.port != 0 ? var.spec.port : null

  # Multi-AZ RDS cluster shape (community mysql/postgres engines): AWS
  # manages one writer + two readers internally, sized here. The CEL
  # rules guarantee these are set exactly when the engine calls for
  # them, and that the folded instances list stays empty.
  db_cluster_instance_class = var.spec.db_cluster_instance_class != "" ? var.spec.db_cluster_instance_class : null
  allocated_storage         = var.spec.allocated_storage_gb != 0 ? var.spec.allocated_storage_gb : null
  iops                      = var.spec.iops != 0 ? var.spec.iops : null
  storage_type              = var.spec.storage_type != "" ? var.spec.storage_type : null

  database_name   = var.spec.database_name != "" ? var.spec.database_name : null
  master_username = var.spec.master_username != "" ? var.spec.master_username : null

  # The three-way password contract (CEL enforces exactly one strategy):
  # AWS-managed secret (recommended -- no secret in manifest or state)
  # or a directly supplied password. manage_master_user_password must be
  # forwarded ONLY when true: an explicit false conflicts with
  # master_password in the provider's ConflictsWith machinery.
  manage_master_user_password   = var.spec.manage_master_user_password ? true : null
  master_user_secret_kms_key_id = var.spec.master_user_secret_kms_key_id != "" ? var.spec.master_user_secret_kms_key_id : null
  master_password               = var.spec.master_password != "" ? var.spec.master_password : null

  # Storage encryption is a one-way door: it can only be chosen at
  # create time, which is why the spec recommends it on by default.
  storage_encrypted = var.spec.storage_encrypted
  kms_key_id        = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # Backups: Aurora backups are continuous; retention bounds the
  # point-in-time recovery window. 0 keeps the AWS default (1 day).
  backup_retention_period      = var.spec.backup_retention_period != 0 ? var.spec.backup_retention_period : null
  preferred_backup_window      = var.spec.preferred_backup_window != "" ? var.spec.preferred_backup_window : null
  preferred_maintenance_window = var.spec.preferred_maintenance_window != "" ? var.spec.preferred_maintenance_window : null
  copy_tags_to_snapshot        = var.spec.copy_tags_to_snapshot
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

  # Aurora MySQL in-place rewind. 0 disables; enabling later is not
  # supported by AWS, so production Aurora MySQL wants this at create.
  backtrack_window = var.spec.backtrack_window_seconds != 0 ? var.spec.backtrack_window_seconds : null

  iam_database_authentication_enabled = var.spec.iam_database_authentication_enabled
  # Engine feature roles are managed as one association resource per
  # spec.iam_roles entry (role_associations.tf) -- never this resource's
  # inline iam_roles argument, which cannot carry feature names and,
  # per the provider's own warning, overwrites association resources
  # when the two are mixed.

  # The Data API: SQL over HTTPS with IAM auth -- the natural fit for
  # Lambda and other connection-averse callers.
  enable_http_endpoint = var.spec.enable_http_endpoint

  enabled_cloudwatch_logs_exports = length(var.spec.enabled_cloudwatch_logs_exports) > 0 ? var.spec.enabled_cloudwatch_logs_exports : null

  # Cluster-level observability; per-instance overrides live on the
  # folded instance entries.
  performance_insights_enabled          = var.spec.performance_insights_enabled ? true : null
  performance_insights_kms_key_id       = var.spec.performance_insights_kms_key_id != "" ? var.spec.performance_insights_kms_key_id : null
  performance_insights_retention_period = var.spec.performance_insights_retention_period != 0 ? var.spec.performance_insights_retention_period : null
  monitoring_interval                   = var.spec.monitoring_interval != 0 ? var.spec.monitoring_interval : null
  monitoring_role_arn                   = var.spec.monitoring_role_arn != "" ? var.spec.monitoring_role_arn : null
  database_insights_mode                = var.spec.database_insights_mode != "" ? var.spec.database_insights_mode : null

  # Aurora Serverless v2: provisioned mode + this block + db.serverless
  # instances. min_capacity 0 enables automatic pause -- compute cost
  # drops to zero while idle, resumed on the next connection.
  dynamic "serverlessv2_scaling_configuration" {
    for_each = var.spec.serverless_v2_scaling != null ? [var.spec.serverless_v2_scaling] : []
    content {
      min_capacity             = serverlessv2_scaling_configuration.value.min_capacity
      max_capacity             = serverlessv2_scaling_configuration.value.max_capacity
      seconds_until_auto_pause = serverlessv2_scaling_configuration.value.seconds_until_auto_pause != 0 ? serverlessv2_scaling_configuration.value.seconds_until_auto_pause : null
    }
  }

  # Legacy Aurora Serverless v1 (engine_mode "serverless") -- AWS owns
  # the compute entirely; the folded instances list stays empty (CEL).
  dynamic "scaling_configuration" {
    for_each = var.spec.serverless_v1_scaling != null ? [var.spec.serverless_v1_scaling] : []
    content {
      auto_pause               = scaling_configuration.value.auto_pause
      min_capacity             = scaling_configuration.value.min_capacity != 0 ? scaling_configuration.value.min_capacity : null
      max_capacity             = scaling_configuration.value.max_capacity != 0 ? scaling_configuration.value.max_capacity : null
      seconds_until_auto_pause = scaling_configuration.value.seconds_until_auto_pause != 0 ? scaling_configuration.value.seconds_until_auto_pause : null
      timeout_action           = scaling_configuration.value.timeout_action != "" ? scaling_configuration.value.timeout_action : null
    }
  }

  # Kerberos authentication through an AWS Managed Microsoft AD --
  # clusters only support the managed-directory shape (the pair is
  # CEL-coupled; self-managed AD is an instance-kind capability).
  domain               = var.spec.domain != "" ? var.spec.domain : null
  domain_iam_role_name = var.spec.domain_iam_role_name != "" ? var.spec.domain_iam_role_name : null

  # Tri-state: null keeps the AWS default (true). Per-instance
  # auto_minor_version_upgrade on the folded instances overrides this
  # for that instance.
  auto_minor_version_upgrade = var.spec.auto_minor_version_upgrade

  # Create-time restore sources (mutually exclusive, CEL-enforced):
  # from a snapshot, from another cluster's continuous backup
  # (point-in-time restore / copy-on-write fast clone), or from a
  # Percona XtraBackup in S3 (aurora-mysql migration on-ramp).
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

  dynamic "restore_to_point_in_time" {
    for_each = var.spec.restore_to_point_in_time != null ? [var.spec.restore_to_point_in_time] : []
    content {
      source_cluster_identifier  = restore_to_point_in_time.value.source_cluster_identifier != "" ? restore_to_point_in_time.value.source_cluster_identifier : null
      source_cluster_resource_id = restore_to_point_in_time.value.source_cluster_resource_id != "" ? restore_to_point_in_time.value.source_cluster_resource_id : null
      restore_to_time            = restore_to_point_in_time.value.restore_to_time != "" ? restore_to_point_in_time.value.restore_to_time : null
      use_latest_restorable_time = restore_to_point_in_time.value.use_latest_restorable_time ? true : null
      restore_type               = restore_to_point_in_time.value.restore_type != "" ? restore_to_point_in_time.value.restore_type : null
    }
  }

  # Cross-region replication and Aurora Global Database membership.
  replication_source_identifier = var.spec.replication_source_identifier != "" ? var.spec.replication_source_identifier : null
  source_region                 = var.spec.source_region != "" ? var.spec.source_region : null
  global_cluster_identifier     = var.spec.global_cluster_identifier != "" ? var.spec.global_cluster_identifier : null
  # Forwarded only when true: the provider defaults both to false, and a
  # bare false on a non-global cluster is a pointless diff.
  enable_global_write_forwarding = var.spec.enable_global_write_forwarding ? true : null
  enable_local_write_forwarding  = var.spec.enable_local_write_forwarding ? true : null

  # Parameter groups: the managed inline group, an existing referenced
  # group, or the engine default. db_instance_parameter_group_name is
  # only consulted by AWS during a major version upgrade.
  db_cluster_parameter_group_name  = local.effective_cluster_parameter_group
  db_instance_parameter_group_name = var.spec.db_instance_parameter_group_name != "" ? var.spec.db_instance_parameter_group_name : null

  ca_certificate_identifier = var.spec.ca_certificate_identifier != "" ? var.spec.ca_certificate_identifier : null

  apply_immediately           = var.spec.apply_immediately
  allow_major_version_upgrade = var.spec.allow_major_version_upgrade

  tags = local.aws_tags
}
