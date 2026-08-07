# OpenStackLoadBalancerMember

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackLoadBalancerMemberSpec defines the configuration for an Octavia pool member
in OpenStack. A member represents a backend server that receives traffic from the pool's
load-balancing algorithm.

Terraform resource: openstack_lb_member_v2
Pulumi resource:    loadbalancer.Member

Validations:
- protocol_port must be between 1 and 65535.
- weight must be between 0 and 256 (when set).

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackLoadBalancerMember
metadata:
  name: test-member
spec:
  poolId:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  address: "10.0.0.10"
  protocolPort: 8080
  weight: 1
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.poolId` | `string \| valueFrom` | yes |  | OpenStackLoadBalancerPool (`status.outputs.pool_id`) |
| `spec.address` | `string` | yes |  |  |
| `spec.protocolPort` | `int32` | yes |  |  |
| `spec.subnetId` | `string \| valueFrom` |  |  | OpenStackSubnet (`status.outputs.subnet_id`) |
| `spec.weight` | `int32` |  |  |  |
| `spec.adminStateUp` | `bool` |  | `true` |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.poolId

`string | valueFrom` · required

(Required) The pool to add this member to.
ForceNew: changing this requires recreating the member.

FK: OpenStackLoadBalancerPool.status.outputs.pool_id

- references: OpenStackLoadBalancerPool (`status.outputs.pool_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackLoadBalancerPool, name: <that resource's name>, fieldPath: status.outputs.pool_id}} -- a bare string does not parse

### spec.address

`string` · required

(Required) The IP address of the backend server.
This is the actual IP that receives forwarded traffic from the load balancer.
ForceNew: changing this requires recreating the member.

This is a plain string (not StringValueOrRef) because member backends can be
any IP -- VMs, containers, bare metal, or external services.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.protocolPort

`int32` · required

(Required) The port on the backend server that accepts traffic.
ForceNew: changing this requires recreating the member.
Must be between 1 and 65535.

- rule: protocol_port must be between 1 and 65535
- rule: {"required":true}

### spec.subnetId

`string | valueFrom`

(Optional) The subnet where the member resides.
Used by Octavia for L3 routing when the member is on a different subnet
than the VIP. ForceNew: changing this requires recreating the member.

FK: OpenStackSubnet.status.outputs.subnet_id

- references: OpenStackSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.weight

`int32` · optional (explicit presence)

(Optional) Weight of this member for weighted load balancing.
A member with weight 0 receives no traffic (useful for draining).
Valid range: 0-256. Default: 1 (set by Octavia when unspecified).

Uses optional int32 because proto3 int32 defaults to 0, which means
"disabled member" in Octavia. optional lets us distinguish "not set"
(Octavia picks default 1) from "explicitly 0" (drain).

- rule: weight must be between 0 and 256

### spec.adminStateUp

`bool` · optional (explicit presence)

(Optional) Administrative state of the member.
When false, the member is removed from the pool's rotation. Default: true.

- default: `true`

### spec.tags

`[]string`

(Optional) Tags applied to the member in OpenStack.
Must be unique within this resource.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

(Optional) Override the region from the provider configuration.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackLoadBalancerMember, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.member_id` | `string` | member_id is the unique identifier (UUID) of the member. |
| `status.outputs.name` | `string` | name is the name of the member (derived from metadata.name). |
| `status.outputs.address` | `string` | address is the IP address of the backend server. |
| `status.outputs.protocol_port` | `int32` | protocol_port is the port on the backend server. |
| `status.outputs.weight` | `int32` | weight is the member's weight in the load-balancing algorithm. |
| `status.outputs.region` | `string` | region is the OpenStack region where the member was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.poolId` | OpenStackLoadBalancerPool | `status.outputs.pool_id` |
| `spec.subnetId` | OpenStackSubnet | `status.outputs.subnet_id` |

## See Also

- [Overview](../README.md)
