# GcpGkeNodePool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `gcp.planton.dev/v1alpha1`

GcpGkeNodePoolSpec defines a GKE node pool (`google_container_node_pool`) —
a group of Compute Engine VMs with one shared configuration, attached to a
GcpGkeCluster.

Node pools are the compute half of the GKE composition boundary: the
cluster owns the control plane and cluster-wide configuration; every node
pool is a first-class resource with its own lifecycle, referencing the
cluster by its name output. Production clusters run several pools —
general-purpose on-demand, scale-to-zero Spot for fault-tolerant batch,
GPU pools for ML — each sized and tainted for its workloads. Autopilot
clusters take no node pools at all (GKE manages nodes).

The node pool inherits its project and location from the parent cluster;
both resolve by reference so a manifest never has to repeat what the
cluster already knows.

## Example

```yaml
# Development manifest for GcpGkeNodePool — exercises a broad slice of the
# spec (autoscaling, spot, taints, labels, upgrade settings, kubelet tuning)
# against an existing cluster.
#
# Usage: planton tofu plan --manifest catalog/gcp/gcpgkenodepool/e2e/manifest.yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpGkeNodePool
metadata:
  name: hack-node-pool
  id: gkenp-hack-001
  org: planton-dev
  env: dev
spec:
  # project_id omitted: the pool lands in the provider's default project.
  clusterName:
    value: hack-gke-cluster
  location:
    value: us-central1
  nodePoolName: hack-batch-pool
  nodeLocations:
    - us-central1-a
    - us-central1-b
  autoscaling:
    minNodes: 0
    maxNodes: 5
    locationPolicy: ANY
  management:
    autoRepair: true
    autoUpgrade: true
  upgradeSettings:
    maxSurge: 2
    maxUnavailable: 1
    strategy: SURGE
  nodeConfig:
    machineType: n2-standard-4
    diskSizeGb: 100
    diskType: pd-balanced
    imageType: COS_CONTAINERD
    spot: true
    labels:
      workload-class: batch
    taints:
      - key: workload-class
        value: batch
        effect: NO_SCHEDULE
    tags:
      - batch-nodes
    gcfsEnabled: true
    gvnicEnabled: true
    shieldedInstanceConfig:
      enableSecureBoot: true
    kubeletConfig:
      podPidsLimit: 4096
      insecureKubeletReadonlyPortEnabled: "FALSE"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.projectId` | `string \| valueFrom` |  |  | GcpProject (`status.outputs.project_id`) |
