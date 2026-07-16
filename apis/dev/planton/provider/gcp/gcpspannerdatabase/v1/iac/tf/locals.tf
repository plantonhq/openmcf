locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Database name defaults to metadata.name when database_name is omitted.
  database_name = var.spec.database_name != "" ? var.spec.database_name : var.metadata.name

  database_dialect         = var.spec.database_dialect != "" ? var.spec.database_dialect : null
  version_retention_period = var.spec.version_retention_period != "" ? var.spec.version_retention_period : null
  default_time_zone        = var.spec.default_time_zone != "" ? var.spec.default_time_zone : null
}
