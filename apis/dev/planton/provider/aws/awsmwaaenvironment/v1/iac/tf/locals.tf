locals {
  # The environment name is metadata.name -- create-only in AWS (ForceNew),
  # and the basis both engines share so a manifest deploys identically on
  # either.
  environment_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsMwaaEnvironment"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
