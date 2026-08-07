# AliCloudKubernetesNodePool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudKubernetesNodePoolSpec defines the configuration for an ACK node pool.

A node pool is a group of worker nodes within an ACK Managed Kubernetes
cluster that share the same instance type, scaling policy, and node
configuration. Each node pool has its own lifecycle and can be scaled
independently of the cluster.

This component wraps a single provider resource:
  Terraform: alicloud_cs_kubernetes_node_pool
  Pulumi:    cs.NodePool

Node pools inherit the VPC and Kubernetes version from the parent cluster.
The region is determined by the VSwitches; no region field is needed on the
node pool itself (the provider region comes from the parent cluster's
provider configuration).

For cost-optimized workloads, use spot_strategy with multiple instance_types
to spread across spot instance pools.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudKubernetesNodePool
metadata:
  name: alicloudkubernetesnodepool-demo
spec:
  region: cn-hangzhou
  clusterId:
    value: c-test-cluster
  name: demo-pool
  vswitchIds:
    - value: vsw-aaa111
    - value: vsw-bbb222
  instanceTypes:
    - ecs.g7.xlarge
  desiredSize: 2
  keyName: test-keypair
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.clusterId` | `string \| valueFrom` | yes |  | AliCloudKubernetesCluster (`status.outputs.cluster_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.vswitchIds` | `[]string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.instanceTypes` | `[]string` | yes |  |  |
| `spec.desiredSize` | `int32` |  |  |  |
| `spec.imageType` | `string` |  | `AliyunLinux3` |  |
| `spec.systemDisk` | `AliCloudKubernetesNodePoolSystemDisk` |  |  |  |
| `spec.systemDisk.category` | `string` |  | `cloud_essd` |  |
| `spec.systemDisk.size` | `int32` |  | `120` |  |
| `spec.systemDisk.performanceLevel` | `string` |  |  |  |
| `spec.systemDisk.encrypted` | `bool` |  |  |  |
| `spec.systemDisk.kmsKeyId` | `string` |  |  |  |
| `spec.dataDisks` | `[]AliCloudKubernetesNodePoolDataDisk` |  |  |  |
| `spec.dataDisks[].category` | `string` |  | `cloud_essd` |  |
| `spec.dataDisks[].size` | `int32` | yes |  |  |
| `spec.dataDisks[].name` | `string` |  |  |  |
| `spec.dataDisks[].performanceLevel` | `string` |  |  |  |
| `spec.dataDisks[].encrypted` | `string` |  |  |  |
| `spec.dataDisks[].kmsKeyId` | `string` |  |  |  |
| `spec.securityGroupIds` | `[]string \| valueFrom` |  |  | AliCloudSecurityGroup (`status.outputs.security_group_id`) |
| `spec.internetMaxBandwidthOut` | `int32` |  |  |  |
| `spec.internetChargeType` | `string` |  |  |  |
| `spec.keyName` | `string` |  |  |  |
| `spec.password` | `string` (sensitive) |  |  |  |
| `spec.labels` | `map<string, string>` |  |  |  |
| `spec.taints` | `[]AliCloudKubernetesNodePoolTaint` |  |  |  |
| `spec.taints[].key` | `string` | yes |  |  |
| `spec.taints[].value` | `string` |  |  |  |
| `spec.taints[].effect` | `string` |  |  |  |
| `spec.cpuPolicy` | `string` |  |  |  |
| `spec.runtimeName` | `string` |  |  |  |
| `spec.runtimeVersion` | `string` |  |  |  |
| `spec.unschedulable` | `bool` |  |  |  |
| `spec.userData` | `string` |  |  |  |
| `spec.installCloudMonitor` | `bool` |  | `true` |  |
| `spec.scalingConfig` | `AliCloudKubernetesNodePoolScalingConfig` |  |  |  |
| `spec.scalingConfig.enable` | `bool` |  | `true` |  |
| `spec.scalingConfig.minSize` | `int32` |  |  |  |
| `spec.scalingConfig.maxSize` | `int32` |  |  |  |
| `spec.scalingConfig.type` | `string` |  |  |  |
| `spec.multiAzPolicy` | `string` |  |  |  |
| `spec.management` | `AliCloudKubernetesNodePoolManagement` |  |  |  |
| `spec.management.enable` | `bool` |  | `true` |  |
| `spec.management.autoRepair` | `bool` |  |  |  |
| `spec.management.autoUpgrade` | `bool` |  |  |  |
| `spec.management.maxUnavailable` | `int32` |  |  |  |
| `spec.spotStrategy` | `string` |  |  |  |
| `spec.spotPriceLimits` | `[]AliCloudKubernetesNodePoolSpotPriceLimit` |  |  |  |
| `spec.spotPriceLimits[].instanceType` | `string` |  |  |  |
| `spec.spotPriceLimits[].priceLimit` | `string` |  |  |  |
| `spec.instanceChargeType` | `string` |  | `PostPaid` |  |
| `spec.period` | `int32` |  |  |  |
| `spec.autoRenew` | `bool` |  |  |  |
| `spec.autoRenewPeriod` | `int32` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.ramRoleName` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region. Must match the parent cluster's region.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.clusterId

