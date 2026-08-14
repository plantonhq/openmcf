locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerEndpoint"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The endpoint's AWS name derives from metadata.name.
  endpoint_name = var.metadata.name

  # Variant names default deterministically per position ("variant-0",
  # "variant-1", ...) so re-plans never regenerate names - the Pulumi
  # module derives the identical defaults.
  production_variants = [
    for i, v in var.spec.production_variants :
    merge(v, { resolved_name = v.variant_name != "" ? v.variant_name : "variant-${i}" })
  ]
  shadow_variants = [
    for i, v in var.spec.shadow_variants :
    merge(v, { resolved_name = v.variant_name != "" ? v.variant_name : "shadow-variant-${i}" })
  ]
}
