# OpenStackRouterInterface

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackRouterInterfaceSpec defines the configuration for attaching an
OpenStack Neutron router to a subnet.

A router interface is a "join" resource -- it connects a router (L3) to a
subnet (L2) by creating a port on the subnet and attaching it to the router.
Without a router interface, a subnet has no route to other subnets or to
external networks.

All fields on the underlying resource are ForceNew: any change recreates
the interface. This is expected for declarative infrastructure.

The resource name in Planton is derived from metadata.name. OpenStack router
interfaces do not have a user-visible "name" attribute -- the resource is
identified by the port it creates.

Terraform resource: openstack_networking_router_interface_v2
Pulumi resource: openstack.networking.RouterInterface

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackRouterInterface
metadata:
  name: test-router-interface
spec:
  router_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  subnet_id:
    value: "b2c3d4e5-f6a7-8901-bcde-f12345678901"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.routerId` | `string \| valueFrom` | yes |  | OpenStackRouter (`status.outputs.router_id`) |
| `spec.subnetId` | `string \| valueFrom` | yes |  | OpenStackSubnet (`status.outputs.subnet_id`) |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.routerId

`string | valueFrom` · required

router_id is the ID of the router to attach the subnet to.
This is a required foreign key -- every router interface must reference
exactly one router.
Can reference an OpenStackRouter resource's output or be a literal router UUID.

- references: OpenStackRouter (`status.outputs.router_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackRouter, name: <that resource's name>, fieldPath: status.outputs.router_id}} -- a bare string does not parse

### spec.subnetId

`string | valueFrom` · required

subnet_id is the ID of the subnet to connect to the router.
This is a required foreign key -- every router interface must reference
exactly one subnet.
Can reference an OpenStackSubnet resource's output or be a literal subnet UUID.

- references: OpenStackSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.region

`string`

region overrides the region from the provider config for this router interface.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackRouterInterface, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.port_id` | `string` | port_id is the UUID of the port created by the router interface attachment. This is also the Terraform resource ID for the router interface. |
| `status.outputs.router_id` | `string` | router_id is the UUID of the router this interface is attached to. |
| `status.outputs.subnet_id` | `string` | subnet_id is the UUID of the subnet connected to the router. |
| `status.outputs.region` | `string` | region is the OpenStack region where the router interface was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.routerId` | OpenStackRouter | `status.outputs.router_id` |
| `spec.subnetId` | OpenStackSubnet | `status.outputs.subnet_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
