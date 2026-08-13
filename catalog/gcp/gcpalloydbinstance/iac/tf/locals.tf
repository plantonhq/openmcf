locals {
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  instance_type     = var.spec.instance_type != "" ? var.spec.instance_type : "READ_POOL"
  availability_type = var.spec.availability_type != "" ? var.spec.availability_type : null
  display_name      = var.spec.display_name != "" ? var.spec.display_name : null

  read_pool_config = (
    var.spec.read_pool_config != null && var.spec.read_pool_config.node_count > 0
  ) ? var.spec.read_pool_config : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.spec.instance_id
    "planton-ai_kind"     = "gcpalloydbinstance"
  }

  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "planton-ai_organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "planton-ai_environment" = var.metadata.env } : {}

  id_label = (
    var.metadata.id != null && var.metadata.id != ""
  ) ? { "planton-ai_id" = var.metadata.id } : {}

  # User labels first so the platform's attribution labels win on key
  # conflicts — the catalog-wide merge order.
  labels = merge(var.spec.labels, local.base_labels, local.org_label, local.env_label, local.id_label)
}
