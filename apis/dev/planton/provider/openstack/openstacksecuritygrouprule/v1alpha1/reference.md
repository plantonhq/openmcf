# OpenStackSecurityGroupRule

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackSecurityGroupRuleSpec defines the configuration for a standalone
OpenStack Neutron security group rule.

This is the standalone counterpart to inline rules defined in
OpenStackSecurityGroup.rules[]. Use this component when individual rules
need to be independently managed and visible as separate nodes in InfraChart
DAG visualizations. The key advantage over inline rules is that the
security_group_id and remote_group_id fields support StringValueOrRef,
enabling cross-resource references that the InfraChart engine can resolve.

All fields on the underlying OpenStack resource are ForceNew: any change
recreates the rule. This is expected for declarative infrastructure.

The resource name in Planton is derived from metadata.name.

Terraform resource: openstack_networking_secgroup_rule_v2
Pulumi resource: openstack.networking.SecGroupRule

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackSecurityGroupRule
metadata:
  name: test-allow-ssh
spec:
  security_group_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  direction: ingress
  ethertype: IPv4
  protocol: tcp
  port_range_min: 22
  port_range_max: 22
  remote_ip_prefix: "0.0.0.0/0"
  description: "Test rule -- allow SSH from anywhere"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.securityGroupId` | `string \| valueFrom` | yes |  | OpenStackSecurityGroup (`status.outputs.security_group_id`) |
| `spec.direction` | `string` |  |  |  |
| `spec.ethertype` | `string` |  |  |  |
| `spec.protocol` | `string` |  |  |  |
| `spec.portRangeMin` | `int32` |  |  |  |
| `spec.portRangeMax` | `int32` |  |  |  |
| `spec.remoteIpPrefix` | `string` |  |  |  |
| `spec.remoteGroupId` | `string \| valueFrom` |  |  | OpenStackSecurityGroup (`status.outputs.security_group_id`) |
| `spec.description` | `string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.securityGroupId

`string | valueFrom` · required

security_group_id is the ID of the security group this rule belongs to.
This is a required foreign key -- every standalone rule must reference
exactly one security group.
Can reference an OpenStackSecurityGroup resource's output or be a literal UUID.

- references: OpenStackSecurityGroup (`status.outputs.security_group_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.direction

`string`

direction specifies whether the rule applies to incoming or outgoing traffic.
Must be "ingress" or "egress".

- rule: {"string":{"in":["ingress","egress"]}}

### spec.ethertype

`string`

ethertype specifies the layer-3 protocol type for the rule.
Must be "IPv4" or "IPv6".

- rule: {"string":{"in":["IPv4","IPv6"]}}

### spec.protocol

`string`

protocol is the IP protocol for the rule.
Common values: "tcp", "udp", "icmp", "icmpv6".
Also accepts any IANA protocol name or number (0-255).
If omitted, the rule applies to all protocols.

### spec.portRangeMin

`int32` · optional (explicit presence)

port_range_min is the minimum port number for the rule.
For TCP/UDP: the start of the port range (0-65535).
For ICMP: the ICMP type (0-255). Type 0 = Echo Reply, Type 8 = Echo Request.
Must be set together with port_range_max. Requires protocol to be set.

### spec.portRangeMax

`int32` · optional (explicit presence)

port_range_max is the maximum port number for the rule.
For TCP/UDP: the end of the port range (0-65535).
For ICMP: the ICMP code (0-255). Code 0 is valid and commonly used.
Must be set together with port_range_min. Requires protocol to be set.

### spec.remoteIpPrefix

`string`

remote_ip_prefix restricts the rule to traffic from/to a specific CIDR.
For ingress: the source IP range. For egress: the destination IP range.
Example: "0.0.0.0/0" (all IPv4), "10.0.0.0/8" (private range), "203.0.113.0/24".
Mutually exclusive with remote_group_id.

### spec.remoteGroupId

`string | valueFrom`

remote_group_id restricts the rule to traffic from/to instances in another
security group (or the same security group for self-referencing rules).
This is an optional foreign key to OpenStackSecurityGroup -- the key advantage
of standalone rules over inline rules, enabling InfraChart DAG wiring for
cross-security-group references.
Mutually exclusive with remote_ip_prefix.

- references: OpenStackSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.description

`string`

description is a human-readable description of the rule.
Stored on the OpenStack rule resource and visible in Horizon and API responses.

### spec.region

`string`

region overrides the region from the provider config for this rule.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Validation Rules

- `port_range.both_or_neither`: port_range_min and port_range_max must both be set or both unset
- `port_range.requires_protocol`: protocol is required when port ranges are specified
- `remote_source.mutual_exclusion`: remote_group_id and remote_ip_prefix are mutually exclusive

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackSecurityGroupRule, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.rule_id` | `string` | rule_id is the unique identifier (UUID) of the security group rule in OpenStack. This is the Terraform resource ID. |
| `status.outputs.security_group_id` | `string` | security_group_id is the UUID of the security group this rule belongs to. |
| `status.outputs.direction` | `string` | direction is the direction of the rule ("ingress" or "egress"). |
| `status.outputs.protocol` | `string` | protocol is the IP protocol of the rule (e.g., "tcp", "udp", "icmp"). Empty string means the rule applies to all protocols. |
| `status.outputs.port_range_min` | `int32` | port_range_min is the lower bound of the port range (or ICMP type). |
| `status.outputs.port_range_max` | `int32` | port_range_max is the upper bound of the port range (or ICMP code). |
| `status.outputs.region` | `string` | region is the OpenStack region where the rule was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.securityGroupId` | OpenStackSecurityGroup | `status.outputs.security_group_id` |
| `spec.remoteGroupId` | OpenStackSecurityGroup | `status.outputs.security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
