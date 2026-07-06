# Enable the Firestore API — backup schedules are managed through the
# Firestore Admin API. disable_on_destroy stays false: tearing down one
# schedule must never disable the API for everything else in the project.
resource "google_project_service" "firestore_api" {
  project = local.project_id
  service = "firestore.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# Firestore backup schedule — periodic managed backups with retention.
# A database supports one daily and one weekly schedule; the
# daily-plus-weekly pattern is two of these resources on the same database.
# Recurrence (daily or weekly day) is immutable; retention updates in place.
# Backups already taken outlive the schedule — deleting this resource stops
# future backups but never deletes existing ones.
resource "google_firestore_backup_schedule" "this" {
  project  = local.project_id
  database = var.spec.database

  # Applies to backups created AFTER a change; existing backups keep the
  # retention they were created with.
  retention = var.spec.retention

  dynamic "daily_recurrence" {
    for_each = local.is_daily ? [1] : []
    content {}
  }

  dynamic "weekly_recurrence" {
    for_each = local.is_daily ? [] : [var.spec.weekly_recurrence]
    content {
      day = weekly_recurrence.value.day
    }
  }

  depends_on = [google_project_service.firestore_api]
}
