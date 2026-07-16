locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  cluster_name = var.spec.cluster_name
  region       = var.spec.region

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.cluster_name
    "planton-ai_kind"     = "gcpdataproccluster"
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
  #
  # The Dataproc API rejects user labels on virtual (GKE-based) clusters,
  # so labels are sent only for the GCE arm — both engines behave
  # identically (the spec validation already rejects user labels with the
  # virtual arm; the platform attribution set is dropped there too because
  # the API restriction applies to ALL labels on virtual clusters).
  final_labels = merge(
    var.spec.labels,
    local.base_labels,
    local.org_label,
    local.env_label,
    local.id_label,
  )
}
