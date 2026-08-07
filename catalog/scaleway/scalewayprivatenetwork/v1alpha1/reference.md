# ScalewayPrivateNetwork

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayPrivateNetworkSpec defines the specification for a Scaleway Private Network.

A Scaleway Private Network is a regional, Layer 2 network that lives inside a VPC
and provides secure, private connectivity between Scaleway resources. It is the
primary attachment point for Kapsule clusters, RDB instances, Redis clusters,
Load Balancers, Instances, and other resources that need private communication.

Private Networks include built-in DHCP, so resources attached to a Private Network
receive IP addresses automatically. You may optionally specify an IPv4 CIDR to
control the address range; if omitted, Scaleway's IPAM allocates one for you.

In the Scaleway resource graph, Private Network is the "universal connector":
almost every infrastructure resource attaches to a Private Network rather than
directly to a VPC.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.vpcId` | `string \| valueFrom` | yes |  | ScalewayVpc (`status.outputs.vpc_id`) |
| `spec.region` | `string` | yes |  |  |
| `spec.ipv4Subnet` | `string` |  |  |  |
| `spec.ipv6Subnets` | `[]string` |  |  |  |
| `spec.enableDefaultRoutePropagation` | `bool` |  | `false` |  |

## Field Details

### spec.vpcId

`string | valueFrom` · required

The VPC in which to create this Private Network.

Can be a literal VPC UUID or a reference to a ScalewayVpc resource's output.
When used inside an infra chart, this is typically wired via valueFrom:

  vpcId:
    valueFrom:
      kind: ScalewayVpc
      name: my-vpc
      fieldPath: status.outputs.vpc_id

The Private Network's region must match the VPC's region.

- references: ScalewayVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.region

`string` · required

The Scaleway region where the Private Network will be created.
Examples: "fr-par", "nl-ams", "pl-waw"

Must match the region of the parent VPC. This field is required and
cannot be changed after creation.

- rule: {"required":true}

### spec.ipv4Subnet

`string`

IPv4 subnet in CIDR notation to associate with this Private Network.
Examples: "192.168.0.0/24", "10.0.1.0/24", "172.16.0.0/22"

If omitted, Scaleway's IPAM automatically allocates a subnet from a
default range. Specifying a subnet gives you control over the address
space, which is important when multiple Private Networks in the same
VPC need non-overlapping ranges for routing to work correctly.

The allocated (or auto-assigned) CIDR is always available in stack
outputs as ipv4_subnet_cidr.

### spec.ipv6Subnets

`[]string`

IPv6 subnets in CIDR notation to associate with this Private Network.
Examples: "fd46:78ab:30b8:177c::/64"

Optional. Multiple IPv6 subnets can be attached for dual-stack networking.
Most production deployments use IPv4 only; IPv6 is primarily useful for
workloads that must be reachable on IPv6 or for advanced networking scenarios.

### spec.enableDefaultRoutePropagation

`bool`

Whether to propagate default IPv4 and IPv6 routes for this Private Network.

When enabled, resources in this Private Network receive the VPC's default
routes, enabling them to communicate with resources in other Private Networks
within the same VPC (provided the VPC has routing enabled).

Default: false

- default: `false`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayPrivateNetwork, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.private_network_id` | `string` | The unique identifier (UUID) of the created Scaleway Private Network. This is the primary output referenced by downstream resources via StringValueOrRef. In infra charts, downstream resources wire to this value using: privateNetworkId: valueFrom: kind: ScalewayPrivateNetwork name: my-network fieldPath: status.outputs.private_network_id |
| `status.outputs.ipv4_subnet_cidr` | `string` | The IPv4 subnet CIDR associated with this Private Network. If ipv4_subnet was specified in the spec, this reflects the requested CIDR. If ipv4_subnet was omitted, this contains the CIDR auto-allocated by Scaleway's IPAM service. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | ScalewayVpc | `status.outputs.vpc_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| ScalewayInstance | `spec.privateNetworkId` | `status.outputs.private_network_id` |
| ScalewayKapsuleCluster | `spec.privateNetworkId` | `status.outputs.private_network_id` |
| ScalewayLoadBalancer | `spec.privateNetworkId` | `status.outputs.private_network_id` |
| ScalewayMongodbInstance | `spec.privateNetworkId` | `status.outputs.private_network_id` |
| ScalewayPublicGateway | `spec.privateNetworkId` | `status.outputs.private_network_id` |
| ScalewayRdbInstance | `spec.privateNetworkId` | `status.outputs.private_network_id` |
| ScalewayRedisCluster | `spec.privateNetworkId` | `status.outputs.private_network_id` |
| ScalewayServerlessContainer | `spec.privateNetworkId` | `status.outputs.private_network_id` |
| ScalewayServerlessFunction | `spec.privateNetworkId` | `status.outputs.private_network_id` |

## See Also

- [Overview](./README.md)
