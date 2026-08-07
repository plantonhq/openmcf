# Azure AKS Node Pool

Creates an AKS node pool -- a scale set of worker nodes attached to an existing AKS cluster by ARM ID -- at the full `azurerm` v4.80 surface: mode, OS, spot economics, autoscaling, disks, GPU, kubelet and Linux OS tuning, and upgrade rollout control.

## What Gets Created

When you deploy an AzureAksNodePool resource, Planton provisions:

- **Node Pool** — an `azurerm_kubernetes_cluster_node_pool` attached to the referenced cluster

Node pools are the unit of compute shape: general, memory-optimized, GPU, spot, or Windows pools each live as their own resource with an independent lifecycle. The cluster carries only its mandatory default (system) pool; everything else is one of these.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An existing AKS cluster** (an `AzureAksCluster` in composed environments -- the pool references its `cluster_id` output)
- **Container service write rights**: `Microsoft.ContainerService/managedClusters/agentPools/write`
- Windows pools require the cluster to carry a `windowsProfile`

## Quick Start

Create a file `pool.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureAksNodePool
metadata:
  name: general-pool
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureAksNodePool.general-pool
spec:
  kubernetesClusterId:
    valueFrom:
      name: prod-aks
  name: general
  vmSize: Standard_D4s_v5
  autoScalingEnabled: true
  minCount: 2
  maxCount: 10
```

Deploy:

```shell
planton apply -f pool.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `kubernetesClusterId` | `StringValueOrRef` | ARM ID of the parent cluster. Defaults to referencing an `AzureAksCluster`'s `cluster_id` output. Changing it replaces the pool. | Required |
| `name` | `string` | Pool name: 1-12 lowercase letters and numbers, starting with a letter (Windows pools: at most 6 characters). | Required |
| `vmSize` | `string` | Azure VM size, e.g. `Standard_D4s_v5` (general), `Standard_E8s_v5` (memory), `Standard_NC24ads_A100_v4` (GPU). | Required |

### Sizing and Scheduling

| Field | Type | Description |
|-------|------|-------------|
| `mode` | `enum` | `USER` (default: workloads, may park at zero, may be spot/Windows) or `SYSTEM` (cluster-critical pods; Linux, on-demand, at least 1 node). |
| `nodeCount` | `int` | Fixed count (0-1000) without autoscaling; initial count with it. USER pools may sit at 0 (parked). |
| `autoScalingEnabled` | `bool` | Cluster autoscaler owns the count between `minCount` and `maxCount`. |
| `minCount` / `maxCount` | `int` | Autoscaling bounds (0-1000; USER pools may scale to zero). |
| `maxPods` | `int` | Max pods per node. Unset = Azure's plugin-dependent default. |
| `nodeLabels` | `map(string)` | Kubernetes labels for scheduling, e.g. `{"workload": "gpu"}`. |
| `nodeTaints` | `list(string)` | Kubernetes taints, each `key=value:Effect`, e.g. `sku=gpu:NoSchedule`. |
| `zones` | `list(string)` | Availability zones (`"1"`, `"2"`, `"3"`). Empty = regional placement. |

### Spot Economics

| Field | Type | Description |
|-------|------|-------------|
| `priority` | `enum` | `REGULAR` (default) or `SPOT` -- 30-90% discount, evictable, USER pools only. Fixed at creation. |
| `evictionPolicy` | `enum` | What eviction does to a spot node: `EVICTION_DELETE` (default; billing stops) or `EVICTION_DEALLOCATE` (disks kept, bill while stopped). SPOT only. |
| `spotMaxPrice` | `double` | Ceiling price per node-hour in USD. Unset = -1: pay up to on-demand, never price-evicted. SPOT only. |

Spot pools automatically carry the `kubernetes.azure.com/scalesetpriority=spot:NoSchedule` taint and do not support upgrade surge settings.

### OS, Disks, and Placement

| Field | Type | Description |
|-------|------|-------------|
| `osType` / `osSku` | `enum` | `LINUX` (default) or `WINDOWS`; image SKU (Ubuntu, Azure Linux, Windows 2019/2022, version-pinned variants). |
| `orchestratorVersion` | `string` | Node Kubernetes version. Unset follows the control plane; may lag up to two minors for canarying. |
| `osDiskSizeGb` / `osDiskType` / `kubeletDiskType` / `ultraSsdEnabled` | | Disk shape: managed vs ephemeral OS disk, kubelet state placement, Ultra SSD. |
| `vnetSubnetId` / `podSubnetId` | `StringValueOrRef` | Optional dedicated subnets (reference `AzureSubnet`). Unset inherits the cluster's network. |
| `fipsEnabled` / `hostEncryptionEnabled` | `bool` | FIPS 140-2 images; host-based encryption. Changing either rotates the pool. |
| `nodePublicIpEnabled` / `nodePublicIpPrefixId` | | Per-node public IPs, optionally from an `AzurePublicIpPrefix`. |
| `gpuInstance` / `gpuDriver` | `enum` | MIG partitioning profile; NVIDIA driver install control. |
| `proximityPlacementGroupId` / `hostGroupId` / `capacityReservationGroupId` | `string` | Placement: co-location, dedicated hosts, reserved capacity. |
| `scaleDownMode` / `snapshotId` / `workloadRuntime` / `temporaryNameForRotation` | | Scale-down disposition, snapshot source, OCI/Kata/WASM runtime, rotation stand-in name. |

### Tuning Blocks

| Field | Description |
|-------|-------------|
| `upgradeSettings` | Rollout control: `maxSurge` XOR `maxUnavailable`, drain timeout, node soak, undrainable-node behavior. Not valid on spot pools. |
| `kubeletConfig` | CPU manager/CFS quota, image GC thresholds, topology manager, unsafe sysctls allowlist, container log caps, pod PID limit. |
| `linuxOsConfig` | Transparent hugepages, swap file, and the full sysctl surface. Linux pools only. |
| `nodeNetworkProfile` | Allowed host ports, application security groups, node public IP tags. |
| `windowsProfile` | `outboundNatEnabled` (Windows pools only; fixed at creation). |
| `tags` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Spot Pool for Fault-Tolerant Workloads

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureAksNodePool
metadata:
  name: spot-pool
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureAksNodePool.spot-pool
spec:
  kubernetesClusterId:
    valueFrom:
      name: prod-aks
  name: spot
  vmSize: Standard_D4s_v5
  priority: SPOT
  evictionPolicy: EVICTION_DELETE
  autoScalingEnabled: true
  minCount: 0
  maxCount: 20
```

### GPU Pool Reserved by Taint

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureAksNodePool
metadata:
  name: gpu-pool
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureAksNodePool.gpu-pool
spec:
  kubernetesClusterId:
    valueFrom:
      name: prod-aks
  name: gpu
  vmSize: Standard_NC4as_T4_v3
  nodeCount: 1
  gpuDriver: INSTALL
  nodeLabels:
    workload: gpu
  nodeTaints:
    - "sku=gpu:NoSchedule"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `node_pool_id` | `string` | Full ARM ID of the agent pool |
| `node_pool_name` | `string` | The pool's name (surfaces in the `kubernetes.azure.com/agentpool` node label) |
| `node_image_version` | `string` | Node OS image version actually rolled out -- audit patch currency across pools |

## Related Components

- [AzureAksCluster](/docs/catalog/azure/aks-cluster) — the parent cluster (carries only the mandatory default pool)
- [AzureSubnet](/docs/catalog/azure/subnet) — optional dedicated subnet for this pool's nodes
- [AzurePublicIpPrefix](/docs/catalog/azure/public-ip-prefix) — allowlistable CIDR for node public IPs
