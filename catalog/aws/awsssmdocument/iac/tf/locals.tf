locals {
  # The AWS document name is metadata.name on both engines (document
  # names allow letters, digits, underscores, hyphens, and periods -
  # hyphenated names fit). Changing it forces replacement.
  document_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSsmDocument"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
