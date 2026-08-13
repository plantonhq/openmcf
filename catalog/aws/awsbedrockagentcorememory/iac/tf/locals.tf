locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockAgentCoreMemory"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Optional single-entry surfaces render only when declared.
  has_kinesis = var.spec.kinesis_delivery != null

  # Strategies keyed by their stable entry names (the for_each keys both
  # engines share and the output map keys).
  strategies = { for s in var.spec.strategies : s.name => s }
}