| `spec.clusterName` | `string \| valueFrom` | yes |  | GcpGkeCluster (`status.outputs.name`) |
| `spec.location` | `string \| valueFrom` | yes |  | GcpGkeCluster (`status.outputs.location`) |
| `spec.nodePoolName` | `string` |  |  |  |
| `spec.nodeLocations` | `[]string` |  |  |  |
| `spec.version` | `string` |  |  |  |
| `spec.maxPodsPerNode` | `int32` |  |  |  |
| `spec.initialNodeCount` | `uint32` |  |  |  |
| `spec.nodeCount` | `uint32` |  |  |  |
| `spec.autoscaling` | `GcpGkeNodePoolAutoscaling` |  |  |  |
| `spec.autoscaling.minNodes` | `uint32` |  |  |  |
| `spec.autoscaling.maxNodes` | `uint32` |  |  |  |
| `spec.autoscaling.totalMinNodes` | `uint32` |  |  |  |
| `spec.autoscaling.totalMaxNodes` | `uint32` |  |  |  |
| `spec.autoscaling.locationPolicy` | `string` |  |  |  |
| `spec.management` | `GcpGkeNodePoolManagement` |  |  |  |
| `spec.management.autoRepair` | `bool` |  | `true` |  |
| `spec.management.autoUpgrade` | `bool` |  | `true` |  |
| `spec.upgradeSettings` | `GcpGkeNodePoolUpgradeSettings` |  |  |  |
| `spec.upgradeSettings.maxSurge` | `uint32` |  |  |  |
| `spec.upgradeSettings.maxUnavailable` | `uint32` |  |  |  |
| `spec.upgradeSettings.strategy` | `string` |  |  |  |
| `spec.upgradeSettings.blueGreenSettings` | `GcpGkeNodePoolBlueGreenSettings` |  |  |  |
| `spec.upgradeSettings.blueGreenSettings.standardRolloutPolicy` | `GcpGkeNodePoolStandardRolloutPolicy` | yes |  |  |
| `spec.upgradeSettings.blueGreenSettings.standardRolloutPolicy.batchPercentage` | `float` |  |  |  |
| `spec.upgradeSettings.blueGreenSettings.standardRolloutPolicy.batchNodeCount` | `uint32` |  |  |  |
| `spec.upgradeSettings.blueGreenSettings.standardRolloutPolicy.batchSoakDuration` | `string` |  |  |  |
| `spec.upgradeSettings.blueGreenSettings.nodePoolSoakDuration` | `string` |  |  |  |
| `spec.placementPolicy` | `GcpGkeNodePoolPlacementPolicy` |  |  |  |
| `spec.placementPolicy.type` | `string` | yes |  |  |
| `spec.placementPolicy.policyName` | `string` |  |  |  |
| `spec.placementPolicy.tpuTopology` | `string` |  |  |  |
| `spec.queuedProvisioningEnabled` | `bool` |  |  |  |
| `spec.networkConfig` | `GcpGkeNodePoolNetworkConfig` |  |  |  |
| `spec.networkConfig.createPodRange` | `bool` |  |  |  |
| `spec.networkConfig.podRange` | `string` |  |  |  |
| `spec.networkConfig.podIpv4CidrBlock` | `string` |  |  |  |
| `spec.networkConfig.enablePrivateNodes` | `bool` |  |  |  |
| `spec.networkConfig.totalEgressBandwidthTier` | `string` |  |  |  |
| `spec.networkConfig.podCidrOverprovisionDisabled` | `bool` |  |  |  |
| `spec.nodeConfig` | `GcpGkeNodePoolNodeConfig` |  |  |  |
| `spec.nodeConfig.machineType` | `string` |  | `e2-medium` |  |
| `spec.nodeConfig.diskSizeGb` | `uint32` |  |  |  |
| `spec.nodeConfig.diskType` | `string` |  |  |  |
| `spec.nodeConfig.imageType` | `string` |  | `COS_CONTAINERD` |  |
| `spec.nodeConfig.serviceAccount` | `string \| valueFrom` |  |  | GcpServiceAccount (`status.outputs.email`) |
| `spec.nodeConfig.oauthScopes` | `[]string` |  |  |  |
| `spec.nodeConfig.labels` | `map<string, string>` |  |  |  |
| `spec.nodeConfig.resourceLabels` | `map<string, string>` |  |  |  |
| `spec.nodeConfig.tags` | `[]string` |  |  |  |
| `spec.nodeConfig.metadata` | `map<string, string>` |  |  |  |
| `spec.nodeConfig.taints` | `[]GcpGkeNodePoolTaint` |  |  |  |
| `spec.nodeConfig.taints[].key` | `string` | yes |  |  |
| `spec.nodeConfig.taints[].value` | `string` | yes |  |  |
| `spec.nodeConfig.taints[].effect` | `string` | yes |  |  |
| `spec.nodeConfig.spot` | `bool` |  |  |  |
| `spec.nodeConfig.preemptible` | `bool` |  |  |  |
| `spec.nodeConfig.guestAccelerators` | `[]GcpGkeNodePoolGuestAccelerator` |  |  |  |
| `spec.nodeConfig.guestAccelerators[].type` | `string` | yes |  |  |
| `spec.nodeConfig.guestAccelerators[].count` | `uint32` |  |  |  |
| `spec.nodeConfig.guestAccelerators[].gpuPartitionSize` | `string` |  |  |  |
| `spec.nodeConfig.guestAccelerators[].gpuDriverVersion` | `string` |  |  |  |
| `spec.nodeConfig.guestAccelerators[].gpuSharingConfig` | `GcpGkeNodePoolGpuSharingConfig` |  |  |  |
| `spec.nodeConfig.guestAccelerators[].gpuSharingConfig.gpuSharingStrategy` | `string` | yes |  |  |
| `spec.nodeConfig.guestAccelerators[].gpuSharingConfig.maxSharedClientsPerGpu` | `uint32` |  |  |  |
| `spec.nodeConfig.shieldedInstanceConfig` | `GcpGkeNodePoolShieldedInstanceConfig` |  |  |  |
| `spec.nodeConfig.shieldedInstanceConfig.enableSecureBoot` | `bool` |  |  |  |
| `spec.nodeConfig.shieldedInstanceConfig.enableIntegrityMonitoring` | `bool` |  | `true` |  |
| `spec.nodeConfig.confidentialNodes` | `GcpGkeNodePoolConfidentialNodes` |  |  |  |
| `spec.nodeConfig.confidentialNodes.enabled` | `bool` |  |  |  |
| `spec.nodeConfig.confidentialNodes.confidentialInstanceType` | `string` |  |  |  |
| `spec.nodeConfig.minCpuPlatform` | `string` |  |  |  |
| `spec.nodeConfig.localSsdCount` | `uint32` |  |  |  |
| `spec.nodeConfig.ephemeralStorageLocalSsd` | `GcpGkeNodePoolEphemeralStorageLocalSsd` |  |  |  |
| `spec.nodeConfig.ephemeralStorageLocalSsd.localSsdCount` | `uint32` |  |  |  |
| `spec.nodeConfig.ephemeralStorageLocalSsd.dataCacheCount` | `uint32` |  |  |  |
| `spec.nodeConfig.localNvmeSsdBlock` | `GcpGkeNodePoolLocalNvmeSsdBlock` |  |  |  |
| `spec.nodeConfig.localNvmeSsdBlock.localSsdCount` | `uint32` |  |  |  |
| `spec.nodeConfig.gcfsEnabled` | `bool` |  |  |  |
| `spec.nodeConfig.gvnicEnabled` | `bool` |  |  |  |
| `spec.nodeConfig.fastSocketEnabled` | `bool` |  |  |  |
| `spec.nodeConfig.bootDiskKmsKey` | `string \| valueFrom` |  |  | GcpKmsKey (`status.outputs.key_id`) |
| `spec.nodeConfig.workloadMetadataMode` | `string` |  |  |  |
| `spec.nodeConfig.reservationAffinity` | `GcpGkeNodePoolReservationAffinity` |  |  |  |
| `spec.nodeConfig.reservationAffinity.consumeReservationType` | `string` | yes |  |  |
| `spec.nodeConfig.reservationAffinity.key` | `string` |  |  |  |
| `spec.nodeConfig.reservationAffinity.values` | `[]string` |  |  |  |
| `spec.nodeConfig.secondaryBootDisks` | `[]GcpGkeNodePoolSecondaryBootDisk` |  |  |  |
| `spec.nodeConfig.secondaryBootDisks[].diskImage` | `string` | yes |  |  |
| `spec.nodeConfig.secondaryBootDisks[].mode` | `string` |  |  |  |
| `spec.nodeConfig.kubeletConfig` | `GcpGkeNodePoolKubeletConfig` |  |  |  |
| `spec.nodeConfig.kubeletConfig.cpuManagerPolicy` | `string` |  |  |  |
| `spec.nodeConfig.kubeletConfig.cpuCfsQuota` | `bool` |  |  |  |
| `spec.nodeConfig.kubeletConfig.cpuCfsQuotaPeriod` | `string` |  |  |  |
| `spec.nodeConfig.kubeletConfig.podPidsLimit` | `int64` |  |  |  |
| `spec.nodeConfig.kubeletConfig.insecureKubeletReadonlyPortEnabled` | `string` |  |  |  |
| `spec.nodeConfig.kubeletConfig.maxParallelImagePulls` | `int64` |  |  |  |
| `spec.nodeConfig.kubeletConfig.containerLogMaxSize` | `string` |  |  |  |
| `spec.nodeConfig.kubeletConfig.containerLogMaxFiles` | `int64` |  |  |  |
| `spec.nodeConfig.kubeletConfig.imageGcLowThresholdPercent` | `int64` |  |  |  |
| `spec.nodeConfig.kubeletConfig.imageGcHighThresholdPercent` | `int64` |  |  |  |
| `spec.nodeConfig.linuxNodeConfig` | `GcpGkeNodePoolLinuxNodeConfig` |  |  |  |
| `spec.nodeConfig.linuxNodeConfig.sysctls` | `map<string, string>` |  |  |  |
| `spec.nodeConfig.linuxNodeConfig.cgroupMode` | `string` |  |  |  |
| `spec.nodeConfig.linuxNodeConfig.hugepagesConfig` | `GcpGkeNodePoolHugepagesConfig` |  |  |  |
| `spec.nodeConfig.linuxNodeConfig.hugepagesConfig.hugepageSize2m` | `int64` |  |  |  |
| `spec.nodeConfig.linuxNodeConfig.hugepagesConfig.hugepageSize1g` | `int64` |  |  |  |
| `spec.nodeConfig.loggingVariant` | `string` |  |  |  |
| `spec.nodeConfig.flexStart` | `bool` |  |  |  |
| `spec.nodeConfig.maxRunDuration` | `string` |  |  |  |

