# AliCloudVswitch

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudVswitchSpec defines the configuration for an Alibaba Cloud VSwitch.

A VSwitch is the subnet-equivalent resource in Alibaba Cloud networking.
It carves out a CIDR range within a VPC and pins it to a single availability
zone. ECS instances, RDS databases, container clusters, NAT gateways, and
most other VPC-aware resources are deployed into a VSwitch.

Each VSwitch belongs to exactly one VPC and one availability zone. The CIDR
block must be a subset of the parent VPC's CIDR block. Changing the VPC,
zone, or CIDR after creation requires destroying and recreating the VSwitch.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudVswitch
metadata:
  name: alicloudvswitch-demo
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-replace-me
  zoneId: cn-hangzhou-a
  cidrBlock: "10.0.0.0/24"
  vswitchName: planton-demo-vswitch
  description: Demo VSwitch for local testing
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.zoneId` | `string` | yes |  |  |
| `spec.cidrBlock` | `string` | yes |  |  |
| `spec.vswitchName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.enableIpv6` | `bool` |  |  |  |
| `spec.ipv6CidrBlockMask` | `int32` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the VSwitch will be created.
Must match the region of the parent VPC.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

VPC ID that this VSwitch belongs to.
The VSwitch's CIDR block must fall within the VPC's CIDR range.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.zoneId

`string` · required

Availability zone for the VSwitch.
Format: "{region}-{letter}" (e.g., "cn-hangzhou-a", "cn-hangzhou-b").
The zone must exist in the VSwitch's region. All resources deployed into
this VSwitch run in this zone.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.cidrBlock

`string` · required

IPv4 CIDR block for the VSwitch.
Must be a subnet of the parent VPC's CIDR block.
Mask length: 16-29 (e.g., "10.0.0.0/24" gives 256 addresses).
This field is immutable after creation.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vswitchName

`string` · required

VSwitch name. 1-128 characters; cannot start with http:// or https://.
Maps to the provider field `vswitch_name`.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.description

`string`

Human-readable description of the VSwitch.
1-256 characters; cannot start with http:// or https://.

### spec.enableIpv6

`bool`

Enable IPv6 for this VSwitch.
The parent VPC must have IPv6 enabled (enable_ipv6 = true) for this
to take effect. When enabled, an IPv6 CIDR block is allocated to the
VSwitch based on ipv6_cidr_block_mask.
Default: false

### spec.ipv6CidrBlockMask

`int32`

IPv6 CIDR block mask for the VSwitch.
Only meaningful when enable_ipv6 is true and the parent VPC has IPv6
enabled. Valid range: 0-255. The mask selects a /64 segment from the
VPC's /56 IPv6 allocation.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":255,"gte":0}}

### spec.tags

`map<string, string>`

Tags to apply to the VSwitch.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudVswitch, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vswitch_id` | `string` | The VSwitch ID assigned by Alibaba Cloud. Referenced by downstream components (NatGateway, EcsInstance, AckManagedCluster, RdsInstance, PolardbCluster, RedisInstance, etc.) via StringValueOrRef. |
| `status.outputs.vswitch_name` | `string` | The VSwitch name as created. |
| `status.outputs.cidr_block` | `string` | The IPv4 CIDR block of the VSwitch. Useful for security group rules and network planning. |
| `status.outputs.zone_id` | `string` | The availability zone in which the VSwitch resides. |
| `status.outputs.ipv6_cidr_block` | `string` | The IPv6 CIDR block allocated to the VSwitch. Only populated when IPv6 is enabled on both the parent VPC and this VSwitch. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AliCloudApplicationLoadBalancer | `spec.zoneMappings[].vswitchId` | `status.outputs.vswitch_id` |
| AliCloudEcsInstance | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudKubernetesCluster | `spec.vswitchIds` | `status.outputs.vswitch_id` |
| AliCloudKubernetesCluster | `spec.podVswitchIds` | `status.outputs.vswitch_id` |
| AliCloudKubernetesNodePool | `spec.vswitchIds` | `status.outputs.vswitch_id` |
| AliCloudMongodbInstance | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudNasFileSystem | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudNatGateway | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudNatGateway | `spec.snatEntries[].sourceVswitchId` | `status.outputs.vswitch_id` |
| AliCloudNetworkLoadBalancer | `spec.zoneMappings[].vswitchId` | `status.outputs.vswitch_id` |
| AliCloudPolardbCluster | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudRdsInstance | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudRedisInstance | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudRocketmqInstance | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudSaeApplication | `spec.vswitchId` | `status.outputs.vswitch_id` |
| AliCloudVpnGateway | `spec.vswitchId` | `status.outputs.vswitch_id` |

## See Also

- [Overview](../README.md)
