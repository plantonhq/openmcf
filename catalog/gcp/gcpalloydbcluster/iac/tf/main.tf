# Enable the AlloyDB API so a fresh project can host clusters.
# disable_on_destroy is false: tearing down one cluster must never disable
# the API for everything else in the project.
resource "google_project_service" "alloydb_api" {
  project = local.project_id
  service = "alloydb.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_alloydb_cluster" "cluster" {
  cluster_id = var.spec.cluster_name
  location   = var.spec.location
  project    = local.project_id
  labels     = local.labels

  # Always sent explicitly (true or false) so the spec is the single source
  # of truth for the destroy guard: the provider defaults it to true, and a
  # send-only-when-set wiring could never turn protection OFF once on.
  deletion_protection = var.spec.deletion_protection

  # AlloyDB's cluster-level destroy behavior — note the value set differs
  # from most GCP resources: DEFAULT | FORCE | PREVENT | ABANDON, where
  # FORCE also deletes any instances still in the cluster.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  skip_await_major_version_upgrade = var.spec.skip_await_major_version_upgrade

  # Dataplex Universal Catalog integration; GCP enables it by default when
  # the block is absent, so it is only rendered when the spec opts in or
  # out explicitly.
  dynamic "dataplex_config" {
    for_each = var.spec.dataplex_config != null ? [var.spec.dataplex_config] : []
    content {
      enabled = dataplex_config.value.enabled
    }
  }

  # Restore provenance — at most one source (spec-enforced), all ForceNew:
  # the restore choice exists only at create time.
  dynamic "restore_backup_source" {
    for_each = var.spec.restore_backup_source != null ? [var.spec.restore_backup_source] : []
    content {
      backup_name = restore_backup_source.value.backup_name
    }
  }

  dynamic "restore_continuous_backup_source" {
    for_each = var.spec.restore_continuous_backup_source != null ? [var.spec.restore_continuous_backup_source] : []
    content {
      cluster       = restore_continuous_backup_source.value.cluster
      point_in_time = restore_continuous_backup_source.value.point_in_time
    }
  }

  dynamic "restore_backupdr_backup_source" {
    for_each = var.spec.restore_backupdr_backup_source != null ? [var.spec.restore_backupdr_backup_source] : []
    content {
      backup = restore_backupdr_backup_source.value.backup
    }
  }

  dynamic "restore_backupdr_pitr_source" {
    for_each = var.spec.restore_backupdr_pitr_source != null ? [var.spec.restore_backupdr_pitr_source] : []
    content {
      data_source   = restore_backupdr_pitr_source.value.data_source
      point_in_time = restore_backupdr_pitr_source.value.point_in_time
    }
  }

  dynamic "network_config" {
    for_each = var.spec.network != "" ? [1] : []
    content {
      network            = var.spec.network
      allocated_ip_range = var.spec.allocated_ip_range != "" ? var.spec.allocated_ip_range : null
    }
  }

  dynamic "psc_config" {
    # try() guards the null object: HCL's && does not short-circuit.
    for_each = try(var.spec.psc_config.psc_enabled, false) ? [1] : []
    content {
      psc_enabled = true
    }
  }

  cluster_type = var.spec.cluster_type != "" ? var.spec.cluster_type : null

  dynamic "secondary_config" {
    for_each = var.spec.secondary_config != null ? [var.spec.secondary_config] : []
    content {
      primary_cluster_name = secondary_config.value.primary_cluster_name
    }
  }

  annotations       = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  subscription_type = var.spec.subscription_type != "" ? var.spec.subscription_type : null

  database_version = var.spec.database_version != "" ? var.spec.database_version : null
  display_name     = var.spec.display_name != "" ? var.spec.display_name : null

  dynamic "initial_user" {
    for_each = var.spec.initial_user != null ? [var.spec.initial_user] : []
    content {
      password = initial_user.value.password
      user     = initial_user.value.user != "" ? initial_user.value.user : null
    }
  }

  dynamic "automated_backup_policy" {
    for_each = var.spec.automated_backup_policy != null ? [var.spec.automated_backup_policy] : []
    content {
      enabled       = automated_backup_policy.value.enabled
      backup_window = automated_backup_policy.value.backup_window != "" ? automated_backup_policy.value.backup_window : null
      location      = automated_backup_policy.value.location != "" ? automated_backup_policy.value.location : null
      labels        = length(automated_backup_policy.value.labels) > 0 ? automated_backup_policy.value.labels : null

      dynamic "quantity_based_retention" {
        for_each = automated_backup_policy.value.quantity_based_retention_count > 0 ? [1] : []
        content {
          count = automated_backup_policy.value.quantity_based_retention_count
        }
      }

      dynamic "time_based_retention" {
        for_each = automated_backup_policy.value.time_based_retention_period != "" ? [1] : []
        content {
          retention_period = automated_backup_policy.value.time_based_retention_period
        }
      }

      dynamic "weekly_schedule" {
        for_each = automated_backup_policy.value.weekly_schedule != null ? [automated_backup_policy.value.weekly_schedule] : []
        content {
          days_of_week = length(weekly_schedule.value.days_of_week) > 0 ? weekly_schedule.value.days_of_week : null

          start_times {
            hours = weekly_schedule.value.start_hour
          }
        }
      }

      dynamic "encryption_config" {
        for_each = automated_backup_policy.value.encryption_kms_key_name != "" ? [1] : []
        content {
          kms_key_name = automated_backup_policy.value.encryption_kms_key_name
        }
      }
    }
  }

  dynamic "continuous_backup_config" {
    for_each = var.spec.continuous_backup_config != null ? [var.spec.continuous_backup_config] : []
    content {
      enabled              = continuous_backup_config.value.enabled
      recovery_window_days = continuous_backup_config.value.recovery_window_days > 0 ? continuous_backup_config.value.recovery_window_days : null

      dynamic "encryption_config" {
        for_each = continuous_backup_config.value.encryption_kms_key_name != "" ? [1] : []
        content {
          kms_key_name = continuous_backup_config.value.encryption_kms_key_name
        }
      }
    }
  }

  dynamic "encryption_config" {
    for_each = var.spec.kms_key_name != "" ? [1] : []
    content {
      kms_key_name = var.spec.kms_key_name
    }
  }

  dynamic "maintenance_update_policy" {
    for_each = var.spec.maintenance_window != null ? [var.spec.maintenance_window] : []
    content {
      maintenance_windows {
        day = maintenance_update_policy.value.day
        start_time {
          hours = maintenance_update_policy.value.start_hour
        }
      }
    }
  }

  depends_on = [google_project_service.alloydb_api]
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.cluster.name
  instance_id   = var.spec.primary_instance.instance_id
  instance_type = "PRIMARY"
  labels        = local.labels

  depends_on = [google_alloydb_cluster.cluster, google_project_service.alloydb_api]

  dynamic "machine_config" {
    for_each = var.spec.primary_instance.cpu_count > 0 || var.spec.primary_instance.machine_type != "" ? [1] : []
    content {
      cpu_count    = var.spec.primary_instance.cpu_count > 0 ? var.spec.primary_instance.cpu_count : null
      machine_type = var.spec.primary_instance.machine_type != "" ? var.spec.primary_instance.machine_type : null
    }
  }

  availability_type = var.spec.primary_instance.availability_type != "" ? var.spec.primary_instance.availability_type : null
  database_flags    = length(var.spec.primary_instance.database_flags) > 0 ? var.spec.primary_instance.database_flags : null
  display_name      = var.spec.primary_instance.display_name != "" ? var.spec.primary_instance.display_name : null

  dynamic "query_insights_config" {
    for_each = var.spec.primary_instance.query_insights_config != null ? [var.spec.primary_instance.query_insights_config] : []
    content {
      query_plans_per_minute  = query_insights_config.value.query_plans_per_minute
      query_string_length     = query_insights_config.value.query_string_length
      record_application_tags = query_insights_config.value.record_application_tags
      record_client_address   = query_insights_config.value.record_client_address
    }
  }

  dynamic "client_connection_config" {
    for_each = var.spec.primary_instance.require_connectors || var.spec.primary_instance.ssl_mode != "" ? [1] : []
    content {
      require_connectors = var.spec.primary_instance.require_connectors

      dynamic "ssl_config" {
        for_each = var.spec.primary_instance.ssl_mode != "" ? [1] : []
        content {
          ssl_mode = var.spec.primary_instance.ssl_mode
        }
      }
    }
  }

  # Stop/start lever: NEVER stops the primary's compute (storage and
  # configuration survive); ALWAYS restarts it.
  activation_policy = var.spec.primary_instance.activation_policy != "" ? var.spec.primary_instance.activation_policy : null

  # Client tool metadata, paired with the computed effective_annotations.
  annotations = length(var.spec.primary_instance.annotations) > 0 ? var.spec.primary_instance.annotations : null

  # ZONAL primaries only — GCP rejects it on REGIONAL instances
  # (spec-enforced pairing). Changing it live-migrates the primary.
  gce_zone = var.spec.primary_instance.gce_zone != "" ? var.spec.primary_instance.gce_zone : null

  # AlloyDB managed connection pooling (built-in pooler).
  dynamic "connection_pool_config" {
    for_each = var.spec.primary_instance.connection_pool_config != null ? [var.spec.primary_instance.connection_pool_config] : []
    content {
      enabled = connection_pool_config.value.enabled
      flags   = length(connection_pool_config.value.flags) > 0 ? connection_pool_config.value.flags : null
    }
  }

  # Public-IP / PSA-range surface on the bundled primary — the same
  # contract the standalone GcpAlloydbInstance kind models.
  dynamic "network_config" {
    for_each = var.spec.primary_instance.enable_public_ip || var.spec.primary_instance.enable_outbound_public_ip || length(var.spec.primary_instance.authorized_external_networks) > 0 || var.spec.primary_instance.allocated_ip_range_override != "" ? [1] : []
    content {
      enable_public_ip          = var.spec.primary_instance.enable_public_ip
      enable_outbound_public_ip = var.spec.primary_instance.enable_outbound_public_ip

      # Immutable: a different PSA range recreates the primary.
      allocated_ip_range_override = var.spec.primary_instance.allocated_ip_range_override != "" ? var.spec.primary_instance.allocated_ip_range_override : null

      dynamic "authorized_external_networks" {
        for_each = var.spec.primary_instance.authorized_external_networks
        content {
          cidr_range = authorized_external_networks.value.cidr_range
        }
      }
    }
  }

  # PSC on the bundled primary (meaningful only on PSC clusters).
  dynamic "psc_instance_config" {
    for_each = var.spec.primary_instance.psc_instance_config != null ? [var.spec.primary_instance.psc_instance_config] : []
    content {
      allowed_consumer_projects = length(psc_instance_config.value.allowed_consumer_projects) > 0 ? psc_instance_config.value.allowed_consumer_projects : null

      dynamic "psc_auto_connections" {
        for_each = psc_instance_config.value.psc_auto_connections
        content {
          consumer_network = psc_auto_connections.value.consumer_network != "" ? psc_auto_connections.value.consumer_network : null
          consumer_project = psc_auto_connections.value.consumer_project != "" ? psc_auto_connections.value.consumer_project : null
        }
      }

      dynamic "psc_interface_configs" {
        for_each = psc_instance_config.value.psc_interface_configs
        content {
          network_attachment_resource = psc_interface_configs.value.network_attachment_resource
        }
      }
    }
  }

  # The PRIMARY INSTANCE's own destroy behavior (the cluster's
  # deletion_policy is separate — AlloyDB gives each resource its own).
  deletion_policy = var.spec.primary_instance.deletion_policy != "" ? var.spec.primary_instance.deletion_policy : null
}
