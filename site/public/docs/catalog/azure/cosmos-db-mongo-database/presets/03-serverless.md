---
title: "Serverless Database"
description: "This preset creates a MongoDB API database on a serverless Cosmos DB account -- one whose AzureCosmosdbAccount declares the ENABLE_SERVERLESS capability alongside ENABLE_MONGO. Serverless accounts..."
type: "preset"
rank: "03"
presetSlug: "03-serverless"
componentSlug: "cosmos-db-mongo-database"
componentTitle: "Cosmos DB Mongo Database"
provider: "azure"
icon: "package"
order: 3
---

# Serverless Database

This preset creates a MongoDB API database on a serverless Cosmos DB
account -- one whose AzureCosmosdbAccount declares the
ENABLE_SERVERLESS capability alongside ENABLE_MONGO. Serverless
accounts bill per request instead of per provisioned RU/s, so the
database (and its collections) must not declare throughput at all:
Azure rejects provisioned throughput on serverless at apply.

## When to Use

- Intermittent or unpredictable traffic where provisioned capacity
  would sit idle (internal tools, prototypes, event-driven jobs)
- New workloads whose traffic profile is still unknown -- measure on
  serverless, move to provisioned throughput when the profile settles
- Cost floors matter more than throughput ceilings

## Key Configuration Choices

- **No `throughput`, no `autoscaleMaxThroughput`** -- a hard
  requirement on serverless, not a stylistic choice; Azure rejects
  either at apply
- **The billing mode lives on the ACCOUNT** -- the referenced
  AzureCosmosdbAccount must carry the ENABLE_SERVERLESS capability;
  nothing on the database selects serverless
- **Collections inherit the constraint** -- AzureCosmosdbMongoCollection
  resources in this database must also leave both throughput fields
  unset

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<serverless-account-resource-name>` | A MONGO_DB-kind AzureCosmosdbAccount with ENABLE_MONGO and ENABLE_SERVERLESS | Your Cosmos composition |
| `<database-name>` | 1-255 characters, unique within the account | Your naming convention |

## Downstream Wiring

Collections reference the database the usual way -- and, on
serverless, also omit throughput:

```yaml
# On an AzureCosmosdbMongoCollection
mongoDatabaseId:
  valueFrom:
    kind: AzureCosmosdbMongoDatabase
    name: my-serverless-database
    fieldPath: status.outputs.mongo_database_id
```
