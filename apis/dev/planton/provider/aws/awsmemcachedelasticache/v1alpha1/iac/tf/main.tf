# The ElastiCache subnet group is a named list of subnets — pure glue
# with no independent lifecycle, so it lives inside the module. The
# referenced subnets themselves are first-class AwsSubnet nodes this
# module never modifies.
resource "aws_elasticache_subnet_group" "this" {
  count = local.manage_subnet_group ? 1 : 0

  name        = local.cluster_identifier
  description = "ElastiCache subnet group for ${local.cluster_identifier}"
  subnet_ids  = local.subnet_ids
  tags        = local.aws_tags
}

# A managed parameter group exists only when the spec carries inline
# parameters — a named parameter list is configuration owned by exactly
# this cluster, not a composable node.
resource "aws_elasticache_parameter_group" "this" {
  count = local.manage_parameter_group ? 1 : 0

  name        = local.cluster_identifier
  family      = local.parameter_group_family
  description = "Custom parameter group for ${local.cluster_identifier}"

  dynamic "parameter" {
    for_each = local.parameters
    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }

  tags = local.aws_tags
}

# The cluster composes onto its neighbors instead of embedding them:
# subnets, security groups, and parameter groups attach by reference,
# and client ingress rules live on the referenced AwsSecurityGroup
# nodes — this module never creates or mutates a resource that deserves
# to be its own node.
#
# Create-only in AWS: the cluster identifier, engine, port, subnet group,
# network_type, and node type (vertical scaling forces recreation).
# Everything else updates in place (immediately or at the next maintenance
# window, per apply_immediately).
resource "aws_elasticache_cluster" "this" {
  cluster_id      = local.cluster_identifier
  engine          = "memcached"
  engine_version  = local.engine_version
  node_type       = var.spec.node_type
  num_cache_nodes = local.num_cache_nodes
  port            = local.port

  # AZ mode
  az_mode                      = local.az_mode
  preferred_availability_zones = length(local.preferred_availability_zones) > 0 ? local.preferred_availability_zones : null

  # Encryption
  transit_encryption_enabled = local.transit_encryption_enabled

  # Networking: the subnet group managed here (or referenced), security
  # groups by reference, and optional dual-stack addressing.
  subnet_group_name  = local.effective_subnet_group_name
  security_group_ids = length(local.sg_ids) > 0 ? local.sg_ids : null
  network_type       = local.network_type
  ip_discovery       = local.ip_discovery

  # Parameter groups: the managed inline group, an existing referenced
  # group, or the engine default.
  parameter_group_name = local.effective_parameter_group_name

  # Maintenance
  maintenance_window         = local.maintenance_window
  apply_immediately          = local.apply_immediately
  auto_minor_version_upgrade = local.auto_minor_version_upgrade ? "true" : "false"

  # Notifications
  notification_topic_arn = local.notification_topic_arn

  tags = local.aws_tags
}
