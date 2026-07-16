# Enable the AlloyDB API so a fresh project can host instances.
# disable_on_destroy is false: tearing down one instance must never disable
# the API for everything else in the project.
resource "google_project_service" "alloydb_api" {
  project = local.project_id
  service = "alloydb.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# An AlloyDB instance on an existing cluster — typically a READ_POOL for
# read scaling, but PRIMARY and SECONDARY types are supported for advanced
# topologies.
resource "google_alloydb_instance" "this" {
  cluster       = var.spec.cluster
  instance_id   = var.spec.instance_id
  instance_type = local.instance_type
  labels        = local.labels

  depends_on = [google_project_service.alloydb_api]

  dynamic "machine_config" {
    for_each = var.spec.cpu_count > 0 || var.spec.machine_type != "" ? [1] : []
    content {
      cpu_count    = var.spec.cpu_count > 0 ? var.spec.cpu_count : null
      machine_type = var.spec.machine_type != "" ? var.spec.machine_type : null
    }
  }

  dynamic "read_pool_config" {
    for_each = local.read_pool_config != null ? [local.read_pool_config] : []
    content {
      node_count = read_pool_config.value.node_count
    }
  }

  availability_type = local.availability_type
  database_flags    = length(var.spec.database_flags) > 0 ? var.spec.database_flags : null
  display_name      = local.display_name

  dynamic "query_insights_config" {
    for_each = var.spec.query_insights_config != null ? [var.spec.query_insights_config] : []
    content {
      query_plans_per_minute  = query_insights_config.value.query_plans_per_minute
      query_string_length     = query_insights_config.value.query_string_length
      record_application_tags = query_insights_config.value.record_application_tags
      record_client_address   = query_insights_config.value.record_client_address
    }
  }

  dynamic "client_connection_config" {
    for_each = var.spec.require_connectors || var.spec.ssl_mode != "" ? [1] : []
    content {
      require_connectors = var.spec.require_connectors

      dynamic "ssl_config" {
        for_each = var.spec.ssl_mode != "" ? [1] : []
        content {
          ssl_mode = var.spec.ssl_mode
        }
      }
    }
  }

  # Managed connection pooling is not modeled: the released google provider
  # does not expose it for AlloyDB instances.
  activation_policy = var.spec.activation_policy != "" ? var.spec.activation_policy : null

  dynamic "network_config" {
    for_each = var.spec.enable_public_ip || var.spec.enable_outbound_public_ip || length(var.spec.authorized_external_networks) > 0 ? [1] : []
    content {
      enable_public_ip           = var.spec.enable_public_ip
      enable_outbound_public_ip  = var.spec.enable_outbound_public_ip

      dynamic "authorized_external_networks" {
        for_each = var.spec.authorized_external_networks
        content {
          cidr_range = authorized_external_networks.value.cidr_range
        }
      }
    }
  }

  dynamic "psc_instance_config" {
    for_each = var.spec.psc_instance_config != null ? [var.spec.psc_instance_config] : []
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
}
