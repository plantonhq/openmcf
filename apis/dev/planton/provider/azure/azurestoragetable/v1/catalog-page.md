# Azure Storage Table

Creates a Storage table inside an AzureStorageAccount -- the serverless NoSQL key/value store of Azure storage. Applications store schemaless entities by partition key + row key at petabyte scale, with no capacity provisioning.

## What Gets Created

When you deploy an AzureStorageTable resource, Planton provisions:

- **Storage Table** -- an `azurerm_storage_table` on the referenced account, with your stored access policies

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureStorageAccount** to create the table in (referenced through `storageAccountId`), with shared keys enabled (Azure's default) -- the provider drives table creation and ACLs through the data plane with shared-key authorization

## Quick Start

Create a file `table.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureStorageTable
metadata:
  name: app-entities
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageTable.app-entities
spec:
  storageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: my-app-storage
      fieldPath: status.outputs.storage_account_id
  tableName: AppEntities
```

Deploy:

```shell
planton apply -f table.yaml
```

Table names are stricter than other storage names: 3-63 alphanumerics starting with a letter, no hyphens, and never the literal word "table". When you need global distribution, throughput SLAs, or secondary indexes, reach for Cosmos DB's Table API instead.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `table_id` | The ARM id data-plane role assignments (Storage Table Data Reader/Contributor) scope to |
| `table_name` | What SDK clients, Functions table bindings, and app settings reference |
| `storage_account_name` | The account/table pair, without a second reference |

Client URLs compose from the ACCOUNT's endpoint output plus this table's name: `{primary_table_endpoint}{table_name}`.

## Related Resources

- [Azure Storage Account](/docs/catalog/azure/azurestorageaccount) -- the parent account
- [Azure Role Assignment](/docs/catalog/azure/azureroleassignment) -- table-scoped data-plane grants
- [Azure Cosmos DB Account](/docs/catalog/azure/azurecosmosdbaccount) -- the premium NoSQL sibling
