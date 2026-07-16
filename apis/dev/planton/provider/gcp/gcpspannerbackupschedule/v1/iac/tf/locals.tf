locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Schedule name defaults to metadata.name when schedule_name is omitted.
  schedule_name = var.spec.schedule_name != "" ? var.spec.schedule_name : var.metadata.name

  # The spec's backup_type maps to the provider's pair of empty marker
  # blocks (full_backup_spec / incremental_backup_spec) — exactly one must
  # be present on the resource.
  is_incremental = var.spec.backup_type == "INCREMENTAL"
}
