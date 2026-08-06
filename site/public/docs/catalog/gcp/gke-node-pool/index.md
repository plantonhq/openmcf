---
title: "GKE Node Pool"
description: "GKE Node Pool deployment documentation"
icon: "package"
order: 100
componentName: "gcpgkenodepool"
---

# GCP GKE Node Pool

Deploys a Google Kubernetes Engine node pool — a group of identically configured Compute Engine VMs attached to a GKE cluster. This component is the compute companion to [GcpGkeCluster](/docs/catalog/gcp/gke-cluster): the cluster manages the control plane; node pools carry the workloads, each sized, tainted, and priced for its purpose.

## What Gets Created

When you deploy a GcpGkeNodePool resource, Planton provisions:

- **GKE node pool** — a `google_container_node_pool` in the parent cluster, addressed by the cluster's `name` and `location` outputs; the Kubernetes Engine API is enabled automatically so a fresh project works on the first deploy
- **Sizing** — a fixed node count, or cluster-autoscaler management with per-zone or total bounds (including scale-to-zero)
- **Node VMs** — machine type, boot disk, OS image, local SSDs, GPUs with GKE-managed drivers, shielded/confidential settings, and CMEK boot-disk encryption as configured
- **Scheduling surface** — Kubernetes node labels and taints, GCE network tags, compact placement, and Spot/preemptible capacity
- **Platform attribution** — organization, environment, and resource labels applied as GCE resource labels on the node VMs

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GKE cluster** — a [GcpGkeCluster](/docs/catalog/gcp/gke-cluster) resource (Standard mode; Autopilot clusters take no node pools)
- **A service account** ([GcpServiceAccount](/docs/catalog/gcp/service-account)) for production pools — the Compute Engine default SA is used otherwise
- **A Cloud KMS key** ([GcpKmsKey](/docs/catalog/gcp/kms-key)) if using CMEK boot-disk encryption

## Quick Start

Create a file `node-pool.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeNodePool
metadata:
  name: general-pool
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpGkeNodePool.general-pool
spec:
  clusterName:
    valueFrom:
      kind: GcpGkeCluster
      name: my-cluster
      fieldPath: status.outputs.name
  location:
    valueFrom:
      kind: GcpGkeCluster
      name: my-cluster
      fieldPath: status.outputs.location
  autoscaling:
    minNodes: 1
    maxNodes: 3
```

Deploy:

```shell
planton apply -f node-pool.yaml
```

This creates an autoscaled pool of `e2-medium` nodes on Container-Optimized OS with auto-repair and auto-upgrade on — GKE's defaults, made explicit.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `clusterName` | `StringValueOrRef` | Parent cluster's cloud name. Reference the cluster's `name` output. Immutable. | Required |
| `location` | `StringValueOrRef` | Parent cluster's region or zone. Reference the cluster's `location` output. Immutable. | Required |

### Core Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project (must be the cluster's project). Can reference a GcpProject resource. |
| `nodePoolName` | `string` | `metadata.name` | Pool name in GKE (1–40 chars, lowercase). Immutable. |
| `nodeLocations` | `string[]` | cluster defaults | Zones this pool's nodes run in (subset of the cluster's region). Mutable. |
| `version` | `string` | GKE-managed | Explicit node Kubernetes version — pin only with `autoUpgrade` off (they fight). |
| `maxPodsPerNode` | `int` | cluster default (110) | Pods-per-node override (8–256). Immutable. |
| `initialNodeCount` | `int` | — | Starting size for autoscaled pools (per zone). Immutable. |

### Sizing

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `nodeCount` | `int` | GKE default (3) | Fixed size (per zone for regional clusters). Exclusive with `autoscaling`. |
| `autoscaling.minNodes` / `maxNodes` | `int` | — | Per-zone bounds; `minNodes: 0` allows scale-to-zero. Exclusive with total bounds. |
| `autoscaling.totalMinNodes` / `totalMaxNodes` | `int` | — | Bounds across ALL zones — an absolute cost cap regardless of zone spread. |
| `autoscaling.locationPolicy` | `string` | `BALANCED` | `BALANCED` evens out zones; `ANY` prefers reservations and cuts Spot preemption risk. |

### Lifecycle

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `management.autoRepair` | `bool` | `true` | Automatically repair unhealthy nodes. |
| `management.autoUpgrade` | `bool` | `true` | Track the control-plane Kubernetes version. |
| `upgradeSettings.maxSurge` / `maxUnavailable` | `int` | GKE default (1/0) | Surge rollout dials — how many nodes upgrade at once. |
| `upgradeSettings.strategy` | `string` | `SURGE` | `SURGE` (rolling) or `BLUE_GREEN` (full green pool, migrate, soak, delete blue). |
| `upgradeSettings.blueGreenSettings` | object | — | Batch percentage XOR node count, batch soak, and the final soak (rollback window). |
| `placementPolicy` | object | — | `COMPACT` physical co-location for HPC/ML; carries `tpuTopology` for TPU pools. Immutable. |
| `queuedProvisioningEnabled` | `bool` | `false` | Nodes obtainable only via the ProvisioningRequest API (Dynamic Workload Scheduler). Immutable. |

