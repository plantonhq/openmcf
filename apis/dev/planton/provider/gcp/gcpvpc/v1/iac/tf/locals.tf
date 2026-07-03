locals {
  # Resource-identity labels built exactly like the Pulumi module's
  # gcplabelkeys set: the "planton-ai_" prefix stands in for "planton.ai/"
  # because GCP label keys reject dots and slashes. google_compute_network
  # itself accepts no labels, so nothing consumes this map today — it is kept
  # key-for-key aligned with the Pulumi engine so any future labeled
  # sub-resource inherits an identical identity set on both engines.
  gcp_labels = merge(
    {
      "planton-ai_resource" = "true"
      "planton-ai_name"     = var.spec.network_name
      "planton-ai_kind"     = "gcpvpc"
    },
    var.metadata.org != null && var.metadata.org != "" ? { "planton-ai_organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton-ai_environment" = var.metadata.env } : {},
    var.metadata.id != null && var.metadata.id != "" ? { "planton-ai_id" = var.metadata.id } : {}
  )

  routing_mode = var.spec.routing_mode

  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project (null lets the provider resolve it; "" would be sent
  # verbatim and rejected by the API).
  project_id = var.spec.project_id != "" ? var.spec.project_id : null
}
