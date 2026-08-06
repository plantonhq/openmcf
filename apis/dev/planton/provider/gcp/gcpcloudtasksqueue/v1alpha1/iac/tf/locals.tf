locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  queue_name = var.spec.queue_name
  location   = var.spec.location

  # Cloud Tasks queues have no labels surface in the API — no platform
  # attribution labels are stamped, identically on both engines.
}
