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
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The cloud-side name defaults to metadata.name when the spec leaves
  # certificate_name empty — the same naming basis every kind uses.
  certificate_name = (
    var.spec.certificate_name != null && var.spec.certificate_name != ""
    ? var.spec.certificate_name
    : var.metadata.name
  )

  # Empty region means a GLOBAL certificate; a set region selects the
  # regional resource. Both branches in main.tf key off this.
  is_regional = var.spec.region != null && var.spec.region != ""
}
