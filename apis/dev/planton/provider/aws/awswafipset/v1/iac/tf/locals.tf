locals {
  # Resource-identity tags follow the catalog convention. The Name tag is
  # what the WAF console displays alongside the set.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsWafIpSet"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Null-when-unset so the provider omits the argument instead of sending an
  # empty string AWS would reject.
  description = var.spec.description != "" ? var.spec.description : null
}
