locals {
  # metadata.name is the bus's cloud name. AWS reserves the name "default"
  # for the account's built-in bus — creating another is rejected, so the
  # module fails fast at plan time (see the precondition in main.tf) instead
  # of surfacing an opaque API error at apply.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEventBridgeBus"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Partner event source — null when not configured.
  event_source_name = var.spec.event_source_name != "" ? var.spec.event_source_name : null

  # The resource policy arrives from the tfvars layer as a nested object (the
  # spec models it as a Struct); the provider wants a JSON string.
  resource_policy = var.spec.resource_policy != null ? jsonencode(var.spec.resource_policy) : null
}