## Field Details

### spec.projectId

`string | valueFrom`

The GCP project the node pool is created in. Must be the parent
cluster's project — node pools cannot live in a different project than
their cluster. Accepts a literal project ID or a reference to a
GcpProject resource. If omitted, the provider's default project is used.

- references: GcpProject (`status.outputs.project_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpProject, name: <that resource's name>, fieldPath: status.outputs.project_id}} -- a bare string does not parse

### spec.clusterName

`string | valueFrom` · required

Name of the parent GKE cluster as created in GCP. Resolves from the
cluster's name output — the name GCP actually assigned — so it stays
correct even when the cluster's cloud name differs from its Planton
metadata.name. Immutable.

- references: GcpGkeCluster (`status.outputs.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpGkeCluster, name: <that resource's name>, fieldPath: status.outputs.name}} -- a bare string does not parse

### spec.location

`string | valueFrom` · required

Location of the parent cluster — a region ("us-central1") for regional
clusters or a zone ("us-central1-a") for zonal ones. Must match the
cluster's own location; resolves from the cluster's location output by
default. Immutable. For a regional cluster, size limits in autoscaling
are PER ZONE (see autoscaling), and the pool gets one managed instance
group per zone.

- references: GcpGkeCluster (`status.outputs.location`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpGkeCluster, name: <that resource's name>, fieldPath: status.outputs.location}} -- a bare string does not parse

### spec.nodePoolName

`string`

Name of the node pool in GKE. Immutable. If not specified, defaults to
metadata.name. Must be 1-40 characters: lowercase letters, digits, and
hyphens; starting with a letter and ending with a letter or digit.
Example: "general-pool", "spot-batch", "gpu-a100"

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^[a-z]([a-z0-9-]{0,38}[a-z0-9])?$"}}

### spec.nodeLocations

`[]string`

Zones this pool's nodes run in, e.g. ["us-central1-a", "us-central1-b"].
Must be within the cluster's region. If unspecified, the cluster-level
node_locations apply (all zones in the region for a regional cluster).
Mutable. Narrowing a pool to fewer zones cuts cost and inter-zone
traffic for workloads that tolerate zonal risk.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","repeated":{"unique":true,"items":{"string":{"pattern":"^[a-z]+-[a-z]+[0-9]+-[a-z]$"}}}}

### spec.version

`string`

Kubernetes version for the nodes. If empty (recommended), GKE picks the
version and auto-upgrade keeps it current. Setting an explicit version
while management.auto_upgrade is on makes the two fight — pin a version
only with auto-upgrade off, and expect to own upgrades manually from
then on.

### spec.maxPodsPerNode

`int32` · optional (explicit presence)

Maximum pods per node in this pool (8-256), overriding the cluster's
default (110). Lower values shrink the per-node pod CIDR slice so the
pod range stretches across more nodes. Immutable. Only effective on
VPC-native clusters (which every Planton GKE cluster is).

- rule: {"int32":{"lte":256,"gte":8}}

### spec.initialNodeCount

`uint32` · optional (explicit presence)

Number of nodes the pool starts with when autoscaling manages the size
(per zone for regional clusters). Only meaningful with autoscaling —
a fixed-size pool's size IS node_count. Immutable: changing it forces
pool recreation, so set it once and let the autoscaler own the size
from then on.

### spec.nodeCount

`uint32`

Fixed number of nodes (per zone for regional clusters). The pool
stays at this size until you change it.

### spec.autoscaling

`GcpGkeNodePoolAutoscaling`

Cluster-autoscaler management of the pool size between bounds.

- rule: use per-zone limits (min_nodes/max_nodes) or total limits (total_min_nodes/total_max_nodes), not both — they are mutually exclusive addressing modes
- rule: autoscaling needs a maximum: set max_nodes (per-zone mode) or total_max_nodes (total mode)
- rule: min_nodes cannot exceed max_nodes
- rule: total_min_nodes cannot exceed total_max_nodes

### spec.autoscaling.minNodes

`uint32` · optional (explicit presence)

Minimum nodes PER ZONE. 0 allows scale-to-zero — the pattern for Spot
and GPU pools that should cost nothing while idle.

### spec.autoscaling.maxNodes

`uint32` · optional (explicit presence)

Maximum nodes PER ZONE. A regional cluster in 3 zones with max_nodes=4
can reach 12 nodes.

### spec.autoscaling.totalMinNodes

`uint32` · optional (explicit presence)

Minimum nodes across ALL zones (total addressing mode).

### spec.autoscaling.totalMaxNodes

`uint32` · optional (explicit presence)

Maximum nodes across ALL zones (total addressing mode) — an absolute
cost cap regardless of zone spread.

### spec.autoscaling.locationPolicy

`string`

Scale-up algorithm: BALANCED spreads new nodes to even out zone sizes;
ANY prioritizes unused reservations and reduces Spot preemption risk —
prefer ANY for Spot pools.

- rule: location_policy must be empty, BALANCED, or ANY

### spec.management

`GcpGkeNodePoolManagement`

Auto-repair and auto-upgrade. If omitted, both default to true — GKE's
own defaults and the posture release channels expect.

### spec.management.autoRepair

`bool` · optional (explicit presence)

Automatically repair nodes that fail health checks. Keep on; a pool of
broken nodes heals itself instead of paging you.

- default: `true`

### spec.management.autoUpgrade

`bool` · optional (explicit presence)

Automatically upgrade node Kubernetes versions to track the control
plane. Keep on unless you pin `version` and own upgrades manually
(required for pools on a release channel).

- default: `true`

### spec.upgradeSettings

`GcpGkeNodePoolUpgradeSettings`

How node upgrades roll through the pool: surge (default — add
max_surge nodes, drain max_unavailable at a time) or blue-green
(provision a full green pool, shift workloads, soak, delete blue).
If omitted, GKE's default surge settings (max_surge=1,
max_unavailable=0) apply.

- rule: blue_green_settings apply only when strategy is BLUE_GREEN
- rule: max_surge/max_unavailable apply to the SURGE strategy — remove them when strategy is BLUE_GREEN

### spec.upgradeSettings.maxSurge

`uint32` · optional (explicit presence)

Additional nodes added during a surge upgrade (0 or more). Higher means
faster upgrades at temporary extra cost. max_surge + max_unavailable
must be at least 1 and at most 20 combined.

### spec.upgradeSettings.maxUnavailable

`uint32` · optional (explicit presence)

Nodes that may be simultaneously unavailable during a surge upgrade.
Higher trades workload disruption for speed.

### spec.upgradeSettings.strategy

`string`

SURGE (default): rolling replacement, cheapest. BLUE_GREEN: provision a
complete new node set, migrate, soak, then delete the old one — safest
rollback story, double capacity while in flight.

- rule: strategy must be empty, SURGE, or BLUE_GREEN

### spec.upgradeSettings.blueGreenSettings

`GcpGkeNodePoolBlueGreenSettings`

Blue-green rollout pacing. Only with strategy BLUE_GREEN.

### spec.upgradeSettings.blueGreenSettings.standardRolloutPolicy

`GcpGkeNodePoolStandardRolloutPolicy` · required

How the blue pool drains, batch by batch.

- rule: {"required":true}
- rule: set batch_percentage or batch_node_count, not both

### spec.upgradeSettings.blueGreenSettings.standardRolloutPolicy.batchPercentage

`float` · optional (explicit presence)

Fraction of blue nodes drained per batch, 0.0-1.0.

- rule: {"float":{"lte":1,"gte":0}}

### spec.upgradeSettings.blueGreenSettings.standardRolloutPolicy.batchNodeCount

`uint32` · optional (explicit presence)

Number of blue nodes drained per batch.

### spec.upgradeSettings.blueGreenSettings.standardRolloutPolicy.batchSoakDuration

`string`

Soak time after each batch drains before the next starts. Duration in
seconds format, e.g. "600s".

- rule: batch_soak_duration must be a seconds-format duration like "600s"

### spec.upgradeSettings.blueGreenSettings.nodePoolSoakDuration

`string`

Time after the entire blue pool is drained before it is deleted —
the rollback window. Duration in seconds format, e.g. "3600s".

- rule: node_pool_soak_duration must be a seconds-format duration like "3600s"

### spec.placementPolicy

`GcpGkeNodePoolPlacementPolicy`

Compact placement: co-locates nodes physically for low inter-node
latency (tightly-coupled HPC/ML workloads) — and carries the TPU
topology for TPU pools. Immutable.

### spec.placementPolicy.type

`string` · required

COMPACT places nodes close together for minimal inter-node latency —
for tightly-coupled HPC and multi-node ML training. Requires machine
families that support compact placement (C2/C3/A2/A3...).

- rule: {"required":true,"string":{"in":["COMPACT"]}}

### spec.placementPolicy.policyName

`string`

Optional user-supplied compute resource policy to place nodes under.
Must be in the same project and region as the pool. When empty, GKE
creates and owns the placement policy.

### spec.placementPolicy.tpuTopology

`string`

TPU placement topology, e.g. "2x2x2" — TPU pools only.
https://cloud.google.com/kubernetes-engine/docs/concepts/plan-tpus#topology

### spec.queuedProvisioningEnabled

`bool`

Nodes are obtainable only through the ProvisioningRequest API (Dynamic
Workload Scheduler queued provisioning) — the pool provisions capacity
in whole batches when it becomes available, for large atomic workloads
like multi-node training jobs. Immutable.

### spec.networkConfig

`GcpGkeNodePoolNetworkConfig`

Pool-level networking overrides (pod range, private nodes, network
performance). If omitted, the cluster-level defaults apply.

- rule: create_pod_range names the NEW range to create — set pod_range (the range name) alongside it, or set neither to use the cluster's pod range

### spec.networkConfig.createPodRange

`bool`

Create a NEW secondary range for this pool's pods (named by pod_range,
sized by pod_ipv4_cidr_block) instead of drawing from the cluster's pod
range. Immutable. Dedicated pod ranges isolate address planning per
pool — useful when one pool's churn would exhaust the shared range.

