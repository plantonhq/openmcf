locals {
  # The profile name is metadata.name -- the naming basis both engines
  # share. AWS allows 1-64 characters.
  profile_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockInferenceProfile"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
