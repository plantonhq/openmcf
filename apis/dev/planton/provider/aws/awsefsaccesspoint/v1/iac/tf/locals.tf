locals {
  # An access point has no name argument at all — the Name tag IS its console
  # display name, so metadata.name is the resource's only human-readable
  # identity (the same basis the Pulumi module uses).
  resource_name = var.metadata.name

  # Resource-identity tags follow the catalog convention; user labels merge in
  # without being able to override the identity keys.
  aws_tags = merge(try(var.metadata.labels, {}), {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEfsAccessPoint"
    "planton.ai/resource-id"   = var.metadata.id
  })
}
