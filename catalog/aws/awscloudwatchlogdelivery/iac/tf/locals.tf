locals {
  # Resource-identity tags match the Pulumi module key-for-key (applied
  # to the taggable four: source, destinations, deliveries, and the
  # cross-account destination; the two policy resources are untaggable
  # at AWS). NOTE: AWS does not return tags on Get for the vended
  # family - the provider tracks them via the tagging API.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsCloudwatchLogDelivery"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # `destinations` arrives untyped (`any`): its entries carry a
  # free-form policy Struct, and a heterogeneous member collapses the
  # generated list type (the inline-policies class). Normalize each
  # entry here so the resources below read typed locals.
  destinations = var.spec.vended != null ? {
    for d in var.spec.vended.destinations : d.name => {
      name                      = d.name
      destination_resource_arn  = try(d.destination_resource_arn, "")
      delivery_destination_type = try(d.delivery_destination_type, "")
      output_format             = try(d.output_format, "")
      policy                    = try(d.policy, null)
    }
  } : {}

  # Destination policies render only for owned destinations that carry
  # one.
  destination_policies = {
    for name, d in local.destinations : name => jsonencode(d.policy) if d.policy != null
  }

  deliveries = var.spec.vended != null ? {
    for d in var.spec.vended.deliveries : d.name => d
  } : {}

  cross_account = var.spec.cross_account_destination
}
