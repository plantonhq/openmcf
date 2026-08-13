locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsRestApiDomain"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The resolved endpoint type: the spec's choice, or REGIONAL when
  # endpoint_configuration is omitted (the right default for almost
  # every new domain). Certificate fan-in keys off it.
  endpoint_type = var.spec.endpoint_configuration != null && try(var.spec.endpoint_configuration.type, "") != "" ? var.spec.endpoint_configuration.type : "REGIONAL"

  # Mappings keyed by base path ("(root)" for the empty path) -- the
  # same keys as the Pulumi loop and the output map.
  base_path_mappings = { for m in var.spec.base_path_mappings : (m.base_path != "" ? m.base_path : "(root)") => m }

  # Associations keyed by the VPC endpoint they grant.
  access_associations = { for a in var.spec.access_associations : a.vpc_endpoint_id => a }
}
