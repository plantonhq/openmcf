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
    "planton.ai/resource-kind" = "AwsRedshiftCluster"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # A subnet group is managed here only when the spec brings raw subnets;
  # an existing group name short-circuits it. The group itself is pure
  # glue (a named list of subnets), which is why it stays inside this
  # module instead of being its own node.
  manage_subnet_group = var.spec.cluster_subnet_group_name == "" && length(var.spec.subnet_ids) > 0

  # A parameter group is managed here only for inline parameters.
  # Bringing an existing group name and inline parameters are mutually
  # exclusive (CEL-enforced). Unlike the RDS-shaped kinds there is no
  # family derivation from an engine version -- the family is
  # redshift-1.0 unless parameter_group_family selects another.
  manage_parameter_group = length(var.spec.parameters) > 0

  # The parameter group the cluster actually uses: the managed group, an
  # explicitly referenced existing group, or null for the Redshift
  # default group of the cluster's generation.
  effective_parameter_group = (
    local.manage_parameter_group ? aws_redshift_parameter_group.this[0].name :
    var.spec.cluster_parameter_group_name != "" ? var.spec.cluster_parameter_group_name :
    null
  )

  # The subnet group name endpoint accesses fall back to when an entry
  # does not bring its own -- the managed group or the referenced
  # existing one (the subnets_or_group CEL guarantees one exists).
  effective_subnet_group = (
    local.manage_subnet_group ? aws_redshift_subnet_group.this[0].name :
    var.spec.cluster_subnet_group_name
  )

  # Usage limits key their resources (and the exported ID map) by the
  # (feature, limit type, period) triple, with an unset period rendered
  # as monthly -- the same normalization the spec's uniqueness CEL
  # applies, so validate-time uniqueness IS plan-time key uniqueness.
  usage_limits_by_key = {
    for l in var.spec.usage_limits :
    "${l.feature_type}/${l.limit_type}/${l.period != "" ? l.period : "monthly"}" => l
  }
}
