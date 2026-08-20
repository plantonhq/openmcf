locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockAgentCoreTools"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Arms keyed by their stable entry names (the for_each keys both
  # engines share and the output map keys).
  browsers          = { for b in var.spec.browsers : b.name => b }
  browser_profiles  = { for p in var.spec.browser_profiles : p.name => p }
  code_interpreters = { for c in var.spec.code_interpreters : c.name => c }
}
