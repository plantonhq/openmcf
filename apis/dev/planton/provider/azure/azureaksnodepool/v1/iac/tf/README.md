# AzureAksNodePool Terraform Module

## Overview

This Terraform module provisions an AKS node pool using the `azurerm`
provider (`~> 4.0`). It creates a single
`azurerm_kubernetes_cluster_node_pool` attached to an existing cluster by
ARM ID -- the full azurerm v4.80 surface: mode, OS, spot economics,
autoscaling, disks, GPU, kubelet/Linux OS tuning, and upgrade rollout
control.

Many shape changes (vm_size, os_disk_type, fips, host encryption) rotate
the pool; `temporary_name_for_rotation` lets AKS stand up a replacement
pool first instead of tearing this one down. Spot pools trade 30-90% cost
savings for evictability and automatically carry the scalesetpriority spot
taint.

## Resources Created

- `azurerm_kubernetes_cluster_node_pool.main` -- the node pool

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Node pool specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `kubernetes_cluster_id` | yes | ARM ID of the parent cluster |
| `name` | yes | Pool name: 1-12 lowercase alphanumerics starting with a letter (Windows: max 6) |
| `vm_size` | yes | Azure VM size, e.g. `Standard_D4s_v5` |
| `mode` | no | `USER` (default) or `SYSTEM` |
| `os_type` / `os_sku` | no | `LINUX` (default) or `WINDOWS`; image SKU |
| `node_count` | no | Fixed count 0-1000 (USER pools may park at 0); initial count with autoscaling |
| `auto_scaling_enabled` / `min_count` / `max_count` | no | Cluster-autoscaler bounds |
| `priority` / `eviction_policy` / `spot_max_price` | no | Spot economics (SPOT pools only; unset price = -1) |
| `node_labels` / `node_taints` / `zones` | no | Scheduling and zonal placement |
| `vnet_subnet_id` / `pod_subnet_id` | no | Optional dedicated subnets (unset inherits the cluster's network) |
| `orchestrator_version` | no | Node Kubernetes version (unset follows the control plane) |
| `os_disk_size_gb` / `os_disk_type` / `kubelet_disk_type` / `ultra_ssd_enabled` | no | Disk shape |
| `fips_enabled` / `host_encryption_enabled` | no | Compliance hardening (rotate the pool) |
| `node_public_ip_enabled` / `node_public_ip_prefix_id` | no | Per-node public IPs, optionally from a prefix |
| `gpu_instance` / `gpu_driver` | no | MIG profile; NVIDIA driver install control |
| `proximity_placement_group_id` / `host_group_id` / `capacity_reservation_group_id` | no | Placement |
| `scale_down_mode` / `snapshot_id` / `workload_runtime` / `temporary_name_for_rotation` | no | Lifecycle |
| `upgrade_settings` | no | `max_surge` XOR `max_unavailable`, drain timeout, soak, undrainable behavior (not valid on spot pools) |
| `kubelet_config` / `linux_os_config` / `node_network_profile` / `windows_profile` | no | Tuning blocks |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `node_pool_id` | Full ARM ID of the agent pool |
| `node_pool_name` | The pool's name as deployed |
| `node_image_version` | Node OS image version the pool is running |

## Usage

```hcl
module "app_node_pool" {
  source = "./iac/tf"

  metadata = {
    name = "general-pool"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    kubernetes_cluster_id = "/subscriptions/.../managedClusters/prod-aks"
    name                  = "general"
    vm_size               = "Standard_D4s_v5"
    auto_scaling_enabled  = true
    min_count             = 2
    max_count             = 10
  }
}
```

## Required Permissions

The deploying credential needs
`Microsoft.ContainerService/managedClusters/agentPools/write` -- held via
Azure Kubernetes Service Contributor, Contributor, or Owner. A dedicated
subnet requires the cluster identity to hold Network Contributor on it.
