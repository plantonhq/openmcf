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
    "planton.ai/resource-kind" = "AwsRdsCluster"
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

  # The parameter-group family is derived from engine + engine_version
  # (inline parameters require a pinned version, CEL-enforced, so the
  # derivation never sees an empty string):
  #   aurora-postgresql 16.4                -> aurora-postgresql16
  #   postgres          16.4                -> postgres16
  #   aurora-mysql      8.0.mysql_aurora... -> aurora-mysql8.0
  #   mysql             8.0.39              -> mysql8.0
  # PostgreSQL families are keyed by major version; MySQL families by
  # major.minor -- AWS's own family naming, not a convention of ours.
  engine_family = (
    var.spec.engine == "aurora-postgresql" ? "aurora-postgresql${split(".", var.spec.engine_version)[0]}" :
    var.spec.engine == "postgres" ? "postgres${split(".", var.spec.engine_version)[0]}" :
    var.spec.engine == "aurora-mysql" ? "aurora-mysql${join(".", slice(split(".", var.spec.engine_version), 0, 2))}" :
    "mysql${join(".", slice(split(".", var.spec.engine_version), 0, 2))}"
  )

  # The parameter group the cluster actually uses: the managed group, an
  # explicitly referenced existing group, or null for the engine default.
  effective_cluster_parameter_group = (
    local.manage_parameter_group ? aws_rds_cluster_parameter_group.this[0].name :
    var.spec.db_cluster_parameter_group_name != "" ? var.spec.db_cluster_parameter_group_name :
    null
  )
}
