# AzureApplicationSecurityGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureApplicationSecurityGroupSpec** defines an Azure Application
Security Group (ASG): a named, logical grouping of network interfaces
that network security group rules can target by name instead of by IP
address. An ASG lets a security policy say "allow the web tier to reach
the app tier" rather than "allow 10.0.1.0/24 to reach 10.0.2.0/24" --
the rule follows the workload as instances scale in and out and change
addresses, so micro-segmentation stops being tied to a fragile CIDR
plan.

The ASG itself is deliberately empty: it holds no members. Membership is
declared from the member's side -- a network interface (or a VM scale
set's IP configuration) lists the ASGs it joins, and an NSG security
rule references source/destination ASGs. That inversion is what makes
the ASG a stable, first-class composition anchor: it is created once and
referenced by many members and rules, each with its own lifecycle.

Everything except tags is fixed at creation; changing the name or region
replaces the group (and every rule and NIC that referenced it must be
re-pointed), so name it after the workload role it represents ("web",
"app-tier", "db").

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureApplicationSecurityGroup
metadata:
  name: test-application-security-group
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: web-tier
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region the group lives in, e.g. "eastus", "westeurope".
An ASG can only be referenced by network interfaces in the same
region. Changing the region replaces the group. Fixed at creation.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the ASG is created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the ASG, unique within the resource group. 1-80
characters (alphanumerics, underscores, periods, and hyphens; must
start with a letter or number and end with a letter, number, or
underscore). Changing the name replaces the group. Name it after the
workload role it represents ("web", "app-tier", "db") so security
rules read as intent.

- rule: Application security group names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.tags

`map<string, string>`

Free-form tags applied to the group, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Tags are Azure's governance surface -- Azure
Policy enforces them and Cost Management groups by them. The only
thing on an ASG that updates in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureApplicationSecurityGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.application_security_group_id` | `string` | The Azure Resource Manager ID of the application security group. This is the composition seam: network interfaces (application_security_group_ids), VM scale set IP configurations, and NSG security rules (source/destination_application_security_group_ids) reference this to declare membership or target the group in a rule. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/applicationSecurityGroups/{name} |
| `status.outputs.application_security_group_name` | `string` | The name of the application security group resource. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureNetworkInterface | `spec.applicationSecurityGroupIds` | `status.outputs.application_security_group_id` |
| AzureNetworkSecurityGroup | `spec.securityRules[].sourceApplicationSecurityGroupIds` | `status.outputs.application_security_group_id` |
| AzureNetworkSecurityGroup | `spec.securityRules[].destinationApplicationSecurityGroupIds` | `status.outputs.application_security_group_id` |
| AzurePrivateEndpoint | `spec.applicationSecurityGroupIds` | `status.outputs.application_security_group_id` |
| AzureVirtualMachineScaleSet | `spec.networkInterfaces[].ipConfigurations[].applicationSecurityGroupIds` | `status.outputs.application_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
