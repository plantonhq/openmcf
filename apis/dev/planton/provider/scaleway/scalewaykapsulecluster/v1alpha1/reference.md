# ScalewayKapsuleCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayKapsuleClusterSpec defines the specification for a Scaleway Kapsule
managed Kubernetes cluster.

Kapsule is Scaleway's managed Kubernetes service. It provisions a fully
managed control plane (API server, etcd, scheduler, controller-manager)
and delegates worker nodes to separate node pools.

This is a **composite resource** that bundles:
  1. The Kapsule cluster (managed control plane).
  2. A default node pool (so the cluster is immediately usable).

Additional node pools can be added via separate `ScalewayKapsulePool`
resources that reference this cluster's `status.outputs.cluster_id`.

Kapsule clusters are **regional** resources (e.g., "fr-par"). Node pools
within the cluster can be placed in specific zones within the region.

**Composition pattern**: The cluster requires a Private Network reference
via `StringValueOrRef`. Downstream resources like `ScalewayKapsulePool`
reference `status.outputs.cluster_id`. Infra charts use
`status.outputs.kubeconfig` and `status.outputs.cluster_ca_certificate`
to configure the Kubernetes provider for addon deployments.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.kubernetesVersion` | `string` | yes |  |  |
| `spec.cni` | `string` | yes | `cilium` |  |
| `spec.privateNetworkId` | `string \| valueFrom` | yes |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.type` | `string` |  | `kapsule` |  |
| `spec.description` | `string` |  |  |  |
| `spec.deleteAdditionalResources` | `bool` |  | `true` |  |
| `spec.autoUpgrade` | `ScalewayKapsuleAutoUpgrade` |  |  |  |
| `spec.autoUpgrade.enable` | `bool` | yes |  |  |
| `spec.autoUpgrade.maintenanceWindowStartHour` | `int32` | yes |  |  |
| `spec.autoUpgrade.maintenanceWindowDay` | `string` | yes |  |  |
| `spec.autoscalerConfig` | `ScalewayKapsuleAutoscalerConfig` |  |  |  |
| `spec.autoscalerConfig.disableScaleDown` | `bool` |  |  |  |
| `spec.autoscalerConfig.scaleDownDelayAfterAdd` | `string` |  |  |  |
| `spec.autoscalerConfig.scaleDownUnneededTime` | `string` |  |  |  |
| `spec.autoscalerConfig.estimator` | `string` |  |  |  |
| `spec.autoscalerConfig.expander` | `string` |  |  |  |
| `spec.autoscalerConfig.scaleDownUtilizationThreshold` | `double` |  |  |  |
| `spec.autoscalerConfig.maxGracefulTerminationSec` | `int32` |  |  |  |
| `spec.autoscalerConfig.ignoreDaemonsetsUtilization` | `bool` |  |  |  |
| `spec.autoscalerConfig.balanceSimilarNodeGroups` | `bool` |  |  |  |
| `spec.autoscalerConfig.expendablePodsPriorityCutoff` | `int32` |  |  |  |
| `spec.featureGates` | `[]string` |  |  |  |
| `spec.admissionPlugins` | `[]string` |  |  |  |
| `spec.podCidr` | `string` |  |  |  |
| `spec.serviceCidr` | `string` |  |  |  |
| `spec.defaultNodePool` | `ScalewayKapsuleDefaultNodePool` | yes |  |  |
| `spec.defaultNodePool.name` | `string` |  |  |  |
| `spec.defaultNodePool.nodeType` | `string` | yes |  |  |
| `spec.defaultNodePool.size` | `int32` | yes |  |  |
| `spec.defaultNodePool.autoScale` | `bool` |  |  |  |
| `spec.defaultNodePool.minSize` | `int32` |  |  |  |
| `spec.defaultNodePool.maxSize` | `int32` |  |  |  |
| `spec.defaultNodePool.autohealing` | `bool` |  |  |  |
| `spec.defaultNodePool.containerRuntime` | `string` |  | `containerd` |  |
| `spec.defaultNodePool.rootVolumeType` | `string` |  |  |  |
| `spec.defaultNodePool.rootVolumeSizeInGb` | `int32` |  |  |  |
| `spec.defaultNodePool.publicIpDisabled` | `bool` |  |  |  |
| `spec.defaultNodePool.upgradePolicy` | `ScalewayKapsuleNodePoolUpgradePolicy` |  |  |  |
| `spec.defaultNodePool.upgradePolicy.maxSurge` | `int32` |  |  |  |
| `spec.defaultNodePool.upgradePolicy.maxUnavailable` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

