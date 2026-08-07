locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API. This project only scopes the in-module
  # API enablement — the connection itself is addressed by the network.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Empty service falls through to the Google managed-services producer —
  # the one behind Cloud SQL, AlloyDB, Memorystore, and Filestore private IP.
  service = var.spec.service != "" ? var.spec.service : "servicenetworking.googleapis.com"
}
