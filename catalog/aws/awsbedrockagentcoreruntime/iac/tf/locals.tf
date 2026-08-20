locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockAgentCoreRuntime"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Optional single-entry surfaces render only when declared.
  has_code_artifact = var.spec.artifact.code != null
  has_lifecycle     = var.spec.lifecycle != null
  has_jwt           = var.spec.custom_jwt_authorizer != null

  # Endpoints keyed by their stable entry names (the for_each keys both
  # engines share and the output map keys).
  endpoints = { for e in var.spec.endpoints : e.name => e }
}
