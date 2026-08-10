locals {
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # User labels first so the platform's attribution labels win on key
  # conflicts — the catalog-wide merge order. Applied to the cluster AND
  # its bundled primary instance.
  labels = merge(
    var.spec.labels,
    {
      "planton-ai_resource" = "true"
      "planton-ai_name"     = var.spec.cluster_name
      "planton-ai_kind"     = "gcpalloydbcluster"
    },
    var.metadata.org != "" ? { "planton-ai_organization" = var.metadata.org } : {},
    var.metadata.env != "" ? { "planton-ai_environment" = var.metadata.env } : {},
    var.metadata.id != "" ? { "planton-ai_id" = var.metadata.id } : {},
  )
}