The Scaleway region where the cluster will be created.
Examples: "fr-par", "nl-ams", "pl-waw"

The region determines which data centers are available for node pools.
All node pools in this cluster must be in zones within this region.

This field is required and cannot be changed after creation.

- rule: {"required":true}

### spec.kubernetesVersion

`string` · required

Kubernetes version for the cluster.

Can be a minor version (e.g., "1.32") or a full patch version (e.g.,
"1.32.3"). When auto-upgrade is enabled, use a minor version so
Scaleway can automatically apply patch upgrades.

Available versions depend on the region and can be checked via the
Scaleway API or console.

- rule: {"required":true}

### spec.cni

`string` · required

Container Network Interface (CNI) plugin for pod networking.

Options:
  - "cilium" (recommended) -- eBPF-based. High performance, advanced
    network policies, service mesh integration, Hubble observability.
  - "calico" -- Mature, widely adopted. Standard Kubernetes network
    policies. Good for teams already familiar with Calico.

IMPORTANT: This field cannot be changed after creation. Changing the
CNI requires recreating the entire cluster.

- default: `cilium`
- rule: {"required":true}

### spec.privateNetworkId

`string | valueFrom` · required

The Private Network to attach the cluster to.

All Kapsule clusters require a Private Network. Nodes communicate
with the control plane and with each other over this network.

Can be a literal Private Network UUID or a reference to a
ScalewayPrivateNetwork resource's output.

In infra charts, this is typically wired via valueFrom:

  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: my-network
      fieldPath: status.outputs.private_network_id

IMPORTANT: This field cannot be changed after creation.

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.type

`string`

Kapsule cluster type.

Options:
  - "kapsule" (default) -- Mutualized (shared) control plane. Suitable
    for most workloads. No additional cost for the control plane.
  - "kapsule-dedicated-4"  -- Dedicated control plane, 4 nodes max.
  - "kapsule-dedicated-8"  -- Dedicated control plane, 8 nodes max.
  - "kapsule-dedicated-16" -- Dedicated control plane, 16 nodes max.

Dedicated control planes provide isolated API servers with guaranteed
resources. Use for production workloads requiring API server SLAs
or strict multi-tenancy isolation.

This field can be changed after creation (upgrade from mutualized
to dedicated, or change dedicated tier).

- default: `kapsule`

### spec.description

`string`

Human-readable description for the cluster.

Optional. Shown in the Scaleway console for identification.

### spec.deleteAdditionalResources

`bool`

Whether to delete additional resources (LBs, volumes, routes) created
by Kubernetes when the cluster is destroyed.

When true, Scaleway automatically cleans up load balancers (from
Services of type LoadBalancer), persistent volumes (from PVCs), and
other resources that Kubernetes provisioned during the cluster's
lifetime. When false, these resources are orphaned and must be
cleaned up manually.

Default: true. Set to false only when you need to preserve
Kubernetes-managed resources after cluster deletion (e.g., data
volumes for migration).

- default: `true`

### spec.autoUpgrade

`ScalewayKapsuleAutoUpgrade`

Automatic patch version upgrade configuration.

When enabled, Scaleway automatically upgrades the cluster to the
latest patch version within the current minor version during the
specified maintenance window. For example, if the cluster is on
1.32.1, it may be upgraded to 1.32.3 during the window.

Optional. If omitted, auto-upgrade is disabled and the cluster
stays on the exact version specified in `kubernetes_version`.

### spec.autoUpgrade.enable

`bool` · required

Whether auto-upgrade is enabled.

- rule: {"required":true}

### spec.autoUpgrade.maintenanceWindowStartHour

`int32` · required

UTC hour (0-23) when the maintenance window starts.

