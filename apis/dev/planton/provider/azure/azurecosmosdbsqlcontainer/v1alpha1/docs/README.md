# AzureCosmosdbSqlContainer -- Design Research

## The Resource

A Cosmos DB SQL container
(`Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers`) is
the unit of storage, indexing, and scale-out in the SQL (NoSQL) API --
documents live in containers, Azure partitions them by the container's
partition key, and RU throughput is billed per container (or shared
from the database). The component maps onto
`azurerm_cosmosdb_sql_container` (azurerm v4.x,
`internal/services/cosmos/cosmosdb_sql_container_resource.go`),
parity-verified against pulumi-azure v6 (`cosmosdb.SqlContainer`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `resource_group_name` + `account_name` + `database_name` | `sql_database_id` | The provider addresses children by the name tuple; the spec models a single parent ARM-id reference and BOTH modules parse the three names from it with identical anchored semantics -- no redundant, contradictable state, and no parity divergence |
| `name` | `container_name` | Required, ForceNew, 1-255 chars (Azure's only Cosmos entity-name constraint) |
| `partition_key_paths` | `partition_key_paths` | Required, ForceNew, 1-3 paths each starting with "/"; the spec adds ARM's hierarchy contract (HASH = exactly one path, MULTI_HASH = several) that the provider leaves to apply time |
| `partition_key_kind` | enum | Hash (default) / MultiHash; ForceNew |
| `partition_key_version` | `partition_key_version` | 1-2; the spec enforces MULTI_HASH-requires-version-2 (ARM's contract, an apply-time error in the provider) |
| `throughput` / `autoscale_settings.max_throughput` | `throughput` XOR `autoscale_max_throughput` | The provider's ConflictsWith as a message rule; min 400 step 100 / min 1000 step 1000 as field rules |
| `default_ttl` | `default_ttl` | >= -1; -1 = TTL on with per-document expiry only |
| `analytical_storage_ttl` | `analytical_storage_ttl` | >= -1; disabling on an existing container forces replacement (provider CustomizeDiff) |
| `unique_key` | `unique_keys` | ForceNew; uniqueness is scoped to the logical partition |
| `indexing_policy` | message | consistent/none modes (the provider excludes the service's legacy automatic/lazy); included/excluded/composite/spatial paths; updatable in place |
| `conflict_resolution_policy` | message | ForceNew; LastWriterWins/Custom with the per-mode field pairing as a message rule |

## Decomposition Decisions

- **First-class kind, not a fold**: containers are many-per-database
  with independent lifecycles, their own throughput billing, and their
  own data-plane RBAC scope -- the classic split-test pass.
- **Partition key version is modeled explicitly** (1 or 2) rather than
  derived from the kind: version is ARM state with its own semantics
  (2 = large-key support, required for MultiHash), and deriving it
  would hide a ForceNew trigger.

## Recorded Skips (with reasons)

- **`azurerm_cosmosdb_sql_trigger` / `sql_stored_procedure` /
  `sql_function`** -- JavaScript code artifacts keyed off the
  container: application content, not infrastructure (the
  table-entities precedent). Applications deploy their own server-side
  code.
- **`spatial_index.types`** -- computed-only in the provider (the
  service infers geometry types); nothing to declare.
- **No endpoint/credential outputs** -- connectivity and keys live on
  the ACCOUNT; the container is addressed by name inside the account's
  connection.

## Operational Behavior Worth Knowing

- **The partition key is forever**: changing paths, kind, or version
  replaces the container; migrations copy data to a new container.
- **A declared indexing policy replaces the default wholesale**: when
  any included/excluded path is set, include "/*" explicitly and carve
  out exceptions -- otherwise queries lose their indexes.
- **Composite indexes must match query ORDER BY order** (or its exact
  reverse) to serve the query.
- **Dedicated vs shared throughput is a create-time choice in
  practice**: a container created sharing the database's throughput
  cannot later be flipped to dedicated in place.
- **Serverless accounts reject provisioned throughput** at apply --
  leave both throughput fields unset.
- **Conflict resolution only matters on multi-region-write accounts**,
  and it is fixed at creation -- set it up front on containers that
  might ever go active-active.

## Composition

- `sql_database_id` -> `AzureCosmosdbSqlDatabase.status.outputs.sql_database_id`
- `sql_container_id` output <- container-scoped data-plane RBAC and
  management references
