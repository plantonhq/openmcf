locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerPipeline"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The pipeline's AWS name derives from metadata.name; the display name
  # defaults to it (the provider REQUIRES a display name).
  pipeline_name = var.metadata.name
  display_name  = var.spec.display_name != "" ? var.spec.display_name : var.metadata.name
}
