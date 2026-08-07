# AliCloudEcsInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudEcsInstanceSpec defines the configuration for an Alibaba Cloud
Elastic Compute Service (ECS) instance.

ECS is the fundamental compute building block on Alibaba Cloud, analogous
to AWS EC2 or Azure Virtual Machines. Each instance runs a single OS image
with configurable CPU, memory, system disk, optional data disks, and
networking (VPC placement, security groups, optional public IP).

Data disks are created inline with the instance using the provider's
built-in `data_disks` block, keeping disk lifecycle tightly coupled to
the instance (consistent with DD07 composite bundling).

Provider resources:
  Terraform: alicloud_instance
  Pulumi:    ecs.Instance

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudEcsInstance
metadata:
  name: alicloudecsinstance-demo
spec:
  region: cn-hangzhou
  instanceType: ecs.g7.large
  imageId: ubuntu_22_04_x64_20G_alibase_20230515.vhd
  vswitchId:
    value: vsw-demo123
  securityGroupIds:
    - value: sg-demo123
  keyName: demo-keypair
  systemDisk:
    category: cloud_essd
    size: 40
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.securityGroupIds` | `[]string \| valueFrom` | yes |  | AliCloudSecurityGroup (`status.outputs.security_group_id`) |
| `spec.instanceType` | `string` | yes |  |  |
| `spec.imageId` | `string` | yes |  |  |
| `spec.instanceName` | `string` |  |  |  |
| `spec.hostName` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.systemDisk` | `AliCloudEcsSystemDisk` |  |  |  |
| `spec.systemDisk.category` | `string` |  | `cloud_essd` |  |
| `spec.systemDisk.size` | `int32` |  | `40` |  |
| `spec.systemDisk.performanceLevel` | `string` |  |  |  |
| `spec.systemDisk.encrypted` | `bool` |  |  |  |
| `spec.systemDisk.kmsKeyId` | `string` |  |  |  |
| `spec.dataDisks` | `[]AliCloudEcsDataDisk` |  |  |  |
| `spec.dataDisks[].size` | `int32` | yes |  |  |
| `spec.dataDisks[].category` | `string` |  | `cloud_essd` |  |
| `spec.dataDisks[].name` | `string` |  |  |  |
| `spec.dataDisks[].performanceLevel` | `string` |  |  |  |
| `spec.dataDisks[].encrypted` | `bool` |  |  |  |
| `spec.dataDisks[].kmsKeyId` | `string` |  |  |  |
| `spec.dataDisks[].snapshotId` | `string` |  |  |  |
| `spec.dataDisks[].deleteWithInstance` | `bool` |  | `true` |  |
| `spec.dataDisks[].description` | `string` |  |  |  |
| `spec.keyName` | `string` |  |  |  |
| `spec.password` | `string` (sensitive) |  |  |  |
| `spec.internetMaxBandwidthOut` | `int32` |  | `0` |  |
| `spec.internetChargeType` | `string` |  |  |  |
| `spec.instanceChargeType` | `string` |  | `PostPaid` |  |
| `spec.period` | `int32` |  |  |  |
| `spec.periodUnit` | `string` |  |  |  |
| `spec.spotStrategy` | `string` |  |  |  |
| `spec.spotPriceLimit` | `double` |  |  |  |
| `spec.userData` | `string` |  |  |  |
| `spec.roleName` | `string` |  |  |  |
| `spec.deletionProtection` | `bool` |  |  |  |
| `spec.securityEnhancementStrategy` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the ECS instance will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vswitchId

`string | valueFrom` · required

VSwitch ID that determines the VPC, availability zone, and subnet CIDR
for the instance's primary network interface.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.securityGroupIds

`[]string | valueFrom` · required

Security group IDs to associate with the instance. At least one is
required. An instance can belong to up to 5 security groups.

