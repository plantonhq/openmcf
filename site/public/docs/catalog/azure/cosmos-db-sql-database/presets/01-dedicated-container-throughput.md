---
title: "Dedicated Container Throughput"
description: "This preset creates a database with NO throughput of its own -- a pure namespace. Each AzureCosmosdbSqlContainer inside it provisions its own dedicated RU/s, so workloads are isolated: a hot..."
type: "preset"
rank: "01"
presetSlug: "01-dedicated-container-throughput"
componentSlug: "cosmos-db-sql-database"
componentTitle: "Cosmos DB SQL Database"
provider: "azure"
icon: "package"
order: 1
---

# Dedicated Container Throughput

This preset creates a database with NO throughput of its own -- a pure
namespace. Each AzureCosmosdbSqlContainer inside it provisions its own
dedicated RU/s, so workloads are isolated: a hot container throttles
itself, never its siblings. This is the production default.

## When to Use

- Production workloads where containers must not compete for capacity
- Containers with meaningfully different traffic profiles (a hot
  operational store next to a quiet audit log)
- Any database whose containers you want to scale and bill individually

## Key Configuration Choices

- **No `throughput`, no `autoscaleMaxThroughput`** -- the deliberate
  choice, not an omission: unset means containers bring their own RU/s.
  A database-level budget would be SHARED by every container that does
  not bring its own, coupling noisy neighbors
- **Containers declare their RU/s** -- set `throughput` or
  `autoscaleMaxThroughput` on each AzureCosmosdbSqlContainer instead
- **The name is the blast radius** -- renaming replaces the database
  and every container in it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | The AzureCosmosdbAccount's Planton resource name (a GLOBAL_DOCUMENT_DB account) | Your Cosmos composition |
| `<database-name>` | 1-255 characters, unique within the account | Your naming convention |

## Downstream Wiring

Containers reference this database's ARM id:

```yaml
# On an AzureCosmosdbSqlContainer
sqlDatabaseId:
  valueFrom:
    kind: AzureCosmosdbSqlDatabase
    name: my-app-database
    fieldPath: status.outputs.sql_database_id
```
