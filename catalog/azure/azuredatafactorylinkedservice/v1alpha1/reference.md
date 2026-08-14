# AzureDataFactoryLinkedService

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a Databricks
# connection exercising the variant's full surface -- Key-Vault-held
# access token (the safe pattern), a fully-specified job cluster with
# autoscaling, Spark configuration, tags, and init scripts.
# References are literal values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDataFactoryLinkedService
metadata:
  name: test-linked-service
  id: test-linked-service
  org: test-org
  env: test
spec:
  dataFactoryId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DataFactory/factories/test-df
  name: spark-workspace
  description: Databricks over a Key-Vault-held token, running per-job clusters.
  annotations:
    - team:data
  parameters:
    env: test
  azureDatabricks:
    adbDomain: https://adb-1234567890123456.7.azuredatabricks.net
    keyVaultPassword:
      linkedServiceName:
        value: secrets-vault
      secretName: databricks-token
    newClusterConfig:
      nodeType: Standard_DS3_v2
      clusterVersion: 16.4.x-scala2.12
      minNumberOfWorkers: 2
      maxNumberOfWorkers: 8
      driverNodeType: Standard_DS4_v2
      sparkConfig:
        spark.speculation: "true"
      sparkEnvironmentVariables:
        PIPELINE_ENV: test
      customTags:
        team: data
      initScripts:
        - dbfs:/init/install-libs.sh
      logDestination: dbfs:/cluster-logs
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.dataFactoryId` | `string \| valueFrom` | yes |  | AzureDataFactory (`status.outputs.data_factory_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.annotations` | `[]string` |  |  |  |
| `spec.parameters` | `map<string, string>` |  |  |  |
| `spec.additionalProperties` | `map<string, string>` |  |  |  |
| `spec.integrationRuntimeName` | `string \| valueFrom` |  |  | AzureDataFactoryIntegrationRuntime (`status.outputs.integration_runtime_name`) |
| `spec.azureBlobStorage` | `AzureDataFactoryLinkedServiceAzureBlobStorage` |  |  |  |
| `spec.azureBlobStorage.connectionString` | `string` (sensitive) |  |  |  |
| `spec.azureBlobStorage.connectionStringInsecure` | `string` |  |  |  |
| `spec.azureBlobStorage.sasUri` | `string` (sensitive) |  |  |  |
| `spec.azureBlobStorage.serviceEndpoint` | `string \| valueFrom` (sensitive) |  |  | AzureStorageAccount (`status.outputs.primary_blob_endpoint`) |
| `spec.azureBlobStorage.sasTokenLinkedKeyVaultKey` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.azureBlobStorage.sasTokenLinkedKeyVaultKey.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureBlobStorage.sasTokenLinkedKeyVaultKey.secretName` | `string` | yes |  |  |
| `spec.azureBlobStorage.servicePrincipalLinkedKeyVaultKey` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.azureBlobStorage.servicePrincipalLinkedKeyVaultKey.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureBlobStorage.servicePrincipalLinkedKeyVaultKey.secretName` | `string` | yes |  |  |
| `spec.azureBlobStorage.storageKind` | `string` |  |  |  |
| `spec.azureBlobStorage.useManagedIdentity` | `bool` |  | `false` |  |
| `spec.azureBlobStorage.servicePrincipalId` | `string` |  |  |  |
| `spec.azureBlobStorage.servicePrincipalKey` | `string` (sensitive) |  |  |  |
| `spec.azureBlobStorage.tenantId` | `string` |  |  |  |
| `spec.azureDatabricks` | `AzureDataFactoryLinkedServiceAzureDatabricks` |  |  |  |
| `spec.azureDatabricks.adbDomain` | `string` | yes |  |  |
| `spec.azureDatabricks.msiWorkspaceId` | `string` |  |  |  |
| `spec.azureDatabricks.accessToken` | `string` (sensitive) |  |  |  |
| `spec.azureDatabricks.keyVaultPassword` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.azureDatabricks.keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureDatabricks.keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.azureDatabricks.existingClusterId` | `string` |  |  |  |
| `spec.azureDatabricks.newClusterConfig` | `AzureDataFactoryLinkedServiceDatabricksNewCluster` |  |  |  |
| `spec.azureDatabricks.newClusterConfig.nodeType` | `string` | yes |  |  |
| `spec.azureDatabricks.newClusterConfig.clusterVersion` | `string` | yes |  |  |
| `spec.azureDatabricks.newClusterConfig.minNumberOfWorkers` | `int32` |  | `1` |  |
| `spec.azureDatabricks.newClusterConfig.maxNumberOfWorkers` | `int32` |  |  |  |
| `spec.azureDatabricks.newClusterConfig.driverNodeType` | `string` |  |  |  |
| `spec.azureDatabricks.newClusterConfig.logDestination` | `string` |  |  |  |
| `spec.azureDatabricks.newClusterConfig.sparkConfig` | `map<string, string>` |  |  |  |
| `spec.azureDatabricks.newClusterConfig.sparkEnvironmentVariables` | `map<string, string>` |  |  |  |
| `spec.azureDatabricks.newClusterConfig.customTags` | `map<string, string>` |  |  |  |
| `spec.azureDatabricks.newClusterConfig.initScripts` | `[]string` |  |  |  |
| `spec.azureDatabricks.instancePool` | `AzureDataFactoryLinkedServiceDatabricksInstancePool` |  |  |  |
| `spec.azureDatabricks.instancePool.instancePoolId` | `string` | yes |  |  |
| `spec.azureDatabricks.instancePool.clusterVersion` | `string` | yes |  |  |
| `spec.azureDatabricks.instancePool.minNumberOfWorkers` | `int32` |  | `1` |  |
| `spec.azureDatabricks.instancePool.maxNumberOfWorkers` | `int32` |  |  |  |
| `spec.azureFileStorage` | `AzureDataFactoryLinkedServiceAzureFileStorage` |  |  |  |
| `spec.azureFileStorage.connectionString` | `string` (sensitive) | yes |  |  |
| `spec.azureFileStorage.fileShare` | `string` |  |  |  |
| `spec.azureFileStorage.host` | `string` |  |  |  |
| `spec.azureFileStorage.userId` | `string` |  |  |  |
| `spec.azureFileStorage.password` | `string` (sensitive) |  |  |  |
| `spec.azureFileStorage.keyVaultPassword` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.azureFileStorage.keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureFileStorage.keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.azureFunction` | `AzureDataFactoryLinkedServiceAzureFunction` |  |  |  |
| `spec.azureFunction.url` | `string` | yes |  |  |
| `spec.azureFunction.key` | `string` (sensitive) |  |  |  |
| `spec.azureFunction.keyVaultKey` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.azureFunction.keyVaultKey.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureFunction.keyVaultKey.secretName` | `string` | yes |  |  |
| `spec.azureSearch` | `AzureDataFactoryLinkedServiceAzureSearch` |  |  |  |
| `spec.azureSearch.url` | `string \| valueFrom` | yes |  | AzureSearchService (`status.outputs.endpoint`) |
| `spec.azureSearch.searchServiceKey` | `string` (sensitive) | yes |  |  |
| `spec.azureSqlDatabase` | `AzureDataFactoryLinkedServiceAzureSqlDatabase` |  |  |  |
| `spec.azureSqlDatabase.connectionString` | `string` (sensitive) |  |  |  |
| `spec.azureSqlDatabase.keyVaultConnectionString` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.azureSqlDatabase.keyVaultConnectionString.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureSqlDatabase.keyVaultConnectionString.secretName` | `string` | yes |  |  |
| `spec.azureSqlDatabase.keyVaultPassword` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.azureSqlDatabase.keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.azureSqlDatabase.keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.azureSqlDatabase.useManagedIdentity` | `bool` |  | `false` |  |
| `spec.azureSqlDatabase.servicePrincipalId` | `string` |  |  |  |
| `spec.azureSqlDatabase.servicePrincipalKey` | `string` (sensitive) |  |  |  |
| `spec.azureSqlDatabase.tenantId` | `string` |  |  |  |
| `spec.azureSqlDatabase.credentialName` | `string` |  |  |  |
| `spec.azureTableStorage` | `AzureDataFactoryLinkedServiceAzureTableStorage` |  |  |  |
| `spec.azureTableStorage.connectionString` | `string` (sensitive) | yes |  |  |
| `spec.cosmosdb` | `AzureDataFactoryLinkedServiceCosmosdb` |  |  |  |
| `spec.cosmosdb.connectionString` | `string` (sensitive) |  |  |  |
| `spec.cosmosdb.accountEndpoint` | `string` |  |  |  |
| `spec.cosmosdb.accountKey` | `string` (sensitive) |  |  |  |
| `spec.cosmosdb.database` | `string` |  |  |  |
| `spec.cosmosdbMongoapi` | `AzureDataFactoryLinkedServiceCosmosdbMongoapi` |  |  |  |
| `spec.cosmosdbMongoapi.connectionString` | `string` (sensitive) | yes |  |  |
| `spec.cosmosdbMongoapi.database` | `string` |  |  |  |
| `spec.cosmosdbMongoapi.serverVersionIs32OrHigher` | `bool` |  | `false` |  |
| `spec.custom` | `AzureDataFactoryLinkedServiceCustom` |  |  |  |
| `spec.custom.type` | `string` | yes |  |  |
| `spec.custom.typePropertiesJson` | `string` | yes |  |  |
| `spec.custom.integrationRuntimeParameters` | `map<string, string>` |  |  |  |
| `spec.dataLakeStorageGen2` | `AzureDataFactoryLinkedServiceDataLakeStorageGen2` |  |  |  |
| `spec.dataLakeStorageGen2.url` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.primary_dfs_endpoint`) |
| `spec.dataLakeStorageGen2.useManagedIdentity` | `bool` |  | `false` |  |
| `spec.dataLakeStorageGen2.storageAccountKey` | `string` (sensitive) |  |  |  |
| `spec.dataLakeStorageGen2.servicePrincipalId` | `string` |  |  |  |
| `spec.dataLakeStorageGen2.servicePrincipalKey` | `string` (sensitive) |  |  |  |
| `spec.dataLakeStorageGen2.tenant` | `string` |  |  |  |
| `spec.keyVault` | `AzureDataFactoryLinkedServiceKeyVault` |  |  |  |
| `spec.keyVault.keyVaultId` | `string \| valueFrom` | yes |  | AzureKeyVault (`status.outputs.key_vault_id`) |
| `spec.kusto` | `AzureDataFactoryLinkedServiceKusto` |  |  |  |
| `spec.kusto.kustoEndpoint` | `string` | yes |  |  |
| `spec.kusto.kustoDatabaseName` | `string` | yes |  |  |
| `spec.kusto.useManagedIdentity` | `bool` |  | `false` |  |
| `spec.kusto.servicePrincipalId` | `string` |  |  |  |
| `spec.kusto.servicePrincipalKey` | `string` (sensitive) |  |  |  |
| `spec.kusto.tenant` | `string` |  |  |  |
| `spec.mysql` | `AzureDataFactoryLinkedServiceMysql` |  |  |  |
| `spec.mysql.connectionString` | `string` (sensitive) | yes |  |  |
| `spec.mysql.driverVersion` | `string` |  | `V2` |  |
| `spec.odata` | `AzureDataFactoryLinkedServiceOdata` |  |  |  |
| `spec.odata.url` | `string` | yes |  |  |
| `spec.odata.basicAuthentication` | `AzureDataFactoryLinkedServiceBasicAuth` |  |  |  |
| `spec.odata.basicAuthentication.username` | `string` | yes |  |  |
| `spec.odata.basicAuthentication.password` | `string` (sensitive) | yes |  |  |
| `spec.odbc` | `AzureDataFactoryLinkedServiceOdbc` |  |  |  |
| `spec.odbc.connectionString` | `string` (sensitive) | yes |  |  |
| `spec.odbc.basicAuthentication` | `AzureDataFactoryLinkedServiceBasicAuth` |  |  |  |
| `spec.odbc.basicAuthentication.username` | `string` | yes |  |  |
| `spec.odbc.basicAuthentication.password` | `string` (sensitive) | yes |  |  |
| `spec.postgresql` | `AzureDataFactoryLinkedServicePostgresql` |  |  |  |
| `spec.postgresql.connectionString` | `string` (sensitive) | yes |  |  |
| `spec.sftp` | `AzureDataFactoryLinkedServiceSftp` |  |  |  |
| `spec.sftp.authenticationType` | `string` | yes |  |  |
| `spec.sftp.host` | `string` | yes |  |  |
| `spec.sftp.port` | `int32` | yes |  |  |
| `spec.sftp.username` | `string` | yes |  |  |
| `spec.sftp.password` | `string` (sensitive) |  |  |  |
| `spec.sftp.keyVaultPassword` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.sftp.keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sftp.keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.sftp.privateKeyContentBase64` | `string` (sensitive) |  |  |  |
| `spec.sftp.keyVaultPrivateKeyContentBase64` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.sftp.keyVaultPrivateKeyContentBase64.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sftp.keyVaultPrivateKeyContentBase64.secretName` | `string` | yes |  |  |
| `spec.sftp.privateKeyPath` | `string` |  |  |  |
| `spec.sftp.privateKeyPassphrase` | `string` (sensitive) |  |  |  |
| `spec.sftp.keyVaultPrivateKeyPassphrase` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.sftp.keyVaultPrivateKeyPassphrase.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sftp.keyVaultPrivateKeyPassphrase.secretName` | `string` | yes |  |  |
| `spec.sftp.skipHostKeyValidation` | `bool` |  |  |  |
| `spec.sftp.hostKeyFingerprint` | `string` |  |  |  |
| `spec.snowflake` | `AzureDataFactoryLinkedServiceSnowflake` |  |  |  |
| `spec.snowflake.connectionString` | `string` (sensitive) | yes |  |  |
| `spec.snowflake.keyVaultPassword` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.snowflake.keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.snowflake.keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.sqlManagedInstance` | `AzureDataFactoryLinkedServiceSqlManagedInstance` |  |  |  |
| `spec.sqlManagedInstance.connectionString` | `string` (sensitive) |  |  |  |
| `spec.sqlManagedInstance.keyVaultConnectionString` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.sqlManagedInstance.keyVaultConnectionString.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sqlManagedInstance.keyVaultConnectionString.secretName` | `string` | yes |  |  |
| `spec.sqlManagedInstance.keyVaultPassword` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.sqlManagedInstance.keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sqlManagedInstance.keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.sqlManagedInstance.servicePrincipalId` | `string` |  |  |  |
| `spec.sqlManagedInstance.servicePrincipalKey` | `string` (sensitive) |  |  |  |
| `spec.sqlManagedInstance.tenant` | `string` |  |  |  |
| `spec.sqlServer` | `AzureDataFactoryLinkedServiceSqlServer` |  |  |  |
| `spec.sqlServer.connectionString` | `string` (sensitive) |  |  |  |
| `spec.sqlServer.keyVaultConnectionString` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.sqlServer.keyVaultConnectionString.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sqlServer.keyVaultConnectionString.secretName` | `string` | yes |  |  |
| `spec.sqlServer.keyVaultPassword` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.sqlServer.keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.sqlServer.keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.sqlServer.userName` | `string` |  |  |  |
| `spec.synapse` | `AzureDataFactoryLinkedServiceSynapse` |  |  |  |
| `spec.synapse.connectionString` | `string` (sensitive) | yes |  |  |
| `spec.synapse.keyVaultPassword` | `AzureDataFactoryLinkedServiceKeyVaultSecretRef` |  |  |  |
| `spec.synapse.keyVaultPassword.linkedServiceName` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.synapse.keyVaultPassword.secretName` | `string` | yes |  |  |
| `spec.web` | `AzureDataFactoryLinkedServiceWeb` |  |  |  |
| `spec.web.url` | `string` | yes |  |  |
| `spec.web.authenticationType` | `string` | yes |  |  |
| `spec.web.username` | `string` |  |  |  |
| `spec.web.password` | `string` (sensitive) |  |  |  |

