locals {
  # The AWS cluster name is create-time immutable -- metadata.name is the
  # naming basis both engines share so a manifest deploys identically on
  # either. AWS caps it at 40 characters. The module-managed subnet group
  # and parameter group derive their names from it too, so everything the
  # module owns is discoverable by one name.
  cluster_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsMemorydbCluster"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The folded arms create module-owned groups; the bring-your-own name
  # arms are mutually exclusive with them by CEL, so exactly one source
  # feeds each cluster argument.
  create_subnet_group    = length(var.spec.subnet_ids) > 0
  create_parameter_group = length(var.spec.parameters) > 0

  effective_subnet_group_name    = local.create_subnet_group ? aws_memorydb_subnet_group.this[0].name : (var.spec.subnet_group_name != "" ? var.spec.subnet_group_name : null)
  effective_parameter_group_name = local.create_parameter_group ? aws_memorydb_parameter_group.this[0].name : (var.spec.parameter_group_name != "" ? var.spec.parameter_group_name : null)
}
