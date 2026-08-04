# ScalewayInstanceSecurityGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1`

ScalewayInstanceSecurityGroupSpec defines the specification for a Scaleway
Instance Security Group -- a stateful (by default) firewall that controls
inbound and outbound traffic to Scaleway Instances.

Security groups in Scaleway are zonal resources that operate at the Instance
level. Unlike VPC/Private Network firewalls in some other cloud providers,
Scaleway security groups are NOT attached to networks. Instead, they are
assigned to individual Instances via the Instance's `security_group_id` field.

Key Scaleway semantics:
  - Default policies ("accept" or "drop") control what happens to traffic
    that matches NO rule. By default, both inbound and outbound default
    policies are "accept" (allow-all), meaning rules are used to DROP
    specific traffic. Set a default policy to "drop" to switch to an
    allowlist model where only explicitly accepted traffic is permitted.
  - Rules are evaluated in order. The first matching rule wins.
  - Actions are "accept" or "drop" (not "allow"/"deny" as in some providers).
  - Protocols are uppercase: "TCP", "UDP", "ICMP", "ANY".
  - Stateful mode (default) automatically permits return traffic for
    accepted connections. Disable only for advanced stateless routing.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zone` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.stateful` | `bool` |  | `true` |  |
| `spec.inboundDefaultPolicy` | `string` |  | `accept` |  |
| `spec.outboundDefaultPolicy` | `string` |  | `accept` |  |
| `spec.enableDefaultSecurity` | `bool` |  | `true` |  |
| `spec.inboundRules` | `[]ScalewaySecurityGroupInboundRule` |  |  |  |
| `spec.inboundRules[].action` | `string` | yes |  |  |
| `spec.inboundRules[].protocol` | `string` |  | `TCP` |  |
| `spec.inboundRules[].portRange` | `string` |  |  |  |
| `spec.inboundRules[].ipRange` | `string` |  |  |  |
| `spec.outboundRules` | `[]ScalewaySecurityGroupOutboundRule` |  |  |  |
| `spec.outboundRules[].action` | `string` | yes |  |  |
| `spec.outboundRules[].protocol` | `string` |  | `TCP` |  |
| `spec.outboundRules[].portRange` | `string` |  |  |  |
| `spec.outboundRules[].ipRange` | `string` |  |  |  |

## Field Details

### spec.zone

`string` · required

The Scaleway zone where the security group will be created.
Examples: "fr-par-1", "nl-ams-1", "pl-waw-1"

Security groups are zonal resources. The zone must match the zone of
the Instances that will use this security group.

This field is required and cannot be changed after creation.

- rule: {"required":true}

### spec.description

`string`

Optional human-readable description of the security group's purpose.
Example: "Web tier firewall -- allows HTTP/HTTPS inbound, all outbound"

### spec.stateful

`bool`

Whether the security group is stateful.

When true (default), return traffic for accepted connections is
automatically permitted. For example, if an inbound rule accepts TCP
port 80, the outbound response packets are automatically allowed
without needing an explicit outbound rule.

Set to false only for advanced use cases such as high-throughput
stateless routing or network appliances where you need full control
over both directions.

Default: true

- default: `true`

### spec.inboundDefaultPolicy

`string`

Default policy for inbound traffic that matches no rule.

Possible values:
  - "accept" -- Allow all inbound traffic unless a rule drops it.
                Use this when you want a denylist model (block specific traffic).
  - "drop"   -- Drop all inbound traffic unless a rule accepts it.
                Use this when you want an allowlist model (only permit specific traffic).

Default: "accept"

Recommendation: Use "drop" for production workloads and define explicit
accept rules for known traffic patterns (SSH, HTTP, HTTPS, etc.).

- default: `accept`

### spec.outboundDefaultPolicy

`string`

Default policy for outbound traffic that matches no rule.

Possible values:
  - "accept" -- Allow all outbound traffic unless a rule drops it.
  - "drop"   -- Drop all outbound traffic unless a rule accepts it.

