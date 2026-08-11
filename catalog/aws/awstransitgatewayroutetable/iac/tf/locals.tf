locals {
  # The route table's cloud name is the resource's metadata.name, carried as
  # the Name tag -- route tables have no name attribute of their own, so the
  # tag IS the console identity. Same basis as the Pulumi module.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsTransitGatewayRouteTable"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Folded routing-domain members, each materialized as its own provider
  # resource keyed by a value that is unique within the table (the spec
  # CELs and comments enforce/record uniqueness), so adding or removing one
  # member never churns its neighbors:
  # - associations/propagations are keyed by the attachment ID itself;
  # - routes by destination CIDR; prefix list references by list ID.
  associations_map = {
    for association in var.spec.associations : association.attachment_id => association
  }
  propagations = toset(var.spec.propagations)

  routes_map = {
    for route in var.spec.routes : route.destination_cidr_block => route
  }

  prefix_list_references_map = {
    for ref in var.spec.prefix_list_references : ref.prefix_list_id => ref
  }
}
