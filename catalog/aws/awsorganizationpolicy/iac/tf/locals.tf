locals {
  # Resource-identity tags match the Pulumi module key-for-key. The
  # policy is the kind's one taggable surface (attachments are
  # untaggable).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsOrganizationPolicy"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Attachments key by their resolved target (each imports as
  # "{target_id}:{policy_id}").
  attachments = {
    for a in var.spec.attachments : a.target_id => a
  }
}
