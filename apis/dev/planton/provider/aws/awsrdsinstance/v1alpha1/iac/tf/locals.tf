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
}
