# AzureFirewallPolicyRuleCollectionGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFirewallPolicyRuleCollectionGroupSpec** defines a rule collection
group inside an Azure Firewall Policy: an ordered document of rule
collections (application, network, and DNAT) that the policy's firewalls
enforce. A policy carries many groups -- typically one per team or per
application -- each deployed and updated independently, which is exactly
why the group is its own resource rather than a fold inside the policy.

**Processing order**: groups are evaluated by group priority, then
collections by collection priority. Across collection TYPES Azure always
processes DNAT rules first, then network rules, then application rules,
regardless of priorities -- priorities order collections WITHIN a type.
Lower numbers run first (100 is the highest precedence).

**Rules fold inside the group.** A rule has no ARM identity of its own
and nothing references an individual rule -- the group is the unit of
deployment, so the rules travel with it as one ordered document.

**ForceNew fields**: `name`, `firewall_policy_id`. Priority and all
collections update in place.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFirewallPolicyRuleCollectionGroup
metadata:
  name: test-rule-collection-group
spec:
  firewallPolicyId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/firewallPolicies/egress-baseline
  name: platform-baseline
  priority: 500
  applicationRuleCollections:
    - name: web-egress
      priority: 300
      action: ALLOW
      rules:
        - name: allow-github
          description: package pulls from github
          protocols:
            - type: HTTPS
              port: 443
          sourceAddresses:
            - "10.0.0.0/16"
          destinationFqdns:
            - "github.com"
            - "*.githubusercontent.com"
        - name: allow-ubuntu-updates
          protocols:
            - type: HTTP
              port: 80
            - type: HTTPS
              port: 443
          sourceAddresses:
            - "10.0.0.0/16"
          destinationFqdnTags:
            - "AzureKubernetesService"
  networkRuleCollections:
    - name: core-egress
      priority: 200
      action: ALLOW
      rules:
        - name: allow-dns
          protocols:
            - UDP
          sourceAddresses:
            - "10.0.0.0/16"
          destinationAddresses:
            - "8.8.8.8"
            - "8.8.4.4"
          destinationPorts:
            - "53"
        - name: allow-ntp
          protocols:
            - UDP
          sourceAddresses:
            - "10.0.0.0/16"
          destinationFqdns:
            - "time.windows.com"
          destinationPorts:
            - "123"
  natRuleCollections:
    - name: inbound-dnat
      priority: 100
      rules:
        - name: rdp-to-jumpbox
          protocols:
            - TCP
          sourceAddresses:
            - "*"
          destinationAddress: "203.0.113.10"
          destinationPorts:
            - "3389"
          translatedAddress: "10.0.1.4"
          translatedPort: 3389
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.firewallPolicyId` | `string \| valueFrom` | yes |  | AzureFirewallPolicy (`status.outputs.firewall_policy_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.priority` | `int32` | yes |  |  |
| `spec.applicationRuleCollections` | `[]AzureFirewallPolicyApplicationRuleCollection` |  |  |  |
| `spec.applicationRuleCollections[].name` | `string` | yes |  |  |
| `spec.applicationRuleCollections[].priority` | `int32` | yes |  |  |
| `spec.applicationRuleCollections[].action` | `enum` | yes |  |  |
| `spec.applicationRuleCollections[].rules` | `[]AzureFirewallPolicyApplicationRule` | yes |  |  |
| `spec.applicationRuleCollections[].rules[].name` | `string` | yes |  |  |
| `spec.applicationRuleCollections[].rules[].description` | `string` |  |  |  |
| `spec.applicationRuleCollections[].rules[].protocols` | `[]AzureFirewallPolicyApplicationProtocol` |  |  |  |
| `spec.applicationRuleCollections[].rules[].protocols[].type` | `enum` | yes |  |  |
| `spec.applicationRuleCollections[].rules[].protocols[].port` | `int32` | yes |  |  |
| `spec.applicationRuleCollections[].rules[].sourceAddresses` | `[]string` |  |  |  |
| `spec.applicationRuleCollections[].rules[].sourceIpGroups` | `[]string \| valueFrom` |  |  | AzureIpGroup (`status.outputs.ip_group_id`) |
| `spec.applicationRuleCollections[].rules[].destinationAddresses` | `[]string` |  |  |  |
| `spec.applicationRuleCollections[].rules[].destinationFqdns` | `[]string` |  |  |  |
| `spec.applicationRuleCollections[].rules[].destinationUrls` | `[]string` |  |  |  |
| `spec.applicationRuleCollections[].rules[].destinationFqdnTags` | `[]string` |  |  |  |
| `spec.applicationRuleCollections[].rules[].terminateTls` | `bool` |  |  |  |
| `spec.applicationRuleCollections[].rules[].webCategories` | `[]string` |  |  |  |
| `spec.applicationRuleCollections[].rules[].httpHeaders` | `[]AzureFirewallPolicyHttpHeader` |  |  |  |
| `spec.applicationRuleCollections[].rules[].httpHeaders[].name` | `string` | yes |  |  |
| `spec.applicationRuleCollections[].rules[].httpHeaders[].value` | `string` | yes |  |  |
| `spec.networkRuleCollections` | `[]AzureFirewallPolicyNetworkRuleCollection` |  |  |  |
| `spec.networkRuleCollections[].name` | `string` | yes |  |  |
| `spec.networkRuleCollections[].priority` | `int32` | yes |  |  |
| `spec.networkRuleCollections[].action` | `enum` | yes |  |  |
| `spec.networkRuleCollections[].rules` | `[]AzureFirewallPolicyNetworkRule` | yes |  |  |
| `spec.networkRuleCollections[].rules[].name` | `string` | yes |  |  |
| `spec.networkRuleCollections[].rules[].description` | `string` |  |  |  |
| `spec.networkRuleCollections[].rules[].protocols` | `[]enum` | yes |  |  |
| `spec.networkRuleCollections[].rules[].sourceAddresses` | `[]string` |  |  |  |
| `spec.networkRuleCollections[].rules[].sourceIpGroups` | `[]string \| valueFrom` |  |  | AzureIpGroup (`status.outputs.ip_group_id`) |
| `spec.networkRuleCollections[].rules[].destinationAddresses` | `[]string` |  |  |  |
| `spec.networkRuleCollections[].rules[].destinationIpGroups` | `[]string \| valueFrom` |  |  | AzureIpGroup (`status.outputs.ip_group_id`) |
| `spec.networkRuleCollections[].rules[].destinationFqdns` | `[]string` |  |  |  |
| `spec.networkRuleCollections[].rules[].destinationPorts` | `[]string` | yes |  |  |
| `spec.natRuleCollections` | `[]AzureFirewallPolicyNatRuleCollection` |  |  |  |
| `spec.natRuleCollections[].name` | `string` | yes |  |  |
| `spec.natRuleCollections[].priority` | `int32` | yes |  |  |
| `spec.natRuleCollections[].rules` | `[]AzureFirewallPolicyNatRule` | yes |  |  |
| `spec.natRuleCollections[].rules[].name` | `string` | yes |  |  |
| `spec.natRuleCollections[].rules[].description` | `string` |  |  |  |
| `spec.natRuleCollections[].rules[].protocols` | `[]enum` | yes |  |  |
| `spec.natRuleCollections[].rules[].sourceAddresses` | `[]string` |  |  |  |
| `spec.natRuleCollections[].rules[].sourceIpGroups` | `[]string \| valueFrom` |  |  | AzureIpGroup (`status.outputs.ip_group_id`) |
| `spec.natRuleCollections[].rules[].destinationAddress` | `string` |  |  |  |
| `spec.natRuleCollections[].rules[].destinationPorts` | `[]string` |  |  |  |
| `spec.natRuleCollections[].rules[].translatedAddress` | `string` |  |  |  |
| `spec.natRuleCollections[].rules[].translatedFqdn` | `string` |  |  |  |
| `spec.natRuleCollections[].rules[].translatedPort` | `int32` | yes |  |  |

## Field Details

### spec.firewallPolicyId

`string | valueFrom` · required

The firewall policy this group belongs to -- references an
AzureFirewallPolicy's ARM id. The policy is the group's parent: the
group deploys into it, and the policy's resource group and name are
derived from this id. Changing it replaces the group.

- references: AzureFirewallPolicy (`status.outputs.firewall_policy_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFirewallPolicy, name: <that resource's name>, fieldPath: status.outputs.firewall_policy_id}} -- a bare string does not parse

