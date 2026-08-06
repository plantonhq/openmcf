locals {
  # AWS limits EKS cluster names to 100 characters; truncate
  # deterministically so the same manifest always yields the same name.
  cluster_name = substr(var.metadata.name, 0, 100)

  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = local.cluster_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEksCluster"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # EKS Auto Mode is all-or-nothing across three AWS blocks (compute, block
  # storage, elastic load balancing) -- the API requires them enabled or
  # disabled together, which is why the spec models one toggle. This flag
  # drives all three blocks below. The ternary (not `!= null &&`) matters:
  # HCL's && does NOT short-circuit, so the attribute access on a null
  # object would error before the null check could save it.
  auto_mode_enabled = var.spec.auto_mode != null ? var.spec.auto_mode.enabled : false
}
