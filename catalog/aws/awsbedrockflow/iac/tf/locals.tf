locals {
  # The flow name is metadata.name -- the naming basis both engines share
  # so a manifest deploys identically on either.
  flow_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockFlow"
    "planton.ai/resource-id"   = var.metadata.id
  }

  has_definition = var.spec.definition != null

  # The structural node classes whose AWS configuration union member is an
  # EMPTY object derived from the type (the Loop family has no expressible
  # member at the pinned provider and renders no configuration at all).
  empty_config_types = ["Input", "Output", "Iterator", "Collector"]
}
