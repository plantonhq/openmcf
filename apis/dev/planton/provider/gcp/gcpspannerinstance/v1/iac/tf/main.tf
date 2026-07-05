# Enable the Spanner API so a fresh project can host instances.
resource "google_project_service" "spanner_api" {
  project = local.project_id
  service = "spanner.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Spanner instance — the unit of compute/storage allocation that pins a
# geographic topology (config) and a capacity envelope shared by every
# database on it. name, config, and project are immutable; capacity, edition,
# and autoscaling all update in place (online, no downtime).
resource "google_spanner_instance" "this" {
  name         = local.instance_name
  project      = local.project_id
  config       = var.spec.config
  display_name = var.spec.display_name

  # Exactly one capacity method is enforced by spec validation. When all are
  # null on a PROVISIONED instance, the API defaults to 1 node.
  num_nodes        = local.num_nodes
  processing_units = local.processing_units

  instance_type                = local.instance_type
  edition                      = local.edition
  default_backup_schedule_type = local.default_backup_schedule_type

  # When true, destroy deletes all backups held on the instance first; when
  # false (the safe default), destroy fails while any backup exists.
  force_destroy = var.spec.force_destroy

  labels = local.final_labels

  dynamic "autoscaling_config" {
    for_each = var.spec.autoscaling_config != null ? [var.spec.autoscaling_config] : []
    content {
      autoscaling_limits {
        min_nodes            = autoscaling_config.value.autoscaling_limits.min_nodes > 0 ? autoscaling_config.value.autoscaling_limits.min_nodes : null
        max_nodes            = autoscaling_config.value.autoscaling_limits.max_nodes > 0 ? autoscaling_config.value.autoscaling_limits.max_nodes : null
        min_processing_units = autoscaling_config.value.autoscaling_limits.min_processing_units > 0 ? autoscaling_config.value.autoscaling_limits.min_processing_units : null
        max_processing_units = autoscaling_config.value.autoscaling_limits.max_processing_units > 0 ? autoscaling_config.value.autoscaling_limits.max_processing_units : null
      }

      dynamic "autoscaling_targets" {
        for_each = autoscaling_config.value.autoscaling_targets != null ? [autoscaling_config.value.autoscaling_targets] : []
        content {
          high_priority_cpu_utilization_percent = autoscaling_targets.value.high_priority_cpu_utilization_percent > 0 ? autoscaling_targets.value.high_priority_cpu_utilization_percent : null
          storage_utilization_percent           = autoscaling_targets.value.storage_utilization_percent > 0 ? autoscaling_targets.value.storage_utilization_percent : null
        }
      }

      # Per-replica-location node bounds for multi-region instances: a
      # read-heavy region scales independently instead of sizing every
      # region for the hottest one. The spec flattens the provider's
      # single-field replica_selection wrapper to replica_location.
      dynamic "asymmetric_autoscaling_options" {
        for_each = autoscaling_config.value.asymmetric_autoscaling_options
        content {
          replica_selection {
            location = asymmetric_autoscaling_options.value.replica_location
          }
          overrides {
            autoscaling_limits {
              min_nodes = asymmetric_autoscaling_options.value.overrides.min_nodes
              max_nodes = asymmetric_autoscaling_options.value.overrides.max_nodes
            }
          }
        }
      }
    }
  }

  depends_on = [google_project_service.spanner_api]
}
