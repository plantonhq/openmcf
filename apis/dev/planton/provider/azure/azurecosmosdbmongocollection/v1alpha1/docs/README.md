# AzureCosmosdbMongoCollection -- Design Research

## The Resource

A Cosmos DB MongoDB API collection
(`Microsoft.DocumentDB/databaseAccounts/mongodbDatabases/collections`) is
the unit of storage and scale-out where documents live. Azure
partitions documents by the collection's shard key and indexes them per
the declared index list; throughput, TTL, and analytical retention are
set per collection. The component maps onto
`azurerm_cosmosdb_mongo_collection` (azurerm v4.x,
`internal/services/cosmos/cosmosdb_mongo_collection_resource.go`),
parity-verified against pulumi-azure v6 (`cosmosdb.MongoCollection`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `resource_group_name` + `account_name` + `database_name` | `mongo_database_id` | azurerm addresses Cosmos children by the (resource group, account, database, name) quartet; the spec models a single parent ARM-id reference and both modules parse the quartet from it identically. ForceNew |
| `name` | `collection_name` | Required, ForceNew, 1-255 characters; renaming replaces the collection and its data |
| `shard_key` | `shard_key` | Optional; ForceNew; unset creates an unsharded collection confined to one partition |
| `throughput` | `throughput` | Optional; minimum 400 RU/s in increments of 100; XOR with autoscale |
| `autoscale_settings.max_throughput` | `autoscale_max_throughput` | Optional; minimum 1000 RU/s in increments of 1000 |
| `default_ttl_seconds` | `default_ttl_seconds` | Optional; -1 enables TTL without default expiry; 0 rejected by ARM and spec |
| `analytical_storage_ttl` | `analytical_storage_ttl` | Optional; one-way disable semantics documented |
| `index` blocks | `indexes` | keys + unique; Azure requires the `_id` index unconditionally ("index with '_id' key is required") -- mirrored as a spec CEL so the failure surfaces at validation, not apply |

No tags field: ARM does not support tags on Cosmos child resources.

## Decomposition Decisions

- **First-class kind, not a fold**: collections are many-per-database with
  independent lifecycles and optional dedicated throughput billing.
- **The parent is a single ARM-id FK**: the database's `mongo_database_id`
  output is the one authoritative reference; resource-group, account, and
  database names azurerm wants are derived, never asked for twice.

## Recorded Skips (with reasons)

Nothing skipped: `azurerm_cosmosdb_mongo_collection` exposes exactly the
name, parent addressing, shard key, throughput, TTL, analytical TTL,
and index surface modeled in the spec.

## Deferred Follow-ups

- **SQL data-plane RBAC pair** (`SqlRoleDefinition` / `SqlRoleAssignment`)
  -- Entra keyless auth, separate from this Mongo collection surface
- **Mongo data-plane RBAC** -- evaluate on adoption demand
