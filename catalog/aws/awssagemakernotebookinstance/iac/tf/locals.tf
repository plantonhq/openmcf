locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerNotebookInstance"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The instance's AWS name derives from metadata.name; the folded
  # lifecycle configuration rides a stable derived name.
  notebook_name         = var.metadata.name
  lifecycle_config_name = "${var.metadata.name}-lifecycle"

  has_lifecycle = var.spec.lifecycle_config != null
}
