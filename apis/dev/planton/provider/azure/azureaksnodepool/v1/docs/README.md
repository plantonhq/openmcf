# AzureAksNodePool -- Design Research

## The Resource

An AKS agent pool (`Microsoft.ContainerService/managedClusters/agentPools`)
is a scale set of worker nodes attached to a managed cluster. The component
maps 1:1 onto `azurerm_kubernetes_cluster_node_pool` (azurerm v4.80,
`internal/services/containers/kubernetes_cluster_node_pool_resource.go`),
parity-verified against pulumi-azure v6.38
(`containerservice.KubernetesClusterNodePool`).

## Shape Decisions (the ones that define the component)

- **The parent is a single `kubernetes_cluster_id` ARM-id FK** (→ the
  cluster's `cluster_id` output). azurerm's own resource takes the cluster
  by ARM id; resource group and cluster name are embedded in it. The
  previous redundant `cluster_name` + `resource_group` pair could
  contradict the referenced cluster; the single id cannot.
- **`node_count` is optional (0-1000)**, replacing the invented required
  `initial_node_count > 0`. azurerm's real contract: USER pools may sit at
  0 (parked -- also how spot pools scale to zero), SYSTEM pools need at
  least 1, and with autoscaling the autoscaler owns the count.
- **Spot is modeled as azurerm models it** -- `priority` +
  `eviction_policy` + `spot_max_price` -- replacing the old `spot_enabled`
  bool that silently fixed the eviction policy and bid. Unset max price
  means -1: pay up to the on-demand price, never price-evicted (the
  setting nearly everyone wants).

## CEL Contracts (mirroring azurerm's real rules)

- Autoscaling bounds: with `auto_scaling_enabled`, `min <= max`; without
  it, both stay unset.
- SYSTEM pools are Linux and on-demand (never spot).
- `eviction_policy` and `spot_max_price` require priority SPOT; the price
  is -1 or a positive dollar amount.
- Spot pools cannot carry `upgrade_settings` surge/unavailable values
  (azurerm: spot does not support maxSurge/maxUnavailable).
- `max_surge` XOR `max_unavailable` -- mutually exclusive rollout
  strategies.
- `os_sku` must match `os_type`; Windows pool names are at most 6
  characters; `windows_profile` only on Windows pools; `linux_os_config`
  only on Linux pools.
- `node_public_ip_prefix_id` requires `node_public_ip_enabled`.

## Design Decisions

- **Both engines derive RG/cluster name from the ARM id** rather than
  asking the user to repeat them -- the established parent-child precedent
  in this catalog (subnets, peerings, DNS links).
- **Zero-node E2E proof:** the minimal scenario deploys a `node_count: 0`
  USER pool -- ARM validates the full cluster wiring and creates the agent
  pool without launching a VM, keeping the live suite fast and cheap.
- **No `mode`-specific sub-kinds.** azurerm keeps SYSTEM/USER as a field
  on one resource; so does the spec. The SYSTEM constraints are CELs, not
  a separate kind.

## Operational Behavior Worth Knowing

- Many shape changes (vm_size, os_disk_type, fips, host encryption) rotate
  the pool. `temporary_name_for_rotation` lets AKS stand up a replacement
  first -- set it proactively on production pools.
- Spot pools automatically carry the
  `kubernetes.azure.com/scalesetpriority=spot:NoSchedule` taint; AKS
  replaces evicted nodes as capacity returns. Eviction policy and max
  price are fixed at creation.
- Node Kubernetes versions may lag the control plane by up to two minors:
  `orchestrator_version` is the seam for canarying node upgrades pool by
  pool.
- Windows pools require the cluster to carry a `windows_profile`, and are
  USER mode only.

## Composition

- `kubernetes_cluster_id` → `AzureAksCluster.status.outputs.cluster_id`
- `vnet_subnet_id` / `pod_subnet_id` → `AzureSubnet.status.outputs.subnet_id`
- `node_public_ip_prefix_id` →
  `AzurePublicIpPrefix.status.outputs.public_ip_prefix_id`
- Nothing downstream deploys INTO a pool (workloads target pools via
  Kubernetes labels/taints, not ARM references) -- the outputs are the
  pool's own identifiers plus the node image actually rolled out.
