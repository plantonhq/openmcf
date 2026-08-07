# AzureIpGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureIpGroupSpec** defines an Azure IP Group: a named, reusable set of
IP addresses and CIDR ranges that Azure Firewall and Firewall Policy
rules reference by ARM id instead of repeating literal address lists.
An IP Group lets a rule say "allow the branch offices to reach the
datacenter ranges" rather than enumerating twenty CIDRs inline in every
rule -- update the group once and every rule that references it follows,
so address management stops being copy-paste across rule collections.

The group itself is deliberately passive: it holds addresses and nothing
else. Consumption is declared from the rule's side -- a firewall policy
rule lists source/destination IP Groups, and an intrusion-detection
traffic bypass references them the same way. That inversion is what
makes the IP Group a stable, first-class composition anchor: it is
created once and referenced by many rules, each with its own lifecycle.

Azure caps how widely a single group can be shared (per Azure Firewall
limits: an IP Group can be referenced from at most 100 Firewall Policies
and carries at most 5,000 individual entries); design address sets
around intent ("branch-offices", "on-prem-datacenter", "scanners") so
rules read as policy statements.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureIpGroup
metadata:
  name: test-ip-group
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: branch-offices
  cidrs:
    - "203.0.113.7"
    - "10.10.0.0/16"
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.cidrs` | `[]string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region the group lives in, e.g. "eastus", "westeurope".
IP Groups are regional resources but can be referenced by firewall
policies in any region. Changing the region replaces the group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the IP Group is created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the IP Group, unique within the resource group. 1-80
characters following ARM's Microsoft.Network naming contract
(alphanumerics, underscores, periods, and hyphens; must start with a
letter or number and end with a letter, number, or underscore).
Changing the name replaces the group -- and every rule that referenced
it must be re-pointed -- so name it after the address set's intent
("branch-offices", "on-prem-ranges", "blocked-scanners").

- rule: IP Group names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.cidrs

`[]string`

The IP addresses and CIDR ranges in the group, e.g. "203.0.113.7" or
"10.10.0.0/16". Azure accepts single addresses and CIDR blocks (up to
5,000 entries per group). An empty group is legal -- useful as a
placeholder anchor that rules reference before the address plan is
final -- and entries update in place: adding or removing an address
immediately retargets every rule that references the group.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.tags

`map<string, string>`

Free-form tags applied to the group, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Tags are Azure's governance surface -- Azure
Policy enforces them and Cost Management groups by them.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureIpGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.ip_group_id` | `string` | The Azure Resource Manager ID of the IP Group. This is the composition seam: firewall policy rules (source_ip_groups / destination_ip_groups) and intrusion-detection traffic bypasses reference this to target the group's address set in a rule. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/ipGroups/{name} |
| `status.outputs.ip_group_name` | `string` | The name of the IP Group resource. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFirewallPolicy | `spec.intrusionDetection.trafficBypass[].sourceIpGroups` | `status.outputs.ip_group_id` |
| AzureFirewallPolicy | `spec.intrusionDetection.trafficBypass[].destinationIpGroups` | `status.outputs.ip_group_id` |
| AzureFirewallPolicyRuleCollectionGroup | `spec.applicationRuleCollections[].rules[].sourceIpGroups` | `status.outputs.ip_group_id` |
| AzureFirewallPolicyRuleCollectionGroup | `spec.networkRuleCollections[].rules[].sourceIpGroups` | `status.outputs.ip_group_id` |
| AzureFirewallPolicyRuleCollectionGroup | `spec.networkRuleCollections[].rules[].destinationIpGroups` | `status.outputs.ip_group_id` |
| AzureFirewallPolicyRuleCollectionGroup | `spec.natRuleCollections[].rules[].sourceIpGroups` | `status.outputs.ip_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
