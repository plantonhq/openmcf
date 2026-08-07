locals {
  # The workgroup name is metadata.name -- create-only in AWS, and the
  # basis both engines share so a manifest deploys identically on either.
  workgroup_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsRedshiftServerlessWorkgroup"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
