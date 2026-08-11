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

  # Client-side destroy behavior: DELETE (default), PREVENT (destroy
  # fails), or ABANDON (drop from state, keep the instance running).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

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
          total_cpu_utilization_percent         = autoscaling_targets.value.total_cpu_utilization_percent > 0 ? autoscaling_targets.value.total_cpu_utilization_percent : null
        }
      }

      # Per-replica autoscaling tuning for multi-region instances: a
      # read-heavy region scales independently instead of sizing every
      # region for the hottest one. The spec flattens the provider's
      # single-field replica_selection wrapper to replica_location and
      # the overrides' autoscaling_limits wrapper onto the overrides
      # message; the limits block is emitted only when a bounds family
      # (nodes or processing units) is actually set, so a targets-only
      # override never sends an empty block.
      dynamic "asymmetric_autoscaling_options" {
        for_each = autoscaling_config.value.asymmetric_autoscaling_options
        content {
          replica_selection {
            location = asymmetric_autoscaling_options.value.replica_location
          }
          overrides {
            dynamic "autoscaling_limits" {
              for_each = (asymmetric_autoscaling_options.value.overrides.min_nodes > 0 || asymmetric_autoscaling_options.value.overrides.max_nodes > 0 || asymmetric_autoscaling_options.value.overrides.min_processing_units > 0 || asymmetric_autoscaling_options.value.overrides.max_processing_units > 0) ? [asymmetric_autoscaling_options.value.overrides] : []
              content {
                min_nodes            = autoscaling_limits.value.min_nodes > 0 ? autoscaling_limits.value.min_nodes : null
                max_nodes            = autoscaling_limits.value.max_nodes > 0 ? autoscaling_limits.value.max_nodes : null
                min_processing_units = autoscaling_limits.value.min_processing_units > 0 ? autoscaling_limits.value.min_processing_units : null
                max_processing_units = autoscaling_limits.value.max_processing_units > 0 ? autoscaling_limits.value.max_processing_units : null
              }
            }

            autoscaling_target_high_priority_cpu_utilization_percent = asymmetric_autoscaling_options.value.overrides.autoscaling_target_high_priority_cpu_utilization_percent > 0 ? asymmetric_autoscaling_options.value.overrides.autoscaling_target_high_priority_cpu_utilization_percent : null
            autoscaling_target_total_cpu_utilization_percent         = asymmetric_autoscaling_options.value.overrides.autoscaling_target_total_cpu_utilization_percent > 0 ? asymmetric_autoscaling_options.value.overrides.autoscaling_target_total_cpu_utilization_percent : null
            disable_high_priority_cpu_autoscaling                    = asymmetric_autoscaling_options.value.overrides.disable_high_priority_cpu_autoscaling ? true : null
            disable_total_cpu_autoscaling                            = asymmetric_autoscaling_options.value.overrides.disable_total_cpu_autoscaling ? true : null
          }
        }
      }
    }
  }

  depends_on = [google_project_service.spanner_api]
}
