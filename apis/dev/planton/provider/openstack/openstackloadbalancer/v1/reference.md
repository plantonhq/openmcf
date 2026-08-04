# OpenStackLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackLoadBalancerSpec defines the configuration for an Octavia load balancer
in OpenStack. This provisions the VIP (Virtual IP) endpoint on a specified subnet,
which listeners and pools attach to for traffic distribution.

Terraform resource: openstack_lb_loadbalancer_v2
Pulumi resource:    loadbalancer.LoadBalancer

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackLoadBalancer
metadata:
  name: test-lb
spec:
  vip_subnet_id:
    value: "e0a1f622-9aab-4a48-8c8c-3b0c7e2a9b1d"
  description: "Test load balancer for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.vipSubnetId` | `string \| valueFrom` | yes |  | OpenStackSubnet (`status.outputs.subnet_id`) |
| `spec.vipAddress` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.adminStateUp` | `bool` |  | `true` |  |
| `spec.flavorId` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.vipSubnetId

`string | valueFrom` · required

(Required) The subnet on which to allocate the VIP address.
This determines the network segment where the load balancer's virtual IP lives.
The subnet must already exist and have available IP addresses.

FK: OpenStackSubnet.status.outputs.subnet_id

- references: OpenStackSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.vipAddress

`string`

(Optional) A specific IP address to request for the VIP.
Must be within the CIDR range of the specified subnet.
If omitted, Octavia auto-allocates an available IP from the subnet.
ForceNew: changing this requires recreating the load balancer.

### spec.description

`string`

(Optional) Human-readable description of the load balancer.

### spec.adminStateUp

`bool` · optional (explicit presence)

(Optional) Administrative state of the load balancer.
When false, the LB stops accepting traffic. Default: true.

- default: `true`

### spec.flavorId

`string`

(Optional) The ID of an Octavia flavor to use for the load balancer.
Flavors define resource limits (bandwidth, connections, etc.).
ForceNew: changing this requires recreating the load balancer.

### spec.tags

`[]string`

(Optional) Tags applied to the load balancer in OpenStack.
Must be unique within this resource.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

(Optional) Override the region from the provider configuration.
Use this when the load balancer must be created in a specific region
that differs from the provider's default.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.loadbalancer_id` | `string` | loadbalancer_id is the unique identifier (UUID) of the load balancer. This is the primary output used as a foreign key by listeners. |
| `status.outputs.name` | `string` | name is the name of the load balancer (derived from metadata.name). |
| `status.outputs.vip_address` | `string` | vip_address is the Virtual IP address allocated to the load balancer. This is the IP that clients connect to. |
| `status.outputs.vip_port_id` | `string` | vip_port_id is the Neutron port ID of the VIP. Useful for attaching security groups or floating IPs to the VIP. |
| `status.outputs.vip_subnet_id` | `string` | vip_subnet_id is the subnet where the VIP was allocated. |
| `status.outputs.region` | `string` | region is the OpenStack region where the load balancer was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vipSubnetId` | OpenStackSubnet | `status.outputs.subnet_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackLoadBalancerListener | `spec.loadbalancerId` | `status.outputs.loadbalancer_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
