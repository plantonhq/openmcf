variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for the GCP GKE node pool"
  type = object({
    # The GCP project the pool is created in (must be the cluster's
    # project). The CLI's tfvars converter resolves StringValueOrRef fields
    # to their literal string before the module runs, so this arrives as a
    # plain string. Empty falls back to the provider's default project.
    project_id = optional(string, "")

    # Parent cluster name (plain string after ref resolution). Immutable.
    cluster_name = string

    # Parent cluster location — region or zone (plain string after ref
    # resolution). Immutable.
    location = string

    # Pool name in GKE. Empty means "use metadata.name". Immutable.
    node_pool_name = optional(string, "")

    # Zones the pool's nodes run in; empty inherits the cluster's
    # node_locations.
    node_locations = optional(list(string), [])

    # Explicit node Kubernetes version; empty lets GKE/auto-upgrade drive.
    version = optional(string, "")

    # Pods-per-node override (8-256); null inherits the cluster default.
    max_pods_per_node = optional(number, null)

    # Starting size for autoscaled pools (per zone). Immutable.
    initial_node_count = optional(number, null)

    # Fixed pool size (per zone). Mutually exclusive with autoscaling —
    # the proto oneof guarantees at most one arrives.
    node_count = optional(number, null)

    # Cluster-autoscaler bounds: per-zone (min/max) XOR total limits.
    autoscaling = optional(object({
      min_nodes       = optional(number, null)
      max_nodes       = optional(number, null)
      total_min_nodes = optional(number, null)
      total_max_nodes = optional(number, null)
      location_policy = optional(string, "")
    }), null)

    # Auto-repair / auto-upgrade; both default true (GKE's own defaults).
    management = optional(object({
      auto_repair  = optional(bool, true)
      auto_upgrade = optional(bool, true)
    }), null)

    # Surge or blue-green upgrade rollout controls.
    upgrade_settings = optional(object({
      max_surge       = optional(number, null)
      max_unavailable = optional(number, null)
      strategy        = optional(string, "")
      blue_green_settings = optional(object({
        standard_rollout_policy = object({
          batch_percentage    = optional(number, null)
          batch_node_count    = optional(number, null)
          batch_soak_duration = optional(string, "")
        })
        node_pool_soak_duration = optional(string, "")
      }), null)
    }), null)

    # Compact placement / TPU topology. Immutable.
    placement_policy = optional(object({
      type         = string
      policy_name  = optional(string, "")
      tpu_topology = optional(string, "")
    }), null)

    # Dynamic Workload Scheduler queued provisioning. Immutable.
    queued_provisioning_enabled = optional(bool, false)

    # Pool-level networking overrides.
    network_config = optional(object({
      create_pod_range                = optional(bool, false)
      pod_range                       = optional(string, "")
      pod_ipv4_cidr_block             = optional(string, "")
      enable_private_nodes            = optional(bool, null)
      total_egress_bandwidth_tier     = optional(string, "")
      pod_cidr_overprovision_disabled = optional(bool, false)
    }), null)

    # Node VM configuration shared by every node in the pool.
    node_config = optional(object({
      machine_type = optional(string, "e2-medium")
      disk_size_gb = optional(number, null)
      disk_type    = optional(string, "")
      image_type   = optional(string, "COS_CONTAINERD")

      # Service account email (plain string after ref resolution).
      service_account = optional(string, "")
      oauth_scopes    = optional(list(string), [])

      labels          = optional(map(string), {})
      resource_labels = optional(map(string), {})
      tags            = optional(list(string), [])
      metadata        = optional(map(string), {})

      taints = optional(list(object({
        key    = string
        value  = string
        effect = string
      })), [])

      spot        = optional(bool, false)
      preemptible = optional(bool, false)

      guest_accelerators = optional(list(object({
        type               = string
        count              = number
        gpu_partition_size = optional(string, "")
        gpu_driver_version = optional(string, "")
        gpu_sharing_config = optional(object({
          gpu_sharing_strategy       = string
          max_shared_clients_per_gpu = number
        }), null)
      })), [])

      shielded_instance_config = optional(object({
        enable_secure_boot          = optional(bool, false)
        enable_integrity_monitoring = optional(bool, true)
      }), null)

      confidential_nodes = optional(object({
        enabled                    = optional(bool, false)
        confidential_instance_type = optional(string, "")
      }), null)

      min_cpu_platform = optional(string, "")
      local_ssd_count  = optional(number, null)

      ephemeral_storage_local_ssd = optional(object({
        local_ssd_count  = number
        data_cache_count = optional(number, null)
      }), null)

      local_nvme_ssd_block = optional(object({
        local_ssd_count = number
      }), null)

      gcfs_enabled        = optional(bool, false)
      gvnic_enabled       = optional(bool, false)
      fast_socket_enabled = optional(bool, false)

      # KMS key path (plain string after ref resolution).
      boot_disk_kms_key = optional(string, "")

      workload_metadata_mode = optional(string, "")

      reservation_affinity = optional(object({
        consume_reservation_type = string
        key                      = optional(string, "")
        values                   = optional(list(string), [])
      }), null)

      secondary_boot_disks = optional(list(object({
        disk_image = string
        mode       = optional(string, "")
      })), [])

      kubelet_config = optional(object({
        cpu_manager_policy                     = optional(string, "")
        cpu_cfs_quota                          = optional(bool, null)
        cpu_cfs_quota_period                   = optional(string, "")
        pod_pids_limit                         = optional(number, null)
        insecure_kubelet_readonly_port_enabled = optional(string, "")
        max_parallel_image_pulls               = optional(number, null)
        container_log_max_size                 = optional(string, "")
        container_log_max_files                = optional(number, null)
        image_gc_low_threshold_percent         = optional(number, null)
        image_gc_high_threshold_percent        = optional(number, null)
      }), null)

      linux_node_config = optional(object({
        sysctls     = optional(map(string), {})
        cgroup_mode = optional(string, "")
        hugepages_config = optional(object({
          hugepage_size_2m = optional(number, null)
          hugepage_size_1g = optional(number, null)
        }), null)
      }), null)

      logging_variant = optional(string, "")
      flex_start      = optional(bool, false)
      max_run_duration = optional(string, "")
    }), null)
  })
}
