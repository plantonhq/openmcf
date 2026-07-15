locals {
  # Resource-identity tags follow the catalog convention. SESv2 tags land
  # on the configuration set itself; event destinations are sub-resources
  # and carry no tags of their own.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSesConfigurationSet"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # sending_enabled is tri-state in the contract (proto optional): null
  # means "not specified", and the catalog default is TRUE -- only an
  # explicit sending_enabled=false pauses the set.
  sending_enabled = coalesce(var.spec.sending_enabled, true)

  # Event destinations are keyed by name so each maps to its own AWS
  # sub-resource and entries add/remove independently.
  event_destinations = { for d in coalesce(var.spec.event_destinations, []) : d.name => d }
}
