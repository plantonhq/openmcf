locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Empty optional strings become null so the provider applies the engine
  # default instead of sending an empty value the API would reject.
  charset   = var.spec.charset != "" ? var.spec.charset : null
  collation = var.spec.collation != "" ? var.spec.collation : null
}
