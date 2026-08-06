locals {
  # The gateway's cloud name is the resource's metadata.name, carried as the
  # Name tag -- Transit Gateways have no name attribute of their own, so the
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
    "planton.ai/resource-kind" = "AwsTransitGateway"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The AWS API speaks "enable"/"disable" strings for the gateway dials while
  # the spec speaks booleans. Tri-state dials (proto `optional bool`) map
  # null -> null so an omitted dial falls through to the provider/AWS default
  # instead of being pinned; plain-bool dials always send their value.
  default_route_table_association = var.spec.default_route_table_association == null ? null : (var.spec.default_route_table_association ? "enable" : "disable")
  default_route_table_propagation = var.spec.default_route_table_propagation == null ? null : (var.spec.default_route_table_propagation ? "enable" : "disable")
  dns_support                     = var.spec.dns_support == null ? null : (var.spec.dns_support ? "enable" : "disable")
  vpn_ecmp_support                = var.spec.vpn_ecmp_support == null ? null : (var.spec.vpn_ecmp_support ? "enable" : "disable")
  encryption_support              = var.spec.encryption_support == null ? null : (var.spec.encryption_support ? "enable" : "disable")
}
