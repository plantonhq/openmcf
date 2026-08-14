locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerModelRegistry"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The group's AWS name derives from metadata.name.
  group_name = var.metadata.name

  has_policy = var.spec.resource_policy != null
}