## Field Details

### spec.dataFactoryId

`string | valueFrom` · required

- references: AzureDataFactory (`status.outputs.data_factory_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactory, name: <that resource's name>, fieldPath: status.outputs.data_factory_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Linked service names must not consist entirely of the characters - . + ? / < > * % & : \
- rule: {"required":true}

### spec.description

`string`

### spec.annotations

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.parameters

`map<string, string>`

### spec.additionalProperties

`map<string, string>`

### spec.integrationRuntimeName

`string | valueFrom`

- references: AzureDataFactoryIntegrationRuntime (`status.outputs.integration_runtime_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryIntegrationRuntime, name: <that resource's name>, fieldPath: status.outputs.integration_runtime_name}} -- a bare string does not parse

### spec.azureBlobStorage

`AzureDataFactoryLinkedServiceAzureBlobStorage`

- rule: Set exactly one of connection_string, connection_string_insecure, sas_uri, or service_endpoint
- rule: use_managed_identity and service_principal_id are mutually exclusive

### spec.azureBlobStorage.connectionString

`string` · sensitive

### spec.azureBlobStorage.connectionStringInsecure

`string`

### spec.azureBlobStorage.sasUri

`string` · sensitive

### spec.azureBlobStorage.serviceEndpoint

`string | valueFrom` · sensitive

- references: AzureStorageAccount (`status.outputs.primary_blob_endpoint`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.primary_blob_endpoint}} -- a bare string does not parse

