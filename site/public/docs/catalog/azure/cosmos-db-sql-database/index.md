---
title: "Cosmos DB SQL Database"
description: "Cosmos DB SQL Database deployment documentation"
icon: "package"
order: 100
componentName: "azurecosmosdbsqldatabase"
---

# Azure Cosmos DB SQL Database

Creates a SQL (NoSQL) API database inside an AzureCosmosdbAccount -- the namespace containers live in and the boundary for shared throughput. A database either provisions RU/s that its containers share, or provisions nothing and lets each container bring dedicated throughput.

## What Gets Created

When you deploy an AzureCosmosdbSqlDatabase resource, Planton provisions:

- **Cosmos DB SQL Database** -- an `azurerm_cosmosdb_sql_database` on the referenced account, with optional shared throughput (fixed RU/s or autoscale)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureCosmosdbAccount** to create the database in (referenced through `cosmosdbAccountId`); the account must be a GLOBAL_DOCUMENT_DB (SQL API) account

## Quick Start

Create a file `database.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbSqlDatabase
metadata:
  name: app-data
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureCosmosdbSqlDatabase.app-data
spec:
  cosmosdbAccountId:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: my-app-cosmos
      fieldPath: status.outputs.cosmosdb_account_id
  databaseName: app-data
```

Deploy:

```shell
planton apply -f database.yaml
```

This database carries no throughput of its own -- each container brings dedicated RU/s, the common production shape because it isolates workloads from each other. To pool one budget across many small containers instead, set `throughput` (fixed, minimum 400 RU/s) or `autoscaleMaxThroughput` (ceiling, minimum 1000 RU/s) -- never both. On serverless accounts leave both unset; Azure rejects provisioned throughput there at apply.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `sql_database_id` | The ARM id -- what AzureCosmosdbSqlContainer's `sqlDatabaseId` references, and the management-plane scope for database-level operations |
| `sql_database_name` | What SDK calls reference inside the account's connection |
| `cosmosdb_account_name` | The account/database pair, without a second reference |

No endpoint or credential outputs on purpose: connectivity and keys live on the ACCOUNT (AzureCosmosdbAccount's endpoint and key outputs); the database is addressed inside that connection by name.

## Related Resources

- [Azure Cosmos DB Account](/docs/catalog/azure/cosmos-db-account) -- the parent account
- [Azure Cosmos DB SQL Container](/docs/catalog/azure/cosmos-db-sql-container) -- the containers that live in this database
