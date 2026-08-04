# AwsSubnet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsSubnetSpec defines a single subnet within an AWS Virtual Private Cloud (VPC).

A subnet is a contiguous range of IP addresses scoped to exactly one
availability zone. Whether a subnet is "public" or "private" is not a property
of the subnet itself -- it is determined by its route table: a subnet whose
route table has a default route to an internet gateway is public; one whose
default route points at a NAT gateway (or which has no internet route) is
private. Routing is therefore folded into this spec.

Two mutually exclusive ways to attach a route table are supported:
  - routes: inline rules from which a dedicated route table is created and
    owned by this subnet.
  - route_table_id: an existing, externally-managed route table to associate.
If neither is set, the subnet uses the VPC's main route table. The
route_table_id is always exported as a stack output so downstream resources
can reference whichever table ended up associated.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSubnet
metadata:
  name: awssubnet-demo
spec:
  region: us-west-2
  vpcId:
    value: "vpc-0abc1234def567890"
  availabilityZone: "us-west-2a"
  cidrBlock: "10.0.1.0/24"
  mapPublicIpOnLaunch: true
  routes:
    - destinationCidrBlock: "0.0.0.0/0"
      targetType: internet_gateway
      targetId:
        value: "igw-0abc1234def567890"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AwsVpc (`status.outputs.vpc_id`) |
| `spec.availabilityZone` | `string` | yes |  |  |
| `spec.cidrBlock` | `string` | yes |  |  |
| `spec.mapPublicIpOnLaunch` | `bool` |  |  |  |
| `spec.assignIpv6AddressOnCreation` | `bool` |  |  |  |
| `spec.ipv6CidrBlock` | `string` |  |  |  |
| `spec.enableDns64` | `bool` |  |  |  |
| `spec.enableResourceNameDnsARecordOnLaunch` | `bool` |  |  |  |
| `spec.enableResourceNameDnsAaaaRecordOnLaunch` | `bool` |  |  |  |
| `spec.privateDnsHostnameTypeOnLaunch` | `string` |  |  |  |
| `spec.routeTableId` | `string \| valueFrom` |  |  |  |
| `spec.routes` | `[]AwsSubnetRoute` |  |  |  |
| `spec.routes[].destinationCidrBlock` | `string` |  |  |  |
| `spec.routes[].destinationIpv6CidrBlock` | `string` |  |  |  |
| `spec.routes[].destinationPrefixListId` | `string` |  |  |  |
| `spec.routes[].targetType` | `enum` |  |  |  |
| `spec.routes[].targetId` | `string \| valueFrom` | yes |  |  |

## Field Details

### spec.region

`string` · required

AWS region the subnet is created in. Must match the region of the parent
VPC (a subnet cannot span regions). Example: "us-west-2", "eu-west-1".
This drives provider construction, so it is required even though the subnet
logically inherits the VPC's region.

