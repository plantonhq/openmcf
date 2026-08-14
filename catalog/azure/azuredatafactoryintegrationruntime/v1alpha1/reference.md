# AzureDataFactoryIntegrationRuntime

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# The canonical deep-shape fixture for offline validation: the
# managed data-flow compute at full depth (every azure-variant field
# set), with a literal factory ID so the manifest validates
# standalone. The live lanes use the scenarios/ files instead.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDataFactoryIntegrationRuntime
metadata:
  name: dataflow-compute
  id: dataflow-compute
  org: planton-oss
  env: e2e
spec:
  dataFactoryId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.DataFactory/factories/app-df
  name: dataflow-compute
  description: The managed Spark compute mapping data flows run on, warm-pooled inside the managed virtual network.
  azure:
    region: eastus
    cleanupEnabled: false
    computeType: MemoryOptimized
    coreCount: 16
    timeToLiveMin: 15
    virtualNetworkEnabled: true
    interactiveAuthoringTimeToLiveInMinutes: 30
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.dataFactoryId` | `string \| valueFrom` | yes |  | AzureDataFactory (`status.outputs.data_factory_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.azure` | `AzureDataFactoryIntegrationRuntimeAzure` |  |  |  |
| `spec.azure.region` | `string` | yes |  |  |
| `spec.azure.cleanupEnabled` | `bool` |  | `true` |  |
| `spec.azure.computeType` | `string` |  |  |  |
| `spec.azure.coreCount` | `int32` |  |  |  |
| `spec.azure.timeToLiveMin` | `int32` |  |  |  |
| `spec.azure.interactiveAuthoringTimeToLiveInMinutes` | `int32` |  |  |  |
| `spec.azure.virtualNetworkEnabled` | `bool` |  |  |  |
| `spec.azureSsis` | `AzureDataFactoryIntegrationRuntimeAzureSsis` |  |  |  |
| `spec.azureSsis.region` | `string` | yes |  |  |
| `spec.azureSsis.nodeSize` | `string` | yes |  |  |
| `spec.azureSsis.numberOfNodes` | `int32` |  |  |  |
| `spec.azureSsis.maxParallelExecutionsPerNode` | `int32` |  |  |  |
| `spec.azureSsis.edition` | `string` |  |  |  |
| `spec.azureSsis.licenseType` | `string` |  |  |  |
| `spec.azureSsis.credentialName` | `string` |  |  |  |
| `spec.azureSsis.catalogInfo` | `AzureDataFactoryIntegrationRuntimeSsisCatalogInfo` |  |  |  |
| `spec.azureSsis.catalogInfo.serverEndpoint` | `string` | yes |  |  |
| `spec.azureSsis.catalogInfo.administratorLogin` | `string` |  |  |  |
| `spec.azureSsis.catalogInfo.administratorPassword` | `string` (sensitive) |  |  |  |
| `spec.azureSsis.catalogInfo.pricingTier` | `string` |  |  |  |
| `spec.azureSsis.catalogInfo.elasticPoolName` | `string` |  |  |  |
| `spec.azureSsis.catalogInfo.dualStandbyPairName` | `string` |  |  |  |
| `spec.azureSsis.customSetupScript` | `AzureDataFactoryIntegrationRuntimeSsisCustomSetupScript` |  |  |  |
| `spec.azureSsis.customSetupScript.blobContainerUri` | `string` | yes |  |  |
| `spec.azureSsis.customSetupScript.sasToken` | `string` (sensitive) | yes |  |  |
| `spec.azureSsis.expressCustomSetup` | `AzureDataFactoryIntegrationRuntimeSsisExpressCustomSetup` |  |  |  |
| `spec.azureSsis.expressCustomSetup.environment` | `map<string, string>` |  |  |  |
| `spec.azureSsis.expressCustomSetup.powershellVersion` | `string` |  |  |  |
| `spec.azureSsis.expressCustomSetup.commandKey` | `[]AzureDataFactoryIntegrationRuntimeSsisCommandKey` |  |  |  |
| `spec.azureSsis.expressCustomSetup.commandKey[].targetName` | `string` | yes |  |  |
| `spec.azureSsis.expressCustomSetup.commandKey[].userName` | `string` | yes |  |  |
| `spec.azureSsis.expressCustomSetup.commandKey[].password` | `string` (sensitive) |  |  |  |
| `spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword` | `AzureDataFactoryIntegrationRuntimeSsisKeyVaultSecretReference` |  |  |  |
| `spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.parameters` | `map<string, string>` |  |  |  |
| `spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.secretVersion` | `string` |  |  |  |
| `spec.azureSsis.expressCustomSetup.component` | `[]AzureDataFactoryIntegrationRuntimeSsisComponent` |  |  |  |
| `spec.azureSsis.expressCustomSetup.component[].name` | `string` | yes |  |  |
| `spec.azureSsis.expressCustomSetup.component[].license` | `string` (sensitive) |  |  |  |
| `spec.azureSsis.expressCustomSetup.component[].keyVaultLicense` | `AzureDataFactoryIntegrationRuntimeSsisKeyVaultSecretReference` |  |  |  |
| `spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.secretName` | `string` | yes |  |  |
| `spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.parameters` | `map<string, string>` |  |  |  |
| `spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.secretVersion` | `string` |  |  |  |
| `spec.azureSsis.expressVnetIntegration` | `AzureDataFactoryIntegrationRuntimeSsisExpressVnetIntegration` |  |  |  |
| `spec.azureSsis.expressVnetIntegration.subnetId` | `string \| valueFrom` | yes |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.azureSsis.vnetIntegration` | `AzureDataFactoryIntegrationRuntimeSsisVnetIntegration` |  |  |  |
| `spec.azureSsis.vnetIntegration.vnetId` | `string \| valueFrom` |  |  | AzureVirtualNetwork (`status.outputs.virtual_network_id`) |
| `spec.azureSsis.vnetIntegration.subnetId` | `string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.azureSsis.vnetIntegration.subnetName` | `string` |  |  |  |
| `spec.azureSsis.vnetIntegration.publicIps` | `[]string \| valueFrom` |  |  | AzurePublicIp (`status.outputs.public_ip_id`) |
| `spec.azureSsis.packageStore` | `[]AzureDataFactoryIntegrationRuntimeSsisPackageStore` |  |  |  |
| `spec.azureSsis.packageStore[].name` | `string` | yes |  |  |
| `spec.azureSsis.packageStore[].linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureSsis.copyComputeScale` | `AzureDataFactoryIntegrationRuntimeSsisCopyComputeScale` |  |  |  |
| `spec.azureSsis.copyComputeScale.dataIntegrationUnit` | `int32` |  |  |  |
| `spec.azureSsis.copyComputeScale.timeToLive` | `int32` |  |  |  |
| `spec.azureSsis.pipelineExternalComputeScale` | `AzureDataFactoryIntegrationRuntimeSsisPipelineExternalComputeScale` |  |  |  |
| `spec.azureSsis.pipelineExternalComputeScale.numberOfExternalNodes` | `int32` |  |  |  |
| `spec.azureSsis.pipelineExternalComputeScale.numberOfPipelineNodes` | `int32` |  |  |  |
| `spec.azureSsis.pipelineExternalComputeScale.timeToLive` | `int32` |  |  |  |
| `spec.azureSsis.proxy` | `AzureDataFactoryIntegrationRuntimeSsisProxy` |  |  |  |
| `spec.azureSsis.proxy.selfHostedIntegrationRuntimeName` | `string \| valueFrom` | yes |  | AzureDataFactoryIntegrationRuntime (`status.outputs.integration_runtime_name`) |
| `spec.azureSsis.proxy.stagingStorageLinkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureSsis.proxy.path` | `string` |  |  |  |
| `spec.selfHosted` | `AzureDataFactoryIntegrationRuntimeSelfHosted` |  |  |  |
| `spec.selfHosted.rbacAuthorization` | `AzureDataFactoryIntegrationRuntimeSelfHostedRbacAuthorization` |  |  |  |
| `spec.selfHosted.rbacAuthorization.resourceId` | `string \| valueFrom` | yes |  | AzureDataFactoryIntegrationRuntime (`status.outputs.integration_runtime_id`) |
| `spec.selfHosted.selfContainedInteractiveAuthoringEnabled` | `bool` |  |  |  |

