locals {
  # AWS limits node group names to 63 characters; truncate deterministically
  # so the same manifest always yields the same name.
  node_group_name = substr(var.metadata.name, 0, 63)

  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = local.node_group_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEksNodeGroup"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The proto enum arrives as its lowercase value name (pruned entirely for
  # the zero value on_demand); AWS expects the uppercase API constant.
  capacity_type = {
    ""               = "ON_DEMAND"
    "on_demand"      = "ON_DEMAND"
    "spot"           = "SPOT"
    "capacity_block" = "CAPACITY_BLOCK"
  }[var.spec.capacity_type]
}
