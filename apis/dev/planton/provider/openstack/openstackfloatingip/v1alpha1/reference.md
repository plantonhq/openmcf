# OpenStackFloatingIp

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackFloatingIpSpec defines the configuration for an OpenStack Neutron
floating IP allocation.

A floating IP provides external connectivity to instances or ports on tenant
networks. It is allocated from an external (provider) network and can
optionally be associated with a port for immediate connectivity. For
DAG-visible association in InfraCharts, allocate without port_id and use the
separate OpenStackFloatingIpAssociate component.

The floating IP name is derived from metadata.name (for Planton identity only;
floating IPs in OpenStack do not have a name attribute).

Terraform resource: openstack_networking_floatingip_v2
Pulumi resource: openstack.networking.FloatingIp

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackFloatingIp
metadata:
  name: test-fip
spec:
  floating_network_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  description: "Test floating IP for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.floatingNetworkId` | `string \| valueFrom` | yes |  | OpenStackNetwork (`status.outputs.network_id`) |
| `spec.portId` | `string \| valueFrom` |  |  | OpenStackNetworkPort (`status.outputs.port_id`) |
| `spec.fixedIp` | `string` |  |  |  |
| `spec.subnetId` | `string` |  |  |  |
| `spec.address` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.floatingNetworkId

`string | valueFrom` · required

floating_network_id is the ID of the external (provider) network from which
the floating IP is allocated. This is the pool of public IP addresses.
Maps to the Terraform "pool" attribute.
Can reference an OpenStackNetwork resource's output or be a literal network UUID.

- references: OpenStackNetwork (`status.outputs.network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackNetwork, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.portId

`string | valueFrom`

port_id is the ID of an existing port to associate with this floating IP.
When set, the floating IP is immediately bound to the port, providing
external connectivity to whatever is attached to that port (typically an instance).
Optional: omit for allocation-only (use OpenStackFloatingIpAssociate separately).
Can reference an OpenStackNetworkPort resource's output or be a literal port UUID.

- references: OpenStackNetworkPort (`status.outputs.port_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackNetworkPort, name: <that resource's name>, fieldPath: status.outputs.port_id}} -- a bare string does not parse

### spec.fixedIp

`string`

fixed_ip specifies which fixed IP address on the port to associate the
floating IP with. Only relevant when port_id is set and the port has
multiple IP addresses. If the port has a single IP, this can be omitted.

### spec.subnetId

`string`

subnet_id is the UUID of a subnet within the external network from which
to allocate the floating IP. This references an admin-managed subnet on the
provider network. If omitted, OpenStack allocates from any available subnet.

### spec.address

`string`

address requests a specific floating IP address from the pool.
If omitted, OpenStack allocates any available address.
This is a create-time setting (ForceNew in Terraform).
Use case: DNS pre-configuration, firewall whitelisting, IP reservation.

### spec.description

`string`

description is a human-readable description of the floating IP.
Stored on the OpenStack resource and visible in Horizon and API responses.

### spec.tags

`[]string`

tags are string tags to associate with the floating IP in OpenStack.
Tags are stored on the OpenStack resource and can be used for filtering
and organization in the OpenStack API and Horizon dashboard.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

region overrides the region from the provider config for this floating IP.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Validation Rules

- `fixed_ip.requires_port_id`: fixed_ip can only be set when port_id is configured

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackFloatingIp, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.floating_ip_id` | `string` | floating_ip_id is the unique identifier (UUID) of the floating IP resource in OpenStack. This is the Terraform resource ID. |
| `status.outputs.address` | `string` | address is the allocated floating IP address (e.g., "203.0.113.42"). This is the primary output -- consumed by OpenStackFloatingIpAssociate as a required FK target and by users for DNS configuration, firewall rules, etc. |
| `status.outputs.floating_network_id` | `string` | floating_network_id is the UUID of the external network the floating IP was allocated from. |
| `status.outputs.port_id` | `string` | port_id is the UUID of the port this floating IP is associated with. Empty if the floating IP is allocated but not associated (allocation-only mode). |
| `status.outputs.fixed_ip` | `string` | fixed_ip is the fixed IP address on the port that the floating IP is mapped to. Empty if no port association exists. |
| `status.outputs.region` | `string` | region is the OpenStack region where the floating IP was allocated. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.floatingNetworkId` | OpenStackNetwork | `status.outputs.network_id` |
| `spec.portId` | OpenStackNetworkPort | `status.outputs.port_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackFloatingIpAssociate | `spec.floatingIp` | `status.outputs.address` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
