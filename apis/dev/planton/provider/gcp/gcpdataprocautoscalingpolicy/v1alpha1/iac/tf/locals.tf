locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  policy_id = var.spec.policy_id
  location  = var.spec.location

  # The autoscaling-policy resource has no labels surface in the Dataproc
  # API — no platform attribution labels are stamped, identically on both
  # engines.
}
