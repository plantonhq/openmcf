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
    "planton.ai/resource-kind" = "AwsDocumentDb"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # A subnet group is managed here only when the spec brings raw subnets;
  # an existing group name short-circuits it. The group itself is pure
  # glue (a named list of subnets), which is why it stays inside this
  # module instead of being its own node.
  manage_subnet_group = var.spec.db_subnet_group_name == "" && length(var.spec.subnet_ids) > 0

  # A cluster parameter group is managed here only for inline parameters.
  # Bringing an existing group name and inline parameters are mutually
  # exclusive (CEL-enforced).
  manage_parameter_group = length(var.spec.parameters) > 0

  # The parameter-group family is derived from the pinned engine_version:
  # "5.0.0" -> docdb5.0. DocumentDB families are keyed by major.minor --
  # AWS's own family naming, not a convention of ours. The slice is
  # bounded because locals evaluate eagerly even when the family is
  # unused: an unpinned engine_version (valid whenever there are no
  # inline parameters -- CEL ties the two together) must not error here.
  engine_version_parts = split(".", var.spec.engine_version)
  engine_family        = "docdb${join(".", slice(local.engine_version_parts, 0, min(2, length(local.engine_version_parts))))}"

  # The parameter group the cluster actually uses: the managed group, an
  # explicitly referenced existing group, or null for the engine default.
  effective_cluster_parameter_group = (
    local.manage_parameter_group ? aws_docdb_cluster_parameter_group.this[0].name :
    var.spec.db_cluster_parameter_group_name != "" ? var.spec.db_cluster_parameter_group_name :
    null
  )
}
