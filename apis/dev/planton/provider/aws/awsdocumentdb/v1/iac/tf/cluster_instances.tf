# The cluster's compute: one aws_docdb_cluster_instance per folded
# instances entry, keyed by the entry's name so adding or removing a
# reader is an in-place update that never touches the cluster or its
# siblings. The instance with the lowest promotion tier that is
# available becomes the writer; all others serve reads from the shared
# cluster storage.
#
# The engine version is inherited from the cluster by AWS itself --
# DocumentDB stamps every instance with the cluster's resolved version,
# and a version change rolls the instances in the same apply.
resource "aws_docdb_cluster_instance" "this" {
  for_each = { for instance in var.spec.instances : instance.name => instance }

  identifier         = "${local.cluster_identifier}-${each.value.name}"
  cluster_identifier = aws_docdb_cluster.this.id

  # "db.serverless" makes this a DocumentDB Serverless instance scaling
  # within the cluster's serverless_v2_scaling bounds; any provisioned
  # class fixes its size.
  instance_class = each.value.instance_class

  # Failover priority (0 promoted first). AWS default is already 0 --
  # forwarded unconditionally since 0 is the meaningful base tier.
  promotion_tier = each.value.promotion_tier

  # Empty lets AWS spread instances across the cluster's zones -- the
  # right call almost always; a pin is create-only.
  availability_zone = each.value.availability_zone != "" ? each.value.availability_zone : null

  # Tri-state: null keeps the AWS default (true).
  auto_minor_version_upgrade = each.value.auto_minor_version_upgrade

  # Performance Insights is instance-scoped on DocumentDB (there is no
  # cluster-level setting).
  enable_performance_insights     = each.value.performance_insights_enabled ? true : null
  performance_insights_kms_key_id = each.value.performance_insights_kms_key_id != "" ? each.value.performance_insights_kms_key_id : null

  # A per-instance maintenance window; empty inherits AWS scheduling.
  preferred_maintenance_window = each.value.preferred_maintenance_window != "" ? each.value.preferred_maintenance_window : null

  ca_cert_identifier = each.value.ca_cert_identifier != "" ? each.value.ca_cert_identifier : null

  copy_tags_to_snapshot = each.value.copy_tags_to_snapshot

  tags = local.aws_tags
}
