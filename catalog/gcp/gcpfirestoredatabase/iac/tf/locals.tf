locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  database_name = var.spec.database_name

  # Empty-string enums translate to null so the provider applies its own
  # defaults instead of receiving an empty value verbatim.
  concurrency_mode            = var.spec.concurrency_mode != "" ? var.spec.concurrency_mode : null
  pitr_enablement             = var.spec.point_in_time_recovery_enablement != "" ? var.spec.point_in_time_recovery_enablement : null
  database_edition            = var.spec.database_edition != "" ? var.spec.database_edition : null
  kms_key_name                = var.spec.kms_key_name != "" ? var.spec.kms_key_name : null
  app_engine_integration_mode = var.spec.app_engine_integration_mode != "" ? var.spec.app_engine_integration_mode : null

  firestore_data_access_mode          = var.spec.firestore_data_access_mode != "" ? var.spec.firestore_data_access_mode : null
  mongodb_compatible_data_access_mode = var.spec.mongodb_compatible_data_access_mode != "" ? var.spec.mongodb_compatible_data_access_mode : null
  realtime_updates_mode               = var.spec.realtime_updates_mode != "" ? var.spec.realtime_updates_mode : null

  # DELETE is the Planton default (the spec contract): the provider's own
  # default is ABANDON, which would leave the database running unmanaged
  # after a destroy. An empty value therefore becomes DELETE, never null.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : "DELETE"
}
