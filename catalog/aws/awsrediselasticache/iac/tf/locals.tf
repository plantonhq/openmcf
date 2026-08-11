locals {
  # The replication group identifier is metadata.name -- create-only in
  # AWS, and the basis both engines share so a manifest deploys
  # identically on either.
  replication_group_id = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsRedisElasticache"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # References arrive pre-flattened to plain strings (the generator
  # contract lowers StringValueOrRef to string; the platform resolves
  # valueFrom before the module runs).
  subnet_ids = coalesce(try(var.spec.subnet_ids, []), [])

  # A subnet group is managed here only when the spec brings raw subnets;
  # an existing group name short-circuits it. The group itself is pure
  # glue (a named list of subnets), which is why it stays inside this
  # module instead of being its own node.
  manage_subnet_group = try(var.spec.subnet_group_name, "") == "" && length(local.subnet_ids) > 0

  # A parameter group is managed here only for inline parameters.
  # Bringing an existing group name and inline parameters are mutually
  # exclusive (CEL-enforced).
  manage_parameter_group = length(coalesce(try(var.spec.parameters, []), [])) > 0 && try(var.spec.parameter_group_name, "") == ""

  # The parameter group the cluster actually uses: the managed group, an
  # explicitly referenced existing group, or null for the engine default.
  effective_parameter_group_name = (
    local.manage_parameter_group ? aws_elasticache_parameter_group.this[0].name :
    try(var.spec.parameter_group_name, "") != "" ? var.spec.parameter_group_name :
    null
  )

  # Topology
  num_cache_clusters = try(var.spec.num_cache_clusters, 0)
  num_node_groups    = try(var.spec.num_node_groups, 0)
  is_clustered       = local.num_node_groups > 0

  # Networking
  sg_ids = coalesce(try(var.spec.security_group_ids, []), [])

  # Encryption. The two enable flags are PRESENCE-typed (optional bool):
  # unset means "let AWS apply its engine default" and must be OMITTED —
  # the provider rejects the arguments' mere presence alongside
  # global_replication_group_id, so rendering an explicit false would
  # break the global-datastore join. The provider types them as nullable
  # bools (strings), hence the "true"/"false" rendering when set.
  at_rest_encryption_enabled = try(var.spec.at_rest_encryption_enabled, null) == null ? null : (var.spec.at_rest_encryption_enabled ? "true" : "false")
  transit_encryption_enabled = try(var.spec.transit_encryption_enabled, null)
  transit_encryption_mode    = try(var.spec.transit_encryption_mode, "") != "" ? var.spec.transit_encryption_mode : null
  kms_key_id                 = try(var.spec.kms_key_id, "") != "" ? var.spec.kms_key_id : null

  # Authentication: legacy AUTH token or RBAC user groups (mutually
  # exclusive, CEL-enforced). Empty string means "not set".
  auth_token     = try(var.spec.auth_token, "") != "" ? var.spec.auth_token : null
  user_group_ids = [for u in coalesce(try(var.spec.user_group_ids, []), []) : u if u != ""]

  # Maintenance and snapshots
  maintenance_window        = try(var.spec.maintenance_window, "") != "" ? var.spec.maintenance_window : null
  snapshot_retention_limit  = try(var.spec.snapshot_retention_limit, 0) > 0 ? var.spec.snapshot_retention_limit : null
  snapshot_window           = try(var.spec.snapshot_window, "") != "" ? var.spec.snapshot_window : null
  final_snapshot_identifier = try(var.spec.final_snapshot_identifier, "") != "" ? var.spec.final_snapshot_identifier : null
  apply_immediately         = coalesce(try(var.spec.apply_immediately, null), false)

  # Logging
  log_configs = coalesce(try(var.spec.log_delivery_configurations, []), [])

  # Advanced. auto_minor_version_upgrade is PRESENCE-typed: AWS enables
  # minor upgrades by default, so unset must be OMITTED (AWS decides),
  # never forced to "false". The provider types it as a nullable bool
  # (string), hence the "true"/"false" rendering when set.
  notification_topic_arn     = try(var.spec.notification_topic_arn, "") != "" ? var.spec.notification_topic_arn : null
  auto_minor_version_upgrade = try(var.spec.auto_minor_version_upgrade, null) == null ? null : (var.spec.auto_minor_version_upgrade ? "true" : "false")
  data_tiering_enabled       = coalesce(try(var.spec.data_tiering_enabled, null), false)
  port                       = try(var.spec.port, null)
}