### spec.azureBlobStorage.sasTokenLinkedKeyVaultKey

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.azureBlobStorage.sasTokenLinkedKeyVaultKey.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureBlobStorage.sasTokenLinkedKeyVaultKey.secretName

`string` · required

- rule: {"required":true}

### spec.azureBlobStorage.servicePrincipalLinkedKeyVaultKey

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.azureBlobStorage.servicePrincipalLinkedKeyVaultKey.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureBlobStorage.servicePrincipalLinkedKeyVaultKey.secretName

`string` · required

- rule: {"required":true}

### spec.azureBlobStorage.storageKind

`string`

- rule: {"string":{"in":["","Storage","StorageV2","BlobStorage","BlockBlobStorage"]}}

### spec.azureBlobStorage.useManagedIdentity

`bool` · optional (explicit presence)

- default: `false`

### spec.azureBlobStorage.servicePrincipalId

`string`

- rule: service_principal_id must be a UUID

### spec.azureBlobStorage.servicePrincipalKey

`string` · sensitive

### spec.azureBlobStorage.tenantId

`string`

### spec.azureDatabricks

`AzureDataFactoryLinkedServiceAzureDatabricks`

- rule: Set exactly one of msi_workspace_id, access_token, or key_vault_password
- rule: Set exactly one of existing_cluster_id, new_cluster_config, or instance_pool