### Networking

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `networkConfig.createPodRange` + `podRange` + `podIpv4CidrBlock` | mixed | cluster range | Dedicated per-pool pod range — create a new one or name an existing subnetwork range. Immutable. |
| `networkConfig.enablePrivateNodes` | `bool` | cluster setting | Per-pool override for internal-only node IPs. |
| `networkConfig.totalEgressBandwidthTier` | `string` | — | `TIER_1` unlocks up to 100 Gbps; requires `gvnicEnabled`. |
| `networkConfig.podCidrOverprovisionDisabled` | `bool` | `false` | Disables the 2× per-node pod-CIDR overprovisioning. Immutable. |

### Node Configuration (`nodeConfig`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `machineType` | `string` | `e2-medium` | Compute Engine machine type. Immutable (changes replace the nodes). |
| `diskSizeGb` | `int` | GKE default (100) | Boot disk per node, min 10 GB. |
| `diskType` | `string` | GKE default (`pd-balanced`) | `pd-standard`, `pd-balanced`, `pd-ssd`, or hyperdisk variants. |
| `imageType` | `string` | `COS_CONTAINERD` | Node OS: `COS_CONTAINERD`, `UBUNTU_CONTAINERD`, `WINDOWS_LTSC_CONTAINERD`. |
| `serviceAccount` | `StringValueOrRef` | Compute default SA | Node identity (GcpServiceAccount reference). Immutable. |
| `oauthScopes` | `string[]` | GKE defaults | Node scopes; with Workload Identity the defaults are usually right. |
| `labels` | `map` | `{}` | Kubernetes node labels — what nodeSelector/affinity match. |
| `resourceLabels` | `map` | `{}` | GCE billing/inventory labels, merged beneath the platform labels. |
| `tags` | `string[]` | `[]` | GCE network tags for VPC firewall rules (additive to GKE's cluster tag). |
| `taints[]` | list | `[]` | Scheduling fences: `key`/`value`/`effect` (`NO_SCHEDULE`, `PREFER_NO_SCHEDULE`, `NO_EXECUTE`). |
| `spot` | `bool` | `false` | Spot VMs (60–91% discount, 30s preemption notice). Exclusive with `preemptible`. Immutable. |
| `preemptible` | `bool` | `false` | Legacy preemptible VMs (24h lifetime) — prefer `spot`. Immutable. |
| `guestAccelerators[]` | list | `[]` | GPUs: `type`, `count`, `gpuDriverVersion` (GKE-managed install), MIG `gpuPartitionSize`, `gpuSharingConfig`. |
| `shieldedInstanceConfig` | object | GCP defaults | `enableSecureBoot` (default false) + `enableIntegrityMonitoring` (default true). Immutable. |
| `confidentialNodes` | object | — | Hardware memory encryption: `SEV`, `SEV_SNP`, `TDX`. Immutable. |
| `bootDiskKmsKey` | `StringValueOrRef` | Google-managed | CMEK for node boot disks (GcpKmsKey reference). Immutable. |
| `minCpuPlatform` | `string` | — | CPU generation floor, e.g. `Intel Ice Lake`. Immutable. |
| `localSsdCount` | `int` | `0` | SCSI local SSDs (legacy knob — prefer the NVMe surfaces below). Immutable. |
| `ephemeralStorageLocalSsd` | object | — | Back emptyDir/logs/layers with local NVMe (+ optional GKE Data Cache count). Immutable. |
| `localNvmeSsdBlock` | object | — | Raw-block local NVMe for workloads managing their own filesystem. Immutable. |
| `gcfsEnabled` | `bool` | `false` | Image streaming — large images start minutes faster (COS only). |
| `gvnicEnabled` | `bool` | `false` | Google Virtual NIC — required for TIER_1 bandwidth. Immutable. |
| `fastSocketEnabled` | `bool` | `false` | NCCL Fast Socket for distributed GPU training. Requires gVNIC. |
| `workloadMetadataMode` | `string` | cluster default | `GKE_METADATA` (Workload Identity) or `GCE_METADATA` (legacy). |
| `reservationAffinity` | object | — | Consume Compute reservations: `ANY_RESERVATION`, `SPECIFIC_RESERVATION` (+ key/values), `NO_RESERVATION`. Immutable. |
| `secondaryBootDisks[]` | list | `[]` | Preload container images/data from disk images (`CONTAINER_IMAGE_CACHE`). Immutable. |
| `kubeletConfig` | object | GKE defaults | CPU manager, CFS quota, pod PID limit, read-only port posture, image pulls, log rotation, image GC. |
| `linuxNodeConfig` | object | GKE defaults | Allowlisted sysctls, cgroup mode, hugepages pre-allocation. |
| `loggingVariant` | `string` | `DEFAULT` | `MAX_THROUGHPUT` (10 MiB/s) for log-heavy pools. |
| `flexStart` | `bool` | `false` | Dynamic Workload Scheduler flex-start capacity (up to 7 days, discounted). Immutable. |
| `maxRunDuration` | `string` | — | Per-node max runtime (`"3600s"` format) before drain-and-delete. Immutable. |

## Examples

### Scale-to-Zero Spot Pool for Batch

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeNodePool
metadata:
  name: spot-batch
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpGkeNodePool.spot-batch
spec:
  clusterName:
    valueFrom:
      kind: GcpGkeCluster
      name: prod-primary
      fieldPath: status.outputs.name
  location:
    valueFrom:
      kind: GcpGkeCluster
      name: prod-primary
      fieldPath: status.outputs.location
  autoscaling:
    minNodes: 0
    maxNodes: 10
    locationPolicy: ANY
  nodeConfig:
    machineType: e2-standard-4
    spot: true
    labels:
      workload-class: batch
    taints:
      - key: workload-class
        value: batch
        effect: NO_SCHEDULE
```

### GPU Inference Pool

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeNodePool
metadata:
  name: gpu-inference
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpGkeNodePool.gpu-inference
spec:
  clusterName:
    valueFrom:
      kind: GcpGkeCluster
      name: prod-primary
      fieldPath: status.outputs.name
  location:
    valueFrom:
      kind: GcpGkeCluster
      name: prod-primary
      fieldPath: status.outputs.location
  autoscaling:
    minNodes: 0
    maxNodes: 4
  nodeConfig:
    machineType: g2-standard-8
    diskSizeGb: 200
    diskType: pd-ssd
    gcfsEnabled: true
    guestAccelerators:
      - type: nvidia-l4
        count: 1
        gpuDriverVersion: DEFAULT
```

### Hardened Production Pool

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeNodePool
metadata:
  name: hardened-pool
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpGkeNodePool.hardened-pool
spec:
  clusterName:
    valueFrom:
      kind: GcpGkeCluster
      name: prod-primary
      fieldPath: status.outputs.name
  location:
    valueFrom:
      kind: GcpGkeCluster
      name: prod-primary
      fieldPath: status.outputs.location
  autoscaling:
    minNodes: 2
    maxNodes: 6
  upgradeSettings:
    maxSurge: 1
    maxUnavailable: 0
  nodeConfig:
    machineType: n2d-standard-8
    serviceAccount:
      valueFrom:
        kind: GcpServiceAccount
        name: gke-nodes
        fieldPath: status.outputs.email
    shieldedInstanceConfig:
      enableSecureBoot: true
    confidentialNodes:
      enabled: true
      confidentialInstanceType: SEV
    bootDiskKmsKey:
      valueFrom:
        kind: GcpKmsKey
        name: node-disks-key
        fieldPath: status.outputs.key_id
    workloadMetadataMode: GKE_METADATA
    kubeletConfig:
      insecureKubeletReadonlyPortEnabled: "FALSE"
      podPidsLimit: 4096
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `node_pool_name` | `string` | Pool name as created in GKE |
| `instance_group_urls` | `string[]` | Managed instance group URLs backing the pool (one per zone) |
| `min_nodes` | `string` | Effective minimum size (autoscaling minimum, or the fixed count) |
| `max_nodes` | `string` | Effective maximum size (autoscaling maximum, or the fixed count) |
| `current_node_count` | `string` | Nodes per zone at the last deploy |
| `node_pool_id` | `string` | `projects/{project}/locations/{location}/clusters/{cluster}/nodePools/{name}` |
| `location` | `string` | The pool's region or zone, exactly as provided in the spec |
| `version` | `string` | Kubernetes version running on the pool's nodes |

## Related Components

- [GcpGkeCluster](/docs/catalog/gcp/gke-cluster) — the control plane this pool attaches to
- [GcpServiceAccount](/docs/catalog/gcp/service-account) — the node identity
- [GcpKmsKey](/docs/catalog/gcp/kms-key) — CMEK key for boot-disk encryption
- [GcpGkeWorkloadIdentityBinding](/docs/catalog/gcp/gke-workload-identity-binding) — grants workloads IAM identities
