---
title: "Cosmos DB MongoDB Database"
description: "Cosmos DB MongoDB Database deployment documentation"
icon: "package"
order: 100
componentName: "azurecosmosdbmongodatabase"
---

# Azure Cosmos DB MongoDB Database

Creates a MongoDB API database inside an AzureCosmosdbAccount -- the namespace collections live in and the boundary for shared throughput. A database either provisions RU/s that its collections share, or provisions nothing and lets each collection bring dedicated throughput.

## What Gets Created

When you deploy an AzureCosmosdbMongoDatabase resource, Planton provisions:

- **Cosmos DB MongoDB Database** -- an `azurerm_cosmosdb_mongo_database` on the referenced account, with optional shared throughput (fixed RU/s or autoscale)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureCosmosdbAccount** to create the database in (referenced through `cosmosdbAccountId`); the account must be a MONGO_DB-kind account with the ENABLE_MONGO capability

## Quick Start

Create a file `database.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbMongoDatabase
metadata:
  name: app-data
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureCosmosdbMongoDatabase.app-data
spec:
  cosmosdbAccountId:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: my-app-cosmos-mongo
      fieldPath: status.outputs.cosmosdb_account_id
  databaseName: app-data
```

Deploy:

```shell
planton apply -f database.yaml
```

This database carries no throughput of its own -- each collection brings dedicated RU/s, the common production shape because it isolates workloads from each other. To pool one budget across many small collections instead, set `throughput` (fixed, minimum 400 RU/s) or `autoscaleMaxThroughput` (ceiling, minimum 1000 RU/s) -- never both. On serverless accounts leave both unset; Azure rejects provisioned throughput there at apply.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `mongo_database_id` | The ARM id -- what AzureCosmosdbMongoCollection's `mongoDatabaseId` references, and the management-plane scope for database-level operations |
| `mongo_database_name` | What MongoDB drivers reference inside the account's connection |
| `cosmosdb_account_name` | The account/database pair, without a second reference |

No endpoint or credential outputs on purpose: connectivity and the MongoDB connection strings live on the ACCOUNT (AzureCosmosdbAccount's outputs); the database is addressed inside that connection by name.

## Related Resources

- [Azure Cosmos DB Account](/docs/catalog/azure/cosmos-db-account) -- the parent account
- [Azure Cosmos DB MongoDB Collection](/docs/catalog/azure/cosmos-db-mongodb-collection) -- the collections that live in this database
