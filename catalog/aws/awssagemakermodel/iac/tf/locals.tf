locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerModel"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The model's AWS name derives from metadata.name (charset-compatible:
  # letters, digits, hyphens, <= 63).
  model_name = var.metadata.name

  # Exactly one of the two container forms is set (spec-validated); the
  # pipeline renders as repeated container blocks.
  pipeline_containers = var.spec.containers
}