## Field Details

### spec.dataFactoryId

`string | valueFrom` · required

- references: AzureDataFactory (`status.outputs.data_factory_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactory, name: <that resource's name>, fieldPath: status.outputs.data_factory_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: {"required":true}

### spec.description

`string`

### spec.azure

`AzureDataFactoryIntegrationRuntimeAzure`

- rule: interactive_authoring_time_to_live_in_minutes requires virtual_network_enabled to be true

### spec.azure.region

`string` · required

- rule: {"required":true}

### spec.azure.cleanupEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.azure.computeType

`string`

- rule: {"string":{"in":["","General","ComputeOptimized","MemoryOptimized"]}}

### spec.azure.coreCount

`int32`

- rule: core_count must be one of 8, 16, 32, 48, 80, 144, or 272

### spec.azure.timeToLiveMin

`int32`

- rule: {"int32":{"gte":0}}

### spec.azure.interactiveAuthoringTimeToLiveInMinutes

`int32`

- rule: interactive_authoring_time_to_live_in_minutes must be one of 10, 30, 60, or 120

### spec.azure.virtualNetworkEnabled

`bool`

### spec.azureSsis

`AzureDataFactoryIntegrationRuntimeAzureSsis`

### spec.azureSsis.region

`string` · required

- rule: {"required":true}

### spec.azureSsis.nodeSize

`string` · required

- rule: {"required":true,"string":{"in":["Standard_D2_v3","Standard_D4_v3","Standard_D8_v3","Standard_D16_v3","Standard_D32_v3","Standard_D64_v3","Standard_E2_v3","Standard_E4_v3","Standard_E8_v3","Standard_E16_v3","Standard_E32_v3","Standard_E64_v3","Standard_D1_v2","Standard_D2_v2","Standard_D3_v2","Standard_D4_v2","Standard_A4_v2","Standard_A8_v2"]}}

### spec.azureSsis.numberOfNodes

`int32`

- rule: number_of_nodes must be between 1 and 10

### spec.azureSsis.maxParallelExecutionsPerNode

`int32`

- rule: max_parallel_executions_per_node must be between 1 and 16

### spec.azureSsis.edition

`string`

- rule: {"string":{"in":["","Standard","Enterprise"]}}

### spec.azureSsis.licenseType

`string`

- rule: {"string":{"in":["","LicenseIncluded","BasePrice"]}}

### spec.azureSsis.credentialName

`string`

### spec.azureSsis.catalogInfo

`AzureDataFactoryIntegrationRuntimeSsisCatalogInfo`

- rule: pricing_tier and elastic_pool_name are mutually exclusive -- SSISDB lands either on a tier or in a pool

### spec.azureSsis.catalogInfo.serverEndpoint

`string` · required

- rule: {"required":true}

### spec.azureSsis.catalogInfo.administratorLogin

`string`

### spec.azureSsis.catalogInfo.administratorPassword

`string` · sensitive

### spec.azureSsis.catalogInfo.pricingTier

`string`

- rule: {"string":{"in":["","Basic","S0","S1","S2","S3","S4","S6","S7","S9","S12","P1","P2","P4","P6","P11","P15","GP_S_Gen5_1","GP_S_Gen5_2","GP_S_Gen5_4","GP_S_Gen5_6","GP_S_Gen5_8","GP_S_Gen5_10","GP_S_Gen5_12","GP_S_Gen5_14","GP_S_Gen5_16","GP_S_Gen5_18","GP_S_Gen5_20","GP_S_Gen5_24","GP_S_Gen5_32","GP_S_Gen5_40","GP_Gen5_2","GP_Gen5_4","GP_Gen5_6","GP_Gen5_8","GP_Gen5_10","GP_Gen5_12","GP_Gen5_14","GP_Gen5_16","GP_Gen5_18","GP_Gen5_20","GP_Gen5_24","GP_Gen5_32","GP_Gen5_40","GP_Gen5_80","BC_Gen5_2","BC_Gen5_4","BC_Gen5_6","BC_Gen5_8","BC_Gen5_10","BC_Gen5_12","BC_Gen5_14","BC_Gen5_16","BC_Gen5_18","BC_Gen5_20","BC_Gen5_24","BC_Gen5_32","BC_Gen5_40","BC_Gen5_80","HS_Gen5_2","HS_Gen5_4","HS_Gen5_6","HS_Gen5_8","HS_Gen5_10","HS_Gen5_12","HS_Gen5_14","HS_Gen5_16","HS_Gen5_18","HS_Gen5_20","HS_Gen5_24","HS_Gen5_32","HS_Gen5_40","HS_Gen5_80"]}}

### spec.azureSsis.catalogInfo.elasticPoolName

`string`

### spec.azureSsis.catalogInfo.dualStandbyPairName

`string`

### spec.azureSsis.customSetupScript

`AzureDataFactoryIntegrationRuntimeSsisCustomSetupScript`

### spec.azureSsis.customSetupScript.blobContainerUri

`string` · required

- rule: {"required":true}

### spec.azureSsis.customSetupScript.sasToken

`string` · required · sensitive

- rule: {"required":true}

### spec.azureSsis.expressCustomSetup

`AzureDataFactoryIntegrationRuntimeSsisExpressCustomSetup`

- rule: Declare at least one of environment, powershell_version, command_key, or component

### spec.azureSsis.expressCustomSetup.environment

`map<string, string>`

### spec.azureSsis.expressCustomSetup.powershellVersion

`string`

### spec.azureSsis.expressCustomSetup.commandKey

`[]AzureDataFactoryIntegrationRuntimeSsisCommandKey`

### spec.azureSsis.expressCustomSetup.commandKey[].targetName

`string` · required

- rule: {"required":true}

### spec.azureSsis.expressCustomSetup.commandKey[].userName

`string` · required

- rule: {"required":true}

### spec.azureSsis.expressCustomSetup.commandKey[].password

`string` · sensitive

### spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword

`AzureDataFactoryIntegrationRuntimeSsisKeyVaultSecretReference`

### spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.parameters

`map<string, string>`

### spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.secretVersion

`string`

### spec.azureSsis.expressCustomSetup.component

`[]AzureDataFactoryIntegrationRuntimeSsisComponent`

### spec.azureSsis.expressCustomSetup.component[].name

`string` · required

- rule: {"required":true}

### spec.azureSsis.expressCustomSetup.component[].license

`string` · sensitive

### spec.azureSsis.expressCustomSetup.component[].keyVaultLicense

`AzureDataFactoryIntegrationRuntimeSsisKeyVaultSecretReference`

### spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.secretName

`string` · required

- rule: {"required":true}

### spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.parameters

`map<string, string>`

### spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.secretVersion

`string`

### spec.azureSsis.expressVnetIntegration

`AzureDataFactoryIntegrationRuntimeSsisExpressVnetIntegration`

### spec.azureSsis.expressVnetIntegration.subnetId

`string | valueFrom` · required

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.azureSsis.vnetIntegration

`AzureDataFactoryIntegrationRuntimeSsisVnetIntegration`

- rule: Set exactly one of vnet_id (with subnet_name) or subnet_id
- rule: subnet_name is required with vnet_id (and is only meaningful there)
- rule: public_ips takes exactly two addresses (or none, letting Azure assign them)

### spec.azureSsis.vnetIntegration.vnetId

`string | valueFrom`

- references: AzureVirtualNetwork (`status.outputs.virtual_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualNetwork, name: <that resource's name>, fieldPath: status.outputs.virtual_network_id}} -- a bare string does not parse

### spec.azureSsis.vnetIntegration.subnetId

`string | valueFrom`

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.azureSsis.vnetIntegration.subnetName

`string`

### spec.azureSsis.vnetIntegration.publicIps

`[]string | valueFrom`

- references: AzurePublicIp (`status.outputs.public_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIp, name: <that resource's name>, fieldPath: status.outputs.public_ip_id}} -- a bare string does not parse

### spec.azureSsis.packageStore

`[]AzureDataFactoryIntegrationRuntimeSsisPackageStore`

### spec.azureSsis.packageStore[].name

`string` · required

- rule: {"required":true}

### spec.azureSsis.packageStore[].linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureSsis.copyComputeScale

`AzureDataFactoryIntegrationRuntimeSsisCopyComputeScale`

### spec.azureSsis.copyComputeScale.dataIntegrationUnit

`int32`

- rule: data_integration_unit must be a multiple of 4 between 4 and 256

### spec.azureSsis.copyComputeScale.timeToLive

`int32`

- rule: time_to_live must be at least 5 minutes

### spec.azureSsis.pipelineExternalComputeScale

`AzureDataFactoryIntegrationRuntimeSsisPipelineExternalComputeScale`

### spec.azureSsis.pipelineExternalComputeScale.numberOfExternalNodes

`int32`

- rule: number_of_external_nodes must be between 1 and 10

### spec.azureSsis.pipelineExternalComputeScale.numberOfPipelineNodes

`int32`

- rule: number_of_pipeline_nodes must be between 1 and 10

### spec.azureSsis.pipelineExternalComputeScale.timeToLive

`int32`

- rule: time_to_live must be at least 5 minutes

### spec.azureSsis.proxy

`AzureDataFactoryIntegrationRuntimeSsisProxy`

### spec.azureSsis.proxy.selfHostedIntegrationRuntimeName

`string | valueFrom` · required

- references: AzureDataFactoryIntegrationRuntime (`status.outputs.integration_runtime_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryIntegrationRuntime, name: <that resource's name>, fieldPath: status.outputs.integration_runtime_name}} -- a bare string does not parse

### spec.azureSsis.proxy.stagingStorageLinkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureSsis.proxy.path

`string`

### spec.selfHosted

`AzureDataFactoryIntegrationRuntimeSelfHosted`

### spec.selfHosted.rbacAuthorization

`AzureDataFactoryIntegrationRuntimeSelfHostedRbacAuthorization`

### spec.selfHosted.rbacAuthorization.resourceId

`string | valueFrom` · required

- references: AzureDataFactoryIntegrationRuntime (`status.outputs.integration_runtime_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryIntegrationRuntime, name: <that resource's name>, fieldPath: status.outputs.integration_runtime_id}} -- a bare string does not parse

### spec.selfHosted.selfContainedInteractiveAuthoringEnabled

`bool`

## Validation Rules

- `azure_data_factory_integration_runtime_exactly_one_variant`: Set exactly one integration runtime variant block -- the variant determines the engine flavor
- `azure_data_factory_integration_runtime_name_format_managed`: Managed runtime names need at least 3 characters -- letters, numbers, and dashes only, starting and ending with a letter or number, no consecutive dashes
- `azure_data_factory_integration_runtime_name_format_self_hosted`: Self-hosted runtime names use letters, numbers, and dashes, starting and ending with a letter or number, with no consecutive dashes

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureDataFactoryIntegrationRuntime, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.integration_runtime_id` | `string` |  |
| `status.outputs.integration_runtime_name` | `string` |  |
| `status.outputs.primary_authorization_key` | `string` |  |
| `status.outputs.secondary_authorization_key` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.dataFactoryId` | AzureDataFactory | `status.outputs.data_factory_id` |
| `spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureSsis.expressVnetIntegration.subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.azureSsis.vnetIntegration.vnetId` | AzureVirtualNetwork | `status.outputs.virtual_network_id` |
| `spec.azureSsis.vnetIntegration.subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.azureSsis.vnetIntegration.publicIps` | AzurePublicIp | `status.outputs.public_ip_id` |
| `spec.azureSsis.packageStore[].linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureSsis.proxy.selfHostedIntegrationRuntimeName` | AzureDataFactoryIntegrationRuntime | `status.outputs.integration_runtime_name` |
| `spec.azureSsis.proxy.stagingStorageLinkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.selfHosted.rbacAuthorization.resourceId` | AzureDataFactoryIntegrationRuntime | `status.outputs.integration_runtime_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureDataFactoryIntegrationRuntime | `spec.azureSsis.proxy.selfHostedIntegrationRuntimeName` | `status.outputs.integration_runtime_name` |
| AzureDataFactoryIntegrationRuntime | `spec.selfHosted.rbacAuthorization.resourceId` | `status.outputs.integration_runtime_id` |
| AzureDataFactoryLinkedService | `spec.integrationRuntimeName` | `status.outputs.integration_runtime_name` |

## See Also

- [Overview](../README.md)
