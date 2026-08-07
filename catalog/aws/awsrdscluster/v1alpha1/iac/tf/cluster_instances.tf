# The cluster's compute: one aws_rds_cluster_instance per folded
# instances entry, keyed by the entry's name so adding or removing a
# reader is an in-place update that never touches the cluster or its
# siblings. The instance with the lowest promotion tier that is
# available becomes the writer; all others serve reads from the shared
# cluster storage.
#
# Engine and version are inherited from the cluster resource by
# attribute reference -- so an unpinned cluster (AWS-default version)
# still stamps its instances with the resolved version, and a version
# change rolls the instances in the same apply.
resource "aws_rds_cluster_instance" "this" {
  for_each = { for instance in var.spec.instances : instance.name => instance }

  identifier         = "${local.cluster_identifier}-${each.value.name}"
  cluster_identifier = aws_rds_cluster.this.id
  engine             = aws_rds_cluster.this.engine
  engine_version     = aws_rds_cluster.this.engine_version

  # "db.serverless" makes this an Aurora Serverless v2 instance scaling
  # within the cluster's serverless_v2_scaling bounds; any provisioned
  # class fixes its size.
  instance_class = each.value.instance_class

  # Failover priority (0 promoted first). AWS default is already 0 --
  # forwarded unconditionally since 0 is the meaningful base tier.
  promotion_tier = each.value.promotion_tier

  # Empty lets AWS spread instances across the cluster's zones -- the
  # right call almost always; a pin is create-only.
  availability_zone = each.value.availability_zone != "" ? each.value.availability_zone : null

  publicly_accessible = each.value.publicly_accessible

  # Instance-LEVEL parameter group (engine tunables scoped to one
  # instance); the cluster-level group lives on the cluster resource.
  db_parameter_group_name = each.value.db_parameter_group_name != "" ? each.value.db_parameter_group_name : null

  # Tri-state: null keeps the AWS default (true).
  auto_minor_version_upgrade = each.value.auto_minor_version_upgrade

  # Tri-state: null inherits the cluster-level Performance Insights
  # posture; an explicit value overrides it for this instance only.
  performance_insights_enabled = each.value.performance_insights_enabled

  # Per-instance Enhanced Monitoring cadence, publishing through the
  # cluster spec's monitoring role (AWS requires the role whenever the
  # interval is set -- the spec's CEL guarantees it).
  monitoring_interval = each.value.monitoring_interval != 0 ? each.value.monitoring_interval : null
  monitoring_role_arn = each.value.monitoring_interval != 0 && var.spec.monitoring_role_arn != "" ? var.spec.monitoring_role_arn : null

  ca_cert_identifier = each.value.ca_cert_identifier != "" ? each.value.ca_cert_identifier : null

  tags = local.aws_tags
}
