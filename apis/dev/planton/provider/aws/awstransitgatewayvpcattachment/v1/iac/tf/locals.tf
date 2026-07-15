locals {
  # The attachment's cloud name is the resource's metadata.name, carried as
  # the Name tag -- attachments have no name attribute of their own, so the
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
    "planton.ai/resource-kind" = "AwsTransitGatewayVpcAttachment"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The AWS API speaks "enable"/"disable" strings for the attachment options
  # while the spec speaks booleans. dns_support and
  # security_group_referencing_support are tri-state (proto `optional bool`):
  # null falls through to the provider/AWS default instead of being pinned.
  dns_support                        = var.spec.dns_support == null ? null : (var.spec.dns_support ? "enable" : "disable")
  security_group_referencing_support = var.spec.security_group_referencing_support == null ? null : (var.spec.security_group_referencing_support ? "enable" : "disable")
}
