locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerMlflowServer"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The server's AWS name derives from metadata.name.
  tracking_server_name = var.metadata.name
}
