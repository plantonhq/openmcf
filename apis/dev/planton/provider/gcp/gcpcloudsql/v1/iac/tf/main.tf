# Enable the Cloud SQL Admin API so a fresh project can host the instance.
# disable_on_destroy is false: tearing down one instance must never disable
# the API for everything else in the project.
resource "google_project_service" "sqladmin_api" {
  project = local.project_id
  service = "sqladmin.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Cloud SQL instance — a managed MySQL, PostgreSQL, or SQL Server server.
# One resource is one instance: a primary, or (with master_instance_name) a
# read replica. Databases and users inside it are separate resources
# (google_sql_database / google_sql_user), composed by instance name.
#
# Lifecycle notes the API enforces:
#   - name/region/CMEK/disk type are immutable; deleted names stay reserved
#     for ~1 week.
#   - database_version upgrades and tier/edition changes apply IN PLACE
#     (with a restart) — no replacement.
#   - disks grow but never shrink; private_network can be set or changed but
#     never removed in place.
#   - a private-IP instance requires the VPC to already carry a service
#     networking connection; the provider prechecks and fails fast because a
#     failed create still burns the reserved instance name.
resource "google_sql_database_instance" "this" {
  name             = var.spec.instance_name
  project          = local.project_id
  region           = var.spec.region
  database_version = var.spec.database_version

  # Set only on read replicas: the primary this instance replicates from.
  master_instance_name = local.master_instance_name

  # CMEK: the key must live in the instance's region. Immutable.
  encryption_key_name = local.encryption_key_name

  # Write-only in GCP — never readable back, never in outputs. Required for
  # SQL Server (spec CEL enforces pre-deploy).
  root_password = local.root_password

  # Engine-side destroy guard: refuses `destroy` at plan time while true.
  # Distinct from settings.deletion_protection_enabled (the API-side guard).
  deletion_protection = var.spec.deletion_protection

  settings {
    tier              = var.spec.tier
    edition           = var.spec.edition
    availability_type = var.spec.availability_type
    activation_policy = var.spec.activation_policy

    disk_type             = local.disk_type
    disk_size             = local.disk_size_gb
    disk_autoresize       = local.disk_auto_resize
    disk_autoresize_limit = local.disk_auto_resize_limit

    # API-side delete guard — blocks deletion from console/gcloud/API too.
    deletion_protection_enabled = var.spec.deletion_protection_enabled

    # Keep automated backups (and PITR logs) after the instance is deleted.
    retain_backups_on_delete = var.spec.retain_backups_on_delete

    connector_enforcement = local.connector_enforcement

    # SQL Server-only knobs (spec CEL restricts them to SQLSERVER engines).
    time_zone = local.time_zone
    collation = local.collation

    enable_google_ml_integration = var.spec.enable_google_ml_integration
    enable_dataplex_integration  = var.spec.enable_dataplex_integration

    user_labels = local.final_labels

    # Emitted only when enabled: the API rejects a data-cache stanza on
    # ENTERPRISE instances (spec CEL already forces ENTERPRISE_PLUS).
    dynamic "data_cache_config" {
      for_each = local.data_cache_enabled ? [1] : []
      content {
        data_cache_enabled = true
      }
    }

    # SQL Server licensing/performance dial.
    dynamic "advanced_machine_features" {
      for_each = var.spec.threads_per_core != null ? [1] : []
      content {
        threads_per_core = var.spec.threads_per_core
      }
    }

    ip_configuration {
      # An omitted spec network block resolves to public IPv4 with no
      # authorized networks — Auth Proxy / connector access only.
      ipv4_enabled    = local.ipv4_enabled
      private_network = local.private_network

      allocated_ip_range                            = local.allocated_ip_range
      enable_private_path_for_google_cloud_services = local.enable_private_path

      ssl_mode       = local.ssl_mode
      server_ca_mode = local.server_ca_mode
      server_ca_pool = local.server_ca_pool

      custom_subject_alternative_names = local.custom_sans

      dynamic "authorized_networks" {
        for_each = local.authorized_networks
        content {
          value           = authorized_networks.value.value
          name            = authorized_networks.value.name != "" ? authorized_networks.value.name : null
          expiration_time = authorized_networks.value.expiration_time != "" ? authorized_networks.value.expiration_time : null
        }
      }

      dynamic "psc_config" {
        for_each = local.psc != null ? [local.psc] : []
        content {
          psc_enabled               = true
          allowed_consumer_projects = psc_config.value.allowed_consumer_projects
          network_attachment_uri    = psc_config.value.network_attachment_uri != "" ? psc_config.value.network_attachment_uri : null

          dynamic "psc_auto_connections" {
            for_each = psc_config.value.auto_connections
            content {
              consumer_network            = psc_auto_connections.value.consumer_network
              consumer_service_project_id = psc_auto_connections.value.consumer_service_project_id != "" ? psc_auto_connections.value.consumer_service_project_id : null
            }
          }
        }
      }
    }

    dynamic "location_preference" {
      for_each = local.location_preference != null ? [local.location_preference] : []
      content {
        zone           = location_preference.value.zone != "" ? location_preference.value.zone : null
        secondary_zone = location_preference.value.secondary_zone != "" ? location_preference.value.secondary_zone : null
      }
    }

    dynamic "backup_configuration" {
      for_each = var.spec.backup != null ? [var.spec.backup] : []
      content {
        enabled    = backup_configuration.value.enabled
        start_time = backup_configuration.value.start_time != "" ? backup_configuration.value.start_time : null
        location   = backup_configuration.value.location != "" ? backup_configuration.value.location : null

        # MySQL PITR mechanism; also required for MySQL replicas and HA.
        binary_log_enabled = backup_configuration.value.binary_log_enabled

        # PostgreSQL / SQL Server PITR mechanism.
        point_in_time_recovery_enabled = backup_configuration.value.point_in_time_recovery_enabled

        transaction_log_retention_days = backup_configuration.value.transaction_log_retention_days

        dynamic "backup_retention_settings" {
          for_each = backup_configuration.value.retained_backups != null ? [1] : []
          content {
            retained_backups = backup_configuration.value.retained_backups
            retention_unit   = "COUNT"
          }
        }
      }
    }

    dynamic "maintenance_window" {
      for_each = var.spec.maintenance_window != null ? [var.spec.maintenance_window] : []
      content {
        day          = maintenance_window.value.day
        hour         = maintenance_window.value.hour
        update_track = maintenance_window.value.update_track != "" ? maintenance_window.value.update_track : null
      }
    }

    dynamic "deny_maintenance_period" {
      for_each = var.spec.deny_maintenance_period != null ? [var.spec.deny_maintenance_period] : []
      content {
        start_date = deny_maintenance_period.value.start_date
        end_date   = deny_maintenance_period.value.end_date
        time       = deny_maintenance_period.value.time
      }
    }

    dynamic "insights_config" {
      # try() because HCL's && does not short-circuit on the nullable block.
      for_each = try(var.spec.insights_config.query_insights_enabled, false) ? [var.spec.insights_config] : []
      content {
        query_insights_enabled  = true
        query_string_length     = insights_config.value.query_string_length
        record_application_tags = insights_config.value.record_application_tags
        record_client_address   = insights_config.value.record_client_address
        query_plans_per_minute  = insights_config.value.query_plans_per_minute
      }
    }

    dynamic "password_validation_policy" {
      for_each = var.spec.password_validation_policy != null ? [var.spec.password_validation_policy] : []
      content {
        enable_password_policy      = password_validation_policy.value.enable_password_policy
        min_length                  = password_validation_policy.value.min_length
        complexity                  = password_validation_policy.value.complexity != "" ? password_validation_policy.value.complexity : null
        reuse_interval              = password_validation_policy.value.reuse_interval
        disallow_username_substring = password_validation_policy.value.disallow_username_substring
        password_change_interval    = password_validation_policy.value.password_change_interval != "" ? password_validation_policy.value.password_change_interval : null
      }
    }

    dynamic "connection_pool_config" {
      # try() because HCL's && does not short-circuit on the nullable block.
      for_each = try(var.spec.connection_pooling.enabled, false) ? [var.spec.connection_pooling] : []
      content {
        connection_pooling_enabled = true

        dynamic "flags" {
          for_each = local.connection_pool_flags_list
          content {
            name  = flags.value.name
            value = flags.value.value
          }
        }
      }
    }

    dynamic "sql_server_audit_config" {
      for_each = var.spec.sql_server_audit_config != null ? [var.spec.sql_server_audit_config] : []
      content {
        bucket             = sql_server_audit_config.value.bucket != "" ? sql_server_audit_config.value.bucket : null
        retention_interval = sql_server_audit_config.value.retention_interval != "" ? sql_server_audit_config.value.retention_interval : null
        upload_interval    = sql_server_audit_config.value.upload_interval != "" ? sql_server_audit_config.value.upload_interval : null
      }
    }

    dynamic "active_directory_config" {
      for_each = var.spec.active_directory_domain != "" ? [1] : []
      content {
        domain = var.spec.active_directory_domain
      }
    }

    dynamic "database_flags" {
      for_each = local.database_flags_list
      content {
        name  = database_flags.value.name
        value = database_flags.value.value
      }
    }
  }

  # Replica behavior + external-source replication channel. Only meaningful
  # with master_instance_name (spec CEL enforces the pairing).
  dynamic "replica_configuration" {
    for_each = var.spec.replica_configuration != null ? [var.spec.replica_configuration] : []
    content {
      failover_target           = replica_configuration.value.failover_target
      cascadable_replica        = replica_configuration.value.cascadable_replica
      username                  = replica_configuration.value.username != "" ? replica_configuration.value.username : null
      password                  = replica_configuration.value.password != "" ? replica_configuration.value.password : null
      ca_certificate            = replica_configuration.value.ca_certificate != "" ? replica_configuration.value.ca_certificate : null
      client_certificate        = replica_configuration.value.client_certificate != "" ? replica_configuration.value.client_certificate : null
      client_key                = replica_configuration.value.client_key != "" ? replica_configuration.value.client_key : null
      dump_file_path            = replica_configuration.value.dump_file_path != "" ? replica_configuration.value.dump_file_path : null
      connect_retry_interval    = replica_configuration.value.connect_retry_interval
      master_heartbeat_period   = replica_configuration.value.master_heartbeat_period
      ssl_cipher                = replica_configuration.value.ssl_cipher != "" ? replica_configuration.value.ssl_cipher : null
      verify_server_certificate = replica_configuration.value.verify_server_certificate
    }
  }

  depends_on = [google_project_service.sqladmin_api]
}
