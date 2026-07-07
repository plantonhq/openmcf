locals {
  # The VPC link's cloud name is the resource's metadata.name -- the same
  # basis the Pulumi module uses. Name is the only mutable attribute on this
  # resource; subnets and security groups replace the link when changed.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsHttpApiVpcLink"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
