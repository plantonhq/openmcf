# OpenStackFloatingIpAssociate

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackFloatingIpAssociateSpec defines the configuration for associating an
existing OpenStack Neutron floating IP with a port.

This is a "join" resource -- it binds a floating IP (allocated via
OpenStackFloatingIp) to a port (created via OpenStackNetworkPort). The
association provides external connectivity to whatever is attached to the
port (typically an instance).

Use this component when the floating IP allocation and port association need
to be separate DAG nodes in an InfraChart. For simple cases where allocation
and association happen together, use OpenStackFloatingIp with its built-in
port_id field instead.

All fields on the underlying resource are ForceNew except port_id: changing
the floating IP or region recreates the association.

The resource name in Planton is derived from metadata.name. OpenStack floating
IP associations do not have a user-visible "name" attribute.

Terraform resource: openstack_networking_floatingip_associate_v2
Pulumi resource: openstack.networking.FloatingIpAssociate

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackFloatingIpAssociate
metadata:
  name: test-fipa
spec:
  floating_ip:
    value: "203.0.113.42"
  port_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.floatingIp` | `string \| valueFrom` | yes |  | OpenStackFloatingIp (`status.outputs.address`) |
| `spec.portId` | `string \| valueFrom` | yes |  | OpenStackNetworkPort (`status.outputs.port_id`) |
| `spec.fixedIp` | `string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.floatingIp

`string | valueFrom` · required

floating_ip is the floating IP address (or ID) to associate.
This targets the *address* output of an OpenStackFloatingIp resource
(e.g., "203.0.113.42"), not its UUID. The Terraform provider accepts
either an IP address or a floating IP UUID for this field.
Can reference an OpenStackFloatingIp resource's output or be a literal address.

- references: OpenStackFloatingIp (`status.outputs.address`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackFloatingIp, name: <that resource's name>, fieldPath: status.outputs.address}} -- a bare string does not parse

### spec.portId

`string | valueFrom` · required

port_id is the ID of the port to associate the floating IP with.
The port must have at least one fixed IP address.
Can reference an OpenStackNetworkPort resource's output or be a literal port UUID.

- references: OpenStackNetworkPort (`status.outputs.port_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackNetworkPort, name: <that resource's name>, fieldPath: status.outputs.port_id}} -- a bare string does not parse

### spec.fixedIp

`string`

fixed_ip specifies which fixed IP address on the port to map the floating
IP to. Only relevant when the port has multiple fixed IP addresses.
If omitted and the port has a single IP, that IP is used automatically.
If omitted and the port has multiple IPs, OpenStack picks the first one.

### spec.region

`string`

region overrides the region from the provider config for this association.
If omitted, the region from the OpenStack provider config is used.
ForceNew: changing the region recreates the association.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackFloatingIpAssociate, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.id` | `string` | id is the Terraform resource identifier for the association. Format is typically the floating IP address itself. |
| `status.outputs.floating_ip` | `string` | floating_ip is the floating IP address that was associated. |
| `status.outputs.port_id` | `string` | port_id is the UUID of the port the floating IP was associated with. |
| `status.outputs.fixed_ip` | `string` | fixed_ip is the fixed IP address on the port that the floating IP maps to. Computed by OpenStack if not explicitly specified in the spec. |
| `status.outputs.region` | `string` | region is the OpenStack region where the association was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.floatingIp` | OpenStackFloatingIp | `status.outputs.address` |
| `spec.portId` | OpenStackNetworkPort | `status.outputs.port_id` |

## See Also

- [Overview](../README.md)
