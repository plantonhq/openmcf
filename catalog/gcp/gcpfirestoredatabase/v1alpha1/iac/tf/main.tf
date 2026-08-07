# Enable the Firestore API — the control plane the database is managed
# through. disable_on_destroy is false: tearing down one database must
# never disable the API for everything else in the project.
resource "google_project_service" "firestore_api" {
  project = local.project_id
  service = "firestore.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Firestore database: the top-level container for collections,
# documents, and indexes. A project supports multiple named databases
# beside the special "(default)" one.
#
# location_id, database_name, database_edition, kms_key_name, and
# app_engine_integration_mode are immutable (ForceNew). type is mutable
# but switching between Native and Datastore Mode is a significant
# operational change. delete_protection_state is GCP's own guard —
# deletion by any client fails while ENABLED.
resource "google_firestore_database" "this" {
  name        = local.database_name
  location_id = var.spec.location_id
  type        = var.spec.type
  project     = local.project_id

  concurrency_mode                  = local.concurrency_mode
  point_in_time_recovery_enablement = local.pitr_enablement
  delete_protection_state           = var.spec.delete_protection_state
  database_edition                  = local.database_edition
  app_engine_integration_mode       = local.app_engine_integration_mode

  # Always DELETE so the IaC tool manages the full lifecycle. Without
  # this, the provider's default ABANDON would leave the database behind
  # on destroy — identical to the Pulumi module.
  deletion_policy = "DELETE"

  dynamic "cmek_config" {
    for_each = local.kms_key_name != null ? [local.kms_key_name] : []
    content {
      kms_key_name = cmek_config.value
    }
  }

  depends_on = [google_project_service.firestore_api]
}