The upgrade process begins at this hour. Choose a low-traffic
period for your workloads. Example: 2 (2:00 AM UTC).

- rule: {"required":true,"int32":{"lte":23,"gte":0}}

### spec.autoUpgrade.maintenanceWindowDay

`string` · required

Day of the week for the maintenance window.

Options: "monday", "tuesday", "wednesday", "thursday", "friday",
"saturday", "sunday", or "any" (Scaleway picks the best day).

Example: "sunday" for weekend maintenance.

- rule: {"required":true}

### spec.autoscalerConfig

`ScalewayKapsuleAutoscalerConfig`

Cluster-wide autoscaler configuration.

Controls HOW the Kubernetes cluster autoscaler behaves when scaling
node pools. Autoscaling itself is toggled per-pool (on the default
node pool below, or on separate ScalewayKapsulePool resources), but
the behavior parameters (delays, thresholds, algorithms) are
configured here at the cluster level.

Optional. If omitted, Scaleway uses sensible defaults (binpacking
estimator, random expander, 10m scale-down delay, 0.5 utilization
threshold).

### spec.autoscalerConfig.disableScaleDown

`bool`

Disable the scale-down behavior entirely.

When true, the autoscaler only scales UP (adds nodes) but never
removes underutilized nodes. Useful during rollouts or when
stability is more important than cost optimization.

Default: false.

### spec.autoscalerConfig.scaleDownDelayAfterAdd

`string`

How long to wait after a scale-up before considering scale-down.

Duration string (e.g., "10m", "15m"). Prevents thrashing by
allowing newly added nodes time to receive workloads.

Default: "10m".

### spec.autoscalerConfig.scaleDownUnneededTime

`string`

How long a node must be underutilized before it's eligible for
scale-down.

Duration string (e.g., "10m", "20m"). Higher values tolerate
temporary dips in utilization.

Default: "10m".

### spec.autoscalerConfig.estimator

`string`

Resource estimation algorithm for scheduling decisions.

Options:
  - "binpacking" (default) -- Estimates how many nodes are needed
    by bin-packing pending pods. Most accurate for heterogeneous
    workloads.

### spec.autoscalerConfig.expander

`string`

Node group expansion strategy when multiple groups can accommodate
pending pods.

Options:
  - "random" (default) -- Pick a random eligible group.
  - "most-pods"  -- Pick the group that schedules the most pods.
  - "least-waste" -- Pick the group with least resource waste.
  - "priority" -- Pick based on user-defined priorities.

### spec.autoscalerConfig.scaleDownUtilizationThreshold

`double`

CPU/memory utilization threshold below which a node is considered
underutilized and eligible for scale-down.

Range: 0.0 to 1.0. For example, 0.5 means a node using less than
50% of its allocatable resources is a scale-down candidate.

Default: 0.5. Lower values are more aggressive (remove more nodes);
higher values are more conservative (keep more headroom).

### spec.autoscalerConfig.maxGracefulTerminationSec

`int32`

Maximum time (in seconds) the autoscaler waits for pod termination
during scale-down.

Pods with long graceful termination periods may need a higher value.

Default: 600 (10 minutes).

### spec.autoscalerConfig.ignoreDaemonsetsUtilization

`bool`

Whether to consider DaemonSet resource usage when calculating node
utilization.

When false (default), DaemonSet resource requests are excluded from
utilization calculations, making scale-down more conservative.

Default: false.

### spec.autoscalerConfig.balanceSimilarNodeGroups

`bool`

Whether to balance the number of nodes across similar node groups.

When true, the autoscaler tries to keep similarly-sized node groups
at the same size. Useful for multi-AZ setups.

Default: false.

### spec.autoscalerConfig.expendablePodsPriorityCutoff

`int32`

Priority cutoff for expendable pods during scale-down.

Pods with priority below this value are considered expendable and
won't prevent scale-down. Default: -10.

### spec.featureGates

`[]string`

Kubernetes feature gates to enable on the cluster.

Feature gates are alpha/beta Kubernetes features that can be toggled.
Example: ["GracefulNodeShutdown", "HPAContainerMetrics"]

