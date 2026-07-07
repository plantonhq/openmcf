locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  network   = var.spec.network
  rule_name = var.spec.rule_name
  direction = var.spec.direction
  action    = var.spec.action
  rules     = var.spec.rules
  priority  = var.spec.priority
}
