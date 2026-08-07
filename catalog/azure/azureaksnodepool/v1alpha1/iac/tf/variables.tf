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
  description = "Azure AKS node pool specification"
  type = object({
    # The parent AKS cluster's ARM ID. The pool is an ARM child of the
    # cluster. References are resolved to a literal ID by the platform
    # before the module runs.
    kubernetes_cluster_id = string

    # Agent-pool name (1-12 lowercase alphanumerics starting with a
    # letter; at most 6 characters for Windows pools).
    name = string

    # Azure VM size for the pool's nodes.
    vm_size = string

    # Purpose (SYSTEM / USER) and OS, as spec enum name strings.
    mode    = optional(string)
    os_type = optional(string)
    os_sku  = optional(string)

    # Sizing: fixed node_count, or autoscaling bounds.
    node_count           = optional(number)
    auto_scaling_enabled = optional(bool, false)
    min_count            = optional(number)
    max_count            = optional(number)
    max_pods             = optional(number)

    # Spot economics (USER pools only).
    priority        = optional(string)
    eviction_policy = optional(string)
    spot_max_price  = optional(number)

    # Scheduling surface.
    node_labels = optional(map(string), {})
    node_taints = optional(list(string), [])
    zones       = optional(list(string), [])

    # Placement network. Unset inherits the cluster's node subnet.
    vnet_subnet_id = optional(string)
    pod_subnet_id  = optional(string)

    # Node Kubernetes version (unset follows the control plane).
    orchestrator_version = optional(string)

    # Disks.
    os_disk_size_gb   = optional(number)
    os_disk_type      = optional(string)
    kubelet_disk_type = optional(string)
    ultra_ssd_enabled = optional(bool, false)

    # Hardening and compliance.
    fips_enabled            = optional(bool, false)
    host_encryption_enabled = optional(bool, false)

    # Node public IPs.
    node_public_ip_enabled   = optional(bool, false)
    node_public_ip_prefix_id = optional(string)

    # GPU.
    gpu_instance = optional(string)
    gpu_driver   = optional(string)

    # Placement groups (plain ARM ids).
    proximity_placement_group_id  = optional(string)
    host_group_id                 = optional(string)
    capacity_reservation_group_id = optional(string)

    # Operations.
    scale_down_mode             = optional(string)
    snapshot_id                 = optional(string)
    workload_runtime            = optional(string)
    temporary_name_for_rotation = optional(string)

    # Upgrade rollout: max_surge XOR max_unavailable (spec-level
    # validation enforces the exclusivity and the spot rules).
    upgrade_settings = optional(object({
      max_surge                     = optional(string)
      max_unavailable               = optional(string)
      drain_timeout_in_minutes      = optional(number)
      node_soak_duration_in_minutes = optional(number)
      undrainable_node_behavior     = optional(string)
    }))

    # Kubelet tuning.
    kubelet_config = optional(object({
      cpu_manager_policy        = optional(string)
      cpu_cfs_quota_enabled     = optional(bool, true)
      cpu_cfs_quota_period      = optional(string)
      image_gc_high_threshold   = optional(number)
      image_gc_low_threshold    = optional(number)
      topology_manager_policy   = optional(string)
      allowed_unsafe_sysctls    = optional(list(string), [])
      container_log_max_size_mb = optional(number)
      container_log_max_files   = optional(number)
      pod_max_pid               = optional(number)
    }))

    # Linux kernel/OS tuning (Linux pools only).
    linux_os_config = optional(object({
      sysctl_config = optional(object({
        fs_aio_max_nr                      = optional(number)
        fs_file_max                        = optional(number)
        fs_inotify_max_user_watches        = optional(number)
        fs_nr_open                         = optional(number)
        kernel_threads_max                 = optional(number)
        net_core_netdev_max_backlog        = optional(number)
        net_core_optmem_max                = optional(number)
        net_core_rmem_default              = optional(number)
        net_core_rmem_max                  = optional(number)
        net_core_somaxconn                 = optional(number)
        net_core_wmem_default              = optional(number)
        net_core_wmem_max                  = optional(number)
        net_ipv4_ip_local_port_range_min   = optional(number)
        net_ipv4_ip_local_port_range_max   = optional(number)
        net_ipv4_neigh_default_gc_thresh1  = optional(number)
        net_ipv4_neigh_default_gc_thresh2  = optional(number)
        net_ipv4_neigh_default_gc_thresh3  = optional(number)
        net_ipv4_tcp_fin_timeout           = optional(number)
        net_ipv4_tcp_keepalive_intvl       = optional(number)
        net_ipv4_tcp_keepalive_probes      = optional(number)
        net_ipv4_tcp_keepalive_time        = optional(number)
        net_ipv4_tcp_max_syn_backlog       = optional(number)
        net_ipv4_tcp_max_tw_buckets        = optional(number)
        net_ipv4_tcp_tw_reuse              = optional(bool, false)
        net_netfilter_nf_conntrack_buckets = optional(number)
        net_netfilter_nf_conntrack_max     = optional(number)
        vm_max_map_count                   = optional(number)
        vm_swappiness                      = optional(number)
        vm_vfs_cache_pressure              = optional(number)
      }))
      transparent_huge_page        = optional(string)
      transparent_huge_page_defrag = optional(string)
      swap_file_size_mb            = optional(number)
    }))

    # Node-level network hardening.
    node_network_profile = optional(object({
      allowed_host_ports = optional(list(object({
        port_start = optional(number)
        port_end   = optional(number)
        protocol   = optional(string)
      })), [])
      application_security_group_ids = optional(list(string), [])
      node_public_ip_tags            = optional(map(string), {})
    }))

    # Windows-pool networking (Windows pools only).
    windows_profile = optional(object({
      outbound_nat_enabled = optional(bool, true)
    }))

    # Free-form user tags, merged over the metadata-derived tags.
    tags = optional(map(string), {})
  })
}