Consult the Kubernetes documentation for available feature gates
at the cluster's version. Enabling unstable feature gates may affect
cluster stability.

Optional. Most clusters don't need custom feature gates.

### spec.admissionPlugins

`[]string`

Kubernetes admission plugins to enable on the cluster.

Admission plugins intercept API server requests after authentication
and authorization but before persistence. Scaleway enables a standard
set by default; this field adds additional plugins.

Example: ["AlwaysPullImages", "NodeRestriction"]

Optional. Most clusters don't need additional admission plugins.

### spec.podCidr

`string`

Pod CIDR for the cluster's pod network.

The IP range allocated to pods. Each node gets a /24 subnet from
this range. Must be large enough to accommodate all pods across
all nodes.

Default: "100.64.0.0/15" (131,072 addresses). Only change this
if you have specific IP planning requirements or conflicts.

IMPORTANT: Cannot be changed after creation.

### spec.serviceCidr

`string`

Service CIDR for Kubernetes services.

The IP range allocated to ClusterIP services. Must not overlap
with pod_cidr or the Private Network's subnet.

Default: "10.32.0.0/20" (4,096 addresses). Only change this if
you have specific IP planning requirements or conflicts.

IMPORTANT: Cannot be changed after creation.

### spec.defaultNodePool

`ScalewayKapsuleDefaultNodePool` · required

The default node pool configuration.

Every Kapsule cluster needs at least one node pool to run workloads.
This embedded pool is created alongside the cluster, giving you a
working cluster from a single resource.

For additional node pools with different instance types, labels, or
taints, create separate `ScalewayKapsulePool` resources that reference
this cluster's `status.outputs.cluster_id`.

- rule: {"required":true}

### spec.defaultNodePool.name

`string`

Pool name. If omitted, defaults to "{cluster-name}-default".

Must be unique within the cluster. Use a descriptive name like
"system", "default", or "general".

IMPORTANT: Cannot be changed after creation.

### spec.defaultNodePool.nodeType

`string` · required

Instance type for worker nodes (required).

Determines CPU, RAM, and local storage for each node. Common types:
  - Development:  "DEV1-M" (3 vCPU, 4 GB RAM)
  - General:      "GP1-XS" (4 vCPU, 16 GB RAM), "GP1-S" (8 vCPU, 32 GB)
  - Production:   "PRO2-S" (2 vCPU, 8 GB), "PRO2-M" (4 vCPU, 16 GB)

See Scaleway pricing page for the full catalog. Instance type
cannot be changed in-place -- changing it requires creating a new
pool and migrating workloads.

IMPORTANT: Cannot be changed after creation.

- rule: {"required":true}

### spec.defaultNodePool.size

`int32` · required

Initial number of nodes in the pool (required).

When autoscaling is disabled, this is the fixed pool size. When
autoscaling is enabled, this is the initial size -- the autoscaler
will adjust between min_size and max_size based on workload demands.

Minimum: 1 (a pool must have at least one node).

- rule: {"required":true,"int32":{"gte":1}}

### spec.defaultNodePool.autoScale

`bool`

Enable the cluster autoscaler for this pool.

When true, Kubernetes automatically adds or removes nodes based on
pending pod resource requests. Requires min_size and max_size to
be configured. The autoscaler's behavior (delays, thresholds) is
controlled by the cluster-level `autoscaler_config`.

Default: false.

### spec.defaultNodePool.minSize

`int32`

Minimum number of nodes when autoscaling is enabled.

The autoscaler will not scale below this number, even if all nodes
are underutilized. Set to at least 1 for availability.

Only meaningful when auto_scale is true.

### spec.defaultNodePool.maxSize

`int32`

Maximum number of nodes when autoscaling is enabled.

The autoscaler will not scale above this number, even if pods are
pending. Controls cost ceiling.

Only meaningful when auto_scale is true.

### spec.defaultNodePool.autohealing

`bool`

Enable autohealing for this pool.

When true, Scaleway automatically detects and replaces unhealthy
nodes. A node is considered unhealthy if its kubelet stops
reporting status for a configurable period.

Recommended for production clusters.

### spec.defaultNodePool.containerRuntime

`string`

Container runtime for pool nodes.

