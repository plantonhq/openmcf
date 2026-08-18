locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBudget"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Actions are keyed by their spec name - the for_each key and the
  # outputs-map key (duplicate names are spec-rejected).
  actions_by_name = { for a in var.spec.actions : a.name => a }
}
