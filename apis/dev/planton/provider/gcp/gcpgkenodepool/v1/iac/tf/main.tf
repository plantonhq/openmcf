# Enable the Kubernetes Engine API so a fresh project can host the pool.
# disable_on_destroy is false: tearing down one node pool must never disable
# the API for everything else in the project (including its own cluster).
resource "google_project_service" "container_api" {
  project = local.project_id
  service = "container.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A GKE node pool — a group of identically configured Compute Engine VMs
# attached to a GcpGkeCluster. The cluster and location arrive as plain
# strings resolved from the cluster's outputs, so the pool addresses its
# parent exactly the way GKE named it — no lookup, no wildcard search.
#
# Lifecycle notes the API enforces:
#   - name, location, initial_node_count, max_pods_per_node,
#     placement_policy, queued_provisioning, and nearly all of node_config
#     (machine type, disks, image, identity, accelerators, shielded/
#     confidential settings, local SSDs) are immutable — changing them
#     replaces the pool (GKE drains and recreates the nodes).
#   - node_count, autoscaling, management, upgrade_settings, node_locations,
#     labels, taints, tags, and resource_labels update in place.
#   - For autoscaled pools the autoscaler owns node_count at runtime;
#     lifecycle.ignore_changes keeps the plan from fighting it.
resource "google_container_node_pool" "this" {
  name     = local.node_pool_name
  project  = local.project_id
  cluster  = var.spec.cluster_name
  location = var.spec.location

  node_locations = length(var.spec.node_locations) > 0 ? var.spec.node_locations : null

  # Explicit version and auto_upgrade fight each other (the API re-upgrades
  # what the plan pins); the spec documents the trade and the field passes
  # through untouched.
  version = local.version

  max_pods_per_node  = var.spec.max_pods_per_node
  initial_node_count = var.spec.initial_node_count

  # Fixed size XOR autoscaling — the proto oneof guarantees at most one
  # arrives. Neither means GKE's default size (3 nodes), unmanaged.
  node_count = var.spec.node_count

  dynamic "autoscaling" {
    for_each = var.spec.autoscaling != null ? [var.spec.autoscaling] : []
    content {
      # Per-zone bounds XOR total bounds (spec-level CEL enforces the
      # exclusivity); nulls drop the unused arm from the API payload.
      min_node_count       = autoscaling.value.min_nodes
      max_node_count       = autoscaling.value.max_nodes
      total_min_node_count = autoscaling.value.total_min_nodes
      total_max_node_count = autoscaling.value.total_max_nodes
      location_policy      = autoscaling.value.location_policy != "" ? autoscaling.value.location_policy : null
    }
  }

  # Auto-repair/auto-upgrade both default true — GKE's own defaults — so an
  # omitted management block and GKE's behavior agree.
  management {
    auto_repair  = try(var.spec.management.auto_repair, true)
    auto_upgrade = try(var.spec.management.auto_upgrade, true)
  }

  # Emitted only when the spec configures it: GKE's default surge settings
  # (max_surge=1, max_unavailable=0) apply otherwise, and omitting the
  # block keeps the plan clean on pools that never touch it.
  dynamic "upgrade_settings" {
    for_each = var.spec.upgrade_settings != null ? [var.spec.upgrade_settings] : []
    content {
      max_surge       = upgrade_settings.value.max_surge
      max_unavailable = upgrade_settings.value.max_unavailable
      strategy        = upgrade_settings.value.strategy != "" ? upgrade_settings.value.strategy : null

      dynamic "blue_green_settings" {
        for_each = upgrade_settings.value.blue_green_settings != null ? [upgrade_settings.value.blue_green_settings] : []
        content {
          node_pool_soak_duration = blue_green_settings.value.node_pool_soak_duration != "" ? blue_green_settings.value.node_pool_soak_duration : null

          standard_rollout_policy {
            batch_percentage    = blue_green_settings.value.standard_rollout_policy.batch_percentage
            batch_node_count    = blue_green_settings.value.standard_rollout_policy.batch_node_count
            batch_soak_duration = blue_green_settings.value.standard_rollout_policy.batch_soak_duration != "" ? blue_green_settings.value.standard_rollout_policy.batch_soak_duration : null
          }
        }
      }
    }
  }

  dynamic "placement_policy" {
    for_each = var.spec.placement_policy != null ? [var.spec.placement_policy] : []
    content {
      type         = placement_policy.value.type
      policy_name  = placement_policy.value.policy_name != "" ? placement_policy.value.policy_name : null
      tpu_topology = placement_policy.value.tpu_topology != "" ? placement_policy.value.tpu_topology : null
    }
  }

  dynamic "queued_provisioning" {
    for_each = var.spec.queued_provisioning_enabled ? [1] : []
    content {
      enabled = true
    }
  }

  dynamic "network_config" {
    for_each = var.spec.network_config != null ? [var.spec.network_config] : []
    content {
      create_pod_range     = network_config.value.create_pod_range
      pod_range            = network_config.value.pod_range != "" ? network_config.value.pod_range : null
      pod_ipv4_cidr_block  = network_config.value.pod_ipv4_cidr_block != "" ? network_config.value.pod_ipv4_cidr_block : null
      enable_private_nodes = network_config.value.enable_private_nodes

      dynamic "network_performance_config" {
        for_each = network_config.value.total_egress_bandwidth_tier != "" ? [1] : []
        content {
          total_egress_bandwidth_tier = network_config.value.total_egress_bandwidth_tier
        }
      }

      dynamic "pod_cidr_overprovision_config" {
        for_each = network_config.value.pod_cidr_overprovision_disabled ? [1] : []
        content {
          disabled = true
        }
      }
    }
  }

  # node_config is always emitted: the platform attribution resource_labels
  # and the disable-legacy-endpoints metadata guard apply to every pool,
  # spec block or not. Everything else inside honors "unset means GKE
  # default" by passing null.
  node_config {
    machine_type    = local.machine_type
    disk_size_gb    = try(local.nc.disk_size_gb, null)
    disk_type       = local.disk_type
    image_type      = local.image_type
    service_account = local.service_account

    # Empty applies GKE's default scopes; with Workload Identity, workload
    # permissions come from IAM on Kubernetes service accounts, not node
    # scopes, so the defaults are normally right.
    oauth_scopes = length(try(local.nc.oauth_scopes, [])) > 0 ? local.nc.oauth_scopes : null

    # Kubernetes node labels (nodeSelector targets) — distinct from the GCE
    # resource labels below, which carry the platform attribution.
    labels          = try(local.nc.labels, {})
    resource_labels = local.final_resource_labels
    tags            = length(try(local.nc.tags, [])) > 0 ? local.nc.tags : null
    metadata        = local.node_metadata

    dynamic "taint" {
      for_each = try(local.nc.taints, [])
      content {
        key    = taint.value.key
        value  = taint.value.value
        effect = taint.value.effect
      }
    }

    # spot supersedes preemptible (no 24h lifetime); spec-level CEL rejects
    # both together, so each passes through independently.
    spot        = try(local.nc.spot, false)
    preemptible = try(local.nc.preemptible, false)

    dynamic "guest_accelerator" {
      for_each = try(local.nc.guest_accelerators, [])
      content {
        type               = guest_accelerator.value.type
        count              = guest_accelerator.value.count
        gpu_partition_size = guest_accelerator.value.gpu_partition_size != "" ? guest_accelerator.value.gpu_partition_size : null

        dynamic "gpu_driver_installation_config" {
          for_each = guest_accelerator.value.gpu_driver_version != "" ? [1] : []
          content {
            gpu_driver_version = guest_accelerator.value.gpu_driver_version
          }
        }

        dynamic "gpu_sharing_config" {
          for_each = guest_accelerator.value.gpu_sharing_config != null ? [guest_accelerator.value.gpu_sharing_config] : []
          content {
            gpu_sharing_strategy       = gpu_sharing_config.value.gpu_sharing_strategy
            max_shared_clients_per_gpu = gpu_sharing_config.value.max_shared_clients_per_gpu
          }
        }
      }
    }

    dynamic "shielded_instance_config" {
      for_each = try(local.nc.shielded_instance_config, null) != null ? [local.nc.shielded_instance_config] : []
      content {
        enable_secure_boot          = shielded_instance_config.value.enable_secure_boot
        enable_integrity_monitoring = shielded_instance_config.value.enable_integrity_monitoring
      }
    }

    dynamic "confidential_nodes" {
      for_each = try(local.nc.confidential_nodes, null) != null ? [local.nc.confidential_nodes] : []
      content {
        enabled                    = confidential_nodes.value.enabled
        confidential_instance_type = confidential_nodes.value.confidential_instance_type != "" ? confidential_nodes.value.confidential_instance_type : null
      }
    }

    min_cpu_platform = local.min_cpu_platform
    local_ssd_count  = try(local.nc.local_ssd_count, null)

    dynamic "ephemeral_storage_local_ssd_config" {
      for_each = try(local.nc.ephemeral_storage_local_ssd, null) != null ? [local.nc.ephemeral_storage_local_ssd] : []
      content {
        local_ssd_count  = ephemeral_storage_local_ssd_config.value.local_ssd_count
        data_cache_count = ephemeral_storage_local_ssd_config.value.data_cache_count
      }
    }

    dynamic "local_nvme_ssd_block_config" {
      for_each = try(local.nc.local_nvme_ssd_block, null) != null ? [local.nc.local_nvme_ssd_block] : []
      content {
        local_ssd_count = local_nvme_ssd_block_config.value.local_ssd_count
      }
    }

    dynamic "gcfs_config" {
      for_each = try(local.nc.gcfs_enabled, false) ? [1] : []
      content {
        enabled = true
      }
    }

    dynamic "gvnic" {
      for_each = try(local.nc.gvnic_enabled, false) ? [1] : []
      content {
        enabled = true
      }
    }

    dynamic "fast_socket" {
      for_each = try(local.nc.fast_socket_enabled, false) ? [1] : []
      content {
        enabled = true
      }
    }

    boot_disk_kms_key = local.boot_disk_kms_key

    dynamic "workload_metadata_config" {
      for_each = try(local.nc.workload_metadata_mode, "") != "" ? [1] : []
      content {
        mode = local.nc.workload_metadata_mode
      }
    }

    dynamic "reservation_affinity" {
      for_each = try(local.nc.reservation_affinity, null) != null ? [local.nc.reservation_affinity] : []
      content {
        consume_reservation_type = reservation_affinity.value.consume_reservation_type
        key                      = reservation_affinity.value.key != "" ? reservation_affinity.value.key : null
        values                   = length(reservation_affinity.value.values) > 0 ? reservation_affinity.value.values : null
      }
    }

    dynamic "secondary_boot_disks" {
      for_each = try(local.nc.secondary_boot_disks, [])
      content {
        disk_image = secondary_boot_disks.value.disk_image
        mode       = secondary_boot_disks.value.mode != "" ? secondary_boot_disks.value.mode : null
      }
    }

    dynamic "kubelet_config" {
      for_each = try(local.nc.kubelet_config, null) != null ? [local.nc.kubelet_config] : []
      content {
        cpu_manager_policy                     = kubelet_config.value.cpu_manager_policy != "" ? kubelet_config.value.cpu_manager_policy : null
        cpu_cfs_quota                          = kubelet_config.value.cpu_cfs_quota
        cpu_cfs_quota_period                   = kubelet_config.value.cpu_cfs_quota_period != "" ? kubelet_config.value.cpu_cfs_quota_period : null
        pod_pids_limit                         = kubelet_config.value.pod_pids_limit
        insecure_kubelet_readonly_port_enabled = kubelet_config.value.insecure_kubelet_readonly_port_enabled != "" ? kubelet_config.value.insecure_kubelet_readonly_port_enabled : null
        max_parallel_image_pulls               = kubelet_config.value.max_parallel_image_pulls
        container_log_max_size                 = kubelet_config.value.container_log_max_size != "" ? kubelet_config.value.container_log_max_size : null
        container_log_max_files                = kubelet_config.value.container_log_max_files
        image_gc_low_threshold_percent         = kubelet_config.value.image_gc_low_threshold_percent
        image_gc_high_threshold_percent        = kubelet_config.value.image_gc_high_threshold_percent
      }
    }

    dynamic "linux_node_config" {
      for_each = try(local.nc.linux_node_config, null) != null ? [local.nc.linux_node_config] : []
      content {
        sysctls     = length(linux_node_config.value.sysctls) > 0 ? linux_node_config.value.sysctls : null
        cgroup_mode = linux_node_config.value.cgroup_mode != "" ? linux_node_config.value.cgroup_mode : null

        dynamic "hugepages_config" {
          for_each = linux_node_config.value.hugepages_config != null ? [linux_node_config.value.hugepages_config] : []
          content {
            hugepage_size_2m = hugepages_config.value.hugepage_size_2m
            hugepage_size_1g = hugepages_config.value.hugepage_size_1g
          }
        }
      }
    }

    logging_variant  = local.logging_variant
    flex_start       = try(local.nc.flex_start, false) ? true : null
    max_run_duration = local.max_run_duration
  }

  # For autoscaled pools the cluster autoscaler owns the live node count;
  # without this, every plan after a scale event would try to reset it.
  lifecycle {
    ignore_changes = [
      node_count
    ]
  }

  depends_on = [google_project_service.container_api]
}
