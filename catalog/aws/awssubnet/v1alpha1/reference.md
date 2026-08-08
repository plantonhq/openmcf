# AwsSubnet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

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
apiVersion: aws.planton.dev/v1alpha1
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
| `spec.availabilityZone` | `string` |  |  |  |
| `spec.cidrBlock` | `string` |  |  |  |
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
| `spec.availabilityZoneId` | `string` |  |  |  |
| `spec.ipv4IpamPoolId` | `string \| valueFrom` |  |  |  |
| `spec.ipv4NetmaskLength` | `int32` |  |  |  |
| `spec.ipv6IpamPoolId` | `string \| valueFrom` |  |  |  |
| `spec.ipv6NetmaskLength` | `int32` |  |  |  |
| `spec.ipv6Native` | `bool` |  |  |  |
| `spec.propagatingVgws` | `[]string` |  |  |  |

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

`string`

The availability zone the subnet lives in (e.g. "us-west-2a"). AWS subnets
are single-AZ; to span AZs, create one AwsSubnet per zone. Immutable:
changing the AZ replaces the subnet.

### spec.cidrBlock

`string`

IPv4 CIDR block for the subnet (e.g. "10.0.1.0/24"). Must fall within the
VPC's CIDR and not overlap any sibling subnet. Note that AWS reserves the
first four and the last IP address in every subnet, so a /28 yields 11
usable addresses. Immutable: changing the CIDR replaces the subnet.

- rule: cidr_block must be an IPv4 CIDR like 10.0.1.0/24

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
- `carrier_gateway`
- `core_network`
- `local_gateway`
- `odb_network`

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

### spec.availabilityZoneId

`string`

### spec.ipv4IpamPoolId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.ipv4NetmaskLength

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":28,"gte":16}}

### spec.ipv6IpamPoolId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.ipv6NetmaskLength

`int32` · optional (explicit presence)

- rule: {"int32":{"in":[44,48,52,56,60,64]}}

### spec.ipv6Native

`bool`

### spec.propagatingVgws

`[]string`

## Validation Rules

- `route_table_mutual_exclusivity`: route_table_id and routes are mutually exclusive; provide one or neither
- `availability_zone_exactly_one`: exactly one of availability_zone or availability_zone_id must be set
- `ipv4_addressing_exactly_one`: set exactly one of cidr_block or ipv4_ipam_pool_id (neither when ipv6_native is true)
- `ipv4_netmask_requires_pool`: ipv4_netmask_length requires ipv4_ipam_pool_id
- `ipv6_netmask_requires_pool`: ipv6_netmask_length requires ipv6_ipam_pool_id
- `ipv6_addressing_at_most_one`: ipv6_cidr_block and ipv6_ipam_pool_id are mutually exclusive
- `ipv6_native_requires_ipv6`: ipv6_native requires ipv6_cidr_block or ipv6_ipam_pool_id
- `propagating_vgws_owned_table_only`: propagating_vgws is not allowed with route_table_id (propagation belongs to the table's owner)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSubnet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.subnet_id` | `string` | The subnet's id (e.g. "subnet-0abc123"). |
| `status.outputs.subnet_arn` | `string` | The subnet's ARN. |
| `status.outputs.availability_zone` | `string` | The availability zone the subnet resides in (e.g. "us-west-2a"). |
| `status.outputs.cidr_block` | `string` | The subnet's IPv4 CIDR block. |
| `status.outputs.route_table_id` | `string` | The id of the route table this subnet owns or references: the inline table created from routes, or the externally referenced route_table_id. EMPTY when neither is set -- the subnet then rides the VPC main route table, which is deliberately not echoed here (attaching to it is an explicit choice); use the AwsVpc's main_route_table_id / default_route_table_id outputs instead. |
| `status.outputs.region` | `string` | The AWS region the subnet was created in. Echoed so downstream tooling and verifiers can target the correct region. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AwsVpc | `status.outputs.vpc_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsAlb | `spec.subnets` | `status.outputs.subnet_id` |
| AwsAppRunnerVpcConnector | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsAutoScalingGroup | `spec.subnets` | `status.outputs.subnet_id` |
| AwsBatchComputeEnvironment | `spec.computeResources.subnetIds` | `status.outputs.subnet_id` |
| AwsClientVpn | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsClientVpn | `spec.routes[].targetSubnetId` | `status.outputs.subnet_id` |
| AwsCodeBuildProject | `spec.vpcConfig.subnetIds` | `status.outputs.subnet_id` |
| AwsDocumentDb | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsEc2Instance | `spec.subnetId` | `status.outputs.subnet_id` |
| AwsEc2Instance | `spec.secondaryNetworkInterfaces[].subnetId` | `status.outputs.subnet_id` |
| AwsEcsService | `spec.network.subnets` | `status.outputs.subnet_id` |
| AwsEksCluster | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsEksFargateProfile | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsEksNodeGroup | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsElasticFileSystem | `spec.mountTargets[].subnetId` | `status.outputs.subnet_id` |
| AwsEventBridgeRule | `spec.targets[].ecsTarget.networkConfiguration.subnets` | `status.outputs.subnet_id` |
| AwsFsxLustreFileSystem | `spec.subnetId` | `status.outputs.subnet_id` |
| AwsFsxOntapFileSystem | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsFsxOntapFileSystem | `spec.preferredSubnetId` | `status.outputs.subnet_id` |
| AwsFsxOntapFileSystem | `spec.routeTableIds` | `status.outputs.route_table_id` |
| AwsFsxOpenzfsFileSystem | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsFsxOpenzfsFileSystem | `spec.preferredSubnetId` | `status.outputs.subnet_id` |
| AwsFsxOpenzfsFileSystem | `spec.routeTableIds` | `status.outputs.route_table_id` |
| AwsFsxWindowsFileSystem | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsFsxWindowsFileSystem | `spec.preferredSubnetId` | `status.outputs.subnet_id` |
| AwsHttpApiVpcLink | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsKinesisFirehose | `spec.opensearch.vpcConfig.subnetIds` | `status.outputs.subnet_id` |
| AwsKinesisFirehose | `spec.opensearchServerless.vpcConfig.subnetIds` | `status.outputs.subnet_id` |
| AwsLambda | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsLaunchTemplate | `spec.networkInterfaces[].subnetId` | `status.outputs.subnet_id` |
| AwsMemcachedElasticache | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsMemorydbCluster | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsMskCluster | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsMskServerlessCluster | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsMwaaEnvironment | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsNatGateway | `spec.subnetId` | `status.outputs.subnet_id` |
| AwsNeptuneCluster | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsNlb | `spec.subnetMappings[].subnetId` | `status.outputs.subnet_id` |
| AwsOpenSearchDomain | `spec.vpcOptions.subnetIds` | `status.outputs.subnet_id` |
| AwsPlantonRunner | `spec.subnets` | `status.outputs.subnet_id` |
| AwsRdsCluster | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsRdsInstance | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsRedisElasticache | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsRedshiftCluster | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsRedshiftServerlessWorkgroup | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsSagemakerDomain | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsServerlessElasticache | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsTransitGatewayVpcAttachment | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsVpcEndpoint | `spec.routeTableIds` | `status.outputs.route_table_id` |
| AwsVpcEndpoint | `spec.subnetIds` | `status.outputs.subnet_id` |
| AwsVpcEndpoint | `spec.subnetConfigurations[].subnetId` | `status.outputs.subnet_id` |

## See Also

- [Overview](../README.md)
