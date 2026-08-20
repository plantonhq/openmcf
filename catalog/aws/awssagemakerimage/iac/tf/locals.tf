locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerImage"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The image's AWS name derives from metadata.name.
  image_name = var.metadata.name

  # Versions keyed by their stable POSITION (the append-only contract
  # taught on the spec) - the for_each keys both engines share and the
  # version_numbers output map keys.
  versions = { for i, v in var.spec.versions : tostring(i) => v }
}
