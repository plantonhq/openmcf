# Enable the Kubernetes Engine API so a fresh project can host the cluster.
# disable_on_destroy is false: tearing down one cluster must never disable
# the API for everything else in the project.
resource "google_project_service" "container_api" {
  project = local.project_id
  service = "container.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A GKE cluster — the managed Kubernetes control plane plus cluster-wide
# configuration. Node pools are separate GcpGkeNodePool resources: the
# default node pool is always removed at create time on Standard clusters,
# so every pool is an explicitly managed, first-class node. Autopilot
# clusters manage nodes themselves and take no node pools.
#
# Lifecycle notes the API enforces:
#   - name, location, description, network, subnetwork, the whole
#     ip_allocation_policy, datapath_provider, default_max_pods_per_node,
#     confidential_nodes, enable_autopilot, and the private control-plane
#     placement fields are immutable — changing any of them replaces the
#     cluster (and everything running on it).
#   - enable_l4_ilb_subsetting is one-way: it can be turned on in place but
#     never off.
#   - deletion_protection is an engine-side guard: while true, a destroy
#     plan fails before touching the cluster.
resource "google_container_cluster" "this" {
  name        = local.cluster_name
  project     = local.project_id
  location    = var.spec.location
  description = local.description

  network    = var.spec.network
  subnetwork = var.spec.subnetwork

  # Nodes may span fewer (regional) or more (zonal) zones than the control
  # plane; empty defers to GKE's location defaults.
  node_locations = length(var.spec.node_locations) > 0 ? var.spec.node_locations : null

  deletion_protection = var.spec.deletion_protection

  # Clusters are always VPC-native (alias IP). The ip_allocation_policy
  # block is emitted even when the spec omits ip_allocation: an empty block
  # tells GKE to create and manage the pod/service secondary ranges itself,
  # while named ranges pin allocation to planned subnetwork ranges.
  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {
    cluster_secondary_range_name  = try(var.spec.ip_allocation.cluster_secondary_range_name, "") != "" ? var.spec.ip_allocation.cluster_secondary_range_name : null
    services_secondary_range_name = try(var.spec.ip_allocation.services_secondary_range_name, "") != "" ? var.spec.ip_allocation.services_secondary_range_name : null
    cluster_ipv4_cidr_block       = try(var.spec.ip_allocation.cluster_ipv4_cidr_block, "") != "" ? var.spec.ip_allocation.cluster_ipv4_cidr_block : null
    services_ipv4_cidr_block      = try(var.spec.ip_allocation.services_ipv4_cidr_block, "") != "" ? var.spec.ip_allocation.services_ipv4_cidr_block : null
    stack_type                    = try(var.spec.ip_allocation.stack_type, "IPV4")

    dynamic "additional_pod_ranges_config" {
      for_each = length(try(var.spec.ip_allocation.additional_pod_range_names, [])) > 0 ? [1] : []
      content {
        pod_range_names = var.spec.ip_allocation.additional_pod_range_names
      }
    }

    dynamic "pod_cidr_overprovision_config" {
      for_each = try(var.spec.ip_allocation.pod_cidr_overprovision_disabled, false) ? [1] : []
      content {
        disabled = true
      }
    }
  }

  # Standard clusters: drop the API-mandated default node pool immediately —
  # node pools are composed as GcpGkeNodePool resources. Autopilot rejects
  # both fields (GKE owns node management there).
  remove_default_node_pool = local.is_autopilot ? null : true
  initial_node_count       = local.is_autopilot ? null : 1

  enable_autopilot = local.is_autopilot ? true : null
  allow_net_admin  = var.spec.allow_net_admin ? true : null

  datapath_provider          = local.datapath_provider
  default_max_pods_per_node  = var.spec.default_max_pods_per_node
  enable_intranode_visibility = var.spec.enable_intranode_visibility
  enable_l4_ilb_subsetting    = var.spec.enable_l4_ilb_subsetting
  enable_fqdn_network_policy  = var.spec.enable_fqdn_network_policy
  enable_cilium_clusterwide_network_policy = var.spec.enable_cilium_clusterwide_network_policy
  enable_multi_networking     = var.spec.enable_multi_networking
  private_ipv6_google_access  = local.private_ipv6_google_access
  in_transit_encryption_config = local.in_transit_encryption
  disable_l4_lb_firewall_reconciliation = var.spec.disable_l4_lb_firewall_reconciliation

  # Calico NetworkPolicy enforcement is two coupled settings: the
  # enforcement block here and the addon below. Both follow the single
  # spec toggle so they can never drift apart. Omitted entirely on
  # Autopilot (which enforces NetworkPolicy natively via Dataplane V2).
  dynamic "network_policy" {
    for_each = var.spec.enable_network_policy ? [1] : []
    content {
      enabled  = true
      provider = "CALICO"
    }
  }

  dynamic "default_snat_status" {
    for_each = var.spec.disable_default_snat ? [1] : []
    content {
      disabled = true
    }
  }

  dynamic "network_performance_config" {
    for_each = var.spec.total_egress_bandwidth_tier != "" ? [1] : []
    content {
      total_egress_bandwidth_tier = var.spec.total_egress_bandwidth_tier
    }
  }

  dynamic "dns_config" {
    for_each = var.spec.dns_config != null ? [var.spec.dns_config] : []
    content {
      cluster_dns                   = dns_config.value.cluster_dns != "" ? dns_config.value.cluster_dns : null
      cluster_dns_scope             = dns_config.value.cluster_dns_scope != "" ? dns_config.value.cluster_dns_scope : null
      cluster_dns_domain            = dns_config.value.cluster_dns_domain != "" ? dns_config.value.cluster_dns_domain : null
      additive_vpc_scope_dns_domain = dns_config.value.additive_vpc_scope_dns_domain != "" ? dns_config.value.additive_vpc_scope_dns_domain : null
    }
  }

  dynamic "gateway_api_config" {
    for_each = var.spec.gateway_api_channel != "" ? [1] : []
    content {
      channel = var.spec.gateway_api_channel
    }
  }

  dynamic "service_external_ips_config" {
    for_each = var.spec.enable_service_external_ips ? [1] : []
    content {
      enabled = true
    }
  }

  # Private topology: private nodes need Cloud NAT (a GcpRouterNat on the
  # network) for outbound internet; a private-only endpoint removes public
  # kubectl access entirely.
  dynamic "private_cluster_config" {
    for_each = var.spec.private_cluster != null ? [var.spec.private_cluster] : []
    content {
      enable_private_nodes    = private_cluster_config.value.enable_private_nodes
      enable_private_endpoint = private_cluster_config.value.enable_private_endpoint
      master_ipv4_cidr_block  = private_cluster_config.value.master_ipv4_cidr_block != "" ? private_cluster_config.value.master_ipv4_cidr_block : null
      private_endpoint_subnetwork = private_cluster_config.value.private_endpoint_subnetwork != "" ? private_cluster_config.value.private_endpoint_subnetwork : null

      dynamic "master_global_access_config" {
        for_each = private_cluster_config.value.enable_master_global_access ? [1] : []
        content {
          enabled = true
        }
      }
    }
  }

  dynamic "master_authorized_networks_config" {
    for_each = var.spec.master_authorized_networks != null ? [var.spec.master_authorized_networks] : []
    content {
      gcp_public_cidrs_access_enabled      = master_authorized_networks_config.value.gcp_public_cidrs_access_enabled
      private_endpoint_enforcement_enabled = master_authorized_networks_config.value.private_endpoint_enforcement_enabled

      dynamic "cidr_blocks" {
        for_each = master_authorized_networks_config.value.cidr_blocks
        content {
          cidr_block   = cidr_blocks.value.cidr_block
          display_name = cidr_blocks.value.display_name != "" ? cidr_blocks.value.display_name : null
        }
      }
    }
  }

  dynamic "control_plane_endpoints_config" {
    for_each = var.spec.control_plane_endpoints != null ? [var.spec.control_plane_endpoints] : []
    content {
      dns_endpoint_config {
        allow_external_traffic = control_plane_endpoints_config.value.dns_endpoint_allow_external_traffic
      }
      ip_endpoints_config {
        enabled = control_plane_endpoints_config.value.ip_endpoints_enabled
      }
    }
  }

  release_channel {
    channel = local.release_channel
  }

  min_master_version = local.min_master_version

  dynamic "maintenance_policy" {
    for_each = var.spec.maintenance_policy != null ? [var.spec.maintenance_policy] : []
    content {
      dynamic "daily_maintenance_window" {
        for_each = maintenance_policy.value.daily_window != null ? [maintenance_policy.value.daily_window] : []
        content {
          start_time = daily_maintenance_window.value.start_time
        }
      }

      dynamic "recurring_window" {
        for_each = maintenance_policy.value.recurring_window != null ? [maintenance_policy.value.recurring_window] : []
        content {
          start_time = recurring_window.value.start_time
          end_time   = recurring_window.value.end_time
          recurrence = recurring_window.value.recurrence
        }
      }

      dynamic "maintenance_exclusion" {
        for_each = maintenance_policy.value.exclusions
        content {
          exclusion_name = maintenance_exclusion.value.exclusion_name
          start_time     = maintenance_exclusion.value.start_time
          end_time       = maintenance_exclusion.value.end_time

          dynamic "exclusion_options" {
            for_each = maintenance_exclusion.value.scope != "" ? [1] : []
            content {
              scope = maintenance_exclusion.value.scope
            }
          }
        }
      }
    }
  }

  # Node auto-provisioning: GKE creates/deletes node pools within the
  # resource limits — bounded by spec-level validation so an enabled NAP
  # always carries limits (an unbounded NAP is an unbounded bill).
  dynamic "cluster_autoscaling" {
    for_each = var.spec.cluster_autoscaling != null ? [var.spec.cluster_autoscaling] : []
    content {
      enabled             = cluster_autoscaling.value.enabled
      autoscaling_profile = cluster_autoscaling.value.autoscaling_profile
      auto_provisioning_locations = length(cluster_autoscaling.value.auto_provisioning_locations) > 0 ? cluster_autoscaling.value.auto_provisioning_locations : null

      dynamic "resource_limits" {
        for_each = cluster_autoscaling.value.resource_limits
        content {
          resource_type = resource_limits.value.resource_type
          minimum       = resource_limits.value.minimum
          maximum       = resource_limits.value.maximum
        }
      }

      dynamic "auto_provisioning_defaults" {
        for_each = cluster_autoscaling.value.auto_provisioning_defaults != null ? [cluster_autoscaling.value.auto_provisioning_defaults] : []
        content {
          service_account   = auto_provisioning_defaults.value.service_account != "" ? auto_provisioning_defaults.value.service_account : null
          oauth_scopes      = length(auto_provisioning_defaults.value.oauth_scopes) > 0 ? auto_provisioning_defaults.value.oauth_scopes : null
          disk_size         = auto_provisioning_defaults.value.disk_size_gb
          disk_type         = auto_provisioning_defaults.value.disk_type != "" ? auto_provisioning_defaults.value.disk_type : null
          image_type        = auto_provisioning_defaults.value.image_type != "" ? auto_provisioning_defaults.value.image_type : null
          min_cpu_platform  = auto_provisioning_defaults.value.min_cpu_platform != "" ? auto_provisioning_defaults.value.min_cpu_platform : null
          boot_disk_kms_key = auto_provisioning_defaults.value.boot_disk_kms_key != "" ? auto_provisioning_defaults.value.boot_disk_kms_key : null

          shielded_instance_config {
            enable_secure_boot          = auto_provisioning_defaults.value.enable_secure_boot
            enable_integrity_monitoring = auto_provisioning_defaults.value.enable_integrity_monitoring
          }

          management {
            auto_upgrade = auto_provisioning_defaults.value.auto_upgrade
            auto_repair  = auto_provisioning_defaults.value.auto_repair
          }
        }
      }
    }
  }

  dynamic "vertical_pod_autoscaling" {
    for_each = var.spec.enable_vertical_pod_autoscaling ? [1] : []
    content {
      enabled = true
    }
  }

  dynamic "pod_autoscaling" {
    for_each = var.spec.hpa_profile != "" ? [1] : []
    content {
      hpa_profile = var.spec.hpa_profile
    }
  }

  # Workload Identity Federation for GKE: the pool name is fixed by the API
  # to PROJECT_ID.svc.id.goog. Autopilot clusters have it always on — the
  # block is suppressed there to avoid a permanent diff.
  dynamic "workload_identity_config" {
    for_each = var.spec.workload_identity_enabled && !local.is_autopilot ? [1] : []
    content {
      workload_pool = local.workload_pool
    }
  }

  enable_shielded_nodes = var.spec.enable_shielded_nodes

  dynamic "database_encryption" {
    for_each = var.spec.database_encryption != null ? [var.spec.database_encryption] : []
    content {
      state    = database_encryption.value.state
      key_name = database_encryption.value.key_name != "" ? database_encryption.value.key_name : null
    }
  }

  dynamic "binary_authorization" {
    for_each = var.spec.binary_authorization_evaluation_mode != "" ? [1] : []
    content {
      evaluation_mode = var.spec.binary_authorization_evaluation_mode
    }
  }

  dynamic "security_posture_config" {
    for_each = var.spec.security_posture != null ? [var.spec.security_posture] : []
    content {
      mode               = security_posture_config.value.mode != "" ? security_posture_config.value.mode : null
      vulnerability_mode = security_posture_config.value.vulnerability_mode != "" ? security_posture_config.value.vulnerability_mode : null
    }
  }

  dynamic "authenticator_groups_config" {
    for_each = var.spec.authenticator_security_group != "" ? [1] : []
    content {
      security_group = var.spec.authenticator_security_group
    }
  }

  enable_legacy_abac = var.spec.enable_legacy_abac

  dynamic "mesh_certificates" {
    for_each = var.spec.enable_mesh_certificates ? [1] : []
    content {
      enable_certificates = true
    }
  }

  dynamic "secret_manager_config" {
    for_each = var.spec.enable_secret_manager_csi ? [1] : []
    content {
      enabled = true
    }
  }

  dynamic "confidential_nodes" {
    for_each = var.spec.confidential_nodes != null ? [var.spec.confidential_nodes] : []
    content {
      enabled                    = confidential_nodes.value.enabled
      confidential_instance_type = confidential_nodes.value.confidential_instance_type != "" ? confidential_nodes.value.confidential_instance_type : null
    }
  }

  dynamic "anonymous_authentication_config" {
    for_each = var.spec.anonymous_authentication_mode != "" ? [1] : []
    content {
      mode = var.spec.anonymous_authentication_mode
    }
  }

  dynamic "identity_service_config" {
    for_each = var.spec.enable_identity_service ? [1] : []
    content {
      enabled = true
    }
  }

  # An explicit empty components list is meaningful: it disables the Cloud
  # Logging/Monitoring integration outright, so the spec's presence (not
  # emptiness) drives whether the block is emitted.
  dynamic "logging_config" {
    for_each = var.spec.logging != null ? [var.spec.logging] : []
    content {
      enable_components = logging_config.value.components
    }
  }

  dynamic "monitoring_config" {
    for_each = var.spec.monitoring != null ? [var.spec.monitoring] : []
    content {
      enable_components = length(monitoring_config.value.components) > 0 ? monitoring_config.value.components : null

      managed_prometheus {
        enabled = monitoring_config.value.managed_prometheus_enabled

        dynamic "auto_monitoring_config" {
          for_each = monitoring_config.value.auto_monitoring_scope != "" ? [1] : []
          content {
            scope = monitoring_config.value.auto_monitoring_scope
          }
        }
      }

      dynamic "advanced_datapath_observability_config" {
        for_each = (monitoring_config.value.advanced_datapath_metrics_enabled || monitoring_config.value.advanced_datapath_relay_enabled) ? [1] : []
        content {
          enable_metrics = monitoring_config.value.advanced_datapath_metrics_enabled
          enable_relay   = monitoring_config.value.advanced_datapath_relay_enabled
        }
      }
    }
  }

  dynamic "notification_config" {
    for_each = var.spec.notification_pubsub != null ? [var.spec.notification_pubsub] : []
    content {
      pubsub {
        enabled = notification_config.value.enabled
        topic   = notification_config.value.topic != "" ? notification_config.value.topic : null

        dynamic "filter" {
          for_each = length(notification_config.value.event_types) > 0 ? [1] : []
          content {
            event_type = notification_config.value.event_types
          }
        }
      }
    }
  }

  dynamic "cost_management_config" {
    for_each = var.spec.enable_cost_management ? [1] : []
    content {
      enabled = true
    }
  }

  dynamic "resource_usage_export_config" {
    for_each = var.spec.resource_usage_export != null ? [var.spec.resource_usage_export] : []
    content {
      enable_network_egress_metering       = resource_usage_export_config.value.enable_network_egress_metering
      enable_resource_consumption_metering = resource_usage_export_config.value.enable_resource_consumption_metering

      bigquery_destination {
        dataset_id = resource_usage_export_config.value.bigquery_dataset_id
      }
    }
  }

  # Addons: emitted when the spec configures them or when Calico network
  # policy needs its companion addon. The network-policy addon toggle always
  # mirrors enable_network_policy (never set on Autopilot).
  dynamic "addons_config" {
    for_each = (var.spec.addons != null || (var.spec.enable_network_policy && !local.is_autopilot)) ? [1] : []
    content {
      dynamic "network_policy_config" {
        for_each = local.is_autopilot ? [] : [1]
        content {
          disabled = !var.spec.enable_network_policy
        }
      }

      http_load_balancing {
        disabled = !try(var.spec.addons.http_load_balancing_enabled, true)
      }

      horizontal_pod_autoscaling {
        disabled = !try(var.spec.addons.horizontal_pod_autoscaling_enabled, true)
      }

      gce_persistent_disk_csi_driver_config {
        enabled = try(var.spec.addons.gce_persistent_disk_csi_driver_enabled, true)
      }

      dynamic "gcp_filestore_csi_driver_config" {
        for_each = try(var.spec.addons.gcp_filestore_csi_driver_enabled, false) ? [1] : []
        content {
          enabled = true
        }
      }

      dynamic "gcs_fuse_csi_driver_config" {
        for_each = try(var.spec.addons.gcs_fuse_csi_driver_enabled, false) ? [1] : []
        content {
          enabled = true
        }
      }

      dynamic "gke_backup_agent_config" {
        for_each = try(var.spec.addons.gke_backup_agent_enabled, false) ? [1] : []
        content {
          enabled = true
        }
      }

      dynamic "dns_cache_config" {
        for_each = try(var.spec.addons.dns_cache_enabled, false) ? [1] : []
        content {
          enabled = true
        }
      }

      dynamic "config_connector_config" {
        for_each = try(var.spec.addons.config_connector_enabled, false) ? [1] : []
        content {
          enabled = true
        }
      }

      dynamic "stateful_ha_config" {
        for_each = try(var.spec.addons.stateful_ha_enabled, false) ? [1] : []
        content {
          enabled = true
        }
      }

      dynamic "ray_operator_config" {
        for_each = try(var.spec.addons.ray_operator_enabled, false) ? [1] : []
        content {
          enabled = true
        }
      }
    }
  }

  dynamic "fleet" {
    for_each = var.spec.fleet_project != "" ? [1] : []
    content {
      project = var.spec.fleet_project
    }
  }

  resource_labels = local.final_labels

  depends_on = [google_project_service.container_api]
}
