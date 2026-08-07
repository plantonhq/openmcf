# Enable the Spanner API so a fresh project can host backup schedules.
resource "google_project_service" "spanner_api" {
  project = local.project_id
  service = "spanner.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Spanner backup schedule — creates backups of one database on a cron
# cadence and retains each for retention_duration. A database commonly
# carries a daily incremental schedule AND a weekly full schedule side by
# side. name, instance, database, and backup type are immutable; cron,
# retention, and encryption update in place.
resource "google_spanner_backup_schedule" "this" {
  name     = local.schedule_name
  project  = local.project_id
  instance = var.spec.instance
  database = var.spec.database

  # Applies to backups created AFTER a change; existing backups keep the
  # retention they were created with.
  retention_duration = var.spec.retention_duration

  spec {
    cron_spec {
      # Evaluated in UTC. Spanner accepts a bounded set of frequencies:
      # every 12 hours, daily, weekly, or monthly.
      text = var.spec.cron
    }
  }

  # Exactly one backup kind, expressed by the provider as a pair of empty
  # marker blocks. INCREMENTAL chains store only changes since the previous
  # backup (cheaper storage, same restore semantics) and require the
  # instance to be ENTERPRISE or ENTERPRISE_PLUS edition.
  dynamic "full_backup_spec" {
    for_each = local.is_incremental ? [] : [1]
    content {}
  }

  dynamic "incremental_backup_spec" {
    for_each = local.is_incremental ? [1] : []
    content {}
  }

  # If omitted, backups use USE_DATABASE_ENCRYPTION (inherit the database's
  # posture). CMEK requires exactly one key shape — enforced pre-deploy by
  # spec validation.
  dynamic "encryption_config" {
    for_each = var.spec.encryption_config != null ? [var.spec.encryption_config] : []
    content {
      encryption_type = encryption_config.value.encryption_type
      kms_key_name    = encryption_config.value.kms_key_name != "" ? encryption_config.value.kms_key_name : null
      kms_key_names   = length(encryption_config.value.kms_key_names) > 0 ? encryption_config.value.kms_key_names : null
    }
  }

  depends_on = [google_project_service.spanner_api]
}
