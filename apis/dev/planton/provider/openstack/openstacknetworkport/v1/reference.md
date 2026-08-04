# OpenStackNetworkPort

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackNetworkPortSpec defines the configuration for an OpenStack Neutron port.

A port represents a virtual switch port on a network. It provides stable
network identity (MAC address, fixed IPs, security groups) that can be
attached to instances, load balancers, or other network-consuming resources.
Explicit port creation is preferred over instance-inline networking when you
need stable IP addresses, multiple security groups, or pre-provisioned
network identities.

The port name is derived from metadata.name.

Terraform resource: openstack_networking_port_v2
Pulumi resource: openstack.networking.Port

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackNetworkPort
metadata:
  name: test-port
spec:
  network_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  fixed_ips:
    - subnet_id:
        value: "b2c3d4e5-f6a7-8901-bcde-f12345678901"
      ip_address: "192.168.1.10"
  description: "Test port for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.networkId` | `string \| valueFrom` | yes |  | OpenStackNetwork (`status.outputs.network_id`) |
| `spec.fixedIps` | `[]FixedIp` |  |  |  |
| `spec.fixedIps[].subnetId` | `string \| valueFrom` |  |  | OpenStackSubnet (`status.outputs.subnet_id`) |
| `spec.fixedIps[].ipAddress` | `string` |  |  |  |
| `spec.securityGroupIds` | `[]string \| valueFrom` |  |  | OpenStackSecurityGroup (`status.outputs.security_group_id`) |
| `spec.noSecurityGroups` | `bool` |  |  |  |
| `spec.adminStateUp` | `bool` |  | `true` |  |
| `spec.macAddress` | `string` |  |  |  |
| `spec.portSecurityEnabled` | `bool` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.networkId

`string | valueFrom` · required

network_id is the ID of the network to create this port on.
This is the defining relationship -- every port belongs to exactly one network.
ForceNew: changing the network recreates the port.
Can reference an OpenStackNetwork resource's output or be a literal network UUID.

- references: OpenStackNetwork (`status.outputs.network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackNetwork, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.fixedIps

`[]FixedIp`

fixed_ips defines the IP address allocations for this port.
Each entry assigns an IP from a subnet on the port's network.
If omitted, OpenStack auto-assigns one IP from any subnet on the network.
Multiple entries create a multi-homed port (one IP per subnet).

### spec.fixedIps[].subnetId

`string | valueFrom`

subnet_id is the ID of the subnet to allocate an IP address from.
Can reference an OpenStackSubnet resource's output (via value_from for
InfraChart DAG wiring) or be a literal subnet UUID.
If omitted, OpenStack auto-selects a subnet on the port's network.

- references: OpenStackSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.fixedIps[].ipAddress

`string`

ip_address requests a specific IP address from the subnet.
If omitted, an available IP from the subnet's allocation pool is assigned.
The IP must belong to the subnet's CIDR and be within an allocation pool.

### spec.securityGroupIds

`[]string | valueFrom`

security_group_ids is the list of security groups to apply to this port.
Each entry can reference an OpenStackSecurityGroup resource's output
(via value_from) or be a literal security group UUID.
If omitted and no_security_groups is false, OpenStack applies the default
security group for the project.
Mutually exclusive with no_security_groups.

- references: OpenStackSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.noSecurityGroups

`bool`

no_security_groups explicitly removes all security groups from this port,
including the default security group that OpenStack normally applies.
Set to true for ports that should have unrestricted traffic (e.g., load
balancer VIPs, network appliance ports).
Mutually exclusive with security_group_ids.

### spec.adminStateUp

`bool` · optional (explicit presence)

admin_state_up controls the administrative state of the port.
When false, the port is administratively down and will not forward traffic.
Default: true

- default: `true`

### spec.macAddress

`string`

mac_address specifies a specific MAC address for this port.
If omitted, OpenStack auto-assigns a MAC from the network's allocation pool.
ForceNew: changing the MAC recreates the port.
Use case: network bonding, DPDK, or license-tied MAC addresses.

### spec.portSecurityEnabled

`bool` · optional (explicit presence)

port_security_enabled controls whether port security is enforced on this port.
When enabled, only traffic matching the port's security groups and allowed
address pairs is permitted. When disabled, all traffic passes regardless of
security groups.
If omitted, inherits from the network's port_security_enabled setting.

### spec.description

`string`

description is a human-readable description of the port.
Stored on the OpenStack resource and visible in Horizon and API responses.

### spec.tags

`[]string`

tags are string tags to associate with the port in OpenStack.
Tags are stored on the OpenStack resource and can be used for filtering
and organization in the OpenStack API and Horizon dashboard.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

region overrides the region from the provider config for this port.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Validation Rules

- `security_groups.mutual_exclusion`: no_security_groups and security_group_ids are mutually exclusive -- use one or neither

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackNetworkPort, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.port_id` | `string` | port_id is the unique identifier (UUID) of the port resource in OpenStack. This is the primary FK target for downstream components. |
| `status.outputs.mac_address` | `string` | mac_address is the MAC address assigned to the port. Auto-generated by OpenStack unless explicitly set in the spec. |
| `status.outputs.all_fixed_ips` | `[]string` | all_fixed_ips is the computed list of all IP addresses assigned to this port. Includes both explicitly requested IPs and auto-assigned IPs. Example: ["192.168.1.10", "10.0.0.5"] |
| `status.outputs.all_security_group_ids` | `[]string` | all_security_group_ids is the computed list of all security group UUIDs applied to this port, including the default SG if no explicit SGs were set. |
| `status.outputs.region` | `string` | region is the OpenStack region where the port was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.networkId` | OpenStackNetwork | `status.outputs.network_id` |
| `spec.fixedIps[].subnetId` | OpenStackSubnet | `status.outputs.subnet_id` |
| `spec.securityGroupIds` | OpenStackSecurityGroup | `status.outputs.security_group_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackFloatingIp | `spec.portId` | `status.outputs.port_id` |
| OpenStackFloatingIpAssociate | `spec.portId` | `status.outputs.port_id` |
| OpenStackInstance | `spec.networks[].port` | `status.outputs.port_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
