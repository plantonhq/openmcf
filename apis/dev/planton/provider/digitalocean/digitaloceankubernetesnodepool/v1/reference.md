# DigitalOceanKubernetesNodePool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1`

DigitalOceanKubernetesNodePoolSpec defines the specification for creating a node pool in an existing DigitalOcean Kubernetes cluster (DOKS).
It focuses on essential parameters, following the 80/20 principle to expose only the most commonly used settings.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.nodePoolName` | `string` | yes |  |  |
| `spec.cluster` | `string \| valueFrom` | yes |  | DigitalOceanKubernetesCluster (`metadata.name`) |
| `spec.size` | `string` | yes |  |  |
| `spec.nodeCount` | `uint32` | yes |  |  |
| `spec.autoScale` | `bool` |  |  |  |
| `spec.minNodes` | `uint32` |  |  |  |
| `spec.maxNodes` | `uint32` |  |  |  |
| `spec.labels` | `map<string, string>` |  |  |  |
| `spec.taints` | `[]DigitalOceanKubernetesNodePoolTaint` |  |  |  |
| `spec.taints[].key` | `string` | yes |  |  |
| `spec.taints[].value` | `string` |  |  |  |
| `spec.taints[].effect` | `string` | yes |  |  |
| `spec.tags` | `[]string` |  |  |  |

## Field Details

### spec.nodePoolName

`string` · required

A name for the node pool. Must be unique within the Kubernetes cluster.

- rule: {"required":true}

### spec.cluster

`string | valueFrom` · required

Reference to the DigitalOcean Kubernetes Cluster in which to create this node pool.
Accepts the cluster's name or a reference to the DigitalOceanKubernetesCluster resource.

- references: DigitalOceanKubernetesCluster (`metadata.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanKubernetesCluster, name: <that resource's name>, fieldPath: metadata.name}} -- a bare string does not parse

### spec.size

`string` · required

The slug identifier for the Droplet size to use for each node (e.g., "s-4vcpu-8gb").
This defines the CPU and memory of the nodes in the pool.

- rule: {"required":true}

### spec.nodeCount

`uint32` · required

The number of nodes to provision in the pool.
Must be at least 1. If auto_scale is enabled, this acts as the initial desired node count.

- rule: {"required":true,"uint32":{"gt":0}}

### spec.autoScale

`bool`

Enable auto-scaling for this node pool.
If true, the platform will manage node count between min_nodes and max_nodes.

### spec.minNodes

`uint32`

Minimum number of nodes when auto-scaling is enabled.
Required if auto_scale is true.

### spec.maxNodes

`uint32`

Maximum number of nodes when auto-scaling is enabled.
Required if auto_scale is true.

### spec.labels

`map<string, string>`

Kubernetes labels to apply to all nodes in this pool.
Labels are key-value pairs used for node selection and workload scheduling.
Example: {"workload": "web", "env": "production"}

### spec.taints

`[]DigitalOceanKubernetesNodePoolTaint`

Kubernetes taints to apply to all nodes in this pool.
Taints prevent pods from being scheduled on these nodes unless they have matching tolerations.
Commonly used for workload isolation (e.g., GPU nodes, dedicated system pools).

### spec.taints[].key

`string` · required

The taint key (e.g., "nvidia.com/gpu", "workload", "dedicated")

- rule: {"required":true}

### spec.taints[].value

`string`

The taint value (e.g., "true", "gpu", "system")

### spec.taints[].effect

`string` · required

The taint effect: NoSchedule, PreferNoSchedule, or NoExecute
- NoSchedule: Pods that don't tolerate this taint will not be scheduled on the node
- PreferNoSchedule: Kubernetes will try to avoid scheduling pods that don't tolerate this taint
- NoExecute: Pods that don't tolerate this taint will be evicted if already running

- rule: {"required":true}

### spec.tags

`[]string`

A list of DigitalOcean tags to apply to the node pool Droplets.
Tags are used for cost attribution and organizational purposes in DigitalOcean's billing and management.
Note: This is different from Kubernetes labels. Tags affect DO billing, labels affect K8s scheduling.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanKubernetesNodePool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.node_pool_id` | `string` | The unique identifier (UUID) of the created node pool. |
| `status.outputs.node_ids` | `[]string` | The IDs of the individual Droplet nodes in the pool. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.cluster` | DigitalOceanKubernetesCluster | `metadata.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
