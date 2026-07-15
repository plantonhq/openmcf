locals {
  # The custom domain's AWS identity is spec.domain_name itself (it IS the
  # resource); metadata.name only feeds the identity tags.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsHttpApiDomain"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Mappings are addressed by their path key. The empty key (the domain
  # root) is a valid, unique key -- the spec CEL guarantees uniqueness -- but
  # for_each needs non-empty map keys, so the root mapping is addressed as
  # "(root)". The alias never reaches AWS; api_mapping_key still sends the
  # real (possibly empty) value.
  api_mapping_map = {
    for m in var.spec.api_mappings : (m.api_mapping_key != "" ? m.api_mapping_key : "(root)") => m
  }
}
