---
title: "Dedicated Collection Throughput"
description: "This preset creates a database with NO throughput of its own -- a pure namespace. Each AzureCosmosdbMongoCollection inside it provisions its own dedicated RU/s, so workloads are isolated: a hot..."
type: "preset"
rank: "01"
presetSlug: "01-dedicated-collection-throughput"
componentSlug: "cosmos-db-mongodb-database"
componentTitle: "Cosmos DB MongoDB Database"
provider: "azure"
icon: "package"
order: 1
---

# Dedicated Collection Throughput

This preset creates a database with NO throughput of its own -- a pure
namespace. Each AzureCosmosdbMongoCollection inside it provisions its
own dedicated RU/s, so workloads are isolated: a hot collection
throttles itself, never its siblings. This is the production default.

## When to Use

- Production workloads where collections must not compete for capacity
- Collections with meaningfully different traffic profiles (a hot
  operational store next to a quiet audit log)
- Any database whose collections you want to scale and bill
  individually

## Key Configuration Choices

- **No `throughput`, no `autoscaleMaxThroughput`** -- the deliberate
  choice, not an omission: unset means collections bring their own
  RU/s. A database-level budget would be SHARED by every collection
  that does not bring its own, coupling noisy neighbors
- **Collections declare their RU/s** -- set `throughput` or
  `autoscaleMaxThroughput` on each AzureCosmosdbMongoCollection instead
- **The name is the blast radius** -- renaming replaces the database
  and every collection in it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | A MONGO_DB-kind AzureCosmosdbAccount with the ENABLE_MONGO capability | Your Cosmos composition |
| `<database-name>` | 1-255 characters, unique within the account | Your naming convention |

## Downstream Wiring

Collections reference this database's ARM id:

```yaml
# On an AzureCosmosdbMongoCollection
mongoDatabaseId:
  valueFrom:
    kind: AzureCosmosdbMongoDatabase
    name: my-app-database
    fieldPath: status.outputs.mongo_database_id
```
