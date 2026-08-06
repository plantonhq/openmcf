# AzureCosmosdbMongoDatabase -- Design Research

## The Resource

A Cosmos DB MongoDB API database
(`Microsoft.DocumentDB/databaseAccounts/mongodbDatabases`) is the
namespace collections live in and the boundary for SHARED throughput --
Azure lets a database provision RU/s that all its collections split, or
provision nothing so each collection brings its own. The component maps
onto `azurerm_cosmosdb_mongo_database` (azurerm v4.x,
`internal/services/cosmos/cosmosdb_mongo_database_resource.go`),
parity-verified against pulumi-azure v6 (`cosmosdb.MongoDatabase`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `resource_group_name` + `account_name` | `cosmosdb_account_id` | azurerm addresses Cosmos children by the (resource group, account, name) trio; the spec models a single parent ARM-id reference and both modules parse the trio from it identically -- no redundant, contradictable state. ForceNew |
| `name` | `database_name` | Required, ForceNew, 1-255 characters (Azure's only constraint on Cosmos entity names); renaming replaces the database and everything in it |
| `throughput` | `throughput` | Optional; minimum 400 RU/s in increments of 100 (spec CELs); sent only when set -- serverless accounts reject provisioned throughput at apply |
| `autoscale_settings.max_throughput` | `autoscale_max_throughput` | Optional; minimum 1000 RU/s in increments of 1000 (spec CELs); Azure scales between 10% and 100% of the ceiling; XOR with `throughput` enforced by a message-level CEL |

No tags field: ARM does not support tags on Cosmos child resources, so
the platform's identity tags live on the account.

## Decomposition Decisions

- **First-class kind, not a fold**: databases are many-per-account with
  independent lifecycles and independent billing (a database's shared
  throughput is its own line item); collections reference databases
  individually, so the database must be independently addressable.
- **The parent is a single ARM-id FK**: the account's
  `cosmosdb_account_id` output is the one authoritative reference; the
  resource-group and account names azurerm wants are derived, never
  asked for twice. The account must be a MONGO_DB-kind account with the
  ENABLE_MONGO capability -- a cross-resource contract Azure enforces
  at apply.

## Recorded Skips (with reasons)

Nothing skipped: `azurerm_cosmosdb_mongo_database` exposes exactly the
name, parent addressing, and throughput surface, and the spec models
all of it.

## Operational Behavior Worth Knowing

- **Shared vs dedicated throughput is a design fork**: database-level
  throughput is SHARED by every collection that does not bring its own
  -- cheap for many small collections, but a noisy neighbor can starve
  its siblings. Dedicated per-collection throughput is the common
  production shape. A database with no throughput is a pure namespace.
- **Serverless accounts take neither field**: on accounts with the
  ENABLE_SERVERLESS capability, Azure rejects provisioned throughput at
  apply -- leave both fields unset.
- **Renaming replaces the resource and its data** -- the name is
  ForceNew, and everything in the database goes with it.
- **A database cannot move between accounts** -- the parent reference
  is fixed at creation.
- **Autoscale has a floor**: Azure scales between 10% and 100% of the
  ceiling, never below -- a 10000 RU/s ceiling always bills at least
  1000 RU/s. Size the ceiling to real peaks.

## Composition

- `cosmosdb_account_id` -> `AzureCosmosdbAccount.status.outputs.cosmosdb_account_id`
  (a MONGO_DB account with ENABLE_MONGO)
- `mongo_database_id` output <- `AzureCosmosdbMongoCollection.mongo_database_id`
  (collections live in this database)
- `mongo_database_name` + the account's MongoDB connection-string
  outputs <- application configuration (MongoDB drivers select the
  database by name inside the account's connection)
