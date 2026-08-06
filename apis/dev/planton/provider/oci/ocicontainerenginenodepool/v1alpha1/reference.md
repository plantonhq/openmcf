# OciContainerEngineNodePool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciContainerEngineNodePoolSpec defines the specification for an Oracle Cloud
Infrastructure Container Engine for Kubernetes (OKE) node pool.

A node pool is a set of worker nodes within an OKE cluster that share the
same configuration: compute shape, OS image, Kubernetes version, and
placement rules. Multiple node pools can be attached to a single cluster,
enabling heterogeneous workloads (e.g., GPU nodes for ML, ARM nodes for
cost-optimized services, preemptible nodes for batch jobs).

Node placement is controlled via placement configs, each specifying an
availability domain, subnet, and optional fault domain constraints. For
regional subnets, provide one placement config per AD using the same subnet.

Deprecated provider features intentionally omitted:
  - node_image_id / node_image_name: use node_source_details instead
  - subnet_ids / quantity_per_subnet: use node_config_details instead

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.clusterId` | `string \| valueFrom` | yes |  | OciContainerEngineCluster (`status.outputs.cluster_id`) |
| `spec.name` | `string` |  |  |  |
| `spec.kubernetesVersion` | `string` |  |  |  |
| `spec.nodeShape` | `string` | yes |  |  |
| `spec.nodeShapeConfig` | `NodeShapeConfig` |  |  |  |
| `spec.nodeShapeConfig.ocpus` | `float` |  |  |  |
| `spec.nodeShapeConfig.memoryInGbs` | `float` |  |  |  |
| `spec.nodeSourceDetails` | `NodeSourceDetails` |  |  |  |
| `spec.nodeSourceDetails.imageId` | `string` | yes |  |  |
| `spec.nodeSourceDetails.bootVolumeSizeInGbs` | `int64` |  |  |  |
| `spec.nodeConfigDetails` | `NodeConfigDetails` | yes |  |  |
| `spec.nodeConfigDetails.placementConfigs` | `[]PlacementConfig` | yes |  |  |
| `spec.nodeConfigDetails.placementConfigs[].availabilityDomain` | `string` | yes |  |  |
| `spec.nodeConfigDetails.placementConfigs[].subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.nodeConfigDetails.placementConfigs[].faultDomains` | `[]string` |  |  |  |
| `spec.nodeConfigDetails.placementConfigs[].capacityReservationId` | `string \| valueFrom` |  |  |  |
| `spec.nodeConfigDetails.placementConfigs[].preemptibleNodeConfig` | `PreemptibleNodeConfig` |  |  |  |
| `spec.nodeConfigDetails.placementConfigs[].preemptibleNodeConfig.isPreserveBootVolume` | `bool` |  |  |  |
| `spec.nodeConfigDetails.size` | `int32` |  |  |  |
| `spec.nodeConfigDetails.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.nodeConfigDetails.kmsKeyId` | `string \| valueFrom` |  |  |  |
| `spec.nodeConfigDetails.isPvEncryptionInTransitEnabled` | `bool` |  |  |  |
| `spec.nodeConfigDetails.podNetworkOptionDetails` | `PodNetworkOptionDetails` |  |  |  |
| `spec.nodeConfigDetails.podNetworkOptionDetails.cniType` | `enum` |  |  |  |
| `spec.nodeConfigDetails.podNetworkOptionDetails.maxPodsPerNode` | `int32` |  |  |  |
| `spec.nodeConfigDetails.podNetworkOptionDetails.podNsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.nodeConfigDetails.podNetworkOptionDetails.podSubnetIds` | `[]string \| valueFrom` |  |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.sshPublicKey` | `string` |  |  |  |
| `spec.initialNodeLabels` | `[]NodeLabel` |  |  |  |
| `spec.initialNodeLabels[].key` | `string` |  |  |  |
| `spec.initialNodeLabels[].value` | `string` |  |  |  |
| `spec.nodeMetadata` | `map<string, string>` |  |  |  |
| `spec.nodeEvictionSettings` | `NodeEvictionSettings` |  |  |  |
| `spec.nodeEvictionSettings.evictionGraceDuration` | `string` |  |  |  |
| `spec.nodeEvictionSettings.isForceActionAfterGraceDuration` | `bool` |  |  |  |
| `spec.nodeEvictionSettings.isForceDeleteAfterGraceDuration` | `bool` |  |  |  |
| `spec.nodePoolCyclingDetails` | `NodePoolCyclingDetails` |  |  |  |
| `spec.nodePoolCyclingDetails.isNodeCyclingEnabled` | `bool` |  |  |  |
| `spec.nodePoolCyclingDetails.maximumSurge` | `string` |  |  |  |
| `spec.nodePoolCyclingDetails.maximumUnavailable` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the node pool will be created.
Changing this after creation forces node pool recreation.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.clusterId

