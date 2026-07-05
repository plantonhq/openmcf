locals {
  # AWS limits launch template names to 125 characters; truncate
  # deterministically so the same manifest always yields the same name.
  launch_template_name = substr(var.metadata.name, 0, 125)

  # Resource-identity tags, matching the Pulumi module key-for-key. They are
  # applied in three places on purpose: on the template itself, and via
  # tag_specifications on the instances and volumes each launch creates -- a
  # launch template's tags do NOT propagate to what it launches, so untagged
  # fleet instances would otherwise escape cost-allocation and orphan-cleanup
  # queries.
  aws_tags = {
    "Name"                     = local.launch_template_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsLaunchTemplate"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