### spec.azureDatabricks.adbDomain

`string` · required

- rule: {"required":true}

### spec.azureDatabricks.msiWorkspaceId

`string`

### spec.azureDatabricks.accessToken

`string` · sensitive

### spec.azureDatabricks.keyVaultPassword

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.azureDatabricks.keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureDatabricks.keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.azureDatabricks.existingClusterId

`string`

### spec.azureDatabricks.newClusterConfig

`AzureDataFactoryLinkedServiceDatabricksNewCluster`

- rule: max_number_of_workers cannot be less than min_number_of_workers

### spec.azureDatabricks.newClusterConfig.nodeType

`string` · required

- rule: {"required":true}

### spec.azureDatabricks.newClusterConfig.clusterVersion

`string` · required

- rule: {"required":true}

### spec.azureDatabricks.newClusterConfig.minNumberOfWorkers

`int32` · optional (explicit presence)

- default: `1`
- rule: {"int32":{"lte":25000,"gte":1}}

### spec.azureDatabricks.newClusterConfig.maxNumberOfWorkers

`int32`

- rule: {"int32":{"lte":25000,"gte":0}}

### spec.azureDatabricks.newClusterConfig.driverNodeType

`string`

### spec.azureDatabricks.newClusterConfig.logDestination

`string`

