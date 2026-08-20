locals {
  # Resource-identity tags match the Pulumi module key-for-key. They
  # land on the OWNED GROUP only - the schedule itself is untaggable
  # at AWS (the deliberate tag-convention absence).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEventBridgeScheduler"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
