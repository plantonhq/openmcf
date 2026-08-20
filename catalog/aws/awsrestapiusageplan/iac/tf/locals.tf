locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsRestApiUsagePlan"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Keys keyed by their stable entry names (the for_each keys both
  # engines share and the output map keys).
  api_keys = { for k in var.spec.api_keys : k.name => k }
}
