locals {
  group_name = var.metadata.name

  # Resource-identity tags, matching the Pulumi module key-for-key. On an
  # auto-scaling group these are emitted through the native tag blocks with
  # propagate_at_launch=true, so every launched instance carries them --
  # fleet members never escape cost-allocation and orphan-cleanup queries.
  aws_tags = {
    "Name"                     = local.group_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsAutoScalingGroup"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
