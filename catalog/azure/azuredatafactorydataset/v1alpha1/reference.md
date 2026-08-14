# AzureDataFactoryDataset

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a delimited
# text (CSV) dataset exercising the variant's full surface -- a
# dynamic-path blob location, every parse setting, compression, and
# declared columns. References are literal values so the manifest
# validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDataFactoryDataset
metadata:
  name: test-dataset
  id: test-dataset
  org: test-org
  env: test
spec:
  dataFactoryId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DataFactory/factories/test-df
  name: orders-csv
  description: Daily order exports -- CSV files landing in blob storage.
  linkedServiceName:
    value: lakehouse-blob
  annotations:
    - team:data
  parameters:
    runDate: ""
  folder: ingest/orders
  delimitedText:
    azureBlobStorageLocation:
      container: landing
      path: raw/orders/@{dataset().runDate}
      dynamicPathEnabled: true
      filename: orders.csv
    columnDelimiter: ";"
    rowDelimiter: "\n"
    quoteCharacter: "'"
    escapeCharacter: "/"
    encoding: UTF-8
    firstRowAsHeader: true
    nullValue: "NULL"
    compressionCodec: gzip
    compressionLevel: Fastest
    schemaColumn:
      - name: order_id
        type: Int64
      - name: placed_at
        type: DateTime
        description: Order timestamp
      - name: total
        type: Decimal
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.dataFactoryId` | `string \| valueFrom` | yes |  | AzureDataFactory (`status.outputs.data_factory_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.linkedServiceName` | `string \| valueFrom` |  |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.description` | `string` |  |  |  |
| `spec.annotations` | `[]string` |  |  |  |
| `spec.parameters` | `map<string, string>` |  |  |  |
| `spec.additionalProperties` | `map<string, string>` |  |  |  |
| `spec.folder` | `string` |  |  |  |
| `spec.azureBlob` | `AzureDataFactoryDatasetAzureBlob` |  |  |  |
| `spec.azureBlob.path` | `string` |  |  |  |
| `spec.azureBlob.filename` | `string` |  |  |  |
| `spec.azureBlob.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.azureBlob.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.azureBlob.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.azureBlob.schemaColumn[].name` | `string` | yes |  |  |
| `spec.azureBlob.schemaColumn[].type` | `string` |  |  |  |
| `spec.azureBlob.schemaColumn[].description` | `string` |  |  |  |
| `spec.azureSqlTable` | `AzureDataFactoryDatasetAzureSqlTable` |  |  |  |
| `spec.azureSqlTable.linkedServiceId` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_id`) |
| `spec.azureSqlTable.schema` | `string` |  |  |  |
| `spec.azureSqlTable.table` | `string` |  |  |  |
| `spec.azureSqlTable.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.azureSqlTable.schemaColumn[].name` | `string` | yes |  |  |
| `spec.azureSqlTable.schemaColumn[].type` | `string` |  |  |  |
| `spec.azureSqlTable.schemaColumn[].description` | `string` |  |  |  |
| `spec.binary` | `AzureDataFactoryDatasetBinary` |  |  |  |
| `spec.binary.httpServerLocation` | `AzureDataFactoryDatasetHttpServerLocation` |  |  |  |
| `spec.binary.httpServerLocation.relativeUrl` | `string` | yes |  |  |
| `spec.binary.httpServerLocation.path` | `string` |  |  |  |
| `spec.binary.httpServerLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.binary.httpServerLocation.filename` | `string` |  |  |  |
| `spec.binary.httpServerLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.binary.azureBlobStorageLocation` | `AzureDataFactoryDatasetBlobStorageLocation` |  |  |  |
| `spec.binary.azureBlobStorageLocation.container` | `string` | yes |  |  |
| `spec.binary.azureBlobStorageLocation.dynamicContainerEnabled` | `bool` |  |  |  |
| `spec.binary.azureBlobStorageLocation.path` | `string` |  |  |  |
| `spec.binary.azureBlobStorageLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.binary.azureBlobStorageLocation.filename` | `string` |  |  |  |
| `spec.binary.azureBlobStorageLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.binary.sftpServerLocation` | `AzureDataFactoryDatasetSftpServerLocation` |  |  |  |
| `spec.binary.sftpServerLocation.path` | `string` | yes |  |  |
| `spec.binary.sftpServerLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.binary.sftpServerLocation.filename` | `string` | yes |  |  |
| `spec.binary.sftpServerLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.binary.compression` | `AzureDataFactoryDatasetBinaryCompression` |  |  |  |
| `spec.binary.compression.type` | `string` | yes |  |  |
| `spec.binary.compression.level` | `string` |  |  |  |
| `spec.cosmosdbSqlapi` | `AzureDataFactoryDatasetCosmosdbSqlapi` |  |  |  |
| `spec.cosmosdbSqlapi.collectionName` | `string` |  |  |  |
| `spec.cosmosdbSqlapi.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.cosmosdbSqlapi.schemaColumn[].name` | `string` | yes |  |  |
| `spec.cosmosdbSqlapi.schemaColumn[].type` | `string` |  |  |  |
| `spec.cosmosdbSqlapi.schemaColumn[].description` | `string` |  |  |  |
| `spec.custom` | `AzureDataFactoryDatasetCustom` |  |  |  |
| `spec.custom.linkedService` | `AzureDataFactoryDatasetCustomLinkedService` | yes |  |  |
| `spec.custom.linkedService.name` | `string \| valueFrom` | yes |  | AzureDataFactoryLinkedService (`status.outputs.linked_service_name`) |
| `spec.custom.linkedService.parameters` | `map<string, string>` |  |  |  |
| `spec.custom.type` | `string` | yes |  |  |
| `spec.custom.typePropertiesJson` | `string` | yes |  |  |
| `spec.custom.schemaJson` | `string` |  |  |  |
| `spec.delimitedText` | `AzureDataFactoryDatasetDelimitedText` |  |  |  |
| `spec.delimitedText.httpServerLocation` | `AzureDataFactoryDatasetHttpServerLocation` |  |  |  |
| `spec.delimitedText.httpServerLocation.relativeUrl` | `string` | yes |  |  |
| `spec.delimitedText.httpServerLocation.path` | `string` |  |  |  |
| `spec.delimitedText.httpServerLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.delimitedText.httpServerLocation.filename` | `string` |  |  |  |
| `spec.delimitedText.httpServerLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.delimitedText.azureBlobStorageLocation` | `AzureDataFactoryDatasetBlobStorageLocation` |  |  |  |
| `spec.delimitedText.azureBlobStorageLocation.container` | `string` | yes |  |  |
| `spec.delimitedText.azureBlobStorageLocation.dynamicContainerEnabled` | `bool` |  |  |  |
| `spec.delimitedText.azureBlobStorageLocation.path` | `string` |  |  |  |
| `spec.delimitedText.azureBlobStorageLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.delimitedText.azureBlobStorageLocation.filename` | `string` |  |  |  |
| `spec.delimitedText.azureBlobStorageLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.delimitedText.azureBlobFsLocation` | `AzureDataFactoryDatasetBlobFsLocation` |  |  |  |
| `spec.delimitedText.azureBlobFsLocation.fileSystem` | `string` |  |  |  |
| `spec.delimitedText.azureBlobFsLocation.dynamicFileSystemEnabled` | `bool` |  |  |  |
| `spec.delimitedText.azureBlobFsLocation.path` | `string` |  |  |  |
| `spec.delimitedText.azureBlobFsLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.delimitedText.azureBlobFsLocation.filename` | `string` |  |  |  |
| `spec.delimitedText.azureBlobFsLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.delimitedText.columnDelimiter` | `string` |  |  |  |
| `spec.delimitedText.rowDelimiter` | `string` |  |  |  |
| `spec.delimitedText.quoteCharacter` | `string` |  |  |  |
| `spec.delimitedText.escapeCharacter` | `string` |  |  |  |
| `spec.delimitedText.encoding` | `string` |  |  |  |
| `spec.delimitedText.firstRowAsHeader` | `bool` |  | `false` |  |
| `spec.delimitedText.nullValue` | `string` |  |  |  |
| `spec.delimitedText.compressionCodec` | `string` |  |  |  |
| `spec.delimitedText.compressionLevel` | `string` |  |  |  |
| `spec.delimitedText.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.delimitedText.schemaColumn[].name` | `string` | yes |  |  |
| `spec.delimitedText.schemaColumn[].type` | `string` |  |  |  |
| `spec.delimitedText.schemaColumn[].description` | `string` |  |  |  |
| `spec.http` | `AzureDataFactoryDatasetHttp` |  |  |  |
| `spec.http.relativeUrl` | `string` |  |  |  |
| `spec.http.requestBody` | `string` |  |  |  |
| `spec.http.requestMethod` | `string` |  |  |  |
| `spec.http.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.http.schemaColumn[].name` | `string` | yes |  |  |
| `spec.http.schemaColumn[].type` | `string` |  |  |  |
| `spec.http.schemaColumn[].description` | `string` |  |  |  |
| `spec.json` | `AzureDataFactoryDatasetJson` |  |  |  |
| `spec.json.httpServerLocation` | `AzureDataFactoryDatasetHttpServerLocation` |  |  |  |
| `spec.json.httpServerLocation.relativeUrl` | `string` | yes |  |  |
| `spec.json.httpServerLocation.path` | `string` |  |  |  |
| `spec.json.httpServerLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.json.httpServerLocation.filename` | `string` |  |  |  |
| `spec.json.httpServerLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.json.azureBlobStorageLocation` | `AzureDataFactoryDatasetBlobStorageLocation` |  |  |  |
| `spec.json.azureBlobStorageLocation.container` | `string` | yes |  |  |
| `spec.json.azureBlobStorageLocation.dynamicContainerEnabled` | `bool` |  |  |  |
| `spec.json.azureBlobStorageLocation.path` | `string` |  |  |  |
| `spec.json.azureBlobStorageLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.json.azureBlobStorageLocation.filename` | `string` |  |  |  |
| `spec.json.azureBlobStorageLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.json.encoding` | `string` |  |  |  |
| `spec.json.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.json.schemaColumn[].name` | `string` | yes |  |  |
| `spec.json.schemaColumn[].type` | `string` |  |  |  |
| `spec.json.schemaColumn[].description` | `string` |  |  |  |
| `spec.mysql` | `AzureDataFactoryDatasetMysql` |  |  |  |
| `spec.mysql.tableName` | `string` |  |  |  |
| `spec.mysql.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.mysql.schemaColumn[].name` | `string` | yes |  |  |
| `spec.mysql.schemaColumn[].type` | `string` |  |  |  |
| `spec.mysql.schemaColumn[].description` | `string` |  |  |  |
| `spec.parquet` | `AzureDataFactoryDatasetParquet` |  |  |  |
| `spec.parquet.httpServerLocation` | `AzureDataFactoryDatasetHttpServerLocation` |  |  |  |
| `spec.parquet.httpServerLocation.relativeUrl` | `string` | yes |  |  |
| `spec.parquet.httpServerLocation.path` | `string` |  |  |  |
| `spec.parquet.httpServerLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.parquet.httpServerLocation.filename` | `string` |  |  |  |
| `spec.parquet.httpServerLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.parquet.azureBlobStorageLocation` | `AzureDataFactoryDatasetBlobStorageLocation` |  |  |  |
| `spec.parquet.azureBlobStorageLocation.container` | `string` | yes |  |  |
| `spec.parquet.azureBlobStorageLocation.dynamicContainerEnabled` | `bool` |  |  |  |
| `spec.parquet.azureBlobStorageLocation.path` | `string` |  |  |  |
| `spec.parquet.azureBlobStorageLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.parquet.azureBlobStorageLocation.filename` | `string` |  |  |  |
| `spec.parquet.azureBlobStorageLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.parquet.azureBlobFsLocation` | `AzureDataFactoryDatasetBlobFsLocation` |  |  |  |
| `spec.parquet.azureBlobFsLocation.fileSystem` | `string` |  |  |  |
| `spec.parquet.azureBlobFsLocation.dynamicFileSystemEnabled` | `bool` |  |  |  |
| `spec.parquet.azureBlobFsLocation.path` | `string` |  |  |  |
| `spec.parquet.azureBlobFsLocation.dynamicPathEnabled` | `bool` |  |  |  |
| `spec.parquet.azureBlobFsLocation.filename` | `string` |  |  |  |
| `spec.parquet.azureBlobFsLocation.dynamicFilenameEnabled` | `bool` |  |  |  |
| `spec.parquet.compressionCodec` | `string` |  |  |  |
| `spec.parquet.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.parquet.schemaColumn[].name` | `string` | yes |  |  |
| `spec.parquet.schemaColumn[].type` | `string` |  |  |  |
| `spec.parquet.schemaColumn[].description` | `string` |  |  |  |
| `spec.postgresql` | `AzureDataFactoryDatasetPostgresql` |  |  |  |
| `spec.postgresql.tableName` | `string` |  |  |  |
| `spec.postgresql.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.postgresql.schemaColumn[].name` | `string` | yes |  |  |
| `spec.postgresql.schemaColumn[].type` | `string` |  |  |  |
| `spec.postgresql.schemaColumn[].description` | `string` |  |  |  |
| `spec.snowflake` | `AzureDataFactoryDatasetSnowflake` |  |  |  |
| `spec.snowflake.tableName` | `string` |  |  |  |
| `spec.snowflake.schemaName` | `string` |  |  |  |
| `spec.snowflake.schemaColumn` | `[]AzureDataFactoryDatasetSnowflakeSchemaColumn` |  |  |  |
| `spec.snowflake.schemaColumn[].name` | `string` | yes |  |  |
| `spec.snowflake.schemaColumn[].type` | `string` |  |  |  |
| `spec.snowflake.schemaColumn[].precision` | `int32` |  |  |  |
| `spec.snowflake.schemaColumn[].scale` | `int32` |  |  |  |
| `spec.sqlServerTable` | `AzureDataFactoryDatasetSqlServerTable` |  |  |  |
| `spec.sqlServerTable.tableName` | `string` |  |  |  |
| `spec.sqlServerTable.schemaColumn` | `[]AzureDataFactoryDatasetSchemaColumn` |  |  |  |
| `spec.sqlServerTable.schemaColumn[].name` | `string` | yes |  |  |
| `spec.sqlServerTable.schemaColumn[].type` | `string` |  |  |  |
| `spec.sqlServerTable.schemaColumn[].description` | `string` |  |  |  |

## Field Details

### spec.dataFactoryId

`string | valueFrom` · required

- references: AzureDataFactory (`status.outputs.data_factory_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactory, name: <that resource's name>, fieldPath: status.outputs.data_factory_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Dataset names must not consist entirely of the characters - . + ? / < > * % & : \
- rule: {"required":true}

