locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerFeatureGroup"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The group's AWS name derives from metadata.name.
  feature_group_name = var.metadata.name

  has_online  = var.spec.online_store != null
  has_offline = var.spec.offline_store != null
}