### spec.azureDatabricks.newClusterConfig.sparkConfig

`map<string, string>`

### spec.azureDatabricks.newClusterConfig.sparkEnvironmentVariables

`map<string, string>`

### spec.azureDatabricks.newClusterConfig.customTags

`map<string, string>`

### spec.azureDatabricks.newClusterConfig.initScripts

`[]string`

### spec.azureDatabricks.instancePool

`AzureDataFactoryLinkedServiceDatabricksInstancePool`

- rule: max_number_of_workers cannot be less than min_number_of_workers

### spec.azureDatabricks.instancePool.instancePoolId

`string` · required

- rule: {"required":true}

### spec.azureDatabricks.instancePool.clusterVersion

`string` · required

- rule: {"required":true}

### spec.azureDatabricks.instancePool.minNumberOfWorkers

`int32` · optional (explicit presence)

- default: `1`
- rule: {"int32":{"lte":25000,"gte":1}}

### spec.azureDatabricks.instancePool.maxNumberOfWorkers

`int32`

- rule: {"int32":{"lte":25000,"gte":0}}

### spec.azureFileStorage

`AzureDataFactoryLinkedServiceAzureFileStorage`

### spec.azureFileStorage.connectionString

`string` · required · sensitive

- rule: {"required":true}

### spec.azureFileStorage.fileShare

`string`

### spec.azureFileStorage.host

`string`

### spec.azureFileStorage.userId

`string`

### spec.azureFileStorage.password

`string` · sensitive

### spec.azureFileStorage.keyVaultPassword

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.azureFileStorage.keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureFileStorage.keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.azureFunction

`AzureDataFactoryLinkedServiceAzureFunction`

- rule: Set exactly one of key or key_vault_key

### spec.azureFunction.url

`string` · required

- rule: {"required":true}

### spec.azureFunction.key

`string` · sensitive

### spec.azureFunction.keyVaultKey

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.azureFunction.keyVaultKey.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureFunction.keyVaultKey.secretName

`string` · required

- rule: {"required":true}

### spec.azureSearch

`AzureDataFactoryLinkedServiceAzureSearch`

### spec.azureSearch.url

`string | valueFrom` · required

