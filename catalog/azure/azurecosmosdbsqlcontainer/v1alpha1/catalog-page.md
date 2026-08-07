# Azure Cosmos DB SQL Container

Creates a SQL (NoSQL) API container inside a Cosmos DB database -- the unit of storage, indexing, and scale-out where documents live. Partition key, indexing policy, TTL, unique keys, conflict resolution, and dedicated throughput are set per container.

## What Gets Created

When you deploy an AzureCosmosdbSqlContainer resource, Planton provisions:

- **Cosmos DB SQL Container** -- an `azurerm_cosmosdb_sql_container` in the referenced database, with your chosen partition key, throughput mode, TTLs, unique keys, indexing policy, and conflict-resolution policy

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureCosmosdbSqlDatabase** to create the container in (referenced through `sqlDatabaseId`), living in a GLOBAL_DOCUMENT_DB (SQL API) account

## Quick Start

Create a file `container.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureCosmosdbSqlContainer
metadata:
  name: orders
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureCosmosdbSqlContainer.orders
spec:
  sqlDatabaseId:
    valueFrom:
      kind: AzureCosmosdbSqlDatabase
      name: app-data
      fieldPath: status.outputs.sql_database_id
  containerName: orders
  partitionKeyPaths:
    - /tenantId
  autoscaleMaxThroughput: 4000
```

Deploy:

```shell
planton apply -f container.yaml
```

The partition key is the container's most consequential design decision -- pick a property with high cardinality and even request distribution (tenantId, userId, deviceId), never a timestamp or a low-cardinality flag. It is fixed at creation.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `sql_container_id` | The ARM id -- management identity and container-level data-plane RBAC scope |
| `sql_container_name` | What SDK calls reference inside the database |
| `sql_database_name` | The parent database name, without a second reference |
| `cosmosdb_account_name` | The account name completing the addressing triple |

Connectivity and keys live on the account: consume AzureCosmosdbAccount's `endpoint` and key/connection-string outputs.

## Related Resources

- **AzureCosmosdbSqlDatabase** -- the parent database (and shared-throughput boundary)
- **AzureCosmosdbAccount** -- the account owning regions, consistency, network posture, and keys
- **AzureCosmosdbMongoCollection** -- the MongoDB-API equivalent