### spec.networkConfig.podRange

`string`

The secondary range for this pool's pod IPs. With create_pod_range,
the name given to the new range; without it, the name of an EXISTING
secondary range on the cluster's subnetwork. Immutable.

### spec.networkConfig.podIpv4CidrBlock

`string`

CIDR (e.g. "10.96.0.0/14") or netmask size (e.g. "/14") for the new pod
range when create_pod_range is set; empty lets GKE choose. Immutable.

### spec.networkConfig.enablePrivateNodes

`bool` · optional (explicit presence)

Whether this pool's nodes get only internal IPs, overriding the
cluster-level private-nodes setting. Unset inherits from the cluster.

### spec.networkConfig.totalEgressBandwidthTier

`string`

Network bandwidth tier: TIER_1 unlocks up to 100 Gbps total egress on
supported machine families (N2/N2D/C2/C3...). Requires gVNIC on this
pool (node_config.gvnic_enabled).

- rule: total_egress_bandwidth_tier must be empty, TIER_UNSPECIFIED, or TIER_1

### spec.networkConfig.podCidrOverprovisionDisabled

`bool`

Disables the pod CIDR overprovisioning for this pool (GKE normally
doubles the pod range slice per node). Only relevant with a dedicated
pod range. Immutable.

### spec.nodeConfig

