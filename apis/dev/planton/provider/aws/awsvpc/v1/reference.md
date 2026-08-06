# AwsVpc

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment

**apiVersion**: `aws.planton.dev/v1`

AwsVpcSpec defines an AWS Virtual Private Cloud (VPC): an isolated virtual
network in which other AWS resources are launched.

A VPC is the networking foundation for nearly everything else on AWS, but it
is itself a thin construct: an IP address space (one primary IPv4 CIDR, plus
optional secondary IPv4 CIDRs and an IPv6 CIDR), a tenancy mode, and a few
DNS toggles. The things that make a network useful -- subnets, internet
gateways, NAT gateways, route tables -- are separate, independently ownable
components that reference this VPC. Keeping the VPC a clean building block is
what lets a topology be composed from first-class nodes rather than hidden
inside one opaque resource.

The primary IPv4 CIDR can be specified explicitly (cidr_block) or allocated
from an IPAM pool (ipv4_ipam_pool_id [+ ipv4_netmask_length]). IPv6 is opt-in
and supports either an Amazon-provided /56 (assign_generated_ipv6_cidr_block)
or an IPAM-allocated block (ipv6_ipam_pool_id). The cross-field rules below
mirror what AWS itself enforces.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsVpc
metadata:
  name: awsvpc-demo
