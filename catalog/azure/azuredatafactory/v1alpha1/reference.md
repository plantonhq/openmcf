# AzureDataFactory

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a factory with
# the combined identity, a GitHub repository binding, global
# parameters across the type vocabulary, customer-managed-key
# encryption by reference form, both credential flavors, the managed
# virtual network, and a managed private endpoint on each arm
# (subresource and Private Link Service fqdns). References are literal
# values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDataFactory
metadata:
  name: test-data-factory
  id: test-data-factory
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: test-org-data-factory
  region: eastus
  identity:
    type: SYSTEM_AND_USER_ASSIGNED
    identityIds:
      - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-uai
  githubConfiguration:
    accountName: test-org
    branchName: main
    repositoryName: data-pipelines
    rootFolder: /
    publishingEnabled: true
  globalParameters:
    - name: environment
      type: String
      value: test
    - name: retries
      type: Int
      value: "3"
    - name: regions
      type: Array
      value: '["eastus","westus"]'
  managedVirtualNetworkEnabled: true
  publicNetworkEnabled: false
  customerManagedKey:
    keyVaultKeyId:
      value: https://test-vault.vault.azure.net/keys/df-cmk
    userAssignedIdentityId:
      value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-uai
  userManagedIdentityCredentials:
    - name: etl-identity
      identityId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-uai
      description: identity the ETL linked services run as
  servicePrincipalCredentials:
    - name: legacy-sp
      tenantId: 11111111-2222-3333-4444-555555555555
      servicePrincipalId: 22222222-3333-4444-5555-666666666666
      servicePrincipalKey:
        linkedServiceName: keyvault-ls
        secretName: legacy-sp-key
  managedPrivateEndpoints:
    - name: datalake-blob
      targetResourceId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/testdata
      subresourceName: blob
    - name: partner-pls
      targetResourceId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/privateLinkServices/partner-pls
      fqdns:
        - internal.partner.example
  tags:
    costCenter: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.identity` | `AzureDataFactoryIdentity` |  |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.githubConfiguration` | `AzureDataFactoryGithubConfiguration` |  |  |  |
