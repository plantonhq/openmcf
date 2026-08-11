# ---------------------------------------------------------------------------
# Subnet group (conditional)
# ---------------------------------------------------------------------------

# The ElastiCache subnet group is a named list of subnets -- pure glue
# with no independent lifecycle, so it lives inside the module. The
# referenced subnets themselves are first-class AwsSubnet nodes this
# module never modifies.
resource "aws_elasticache_subnet_group" "this" {
  count = local.manage_subnet_group ? 1 : 0

  name        = local.replication_group_id
  description = "ElastiCache subnet group for ${local.replication_group_id}"
  subnet_ids  = local.subnet_ids
  tags        = local.aws_tags
}

# ---------------------------------------------------------------------------
# Parameter group (conditional)
# ---------------------------------------------------------------------------

# A managed parameter group exists only when the spec carries inline
# parameters -- a named parameter list is configuration owned by exactly
# this cluster, not a composable node. Bringing an existing group name
# and inline parameters are mutually exclusive (CEL-enforced).
resource "aws_elasticache_parameter_group" "this" {
  count = local.manage_parameter_group ? 1 : 0

  name        = local.replication_group_id
  family      = var.spec.parameter_group_family
  description = "Custom parameter group for ${local.replication_group_id}"

  dynamic "parameter" {
    for_each = coalesce(try(var.spec.parameters, []), [])
    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }

  tags = local.aws_tags
}

# ---------------------------------------------------------------------------
# Replication group
# ---------------------------------------------------------------------------

# The replication group composes onto its neighbors instead of embedding
# them: subnets, security groups, KMS keys, and RBAC user groups attach
# by reference, and network ingress rules live on the referenced
# AwsSecurityGroup nodes -- this module never creates or mutates a
# resource that deserves to be its own node.
#
# Create-only in AWS: the replication_group_id, port, network_type,
# at-rest encryption + KMS key, restore sources (snapshot_arns/
# snapshot_name), global datastore membership, durability, and explicit
# shard placement (node_group_configurations). Everything else updates
# in place (immediately or at the next maintenance window, per
# apply_immediately).
resource "aws_elasticache_replication_group" "this" {
  replication_group_id = local.replication_group_id
  description          = var.spec.description

  # Engine settings are inherited from the global primary when joining
  # a global datastore -- forward only when the spec carries them.
  engine         = try(var.spec.engine, "") != "" ? var.spec.engine : null
  engine_version = try(var.spec.engine_version, "") != "" ? var.spec.engine_version : null
  node_type      = try(var.spec.node_type, "") != "" ? var.spec.node_type : null
  port           = local.port

  # Topology — non-clustered
  num_cache_clusters = local.is_clustered ? null : (local.num_cache_clusters > 0 ? local.num_cache_clusters : null)

  # Topology — clustered
  num_node_groups         = local.is_clustered ? local.num_node_groups : null
  replicas_per_node_group = local.is_clustered && try(var.spec.replicas_per_node_group, 0) > 0 ? var.spec.replicas_per_node_group : null

  preferred_cache_cluster_azs = length(coalesce(try(var.spec.preferred_cache_cluster_azs, []), [])) > 0 ? var.spec.preferred_cache_cluster_azs : null

  dynamic "node_group_configuration" {
    for_each = coalesce(try(var.spec.node_group_configurations, []), [])
    content {
      node_group_id              = try(node_group_configuration.value.node_group_id, "") != "" ? node_group_configuration.value.node_group_id : null
      primary_availability_zone  = try(node_group_configuration.value.primary_availability_zone, "") != "" ? node_group_configuration.value.primary_availability_zone : null
      replica_availability_zones = length(coalesce(try(node_group_configuration.value.replica_availability_zones, []), [])) > 0 ? node_group_configuration.value.replica_availability_zones : null
      replica_count              = try(node_group_configuration.value.replica_count, 0) > 0 ? node_group_configuration.value.replica_count : null
      slots                      = try(node_group_configuration.value.slots, "") != "" ? node_group_configuration.value.slots : null
    }
  }

  # High availability
  automatic_failover_enabled = coalesce(try(var.spec.automatic_failover_enabled, null), false)
  multi_az_enabled           = coalesce(try(var.spec.multi_az_enabled, null), false)
  durability                 = try(var.spec.durability, "") != "" ? var.spec.durability : null

  # Global datastore membership
  global_replication_group_id = try(var.spec.global_replication_group_id, "") != "" ? var.spec.global_replication_group_id : null

  # Networking: the subnet group managed here (or referenced), security
  # groups for node-level access, and optional dual-stack settings.
  subnet_group_name  = local.manage_subnet_group ? aws_elasticache_subnet_group.this[0].name : (try(var.spec.subnet_group_name, "") != "" ? var.spec.subnet_group_name : null)
  security_group_ids = length(local.sg_ids) > 0 ? local.sg_ids : null
  network_type       = try(var.spec.network_type, "") != "" ? var.spec.network_type : null
  ip_discovery       = try(var.spec.ip_discovery, "") != "" ? var.spec.ip_discovery : null

  # Encryption. Presence-typed enable flags: null omits the argument
  # entirely (AWS applies its engine default, and a global-datastore
  # secondary MUST omit them — the provider conflicts their presence
  # with global_replication_group_id).
  at_rest_encryption_enabled = local.at_rest_encryption_enabled
  transit_encryption_enabled = local.transit_encryption_enabled
  transit_encryption_mode    = local.transit_encryption_mode
  kms_key_id                 = local.kms_key_id

  # Authentication: legacy AUTH token or RBAC user groups (mutually
  # exclusive, CEL-enforced).
  auth_token                 = local.auth_token
  auth_token_update_strategy = try(var.spec.auth_token_update_strategy, "") != "" ? var.spec.auth_token_update_strategy : null
  user_group_ids             = length(local.user_group_ids) > 0 ? toset(local.user_group_ids) : null

  # Create-time restore sources (mutually exclusive, CEL-enforced)
  snapshot_arns = length(coalesce(try(var.spec.snapshot_arns, []), [])) > 0 ? var.spec.snapshot_arns : null
  snapshot_name = try(var.spec.snapshot_name, "") != "" ? var.spec.snapshot_name : null

  # Maintenance and snapshots
  maintenance_window        = local.maintenance_window
  snapshot_retention_limit  = local.snapshot_retention_limit
  snapshot_window           = local.snapshot_window
  final_snapshot_identifier = local.final_snapshot_identifier
  apply_immediately         = local.apply_immediately

  # Parameter groups: the managed inline group, an existing referenced
  # group, or the engine default.
  parameter_group_name = local.effective_parameter_group_name

  # Logging
  dynamic "log_delivery_configuration" {
    for_each = local.log_configs
    content {
      destination_type = log_delivery_configuration.value.destination_type
      destination      = log_delivery_configuration.value.destination
      log_format       = log_delivery_configuration.value.log_format
      log_type         = log_delivery_configuration.value.log_type
    }
  }

  # Advanced. auto_minor_version_upgrade omits when unset so AWS's
  # enabled-by-default stands; data_tiering renders only when true
  # (matching the Pulumi module — AWS defaults it off).
  notification_topic_arn     = local.notification_topic_arn
  auto_minor_version_upgrade = local.auto_minor_version_upgrade
  data_tiering_enabled       = local.data_tiering_enabled ? true : null
  cluster_mode               = try(var.spec.cluster_mode, "") != "" ? var.spec.cluster_mode : null

  tags = local.aws_tags
}
