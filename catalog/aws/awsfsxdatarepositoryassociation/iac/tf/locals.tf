locals {
  # Associations have no cloud name argument — the console shows the
  # association ID; the Name tag carries metadata.name on the same basis as
  # the Pulumi module.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key (the canonical
  # six-key identity map -- user labels never merge into cloud tags).
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsFsxDataRepositoryAssociation"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The s3 sync block is only emitted when at least one policy has events —
  # an empty block and an omitted one are the same AWS state, but omitting
  # keeps the plan free of empty-block noise.
  has_s3_policies = length(var.spec.auto_import_events) > 0 || length(var.spec.auto_export_events) > 0
}