`string | valueFrom` · required

ACK cluster ID that this node pool belongs to.

- references: AliCloudKubernetesCluster (`status.outputs.cluster_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudKubernetesCluster, name: <that resource's name>, fieldPath: status.outputs.cluster_id}} -- a bare string does not parse

### spec.name

`string` · required

Node pool name. 1-63 characters.
Maps to provider field `node_pool_name`.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63"}}

### spec.vswitchIds

`[]string | valueFrom` · required

VSwitch IDs for worker node placement. 1-5 VSwitches, preferably in
distinct availability zones for high availability.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true,"repeated":{"minItems":"1","maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.instanceTypes

`[]string` · required

ECS instance types for worker nodes. At least one type is required.
Specifying multiple types improves availability, especially for spot
instances where a single type may be unavailable.
Examples: "ecs.g7.xlarge", "ecs.g7.2xlarge", "ecs.c7.xlarge"

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.desiredSize

`int32` · optional (explicit presence)

Desired number of nodes in the node pool.
For auto-scaling pools (scaling_config.enable = true), this sets the
initial count; the auto-scaler may adjust between min_size and max_size.
For fixed-size pools, this is the exact count.
Range: 0-1000.

- rule: {"int32":{"lte":1000,"gte":0}}

### spec.imageType

`string` · optional (explicit presence)

OS image type for the worker nodes.
Default: "AliyunLinux3"

- default: `AliyunLinux3`
- rule: image_type must be one of: AliyunLinux, AliyunLinux3, AliyunLinux3Arm64, AliyunLinuxUEFI, CentOS, Windows, WindowsCore, ContainerOS, Ubuntu, AliyunLinux3ContainerOptimized, Custom

### spec.systemDisk

`AliCloudKubernetesNodePoolSystemDisk`

System disk configuration for the worker nodes.

### spec.systemDisk.category

`string` · optional (explicit presence)

Disk category.
Default: "cloud_essd"

- default: `cloud_essd`
- rule: category must be one of: cloud_efficiency, cloud_ssd, cloud_essd, cloud_auto

### spec.systemDisk.size

`int32` · optional (explicit presence)

Disk size in GiB.
Minimum: 40 GiB. Default: 120 GiB.

- default: `120`
- rule: {"int32":{"lte":500,"gte":40}}

### spec.systemDisk.performanceLevel

`string`

ESSD performance level. Only applies when category is "cloud_essd".
Valid values: "PL0", "PL1", "PL2", "PL3".
Higher levels provide more IOPS and throughput.

- rule: performance_level must be one of: PL0, PL1, PL2, PL3

### spec.systemDisk.encrypted

`bool` · optional (explicit presence)

Whether to encrypt the system disk.

### spec.systemDisk.kmsKeyId

`string`

KMS key ID for disk encryption.
Only relevant when encrypted is true.

### spec.dataDisks

`[]AliCloudKubernetesNodePoolDataDisk`

Additional data disks attached to each worker node.

### spec.dataDisks[].category

`string` · optional (explicit presence)

Disk category.
Default: "cloud_essd"

- default: `cloud_essd`
- rule: category must be one of: cloud_efficiency, cloud_ssd, cloud_essd, cloud_auto

### spec.dataDisks[].size

`int32` · required

Disk size in GiB. Range: 40-32767.

- rule: {"required":true,"int32":{"lte":32767,"gte":40}}

### spec.dataDisks[].name

`string`

Disk name.

### spec.dataDisks[].performanceLevel

`string`

ESSD performance level. Only applies when category is "cloud_essd".
Valid values: "PL0", "PL1", "PL2", "PL3".

- rule: performance_level must be one of: PL0, PL1, PL2, PL3

### spec.dataDisks[].encrypted

`string`

Whether to encrypt the data disk. "true" or "false".

### spec.dataDisks[].kmsKeyId

`string`

KMS key ID for disk encryption.

### spec.securityGroupIds

`[]string | valueFrom`

Security group IDs for the worker nodes.
Immutable after creation.
If omitted, the cluster's default security group is used.

- references: AliCloudSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.internetMaxBandwidthOut

`int32` · optional (explicit presence)

Maximum outbound internet bandwidth in Mbps.
0 means no public IP is allocated (default). 1-100 allocates a public IP.

- rule: {"int32":{"lte":100,"gte":0}}

### spec.internetChargeType

`string` · optional (explicit presence)

Billing method for public internet access.
Only relevant when internet_max_bandwidth_out > 0.
Default: "PayByTraffic"

- rule: internet_charge_type must be one of: PayByBandwidth, PayByTraffic

### spec.keyName

`string`

SSH key pair name for node access.
Mutually exclusive with password.
Recommended over password for managed node pools.

### spec.password

`string` · sensitive

SSH password for node access. 8-30 characters.
Mutually exclusive with key_name. Sensitive value.

### spec.labels

`map<string, string>`

Kubernetes labels applied to all nodes in the pool.
Used for pod scheduling affinity and node selectors.
Example: {"workload-type": "compute", "team": "platform"}

### spec.taints

`[]AliCloudKubernetesNodePoolTaint`

Kubernetes taints applied to all nodes in the pool.
Taints repel pods unless the pod has a matching toleration.

### spec.taints[].key

`string` · required

Taint key. Must be a valid Kubernetes label key.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.taints[].value

`string`

Taint value.

### spec.taints[].effect

`string`

Taint effect that determines the scheduling behavior.
"NoSchedule" -- pods without a toleration are not scheduled.
"PreferNoSchedule" -- soft version; scheduler tries to avoid.
"NoExecute" -- existing pods without a toleration are evicted.

- rule: effect must be one of: NoSchedule, PreferNoSchedule, NoExecute

### spec.cpuPolicy

`string` · optional (explicit presence)

CPU management policy for nodes.
"none" (default) -- standard CFS scheduling.
"static" -- pin exclusive containers to specific CPUs; improves
  determinism for latency-sensitive workloads.

- rule: cpu_policy must be one of: none, static

### spec.runtimeName

`string`

Container runtime name (e.g., "containerd", "Sandboxed-Container.runv").
If omitted, the provider default for the selected image_type is used.

### spec.runtimeVersion

`string`

Container runtime version.
If omitted, the latest version supported by the cluster is used.

### spec.unschedulable

`bool` · optional (explicit presence)

Whether newly added nodes are marked unschedulable.
When true, pods cannot be scheduled on new nodes until the taint is
manually removed. Useful for pre-configuration before accepting workloads.

### spec.userData

`string`

Custom user data script to execute on each node at boot time.
Must be base64-encoded. Maximum 16 KB before encoding.

### spec.installCloudMonitor

`bool` · optional (explicit presence)

Whether to install Alibaba Cloud CloudMonitor agent on nodes.
Default: true

- default: `true`

### spec.scalingConfig

`AliCloudKubernetesNodePoolScalingConfig`

Auto-scaling configuration. When enabled, the cluster auto-scaler
adjusts the number of nodes between min_size and max_size based on
pending pod resource requests.

### spec.scalingConfig.enable

`bool` · optional (explicit presence)

Whether auto-scaling is enabled.
Default: true (when scaling_config is present).

- default: `true`

### spec.scalingConfig.minSize

`int32`

Minimum number of nodes the auto-scaler will maintain.
Range: 0-1000.

- rule: {"int32":{"lte":1000,"gte":0}}

### spec.scalingConfig.maxSize

`int32`

Maximum number of nodes the auto-scaler can scale to.
Range: 0-2000. Must be >= min_size.

- rule: {"int32":{"lte":2000,"gte":0}}

### spec.scalingConfig.type

`string` · optional (explicit presence)

Instance classification for the auto-scaler.
"cpu" (default) -- general-purpose CPU instances.
"gpu" -- GPU instances.
"gpushare" -- shared GPU instances.
"spot" -- spot/preemptible instances.

- rule: type must be one of: cpu, gpu, gpushare, spot

### spec.multiAzPolicy

`string` · optional (explicit presence)

Multi-AZ scheduling policy when vswitch_ids spans multiple zones.
"PRIORITY" -- allocate in the first AZ with available capacity.
"COST_OPTIMIZED" -- allocate in the cheapest AZ.
"BALANCE" -- evenly distribute across AZs.

- rule: multi_az_policy must be one of: PRIORITY, COST_OPTIMIZED, BALANCE

### spec.management

`AliCloudKubernetesNodePoolManagement`

Managed node pool lifecycle management settings.
When enabled, ACK automatically repairs unhealthy nodes, upgrades
kubelet versions, and patches vulnerabilities.

### spec.management.enable

`bool` · optional (explicit presence)

Whether managed node pool features are enabled.
Default: true

- default: `true`

### spec.management.autoRepair

`bool` · optional (explicit presence)

Whether ACK automatically repairs unhealthy nodes.
Unhealthy nodes are detected via node conditions (Ready, DiskPressure, etc.).

### spec.management.autoUpgrade

`bool` · optional (explicit presence)

Whether ACK automatically upgrades kubelet on nodes when the cluster
version is upgraded.

### spec.management.maxUnavailable

`int32` · optional (explicit presence)

Maximum number of nodes that can be unavailable during managed
operations (repair, upgrade, patching).
Range: 0-1000. Default: 1.

- rule: {"int32":{"lte":1000,"gte":0}}

### spec.spotStrategy

`string` · optional (explicit presence)

Spot instance strategy for cost optimization.
"NoSpot" (default) -- only on-demand instances.
"SpotWithPriceLimit" -- spot instances with a price cap.
"SpotAsPriceGo" -- spot instances at market price (cheapest).

- rule: spot_strategy must be one of: NoSpot, SpotWithPriceLimit, SpotAsPriceGo

### spec.spotPriceLimits

`[]AliCloudKubernetesNodePoolSpotPriceLimit`

Per-instance-type price caps for spot instances.
Only relevant when spot_strategy is "SpotWithPriceLimit".

### spec.spotPriceLimits[].instanceType

`string`

ECS instance type (e.g., "ecs.g7.xlarge").
Must be one of the types listed in instance_types.

### spec.spotPriceLimits[].priceLimit

`string`

Maximum hourly price in CNY (e.g., "0.98").
When the spot market price exceeds this, the instance is not created.

### spec.instanceChargeType

`string` · optional (explicit presence)

Billing method for worker node instances.
"PostPaid" (default) -- pay-as-you-go.
"PrePaid" -- subscription; requires period and optionally auto_renew.

- default: `PostPaid`
- rule: instance_charge_type must be one of: PostPaid, PrePaid

### spec.period

`int32` · optional (explicit presence)

Subscription period in months. Required when instance_charge_type is PrePaid.
Valid values: 1, 2, 3, 6, 12.

- rule: period must be one of: 1, 2, 3, 6, 12

### spec.autoRenew

`bool` · optional (explicit presence)

Whether to automatically renew the subscription.
Only relevant when instance_charge_type is PrePaid.

### spec.autoRenewPeriod

`int32` · optional (explicit presence)

Auto-renewal period in months.
Valid values: 1, 2, 3, 6, 12.

- rule: auto_renew_period must be one of: 1, 2, 3, 6, 12

### spec.tags

`map<string, string>`

Tags applied to the ECS instances in the node pool.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping.

### spec.ramRoleName

`string`

RAM role name for worker nodes.
Determines what Alibaba Cloud APIs the node can call (e.g., pulling
images from ACR, writing to SLS). If omitted, the cluster's default
worker RAM role is used.
Immutable after creation.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudKubernetesNodePool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.node_pool_id` | `string` | Node pool ID assigned by Alibaba Cloud. |
| `status.outputs.scaling_group_id` | `string` | Auto Scaling group ID associated with this node pool. Can be used to query scaling activities and node status. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.clusterId` | AliCloudKubernetesCluster | `status.outputs.cluster_id` |
| `spec.vswitchIds` | AliCloudVswitch | `status.outputs.vswitch_id` |
| `spec.securityGroupIds` | AliCloudSecurityGroup | `status.outputs.security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