### spec.linkedServiceName

`string | valueFrom`

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.description

`string`

### spec.annotations

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.parameters

`map<string, string>`

### spec.additionalProperties

`map<string, string>`

### spec.folder

`string`

### spec.azureBlob

`AzureDataFactoryDatasetAzureBlob`

### spec.azureBlob.path

`string`

### spec.azureBlob.filename

`string`

### spec.azureBlob.dynamicPathEnabled

`bool`

### spec.azureBlob.dynamicFilenameEnabled

`bool`

### spec.azureBlob.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.azureBlob.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.azureBlob.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.azureBlob.schemaColumn[].description

`string`

### spec.azureSqlTable

`AzureDataFactoryDatasetAzureSqlTable`

### spec.azureSqlTable.linkedServiceId

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_id}} -- a bare string does not parse

### spec.azureSqlTable.schema

`string`

### spec.azureSqlTable.table

`string`

### spec.azureSqlTable.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.azureSqlTable.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.azureSqlTable.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.azureSqlTable.schemaColumn[].description

`string`

### spec.binary

`AzureDataFactoryDatasetBinary`

- rule: Set exactly one of http_server_location, azure_blob_storage_location, or sftp_server_location
- rule: http_server_location requires both path and filename for the binary format

### spec.binary.httpServerLocation

