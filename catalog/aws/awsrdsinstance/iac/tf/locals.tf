locals {
  # The instance identifier is metadata.name -- the basis both engines
  # share so a manifest deploys identically on either.
  instance_identifier = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsRdsInstance"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # A subnet group is managed here only when the spec brings raw subnets;
  # an existing group name short-circuits it. The group itself is pure
  # glue (a named list of subnets), which is why it stays inside this
  # module instead of being its own node.
  manage_subnet_group = var.spec.db_subnet_group_name == "" && length(var.spec.subnet_ids) > 0

  # The Active Directory block, when present, is one of two shapes
  # (CEL-enforced): AWS-managed (domain + role) or self-managed
  # (fqdn/ou/secret/dns_ips).
  active_directory = var.spec.active_directory

  # A DB parameter group is managed here only for inline parameters; an
  # option group only for inline options. Bringing an existing group
  # name and the inline list are mutually exclusive (CEL-enforced), and
  # both inline lists require engine + engine_version (CEL-enforced),
  # so the derivations below never see an empty string.
  manage_parameter_group = length(var.spec.parameters) > 0
  manage_option_group    = length(var.spec.options) > 0

  engine_major = split(".", var.spec.engine_version)[0]
  engine_minor = length(split(".", var.spec.engine_version)) > 1 ? tonumber(split(".", var.spec.engine_version)[1]) : 0

  # The parameter-group family, per AWS's own family naming:
  #   postgres     16.4          -> postgres16        (major)
  #   mysql        8.0.39        -> mysql8.0          (major.minor)
  #   mariadb      10.11.8       -> mariadb10.11      (major.minor)
  #   oracle-ee    19.0.0.0.ru.. -> oracle-ee-19      (engine-major)
  #   sqlserver-se 16.00.4085..  -> sqlserver-se-16.0 (engine-major.minor;
  #                                  tonumber collapses "00" to 0)
  # db2-* engines take the sqlserver arm (db2-ae-11.5).
  parameter_group_family = (
    var.spec.engine == "postgres" ? "postgres${local.engine_major}" :
    var.spec.engine == "mysql" ? "mysql${local.engine_major}.${local.engine_minor}" :
    var.spec.engine == "mariadb" ? "mariadb${local.engine_major}.${local.engine_minor}" :
    startswith(var.spec.engine, "oracle-") ? "${var.spec.engine}-${local.engine_major}" :
    "${var.spec.engine}-${local.engine_major}.${local.engine_minor}"
  )

  # The option group's major engine version keeps AWS's raw segment
  # convention (sqlserver wants "16.00", not "16.0"; oracle wants the
  # bare major).
  option_major_engine_version = (
    startswith(var.spec.engine, "oracle-") ? local.engine_major :
    join(".", slice(split(".", var.spec.engine_version), 0, min(2, length(split(".", var.spec.engine_version)))))
  )

  # The groups the instance actually uses: the managed group, an
  # explicitly referenced existing group, or null for the engine
  # default.
  effective_parameter_group = (
    local.manage_parameter_group ? aws_db_parameter_group.this[0].name :
    var.spec.parameter_group_name != "" ? var.spec.parameter_group_name :
    null
  )
  effective_option_group = (
    local.manage_option_group ? aws_db_option_group.this[0].name :
    var.spec.option_group_name != "" ? var.spec.option_group_name :
    null
  )
}
