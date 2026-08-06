# CivoKubernetesNodePool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1alpha1`

CivoKubernetesNodePoolSpec defines the specification for creating a node pool in an existing Civo Kubernetes cluster.
It focuses on essential parameters, following the 80/20 principle to expose only the most commonly used settings.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.nodePoolName` | `string` | yes |  |  |
| `spec.cluster` | `string \| valueFrom` | yes |  | CivoKubernetesCluster (`metadata.name`) |
| `spec.size` | `string` | yes |  |  |
| `spec.nodeCount` | `uint32` | yes |  |  |
| `spec.autoScale` | `bool` |  |  |  |
| `spec.minNodes` | `uint32` |  |  |  |
| `spec.maxNodes` | `uint32` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |

## Field Details

### spec.nodePoolName

`string` · required

A name for the node pool. Must be unique within the Civo Kubernetes cluster.

- rule: {"required":true}

### spec.cluster

`string | valueFrom` · required

Reference to the Civo Kubernetes cluster in which to create this node pool.
Accepts the cluster's name or a reference to the CivoKubernetesCluster resource.

- references: CivoKubernetesCluster (`metadata.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoKubernetesCluster, name: <that resource's name>, fieldPath: metadata.name}} -- a bare string does not parse

### spec.size

`string` · required

The slug identifier for the size of each node in the pool (e.g., "g4s.kube.medium").
Defines the CPU and memory of the nodes.

- rule: {"required":true}

### spec.nodeCount

`uint32` · required

The number of nodes to provision in the pool. Must be at least 1.

- rule: {"required":true,"uint32":{"gt":0}}

### spec.autoScale

`bool`

Enable auto-scaling for this node pool. If true, the platform manages node count between min_nodes and max_nodes.

### spec.minNodes

`uint32`

Minimum number of nodes when auto-scaling is enabled. Required if auto_scale is true.

### spec.maxNodes

`uint32`

Maximum number of nodes when auto-scaling is enabled. Required if auto_scale is true.

### spec.tags

`[]string`

A list of tags to apply to the node pool for organizational purposes.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoKubernetesNodePool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.node_pool_id` | `string` | The unique identifier of the created node pool. |
| `status.outputs.node_ids` | `[]string` | The IDs of the individual nodes in the pool. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.cluster` | CivoKubernetesCluster | `metadata.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
