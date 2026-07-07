locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  region      = var.spec.region
  environment = var.spec.environment
  secret_name = var.spec.secret_name

  # Kubernetes Secrets carry no GCP labels surface — no platform
  # attribution labels are stamped, identically on both engines.
}
