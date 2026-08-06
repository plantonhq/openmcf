---
title: "Cosmos DB MongoDB Collection"
description: "Cosmos DB MongoDB Collection deployment documentation"
icon: "package"
order: 100
componentName: "azurecosmosdbmongocollection"
---

# Azure Cosmos DB MongoDB Collection

Creates a MongoDB API collection inside a Cosmos DB database -- the unit of storage and scale-out where documents live. Shard key, indexes, TTL, and dedicated throughput are set per collection.

## What Gets Created

When you deploy an AzureCosmosdbMongoCollection resource, Planton provisions:

- **Cosmos DB MongoDB Collection** -- an `azurerm_cosmosdb_mongo_collection` in the referenced database, with your chosen shard key, throughput mode, TTL, and indexes

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureCosmosdbMongoDatabase** to create the collection in (referenced through `mongoDatabaseId`), living in a MONGO_DB-kind account with the ENABLE_MONGO capability

## Quick Start

Create a file `collection.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureCosmosdbMongoCollection
metadata:
  name: events
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureCosmosdbMongoCollection.events
spec:
  mongoDatabaseId:
    valueFrom:
      kind: AzureCosmosdbMongoDatabase
      name: app-data
      fieldPath: status.outputs.mongo_database_id
  collectionName: events
  shardKey: tenantId
  autoscaleMaxThroughput: 4000
```

Deploy:

```shell
planton apply -f collection.yaml
```

The shard key is the collection's most consequential design decision -- pick a property with high cardinality and even request distribution (tenantId, userId, deviceId). It is fixed at creation.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `mongo_collection_id` | The ARM id -- management identity and collection-level data-plane RBAC scope |
| `mongo_collection_name` | What SDK calls reference inside the database |
| `mongo_database_name` | The parent database name, without a second reference |
| `cosmosdb_account_name` | The account name completing the addressing triple |

Connectivity and keys live on the account: consume AzureCosmosdbAccount's `endpoint` and key/connection-string outputs.

## Related Resources

- **AzureCosmosdbMongoDatabase** -- the parent database (and shared-throughput boundary)
- **AzureCosmosdbAccount** -- the account owning regions, consistency, network posture, and keys
- **AzureCosmosdbSqlContainer** -- the SQL (NoSQL) API equivalent