`AzureDataFactoryDatasetHttpServerLocation`

### spec.binary.httpServerLocation.relativeUrl

`string` · required

- rule: {"required":true}

### spec.binary.httpServerLocation.path

`string`

### spec.binary.httpServerLocation.dynamicPathEnabled

`bool`

### spec.binary.httpServerLocation.filename

`string`

### spec.binary.httpServerLocation.dynamicFilenameEnabled

`bool`

### spec.binary.azureBlobStorageLocation

`AzureDataFactoryDatasetBlobStorageLocation`

### spec.binary.azureBlobStorageLocation.container

`string` · required

- rule: {"required":true}

### spec.binary.azureBlobStorageLocation.dynamicContainerEnabled

`bool`

### spec.binary.azureBlobStorageLocation.path

`string`

### spec.binary.azureBlobStorageLocation.dynamicPathEnabled

`bool`

### spec.binary.azureBlobStorageLocation.filename

`string`

### spec.binary.azureBlobStorageLocation.dynamicFilenameEnabled

`bool`

### spec.binary.sftpServerLocation

`AzureDataFactoryDatasetSftpServerLocation`

### spec.binary.sftpServerLocation.path

`string` · required

- rule: {"required":true}

### spec.binary.sftpServerLocation.dynamicPathEnabled

