locals {
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  labels = merge(
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
