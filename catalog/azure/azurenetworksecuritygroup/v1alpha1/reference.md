# AzureNetworkSecurityGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureNetworkSecurityGroupSpec** defines the configuration for creating
an Azure Network Security Group (NSG): the stateful firewall that filters
inbound and outbound traffic for everything deployed in the subnets and
NICs it guards.

The security rules are folded into the NSG rather than modeled as their
own resource: a rule has no life outside its group, is never referenced
by anything else, and the group's rule set is designed and reviewed as
one unit. An NSG with no rules is still meaningful -- Azure's implicit
defaults then govern (allow VNet-internal traffic and load-balancer
probes, deny all other inbound, allow all outbound).

The subnet-side attachment is deliberately not modeled here: a subnet
declares which NSG guards it (AzureSubnet's network_security_group_id),
matching Azure's model, so one NSG serves many subnets without listing
them.

**Key Azure behavior**: user-defined rules take priorities 100-4096 and
are evaluated lowest-number-first within each direction; the first match
wins. Azure's implicit default rules sit at priorities 65000-65500 and
only apply when nothing else matched.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNetworkSecurityGroup
metadata:
  name: test-nsg
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: test-nsg
  securityRules:
    - name: allow-https
      priority: 100
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "443"
    - name: allow-web-multi
      priority: 200
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRanges:
        - "80"
        - "8080"
      sourceAddressPrefixes:
        - "10.1.0.0/16"
        - "10.2.0.0/16"
    - name: deny-all-inbound
      priority: 4096
      direction: INBOUND
      access: DENY
      protocol: ANY
      destinationPortRange: "*"
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.securityRules` | `[]AzureNetworkSecurityGroupRule` |  |  |  |
| `spec.securityRules[].name` | `string` | yes |  |  |
| `spec.securityRules[].description` | `string` |  |  |  |
| `spec.securityRules[].priority` | `int32` | yes |  |  |
| `spec.securityRules[].direction` | `enum` | yes |  |  |
| `spec.securityRules[].access` | `enum` | yes |  |  |
| `spec.securityRules[].protocol` | `enum` | yes |  |  |
| `spec.securityRules[].sourcePortRange` | `string` |  |  |  |
| `spec.securityRules[].sourcePortRanges` | `[]string` |  |  |  |
| `spec.securityRules[].destinationPortRange` | `string` |  |  |  |
| `spec.securityRules[].destinationPortRanges` | `[]string` |  |  |  |
| `spec.securityRules[].sourceAddressPrefix` | `string` |  |  |  |
| `spec.securityRules[].sourceAddressPrefixes` | `[]string` |  |  |  |
| `spec.securityRules[].sourceApplicationSecurityGroupIds` | `[]string \| valueFrom` |  |  | AzureApplicationSecurityGroup (`status.outputs.application_security_group_id`) |
| `spec.securityRules[].destinationAddressPrefix` | `string` |  |  |  |
| `spec.securityRules[].destinationAddressPrefixes` | `[]string` |  |  |  |
| `spec.securityRules[].destinationApplicationSecurityGroupIds` | `[]string \| valueFrom` |  |  | AzureApplicationSecurityGroup (`status.outputs.application_security_group_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the NSG will be created, e.g. "eastus",
"westeurope". Must match the region of the subnets and NICs it will be
associated with. Changing the region replaces the group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the NSG will be created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the NSG, unique within the resource group. 1-80 characters
(alphanumerics, underscores, periods, and hyphens; must start with a
letter or number and end with a letter, number, or underscore).
Changing the name replaces the group, detaching it from every subnet
and NIC until the replacement is re-attached -- name it after the tier
it guards ("web-tier", "data-tier").

- rule: NSG names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.securityRules

`[]AzureNetworkSecurityGroupRule`

The security rules. Each is a 5-tuple filter (source, destination,
port, protocol, direction) with an access decision and a priority.
Rules update in place and take effect immediately for every subnet and
NIC the group guards. An empty list is meaningful: Azure's implicit
default rules then govern all traffic.

- rule: Set source ports as either source_port_range or source_port_ranges, not both (unset means any)
- rule: Set exactly one of destination_port_range or destination_port_ranges (use "*" for any port)
- rule: Set at most one source addressing style: source_address_prefix, source_address_prefixes, or source_application_security_group_ids (all unset means any)
- rule: Set at most one destination addressing style: destination_address_prefix, destination_address_prefixes, or destination_application_security_group_ids (all unset means any)

### spec.securityRules[].name

`string` · required

The rule's name, unique within the NSG. Shown in the portal, CLI, and
flow logs -- "allow-https-inbound" reads better than "rule1". 1-80
characters. Renaming replaces the rule (a momentary gap in an
otherwise-in-place rule set update).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.securityRules[].description

`string`

Optional human-readable statement of intent, up to 140 characters --
e.g. "Allow HTTPS from the corporate VPN range". Auditors and future
operators read these before they read the tuples.

- rule: {"string":{"maxLen":"140"}}

### spec.securityRules[].priority

`int32` · required

The evaluation priority, 100-4096, unique per direction within the
NSG. Lower numbers evaluate first; the first matching rule decides.
Leave gaps (100, 200, 300...) so rules can be inserted later without
renumbering.

- rule: {"required":true,"int32":{"lte":4096,"gte":100}}

### spec.securityRules[].direction

`enum` · required

Whether the rule filters traffic entering (INBOUND) or leaving
(OUTBOUND) the guarded resources.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_network_security_group_rule_direction_unspecified` -- Not specified -- invalid; every rule must declare its direction.
- `INBOUND` -- Traffic entering the guarded resources.
- `OUTBOUND` -- Traffic leaving the guarded resources.

### spec.securityRules[].access

`enum` · required

Whether matching traffic is permitted (ALLOW) or blocked (DENY).

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_network_security_group_rule_access_unspecified` -- Not specified -- invalid; every rule must declare its decision.
- `ALLOW` -- Permit matching traffic.
- `DENY` -- Block matching traffic.

### spec.securityRules[].protocol

`enum` · required

The network protocol the rule matches. ANY matches everything; AH and
ESP match IPsec traffic (site-to-site VPN scenarios).

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_network_security_group_rule_protocol_unspecified` -- Not specified -- invalid; every rule must declare its protocol.
- `ANY` -- Any protocol (ARM's "*").
- `TCP` -- TCP traffic.
- `UDP` -- UDP traffic.
- `ICMP` -- ICMP traffic (ping, traceroute).
- `AH` -- IPsec Authentication Header traffic.
- `ESP` -- IPsec Encapsulating Security Payload traffic.

### spec.securityRules[].sourcePortRange

`string` · optional (explicit presence)

The source port or port range: a single port ("443"), a range
("1024-65535"), or "*" for any. Unset (with source_port_ranges empty)
means any -- the right choice for almost every rule, since client
source ports are ephemeral. Never combined with source_port_ranges.

### spec.securityRules[].sourcePortRanges

`[]string`

Multiple source ports/ranges in one rule. Never combined with
source_port_range.

### spec.securityRules[].destinationPortRange

`string` · optional (explicit presence)

The destination port or port range -- the field that says what service
the rule is about: "22" (SSH), "443" (HTTPS), "5432" (PostgreSQL), or
"*" for any. Exactly one of destination_port_range or
destination_port_ranges must be set.

### spec.securityRules[].destinationPortRanges

`[]string`

Multiple destination ports/ranges in one rule (e.g. ["80", "443"]).
Exactly one of destination_port_range or destination_port_ranges must
be set.

### spec.securityRules[].sourceAddressPrefix

`string` · optional (explicit presence)

The source as a single prefix: a CIDR ("10.0.0.0/8"), an IP, an Azure
service tag ("VirtualNetwork", "AzureLoadBalancer", "Internet"), or
"*". Service tags and "*" only work here, not in the plural form.
At most one source addressing style may be set; all unset means any.

### spec.securityRules[].sourceAddressPrefixes

`[]string`

The source as multiple CIDRs/IPs (e.g. several VPN ranges in one
rule). Service tags are not accepted in the list form. At most one
source addressing style may be set.

### spec.securityRules[].sourceApplicationSecurityGroupIds

`[]string | valueFrom`

The source as application security group membership -- identity-based
addressing that follows workloads as they scale instead of pinning
CIDRs. Each entry is an application security group by ARM ID, or a
reference to an AzureApplicationSecurityGroup's output (up to 10). At
most one source addressing style may be set.

- references: AzureApplicationSecurityGroup (`status.outputs.application_security_group_id`)
- rule: {"repeated":{"maxItems":"10"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureApplicationSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.application_security_group_id}} -- a bare string does not parse

### spec.securityRules[].destinationAddressPrefix

`string` · optional (explicit presence)

The destination as a single prefix: a CIDR, an IP, a service tag, or
"*". Service tags and "*" only work here, not in the plural form. At
most one destination addressing style may be set; all unset means any.

### spec.securityRules[].destinationAddressPrefixes

`[]string`

The destination as multiple CIDRs/IPs. Service tags are not accepted
in the list form. At most one destination addressing style may be set.

### spec.securityRules[].destinationApplicationSecurityGroupIds

`[]string | valueFrom`

The destination as application security group membership. Each entry
is an application security group by ARM ID, or a reference to an
AzureApplicationSecurityGroup's output (up to 10). At most one
destination addressing style may be set.

- references: AzureApplicationSecurityGroup (`status.outputs.application_security_group_id`)
- rule: {"repeated":{"maxItems":"10"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureApplicationSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.application_security_group_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

Free-form tags applied to the NSG, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Tags are Azure's governance surface -- Azure
Policy enforces them and Microsoft Cost Management groups by them.
Updatable in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureNetworkSecurityGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.network_security_group_id` | `string` | The Azure Resource Manager ID of the NSG. This is the primary output: AzureSubnet's network_security_group_id references it to attach the group to a subnet. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/networkSecurityGroups/{name} |
| `status.outputs.network_security_group_name` | `string` | The name of the Network Security Group. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.securityRules[].sourceApplicationSecurityGroupIds` | AzureApplicationSecurityGroup | `status.outputs.application_security_group_id` |
| `spec.securityRules[].destinationApplicationSecurityGroupIds` | AzureApplicationSecurityGroup | `status.outputs.application_security_group_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureNetworkInterface | `spec.networkSecurityGroupId` | `status.outputs.network_security_group_id` |
| AzureSubnet | `spec.networkSecurityGroupId` | `status.outputs.network_security_group_id` |
| AzureVirtualMachineScaleSet | `spec.networkInterfaces[].networkSecurityGroupId` | `status.outputs.network_security_group_id` |

## See Also

- [Overview](../README.md)
