# AliCloudVpc

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudVpcSpec defines the configuration for an Alibaba Cloud Virtual
Private Cloud (VPC).

A VPC is the networking foundation for nearly every other Alibaba Cloud
resource. It provides an isolated virtual network with its own CIDR block,
route table, and virtual router. VSwitches (subnets), security groups, NAT
gateways, load balancers, database instances, and Kubernetes clusters are
all deployed into a VPC.

This component creates a single VPC. VSwitches, NAT gateways, and other
networking resources are managed as separate components, keeping the VPC
itself a clean, composable building block.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudVpc
metadata:
  name: alicloudvpc-demo
spec:
  region: cn-hangzhou
  vpcName: alicloudvpc-demo
  cidrBlock: "10.0.0.0/16"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcName` | `string` | yes |  |  |
| `spec.cidrBlock` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.enableIpv6` | `bool` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the VPC will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpcName

`string` · required

VPC name. 1-128 characters; cannot start with http:// or https://.
Maps to the provider field `vpc_name`.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.cidrBlock

`string` · required

Primary IPv4 CIDR block for the VPC.
Must be a valid private CIDR in one of the following ranges:
  10.0.0.0/8   - 10.255.255.255
  172.16.0.0/12 - 172.31.255.255
  192.168.0.0/16 - 192.168.255.255
Mask length: 8-28.
Example: "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.description

`string`

Human-readable description of the VPC.
1-256 characters; cannot start with http:// or https://.

### spec.enableIpv6

`bool`

Enable IPv6 for this VPC.
When enabled, Alibaba Cloud allocates an IPv6 CIDR block (/56) to the VPC.
VSwitches within this VPC can then be assigned IPv6 CIDR blocks.
Default: false

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the VPC is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the VPC.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudVpc, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpc_id` | `string` | The VPC ID assigned by Alibaba Cloud. Referenced by downstream components (VSwitch, SecurityGroup, NAT, ACK, etc.) via StringValueOrRef. |
| `status.outputs.vpc_name` | `string` | The VPC name as created. |
| `status.outputs.cidr_block` | `string` | The primary IPv4 CIDR block of the VPC. Useful for downstream VSwitch CIDR planning and security group rules. |
| `status.outputs.router_id` | `string` | The virtual router ID automatically created with the VPC. Each VPC gets one VRouter that manages route tables. |
| `status.outputs.route_table_id` | `string` | The system route table ID associated with the VPC. Every VPC has a default system route table that is auto-created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AliCloudApplicationLoadBalancer | `spec.vpcId` | `status.outputs.vpc_id` |
| AliCloudCenInstance | `spec.attachments[].childInstanceId` | `status.outputs.vpc_id` |
| AliCloudFunction | `spec.vpcConfig.vpcId` | `status.outputs.vpc_id` |
| AliCloudNasFileSystem | `spec.vpcId` | `status.outputs.vpc_id` |
| AliCloudNatGateway | `spec.vpcId` | `status.outputs.vpc_id` |
| AliCloudNetworkLoadBalancer | `spec.vpcId` | `status.outputs.vpc_id` |
| AliCloudPrivateDnsZone | `spec.vpcAttachments[].vpcId` | `status.outputs.vpc_id` |
| AliCloudRocketmqInstance | `spec.vpcId` | `status.outputs.vpc_id` |
| AliCloudSaeApplication | `spec.vpcId` | `status.outputs.vpc_id` |
| AliCloudSecurityGroup | `spec.vpcId` | `status.outputs.vpc_id` |
| AliCloudVpnGateway | `spec.vpcId` | `status.outputs.vpc_id` |
| AliCloudVswitch | `spec.vpcId` | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
