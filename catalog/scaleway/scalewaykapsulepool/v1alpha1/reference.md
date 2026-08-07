# ScalewayKapsulePool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayKapsulePoolSpec defines the specification for an additional node
pool in an existing Scaleway Kapsule Kubernetes cluster.

Node pools provide compute capacity for Kubernetes workloads. Each pool
consists of identical nodes (same instance type, root volume, container
runtime) and can be independently scaled, upgraded, and configured.

This is a **standalone resource** (not composite) that creates a single
`scaleway_k8s_pool` resource. It depends on an existing
ScalewayKapsuleCluster, referenced via `cluster_id`.

**Kubernetes labels and taints** are first-class fields in this spec.
Under the hood, they are applied via Scaleway's Cloud Controller Manager
(CCM) tag convention:
  - Labels: Pool tag `noprefix={key}={value}` → K8s node label `{key}={value}`
  - Taints: Pool tag `taint=noprefix={key}={value}:{Effect}` → K8s taint

This abstraction gives users the same clean experience as other providers
(DigitalOcean, AWS) while using the correct Scaleway mechanism.

**Composition pattern**: The pool requires a cluster reference via
`StringValueOrRef`. Infra charts wire this using `valueFrom`:

  clusterId:
    valueFrom:
      kind: ScalewayKapsuleCluster
      name: my-cluster
      fieldPath: status.outputs.cluster_id

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.clusterId` | `string \| valueFrom` | yes |  | ScalewayKapsuleCluster (`status.outputs.cluster_id`) |
| `spec.nodeType` | `string` | yes |  |  |
| `spec.size` | `int32` | yes |  |  |
| `spec.autoScale` | `bool` |  |  |  |
| `spec.minSize` | `int32` |  |  |  |
| `spec.maxSize` | `int32` |  |  |  |
| `spec.autohealing` | `bool` |  |  |  |
| `spec.containerRuntime` | `string` |  | `containerd` |  |
| `spec.rootVolumeType` | `string` |  |  |  |
| `spec.rootVolumeSizeInGb` | `int32` |  |  |  |
| `spec.publicIpDisabled` | `bool` |  |  |  |
| `spec.zone` | `string` |  |  |  |
| `spec.placementGroupId` | `string` |  |  |  |
| `spec.kubernetesLabels` | `map<string, string>` |  |  |  |
| `spec.taints` | `[]ScalewayKapsulePoolTaint` |  |  |  |
| `spec.taints[].key` | `string` | yes |  |  |
| `spec.taints[].value` | `string` |  |  |  |
| `spec.taints[].effect` | `string` | yes |  |  |
| `spec.upgradePolicy` | `ScalewayKapsulePoolUpgradePolicy` |  |  |  |
| `spec.upgradePolicy.maxSurge` | `int32` |  |  |  |
| `spec.upgradePolicy.maxUnavailable` | `int32` |  |  |  |
| `spec.kubeletArgs` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Scaleway region where the pool will be created.

Must match the parent cluster's region. All nodes in this pool will
be placed in zones within this region.

Examples: "fr-par", "nl-ams", "pl-waw"

IMPORTANT: Cannot be changed after creation.

- rule: {"required":true}

### spec.clusterId

`string | valueFrom` · required

Reference to the Kapsule cluster in which to create this node pool.

Can be a literal cluster ID or a reference to a ScalewayKapsuleCluster
resource's output. In infra charts, this is wired via `valueFrom`.

IMPORTANT: Cannot be changed after creation.

- references: ScalewayKapsuleCluster (`status.outputs.cluster_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayKapsuleCluster, name: <that resource's name>, fieldPath: status.outputs.cluster_id}} -- a bare string does not parse

### spec.nodeType

`string` · required

Instance type for worker nodes (required).

Determines CPU, RAM, and local storage for each node. Common types:
  - Development:  "DEV1-M" (3 vCPU, 4 GB RAM)
  - General:      "GP1-XS" (4 vCPU, 16 GB RAM), "GP1-S" (8 vCPU, 32 GB)
  - Production:   "PRO2-S" (2 vCPU, 8 GB), "PRO2-M" (4 vCPU, 16 GB)

See Scaleway pricing page for the full catalog. Instances with
insufficient memory (DEV1-S, PLAY2-PICO, STARDUST) are not eligible.

IMPORTANT: Cannot be changed after creation. To change instance types,
create a new pool and migrate workloads.

- rule: {"required":true}

### spec.size

`int32` · required

Number of nodes in the pool (required).

When autoscaling is disabled, this is the fixed pool size. When
autoscaling is enabled, this is the initial size -- the autoscaler
will adjust between min_size and max_size based on workload demands.

Minimum: 1 (a pool must have at least one node).

Note: When autoscaling is enabled, updates to this field are ignored
by the provider -- the autoscaler controls the actual size.

- rule: {"required":true,"int32":{"gte":1}}

### spec.autoScale

`bool`

Enable the cluster autoscaler for this pool.

When true, Kubernetes automatically adds or removes nodes based on
pending pod resource requests. Requires min_size and max_size to
be configured. The autoscaler's behavior (delays, thresholds) is
controlled by the cluster-level `autoscaler_config` on the parent
ScalewayKapsuleCluster.

Default: false (fixed-size pool).

### spec.minSize

`int32`

Minimum number of nodes when autoscaling is enabled.

The autoscaler will not scale below this number, even if all nodes
are underutilized. Set to at least 1 for availability.

Only meaningful when auto_scale is true.

### spec.maxSize

`int32`

Maximum number of nodes when autoscaling is enabled.

The autoscaler will not scale above this number, even if pods are
pending. Controls cost ceiling.

Only meaningful when auto_scale is true.

### spec.autohealing

`bool`

Enable autohealing for this pool.

