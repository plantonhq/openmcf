locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  location     = var.spec.location
  display_name = var.spec.display_name

  # An empty metadata block ({} with no contents_delta_uri and default
  # config semantics) is never sent — the index always carries config, so
  # metadata is always present. This local decides whether the
  # contents_delta_uri sub-fields ride along.
  has_contents = var.spec.contents_delta_uri != ""

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = lower(var.metadata.name)
    "planton-ai_kind"     = "gcpvertexaiindex"
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