`string | valueFrom` · required

OCID of the OKE cluster to which this node pool is attached.
Changing this after creation forces node pool recreation.

- references: OciContainerEngineCluster (`status.outputs.cluster_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciContainerEngineCluster, name: <that resource's name>, fieldPath: status.outputs.cluster_id}} -- a bare string does not parse

### spec.name

`string`

Human-readable name for the node pool shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.kubernetesVersion

`string`

Kubernetes version for the nodes. When omitted, inherits the cluster's
Kubernetes version. Set explicitly to pin a specific version or to
perform a rolling version upgrade of worker nodes independently of the
control plane.
Example: "v1.28.2".

### spec.nodeShape

`string` · required

Compute shape for all nodes in this pool.
Example: "VM.Standard.E4.Flex", "VM.Standard.A1.Flex", "VM.GPU.A10.1".
For flex shapes, also set node_shape_config to specify OCPUs and memory.

- rule: {"string":{"minLen":"1"}}

### spec.nodeShapeConfig

`NodeShapeConfig`

Shape configuration for flex shapes. Required when node_shape is a
flex shape (e.g., VM.Standard.E4.Flex). Ignored for fixed shapes.

### spec.nodeShapeConfig.ocpus

`float`

Number of OCPUs allocated to each node.
Example: 2.0 for a 2-OCPU flex instance.

### spec.nodeShapeConfig.memoryInGbs

`float`

Memory in gigabytes allocated to each node.
Example: 32.0 for 32 GB of RAM.

### spec.nodeSourceDetails

`NodeSourceDetails`

OS image and boot volume configuration for nodes.
When omitted, OKE uses the default Oracle Linux image for the cluster's
Kubernetes version.

### spec.nodeSourceDetails.imageId

`string` · required

OCID of the OCI platform image or custom image for the node OS.
Use `oci ce node-pool-options get` to list available images for
a given Kubernetes version.

- rule: {"string":{"minLen":"1"}}

### spec.nodeSourceDetails.bootVolumeSizeInGbs

`int64`

Boot volume size in gigabytes. Minimum 50 GB. When omitted, uses
the image's default boot volume size (typically 50 GB).

### spec.nodeConfigDetails

`NodeConfigDetails` · required

Node placement, sizing, networking, and encryption configuration.

- rule: {"required":true}

### spec.nodeConfigDetails.placementConfigs

`[]PlacementConfig` · required

Placement configurations determining which availability domains and
subnets receive nodes. Provide one entry per AD for regional subnets,
or one entry per AD-specific subnet.

- rule: {"repeated":{"minItems":"1"}}

### spec.nodeConfigDetails.placementConfigs[].availabilityDomain

`string` · required

Availability domain name where nodes will be launched.
Example: "Uocm:PHX-AD-1".

- rule: {"string":{"minLen":"1"}}

### spec.nodeConfigDetails.placementConfigs[].subnetId

`string | valueFrom` · required

OCID of the subnet in which to place nodes in this AD.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.nodeConfigDetails.placementConfigs[].faultDomains

`[]string`

Fault domains within the AD to constrain node placement.
When omitted, OKE distributes nodes across all fault domains.
Example: ["FAULT-DOMAIN-1", "FAULT-DOMAIN-2"]

### spec.nodeConfigDetails.placementConfigs[].capacityReservationId

`string | valueFrom`

OCID of a compute capacity reservation to use for nodes in this AD.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.nodeConfigDetails.placementConfigs[].preemptibleNodeConfig

`PreemptibleNodeConfig`

Preemptible node configuration. When set, nodes in this placement
use preemptible (spot) instances that can be reclaimed by OCI.
Suitable for fault-tolerant and batch workloads.

### spec.nodeConfigDetails.placementConfigs[].preemptibleNodeConfig.isPreserveBootVolume

`bool` · optional (explicit presence)

Whether to preserve the boot volume when the preemptible instance
is terminated. Defaults to false (boot volume is deleted).

### spec.nodeConfigDetails.size

`int32`

Desired number of nodes in this pool. OKE distributes nodes across
the placement configs as evenly as possible.

- rule: {"int32":{"gt":0}}

### spec.nodeConfigDetails.nsgIds

`[]string | valueFrom`

OCIDs of network security groups applied to the node VNICs.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.nodeConfigDetails.kmsKeyId

`string | valueFrom`

OCID of the KMS key for encrypting boot volumes at rest.
default_kind will be updated when OciKmsKey (R25) is implemented.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.nodeConfigDetails.isPvEncryptionInTransitEnabled

`bool`

Whether to enable in-transit encryption for the data volume's
paravirtualized attachment. Applies to both boot and block volumes.

### spec.nodeConfigDetails.podNetworkOptionDetails

`PodNetworkOptionDetails`

Pod networking configuration. Required when the cluster uses
OCI VCN-native pod networking (oci_vcn_ip_native CNI).

### spec.nodeConfigDetails.podNetworkOptionDetails.cniType

`enum`

CNI plugin type. Must match the cluster's CNI configuration.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `cni_unspecified`
- `flannel_overlay`
- `oci_vcn_ip_native`

### spec.nodeConfigDetails.podNetworkOptionDetails.maxPodsPerNode

`int32`

Maximum number of pods per node. Limited by the number of VNICs
attachable to the node shape. Only applicable for oci_vcn_ip_native.

### spec.nodeConfigDetails.podNetworkOptionDetails.podNsgIds

`[]string | valueFrom`

OCIDs of NSGs applied to pod VNICs. Only applicable for
oci_vcn_ip_native.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.nodeConfigDetails.podNetworkOptionDetails.podSubnetIds

`[]string | valueFrom`

OCIDs of subnets for pod IP allocation. Only applicable for
oci_vcn_ip_native. Can be the same as or different from the node
subnets.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.sshPublicKey

`string`

SSH public key installed on each node for debug access.
The corresponding private key allows SSH to nodes via their private IP
(or public IP if the subnet allows it).

### spec.initialNodeLabels

`[]NodeLabel`

Kubernetes labels applied to each node after it joins the cluster.
Commonly used for scheduling constraints (nodeSelector, affinity rules).
Example: [{ key: "workload-type", value: "gpu" }]

### spec.initialNodeLabels[].key

`string`

### spec.initialNodeLabels[].value

`string`

### spec.nodeMetadata

`map<string, string>`

Key/value pairs added to each underlying OCI compute instance at launch.
Used for cloud-init user data and instance metadata configuration.

### spec.nodeEvictionSettings

`NodeEvictionSettings`

Controls graceful node eviction behavior during node pool operations
(scale-down, version upgrades, shape changes).

### spec.nodeEvictionSettings.evictionGraceDuration

`string`

Maximum time OKE will attempt to evict pods before giving up.
ISO 8601 duration format. Default: PT60M. Range: PT0M to PT60M.
PT0M means delete the node immediately without cordon and drain.

### spec.nodeEvictionSettings.isForceActionAfterGraceDuration

`bool` · optional (explicit presence)

Whether to proceed with the node action if not all pods can be
evicted within the grace period.

### spec.nodeEvictionSettings.isForceDeleteAfterGraceDuration

`bool` · optional (explicit presence)

Whether to delete the underlying compute instance if pods cannot
be fully evicted within the grace period.

### spec.nodePoolCyclingDetails

`NodePoolCyclingDetails`

Rolling update strategy for node pool operations. Controls how many
nodes can be replaced simultaneously during upgrades.

### spec.nodePoolCyclingDetails.isNodeCyclingEnabled

`bool`

Whether node cycling is enabled for this pool.

### spec.nodePoolCyclingDetails.maximumSurge

`string`

Maximum additional nodes that can be temporarily created during
cycling. Accepts an integer ("1") or percentage ("25%").
Default: "1".

### spec.nodePoolCyclingDetails.maximumUnavailable

`string`

Maximum nodes that can be unavailable during cycling.
Accepts an integer ("0") or percentage ("25%").
Default: "0".

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciContainerEngineNodePool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.node_pool_id` | `string` | OCID of the node pool. |
| `status.outputs.kubernetes_version` | `string` | Kubernetes version running on the nodes in this pool. Matches the cluster version when not explicitly overridden. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.clusterId` | OciContainerEngineCluster | `status.outputs.cluster_id` |
| `spec.nodeConfigDetails.placementConfigs[].subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.nodeConfigDetails.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |
| `spec.nodeConfigDetails.podNetworkOptionDetails.podNsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |
| `spec.nodeConfigDetails.podNetworkOptionDetails.podSubnetIds` | OciSubnet | `status.outputs.subnet_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