- references: AzureSearchService (`status.outputs.endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSearchService, name: <that resource's name>, fieldPath: status.outputs.endpoint}} -- a bare string does not parse

### spec.azureSearch.searchServiceKey

`string` · required · sensitive

- rule: {"required":true}

### spec.azureSqlDatabase

`AzureDataFactoryLinkedServiceAzureSqlDatabase`

- rule: Set exactly one of connection_string or key_vault_connection_string
- rule: use_managed_identity and service_principal_id are mutually exclusive
- rule: service_principal_id and service_principal_key must be set together

### spec.azureSqlDatabase.connectionString

`string` · sensitive

### spec.azureSqlDatabase.keyVaultConnectionString

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.azureSqlDatabase.keyVaultConnectionString.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureSqlDatabase.keyVaultConnectionString.secretName

`string` · required

- rule: {"required":true}

### spec.azureSqlDatabase.keyVaultPassword

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.azureSqlDatabase.keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.azureSqlDatabase.keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.azureSqlDatabase.useManagedIdentity

`bool` · optional (explicit presence)

- default: `false`

### spec.azureSqlDatabase.servicePrincipalId

`string`

- rule: service_principal_id must be a UUID

### spec.azureSqlDatabase.servicePrincipalKey

`string` · sensitive

### spec.azureSqlDatabase.tenantId

`string`

### spec.azureSqlDatabase.credentialName

`string`

### spec.azureTableStorage

`AzureDataFactoryLinkedServiceAzureTableStorage`

### spec.azureTableStorage.connectionString

`string` · required · sensitive

- rule: {"required":true}

### spec.cosmosdb

`AzureDataFactoryLinkedServiceCosmosdb`

- rule: Set exactly one connection form -- connection_string, or account_endpoint + account_key (+ database)
- rule: account_endpoint, account_key, and database must all be set for account-detail authentication

### spec.cosmosdb.connectionString

`string` · sensitive

### spec.cosmosdb.accountEndpoint

`string`

### spec.cosmosdb.accountKey

`string` · sensitive

### spec.cosmosdb.database

`string`

### spec.cosmosdbMongoapi

`AzureDataFactoryLinkedServiceCosmosdbMongoapi`

### spec.cosmosdbMongoapi.connectionString

`string` · required · sensitive

- rule: {"required":true}

### spec.cosmosdbMongoapi.database

`string`

### spec.cosmosdbMongoapi.serverVersionIs32OrHigher

`bool` · optional (explicit presence)

- default: `false`

### spec.custom

`AzureDataFactoryLinkedServiceCustom`

### spec.custom.type

`string` · required

- rule: {"required":true}

### spec.custom.typePropertiesJson

`string` · required

- rule: {"required":true}

### spec.custom.integrationRuntimeParameters

`map<string, string>`

### spec.dataLakeStorageGen2

`AzureDataFactoryLinkedServiceDataLakeStorageGen2`

- rule: Set exactly one authentication mode -- use_managed_identity, storage_account_key, or a service principal (service_principal_id + service_principal_key + tenant)
- rule: A service principal needs service_principal_id, service_principal_key, and tenant together

### spec.dataLakeStorageGen2.url

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.primary_dfs_endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.primary_dfs_endpoint}} -- a bare string does not parse

### spec.dataLakeStorageGen2.useManagedIdentity

`bool` · optional (explicit presence)

- default: `false`

### spec.dataLakeStorageGen2.storageAccountKey

`string` · sensitive

### spec.dataLakeStorageGen2.servicePrincipalId

`string`

- rule: service_principal_id must be a UUID

### spec.dataLakeStorageGen2.servicePrincipalKey

`string` · sensitive

### spec.dataLakeStorageGen2.tenant

`string`

### spec.keyVault

`AzureDataFactoryLinkedServiceKeyVault`

### spec.keyVault.keyVaultId

`string | valueFrom` · required

- references: AzureKeyVault (`status.outputs.key_vault_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVault, name: <that resource's name>, fieldPath: status.outputs.key_vault_id}} -- a bare string does not parse

### spec.kusto

`AzureDataFactoryLinkedServiceKusto`

- rule: Set exactly one authentication mode -- use_managed_identity or a service principal (service_principal_id + service_principal_key)
- rule: service_principal_key is required with service_principal_id, and tenant belongs to the service principal mode

### spec.kusto.kustoEndpoint

`string` · required

- rule: {"required":true}

### spec.kusto.kustoDatabaseName

`string` · required

- rule: {"required":true}

### spec.kusto.useManagedIdentity

`bool` · optional (explicit presence)

- default: `false`

### spec.kusto.servicePrincipalId

`string`

- rule: service_principal_id must be a UUID

### spec.kusto.servicePrincipalKey

`string` · sensitive

### spec.kusto.tenant

`string`

### spec.mysql

`AzureDataFactoryLinkedServiceMysql`

### spec.mysql.connectionString

`string` · required · sensitive

- rule: {"required":true}

### spec.mysql.driverVersion

`string` · optional (explicit presence)

- default: `V2`
- rule: {"string":{"in":["V1","V2"]}}

### spec.odata

`AzureDataFactoryLinkedServiceOdata`

### spec.odata.url

`string` · required

- rule: {"required":true}

### spec.odata.basicAuthentication

`AzureDataFactoryLinkedServiceBasicAuth`

### spec.odata.basicAuthentication.username

`string` · required

- rule: {"required":true}

### spec.odata.basicAuthentication.password

`string` · required · sensitive

- rule: {"required":true}

### spec.odbc

`AzureDataFactoryLinkedServiceOdbc`

### spec.odbc.connectionString

`string` · required · sensitive

- rule: {"required":true}

### spec.odbc.basicAuthentication

`AzureDataFactoryLinkedServiceBasicAuth`

### spec.odbc.basicAuthentication.username

`string` · required

- rule: {"required":true}

### spec.odbc.basicAuthentication.password

`string` · required · sensitive

- rule: {"required":true}

### spec.postgresql

`AzureDataFactoryLinkedServicePostgresql`

### spec.postgresql.connectionString

`string` · required · sensitive

- rule: {"required":true}

### spec.sftp

`AzureDataFactoryLinkedServiceSftp`

- rule: Set at most one of password or key_vault_password
- rule: Set at most one of private_key_content_base64, key_vault_private_key_content_base64, or private_key_path
- rule: Set at most one of private_key_passphrase or key_vault_private_key_passphrase
- rule: password / key_vault_password require authentication_type Basic or MultiFactor
- rule: private key fields require authentication_type SshPublicKey or MultiFactor

### spec.sftp.authenticationType

`string` · required

- rule: {"required":true,"string":{"in":["Basic","MultiFactor","SshPublicKey"]}}

### spec.sftp.host

`string` · required

- rule: {"required":true}

### spec.sftp.port

`int32` · required

- rule: {"required":true}

### spec.sftp.username

`string` · required

- rule: {"required":true}

### spec.sftp.password

`string` · sensitive

### spec.sftp.keyVaultPassword

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.sftp.keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sftp.keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.sftp.privateKeyContentBase64

`string` · sensitive

- rule: private_key_content_base64 must be base64-encoded

### spec.sftp.keyVaultPrivateKeyContentBase64

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.sftp.keyVaultPrivateKeyContentBase64.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sftp.keyVaultPrivateKeyContentBase64.secretName

`string` · required

- rule: {"required":true}

### spec.sftp.privateKeyPath

`string`

### spec.sftp.privateKeyPassphrase

`string` · sensitive

### spec.sftp.keyVaultPrivateKeyPassphrase

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.sftp.keyVaultPrivateKeyPassphrase.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sftp.keyVaultPrivateKeyPassphrase.secretName

`string` · required

- rule: {"required":true}

### spec.sftp.skipHostKeyValidation

`bool` · optional (explicit presence)

### spec.sftp.hostKeyFingerprint

`string`

### spec.snowflake

`AzureDataFactoryLinkedServiceSnowflake`

### spec.snowflake.connectionString

`string` · required · sensitive

- rule: {"required":true}

### spec.snowflake.keyVaultPassword

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.snowflake.keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.snowflake.keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.sqlManagedInstance

`AzureDataFactoryLinkedServiceSqlManagedInstance`

- rule: Set exactly one of connection_string or key_vault_connection_string
- rule: service_principal_id, service_principal_key, and tenant must be set together

### spec.sqlManagedInstance.connectionString

`string` · sensitive

### spec.sqlManagedInstance.keyVaultConnectionString

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.sqlManagedInstance.keyVaultConnectionString.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sqlManagedInstance.keyVaultConnectionString.secretName

`string` · required

- rule: {"required":true}

### spec.sqlManagedInstance.keyVaultPassword

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.sqlManagedInstance.keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sqlManagedInstance.keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.sqlManagedInstance.servicePrincipalId

`string`

- rule: service_principal_id must be a UUID

### spec.sqlManagedInstance.servicePrincipalKey

`string` · sensitive

### spec.sqlManagedInstance.tenant

`string`

- rule: tenant must be a UUID

### spec.sqlServer

`AzureDataFactoryLinkedServiceSqlServer`

- rule: Set exactly one of connection_string or key_vault_connection_string

### spec.sqlServer.connectionString

`string` · sensitive

### spec.sqlServer.keyVaultConnectionString

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.sqlServer.keyVaultConnectionString.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sqlServer.keyVaultConnectionString.secretName

`string` · required

- rule: {"required":true}

### spec.sqlServer.keyVaultPassword

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.sqlServer.keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.sqlServer.keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.sqlServer.userName

`string`

### spec.synapse

`AzureDataFactoryLinkedServiceSynapse`

### spec.synapse.connectionString

`string` · required · sensitive

- rule: {"required":true}

### spec.synapse.keyVaultPassword

`AzureDataFactoryLinkedServiceKeyVaultSecretRef`

### spec.synapse.keyVaultPassword.linkedServiceName

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.synapse.keyVaultPassword.secretName

`string` · required

- rule: {"required":true}

### spec.web

`AzureDataFactoryLinkedServiceWeb`

- rule: authentication_type Basic requires username and password

### spec.web.url

`string` · required

- rule: {"required":true}

### spec.web.authenticationType

`string` · required

- rule: {"required":true,"string":{"in":["Anonymous","Basic"]}}

### spec.web.username

`string`

### spec.web.password

`string` · sensitive

## Validation Rules

- `azure_data_factory_linked_service_exactly_one_variant`: Set exactly one connection variant block -- the variant determines the linked service type

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureDataFactoryLinkedService, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.linked_service_id` | `string` |  |
| `status.outputs.linked_service_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.dataFactoryId` | AzureDataFactory | `status.outputs.data_factory_id` |
| `spec.integrationRuntimeName` | AzureDataFactoryIntegrationRuntime | `status.outputs.integration_runtime_name` |
| `spec.azureBlobStorage.serviceEndpoint` | AzureStorageAccount | `status.outputs.primary_blob_endpoint` |
| `spec.azureBlobStorage.sasTokenLinkedKeyVaultKey.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureBlobStorage.servicePrincipalLinkedKeyVaultKey.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureDatabricks.keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureFileStorage.keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureFunction.keyVaultKey.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureSearch.url` | AzureSearchService | `status.outputs.endpoint` |
| `spec.azureSqlDatabase.keyVaultConnectionString.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureSqlDatabase.keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.dataLakeStorageGen2.url` | AzureStorageAccount | `status.outputs.primary_dfs_endpoint` |
| `spec.keyVault.keyVaultId` | AzureKeyVault | `status.outputs.key_vault_id` |
| `spec.sftp.keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sftp.keyVaultPrivateKeyContentBase64.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sftp.keyVaultPrivateKeyPassphrase.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.snowflake.keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sqlManagedInstance.keyVaultConnectionString.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sqlManagedInstance.keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sqlServer.keyVaultConnectionString.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.sqlServer.keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.synapse.keyVaultPassword.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureDataFactoryDataFlow | `spec.sources[].linkedService.name` | `status.outputs.linked_service_name` |
| AzureDataFactoryDataFlow | `spec.sources[].schemaLinkedService.name` | `status.outputs.linked_service_name` |
| AzureDataFactoryDataFlow | `spec.sinks[].linkedService.name` | `status.outputs.linked_service_name` |
| AzureDataFactoryDataFlow | `spec.sinks[].schemaLinkedService.name` | `status.outputs.linked_service_name` |
| AzureDataFactoryDataFlow | `spec.sinks[].rejectedLinkedService.name` | `status.outputs.linked_service_name` |
| AzureDataFactoryDataFlow | `spec.transformations[].linkedService.name` | `status.outputs.linked_service_name` |
| AzureDataFactoryDataset | `spec.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryDataset | `spec.azureSqlTable.linkedServiceId` | `status.outputs.linked_service_id` |
| AzureDataFactoryDataset | `spec.custom.linkedService.name` | `status.outputs.linked_service_name` |
| AzureDataFactoryIntegrationRuntime | `spec.azureSsis.expressCustomSetup.commandKey[].keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryIntegrationRuntime | `spec.azureSsis.expressCustomSetup.component[].keyVaultLicense.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryIntegrationRuntime | `spec.azureSsis.packageStore[].linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryIntegrationRuntime | `spec.azureSsis.proxy.stagingStorageLinkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.azureBlobStorage.sasTokenLinkedKeyVaultKey.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.azureBlobStorage.servicePrincipalLinkedKeyVaultKey.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.azureDatabricks.keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.azureFileStorage.keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.azureFunction.keyVaultKey.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.azureSqlDatabase.keyVaultConnectionString.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.azureSqlDatabase.keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.sftp.keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.sftp.keyVaultPrivateKeyContentBase64.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.sftp.keyVaultPrivateKeyPassphrase.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.snowflake.keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.sqlManagedInstance.keyVaultConnectionString.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.sqlManagedInstance.keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.sqlServer.keyVaultConnectionString.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.sqlServer.keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |
| AzureDataFactoryLinkedService | `spec.synapse.keyVaultPassword.linkedServiceName` | `status.outputs.linked_service_name` |

## See Also

- [Overview](../README.md)
