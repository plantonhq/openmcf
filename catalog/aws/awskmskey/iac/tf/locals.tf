locals {
  # metadata.name drives the Name identity tag -- KMS keys have no AWS name,
  # only a generated ID/ARN.
  key_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsKmsKey"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Aliases keyed by name so each materializes as its own resource and list
  # edits add/remove in place. CEL enforces name uniqueness.
  aliases = { for a in coalesce(var.spec.aliases, []) : a => a }
}