When true, Scaleway automatically detects and replaces unhealthy
nodes. A node is considered unhealthy if its kubelet stops
reporting status for a configurable period.

Recommended for production pools.

### spec.containerRuntime

`string`

Container runtime for pool nodes.

Options:
  - "containerd" (default, recommended) -- Industry-standard container
    runtime. Required for Kubernetes 1.24+.

IMPORTANT: Cannot be changed after creation.

- default: `containerd`

### spec.rootVolumeType

`string`

Root volume type for pool nodes.

Controls the storage backing each node's root filesystem. Options
depend on the instance type and availability zone.

IMPORTANT: Cannot be changed after creation.

### spec.rootVolumeSizeInGb

`int32`

Root volume size in GB for pool nodes.

If omitted, uses the default size for the instance type. Increase
for workloads that pull many large container images or need
significant local ephemeral storage.

IMPORTANT: Cannot be changed after creation.

### spec.publicIpDisabled

`bool`

Disable public IP addresses on pool nodes.

When true, nodes have only private IPs (from the cluster's Private
Network). This is the recommended security posture for production:
nodes are not reachable from the internet.

Requires a Public Gateway or NAT on the Private Network so nodes
can reach external registries and APIs.

IMPORTANT: Cannot be changed after creation.

### spec.zone

`string`

Zone within the region to place pool nodes.

Optional. If omitted, Scaleway chooses the zone automatically.
Use this for zone-specific placement (e.g., "fr-par-1", "fr-par-2").

All nodes in this pool will be in the specified zone. For multi-AZ
deployments, create separate pools in different zones.

IMPORTANT: Cannot be changed after creation.

### spec.placementGroupId

`string`

Placement group ID for anti-affinity scheduling.

Associates pool nodes with a Scaleway Instance placement group. Use
placement groups to spread nodes across different hypervisors for
higher availability.

Must be a valid placement group UUID. Placement groups are created
separately via the Scaleway console or API.

IMPORTANT: Cannot be changed after creation.

### spec.kubernetesLabels

`map<string, string>`

Kubernetes labels to apply to all nodes in this pool.

Labels are key-value pairs used for node selection and workload
scheduling. Pods use `nodeSelector` or node affinity rules to
target nodes with specific labels.

Example: {"workload": "gpu", "team": "ml", "tier": "compute"}

Implementation: Labels are applied via Scaleway's Cloud Controller
Manager tag convention. Each label generates a pool tag in the
format `noprefix={key}={value}`, which the CCM syncs to K8s nodes
as the label `{key}={value}`.

### spec.taints

`[]ScalewayKapsulePoolTaint`

Kubernetes taints to apply to all nodes in this pool.

Taints prevent pods from being scheduled on these nodes unless they
have matching tolerations. Commonly used for workload isolation
(e.g., GPU nodes, dedicated system pools, batch processing).

Implementation: Taints are applied via Scaleway's Cloud Controller
Manager tag convention. Each taint generates a pool tag in the
format `taint=noprefix={key}={value}:{Effect}`, which the CCM
syncs to K8s nodes as the taint `{key}={value}:{Effect}`.

### spec.taints[].key

`string` · required

The taint key.

Examples: "nvidia.com/gpu", "dedicated", "workload", "node-role"

- rule: {"required":true}

### spec.taints[].value

`string`

The taint value.

Examples: "true", "gpu", "batch", "system"

### spec.taints[].effect

`string` · required

The taint effect.

Must be one of:
  - "NoSchedule" -- Pods that don't tolerate this taint will not be
    scheduled on the node. Existing pods are not evicted.
  - "PreferNoSchedule" -- Kubernetes will try to avoid scheduling
    pods that don't tolerate this taint, but it's not guaranteed.
  - "NoExecute" -- Pods that don't tolerate this taint will be
    evicted if already running, and not scheduled if pending.

- rule: {"required":true}

### spec.upgradePolicy

`ScalewayKapsulePoolUpgradePolicy`

Node pool upgrade policy for rolling updates.

Controls how nodes are replaced during Kubernetes version upgrades
or pool configuration changes.

Optional. If omitted, Scaleway uses defaults (max_surge=0,
max_unavailable=1 -- one node at a time).

### spec.upgradePolicy.maxSurge

`int32`

Maximum number of extra nodes created during an upgrade.

Surge nodes are temporary workers that accept workloads while
existing nodes are drained and replaced. Higher values speed up
upgrades but temporarily increase cost.

Default: 0 (no surge nodes).

### spec.upgradePolicy.maxUnavailable

`int32`

Maximum number of nodes that can be unavailable simultaneously
during an upgrade.

Controls the disruption budget. Setting this to 1 means nodes
are replaced one at a time (safest but slowest).

Default: 1.

### spec.kubeletArgs

`map<string, string>`

Custom kubelet arguments for pool nodes.

Power-user escape hatch for setting kubelet flags not exposed as
first-class fields. Example: {"maxPods": "150"}.

WARNING: Use with caution. Incorrect kubelet arguments can prevent
nodes from joining the cluster or cause instability.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayKapsulePool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.pool_id` | `string` | The unique identifier of the created node pool. Format: regional ID (e.g., "fr-par/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"). Useful for management, monitoring, and API calls. |
| `status.outputs.pool_version` | `string` | The actual Kubernetes version running on pool nodes. May differ from the cluster's version during rolling upgrades or when the pool was created at a different time than the cluster. |
| `status.outputs.current_size` | `int32` | The actual number of nodes currently in the pool. When autoscaling is enabled, this may differ from the `size` field in the spec. Reflects the autoscaler's current decision. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.clusterId` | ScalewayKapsuleCluster | `status.outputs.cluster_id` |

## See Also

- [Overview](../README.md)