- references: AliCloudSecurityGroup (`status.outputs.security_group_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.instanceType

`string` · required

ECS instance type that determines CPU and memory allocation.
Must start with "ecs." prefix.
Examples: "ecs.g7.large" (2 vCPU, 8 GiB), "ecs.c7.xlarge" (4 vCPU, 8 GiB).

- rule: instance_type must start with 'ecs.'
- rule: {"required":true}

### spec.imageId

`string` · required

Operating system image ID. Determines the OS, architecture, and
pre-installed software on the instance.
Example: "ubuntu_22_04_x64_20G_alibase_20230515.vhd"

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.instanceName

`string`

Instance display name. 2-128 characters.
If omitted, defaults to the metadata.name.

- rule: instance_name must be between 2 and 128 characters when set

### spec.hostName

`string`

OS-level hostname. If omitted, Alibaba Cloud generates one.

### spec.description

`string`

Instance description. 2-256 characters.

- rule: description must be between 2 and 256 characters when set

### spec.systemDisk

`AliCloudEcsSystemDisk`

System disk configuration. Every ECS instance has exactly one system
disk that holds the OS. If omitted, defaults are applied (cloud_essd, 40 GB).

### spec.systemDisk.category

`string` · optional (explicit presence)

Disk category. Determines the underlying storage technology and
performance characteristics.
Default: "cloud_essd"

- default: `cloud_essd`
- rule: category must be one of: cloud_efficiency, cloud_ssd, cloud_essd, cloud_auto, cloud_essd_entry

### spec.systemDisk.size

`int32` · optional (explicit presence)

Disk size in GB. The minimum depends on the image; typical minimum
is 20 GB for Linux and 40 GB for Windows.
Default: 40

- default: `40`
- rule: {"int32":{"gte":20}}

### spec.systemDisk.performanceLevel

`string` · optional (explicit presence)

ESSD performance level. Only applicable when category is "cloud_essd".
Higher levels deliver more IOPS and throughput at higher cost.

- rule: performance_level must be one of: PL0, PL1, PL2, PL3

### spec.systemDisk.encrypted

`bool` · optional (explicit presence)

Enable server-side encryption for the system disk.

### spec.systemDisk.kmsKeyId

`string`

KMS key ID for disk encryption. Only applicable when encrypted is true.

### spec.dataDisks

`[]AliCloudEcsDataDisk`

Additional data disks to create and attach to the instance. Up to 16
data disks are supported. Disks are created inline with the instance
and their lifecycle is tied to it by default.

- rule: {"repeated":{"maxItems":"16"}}

### spec.dataDisks[].size

`int32` · required

Disk size in GB. Required for each data disk.

- rule: {"required":true,"int32":{"gte":20}}

### spec.dataDisks[].category

`string` · optional (explicit presence)

Disk category.
Default: "cloud_essd"

- default: `cloud_essd`
- rule: category must be one of: cloud_efficiency, cloud_ssd, cloud_essd, cloud_auto, cloud_essd_entry

### spec.dataDisks[].name

`string`

Disk display name.

### spec.dataDisks[].performanceLevel

`string` · optional (explicit presence)

ESSD performance level. Only applicable when category is "cloud_essd".

- rule: performance_level must be one of: PL0, PL1, PL2, PL3

### spec.dataDisks[].encrypted

`bool` · optional (explicit presence)

Enable server-side encryption for this data disk.

### spec.dataDisks[].kmsKeyId

`string`

KMS key ID for disk encryption.

### spec.dataDisks[].snapshotId

`string`

Snapshot ID to create this disk from.

### spec.dataDisks[].deleteWithInstance

`bool` · optional (explicit presence)

Whether this disk is deleted when the instance is released.
Default: true

- default: `true`

### spec.dataDisks[].description

`string`

Disk description.

### spec.keyName

`string`

SSH key pair name for passwordless authentication. Mutually exclusive
authentication option with password.

### spec.password

`string` · sensitive

Instance login password. 8-30 characters, must contain at least three
of: uppercase, lowercase, digits, special characters.
Mutually exclusive authentication option with key_name.

- rule: password must be between 8 and 30 characters when set

### spec.internetMaxBandwidthOut

`int32` · optional (explicit presence)

Maximum outbound internet bandwidth in Mbps. A value > 0 causes Alibaba
Cloud to allocate a public IP address to the instance.
Default: 0 (no public IP).

- default: `0`
- rule: {"int32":{"lte":100,"gte":0}}

### spec.internetChargeType

`string` · optional (explicit presence)

Billing method for internet traffic. Only meaningful when
internet_max_bandwidth_out > 0.

- rule: internet_charge_type must be one of: PayByTraffic, PayByBandwidth

### spec.instanceChargeType

`string` · optional (explicit presence)

Billing method for the instance itself.
Default: "PostPaid" (pay-as-you-go).

- default: `PostPaid`
- rule: instance_charge_type must be one of: PostPaid, PrePaid

### spec.period

`int32` · optional (explicit presence)

Subscription period. Only applicable when instance_charge_type is "PrePaid".
Valid values: 1-9, 12, 24, 36, 48, 60.

- rule: period must be one of: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, 60

### spec.periodUnit

`string` · optional (explicit presence)

Subscription period unit. Only applicable when instance_charge_type is "PrePaid".

- rule: period_unit must be one of: Week, Month

### spec.spotStrategy

`string` · optional (explicit presence)

Spot instance bidding strategy. Spot instances can reduce costs up to 90%
but may be reclaimed by Alibaba Cloud when capacity is needed.

- rule: spot_strategy must be one of: NoSpot, SpotAsPriceGo, SpotWithPriceLimit

### spec.spotPriceLimit

`double` · optional (explicit presence)

Maximum hourly price for a spot instance. Only applicable when
spot_strategy is "SpotWithPriceLimit".

### spec.userData

`string`

Cloud-init user data script. Typically base64-encoded. Executed once
at first boot for automated instance provisioning.

### spec.roleName

`string`

RAM role name to attach as the instance profile. Enables the instance
to call Alibaba Cloud APIs without embedded credentials.

### spec.deletionProtection

`bool` · optional (explicit presence)

Prevent accidental instance deletion via console or API.

### spec.securityEnhancementStrategy

`string` · optional (explicit presence)

Cloud security center agent mode. "Active" enables the China-region
security agent; "Deactive" disables it.

- rule: security_enhancement_strategy must be one of: Active, Deactive

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the instance is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the ECS instance and its disks.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudEcsInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The ECS instance ID assigned by Alibaba Cloud (e.g., "i-bp1xxxxx"). |
| `status.outputs.private_ip` | `string` | The primary private IP address assigned to the instance within its VSwitch. |
| `status.outputs.public_ip` | `string` | The public IP address allocated to the instance. Empty when internet_max_bandwidth_out is 0 (no public IP requested). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |
| `spec.securityGroupIds` | AliCloudSecurityGroup | `status.outputs.security_group_id` |

## See Also

- [Overview](../README.md)
