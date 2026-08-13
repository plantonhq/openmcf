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

  # The metric name defaults to metadata.name when the spec leaves
  # metric_name empty — the same naming basis every kind uses.
  metric_name = (
    var.spec.metric_name != null && var.spec.metric_name != ""
    ? var.spec.metric_name
    : var.metadata.name
  )
}
