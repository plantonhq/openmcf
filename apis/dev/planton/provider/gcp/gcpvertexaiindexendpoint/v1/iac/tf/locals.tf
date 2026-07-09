locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  location     = var.spec.location
  display_name = var.spec.display_name

  # The Vertex AI API expects the RELATIVE network form
  # projects/{project}/global/networks/{name} and rejects full compute
  # self-link URLs. GcpVpcNetwork references resolve to the self-link
  # (the kind's canonical output), so both literal URLs and references
  # are normalized here. Stripping is a no-op for values already in
  # relative form — identical normalization to the Pulumi module.
  network = replace(var.spec.network, "/^https://www\\.googleapis\\.com/compute/v1//", "")

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = lower(var.metadata.name)
    "planton-ai_kind"     = "gcpvertexaiindexendpoint"
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
