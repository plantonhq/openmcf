locals {
  # metadata.name is the consumer's cloud name (ForceNew — enhanced fan-out
  # consumers cannot be renamed).
  consumer_name = var.metadata.name

  # Resource-identity tags follow the catalog convention.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsKinesisStreamConsumer"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Resource policy — the Struct arrives from the tfvars layer as a nested
  # object; the provider wants the document as a JSON string.
  resource_policy = var.spec.resource_policy != null ? jsonencode(var.spec.resource_policy) : null
}
