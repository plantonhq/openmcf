locals {
  # The AWS user group id is create-time immutable -- metadata.name is the
  # naming basis both engines share so a manifest deploys identically on either.
  user_group_id = var.metadata.name

  # Membership refs arrive pre-flattened to plain user ids (the generator
  # contract lowers StringValueOrRef to string; the platform resolves
  # valueFrom before the module runs).
  user_ids = coalesce(try(var.spec.user_ids, []), [])

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsElasticacheUserGroup"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
