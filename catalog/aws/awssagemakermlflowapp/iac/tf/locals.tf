locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerMlflowApp"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The app's AWS name derives from metadata.name (updateable in place -
  # the ARN, not the name, is the app's identity).
  app_name = var.metadata.name
}
