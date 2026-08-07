# AzureCosmosdbSqlContainer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureCosmosdbSqlContainerSpec** defines the configuration for creating
a SQL (NoSQL) API container inside a Cosmos DB database: the unit of
storage, indexing, and scale-out. Documents live in containers; Azure
distributes them across physical partitions by the partition key and
indexes them by the container's indexing policy.

**The partition key is the single most consequential design decision**
for Cosmos DB performance and cost. Choose a property with high
cardinality (many distinct values), even request distribution, and
frequent use in query filters -- tenantId, userId, deviceId. Avoid
timestamps (hot partition) and low-cardinality flags. The key is fixed
at creation: changing it means a new container and a data migration.

**Throughput**: a container either shares its database's provisioned
throughput (leave both fields unset), or brings its own dedicated
fixed/autoscale throughput. On serverless accounts neither may be set.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureCosmosdbSqlContainer
metadata:
  name: test-cosmos-sql-container
spec:
  sqlDatabaseId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DocumentDB/databaseAccounts/planton-hack-cosmos/sqlDatabases/app-data
  containerName: orders
  # Exercises the hierarchical (MultiHash) partition-key seam, which
  # requires version 2 -- the spec enforces the pairing.
  partitionKeyPaths:
    - /tenantId
    - /orderId
  partitionKeyKind: MULTI_HASH
  partitionKeyVersion: 2
  # Exercises the fixed-throughput seam (mutually exclusive with
  # autoscale -- the spec rejects both together).
  throughput: 400
  # TTL on with per-document expiry only.
  defaultTtl: -1
  # Exercises the unique-key block rendering.
  uniqueKeys:
    - paths:
        - /orderNumber
  # Exercises every indexing-policy seam: mode, include-all with an
  # excluded payload subtree, a composite index, and a spatial index.
  indexingPolicy:
    indexingMode: CONSISTENT
    includedPaths:
      - path: /*
    excludedPaths:
      - path: /payload/*
    compositeIndexes:
      - entries:
          - path: /tenantId
            order: ASCENDING
          - path: /createdAt
            order: DESCENDING
    spatialIndexes:
      - path: /location/*
  # Exercises the conflict-resolution enum mapping.
  conflictResolutionPolicy:
    mode: LAST_WRITER_WINS
    conflictResolutionPath: /_ts
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.sqlDatabaseId` | `string \| valueFrom` | yes |  | AzureCosmosdbSqlDatabase (`status.outputs.sql_database_id`) |
| `spec.containerName` | `string` | yes |  |  |
| `spec.partitionKeyPaths` | `[]string` | yes |  |  |
| `spec.partitionKeyKind` | `enum` |  |  |  |
| `spec.partitionKeyVersion` | `int32` |  |  |  |
| `spec.throughput` | `int32` |  |  |  |
| `spec.autoscaleMaxThroughput` | `int32` |  |  |  |
| `spec.defaultTtl` | `int32` |  |  |  |
| `spec.analyticalStorageTtl` | `int32` |  |  |  |
| `spec.uniqueKeys` | `[]AzureCosmosdbSqlContainerUniqueKey` |  |  |  |
| `spec.uniqueKeys[].paths` | `[]string` | yes |  |  |
| `spec.indexingPolicy` | `AzureCosmosdbSqlContainerIndexingPolicy` |  |  |  |
| `spec.indexingPolicy.indexingMode` | `enum` |  |  |  |
| `spec.indexingPolicy.includedPaths` | `[]AzureCosmosdbSqlContainerIndexPath` |  |  |  |
| `spec.indexingPolicy.includedPaths[].path` | `string` | yes |  |  |
| `spec.indexingPolicy.excludedPaths` | `[]AzureCosmosdbSqlContainerIndexPath` |  |  |  |
| `spec.indexingPolicy.excludedPaths[].path` | `string` | yes |  |  |
| `spec.indexingPolicy.compositeIndexes` | `[]AzureCosmosdbSqlContainerCompositeIndex` |  |  |  |
| `spec.indexingPolicy.compositeIndexes[].entries` | `[]AzureCosmosdbSqlContainerCompositeIndexEntry` | yes |  |  |
| `spec.indexingPolicy.compositeIndexes[].entries[].path` | `string` | yes |  |  |
| `spec.indexingPolicy.compositeIndexes[].entries[].order` | `enum` |  |  |  |
| `spec.indexingPolicy.spatialIndexes` | `[]AzureCosmosdbSqlContainerSpatialIndex` |  |  |  |
| `spec.indexingPolicy.spatialIndexes[].path` | `string` | yes |  |  |
| `spec.conflictResolutionPolicy` | `AzureCosmosdbSqlContainerConflictResolutionPolicy` |  |  |  |
| `spec.conflictResolutionPolicy.mode` | `enum` | yes |  |  |
| `spec.conflictResolutionPolicy.conflictResolutionPath` | `string` |  |  |  |
| `spec.conflictResolutionPolicy.conflictResolutionProcedure` | `string` |  |  |  |

## Field Details

### spec.sqlDatabaseId

`string | valueFrom` · required

The SQL database the container lives in, by ARM ID. References an
AzureCosmosdbSqlDatabase's sql_database_id output so the account,
database, and containers compose in one manifest set. Fixed at
creation.

- references: AzureCosmosdbSqlDatabase (`status.outputs.sql_database_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureCosmosdbSqlDatabase, name: <that resource's name>, fieldPath: status.outputs.sql_database_id}} -- a bare string does not parse

### spec.containerName

`string` · required

The container's name -- unique within the database, 1-255 characters
(Azure's only constraint on Cosmos entity names). Changing the name
replaces the container and its data.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"255"}}

### spec.partitionKeyPaths

`[]string` · required

The partition key paths. Each starts with "/" and names a document
property. One path is the normal case (kind HASH); up to three
paths form a HIERARCHICAL key (kind MULTI_HASH, e.g. ["/tenantId",
"/userId"]) that routes queries carrying any prefix of the hierarchy
efficiently. Fixed at creation.

- rule: {"repeated":{"minItems":"1","maxItems":"3","items":{"cel":[{"id":"cosmosdb_sql_container_pk_path_format","message":"partition key paths must start with '/'","expression":"this.startsWith('/')"}]}}}

### spec.partitionKeyKind

`enum`

How the partition key hashes documents to partitions. Unspecified
means HASH (a single path). MULTI_HASH enables the hierarchical
(multi-path) key and requires partition key version 2. Fixed at
creation.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_cosmosdb_sql_container_partition_key_kind_unspecified` -- Not specified: Hash.
- `HASH` -- A single-path hash key -- the normal case. Wire value: "Hash".
- `MULTI_HASH` -- A hierarchical (multi-path) key of up to three levels. Requires partition_key_version 2. Wire value: "MultiHash".

### spec.partitionKeyVersion

`int32` · optional (explicit presence)

The partition key definition version. 1 (the classic default)
caps partition key values at 101 bytes; 2 supports large keys and
is REQUIRED for MULTI_HASH hierarchical keys. Fixed at creation
(Azure tolerates only the unset -> 1 transition in place).

- rule: {"int32":{"lte":2,"gte":1}}

### spec.throughput

`int32` · optional (explicit presence)

Fixed dedicated throughput in RU/s for this container. Minimum 400,
in increments of 100. Mutually exclusive with
`autoscale_max_throughput`; leave both unset to share the database's
throughput (or on serverless accounts).

- rule: throughput must be set in increments of 100 RU/s
- rule: {"int32":{"gte":400}}

### spec.autoscaleMaxThroughput

`int32` · optional (explicit presence)

Autoscale ceiling in RU/s for this container: Azure scales between
10% of this value and this value. Minimum 1000, in increments of
1000. Mutually exclusive with `throughput`.

- rule: autoscale_max_throughput must be set in increments of 1000 RU/s
- rule: {"int32":{"gte":1000}}

### spec.defaultTtl

`int32` · optional (explicit presence)

Default time-to-live for documents, in seconds. -1 turns TTL on
with no default expiry (documents expire only if they carry their
own ttl property); a positive value expires documents that many
seconds after their last write unless overridden per document.
Unset leaves TTL off entirely.

- rule: {"int32":{"gte":-1}}

### spec.analyticalStorageTtl

`int32` · optional (explicit presence)

How long documents stay in the ANALYTICAL store, in seconds
(requires analytical storage on the account). -1 keeps analytical
data forever -- the common choice; a positive value ages it out.
Turning analytical storage OFF on an existing container (unsetting
a previously set value) forces a replacement.

- rule: {"int32":{"gte":-1}}

### spec.uniqueKeys

`[]AzureCosmosdbSqlContainerUniqueKey`

Unique key constraints: within each logical partition, no two
documents may share the same values for a constraint's paths.
Uniqueness is scoped to the partition key, not the whole container.
Fixed at creation.

### spec.uniqueKeys[].paths

`[]string` · required

The document paths whose combined values must be unique within each
logical partition, e.g. ["/email"] or ["/firstName", "/lastName"].

- rule: {"repeated":{"minItems":"1"}}

### spec.indexingPolicy

`AzureCosmosdbSqlContainerIndexingPolicy`

How the container indexes documents. Unset applies Azure's default:
consistent indexing of every path. Tuning the policy -- excluding
bulky payload paths, adding composite indexes for multi-property
ORDER BY -- is the main lever for write RU cost and query
performance, and it updates in place.

- rule: included_paths and excluded_paths must be empty when indexing_mode is NONE
- rule: when included_paths or excluded_paths are declared, exactly one of them must carry the root path '/*'

### spec.indexingPolicy.indexingMode

`enum`

CONSISTENT (the default) indexes synchronously with writes; NONE
disables indexing entirely -- only point reads by id work, which
suits pure key-value containers and bulk-load staging.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_cosmosdb_sql_container_indexing_mode_unspecified` -- Not specified: consistent.
- `CONSISTENT` -- Index synchronously with every write -- Azure's default. Wire value: "consistent".
- `NONE` -- No indexing: only point reads by id. Wire value: "none".

### spec.indexingPolicy.includedPaths

`[]AzureCosmosdbSqlContainerIndexPath`

The paths to index, e.g. "/*" (everything -- Azure's default) or
"/tenantId/?". When any included or excluded path is declared, the
policy replaces Azure's default wholesale, so a tuned policy
normally includes "/*" and excludes the bulky exceptions.

### spec.indexingPolicy.includedPaths[].path

`string` · required

The path, in Cosmos DB index-path syntax: "/*" for everything,
"/prop/?" for a scalar, "/prop/*" for a subtree.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.indexingPolicy.excludedPaths

`[]AzureCosmosdbSqlContainerIndexPath`

The paths NOT to index, e.g. "/payload/*" for a large blob property
that is never filtered on -- the main write-RU saver.

### spec.indexingPolicy.excludedPaths[].path

`string` · required

The path, in Cosmos DB index-path syntax: "/*" for everything,
"/prop/?" for a scalar, "/prop/*" for a subtree.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.indexingPolicy.compositeIndexes

`[]AzureCosmosdbSqlContainerCompositeIndex`

Composite indexes: ordered property sets that make multi-property
ORDER BY (and some multi-filter queries) efficient. Each composite
index lists two or more (path, order) entries.

### spec.indexingPolicy.compositeIndexes[].entries

`[]AzureCosmosdbSqlContainerCompositeIndexEntry` · required

The (path, order) entries composing the index, in significance
order. Cosmos DB requires at least two entries for a composite
index to be meaningful; queries must match the declared orders (or
their exact reverse).

- rule: {"repeated":{"minItems":"1"}}

### spec.indexingPolicy.compositeIndexes[].entries[].path

`string` · required

The document path, e.g. "/lastName". Composite-index paths name a
scalar property directly (no "/?" or "/*" suffix).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.indexingPolicy.compositeIndexes[].entries[].order

`enum`

The sort order this entry is indexed in.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_cosmosdb_sql_container_composite_index_order_unspecified` -- Not specified: ascending.
- `ASCENDING` -- Wire value: "Ascending".
- `DESCENDING` -- Wire value: "Descending".

### spec.indexingPolicy.spatialIndexes

`[]AzureCosmosdbSqlContainerSpatialIndex`

Spatial indexes for GeoJSON geometry queries on the given paths.

### spec.indexingPolicy.spatialIndexes[].path

`string` · required

The path holding GeoJSON geometry to index, e.g. "/location/*".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.conflictResolutionPolicy

`AzureCosmosdbSqlContainerConflictResolutionPolicy`

How conflicting writes are resolved on multi-region-write accounts.
Unset applies Azure's default: last-writer-wins on the document's
_ts timestamp. Fixed at creation.

- rule: conflict_resolution_path applies to LAST_WRITER_WINS and conflict_resolution_procedure to CUSTOM

### spec.conflictResolutionPolicy.mode

`enum` · required

LAST_WRITER_WINS picks the document with the highest value at
`conflict_resolution_path` (default /_ts); CUSTOM hands conflicts to
the stored procedure named in `conflict_resolution_procedure` (or to
the conflict feed when none is registered).

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_cosmosdb_sql_container_conflict_resolution_mode_unspecified`
- `LAST_WRITER_WINS` -- The highest value at conflict_resolution_path wins. Wire value: "LastWriterWins".
- `CUSTOM` -- A stored procedure (or the conflict feed) resolves conflicts. Wire value: "Custom".

### spec.conflictResolutionPolicy.conflictResolutionPath

`string`

For LAST_WRITER_WINS: the numeric document path compared to pick the
winner. Unset means Azure's default, "/_ts" (the write timestamp).

### spec.conflictResolutionPolicy.conflictResolutionProcedure

`string`

For CUSTOM: the stored procedure that resolves conflicts, e.g.
"dbs/{db}/colls/{coll}/sprocs/{name}". Unset routes conflicts to the
conflict feed for the application to resolve.

## Validation Rules

- `cosmosdb_sql_container_throughput_xor`: throughput and autoscale_max_throughput are mutually exclusive
- `cosmosdb_sql_container_multihash_paths`: MULTI_HASH hierarchical keys require partition_key_version 2; HASH takes exactly one path

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureCosmosdbSqlContainer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.sql_container_id` | `string` | The Azure Resource Manager ID of the container -- the management- plane identity and the scope for container-level data-plane RBAC. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.DocumentDB/databaseAccounts/{account}/sqlDatabases/{db}/containers/{name} |
| `status.outputs.sql_container_name` | `string` | The container's name -- what SDK calls reference inside the database. |
| `status.outputs.sql_database_name` | `string` | The name of the database the container lives in, parsed from the resolved database ID. |
| `status.outputs.cosmosdb_account_name` | `string` | The name of the Cosmos DB account, parsed from the resolved database ID -- saves consumers a second reference when they need the full account/database/container triple. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.sqlDatabaseId` | AzureCosmosdbSqlDatabase | `status.outputs.sql_database_id` |

## See Also

- [Overview](../README.md)