Default: "accept"

Recommendation: Keep as "accept" unless you need strict egress control.
Most workloads need unrestricted outbound access for package updates,
DNS resolution, and API calls.

- default: `accept`

### spec.enableDefaultSecurity

`bool`

Whether to enable Scaleway's default SMTP security.

When true (default), outbound SMTP traffic on ports 25, 465, and 587
is blocked to prevent spam abuse. This is a Scaleway account-level
protection that applies on top of your rules.

Set to false ONLY if your Scaleway account is authorized for SMTP
sending and the Instance needs to send email directly. Most workloads
should use a third-party email service (e.g., SendGrid, Mailgun)
instead of direct SMTP.

Default: true

- default: `true`

### spec.inboundRules

`[]ScalewaySecurityGroupInboundRule`

Inbound (ingress) rules: traffic allowed or dropped TO instances.

Rules are evaluated in order -- the first matching rule wins.
Traffic that matches no rule is handled by inbound_default_policy.

When inbound_default_policy is "drop", define accept rules for
traffic you want to permit (allowlist model).
When inbound_default_policy is "accept", define drop rules for
traffic you want to block (denylist model).

### spec.inboundRules[].action

`string` · required

Action to take when traffic matches this rule.

Possible values: "accept" or "drop".
This field is required.

- rule: {"required":true,"string":{"pattern":"^(accept|drop)$"}}

### spec.inboundRules[].protocol

`string`

IP protocol for this rule.

Possible values: "TCP", "UDP", "ICMP", "ANY".
Note: Scaleway uses uppercase protocol names.

Default: "TCP" (if omitted).

- default: `TCP`

### spec.inboundRules[].portRange

`string`

Port or port range for this rule.

Formats:
  - Single port: "80", "443", "22"
  - Port range:  "8000-9000", "30000-32767"

If omitted, the rule applies to ALL ports.
Ignored when protocol is "ICMP".

### spec.inboundRules[].ipRange

`string`

Source IP range in CIDR notation.

Examples: "0.0.0.0/0" (all IPv4), "203.0.113.10/32" (single host),
          "10.0.0.0/8" (private range)

If omitted, the rule applies to ALL source IPs.

### spec.outboundRules

`[]ScalewaySecurityGroupOutboundRule`

Outbound (egress) rules: traffic allowed or dropped FROM instances.

Rules are evaluated in order -- the first matching rule wins.
Traffic that matches no rule is handled by outbound_default_policy.

### spec.outboundRules[].action

`string` · required

Action to take when traffic matches this rule.

Possible values: "accept" or "drop".
This field is required.

- rule: {"required":true,"string":{"pattern":"^(accept|drop)$"}}

### spec.outboundRules[].protocol

`string`

IP protocol for this rule.

Possible values: "TCP", "UDP", "ICMP", "ANY".
Note: Scaleway uses uppercase protocol names.

Default: "TCP" (if omitted).

- default: `TCP`

### spec.outboundRules[].portRange

`string`

Port or port range for this rule.

Formats:
  - Single port: "443", "53"
  - Port range:  "1024-65535"

If omitted, the rule applies to ALL ports.
Ignored when protocol is "ICMP".

### spec.outboundRules[].ipRange

`string`

Destination IP range in CIDR notation.

Examples: "0.0.0.0/0" (all IPv4), "10.0.0.0/8" (private range)

If omitted, the rule applies to ALL destination IPs.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayInstanceSecurityGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.security_group_id` | `string` | The unique identifier of the created security group. This is the primary output referenced by downstream resources (e.g., ScalewayInstance) via StringValueOrRef on the security_group_id field. Format: Scaleway zoned ID "{zone}/{uuid}" (e.g., "fr-par-1/11111111-1111-...") |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| ScalewayInstance | `spec.securityGroupId` | `status.outputs.security_group_id` |

## See Also

- [Overview](./README.md)
