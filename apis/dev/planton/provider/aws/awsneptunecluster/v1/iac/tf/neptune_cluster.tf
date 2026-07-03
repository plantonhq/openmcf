# The cluster composes onto its neighbors instead of embedding them:
# subnets, security groups, KMS keys, and IAM roles attach by reference,
# and database ingress rules live on the referenced AwsSecurityGroup
# nodes -- this module never creates or mutates a resource that deserves
# to be its own node. The compute serving queries is the folded
# instances list (cluster_instances.tf); this resource is the
# shared-storage brain: endpoints, backups, encryption, and engine
# lifecycle.
#
# Neptune has no master username or password -- access is network
# reachability plus (optionally) IAM database authentication.
#
# Create-only in AWS: the identifier, port, subnet group, availability
# zones, storage encryption + KMS key, and snapshot_identifier.
# Everything else updates in place (immediately or at the next
# maintenance window, per apply_immediately).
resource "aws_neptune_cluster" "this" {
  cluster_identifier = local.cluster_identifier

  engine = "neptune"
  # Empty pins nothing: AWS picks the current default Neptune version,
  # so an unpinned manifest never goes stale.
  engine_version = var.spec.engine_version != "" ? var.spec.engine_version : null

  # Networking: the subnet group managed here (or referenced), the VPC
  # default SG when no groups are given (AWS's own default), and
  # AWS-picked AZ spread unless explicitly pinned (the list is
  # create-only -- letting AWS choose is almost always right).
  neptune_subnet_group_name = local.manage_subnet_group ? aws_neptune_subnet_group.this[0].name : (var.spec.neptune_subnet_group_name != "" ? var.spec.neptune_subnet_group_name : null)
  vpc_security_group_ids    = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
  availability_zones        = length(var.spec.availability_zones) > 0 ? var.spec.availability_zones : null
  # 0 keeps the AWS default (8182). Create-only -- a port change
  # replaces the cluster.
  port = var.spec.port != 0 ? var.spec.port : null

  # "iopt1" opts into I/O-Optimized storage (engine 1.3+); empty keeps
  # standard per-I/O billing (the AWS default).
  storage_type = var.spec.storage_type != "" ? var.spec.storage_type : null

  # Storage encryption is a one-way door: it can only be chosen at
  # create time, which is why the spec recommends it on by default.
  storage_encrypted = var.spec.storage_encrypted
  kms_key_arn       = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # SigV4-signed requests from IAM identities -- Neptune's only
  # credential mechanism.
  iam_database_authentication_enabled = var.spec.iam_database_authentication_enabled
  # Roles the ENGINE assumes for the S3 bulk loader and Neptune ML. The
  # roles own their policies -- this cluster only associates them (a
  # module never mutates a resource it references).
  iam_roles = length(var.spec.iam_roles) > 0 ? var.spec.iam_roles : null

  # Backups: Neptune backups are continuous; retention bounds the
  # point-in-time recovery window. 0 keeps the AWS default (1 day).
  backup_retention_period      = var.spec.backup_retention_period != 0 ? var.spec.backup_retention_period : null
  preferred_backup_window      = var.spec.preferred_backup_window != "" ? var.spec.preferred_backup_window : null
  preferred_maintenance_window = var.spec.preferred_maintenance_window != "" ? var.spec.preferred_maintenance_window : null
  copy_tags_to_snapshot        = var.spec.copy_tags_to_snapshot

  # Deletion safety: the CEL contract requires a final-snapshot name
  # unless skipping is explicit, so a delete can never fail late on a
  # missing snapshot identifier.
  skip_final_snapshot       = var.spec.skip_final_snapshot
  final_snapshot_identifier = var.spec.final_snapshot_identifier != "" ? var.spec.final_snapshot_identifier : null
  deletion_protection       = var.spec.deletion_protection

  # "audit" and "slowquery" -- both also need their matching cluster
  # parameters (neptune_enable_audit_log / the slow-query threshold
  # parameters) before Neptune emits anything.
  enable_cloudwatch_logs_exports = length(var.spec.enabled_cloudwatch_logs_exports) > 0 ? var.spec.enabled_cloudwatch_logs_exports : null

  # Neptune Serverless: this block + db.serverless instances. NCU
  # bounds are 1-128 on both ends.
  dynamic "serverless_v2_scaling_configuration" {
    for_each = var.spec.serverless_v2_scaling != null ? [var.spec.serverless_v2_scaling] : []
    content {
      min_capacity = serverless_v2_scaling_configuration.value.min_capacity
      max_capacity = serverless_v2_scaling_configuration.value.max_capacity
    }
  }

  # Create-time restore source: a manual or automated cluster snapshot.
  snapshot_identifier = var.spec.snapshot_identifier != "" ? var.spec.snapshot_identifier : null

  # Cross-cluster replication and Neptune global database membership.
  replication_source_identifier = var.spec.replication_source_identifier != "" ? var.spec.replication_source_identifier : null
  global_cluster_identifier     = var.spec.global_cluster_identifier != "" ? var.spec.global_cluster_identifier : null

  # Parameter groups: the managed inline group, an existing referenced
  # group, or the engine default. neptune_instance_parameter_group_name
  # is only consulted by AWS during a major version upgrade (the spec's
  # CEL requires it alongside allow_major_version_upgrade).
  neptune_cluster_parameter_group_name  = local.effective_cluster_parameter_group
  neptune_instance_parameter_group_name = var.spec.neptune_instance_parameter_group_name != "" ? var.spec.neptune_instance_parameter_group_name : null

  apply_immediately           = var.spec.apply_immediately
  allow_major_version_upgrade = var.spec.allow_major_version_upgrade

  tags = local.aws_tags
}