Options:
  - "containerd" (default, recommended) -- Industry-standard container
    runtime. Required for Kubernetes 1.24+.

IMPORTANT: Cannot be changed after creation.

- default: `containerd`

### spec.defaultNodePool.rootVolumeType

`string`

Root volume type for pool nodes.

Controls the storage backing each node's root filesystem. Options
depend on the instance type and availability zone.

IMPORTANT: Cannot be changed after creation.

### spec.defaultNodePool.rootVolumeSizeInGb

`int32`

Root volume size in GB for pool nodes.

If omitted, uses the default size for the instance type. Increase
for workloads that pull many large container images or need
significant local ephemeral storage.

IMPORTANT: Cannot be changed after creation.

### spec.defaultNodePool.publicIpDisabled

`bool`

Disable public IP addresses on pool nodes.

When true, nodes have only private IPs (from the cluster's Private
Network). This is the recommended security posture for production:
nodes are not reachable from the internet.

Requires a Public Gateway or NAT on the Private Network so nodes
can reach external registries and APIs.

Default: false (nodes get public IPs).

### spec.defaultNodePool.upgradePolicy

`ScalewayKapsuleNodePoolUpgradePolicy`

Node pool upgrade policy for rolling updates.

Controls how nodes are replaced during Kubernetes version upgrades
or pool configuration changes.

Optional. If omitted, Scaleway uses defaults (max_surge=0,
max_unavailable=1 -- one node at a time).

### spec.defaultNodePool.upgradePolicy.maxSurge

`int32`

Maximum number of extra nodes created during an upgrade.

Surge nodes are temporary workers that accept workloads while
existing nodes are drained and replaced. Higher values speed up
upgrades but temporarily increase cost.

Default: 0 (no surge nodes).

### spec.defaultNodePool.upgradePolicy.maxUnavailable

`int32`

Maximum number of nodes that can be unavailable simultaneously
during an upgrade.

Controls the disruption budget. Setting this to 1 means nodes
are replaced one at a time (safest but slowest).

Default: 1.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayKapsuleCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | The unique identifier of the created Kapsule cluster. Format: regional ID (e.g., "fr-par/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"). This is the primary output referenced by downstream resources: - ScalewayKapsulePool resources (via StringValueOrRef on cluster_id) - Monitoring and management tools In infra charts, downstream pools wire to this value using: clusterId: valueFrom: kind: ScalewayKapsuleCluster name: my-cluster fieldPath: status.outputs.cluster_id |
| `status.outputs.kubeconfig` | `string` | Raw kubeconfig file content for connecting to the cluster. Contains the API server URL, CA certificate, and authentication token needed to interact with the Kubernetes API. Can be written directly to `~/.kube/config` or passed to kubectl via `KUBECONFIG` environment variable. SENSITIVE: This output contains authentication credentials. Handle with care in CI/CD pipelines and logging. |
| `status.outputs.apiserver_url` | `string` | The URL of the Kubernetes API server. Format: "https://<uuid>.api.k8s.fr-par.scw.cloud:6443" Use this output when configuring Kubernetes providers in IaC tools (Pulumi, Terraform) or when setting up external monitoring and CI/CD systems that need direct API access. |
| `status.outputs.cluster_ca_certificate` | `string` | The CA certificate of the Kubernetes API server (base64-encoded). Used together with `apiserver_url` and a token to configure Kubernetes providers in IaC tools without the full kubeconfig. This is the recommended approach for infra charts that deploy Kubernetes addons (cert-manager, ingress-nginx, etc.) as part of the same deployment DAG. |
| `status.outputs.wildcard_dns` | `string` | DNS wildcard for ready nodes in the cluster. Format: "*.xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.nodes.k8s.fr-par.scw.cloud" Can be used for DNS-based service discovery or as a CNAME target for external-facing services. |
| `status.outputs.default_pool_id` | `string` | The unique identifier of the default node pool. Format: regional ID. Useful for management, monitoring, and distinguishing the default pool from additional pools. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| ScalewayKapsulePool | `spec.clusterId` | `status.outputs.cluster_id` |

## See Also

- [Overview](./README.md)