### spec.name

`string` · required

The group's name, unique within the policy. 2-80 characters; must
begin with a letter or number, end with a letter, number, or
underscore, and may contain only letters, numbers, underscores,
periods, or hyphens. Changing the name replaces the group. Name it
after the team or application whose rules it carries
("platform-baseline", "payments-app").

- rule: Rule collection group names are 2-80 characters, start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"maxLen":"80"}}

### spec.priority

`int32` · required

The group's evaluation priority among the policy's groups: 100-65000,
lower runs first. Leave gaps (100, 200, 300...) so future groups can
slot between existing ones without renumbering.

- rule: {"required":true,"int32":{"lte":65000,"gte":100}}

### spec.applicationRuleCollections

`[]AzureFirewallPolicyApplicationRuleCollection`

Application rule collections: allow/deny by HTTP(S)/MSSQL destination
-- FQDNs, URLs (Premium), FQDN tags, and web categories (Premium).
Application rules are evaluated LAST (after DNAT and network rules).

### spec.applicationRuleCollections[].name

`string` · required

The collection's name, unique within the group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.applicationRuleCollections[].priority

`int32` · required

The collection's priority among the group's collections of the same
type: 100-65000, lower runs first.

- rule: {"required":true,"int32":{"lte":65000,"gte":100}}

