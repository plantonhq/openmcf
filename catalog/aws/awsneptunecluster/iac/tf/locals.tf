locals {
  # The cluster identifier is metadata.name -- create-only in AWS, and the
  # basis both engines share so a manifest deploys identically on either.
  cluster_identifier = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsNeptuneCluster"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # A subnet group is managed here only when the spec brings raw subnets;
  # an existing group name short-circuits it. The group itself is pure
  # glue (a named list of subnets), which is why it stays inside this
  # module instead of being its own node.
  manage_subnet_group = var.spec.neptune_subnet_group_name == "" && length(var.spec.subnet_ids) > 0

  # A cluster parameter group is managed here only for inline parameters.
  # Bringing an existing group name and inline parameters are mutually
  # exclusive (CEL-enforced).
  manage_parameter_group = length(var.spec.parameters) > 0

  # The parameter-group family is derived from the pinned engine_version:
  # "1.4.5.1" -> neptune1.4. Neptune families are keyed by major.minor --
  # AWS's own family naming, not a convention of ours. The slice is
  # bounded because locals evaluate eagerly even when the family is
  # unused: an unpinned engine_version (valid whenever there are no
  # inline parameters -- CEL ties the two together) must not error here.
  engine_version_parts = split(".", var.spec.engine_version)
  engine_family        = "neptune${join(".", slice(local.engine_version_parts, 0, min(2, length(local.engine_version_parts))))}"

  # The parameter group the cluster actually uses: the managed group, an
  # explicitly referenced existing group, or null for the engine default.
  effective_cluster_parameter_group = (
    local.manage_parameter_group ? aws_neptune_cluster_parameter_group.this[0].name :
    var.spec.neptune_cluster_parameter_group_name != "" ? var.spec.neptune_cluster_parameter_group_name :
    null
  )

  # The instance-level twin: a managed instance parameter group exists only
  # for inline instance_parameters (mutually exclusive with bringing an
  # existing name -- CEL-enforced).
  manage_instance_parameter_group = length(var.spec.instance_parameters) > 0

  # The instance parameter group folded instances adopt when they carry no
  # per-instance override: the managed group, or null for the engine
  # default. (The spec-level neptune_instance_parameter_group_name is a
  # cluster-resource argument AWS reads during major version upgrades, not
  # an instance default -- it stays on the cluster resource.)
  effective_instance_parameter_group = (
    local.manage_instance_parameter_group ? aws_neptune_parameter_group.instance[0].name : null
  )
}
