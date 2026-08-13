locals {
  # The cluster identifier is metadata.name — create-only in AWS, and the
  # basis both engines share so a manifest deploys identically on either.
  cluster_identifier = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsMemcachedElasticache"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # A subnet group is managed here only when the spec brings raw subnets;
  # an existing group name short-circuits it. The group itself is pure
  # glue (a named list of subnets), which is why it stays inside this
  # module instead of being its own node.
  manage_subnet_group = var.spec.subnet_group_name == "" && length(var.spec.subnet_ids) > 0

  # A parameter group is managed here only for inline parameters.
  # Bringing an existing group name and inline parameters are mutually
  # exclusive (CEL-enforced).
  manage_parameter_group = length(coalesce(try(var.spec.parameters, []), [])) > 0 && try(var.spec.parameter_group_family, "") != ""

  # Engine — empty pins nothing: AWS picks the engine's current default
  # version, so an unpinned manifest never goes stale.
  engine_version = var.spec.engine_version != "" ? var.spec.engine_version : null

  # Nodes
  num_cache_nodes = coalesce(try(var.spec.num_cache_nodes, null), 1)
  az_mode         = try(var.spec.az_mode, null) != "" ? var.spec.az_mode : null

  # Networking — references arrive pre-flattened to plain strings (the
  # generator contract lowers StringValueOrRef to string).
  subnet_ids   = coalesce(try(var.spec.subnet_ids, []), [])
  sg_ids       = coalesce(try(var.spec.security_group_ids, []), [])
  network_type = try(var.spec.network_type, null) != "" ? var.spec.network_type : null
  ip_discovery = try(var.spec.ip_discovery, null) != "" ? var.spec.ip_discovery : null

  # The subnet group the cluster actually uses: the managed group, an
  # explicitly referenced existing group, or null for the VPC default.
  effective_subnet_group_name = (
    local.manage_subnet_group ? aws_elasticache_subnet_group.this[0].name :
    var.spec.subnet_group_name != "" ? var.spec.subnet_group_name :
    null
  )

  # Encryption
  transit_encryption_enabled = coalesce(try(var.spec.transit_encryption_enabled, null), false)

  # Parameters
  parameter_group_family = try(var.spec.parameter_group_family, "")
  parameters             = coalesce(try(var.spec.parameters, []), [])

  # The parameter group the cluster actually uses: the managed group, an
  # explicitly referenced existing group, or null for the engine default.
  effective_parameter_group_name = (
    local.manage_parameter_group ? aws_elasticache_parameter_group.this[0].name :
    var.spec.parameter_group_name != "" ? var.spec.parameter_group_name :
    null
  )

  # Maintenance. auto_minor_version_upgrade is PRESENCE-typed: AWS (and
  # the provider, whose default is "true") enables minor upgrades by
  # default, so unset must be OMITTED — never forced to "false". The
  # provider types it as a nullable bool (string), hence the
  # "true"/"false" rendering when set.
  maintenance_window         = try(var.spec.maintenance_window, null) != "" ? var.spec.maintenance_window : null
  apply_immediately          = coalesce(try(var.spec.apply_immediately, null), false)
  auto_minor_version_upgrade = try(var.spec.auto_minor_version_upgrade, null) == null ? null : (var.spec.auto_minor_version_upgrade ? "true" : "false")

  # Notifications
  notification_topic_arn = try(var.spec.notification_topic_arn, "") != "" ? var.spec.notification_topic_arn : null

  # Node placement: pin ALL nodes to one AZ, or place nodes per-AZ via
  # the list (mutually exclusive, CEL-enforced).
  availability_zone            = try(var.spec.availability_zone, "") != "" ? var.spec.availability_zone : null
  preferred_availability_zones = coalesce(try(var.spec.preferred_availability_zones, []), [])

  # Port
  port = try(var.spec.port, null)
}
