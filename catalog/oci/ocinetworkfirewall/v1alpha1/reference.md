# OciNetworkFirewall

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciNetworkFirewallSpec defines the specification for an OCI Network
Firewall bundled with an inline firewall policy and its core
sub-resources (address lists, services, service lists, URL lists,
security rules).

The firewall appliance is deployed into a subnet and inspects traffic
according to the security rules defined in the policy. Traffic
matching is based on source/destination IP addresses, TCP/UDP ports,
and URL patterns. Rules are evaluated in order; priority is derived
from the list position of each security rule.

Key behaviors:
  - subnet_id is immutable after creation (ForceNew)
  - ipv4_address and ipv6_address are immutable (ForceNew)
  - availability_domain is immutable (ForceNew)
  - Policy sub-resource names are immutable (ForceNew)
  - Security rules are evaluated in list order (priority_order = index + 1)
  - Sub-resources reference each other by name within the YAML manifest

Excluded from v1:
  - Applications / Application Groups (ICMP type/code matching)
  - Decryption Profiles / Rules / Mapped Secrets (TLS inspection)
  - NAT Rules (specialized NAT handling)
  - Tunnel Inspection Rules (VXLAN inspection)
  - defined_tags, system_tags, freeform_tags (auto from labels)

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.ipv4Address` | `string` |  |  |  |
| `spec.ipv6Address` | `string` |  |  |  |
| `spec.availabilityDomain` | `string` |  |  |  |
| `spec.networkSecurityGroupIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.natConfiguration` | `NatConfiguration` |  |  |  |
| `spec.natConfiguration.mustEnablePrivateNat` | `bool` |  |  |  |
| `spec.shape` | `string` |  |  |  |
| `spec.policy` | `Policy` | yes |  |  |
| `spec.policy.displayName` | `string` |  |  |  |
| `spec.policy.description` | `string` |  |  |  |
| `spec.policy.addressLists` | `[]AddressList` |  |  |  |
| `spec.policy.addressLists[].name` | `string` | yes |  |  |
| `spec.policy.addressLists[].type` | `enum` |  |  |  |
| `spec.policy.addressLists[].addresses` | `[]string` | yes |  |  |
| `spec.policy.addressLists[].description` | `string` |  |  |  |
| `spec.policy.services` | `[]Service` |  |  |  |
| `spec.policy.services[].name` | `string` | yes |  |  |
| `spec.policy.services[].type` | `enum` |  |  |  |
| `spec.policy.services[].portRanges` | `[]PortRange` | yes |  |  |
| `spec.policy.services[].portRanges[].minimumPort` | `int32` |  |  |  |
| `spec.policy.services[].portRanges[].maximumPort` | `int32` |  |  |  |
| `spec.policy.services[].description` | `string` |  |  |  |
| `spec.policy.serviceLists` | `[]ServiceList` |  |  |  |
| `spec.policy.serviceLists[].name` | `string` | yes |  |  |
| `spec.policy.serviceLists[].services` | `[]string` |  |  |  |
| `spec.policy.serviceLists[].description` | `string` |  |  |  |
| `spec.policy.urlLists` | `[]UrlList` |  |  |  |
| `spec.policy.urlLists[].name` | `string` | yes |  |  |
| `spec.policy.urlLists[].urls` | `[]UrlPattern` | yes |  |  |
| `spec.policy.urlLists[].urls[].pattern` | `string` | yes |  |  |
| `spec.policy.urlLists[].description` | `string` |  |  |  |
| `spec.policy.securityRules` | `[]SecurityRule` |  |  |  |
| `spec.policy.securityRules[].name` | `string` | yes |  |  |
| `spec.policy.securityRules[].action` | `enum` |  |  |  |
| `spec.policy.securityRules[].condition` | `SecurityRuleCondition` | yes |  |  |
| `spec.policy.securityRules[].condition.sourceAddresses` | `[]string` |  |  |  |
| `spec.policy.securityRules[].condition.destinationAddresses` | `[]string` |  |  |  |
| `spec.policy.securityRules[].condition.services` | `[]string` |  |  |  |
| `spec.policy.securityRules[].condition.urls` | `[]string` |  |  |  |
| `spec.policy.securityRules[].inspection` | `enum` |  |  |  |
| `spec.policy.securityRules[].description` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the firewall and policy will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.subnetId

`string | valueFrom` · required

OCID of the subnet where the firewall appliance will be deployed.
Immutable after creation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.displayName

`string`

Display name for the firewall. When omitted, the metadata name is used.

### spec.ipv4Address

`string`

Static IPv4 address for the firewall. When omitted, OCI auto-assigns
from the subnet. Immutable after creation.

### spec.ipv6Address

`string`

Static IPv6 address for the firewall. When omitted, OCI auto-assigns
if the subnet supports IPv6. Immutable after creation.

### spec.availabilityDomain

`string`

Availability domain for the firewall placement.
When omitted, OCI selects automatically. Immutable after creation.

### spec.networkSecurityGroupIds

`[]string | valueFrom`

OCIDs of network security groups applied to the firewall.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.natConfiguration

`NatConfiguration`

NAT configuration for the firewall. Controls whether private NAT
is used for egress traffic inspection.

### spec.natConfiguration.mustEnablePrivateNat

`bool`

Whether to enable private NAT for the firewall.
When true, the firewall uses a private IP for NAT operations
instead of a public IP.

### spec.shape

`string`

Firewall shape determining throughput capacity.
When omitted, OCI applies the default shape.

### spec.policy

`Policy` · required

Inline firewall policy defining security rules and their supporting
objects (address lists, services, URL lists). The policy is always
created as part of this component.

- rule: {"required":true}

### spec.policy.displayName

`string`

Display name for the policy. When omitted, defaults to
"{firewall_display_name}-policy".

### spec.policy.description

`string`

Description of the policy.

### spec.policy.addressLists

`[]AddressList`

IP and FQDN address lists referenced by security rules.

### spec.policy.addressLists[].name

`string` · required

Unique name within the policy. Immutable after creation.
Referenced by security rules via source_addresses and
destination_addresses.

- rule: {"string":{"minLen":"1"}}

### spec.policy.addressLists[].type

`enum`

Whether this list contains IP addresses/CIDRs or FQDNs.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `unspecified`
- `ip`
- `fqdn`

### spec.policy.addressLists[].addresses

`[]string` · required

Addresses in the list. IP lists accept CIDRs (e.g., "10.0.0.0/8")
and individual IPs. FQDN lists accept domain names
(e.g., "example.com").

- rule: {"repeated":{"minItems":"1"}}

### spec.policy.addressLists[].description

`string`

Optional description.

### spec.policy.services

`[]Service`

TCP/UDP port definitions referenced by security rules.

### spec.policy.services[].name

`string` · required

Unique name within the policy. Immutable after creation.
Referenced by security rules and service lists.

- rule: {"string":{"minLen":"1"}}

### spec.policy.services[].type

`enum`

Protocol for this service definition.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `service_type_unspecified`
- `tcp_service`
- `udp_service`

### spec.policy.services[].portRanges

`[]PortRange` · required

Port ranges for this service. At least one port range is required.

- rule: {"repeated":{"minItems":"1"}}
- rule: maximum_port must be >= minimum_port when set

### spec.policy.services[].portRanges[].minimumPort

`int32`

Start of the port range (inclusive). Must be 1-65535.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.policy.services[].portRanges[].maximumPort

`int32` · optional (explicit presence)

End of the port range (inclusive). When omitted, equals
minimum_port (single port). Must be >= minimum_port when set.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.policy.services[].description

`string`

Optional description.

### spec.policy.serviceLists

`[]ServiceList`

Groups of services for reuse across multiple security rules.

### spec.policy.serviceLists[].name

`string` · required

Unique name within the policy. Immutable after creation.
Referenced by security rules via the services condition field.

- rule: {"string":{"minLen":"1"}}

### spec.policy.serviceLists[].services

`[]string`

Names of services to include in this list. Each name must match
a Service defined in the policy's services field.

### spec.policy.serviceLists[].description

`string`

Optional description.

### spec.policy.urlLists

`[]UrlList`

URL pattern lists for L7 HTTP(S) traffic inspection.

### spec.policy.urlLists[].name

`string` · required

Unique name within the policy. Immutable after creation.
Referenced by security rules via the urls condition field.

- rule: {"string":{"minLen":"1"}}

### spec.policy.urlLists[].urls

`[]UrlPattern` · required

URL patterns in the list. At least one pattern is required.

- rule: {"repeated":{"minItems":"1"}}

### spec.policy.urlLists[].urls[].pattern

`string` · required

URL pattern string (e.g., "*.example.com",
"malware.example.com/path").

- rule: {"string":{"minLen":"1"}}

### spec.policy.urlLists[].description

`string`

Optional description.

### spec.policy.securityRules

`[]SecurityRule`

Security rules evaluated in list order. Priority is derived from
the position in this list (first rule = highest priority).

- rule: inspection must be set when action is inspect

### spec.policy.securityRules[].name

`string` · required

Unique name within the policy. Immutable after creation.

- rule: {"string":{"minLen":"1"}}

### spec.policy.securityRules[].action

`enum`

Action to take when traffic matches the condition.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `action_unspecified`
- `allow`
- `drop`
- `reject`
- `inspect`

### spec.policy.securityRules[].condition

`SecurityRuleCondition` · required

Traffic matching condition. References address lists, services,
and URL lists by name.

- rule: {"required":true}

### spec.policy.securityRules[].condition.sourceAddresses

`[]string`

Names of address lists matching source IP addresses.

### spec.policy.securityRules[].condition.destinationAddresses

`[]string`

Names of address lists matching destination IP addresses.

### spec.policy.securityRules[].condition.services

`[]string`

Names of services or service lists matching traffic ports.

### spec.policy.securityRules[].condition.urls

`[]string`

Names of URL lists matching HTTP(S) request URLs.

### spec.policy.securityRules[].inspection

`enum`

Inspection type. Required when action is inspect.

Allowed values (use exactly as shown):

- `inspection_unspecified`
- `intrusion_detection`
- `intrusion_prevention`

### spec.policy.securityRules[].description

`string`

Optional description.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciNetworkFirewall, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.firewall_id` | `string` | OCID of the network firewall. |
| `status.outputs.ipv4_address` | `string` | IPv4 address of the firewall appliance. Used for configuring route table entries to direct traffic through the firewall. |
| `status.outputs.policy_id` | `string` | OCID of the firewall policy. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.networkSecurityGroupIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](../README.md)
