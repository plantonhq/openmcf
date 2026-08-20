locals {
  # Resource-identity tags match the Pulumi module key-for-key. Only
  # the account-scoped rule resource is taggable - the organization
  # rule resources carry no tags in the provider.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsConfigRule"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
