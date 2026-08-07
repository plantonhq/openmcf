locals {
  # AWS limits Fargate profile names to 63 characters; truncate
  # deterministically so the same manifest always yields the same name.
  fargate_profile_name = substr(var.metadata.name, 0, 63)

  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = local.fargate_profile_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEksFargateProfile"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
