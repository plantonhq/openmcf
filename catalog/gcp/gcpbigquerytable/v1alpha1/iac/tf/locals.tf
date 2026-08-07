locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Empty-/zero-means-unset scalars translate to null so the API applies its
  # own server-side defaults.
  friendly_name   = var.spec.friendly_name != "" ? var.spec.friendly_name : null
  description     = var.spec.description != "" ? var.spec.description : null
  schema          = var.spec.schema != "" ? var.spec.schema : null
  expiration_time = var.spec.expiration_time > 0 ? var.spec.expiration_time : null
  max_staleness   = var.spec.max_staleness != "" ? var.spec.max_staleness : null
  clustering      = length(var.spec.clustering) > 0 ? var.spec.clustering : null
  kms_key_name    = var.spec.kms_key_name != "" ? var.spec.kms_key_name : null
  resource_tags   = length(var.spec.resource_tags) > 0 ? var.spec.resource_tags : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.spec.table_id
    "planton-ai_kind"     = "gcpbigquerytable"
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