`GcpGkeNodePoolNodeConfig`

The node VM configuration: machine type, disks, identity, scheduling
constraints, accelerators, security posture, and kubelet/OS tuning.
If omitted entirely, GKE defaults apply (e2-medium, 100 GB pd-balanced,
Container-Optimized OS, Compute Engine default service account).

- rule: spot and preemptible are mutually exclusive — spot is the current model (no 24h max lifetime); preemptible is the legacy one
- rule: fast_socket_enabled requires gvnic_enabled — NCCL Fast Socket rides on gVNIC

### spec.nodeConfig.machineType

`string` · optional (explicit presence)

Compute Engine machine type, e.g. "e2-medium", "n2-standard-8",
"a2-highgpu-1g". Defaults to e2-medium (2 vCPU, 4 GB) — fine for
sandboxes, undersized for real workloads. Immutable (changing it
replaces the pool's nodes).

- default: `e2-medium`

### spec.nodeConfig.diskSizeGb

`uint32` · optional (explicit presence)

Boot disk size in GB per node (min 10). GKE's default is 100.

- rule: {"uint32":{"gte":10}}

### spec.nodeConfig.diskType

`string`

Boot disk type. pd-balanced (GKE's current default) suits most pools;
pd-ssd for I/O-heavy node-local work; hyperdisk-balanced on machine
families that require it (C3/C3D and newer).

- rule: disk_type must be empty, pd-standard, pd-balanced, pd-ssd, hyperdisk-balanced, hyperdisk-extreme, or hyperdisk-throughput

### spec.nodeConfig.imageType

`string` · optional (explicit presence)

Node OS image. COS_CONTAINERD (Container-Optimized OS, the default and
GKE's recommendation), UBUNTU_CONTAINERD (when you need Ubuntu
packages/kernel modules), or WINDOWS_LTSC_CONTAINERD for Windows
workloads.

- default: `COS_CONTAINERD`
- rule: image_type must be COS_CONTAINERD, UBUNTU_CONTAINERD, or WINDOWS_LTSC_CONTAINERD

### spec.nodeConfig.serviceAccount

`string | valueFrom`

Service account the node VMs run as. Accepts an email literal or a
reference to a GcpServiceAccount resource. Defaults to the Compute
Engine default SA — create a minimal dedicated SA for production and
grant workload permissions through Workload Identity instead of node
scopes.

- references: GcpServiceAccount (`status.outputs.email`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpServiceAccount, name: <that resource's name>, fieldPath: status.outputs.email}} -- a bare string does not parse

### spec.nodeConfig.oauthScopes

`[]string`

OAuth scopes on the node VMs. Empty applies GKE's defaults
(devstorage.read_only, logging.write, monitoring). With Workload
Identity (the Planton cluster default), workload permissions come from
IAM on Kubernetes service accounts — node scopes only gate node-level
agents, so the defaults are usually right.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.nodeConfig.labels

`map<string, string>`

Kubernetes labels applied to every node in the pool — what nodeSelector
and affinity rules match on. Example: {"workload-class": "batch"}.

### spec.nodeConfig.resourceLabels

`map<string, string>`

GCE resource labels on the node VMs (cloud billing/inventory labels,
not Kubernetes labels). Merged with the standard platform labels;
platform attribution keys win on conflict.

### spec.nodeConfig.tags

`[]string`

GCE network tags on the node VMs — what VPC firewall rules match.
GKE adds its own cluster tag automatically; entries here are additive.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.nodeConfig.metadata

`map<string, string>`

GCE instance metadata key/value pairs. GKE requires
"disable-legacy-endpoints" = "true" and both engines enforce it
beneath any entries set here. Immutable.

### spec.nodeConfig.taints

`[]GcpGkeNodePoolTaint`

Kubernetes taints applied to every node — the scheduling fence that
keeps general workloads off special-purpose pools. Pair each taint
with a matching toleration on the intended workloads. GPU pools get
an automatic nvidia.com/gpu taint from GKE.

### spec.nodeConfig.taints[].key

`string` · required

Taint key, e.g. "workload-class".

- rule: {"required":true}

### spec.nodeConfig.taints[].value

`string` · required

Taint value, e.g. "batch".

- rule: {"required":true}

### spec.nodeConfig.taints[].effect

`string` · required

NO_SCHEDULE fences new pods without a toleration; PREFER_NO_SCHEDULE
is advisory; NO_EXECUTE also evicts running pods without a toleration.

- rule: {"required":true,"string":{"in":["NO_SCHEDULE","PREFER_NO_SCHEDULE","NO_EXECUTE"]}}

### spec.nodeConfig.spot

`bool`

Spot VMs: deeply discounted (60-91%) capacity with no availability
guarantee — nodes can be preempted with 30s notice. The current model
(no 24-hour max lifetime; replaces preemptible). Pair with
scale-to-zero autoscaling, ANY location policy, and taints so only
fault-tolerant workloads land here. Immutable.

### spec.nodeConfig.preemptible

`bool`

Legacy preemptible VMs (24-hour max lifetime). Prefer spot for new
pools. Immutable.

### spec.nodeConfig.guestAccelerators

`[]GcpGkeNodePoolGuestAccelerator`

GPU accelerators attached to every node. The machine type must support
the accelerator (or be an accelerator-optimized family like A2/A3/G2,
where the GPU is implied by the machine type and this block must match
it). GKE taints GPU nodes automatically.

### spec.nodeConfig.guestAccelerators[].type

`string` · required

Accelerator type resource name, e.g. "nvidia-tesla-t4",
"nvidia-l4", "nvidia-a100-80gb". Must be available in the pool's
zones.

- rule: {"required":true}

### spec.nodeConfig.guestAccelerators[].count

`uint32`

Number of accelerator cards per node.

- rule: {"uint32":{"gte":1}}

### spec.nodeConfig.guestAccelerators[].gpuPartitionSize

`string`

NVIDIA MIG partition size, e.g. "1g.5gb" — slices one physical GPU
into isolated instances (A100/H100 families).

### spec.nodeConfig.guestAccelerators[].gpuDriverVersion

`string`

Driver installation: DEFAULT (GKE installs the default driver
version), LATEST (newest available; COS only), or
INSTALLATION_DISABLED (bring your own DaemonSet). If omitted on GKE
1.30.1+, DEFAULT applies.

- rule: gpu_driver_version must be empty, INSTALLATION_DISABLED, DEFAULT, or LATEST

### spec.nodeConfig.guestAccelerators[].gpuSharingConfig

`GcpGkeNodePoolGpuSharingConfig`

GPU sharing: lets multiple pods share one physical GPU.

### spec.nodeConfig.guestAccelerators[].gpuSharingConfig.gpuSharingStrategy

`string` · required

GPU_TIME_SHARING context-switches the GPU between pods; MPS (NVIDIA
Multi-Process Service) runs them concurrently with resource limits.

- rule: {"required":true,"string":{"in":["GPU_TIME_SHARING","MPS"]}}

### spec.nodeConfig.guestAccelerators[].gpuSharingConfig.maxSharedClientsPerGpu

`uint32`

Maximum pods sharing each physical GPU.

- rule: {"uint32":{"gte":1}}

### spec.nodeConfig.shieldedInstanceConfig

`GcpGkeNodePoolShieldedInstanceConfig`

Shielded VM options. GKE's defaults: secure boot off (to tolerate
unsigned third-party kernel modules), integrity monitoring on. Enable
secure boot unless a workload loads unsigned modules. Immutable.

### spec.nodeConfig.shieldedInstanceConfig.enableSecureBoot

`bool`

Verify boot components against a signature baseline. GCP default
false — because unsigned third-party kernel modules fail secure boot.
Turn it on unless you load such modules.

### spec.nodeConfig.shieldedInstanceConfig.enableIntegrityMonitoring

`bool` · optional (explicit presence)

Monitor and attest boot integrity at runtime (GCP default true).

- default: `true`

### spec.nodeConfig.confidentialNodes

`GcpGkeNodePoolConfidentialNodes`

Confidential GKE nodes: hardware memory encryption (AMD SEV / Intel
TDX) for the node VMs. Requires a supporting machine family (N2D/C2D/
C3D...). Immutable.

### spec.nodeConfig.confidentialNodes.enabled

`bool`

Whether confidential nodes are enabled for this pool.

### spec.nodeConfig.confidentialNodes.confidentialInstanceType

`string`

Confidential computing technology: SEV (AMD, the common choice),
SEV_SNP, or TDX (Intel). Empty lets GCP choose for the machine family.

- rule: confidential_instance_type must be empty, SEV, SEV_SNP, or TDX

### spec.nodeConfig.minCpuPlatform

`string`

Minimum CPU platform, e.g. "Intel Ice Lake" — pins scheduling to that
CPU generation or newer for instruction-set or performance floors.
Immutable.

### spec.nodeConfig.localSsdCount

`uint32` · optional (explicit presence)

Local SSDs attached to each node for scratch I/O, automatically
formatted and mounted by GKE (SCSI interface; legacy knob). For NVMe,
prefer ephemeral_storage_local_ssd or local_nvme_ssd_block. Immutable.

### spec.nodeConfig.ephemeralStorageLocalSsd

`GcpGkeNodePoolEphemeralStorageLocalSsd`

Back ephemeral storage (emptyDir, container layers, logs) with local
NVMe SSDs instead of the boot disk — the biggest node-side I/O win for
churn-heavy workloads. Immutable.

### spec.nodeConfig.ephemeralStorageLocalSsd.localSsdCount

`uint32`

Number of local SSDs backing ephemeral storage. Each is 375 GB (or
3000 GB on Z3); count must match what the machine type supports.

- rule: {"uint32":{"gte":1}}

### spec.nodeConfig.ephemeralStorageLocalSsd.dataCacheCount

`uint32` · optional (explicit presence)

Of those, SSDs dedicated to GKE Data Cache (read caching for
persistent volumes).

### spec.nodeConfig.localNvmeSsdBlock

`GcpGkeNodePoolLocalNvmeSsdBlock`

Attach raw-block local NVMe SSDs for workloads that manage their own
filesystem (databases with direct block access). Immutable.

### spec.nodeConfig.localNvmeSsdBlock.localSsdCount

`uint32`

Number of raw-block local NVMe SSDs (375 GB each).

- rule: {"uint32":{"gte":1}}

### spec.nodeConfig.gcfsEnabled

`bool`

Image streaming (GCFS): containers start before the full image is
pulled, pulling data on demand — large-image pools (ML frameworks)
start minutes faster. Requires Container-Optimized OS.

### spec.nodeConfig.gvnicEnabled

`bool`

gVNIC (Google Virtual NIC): higher-throughput networking than virtio;
required for TIER_1 bandwidth and 100+ Gbps machine shapes. Immutable.

### spec.nodeConfig.fastSocketEnabled

`bool`

NCCL Fast Socket: optimizes multi-node collective communication for
distributed GPU training. Requires gvnic_enabled.

### spec.nodeConfig.bootDiskKmsKey

`string | valueFrom`

Customer-managed encryption key (CMEK) for the node boot disks.
Accepts a full crypto key path or a reference to a GcpKmsKey resource.
The Compute Engine service agent needs Encrypter/Decrypter on the key.
Immutable.

- references: GcpKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.nodeConfig.workloadMetadataMode

`string`

How workloads see instance metadata: GKE_METADATA runs the GKE
metadata server per node (required for Workload Identity — the default
on WI clusters and the right answer); GCE_METADATA exposes the raw VM
metadata (legacy, leaks node credentials to pods).

- rule: workload_metadata_mode must be empty, GCE_METADATA, or GKE_METADATA

### spec.nodeConfig.reservationAffinity

`GcpGkeNodePoolReservationAffinity`

Consume Compute Engine reservations: committed capacity for the pool.
ANY_RESERVATION uses any matching reservation; SPECIFIC_RESERVATION
targets one by name (key "compute.googleapis.com/reservation-name",
values [name]); NO_RESERVATION opts out. Immutable.

- rule: SPECIFIC_RESERVATION requires key "compute.googleapis.com/reservation-name" and the reservation name in values

### spec.nodeConfig.reservationAffinity.consumeReservationType

`string` · required

NO_RESERVATION opts out; ANY_RESERVATION consumes any matching
reservation; SPECIFIC_RESERVATION targets one by name.

- rule: {"required":true,"string":{"in":["NO_RESERVATION","ANY_RESERVATION","SPECIFIC_RESERVATION"]}}

### spec.nodeConfig.reservationAffinity.key

`string`

Reservation label key — "compute.googleapis.com/reservation-name" for
SPECIFIC_RESERVATION.

### spec.nodeConfig.reservationAffinity.values

`[]string`

Reservation label values — the reservation name(s).

### spec.nodeConfig.secondaryBootDisks

`[]GcpGkeNodePoolSecondaryBootDisk`

Secondary boot disks that preload container images or data onto every
node — cold-start acceleration for very large images. Immutable.

### spec.nodeConfig.secondaryBootDisks[].diskImage

`string` · required

Disk image to create the secondary boot disk from (a prepared image
containing the container images/data to preload).

- rule: {"required":true}

### spec.nodeConfig.secondaryBootDisks[].mode

`string`

CONTAINER_IMAGE_CACHE serves preloaded container images to the
container runtime. Empty attaches the disk without special handling.

- rule: mode must be empty or CONTAINER_IMAGE_CACHE

### spec.nodeConfig.kubeletConfig

`GcpGkeNodePoolKubeletConfig`

Kubelet tuning: CPU management, PID limits, log rotation, image GC.
Only set what you need; unset fields keep GKE defaults.

### spec.nodeConfig.kubeletConfig.cpuManagerPolicy

`string`

CPU management: "static" gives Guaranteed-QoS pods exclusive cores
(latency-sensitive workloads); "none" (default) shares cores.

- rule: cpu_manager_policy must be empty, static, or none

### spec.nodeConfig.kubeletConfig.cpuCfsQuota

`bool` · optional (explicit presence)

Enforce CPU CFS quota for containers with CPU limits. Disabling trades
throttling for potential node CPU contention.

### spec.nodeConfig.kubeletConfig.cpuCfsQuotaPeriod

`string`

CFS quota period, e.g. "100ms" (kubelet default) — shorter periods
smooth throttling for latency-sensitive workloads.

### spec.nodeConfig.kubeletConfig.podPidsLimit

`int64` · optional (explicit presence)

Maximum processes per pod — a fork-bomb fence for multi-tenant pools.

- rule: {"int64":{"lte":"4194304","gte":"1024"}}

### spec.nodeConfig.kubeletConfig.insecureKubeletReadonlyPortEnabled

`string`

The kubelet's insecure read-only port 10255: FALSE closes it (the
hardened posture new clusters default to); TRUE keeps it open for
legacy monitoring agents that still scrape it.

- rule: insecure_kubelet_readonly_port_enabled must be empty, TRUE, or FALSE

### spec.nodeConfig.kubeletConfig.maxParallelImagePulls

`int64` · optional (explicit presence)

Parallel image pulls (kubelet default serializes fewer) — speeds up
busy nodes that churn many images.

### spec.nodeConfig.kubeletConfig.containerLogMaxSize

`string`

Container log file size before rotation, e.g. "10Mi" (10Mi-500Mi).

### spec.nodeConfig.kubeletConfig.containerLogMaxFiles

`int64` · optional (explicit presence)

Rotated container log files kept per container (2-10).

- rule: {"int64":{"lte":"10","gte":"2"}}

### spec.nodeConfig.kubeletConfig.imageGcLowThresholdPercent

`int64` · optional (explicit presence)

Disk usage percent below which image GC never runs (must be lower
than the high threshold).

- rule: {"int64":{"lte":"85","gte":"10"}}

### spec.nodeConfig.kubeletConfig.imageGcHighThresholdPercent

`int64` · optional (explicit presence)

Disk usage percent above which image GC always runs.

- rule: {"int64":{"lte":"85","gte":"10"}}

### spec.nodeConfig.linuxNodeConfig

`GcpGkeNodePoolLinuxNodeConfig`

Linux node OS tuning: sysctls, cgroup mode, hugepages.

### spec.nodeConfig.linuxNodeConfig.sysctls

`map<string, string>`

Sysctls applied to every node, e.g.
{"net.core.somaxconn": "4096"} — only GKE's allowlisted keys.

### spec.nodeConfig.linuxNodeConfig.cgroupMode

`string`

Container runtime cgroup mode: CGROUP_MODE_V2 (the modern default on
current GKE versions) or CGROUP_MODE_V1 for workloads that still need
v1.

- rule: cgroup_mode must be empty, CGROUP_MODE_UNSPECIFIED, CGROUP_MODE_V1, or CGROUP_MODE_V2

### spec.nodeConfig.linuxNodeConfig.hugepagesConfig

`GcpGkeNodePoolHugepagesConfig`

Hugepage pre-allocation for DPDK/database workloads.

### spec.nodeConfig.linuxNodeConfig.hugepagesConfig.hugepageSize2m

`int64` · optional (explicit presence)

Number of 2MB hugepages.

### spec.nodeConfig.linuxNodeConfig.hugepagesConfig.hugepageSize1g

`int64` · optional (explicit presence)

Number of 1GB hugepages.

### spec.nodeConfig.loggingVariant

`string`

Node system-log throughput: DEFAULT (100 KiB/s) or MAX_THROUGHPUT
(10 MiB/s, costs a little node CPU) for pools with log-heavy workloads.

- rule: logging_variant must be empty, DEFAULT, or MAX_THROUGHPUT

### spec.nodeConfig.flexStart

`bool`

Flex-start (Dynamic Workload Scheduler): nodes are requested when
needed and run up to 7 days at a discount — the on-demand counterpart
to queued provisioning for hard-to-get GPU capacity. Immutable.

### spec.nodeConfig.maxRunDuration

`string`

Maximum runtime of each node, in seconds format (e.g. "3600s"), after
which it is drained and deleted. Pairs with flex-start/Spot batch
pools; leave empty for long-lived pools. Immutable.

- rule: max_run_duration must be a seconds-format duration like "3600s"

## Validation Rules

- `initial_node_count_requires_autoscaling`: initial_node_count only applies to autoscaled pools — a fixed-size pool's size is node_count itself

## Outputs

Reference an output from another manifest as `valueFrom: {kind: GcpGkeNodePool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.node_pool_name` | `string` | Name of the node pool as created in GKE — the handle gcloud commands and in-cluster references use. Matches spec.node_pool_name when set, otherwise metadata.name. |
| `status.outputs.instance_group_urls` | `[]string` | Resource URLs of the managed instance groups backing this pool — one per zone for regional clusters. What instance-group-targeted load-balancer backends and infrastructure automation compose against. |
| `status.outputs.min_nodes` | `string` | Effective minimum size of the pool: the autoscaling minimum (per zone) when autoscaling manages the pool, else the fixed node_count. |
| `status.outputs.max_nodes` | `string` | Effective maximum size of the pool: the autoscaling maximum (per zone) when autoscaling manages the pool, else the fixed node_count. |
| `status.outputs.current_node_count` | `string` | Number of nodes per zone at the time of the last deploy. For autoscaled pools this drifts as the autoscaler works — treat it as a snapshot, not live state. |
| `status.outputs.node_pool_id` | `string` | Fully qualified GKE node pool resource ID: projects/{project}/locations/{location}/clusters/{cluster}/nodePools/{name}. The handle downstream services (e.g. Dataproc on GKE) reference the pool by, and the path API callers address it with. |
| `status.outputs.location` | `string` | The pool's location (the parent cluster's region or zone), exactly as provided in the spec. |
| `status.outputs.version` | `string` | The Kubernetes version running on the pool's nodes. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.projectId` | GcpProject | `status.outputs.project_id` |
| `spec.clusterName` | GcpGkeCluster | `status.outputs.name` |
| `spec.location` | GcpGkeCluster | `status.outputs.location` |
| `spec.nodeConfig.serviceAccount` | GcpServiceAccount | `status.outputs.email` |
| `spec.nodeConfig.bootDiskKmsKey` | GcpKmsKey | `status.outputs.key_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| GcpDataprocCluster | `spec.virtualClusterConfig.kubernetesClusterConfig.gkeClusterConfig.nodePoolTarget[].nodePool` | `status.outputs.node_pool_id` |

## See Also

- [Overview](../README.md)
