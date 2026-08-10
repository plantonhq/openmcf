# AzureResourceGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureResourceGroupSpec** defines the configuration for creating an Azure Resource Group.

Azure Resource Groups are the fundamental organizational unit for Azure resources.
Every Azure resource must belong to a resource group. Resource groups provide:
- **Lifecycle management**: Delete a resource group to delete all contained resources
- **RBAC boundaries**: Assign permissions at the resource group level
- **Cost management**: Track and allocate costs per resource group
- **Tagging**: Apply tags for organizational and billing purposes

This is an intentionally minimal spec. Resource groups are simple containers --
they need only a name and a region. All complexity lives in the resources deployed
into them.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureResourceGroup
metadata:
  name: test-rg
  org: test-org
  env: dev
spec:
  name: test-platform-rg
  region: eastus
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |

## Field Details

### spec.name

`string` · required

The name of the resource group.
Must be unique within the Azure subscription.
Allowed characters: alphanumeric, underscores, hyphens, periods, and parentheses.
Cannot end with a period. Maximum 90 characters.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"90"}}

### spec.region

`string` · required

The Azure region where the resource group will be created.
Examples: "eastus", "westus2", "westeurope", "southeastasia".
Note: Resources within the resource group can be in different regions,
but the resource group's region determines where its metadata is stored.

- rule: {"required":true,"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureResourceGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.resource_group_id` | `string` | The Azure Resource Manager ID of the resource group. Format: /subscriptions/{subscription-id}/resourceGroups/{resource-group-name} |
| `status.outputs.resource_group_name` | `string` | The name of the resource group. This is the primary output referenced by downstream Azure resources via StringValueOrRef resource_group with field_path "status.outputs.resource_group_name". |
| `status.outputs.region` | `string` | The Azure region where the resource group was created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureAiFoundry | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureAksCluster | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureApplicationGateway | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureApplicationInsights | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureApplicationInsightsStandardWebTest | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureApplicationSecurityGroup | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureCognitiveAccount | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureContainerApp | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureContainerAppEnvironment | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureContainerAppJob | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureContainerRegistry | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureCosmosdbAccount | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureDiskEncryptionSet | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureDnsRecord | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureDnsZone | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureEventHubCluster | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureEventHubNamespace | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureExpressRouteCircuit | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureExpressRouteCircuitPeering | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureExpressRouteGateway | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureExpressRoutePort | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureFirewall | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureFirewallPolicy | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureFrontDoorFirewallPolicy | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureFrontDoorProfile | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureFunctionApp | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureIpGroup | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureKeyVault | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureLinuxWebApp | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureLoadBalancer | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureLocalNetworkGateway | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureLogAnalyticsWorkspace | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureMachineLearningWorkspace | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureManagedDisk | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureManagedRedis | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureMonitorActionGroup | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureMonitorActivityLogAlert | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureMonitorActivityLogAlert | `spec.scopes` | `status.outputs.resource_group_id` |
| AzureMonitorMetricAlert | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureMonitorScheduledQueryAlert | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureMssqlServer | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureMysqlFlexibleServer | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureNatGateway | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureNetworkInterface | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureNetworkSecurityGroup | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzurePointToSiteVpnGateway | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzurePostgresqlFlexibleServer | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzurePrivateDnsZone | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzurePrivateEndpoint | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzurePrivateLinkService | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzurePublicIp | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzurePublicIpPrefix | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureRedisCache | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureRoleAssignment | `spec.scope` | `status.outputs.resource_group_id` |
| AzureRoleDefinition | `spec.scope` | `status.outputs.resource_group_id` |
| AzureRoleDefinition | `spec.assignableScopes` | `status.outputs.resource_group_id` |
| AzureRouteTable | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureSearchService | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureServiceBusNamespace | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureServicePlan | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureStorageAccount | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureUserAssignedIdentity | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVirtualHub | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVirtualMachine | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVirtualMachineScaleSet | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVirtualNetwork | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVirtualNetworkGateway | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVirtualNetworkGatewayConnection | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVirtualWan | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVpnGateway | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVpnServerConfiguration | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureVpnSite | `spec.resourceGroup` | `status.outputs.resource_group_name` |
| AzureWebApplicationFirewallPolicy | `spec.resourceGroup` | `status.outputs.resource_group_name` |

## See Also

- [Overview](../README.md)
