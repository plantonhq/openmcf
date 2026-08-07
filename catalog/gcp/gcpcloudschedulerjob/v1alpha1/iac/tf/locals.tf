locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Spec contract: an empty job_name falls back to metadata.name — the same
  # resolution the Pulumi module performs. An explicit conditional (never
  # coalesce()) because coalesce treats "" as a value, not as absent.
  job_name = var.spec.job_name != "" ? var.spec.job_name : var.metadata.name

  location = var.spec.location

  # Cloud Scheduler jobs have no labels surface in the API — no platform
  # attribution labels are stamped, identically on both engines.
}