| `spec.githubConfiguration.accountName` | `string` | yes |  |  |
| `spec.githubConfiguration.branchName` | `string` | yes |  |  |
| `spec.githubConfiguration.gitUrl` | `string` |  |  |  |
| `spec.githubConfiguration.repositoryName` | `string` | yes |  |  |
| `spec.githubConfiguration.rootFolder` | `string` | yes |  |  |
| `spec.githubConfiguration.publishingEnabled` | `bool` |  | `true` |  |
| `spec.vstsConfiguration` | `AzureDataFactoryVstsConfiguration` |  |  |  |
| `spec.vstsConfiguration.accountName` | `string` | yes |  |  |
| `spec.vstsConfiguration.branchName` | `string` | yes |  |  |
| `spec.vstsConfiguration.projectName` | `string` | yes |  |  |
| `spec.vstsConfiguration.repositoryName` | `string` | yes |  |  |
| `spec.vstsConfiguration.rootFolder` | `string` | yes |  |  |
| `spec.vstsConfiguration.tenantId` | `string` | yes |  |  |
| `spec.vstsConfiguration.publishingEnabled` | `bool` |  | `true` |  |
| `spec.globalParameters` | `[]AzureDataFactoryGlobalParameter` |  |  |  |
| `spec.globalParameters[].name` | `string` | yes |  |  |
| `spec.globalParameters[].type` | `string` | yes |  |  |
| `spec.globalParameters[].value` | `string` | yes |  |  |
| `spec.managedVirtualNetworkEnabled` | `bool` |  |  |  |
| `spec.publicNetworkEnabled` | `bool` |  | `true` |  |
| `spec.purviewId` | `string \| valueFrom` |  |  |  |
| `spec.customerManagedKey` | `AzureDataFactoryCustomerManagedKey` |  |  |  |
| `spec.customerManagedKey.keyVaultKeyId` | `string \| valueFrom` | yes |  | AzureKeyVaultKey (`status.outputs.versionless_id`) |
| `spec.customerManagedKey.userAssignedIdentityId` | `string \| valueFrom` | yes |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.userManagedIdentityCredentials` | `[]AzureDataFactoryUserManagedIdentityCredential` |  |  |  |
| `spec.userManagedIdentityCredentials[].name` | `string` | yes |  |  |
| `spec.userManagedIdentityCredentials[].identityId` | `string \| valueFrom` | yes |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.userManagedIdentityCredentials[].description` | `string` |  |  |  |
| `spec.userManagedIdentityCredentials[].annotations` | `[]string` |  |  |  |
| `spec.servicePrincipalCredentials` | `[]AzureDataFactoryServicePrincipalCredential` |  |  |  |
| `spec.servicePrincipalCredentials[].name` | `string` | yes |  |  |
| `spec.servicePrincipalCredentials[].tenantId` | `string` | yes |  |  |
| `spec.servicePrincipalCredentials[].servicePrincipalId` | `string` | yes |  |  |
| `spec.servicePrincipalCredentials[].servicePrincipalKey` | `AzureDataFactoryServicePrincipalKey` |  |  |  |
| `spec.servicePrincipalCredentials[].servicePrincipalKey.linkedServiceName` | `string` | yes |  |  |
| `spec.servicePrincipalCredentials[].servicePrincipalKey.secretName` | `string` | yes |  |  |
| `spec.servicePrincipalCredentials[].servicePrincipalKey.secretVersion` | `string` |  |  |  |
| `spec.servicePrincipalCredentials[].description` | `string` |  |  |  |
| `spec.servicePrincipalCredentials[].annotations` | `[]string` |  |  |  |
| `spec.managedPrivateEndpoints` | `[]AzureDataFactoryManagedPrivateEndpoint` |  |  |  |
| `spec.managedPrivateEndpoints[].name` | `string` | yes |  |  |
| `spec.managedPrivateEndpoints[].targetResourceId` | `string \| valueFrom` | yes |  |  |
| `spec.managedPrivateEndpoints[].subresourceName` | `string` |  |  |  |
| `spec.managedPrivateEndpoints[].fqdns` | `[]string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Factory names must be 3-63 characters of letters, numbers, and hyphens, starting and ending with a letter or number
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.identity

`AzureDataFactoryIdentity`

- rule: identity_ids is required for USER_ASSIGNED and SYSTEM_AND_USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_data_factory_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`
- `SYSTEM_AND_USER_ASSIGNED`

### spec.identity.identityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.githubConfiguration

`AzureDataFactoryGithubConfiguration`

### spec.githubConfiguration.accountName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.githubConfiguration.branchName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.githubConfiguration.gitUrl

`string`

### spec.githubConfiguration.repositoryName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.githubConfiguration.rootFolder

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.githubConfiguration.publishingEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.vstsConfiguration

`AzureDataFactoryVstsConfiguration`

### spec.vstsConfiguration.accountName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vstsConfiguration.branchName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vstsConfiguration.projectName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vstsConfiguration.repositoryName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vstsConfiguration.rootFolder

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vstsConfiguration.tenantId

`string` · required

- rule: {"required":true,"string":{"uuid":true}}

### spec.vstsConfiguration.publishingEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.globalParameters

`[]AzureDataFactoryGlobalParameter`

### spec.globalParameters[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.globalParameters[].type

`string` · required

- rule: {"required":true,"string":{"in":["Array","Bool","Float","Int","Object","String"]}}

### spec.globalParameters[].value

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managedVirtualNetworkEnabled

`bool` · optional (explicit presence)

### spec.publicNetworkEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.purviewId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.customerManagedKey

`AzureDataFactoryCustomerManagedKey`

### spec.customerManagedKey.keyVaultKeyId

`string | valueFrom` · required

