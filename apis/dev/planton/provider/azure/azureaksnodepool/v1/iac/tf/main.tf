# Create the AKS node pool: a scale set of worker nodes attached to an
# existing AKS cluster by ARM ID.
#
# Design and lifecycle notes worth knowing before operating this resource:
# - The pool is the unit of compute shape: general, memory-optimized,
#   GPU, spot, or Windows pools each live as their own resource with an
#   independent lifecycle. The cluster carries only its mandatory default
#   (system) pool.
# - Spot pools (priority SPOT) trade 30-90% cost savings for evictability
#   and automatically carry the scalesetpriority spot taint; AKS replaces
#   evicted nodes as capacity returns. Eviction policy and max price are
#   fixed at creation.
# - Many shape changes (vm_size, os_disk_type, fips, host encryption)
#   rotate the pool. Setting temporary_name_for_rotation lets AKS stand
#   up a replacement pool first instead of tearing this one down --
#   set it proactively on production pools.
# - Node Kubernetes versions may lag the control plane by up to two
#   minor versions: orchestrator_version is the seam for canarying node
#   upgrades pool by pool.
resource "azurerm_kubernetes_cluster_node_pool" "main" {
  name                  = var.spec.name
  kubernetes_cluster_id = var.spec.kubernetes_cluster_id
  vm_size               = var.spec.vm_size

  mode    = local.mode
  os_type = local.os_type
  os_sku  = lookup(local.os_sku_map, coalesce(var.spec.os_sku, "_"), null)

  # With autoscaling, ARM owns node_count after creation; without it,
  # node_count is the pool's fixed size (0 parks the pool).
  node_count           = var.spec.node_count
  auto_scaling_enabled = var.spec.auto_scaling_enabled
  min_count            = var.spec.auto_scaling_enabled ? var.spec.min_count : null
  max_count            = var.spec.auto_scaling_enabled ? var.spec.max_count : null
  max_pods             = var.spec.max_pods

  # Spot economics: only rendered for SPOT pools. An unset max price
  # becomes -1 (pay up to the on-demand price, never price-evicted).
  priority        = local.priority != null ? local.priority : "Regular"
  eviction_policy = local.eviction_policy
  spot_max_price  = local.spot_max_price

  node_labels = length(var.spec.node_labels) > 0 ? var.spec.node_labels : null
  node_taints = length(var.spec.node_taints) > 0 ? var.spec.node_taints : null
  zones       = length(var.spec.zones) > 0 ? var.spec.zones : null

  # Unset inherits the cluster's node subnet -- correct for nearly every
  # pool; a dedicated subnet segments this pool's nodes.
  vnet_subnet_id = var.spec.vnet_subnet_id
  pod_subnet_id  = var.spec.pod_subnet_id

  orchestrator_version = var.spec.orchestrator_version

  os_disk_size_gb   = var.spec.os_disk_size_gb
  os_disk_type      = local.os_disk_type != null ? local.os_disk_type : "Managed"
  kubelet_disk_type = local.kubelet_disk_type
  ultra_ssd_enabled = var.spec.ultra_ssd_enabled

  fips_enabled            = var.spec.fips_enabled
  host_encryption_enabled = var.spec.host_encryption_enabled

  node_public_ip_enabled   = var.spec.node_public_ip_enabled
  node_public_ip_prefix_id = var.spec.node_public_ip_prefix_id

  gpu_instance = lookup(local.gpu_instance_map, coalesce(var.spec.gpu_instance, "_"), null)
  gpu_driver   = local.gpu_driver

  proximity_placement_group_id  = var.spec.proximity_placement_group_id
  host_group_id                 = var.spec.host_group_id
  capacity_reservation_group_id = var.spec.capacity_reservation_group_id

  scale_down_mode  = local.scale_down_mode != null ? local.scale_down_mode : "Delete"
  snapshot_id      = var.spec.snapshot_id
  workload_runtime = local.workload_runtime

  # A stand-in pool AKS rotates through otherwise replace-forcing changes.
  temporary_name_for_rotation = var.spec.temporary_name_for_rotation

  # Upgrade rollout: surge (extra nodes, faster) XOR unavailability
  # (in-place, no surge cost); spot pools set neither.
  dynamic "upgrade_settings" {
    for_each = var.spec.upgrade_settings != null ? [var.spec.upgrade_settings] : []
    content {
      max_surge                     = upgrade_settings.value.max_surge
      max_unavailable               = upgrade_settings.value.max_unavailable
      drain_timeout_in_minutes      = upgrade_settings.value.drain_timeout_in_minutes
      node_soak_duration_in_minutes = upgrade_settings.value.node_soak_duration_in_minutes
      undrainable_node_behavior     = lookup(local.undrainable_node_behavior_map, coalesce(upgrade_settings.value.undrainable_node_behavior, "_"), null)
    }
  }

  dynamic "kubelet_config" {
    for_each = var.spec.kubelet_config != null ? [var.spec.kubelet_config] : []
    content {
      cpu_manager_policy        = lookup(local.cpu_manager_policy_map, coalesce(kubelet_config.value.cpu_manager_policy, "_"), null)
      cpu_cfs_quota_enabled     = kubelet_config.value.cpu_cfs_quota_enabled
      cpu_cfs_quota_period      = kubelet_config.value.cpu_cfs_quota_period
      image_gc_high_threshold   = kubelet_config.value.image_gc_high_threshold
      image_gc_low_threshold    = kubelet_config.value.image_gc_low_threshold
      topology_manager_policy   = lookup(local.topology_manager_policy_map, coalesce(kubelet_config.value.topology_manager_policy, "_"), null)
      allowed_unsafe_sysctls    = length(kubelet_config.value.allowed_unsafe_sysctls) > 0 ? kubelet_config.value.allowed_unsafe_sysctls : null
      container_log_max_size_mb = kubelet_config.value.container_log_max_size_mb
      container_log_max_files   = kubelet_config.value.container_log_max_files
      pod_max_pid               = kubelet_config.value.pod_max_pid
    }
  }

  dynamic "linux_os_config" {
    for_each = var.spec.linux_os_config != null ? [var.spec.linux_os_config] : []
    content {
      transparent_huge_page        = lookup(local.transparent_huge_page_map, coalesce(linux_os_config.value.transparent_huge_page, "_"), null)
      transparent_huge_page_defrag = lookup(local.transparent_huge_page_defrag_map, coalesce(linux_os_config.value.transparent_huge_page_defrag, "_"), null)
      swap_file_size_mb            = linux_os_config.value.swap_file_size_mb

      dynamic "sysctl_config" {
        for_each = linux_os_config.value.sysctl_config != null ? [linux_os_config.value.sysctl_config] : []
        content {
          fs_aio_max_nr                      = sysctl_config.value.fs_aio_max_nr
          fs_file_max                        = sysctl_config.value.fs_file_max
          fs_inotify_max_user_watches        = sysctl_config.value.fs_inotify_max_user_watches
          fs_nr_open                         = sysctl_config.value.fs_nr_open
          kernel_threads_max                 = sysctl_config.value.kernel_threads_max
          net_core_netdev_max_backlog        = sysctl_config.value.net_core_netdev_max_backlog
          net_core_optmem_max                = sysctl_config.value.net_core_optmem_max
          net_core_rmem_default              = sysctl_config.value.net_core_rmem_default
          net_core_rmem_max                  = sysctl_config.value.net_core_rmem_max
          net_core_somaxconn                 = sysctl_config.value.net_core_somaxconn
          net_core_wmem_default              = sysctl_config.value.net_core_wmem_default
          net_core_wmem_max                  = sysctl_config.value.net_core_wmem_max
          net_ipv4_ip_local_port_range_min   = sysctl_config.value.net_ipv4_ip_local_port_range_min
          net_ipv4_ip_local_port_range_max   = sysctl_config.value.net_ipv4_ip_local_port_range_max
          net_ipv4_neigh_default_gc_thresh1  = sysctl_config.value.net_ipv4_neigh_default_gc_thresh1
          net_ipv4_neigh_default_gc_thresh2  = sysctl_config.value.net_ipv4_neigh_default_gc_thresh2
          net_ipv4_neigh_default_gc_thresh3  = sysctl_config.value.net_ipv4_neigh_default_gc_thresh3
          net_ipv4_tcp_fin_timeout           = sysctl_config.value.net_ipv4_tcp_fin_timeout
          net_ipv4_tcp_keepalive_intvl       = sysctl_config.value.net_ipv4_tcp_keepalive_intvl
          net_ipv4_tcp_keepalive_probes      = sysctl_config.value.net_ipv4_tcp_keepalive_probes
          net_ipv4_tcp_keepalive_time        = sysctl_config.value.net_ipv4_tcp_keepalive_time
          net_ipv4_tcp_max_syn_backlog       = sysctl_config.value.net_ipv4_tcp_max_syn_backlog
          net_ipv4_tcp_max_tw_buckets        = sysctl_config.value.net_ipv4_tcp_max_tw_buckets
          net_ipv4_tcp_tw_reuse              = sysctl_config.value.net_ipv4_tcp_tw_reuse
          net_netfilter_nf_conntrack_buckets = sysctl_config.value.net_netfilter_nf_conntrack_buckets
          net_netfilter_nf_conntrack_max     = sysctl_config.value.net_netfilter_nf_conntrack_max
          vm_max_map_count                   = sysctl_config.value.vm_max_map_count
          vm_swappiness                      = sysctl_config.value.vm_swappiness
          vm_vfs_cache_pressure              = sysctl_config.value.vm_vfs_cache_pressure
        }
      }
    }
  }

  dynamic "node_network_profile" {
    for_each = var.spec.node_network_profile != null ? [var.spec.node_network_profile] : []
    content {
      dynamic "allowed_host_ports" {
        for_each = node_network_profile.value.allowed_host_ports
        content {
          port_start = allowed_host_ports.value.port_start
          port_end   = allowed_host_ports.value.port_end
          protocol   = lookup(local.host_port_protocol_map, coalesce(allowed_host_ports.value.protocol, "_"), null)
        }
      }
      application_security_group_ids = length(node_network_profile.value.application_security_group_ids) > 0 ? node_network_profile.value.application_security_group_ids : null
      node_public_ip_tags            = length(node_network_profile.value.node_public_ip_tags) > 0 ? node_network_profile.value.node_public_ip_tags : null
    }
  }

  # Windows-pool outbound NAT control (Windows pools only; fixed at
  # creation).
  dynamic "windows_profile" {
    for_each = var.spec.windows_profile != null ? [var.spec.windows_profile] : []
    content {
      outbound_nat_enabled = windows_profile.value.outbound_nat_enabled
    }
  }

  tags = local.final_tags
}
