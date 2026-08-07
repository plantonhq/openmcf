# AwsNatGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsNatGatewaySpec defines a NAT gateway for an AWS Virtual Private Cloud (VPC).

A NAT (network address translation) gateway gives instances in a private
subnet outbound access without exposing them to inbound connections from the
internet. It is an AWS-managed, highly-available service that lives in a
single subnet and is referenced by other subnets' route tables as the target
of their default route.

There are two kinds, selected by connectivity_type:
  - public  (the common case): the gateway lives in a PUBLIC subnet (one that
    routes to an internet gateway) and is fronted by an Elastic IP. Private
    subnets send 0.0.0.0/0 to it to reach the internet outbound-only.
  - private: the gateway lives in any subnet and has no Elastic IP. It enables
    outbound communication to other VPCs or an on-premises network (via a
    transit/VPN/peering path) while keeping instances unreachable from those
    networks.

A NAT gateway does not, by itself, route anything: a private subnet becomes
NAT-egressed only when its route table sends a default route
(target_type = nat_gateway) to this gateway. The Elastic IP is composed by
reference (allocation_id -> AwsElasticIp), never embedded, so the address has
its own lifecycle and can be reasoned about as a first-class node.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsNatGateway
metadata:
  name: awsnatgateway-demo
spec:
  region: us-west-2
  connectivityType: public
  subnetId:
    value: "subnet-0abc1234def567890"
  allocationId:
    value: "eipalloc-0abc1234def567890"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.connectivityType` | `string` |  | `public` |  |
| `spec.subnetId` | `string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.allocationId` | `string \| valueFrom` |  |  | AwsElasticIp (`status.outputs.allocation_id`) |
| `spec.privateIp` | `string` |  |  |  |
| `spec.secondaryAllocationIds` | `[]string \| valueFrom` |  |  | AwsElasticIp (`status.outputs.allocation_id`) |
| `spec.secondaryPrivateIpAddresses` | `[]string` |  |  |  |
| `spec.secondaryPrivateIpAddressCount` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

AWS region the NAT gateway is created in. Must match the region of the
subnet it lives in. Example: "us-west-2", "eu-west-1". This drives provider
construction, so it is required even though the gateway logically inherits
the subnet's region.

- rule: {"string":{"minLen":"1"}}

### spec.connectivityType

`string`

How the gateway connects: "public" (default) or "private".

- "public": the gateway is placed in a public subnet and assigned an Elastic
  IP (allocation_id is required); it provides outbound internet access for
  private subnets that route to it.
- "private": the gateway has no Elastic IP and provides outbound access only
  to other private networks (peered/transit/VPN), never the internet.

AWS's own default for a new NAT gateway is public; this field is required so
the choice is always explicit (an empty value is rejected) and the two IaC
engines never need to invent a default. ForceNew: changing it replaces the
gateway.

- default: `public`
- rule: connectivity_type must be one of: public, private

### spec.subnetId

`string | valueFrom` · required

The subnet the NAT gateway is created in. Supply a literal subnet-id or
reference an AwsSubnet and the platform resolves its subnet_id output. For a
public gateway this must be a PUBLIC subnet (one whose route table reaches an
internet gateway); for a private gateway any subnet works. Immutable:
changing the subnet replaces the gateway.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.allocationId

`string | valueFrom`

The Elastic IP allocation that gives a public gateway its stable outbound
address. Supply a literal eipalloc-id or reference an AwsElasticIp and the
platform resolves its allocation_id output. Required when
connectivity_type is public; must be empty when it is private. ForceNew:
changing it replaces the gateway.

- references: AwsElasticIp (`status.outputs.allocation_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsElasticIp, name: <that resource's name>, fieldPath: status.outputs.allocation_id}} -- a bare string does not parse

### spec.privateIp

`string`

The private IPv4 address to assign to a private gateway from the subnet's
range. Only valid when connectivity_type is private; leave empty to let AWS
choose. ForceNew: changing it replaces the gateway.

### spec.secondaryAllocationIds

`[]string | valueFrom`

Additional Elastic IP allocations to attach to a public gateway, increasing
the number of available source ports for very high-throughput egress. Each
entry is a literal eipalloc-id or a reference to an AwsElasticIp. Only valid
when connectivity_type is public.

- references: AwsElasticIp (`status.outputs.allocation_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsElasticIp, name: <that resource's name>, fieldPath: status.outputs.allocation_id}} -- a bare string does not parse

### spec.secondaryPrivateIpAddresses

`[]string`

Additional private IPv4 addresses to assign to a private gateway, increasing
the number of available source ports. Only valid when connectivity_type is
private. Mutually exclusive with secondary_private_ip_address_count.

### spec.secondaryPrivateIpAddressCount

`int32`

Number of additional private IPv4 addresses to let AWS assign to a private
gateway (an alternative to listing them explicitly). Only valid when
connectivity_type is private. Mutually exclusive with
secondary_private_ip_addresses.

## Validation Rules

- `public_requires_allocation_id`: allocation_id is required when connectivity_type is public (a public NAT gateway needs an Elastic IP)
- `private_forbids_eip`: allocation_id and secondary_allocation_ids are not allowed when connectivity_type is private
- `public_forbids_private_addressing`: private_ip, secondary_private_ip_addresses, and secondary_private_ip_address_count are only valid when connectivity_type is private
- `secondary_private_ip_exclusive`: secondary_private_ip_addresses and secondary_private_ip_address_count are mutually exclusive

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsNatGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.nat_gateway_id` | `string` | The NAT gateway's id (e.g. "nat-0abc123"). This is the value a subnet route uses as its target_id when target_type is nat_gateway. |
| `status.outputs.public_ip` | `string` | The public IPv4 address of a public gateway (the Elastic IP's address). Empty for a private gateway. |
| `status.outputs.private_ip` | `string` | The private IPv4 address assigned to the gateway within its subnet. |
| `status.outputs.network_interface_id` | `string` | The id of the elastic network interface AWS created for the gateway. |
| `status.outputs.subnet_id` | `string` | The id of the subnet the gateway lives in. |
| `status.outputs.region` | `string` | The AWS region the NAT gateway was created in. Echoed so downstream tooling and verifiers can target the correct region. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.subnetId` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.allocationId` | AwsElasticIp | `status.outputs.allocation_id` |
| `spec.secondaryAllocationIds` | AwsElasticIp | `status.outputs.allocation_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
