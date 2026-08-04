# OpenStackRouter

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackRouterSpec defines the configuration for an OpenStack Neutron router.

A router provides L3 routing between subnets and, optionally, external network
connectivity via a gateway. Routers are the backbone of OpenStack networking --
they connect tenant subnets to each other and to the outside world via SNAT/DNAT.

The router name is derived from metadata.name.

Terraform resource: openstack_networking_router_v2
Pulumi resource: openstack.networking.Router

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackRouter
metadata:
  name: test-router
spec:
  external_network_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  enable_snat: true
  description: "Test router for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.externalNetworkId` | `string \| valueFrom` |  |  | OpenStackNetwork (`status.outputs.network_id`) |
| `spec.adminStateUp` | `bool` |  | `true` |  |
| `spec.enableSnat` | `bool` |  |  |  |
| `spec.distributed` | `bool` |  |  |  |
| `spec.externalFixedIps` | `[]ExternalFixedIp` |  |  |  |
| `spec.externalFixedIps[].subnetId` | `string` |  |  |  |
| `spec.externalFixedIps[].ipAddress` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.externalNetworkId

`string | valueFrom`

external_network_id is the ID of the external (provider) network used as
the router's gateway. When set, the router gains external connectivity and
can perform SNAT for tenant traffic.
Optional: routers without an external gateway provide internal routing only.
Can reference an OpenStackNetwork resource's output or be a literal network UUID.

- references: OpenStackNetwork (`status.outputs.network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackNetwork, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.adminStateUp

`bool` · optional (explicit presence)

admin_state_up controls the administrative state of the router.
When false, the router is administratively disabled and does not forward traffic.
Default: true

- default: `true`

### spec.enableSnat

`bool` · optional (explicit presence)

enable_snat controls whether Source NAT is enabled on the router's external gateway.
When enabled, traffic from tenant subnets is NATed to the router's external IP.
Only valid when external_network_id is configured.
If omitted, OpenStack uses its deployment default (typically true).

### spec.distributed

`bool` · optional (explicit presence)

distributed controls whether the router uses Distributed Virtual Router (DVR) mode.
DVR eliminates the centralized L3 agent bottleneck by distributing routing to
each compute node. This is a create-time setting and cannot be changed after creation.
If omitted, OpenStack uses its deployment default.

### spec.externalFixedIps

`[]ExternalFixedIp`

external_fixed_ips specifies fixed IP addresses to allocate on the external network
for the router's gateway. Each entry can request a specific subnet and/or IP address.
Only valid when external_network_id is configured.
If omitted, OpenStack automatically allocates an IP from the external network.

### spec.externalFixedIps[].subnetId

`string`

subnet_id is the UUID of a subnet on the external network from which to allocate the IP.
This references a subnet managed by the cloud administrator on the provider network.

### spec.externalFixedIps[].ipAddress

`string`

ip_address is the specific IP address to allocate on the external network.
Must be within the range of the specified subnet (or any external subnet if subnet_id is omitted).

### spec.description

`string`

description is a human-readable description of the router.
This is stored on the OpenStack resource and visible in Horizon and API responses.

### spec.tags

`[]string`

tags are string tags to associate with the router in OpenStack.
Tags are stored on the OpenStack resource and can be used for filtering
and organization in the OpenStack API and Horizon dashboard.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

region overrides the region from the provider config for this router.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Validation Rules

- `enable_snat.requires_external_network`: enable_snat can only be set when external_network_id is configured
- `external_fixed_ips.requires_external_network`: external_fixed_ips can only be specified when external_network_id is configured

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackRouter, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.router_id` | `string` | router_id is the unique identifier (UUID) of the router in OpenStack. This is the primary output used as a foreign key by downstream components. |
| `status.outputs.name` | `string` | name is the name of the router (derived from metadata.name). |
| `status.outputs.external_network_id` | `string` | external_network_id is the ID of the external network used as the router's gateway. Empty if no external gateway is configured. |
| `status.outputs.external_gateway_ip` | `string` | external_gateway_ip is the primary external IP address allocated to the router's gateway. This is the first IP from the external_fixed_ips list. Empty if no external gateway is configured. |
| `status.outputs.region` | `string` | region is the OpenStack region where the router was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.externalNetworkId` | OpenStackNetwork | `status.outputs.network_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackRouterInterface | `spec.routerId` | `status.outputs.router_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