### spec.applicationRuleCollections[].action

`enum` · required

Whether matching traffic is allowed or denied.

- rule: {"required":true,"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_filter_action_unspecified` -- Not specified -- invalid; every filter collection declares its action.
- `ALLOW` -- Permit matching traffic.
- `DENY` -- Block matching traffic.

### spec.applicationRuleCollections[].rules

`[]AzureFirewallPolicyApplicationRule` · required

The collection's rules, evaluated in order.

- rule: {"repeated":{"minItems":"1"}}

### spec.applicationRuleCollections[].rules[].name

`string` · required

The rule's name, unique within the collection.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.applicationRuleCollections[].rules[].description

`string`

What this rule permits or blocks -- operator documentation.

### spec.applicationRuleCollections[].rules[].protocols

`[]AzureFirewallPolicyApplicationProtocol`

The L7 protocols (type + port) the rule matches, e.g. HTTPS on 443.
Required for FQDN/URL/category destinations.

### spec.applicationRuleCollections[].rules[].protocols[].type

`enum` · required

The protocol type.

- rule: {"required":true,"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_application_protocol_type_unspecified` -- Not specified -- invalid; every protocol entry declares its type.
- `HTTP` -- Plain HTTP.
- `HTTPS` -- HTTP over TLS. SNI-matched unless the rule terminates TLS.
- `MSSQL` -- The Microsoft SQL Server TDS protocol (used with Azure SQL through the firewall).

### spec.applicationRuleCollections[].rules[].protocols[].port

`int32` · required

The port the protocol runs on: 0-64000 (the provider's enforced
bound), e.g. 80 for HTTP, 443 for HTTPS, 1433 for MSSQL.

- rule: {"required":true,"int32":{"lte":64000,"gte":0}}

### spec.applicationRuleCollections[].rules[].sourceAddresses

`[]string`

Source IP addresses/CIDRs (or "*" for any).

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.applicationRuleCollections[].rules[].sourceIpGroups

`[]string | valueFrom`

Source IP Groups -- references to AzureIpGroup ARM ids, the reusable
alternative to literal source_addresses.

- references: AzureIpGroup (`status.outputs.ip_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureIpGroup, name: <that resource's name>, fieldPath: status.outputs.ip_group_id}} -- a bare string does not parse

### spec.applicationRuleCollections[].rules[].destinationAddresses

`[]string`

Destination IP addresses/CIDRs (or "*"). Use for L7 filtering to
fixed addresses; prefer FQDNs for anything DNS-named.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.applicationRuleCollections[].rules[].destinationFqdns

`[]string`

Destination fully-qualified domain names, e.g. "api.github.com" or
"*.ubuntu.com". FQDN filtering without TLS termination matches the
TLS SNI; with terminate_tls the firewall sees the decrypted host.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.applicationRuleCollections[].rules[].destinationUrls

`[]string`

Destination URLs including paths, e.g.
"github.com/planton/*". URL (path-level) filtering is Premium-only
and requires terminate_tls for HTTPS traffic -- without decryption
the firewall cannot see the path.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.applicationRuleCollections[].rules[].destinationFqdnTags

`[]string`

Azure-curated FQDN tags naming well-known service endpoint sets,
e.g. "WindowsUpdate", "AzureKubernetesService". The tag tracks
Microsoft's endpoint list so you do not have to.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.applicationRuleCollections[].rules[].terminateTls

`bool`

Decrypt, inspect, and re-encrypt matching TLS traffic using the
policy's tls_certificate (Premium only). Required for URL filtering
and full web-category classification of HTTPS traffic.

### spec.applicationRuleCollections[].rules[].webCategories

`[]string`

Azure web categories to match, e.g. "Gambling", "SocialNetworking"
(Premium only). Category classification of HTTPS traffic is
SNI-based unless terminate_tls is on.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.applicationRuleCollections[].rules[].httpHeaders

`[]AzureFirewallPolicyHttpHeader`

HTTP headers to ADD to matching requests before forwarding (e.g. a
token an upstream requires). Applies to HTTP and TLS-terminated
HTTPS traffic.

### spec.applicationRuleCollections[].rules[].httpHeaders[].name

`string` · required

The header name, e.g. "X-Forwarded-For-Org".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.applicationRuleCollections[].rules[].httpHeaders[].value

`string` · required

The header value.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.networkRuleCollections

`[]AzureFirewallPolicyNetworkRuleCollection`

Network rule collections: allow/deny by protocol, source, destination
address/FQDN, and port -- classic L3/L4 filtering. Evaluated after
DNAT rules and before application rules.

### spec.networkRuleCollections[].name

`string` · required

The collection's name, unique within the group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.networkRuleCollections[].priority

`int32` · required

The collection's priority among the group's collections of the same
type: 100-65000, lower runs first.

- rule: {"required":true,"int32":{"lte":65000,"gte":100}}

### spec.networkRuleCollections[].action

`enum` · required

Whether matching traffic is allowed or denied.

- rule: {"required":true,"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_filter_action_unspecified` -- Not specified -- invalid; every filter collection declares its action.
- `ALLOW` -- Permit matching traffic.
- `DENY` -- Block matching traffic.

### spec.networkRuleCollections[].rules

`[]AzureFirewallPolicyNetworkRule` · required

The collection's rules, evaluated in order.

- rule: {"repeated":{"minItems":"1"}}

### spec.networkRuleCollections[].rules[].name

`string` · required

The rule's name, unique within the collection.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.networkRuleCollections[].rules[].description

`string`

What this rule permits or blocks -- operator documentation.

### spec.networkRuleCollections[].rules[].protocols

`[]enum` · required

The protocols the rule matches: ANY, TCP, UDP, or ICMP.

- rule: {"repeated":{"minItems":"1","items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_rule_protocol_unspecified` -- Not specified -- invalid; every rule declares its protocols.
- `ANY` -- Match any protocol (network rules only).
- `TCP` -- Match TCP.
- `UDP` -- Match UDP.
- `ICMP` -- Match ICMP (network rules only).

### spec.networkRuleCollections[].rules[].sourceAddresses

`[]string`

Source IP addresses/CIDRs (or "*" for any).

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.networkRuleCollections[].rules[].sourceIpGroups

`[]string | valueFrom`

Source IP Groups -- references to AzureIpGroup ARM ids.

- references: AzureIpGroup (`status.outputs.ip_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureIpGroup, name: <that resource's name>, fieldPath: status.outputs.ip_group_id}} -- a bare string does not parse

### spec.networkRuleCollections[].rules[].destinationAddresses

`[]string`

Destination IP addresses/CIDRs, "*", or Azure service tags
(e.g. "AzureCloud.EastUS", "Storage").

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.networkRuleCollections[].rules[].destinationIpGroups

`[]string | valueFrom`

Destination IP Groups -- references to AzureIpGroup ARM ids.

- references: AzureIpGroup (`status.outputs.ip_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureIpGroup, name: <that resource's name>, fieldPath: status.outputs.ip_group_id}} -- a bare string does not parse

### spec.networkRuleCollections[].rules[].destinationFqdns

`[]string`

Destination FQDNs, e.g. "time.windows.com". Requires the policy's
DNS proxy so the firewall and clients resolve identically.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.networkRuleCollections[].rules[].destinationPorts

`[]string` · required

Destination ports: a port ("443"), a range ("1000-2000"), or "*".

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.natRuleCollections

`[]AzureFirewallPolicyNatRuleCollection`

DNAT rule collections: translate traffic arriving at the firewall's
public IP to an internal address/FQDN and port. Evaluated FIRST; a
matching DNAT rule also implicitly allows the translated flow.

### spec.natRuleCollections[].name

`string` · required

The collection's name, unique within the group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.natRuleCollections[].priority

`int32` · required

The collection's priority among the group's collections of the same
type: 100-65000, lower runs first.

- rule: {"required":true,"int32":{"lte":65000,"gte":100}}

### spec.natRuleCollections[].rules

`[]AzureFirewallPolicyNatRule` · required

The collection's rules, evaluated in order.

- rule: {"repeated":{"minItems":"1"}}
- rule: Set exactly one of translated_address or translated_fqdn -- the rule needs one translation target
- rule: DNAT rules support only TCP and UDP protocols

### spec.natRuleCollections[].rules[].name

`string` · required

The rule's name, unique within the collection.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.natRuleCollections[].rules[].description

`string`

What this rule translates -- operator documentation.

### spec.natRuleCollections[].rules[].protocols

`[]enum` · required

The protocols the rule matches. DNAT supports only TCP and UDP.

- rule: {"repeated":{"minItems":"1","items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_rule_protocol_unspecified` -- Not specified -- invalid; every rule declares its protocols.
- `ANY` -- Match any protocol (network rules only).
- `TCP` -- Match TCP.
- `UDP` -- Match UDP.
- `ICMP` -- Match ICMP (network rules only).

### spec.natRuleCollections[].rules[].sourceAddresses

`[]string`

Source IP addresses/CIDRs (or "*" for any).

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.natRuleCollections[].rules[].sourceIpGroups

`[]string | valueFrom`

Source IP Groups -- references to AzureIpGroup ARM ids.

- references: AzureIpGroup (`status.outputs.ip_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureIpGroup, name: <that resource's name>, fieldPath: status.outputs.ip_group_id}} -- a bare string does not parse

### spec.natRuleCollections[].rules[].destinationAddress

`string`

The destination the inbound traffic arrives at: one of the
firewall's public IP addresses. A single address or CIDR.

### spec.natRuleCollections[].rules[].destinationPorts

`[]string`

The destination port(s) the inbound traffic arrives on. ARM accepts
exactly ONE entry today (a port or a range, 1-64000; no wildcard) --
the list shape mirrors ARM's own, which may lift the cap.

- rule: {"repeated":{"maxItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.natRuleCollections[].rules[].translatedAddress

`string`

The internal IP address the traffic is translated TO. Exactly one of
translated_address / translated_fqdn.

### spec.natRuleCollections[].rules[].translatedFqdn

`string`

The internal FQDN the traffic is translated TO (requires the
policy's DNS proxy so the name resolves deterministically). Exactly
one of translated_address / translated_fqdn.

### spec.natRuleCollections[].rules[].translatedPort

`int32` · required

The internal port the traffic is translated to, 1-65535.

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFirewallPolicyRuleCollectionGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.rule_collection_group_id` | `string` | The Azure Resource Manager ID of the rule collection group. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/firewallPolicies/{policy}/ruleCollectionGroups/{name} |
| `status.outputs.rule_collection_group_name` | `string` | The name of the rule collection group resource. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.firewallPolicyId` | AzureFirewallPolicy | `status.outputs.firewall_policy_id` |
| `spec.applicationRuleCollections[].rules[].sourceIpGroups` | AzureIpGroup | `status.outputs.ip_group_id` |
| `spec.networkRuleCollections[].rules[].sourceIpGroups` | AzureIpGroup | `status.outputs.ip_group_id` |
| `spec.networkRuleCollections[].rules[].destinationIpGroups` | AzureIpGroup | `status.outputs.ip_group_id` |
| `spec.natRuleCollections[].rules[].sourceIpGroups` | AzureIpGroup | `status.outputs.ip_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