spec:
  region: us-west-2
  cidrBlock: 10.0.0.0/16
  enableDnsHostnames: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.cidrBlock` | `string` |  |  |  |
| `spec.secondaryIpv4CidrBlocks` | `[]string` |  |  |  |
| `spec.ipv4IpamPoolId` | `string` |  |  |  |
| `spec.ipv4NetmaskLength` | `int32` |  |  |  |
| `spec.instanceTenancy` | `string` |  | `default` |  |
| `spec.enableDnsSupport` | `bool` |  | `true` |  |
| `spec.enableDnsHostnames` | `bool` |  | `true` |  |
| `spec.enableNetworkAddressUsageMetrics` | `bool` |  |  |  |
| `spec.assignGeneratedIpv6CidrBlock` | `bool` |  |  |  |
| `spec.ipv6CidrBlock` | `string` |  |  |  |
| `spec.ipv6CidrBlockNetworkBorderGroup` | `string` |  |  |  |
| `spec.ipv6IpamPoolId` | `string` |  |  |  |
| `spec.ipv6NetmaskLength` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the VPC will be created.
Example: "us-west-2", "eu-west-1", "ap-southeast-1".
For the list of regions, see:
https://aws.amazon.com/about-aws/global-infrastructure/regions_az/

- rule: {"string":{"minLen":"1"}}

### spec.cidrBlock

`string`

The primary IPv4 CIDR block for the VPC. This is the VPC's main address
range and is immutable: changing it replaces the VPC. The mask must be
between /16 and /28 (AWS limits), e.g. "10.0.0.0/16" yields 65,536
addresses. Leave empty only when allocating the primary CIDR from IPAM
(set ipv4_ipam_pool_id instead); exactly one of cidr_block or
ipv4_ipam_pool_id must be provided.

### spec.secondaryIpv4CidrBlocks

`[]string`

Additional IPv4 CIDR blocks to associate with the VPC beyond the primary
one (e.g. ["10.1.0.0/16", "100.64.0.0/16"]). Each is associated as its own
resource and can be added or removed without recreating the VPC. Each block
must be non-overlapping and within the AWS /16-/28 mask limits. Useful when
a single /16 runs out of subnet space or when carving a distinct range
(such as the 100.64.0.0/10 shared space) for specific workloads.

### spec.ipv4IpamPoolId

`string`

The IPAM (IP Address Manager) pool to allocate the primary IPv4 CIDR from,
instead of specifying cidr_block directly. Pair with ipv4_netmask_length to
let IPAM choose a block of the requested size, or with cidr_block to take a
specific block from the pool. Immutable: changing it replaces the VPC.

### spec.ipv4NetmaskLength

`int32`

The netmask length of the primary IPv4 CIDR to allocate from
ipv4_ipam_pool_id (e.g. 16 for a /16). Must be between 16 and 28. Requires
ipv4_ipam_pool_id and is mutually exclusive with cidr_block. Immutable:
changing it replaces the VPC.

### spec.instanceTenancy

`string`

The tenancy of instances launched into this VPC: "default" (instances may
run on shared hardware) or "dedicated" (every instance runs on
single-tenant hardware, at higher cost). Leave empty for the AWS default
("default"). Note that AWS only supports changing "dedicated" -> "default"
in place; other tenancy changes are not permitted by the API.

- default: `default`
- rule: instance_tenancy must be one of: default, dedicated

### spec.enableDnsSupport

`bool` · optional (explicit presence)

Whether the Amazon-provided DNS server (the .2 resolver) answers DNS
queries from within the VPC. AWS enables this by default, and disabling it
breaks name resolution for most workloads, so when this field is unset the
VPC keeps DNS resolution ON. Set it explicitly to false only to deliberately
turn the Amazon resolver off.

- default: `true`

### spec.enableDnsHostnames

`bool`

Whether instances with public IP addresses receive public DNS hostnames.
AWS leaves this off by default; turn it on for VPCs whose instances should
be reachable by DNS name. See:
https://docs.aws.amazon.com/vpc/latest/userguide/vpc-dns.html#vpc-dns-hostnames

- default: `true`

### spec.enableNetworkAddressUsageMetrics

`bool`

Whether to enable Network Address Usage (NAU) metrics for the VPC. NAU is a
CloudWatch metric that tracks how the VPC's IP address space is consumed,
useful for capacity planning in large networks. Off by default.

### spec.assignGeneratedIpv6CidrBlock

`bool`

Request an Amazon-provided IPv6 /56 CIDR block for the VPC. This is the
simplest way to make a VPC dual-stack. Mutually exclusive with the IPAM
IPv6 fields (ipv6_ipam_pool_id / ipv6_cidr_block / ipv6_netmask_length).

### spec.ipv6CidrBlock

`string`

An explicit IPv6 CIDR block to allocate from ipv6_ipam_pool_id (e.g.
"2600:1f18:abcd:1200::/56"). Requires ipv6_ipam_pool_id (a bring-your-own
IPv6 block without IPAM is not supported on the VPC resource) and is
mutually exclusive with ipv6_netmask_length.

### spec.ipv6CidrBlockNetworkBorderGroup

`string`

The network border group from which to advertise an Amazon-provided IPv6
CIDR (a Local Zone / Wavelength advertisement scope). Only valid together
with assign_generated_ipv6_cidr_block.

### spec.ipv6IpamPoolId

`string`

The IPAM pool to allocate the IPv6 CIDR from. Pair with ipv6_netmask_length
(the size to allocate) or ipv6_cidr_block (a specific block). Mutually
exclusive with assign_generated_ipv6_cidr_block.

### spec.ipv6NetmaskLength

`int32`

The netmask length of the IPv6 CIDR to allocate from ipv6_ipam_pool_id.
Must be one of 44, 48, 52, 56, or 60. Requires ipv6_ipam_pool_id and is
mutually exclusive with ipv6_cidr_block.

## Validation Rules

- `ipv4_primary_source_required`: set exactly one primary IPv4 source: cidr_block or ipv4_ipam_pool_id
- `ipv4_cidr_block_vs_netmask_exclusive`: cidr_block and ipv4_netmask_length are mutually exclusive
- `ipv4_cidr_block_format`: cidr_block must be a valid IPv4 CIDR, e.g. 10.0.0.0/16
- `ipv4_netmask_length_valid`: ipv4_netmask_length must be between 16 and 28 and requires ipv4_ipam_pool_id
- `ipv6_amazon_vs_ipam_exclusive`: assign_generated_ipv6_cidr_block cannot be combined with the IPAM IPv6 fields (ipv6_ipam_pool_id/ipv6_cidr_block/ipv6_netmask_length)
- `ipv6_cidr_block_requires_ipam`: ipv6_cidr_block requires ipv6_ipam_pool_id
- `ipv6_netmask_length_valid`: ipv6_netmask_length must be one of 44, 48, 52, 56, 60 and requires ipv6_ipam_pool_id
- `ipv6_cidr_block_vs_netmask_exclusive`: ipv6_cidr_block and ipv6_netmask_length are mutually exclusive
- `ipv6_border_group_requires_amazon`: ipv6_cidr_block_network_border_group requires assign_generated_ipv6_cidr_block

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsVpc, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpc_id` | `string` | The ID of the VPC (e.g. "vpc-0abc123"). The primary handle other resources reference via status.outputs.vpc_id. |
| `status.outputs.vpc_arn` | `string` | The ARN of the VPC. |
| `status.outputs.cidr_block` | `string` | The primary IPv4 CIDR block of the VPC (e.g. "10.0.0.0/16"). Reflects the allocated block even when it was assigned from IPAM. |
| `status.outputs.ipv6_cidr_block` | `string` | The IPv6 CIDR block associated with the VPC (e.g. "2600:1f18:abcd:1200::/56"), empty when the VPC is IPv4-only. |
| `status.outputs.owner_id` | `string` | The AWS account ID that owns the VPC. |
| `status.outputs.main_route_table_id` | `string` | The ID of the VPC's main route table. Subnets with no explicit route-table association use this table. |
| `status.outputs.default_security_group_id` | `string` | The ID of the default security group created with the VPC. |
| `status.outputs.default_network_acl_id` | `string` | The ID of the default network ACL created with the VPC. |
| `status.outputs.default_route_table_id` | `string` | The ID of the default route table created with the VPC (equal to the main route table). |
| `status.outputs.region` | `string` | The region the VPC was created in (mirrors spec.region, included for convenience). |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsClientVpn | `spec.vpcId` | `status.outputs.vpc_id` |
| AwsCodeBuildProject | `spec.vpcConfig.vpcId` | `status.outputs.vpc_id` |
| AwsEgressOnlyInternetGateway | `spec.vpcId` | `status.outputs.vpc_id` |
| AwsInternetGateway | `spec.vpcId` | `status.outputs.vpc_id` |
| AwsLbTargetGroup | `spec.vpcId` | `status.outputs.vpc_id` |
| AwsRoute53Zone | `spec.vpcAssociations[].vpcId` | `status.outputs.vpc_id` |
| AwsSagemakerDomain | `spec.vpcId` | `status.outputs.vpc_id` |
| AwsSecurityGroup | `spec.vpcId` | `status.outputs.vpc_id` |
| AwsSubnet | `spec.vpcId` | `status.outputs.vpc_id` |
| AwsTransitGatewayVpcAttachment | `spec.vpcId` | `status.outputs.vpc_id` |
| AwsVpcEndpoint | `spec.vpcId` | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
