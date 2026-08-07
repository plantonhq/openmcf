# OpenStackSecurityGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackSecurityGroupSpec defines the configuration for an OpenStack Neutron
security group with optional inline rules.

A security group acts as a virtual firewall for instances and ports. It contains
a set of rules that control ingress and egress traffic. Security groups are
referenced by instances, network ports, and other networking resources.

The security group name is derived from metadata.name.

This component supports two modes for managing rules:
  - **Inline rules** (via the `rules` field): Convenient for self-contained
    security groups where all rules are defined in one place.
  - **Standalone rules** (via the separate OpenStackSecurityGroupRule component):
    For DAG-visible, independently managed rules in InfraCharts.

Terraform resource: openstack_networking_secgroup_v2 + openstack_networking_secgroup_rule_v2 (for inline rules)
Pulumi resource: openstack.networking.SecurityGroup + openstack.networking.SecGroupRule (for inline rules)

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackSecurityGroup
metadata:
  name: test-sg
spec:
  description: "Test security group for local development"
  delete_default_rules: true
  rules:
    - key: allow-ssh
      direction: ingress
      ethertype: IPv4
      protocol: tcp
      port_range_min: 22
      port_range_max: 22
      remote_ip_prefix: "0.0.0.0/0"
      description: "Allow SSH from anywhere"
    - key: allow-http
      direction: ingress
      ethertype: IPv4
      protocol: tcp
      port_range_min: 80
      port_range_max: 80
      remote_ip_prefix: "0.0.0.0/0"
      description: "Allow HTTP from anywhere"
    - key: egress-all-ipv4
      direction: egress
      ethertype: IPv4
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.description` | `string` |  |  |  |
| `spec.deleteDefaultRules` | `bool` |  |  |  |
| `spec.stateful` | `bool` |  |  |  |
| `spec.rules` | `[]SecurityGroupRule` |  |  |  |
| `spec.rules[].key` | `string` | yes |  |  |
| `spec.rules[].direction` | `string` |  |  |  |
| `spec.rules[].ethertype` | `string` |  |  |  |
| `spec.rules[].protocol` | `string` |  |  |  |
| `spec.rules[].portRangeMin` | `int32` |  |  |  |
| `spec.rules[].portRangeMax` | `int32` |  |  |  |
| `spec.rules[].remoteIpPrefix` | `string` |  |  |  |
| `spec.rules[].remoteGroupId` | `string` |  |  |  |
| `spec.rules[].description` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.description

`string`

description is a human-readable description of the security group.
This is stored on the OpenStack resource and visible in Horizon and API responses.

### spec.deleteDefaultRules

`bool` · optional (explicit presence)

delete_default_rules controls whether OpenStack's automatically created default
egress rules are deleted after the security group is created.
OpenStack creates two default rules on every new security group:
  - Allow all egress IPv4 traffic
  - Allow all egress IPv6 traffic
Set to true to start with a completely empty rule set (zero-trust baseline).
If omitted, the default rules are kept.
This is a create-time setting and cannot be changed after creation.

### spec.stateful

`bool` · optional (explicit presence)

stateful controls whether the security group operates in stateful or stateless mode.
Stateful (default in OpenStack): Return traffic is automatically allowed regardless
of rules. Stateless: Return traffic must be explicitly permitted by rules.
Stateless security groups offer better performance for high-throughput workloads
but require more careful rule management.
If omitted, OpenStack uses its deployment default (typically stateful).
Note: Not all OpenStack deployments support stateless security groups.

### spec.rules

`[]SecurityGroupRule`

rules defines inline security group rules to create alongside the security group.
Each rule is provisioned as a separate openstack_networking_secgroup_rule_v2 resource,
keyed by the rule's `key` field for stable IaC state management.
For DAG-visible, independently managed rules, use the OpenStackSecurityGroupRule
component instead.

- rule: port_range_min and port_range_max must both be set or both unset
- rule: protocol is required when port ranges are specified
- rule: remote_group_id and remote_ip_prefix are mutually exclusive

### spec.rules[].key

`string` · required

key is a unique identifier for this rule within the security group.
Used as the resource key in IaC state management (Terraform for_each key,
Pulumi resource name suffix). Must be unique across all rules in the spec.
Use descriptive, kebab-case names like "allow-ssh", "egress-all-ipv4",
"allow-https-from-lb".

- rule: {"string":{"minLen":"1"}}

### spec.rules[].direction

`string`

direction specifies whether the rule applies to incoming or outgoing traffic.
Must be "ingress" or "egress".

- rule: {"string":{"in":["ingress","egress"]}}

### spec.rules[].ethertype

`string`

ethertype specifies the layer-3 protocol type for the rule.
Must be "IPv4" or "IPv6".

- rule: {"string":{"in":["IPv4","IPv6"]}}

### spec.rules[].protocol

`string`

protocol is the IP protocol for the rule.
Common values: "tcp", "udp", "icmp", "icmpv6".
Also accepts any IANA protocol name or number (0-255).
If omitted, the rule applies to all protocols.

### spec.rules[].portRangeMin

`int32` · optional (explicit presence)

port_range_min is the minimum port number for the rule.
For TCP/UDP: the start of the port range (0-65535).
For ICMP: the ICMP type (0-255). Type 0 = Echo Reply, Type 8 = Echo Request.
Must be set together with port_range_max. Requires protocol to be set.

### spec.rules[].portRangeMax

`int32` · optional (explicit presence)

port_range_max is the maximum port number for the rule.
For TCP/UDP: the end of the port range (0-65535).
For ICMP: the ICMP code (0-255). Code 0 is valid and commonly used.
Must be set together with port_range_min. Requires protocol to be set.

### spec.rules[].remoteIpPrefix

`string`

remote_ip_prefix restricts the rule to traffic from/to a specific CIDR.
For ingress: the source IP range. For egress: the destination IP range.
Example: "0.0.0.0/0" (all IPv4), "10.0.0.0/8" (private range), "203.0.113.0/24".
Mutually exclusive with remote_group_id.

### spec.rules[].remoteGroupId

`string`

remote_group_id restricts the rule to traffic from/to instances in another
security group (or the same security group for self-referencing rules).
This is a literal UUID of an existing security group.
For self-referencing rules in InfraCharts, use the standalone
OpenStackSecurityGroupRule component with a StringValueOrRef FK instead.
Mutually exclusive with remote_ip_prefix.

### spec.rules[].description

`string`

description is a human-readable description of the rule.
Stored on the OpenStack rule resource.

### spec.tags

`[]string`

tags are string tags to associate with the security group in OpenStack.
Tags are stored on the OpenStack resource and can be used for filtering
and organization in the OpenStack API and Horizon dashboard.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

region overrides the region from the provider config for this security group.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Validation Rules

- `rules.unique_keys`: inline rule keys must be unique within the security group

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackSecurityGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.security_group_id` | `string` | security_group_id is the unique identifier (UUID) of the security group in OpenStack. This is the primary output used as a foreign key by downstream components. |
| `status.outputs.name` | `string` | name is the name of the security group (derived from metadata.name). |
| `status.outputs.region` | `string` | region is the OpenStack region where the security group was created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackInstance | `spec.securityGroups` | `status.outputs.name` |
| OpenStackNetworkPort | `spec.securityGroupIds` | `status.outputs.security_group_id` |
| OpenStackSecurityGroupRule | `spec.securityGroupId` | `status.outputs.security_group_id` |
| OpenStackSecurityGroupRule | `spec.remoteGroupId` | `status.outputs.security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
