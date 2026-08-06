# CivoKubernetesCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1alpha1`

CivoKubernetesClusterSpec defines the specification for creating a managed Kubernetes cluster on Civo Cloud (K3s).
It focuses on essential parameters for a production-grade cluster, following the 80/20 principle to expose only the most commonly used settings.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.clusterName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.kubernetesVersion` | `string` | yes |  |  |
| `spec.network` | `string \| valueFrom` | yes |  | CivoVpc (`status.outputs.network_id`) |
| `spec.highlyAvailable` | `bool` |  | `false` |  |
| `spec.autoUpgrade` | `bool` |  |  |  |
| `spec.disableSurgeUpgrade` | `bool` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.defaultNodePool` | `CivoKubernetesClusterDefaultNodePool` | yes |  |  |
| `spec.defaultNodePool.size` | `string` | yes |  |  |
| `spec.defaultNodePool.nodeCount` | `uint32` | yes |  |  |

## Field Details

### spec.clusterName

`string` · required

The name of the Kubernetes cluster.
Constraints: Must be unique per account. No spaces allowed (alphanumeric and hyphens recommended). If left blank, a random name will be assigned.

- rule: {"required":true}

### spec.region

`enum` · required

The Civo region where the cluster will be created.
Determines where the cluster's control plane and nodes are provisioned.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `civo_region_unspecified` -- 0: default / unspecified region
- `lon1` -- london 1
- `lon2` -- london 2
- `fra1` -- frankfurt 1
- `nyc1` -- new york 1
- `phx1` -- phoenix 1
- `mum1` -- mumbai 1

### spec.kubernetesVersion

`string` · required

The Kubernetes version to use for the cluster (semantic versioning).
Must be a supported version on Civo (e.g., "1.26.3").

- rule: {"required":true}

### spec.network

`string | valueFrom` · required

Reference to the Civo network where the cluster will reside.
This must be an existing network in the same region. The network's ID is used for the cluster's networking.

- references: CivoVpc (`status.outputs.network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoVpc, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.highlyAvailable

`bool`

Whether to enable a highly available control plane for the cluster.
If true (when supported), the cluster is created with multiple master nodes for increased availability.
Default: false.

- default: `false`

### spec.autoUpgrade

`bool`

Whether to enable automatic Kubernetes version patch upgrades for the cluster.
If true, the cluster will automatically upgrade to new patch releases of Kubernetes when available.

### spec.disableSurgeUpgrade

`bool`

Whether to disable surge upgrades for the cluster.
If false (default), cluster upgrades may temporarily provision extra resources to minimize downtime during updates.

### spec.tags

`[]string`

A list of tags to apply to the cluster.
Tags help organize and identify the cluster within Civo.

### spec.defaultNodePool

`CivoKubernetesClusterDefaultNodePool` · required

Configuration for the cluster's default node pool.

- rule: {"required":true}

### spec.defaultNodePool.size

`string` · required

The instance size (node flavor) for each node in the default pool (e.g., "g4s.kube.medium").
This defines the CPU and memory for the nodes.

- rule: {"required":true}

### spec.defaultNodePool.nodeCount

`uint32` · required

The number of nodes to provision in the default node pool.
Must be at least 1.

- rule: {"required":true,"uint32":{"gt":0}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoKubernetesCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | The unique identifier (ID) of the created Kubernetes cluster. |
| `status.outputs.kubeconfig_b64` | `string` | A base64-encoded Kubernetes config (kubeconfig) for accessing the cluster. |
| `status.outputs.api_server_endpoint` | `string` | The endpoint URL of the Kubernetes API server for the cluster. |
| `status.outputs.created_at_rfc3339` | `string` | The timestamp when the cluster was created, in RFC 3339 format. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.network` | CivoVpc | `status.outputs.network_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| CivoKubernetesNodePool | `spec.cluster` | `metadata.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
