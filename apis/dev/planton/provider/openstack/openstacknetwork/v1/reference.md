# OpenStackNetwork

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackNetworkSpec defines the configuration for an OpenStack Neutron network.

A network is the foundational networking primitive in OpenStack. It represents an
isolated Layer 2 broadcast domain. Subnets, ports, routers, and security groups
all attach to or reference a network.

The network name is derived from metadata.name.

Terraform resource: openstack_networking_network_v2
Pulumi resource: openstack.networking.Network

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackNetwork
metadata:
  name: test-network
spec:
  description: "Test network for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.description` | `string` |  |  |  |
| `spec.adminStateUp` | `bool` |  | `true` |  |
| `spec.shared` | `bool` |  |  |  |
| `spec.external` | `bool` |  |  |  |
| `spec.mtu` | `int32` |  |  |  |
| `spec.dnsDomain` | `string` |  |  |  |
| `spec.portSecurityEnabled` | `bool` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.description

`string`

description is a human-readable description of the network.
This is stored on the OpenStack resource and visible in Horizon and API responses.

### spec.adminStateUp

`bool` · optional (explicit presence)

admin_state_up controls the administrative state of the network.
When false, the network is administratively down and will not forward traffic.
Default: true

- default: `true`

### spec.shared

`bool`

shared indicates whether this network is shared across all tenants/projects.
Creating shared networks typically requires admin privileges.
Most tenant users should leave this as false (the default).

### spec.external

`bool`

external indicates whether this is an external (provider) network.
External networks are used for floating IP allocation and router gateways.
Creating external networks typically requires admin privileges.

### spec.mtu

`int32`

mtu is the Maximum Transmission Unit value for the network, in bytes.
Common values: 1500 (standard Ethernet), 1450 (VXLAN overlay), 9000 (jumbo frames).
If omitted or set to 0, OpenStack uses the default MTU for the network type.

- rule: {"int32":{"gte":0}}

### spec.dnsDomain

`string`

dns_domain is the DNS domain associated with the network.
When set, ports on this network can have DNS names auto-assigned.
Must end with a dot (.) if specified (e.g., "my-network.example.com.").
Requires the dns-integration extension in Neutron.

- rule: {"string":{"pattern":"^$|\\.$"}}

### spec.portSecurityEnabled

`bool` · optional (explicit presence)

port_security_enabled controls whether port security is enforced on ports
created on this network. When enabled, only traffic matching the port's
security groups and allowed address pairs is permitted.
If omitted, the OpenStack deployment's default is used (typically true).

### spec.tags

`[]string`

tags are string tags to associate with the network in OpenStack.
Tags are stored on the OpenStack resource and can be used for filtering
and organization in the OpenStack API and Horizon dashboard.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

region overrides the region from the provider config for this network.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackNetwork, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.network_id` | `string` | network_id is the unique identifier (UUID) of the network in OpenStack. This is the primary output used as a foreign key by downstream components. |
| `status.outputs.name` | `string` | name is the name of the network (derived from metadata.name). |
| `status.outputs.region` | `string` | region is the OpenStack region where the network was created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackFloatingIp | `spec.floatingNetworkId` | `status.outputs.network_id` |
| OpenStackInstance | `spec.networks[].uuid` | `status.outputs.network_id` |
| OpenStackNetworkPort | `spec.networkId` | `status.outputs.network_id` |
| OpenStackRouter | `spec.externalNetworkId` | `status.outputs.network_id` |
| OpenStackSubnet | `spec.networkId` | `status.outputs.network_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
