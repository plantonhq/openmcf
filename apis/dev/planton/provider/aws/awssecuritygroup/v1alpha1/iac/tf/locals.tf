locals {
  # The group name is metadata.name -- create-only in AWS (ForceNew), and the
  # basis both engines share so a manifest deploys identically on either.
  security_group_name = var.metadata.name

  # Description is required by AWS and immutable after creation.
  description = var.spec.description

  # VPC ID (already resolved from the reference before the module runs).
  vpc_id = var.spec.vpc_id

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSecurityGroup"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