`bool`

### spec.binary.sftpServerLocation.filename

`string` · required

- rule: {"required":true}

### spec.binary.sftpServerLocation.dynamicFilenameEnabled

`bool`

### spec.binary.compression

`AzureDataFactoryDatasetBinaryCompression`

### spec.binary.compression.type

`string` · required

- rule: {"required":true,"string":{"in":["BZip2","Deflate","GZip","Tar","TarGZip","ZipDeflate"]}}

### spec.binary.compression.level

`string`

- rule: {"string":{"in":["","Optimal","Fastest"]}}

### spec.cosmosdbSqlapi

`AzureDataFactoryDatasetCosmosdbSqlapi`

### spec.cosmosdbSqlapi.collectionName

`string`

### spec.cosmosdbSqlapi.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.cosmosdbSqlapi.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.cosmosdbSqlapi.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.cosmosdbSqlapi.schemaColumn[].description

`string`

### spec.custom

`AzureDataFactoryDatasetCustom`

### spec.custom.linkedService

`AzureDataFactoryDatasetCustomLinkedService` · required

- rule: {"required":true}

### spec.custom.linkedService.name

`string | valueFrom` · required

- references: AzureDataFactoryLinkedService (`status.outputs.linked_service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryLinkedService, name: <that resource's name>, fieldPath: status.outputs.linked_service_name}} -- a bare string does not parse

### spec.custom.linkedService.parameters

`map<string, string>`

### spec.custom.type

`string` · required

- rule: {"required":true}

### spec.custom.typePropertiesJson

`string` · required

- rule: {"required":true}

### spec.custom.schemaJson

`string`

### spec.delimitedText

`AzureDataFactoryDatasetDelimitedText`

- rule: Set exactly one of http_server_location, azure_blob_storage_location, or azure_blob_fs_location
- rule: http_server_location requires both path and filename for the delimited text format

### spec.delimitedText.httpServerLocation

`AzureDataFactoryDatasetHttpServerLocation`

### spec.delimitedText.httpServerLocation.relativeUrl

`string` · required

- rule: {"required":true}

### spec.delimitedText.httpServerLocation.path

`string`

### spec.delimitedText.httpServerLocation.dynamicPathEnabled

`bool`

### spec.delimitedText.httpServerLocation.filename

`string`

### spec.delimitedText.httpServerLocation.dynamicFilenameEnabled

`bool`

### spec.delimitedText.azureBlobStorageLocation

`AzureDataFactoryDatasetBlobStorageLocation`

### spec.delimitedText.azureBlobStorageLocation.container

`string` · required

- rule: {"required":true}

### spec.delimitedText.azureBlobStorageLocation.dynamicContainerEnabled

`bool`

### spec.delimitedText.azureBlobStorageLocation.path

`string`

### spec.delimitedText.azureBlobStorageLocation.dynamicPathEnabled

`bool`

### spec.delimitedText.azureBlobStorageLocation.filename

`string`

### spec.delimitedText.azureBlobStorageLocation.dynamicFilenameEnabled

`bool`

### spec.delimitedText.azureBlobFsLocation

`AzureDataFactoryDatasetBlobFsLocation`

### spec.delimitedText.azureBlobFsLocation.fileSystem

`string`

### spec.delimitedText.azureBlobFsLocation.dynamicFileSystemEnabled

`bool`

### spec.delimitedText.azureBlobFsLocation.path

`string`

### spec.delimitedText.azureBlobFsLocation.dynamicPathEnabled

`bool`

### spec.delimitedText.azureBlobFsLocation.filename

`string`

### spec.delimitedText.azureBlobFsLocation.dynamicFilenameEnabled

`bool`

### spec.delimitedText.columnDelimiter

`string`

### spec.delimitedText.rowDelimiter

`string`

### spec.delimitedText.quoteCharacter

`string`

### spec.delimitedText.escapeCharacter

`string`

### spec.delimitedText.encoding

`string`

### spec.delimitedText.firstRowAsHeader

`bool` · optional (explicit presence)

- default: `false`

### spec.delimitedText.nullValue

`string`

### spec.delimitedText.compressionCodec

`string`

- rule: {"string":{"in":["","None","bzip2","gzip","deflate","ZipDeflate","TarGzip","Tar","snappy","lz4"]}}

### spec.delimitedText.compressionLevel

`string`

- rule: {"string":{"in":["","Optimal","Fastest"]}}

### spec.delimitedText.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.delimitedText.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.delimitedText.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.delimitedText.schemaColumn[].description

`string`

### spec.http

`AzureDataFactoryDatasetHttp`

### spec.http.relativeUrl

`string`

### spec.http.requestBody

`string`

### spec.http.requestMethod

`string`

### spec.http.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.http.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.http.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.http.schemaColumn[].description

`string`

### spec.json

`AzureDataFactoryDatasetJson`

- rule: Set exactly one of http_server_location or azure_blob_storage_location
- rule: http_server_location requires both path and filename for the JSON format
- rule: azure_blob_storage_location requires both path and filename for the JSON format

### spec.json.httpServerLocation

`AzureDataFactoryDatasetHttpServerLocation`

### spec.json.httpServerLocation.relativeUrl

`string` · required

- rule: {"required":true}

### spec.json.httpServerLocation.path

`string`

### spec.json.httpServerLocation.dynamicPathEnabled

`bool`

### spec.json.httpServerLocation.filename

`string`

### spec.json.httpServerLocation.dynamicFilenameEnabled

`bool`

### spec.json.azureBlobStorageLocation

`AzureDataFactoryDatasetBlobStorageLocation`

### spec.json.azureBlobStorageLocation.container

`string` · required

- rule: {"required":true}

### spec.json.azureBlobStorageLocation.dynamicContainerEnabled

`bool`

### spec.json.azureBlobStorageLocation.path

`string`

### spec.json.azureBlobStorageLocation.dynamicPathEnabled

`bool`

### spec.json.azureBlobStorageLocation.filename

`string`

### spec.json.azureBlobStorageLocation.dynamicFilenameEnabled

`bool`

### spec.json.encoding

`string`

### spec.json.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.json.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.json.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.json.schemaColumn[].description

`string`

### spec.mysql

`AzureDataFactoryDatasetMysql`

### spec.mysql.tableName

`string`

### spec.mysql.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.mysql.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.mysql.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.mysql.schemaColumn[].description

`string`

### spec.parquet

`AzureDataFactoryDatasetParquet`

- rule: Set exactly one of http_server_location, azure_blob_storage_location, or azure_blob_fs_location
- rule: http_server_location requires filename for the Parquet format

### spec.parquet.httpServerLocation

`AzureDataFactoryDatasetHttpServerLocation`

### spec.parquet.httpServerLocation.relativeUrl

`string` · required

- rule: {"required":true}

### spec.parquet.httpServerLocation.path

`string`

### spec.parquet.httpServerLocation.dynamicPathEnabled

`bool`

### spec.parquet.httpServerLocation.filename

`string`

### spec.parquet.httpServerLocation.dynamicFilenameEnabled

`bool`

### spec.parquet.azureBlobStorageLocation

`AzureDataFactoryDatasetBlobStorageLocation`

### spec.parquet.azureBlobStorageLocation.container

`string` · required

- rule: {"required":true}

### spec.parquet.azureBlobStorageLocation.dynamicContainerEnabled

`bool`

### spec.parquet.azureBlobStorageLocation.path

`string`

### spec.parquet.azureBlobStorageLocation.dynamicPathEnabled

`bool`

### spec.parquet.azureBlobStorageLocation.filename

`string`

### spec.parquet.azureBlobStorageLocation.dynamicFilenameEnabled

`bool`

### spec.parquet.azureBlobFsLocation

`AzureDataFactoryDatasetBlobFsLocation`

### spec.parquet.azureBlobFsLocation.fileSystem

`string`

### spec.parquet.azureBlobFsLocation.dynamicFileSystemEnabled

`bool`

### spec.parquet.azureBlobFsLocation.path

`string`

### spec.parquet.azureBlobFsLocation.dynamicPathEnabled

`bool`

### spec.parquet.azureBlobFsLocation.filename

`string`

### spec.parquet.azureBlobFsLocation.dynamicFilenameEnabled

`bool`

### spec.parquet.compressionCodec

`string`

- rule: {"string":{"in":["","bzip2","gzip","deflate","ZipDeflate","TarGzip","Tar","snappy","lz4"]}}

### spec.parquet.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.parquet.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.parquet.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.parquet.schemaColumn[].description

`string`

### spec.postgresql

`AzureDataFactoryDatasetPostgresql`

### spec.postgresql.tableName

`string`

### spec.postgresql.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.postgresql.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.postgresql.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.postgresql.schemaColumn[].description

`string`

### spec.snowflake

`AzureDataFactoryDatasetSnowflake`

### spec.snowflake.tableName

`string`

### spec.snowflake.schemaName

`string`

### spec.snowflake.schemaColumn

`[]AzureDataFactoryDatasetSnowflakeSchemaColumn`

### spec.snowflake.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.snowflake.schemaColumn[].type

`string`

- rule: {"string":{"in":["","NUMBER","DECIMAL","NUMERIC","INT","INTEGER","BIGINT","SMALLINT","FLOAT","FLOAT4","FLOAT8","DOUBLE","DOUBLE PRECISION","REAL","VARCHAR","CHAR","CHARACTER","STRING","TEXT","BINARY","VARBINARY","BOOLEAN","DATE","DATETIME","TIME","TIMESTAMP","TIMESTAMP_LTZ","TIMESTAMP_NTZ","TIMESTAMP_TZ","VARIANT","OBJECT","ARRAY","GEOGRAPHY"]}}

### spec.snowflake.schemaColumn[].precision

`int32`

- rule: {"int32":{"gte":0}}

### spec.snowflake.schemaColumn[].scale

`int32`

- rule: {"int32":{"gte":0}}

### spec.sqlServerTable

`AzureDataFactoryDatasetSqlServerTable`

### spec.sqlServerTable.tableName

`string`

### spec.sqlServerTable.schemaColumn

`[]AzureDataFactoryDatasetSchemaColumn`

### spec.sqlServerTable.schemaColumn[].name

`string` · required

- rule: {"required":true}

### spec.sqlServerTable.schemaColumn[].type

`string`

- rule: {"string":{"in":["","Byte","Byte[]","Boolean","Date","DateTime","DateTimeOffset","Decimal","Double","Guid","Int16","Int32","Int64","Single","String","TimeSpan"]}}

### spec.sqlServerTable.schemaColumn[].description

`string`

## Validation Rules

- `azure_data_factory_dataset_exactly_one_variant`: Set exactly one dataset variant block -- the variant determines the dataset type
- `azure_data_factory_dataset_linked_service_name_required`: linked_service_name is required for every variant except azure_sql_table and custom
- `azure_data_factory_dataset_linked_service_name_conflicts`: azure_sql_table and custom carry their own linked service reference -- do not also set linked_service_name

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureDataFactoryDataset, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.dataset_id` | `string` |  |
| `status.outputs.dataset_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.dataFactoryId` | AzureDataFactory | `status.outputs.data_factory_id` |
| `spec.linkedServiceName` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |
| `spec.azureSqlTable.linkedServiceId` | AzureDataFactoryLinkedService | `status.outputs.linked_service_id` |
| `spec.custom.linkedService.name` | AzureDataFactoryLinkedService | `status.outputs.linked_service_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureDataFactoryDataFlow | `spec.sources[].dataset.name` | `status.outputs.dataset_name` |
| AzureDataFactoryDataFlow | `spec.sinks[].dataset.name` | `status.outputs.dataset_name` |
| AzureDataFactoryDataFlow | `spec.transformations[].dataset.name` | `status.outputs.dataset_name` |

## See Also

- [Overview](../README.md)
