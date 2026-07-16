locals {
  # The AWS ACL name is create-time immutable -- metadata.name is the naming
  # basis both engines share so a manifest deploys identically on either.
  # AWS caps it at 40 characters.
  acl_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsMemorydbAcl"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
