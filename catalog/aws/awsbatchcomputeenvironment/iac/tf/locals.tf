locals {
  # Resource-identity tags follow the catalog convention. These land on the
  # compute environment itself; tags for the EC2 instances Batch launches
  # are a separate concern carried by spec.compute_resources.resource_tags.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBatchComputeEnvironment"
    "planton.ai/resource-id"   = var.metadata.id
  }

  cr = var.spec.compute_resources

  # The spec's CEL rules guarantee EC2/SPOT-only fields are absent for the
  # Fargate types; these switches exist to build the right provider payload
  # (AWS rejects EC2 knobs on Fargate requests), not to discard user intent.
  is_ec2_family = contains(["EC2", "SPOT"], local.cr.type)
  is_spot       = local.cr.type == "SPOT"
}
