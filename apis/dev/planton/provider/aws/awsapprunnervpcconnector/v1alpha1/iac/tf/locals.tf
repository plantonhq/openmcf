locals {
  # The connector's cloud name is the resource's metadata.name -- the same
  # basis the Pulumi module uses. AWS requires connector names to be 4-40
  # characters; recreating under the same name yields a new revision.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsAppRunnerVpcConnector"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
