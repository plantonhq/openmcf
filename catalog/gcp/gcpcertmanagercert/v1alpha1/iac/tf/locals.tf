locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project (ambient credentials decide).
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # cert_name falls back to metadata.name — explicit conditional, so both
  # engines derive the identical cloud-side name.
  cert_name = var.spec.cert_name != "" ? var.spec.cert_name : var.metadata.name

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.cert_name
    "planton-ai_kind"     = "gcpcertmanagercert"
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

  # User labels first so platform attribution labels win on key conflicts —
  # identical merge order to the Pulumi module.
  final_labels = merge(
    var.spec.labels,
    local.base_labels,
    local.org_label,
    local.env_label,
    local.id_label,
  )
}