- references: AzureKeyVaultKey (`status.outputs.versionless_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultKey, name: <that resource's name>, fieldPath: status.outputs.versionless_id}} -- a bare string does not parse

### spec.customerManagedKey.userAssignedIdentityId

`string | valueFrom` · required

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.userManagedIdentityCredentials

`[]AzureDataFactoryUserManagedIdentityCredential`

### spec.userManagedIdentityCredentials[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.userManagedIdentityCredentials[].identityId

`string | valueFrom` · required

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.userManagedIdentityCredentials[].description

`string`

### spec.userManagedIdentityCredentials[].annotations

`[]string`

### spec.servicePrincipalCredentials

`[]AzureDataFactoryServicePrincipalCredential`

### spec.servicePrincipalCredentials[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.servicePrincipalCredentials[].tenantId

`string` · required

- rule: {"required":true,"string":{"uuid":true}}

### spec.servicePrincipalCredentials[].servicePrincipalId

`string` · required

- rule: {"required":true,"string":{"uuid":true}}

### spec.servicePrincipalCredentials[].servicePrincipalKey

`AzureDataFactoryServicePrincipalKey`

### spec.servicePrincipalCredentials[].servicePrincipalKey.linkedServiceName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.servicePrincipalCredentials[].servicePrincipalKey.secretName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.servicePrincipalCredentials[].servicePrincipalKey.secretVersion

`string`

### spec.servicePrincipalCredentials[].description

`string`

### spec.servicePrincipalCredentials[].annotations

`[]string`

### spec.managedPrivateEndpoints

`[]AzureDataFactoryManagedPrivateEndpoint`

- rule: Exactly one of subresource_name (regular ARM targets) and fqdns (Private Link Service targets) must be set

### spec.managedPrivateEndpoints[].name

`string` · required

- rule: Endpoint names must be 2-80 characters of letters, numbers, dots, hyphens, and underscores, starting with a letter or number and ending with a letter, number, or underscore
- rule: {"required":true}

### spec.managedPrivateEndpoints[].targetResourceId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.managedPrivateEndpoints[].subresourceName

`string`

- rule: subresource_name must be 3-63 characters of letters, numbers, dots, hyphens, and underscores, starting and ending with a letter or number

### spec.managedPrivateEndpoints[].fqdns

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.tags

`map<string, string>`

## Validation Rules

- `data_factory_repo_configuration_exclusive`: At most one of github_configuration and vsts_configuration may be set -- a factory binds to one repository
- `data_factory_cmk_requires_user_assigned_identity`: customer_managed_key requires the identity block with a user-assigned flavor (USER_ASSIGNED or SYSTEM_AND_USER_ASSIGNED)
- `data_factory_global_parameter_names_unique`: global_parameters names must be unique -- the name is the parameter's identity
- `data_factory_credential_names_unique`: credential names must be unique across user_managed_identity_credentials and service_principal_credentials -- credentials share one namespace under the factory
- `data_factory_managed_private_endpoint_names_unique`: managed_private_endpoints names must be unique -- the name is the endpoint's identity
- `data_factory_managed_private_endpoints_require_managed_vnet`: managed_private_endpoints require managed_virtual_network_enabled to be true -- the endpoints live inside the managed virtual network

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureDataFactory, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.data_factory_id` | `string` |  |
| `status.outputs.data_factory_name` | `string` |  |
| `status.outputs.identity_principal_id` | `string` |  |
| `status.outputs.credential_ids` | `map<string, string>` |  |
| `status.outputs.managed_private_endpoint_ids` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.customerManagedKey.keyVaultKeyId` | AzureKeyVaultKey | `status.outputs.versionless_id` |
| `spec.customerManagedKey.userAssignedIdentityId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.userManagedIdentityCredentials[].identityId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureDataFactoryDataFlow | `spec.dataFactoryId` | `status.outputs.data_factory_id` |
| AzureDataFactoryDataset | `spec.dataFactoryId` | `status.outputs.data_factory_id` |
| AzureDataFactoryIntegrationRuntime | `spec.dataFactoryId` | `status.outputs.data_factory_id` |
| AzureDataFactoryLinkedService | `spec.dataFactoryId` | `status.outputs.data_factory_id` |
| AzureDataFactoryPipeline | `spec.dataFactoryId` | `status.outputs.data_factory_id` |
| AzureDataFactoryTrigger | `spec.dataFactoryId` | `status.outputs.data_factory_id` |

## See Also

- [Overview](../README.md)
