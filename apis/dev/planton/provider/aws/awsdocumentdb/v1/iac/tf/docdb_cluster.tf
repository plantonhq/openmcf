# The cluster composes onto its neighbors instead of embedding them:
# subnets, security groups, and KMS keys attach by reference, and
# database ingress rules live on the referenced AwsSecurityGroup nodes --
# this module never creates or mutates a resource that deserves to be
# its own node. The compute serving queries is the folded instances list
# (cluster_instances.tf); this resource is the shared-storage brain:
# endpoints, credentials, backups, encryption, and engine lifecycle.
#
# Create-only in AWS: the identifier, port, subnet group, availability
# zones, master username, storage encryption + KMS key, and restore
# sources. Everything else updates in place (immediately or at the next
# maintenance window, per apply_immediately).
resource "aws_docdb_cluster" "this" {
  cluster_identifier = local.cluster_identifier

  # Empty pins nothing: AWS picks the current default DocumentDB
  # version, so an unpinned manifest never goes stale.
  engine_version = var.spec.engine_version != "" ? var.spec.engine_version : null

  # Networking: the subnet group managed here (or referenced), the VPC
  # default SG when no groups are given (AWS's own default), and
  # AWS-picked AZ spread unless explicitly pinned (the list is
  # create-only -- letting AWS choose is almost always right).
  db_subnet_group_name   = local.manage_subnet_group ? aws_docdb_subnet_group.this[0].name : (var.spec.db_subnet_group_name != "" ? var.spec.db_subnet_group_name : null)
  vpc_security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
  availability_zones     = length(var.spec.availability_zones) > 0 ? var.spec.availability_zones : null
  network_type           = var.spec.network_type != "" ? var.spec.network_type : null
  # 0 keeps the AWS default (27017). Create-only -- a port change
  # replaces the cluster.
  port = var.spec.port != 0 ? var.spec.port : null

  # "iopt1" opts into I/O-Optimized storage; empty keeps standard
  # per-I/O billing (the AWS default).
  storage_type = var.spec.storage_type != "" ? var.spec.storage_type : null

  master_username = var.spec.master_username != "" ? var.spec.master_username : null

  # The password contract (CEL enforces exactly one strategy): the
  # AWS-managed Secrets Manager secret (recommended -- no secret in
  # manifest or state) or a directly supplied password.
  # manage_master_user_password is forwarded ONLY when true: an explicit
  # false conflicts with master_password in the provider's
  # ConflictsWith machinery.
  manage_master_user_password = var.spec.manage_master_user_password ? true : null
  master_password             = var.spec.master_password != "" ? var.spec.master_password : null

  # Storage encryption is a one-way door: it can only be chosen at
  # create time, which is why the spec recommends it on by default.
  storage_encrypted = var.spec.storage_encrypted
  kms_key_id        = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # Backups: DocumentDB backups are continuous; retention bounds the
  # point-in-time recovery window. 0 keeps the AWS default (1 day).
  backup_retention_period      = var.spec.backup_retention_period != 0 ? var.spec.backup_retention_period : null
  preferred_backup_window      = var.spec.preferred_backup_window != "" ? var.spec.preferred_backup_window : null
  preferred_maintenance_window = var.spec.preferred_maintenance_window != "" ? var.spec.preferred_maintenance_window : null

  # Deletion safety: the CEL contract requires a final-snapshot name
  # unless skipping is explicit, so a delete can never fail late on a
  # missing snapshot identifier.
  skip_final_snapshot       = var.spec.skip_final_snapshot
  final_snapshot_identifier = var.spec.final_snapshot_identifier != "" ? var.spec.final_snapshot_identifier : null
  deletion_protection       = var.spec.deletion_protection

  # "audit" and "profiler" -- both also need their matching cluster
  # parameters (audit_logs / profiler) before DocumentDB emits anything.
  enabled_cloudwatch_logs_exports = length(var.spec.enabled_cloudwatch_logs_exports) > 0 ? var.spec.enabled_cloudwatch_logs_exports : null

  # DocumentDB Serverless: this block + db.serverless instances. Adding
  # or modifying scales in place; REMOVING the block from a live cluster
  # replaces it (AWS cannot switch a cluster off serverless).
  dynamic "serverless_v2_scaling_configuration" {
    for_each = var.spec.serverless_v2_scaling != null ? [var.spec.serverless_v2_scaling] : []
    content {
      min_capacity = serverless_v2_scaling_configuration.value.min_capacity
      max_capacity = serverless_v2_scaling_configuration.value.max_capacity
    }
  }

  # Create-time restore sources (mutually exclusive, CEL-enforced):
  # from a snapshot, or from another cluster's continuous backup
  # (point-in-time restore / copy-on-write fast clone).
  snapshot_identifier = var.spec.snapshot_identifier != "" ? var.spec.snapshot_identifier : null

  dynamic "restore_to_point_in_time" {
    for_each = var.spec.restore_to_point_in_time != null ? [var.spec.restore_to_point_in_time] : []
    content {
      source_cluster_identifier  = restore_to_point_in_time.value.source_cluster_identifier
      restore_to_time            = restore_to_point_in_time.value.restore_to_time != "" ? restore_to_point_in_time.value.restore_to_time : null
      use_latest_restorable_time = restore_to_point_in_time.value.use_latest_restorable_time ? true : null
      restore_type               = restore_to_point_in_time.value.restore_type != "" ? restore_to_point_in_time.value.restore_type : null
    }
  }

  # DocumentDB global cluster membership. The first cluster joined
  # becomes the global writer; later joiners are read-only secondaries.
  global_cluster_identifier = var.spec.global_cluster_identifier != "" ? var.spec.global_cluster_identifier : null

  # Parameter groups: the managed inline group, an existing referenced
  # group, or the engine default.
  db_cluster_parameter_group_name = local.effective_cluster_parameter_group

  apply_immediately           = var.spec.apply_immediately
  allow_major_version_upgrade = var.spec.allow_major_version_upgrade

  tags = local.aws_tags
}