- rule: {"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

The VPC this subnet belongs to. Supply a literal vpc-id or reference an
AwsVpc and the platform resolves its vpc_id output. Immutable: changing the
VPC replaces the subnet.

- references: AwsVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.availabilityZone

`string` · required

The availability zone the subnet lives in (e.g. "us-west-2a"). AWS subnets
are single-AZ; to span AZs, create one AwsSubnet per zone. Immutable:
changing the AZ replaces the subnet.

- rule: {"string":{"minLen":"1"}}

### spec.cidrBlock

`string` · required

IPv4 CIDR block for the subnet (e.g. "10.0.1.0/24"). Must fall within the
VPC's CIDR and not overlap any sibling subnet. Note that AWS reserves the
first four and the last IP address in every subnet, so a /28 yields 11
usable addresses. Immutable: changing the CIDR replaces the subnet.

- rule: {"required":true,"string":{"pattern":"^([0-9]{1,3}\\.){3}[0-9]{1,3}/[0-9]{1,2}$"}}

### spec.mapPublicIpOnLaunch

`bool`

When true, instances launched into this subnet receive a public IPv4
address by default. This is a convenience for public subnets; it does NOT
by itself make the subnet routable to the internet (that requires a route
to an internet gateway). Defaults to false (the AWS default).

### spec.assignIpv6AddressOnCreation

`bool`

When true, network interfaces created in this subnet are assigned an IPv6
address from the subnet's IPv6 CIDR on creation. Only meaningful when
ipv6_cidr_block is set. Defaults to false (the AWS default).

### spec.ipv6CidrBlock

`string`

IPv6 CIDR block for a dual-stack subnet (e.g. "2600:1f18:abcd:1200::/64").
Must be a /64 carved from an IPv6 CIDR associated with the parent VPC.
Leave empty for an IPv4-only subnet.

### spec.enableDns64

`bool`

When true, enables DNS64 on the subnet so that instances can reach
IPv4-only destinations from an IPv6-only subnet via NAT64. Requires the
subnet to have an IPv6 CIDR. Defaults to false (the AWS default).

### spec.enableResourceNameDnsARecordOnLaunch

`bool`

When true, instances launched into this subnet get a DNS A record for their
resource name. Used with private_dns_hostname_type_on_launch. Defaults to
false (the AWS default).

### spec.enableResourceNameDnsAaaaRecordOnLaunch

`bool`

When true, instances launched into this subnet get a DNS AAAA record for
their resource name (IPv6). Defaults to false (the AWS default).

### spec.privateDnsHostnameTypeOnLaunch

`string` · optional (explicit presence)

The type of hostname assigned to instances at launch. One of "ip-name"
(hostname derived from the private IPv4 address) or "resource-name"
(hostname derived from the instance's resource id, required for IPv6-only
subnets). Leave empty to use the AWS default for the subnet's address
family.

- rule: private_dns_hostname_type_on_launch must be one of: ip-name, resource-name

### spec.routeTableId

`string | valueFrom`

An existing route table to associate with this subnet. Mutually exclusive
with routes. If neither is set, the VPC's main route table is used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.routes

`[]AwsSubnetRoute`

Inline route rules. When present, a dedicated route table is created, owned
by this subnet, populated with these rules, and associated with the subnet.
Mutually exclusive with route_table_id.

- rule: exactly one of destination_cidr_block, destination_ipv6_cidr_block, or destination_prefix_list_id must be set

### spec.routes[].destinationCidrBlock

`string`

Destination IPv4 CIDR (e.g. "0.0.0.0/0" for the default route).

### spec.routes[].destinationIpv6CidrBlock

`string`

Destination IPv6 CIDR (e.g. "::/0" for the default IPv6 route).

### spec.routes[].destinationPrefixListId

`string`

Destination managed-prefix-list id (e.g. "pl-0123abcd"), for routing to a
curated set of CIDRs such as an AWS service prefix list.

### spec.routes[].targetType

`enum`

The kind of network entity this route targets. Determines which AWS route
attribute target_id maps to (gateway_id, nat_gateway_id, etc.).

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `unspecified`
- `internet_gateway`
- `nat_gateway`
- `transit_gateway`
- `vpc_peering_connection`
- `vpc_endpoint`
- `network_interface`
- `egress_only_internet_gateway`

### spec.routes[].targetId

`string | valueFrom` · required

Identifier of the target. Supply a literal id, or reference the resource
that produces it (e.g. an AwsInternetGateway's internet_gateway_id or an
AwsNatGateway's nat_gateway_id). This single field is intentionally
polymorphic across all of target_type's kinds, so it carries no
default_kind -- the target's kind is given by target_type, and the
producing resource is referenced explicitly via value_from.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

## Validation Rules

- `route_table_mutual_exclusivity`: route_table_id and routes are mutually exclusive; provide one or neither

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSubnet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.subnetId` | `string` | The subnet's id (e.g. "subnet-0abc123"). |
| `status.outputs.subnetArn` | `string` | The subnet's ARN. |
| `status.outputs.availabilityZone` | `string` | The availability zone the subnet resides in (e.g. "us-west-2a"). |
| `status.outputs.cidrBlock` | `string` | The subnet's IPv4 CIDR block. |
| `status.outputs.routeTableId` | `string` | The id of the route table this subnet owns or references: the inline table created from routes, or the externally referenced route_table_id. EMPTY when neither is set -- the subnet then rides the VPC main route table, which is deliberately not echoed here (attaching to it is an explicit choice); use the AwsVpc's main_route_table_id / default_route_table_id outputs instead. |
| `status.outputs.region` | `string` | The AWS region the subnet was created in. Echoed so downstream tooling and verifiers can target the correct region. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AwsVpc | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
