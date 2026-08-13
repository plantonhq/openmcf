locals {
  # An access point has no name argument at all — the Name tag IS its console
  # display name, so metadata.name is the resource's only human-readable
  # identity (the same basis the Pulumi module uses).
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key -- the
  # canonical six-key identity map, no label merge (a merge here would make
  # the two engines tag the same manifest differently).
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEfsAccessPoint"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
