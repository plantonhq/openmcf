# AzureCosmosdbMongoCollection

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureCosmosdbMongoCollectionSpec** defines the configuration for
creating a MongoDB API collection inside a Cosmos DB database: the unit
of storage and scale-out where documents actually live.

**The shard key is the MongoDB face of the partition key** and carries
the same design weight: a property with high cardinality and even
request distribution (tenantId, userId, deviceId). An unsharded
collection is legal but confined to a single physical partition --
fine for small lookup collections, wrong for anything that grows.
The shard key is fixed at creation.

**Throughput**: a collection either shares its database's provisioned
throughput (leave both fields unset), or brings its own dedicated
fixed/autoscale throughput. On serverless accounts neither may be set.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureCosmosdbMongoCollection
metadata:
  name: test-cosmos-mongo-collection
spec:
  mongoDatabaseId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DocumentDB/databaseAccounts/planton-hack-cosmos-mongo/mongodbDatabases/app-data
  collectionName: events
  shardKey: tenantId
  # Exercises the autoscale seam (mutually exclusive with fixed
  # throughput -- the spec rejects both together).
  autoscaleMaxThroughput: 1000
  # TTL on with per-document expiry only (never 0 -- the spec rejects
  # it, matching ARM).
  defaultTtlSeconds: -1
  # Exercises the index block rendering, including the _id unique index
  # Azure requires on every Mongo collection.
  indexes:
    - keys:
        - _id
      unique: true
    - keys:
        - tenantId
        - createdAt
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.mongoDatabaseId` | `string \| valueFrom` | yes |  | AzureCosmosdbMongoDatabase (`status.outputs.mongo_database_id`) |
| `spec.collectionName` | `string` | yes |  |  |
| `spec.shardKey` | `string` |  |  |  |
| `spec.throughput` | `int32` |  |  |  |
| `spec.autoscaleMaxThroughput` | `int32` |  |  |  |
| `spec.defaultTtlSeconds` | `int32` |  |  |  |
| `spec.analyticalStorageTtl` | `int32` |  |  |  |
| `spec.indexes` | `[]AzureCosmosdbMongoCollectionIndex` |  |  |  |
| `spec.indexes[].keys` | `[]string` | yes |  |  |
| `spec.indexes[].unique` | `bool` |  |  |  |

## Field Details

### spec.mongoDatabaseId

`string | valueFrom` · required

The Mongo database the collection lives in, by ARM ID. References an
AzureCosmosdbMongoDatabase's mongo_database_id output so the
account, database, and collections compose in one manifest set.
Fixed at creation.

- references: AzureCosmosdbMongoDatabase (`status.outputs.mongo_database_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureCosmosdbMongoDatabase, name: <that resource's name>, fieldPath: status.outputs.mongo_database_id}} -- a bare string does not parse

### spec.collectionName

`string` · required

The collection's name -- unique within the database, 1-255
characters (Azure's only constraint on Cosmos entity names).
Changing the name replaces the collection and its data.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"255"}}

### spec.shardKey

`string`

The document property documents are sharded (partitioned) by, e.g.
"tenantId". Unset creates an UNSHARDED collection confined to one
physical partition -- acceptable only for small, bounded
collections. Fixed at creation: resharding means a new collection
and a data migration.

### spec.throughput

`int32` · optional (explicit presence)

Fixed dedicated throughput in RU/s for this collection. Minimum
400, in increments of 100. Mutually exclusive with
`autoscale_max_throughput`; leave both unset to share the
database's throughput (or on serverless accounts).

- rule: throughput must be set in increments of 100 RU/s
- rule: {"int32":{"gte":400}}

### spec.autoscaleMaxThroughput

`int32` · optional (explicit presence)

Autoscale ceiling in RU/s for this collection: Azure scales between
10% of this value and this value. Minimum 1000, in increments of
1000. Mutually exclusive with `throughput`.

- rule: autoscale_max_throughput must be set in increments of 1000 RU/s
- rule: {"int32":{"gte":1000}}

### spec.defaultTtlSeconds

`int32` · optional (explicit presence)

Default time-to-live for documents, in seconds (implemented by
Cosmos DB as an expireAfter index on _ts). -1 turns TTL on with no
default expiry (documents expire only if they carry their own ttl);
a positive value expires documents that many seconds after their
last write. 0 is not a valid value; unset leaves TTL off.

- rule: default_ttl_seconds must be -1 (on, no default expiry) or a positive number of seconds

### spec.analyticalStorageTtl

`int32` · optional (explicit presence)

How long documents stay in the ANALYTICAL store, in seconds
(requires analytical storage on the account). -1 keeps analytical
data forever; a positive value ages it out.

- rule: {"int32":{"gte":-1}}

### spec.indexes

`[]AzureCosmosdbMongoCollectionIndex`

Indexes on the collection. Azure REQUIRES an index on ["_id"]
(declare it with unique: true -- the _id index is always unique) --
a collection cannot be created or updated without it, so it is part
of the collection's real configuration, never injected silently.
Compound indexes list multiple keys; `unique` enforces distinct
values across the indexed keys.

### spec.indexes[].keys

`[]string` · required

The document properties the index covers, in order. A single key
for a simple index (e.g. ["_id"]); multiple keys for a compound
index (e.g. ["tenantId", "createdAt"]).

- rule: {"repeated":{"minItems":"1"}}

### spec.indexes[].unique

`bool` · optional (explicit presence)

Whether the index enforces uniqueness across its keys. The ["_id"]
index is always unique.

## Validation Rules

- `cosmosdb_mongo_coll_throughput_xor`: throughput and autoscale_max_throughput are mutually exclusive
- `cosmosdb_mongo_coll_id_index_required`: indexes must include an index on the '_id' key -- Azure rejects a Mongo collection without it

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureCosmosdbMongoCollection, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.mongo_collection_id` | `string` | The Azure Resource Manager ID of the collection -- the management- plane identity ARM reads and policy target. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.DocumentDB/databaseAccounts/{account}/mongodbDatabases/{db}/collections/{name} |
| `status.outputs.mongo_collection_name` | `string` | The collection's name -- what MongoDB drivers reference inside the database. |
| `status.outputs.mongo_database_name` | `string` | The name of the database the collection lives in, parsed from the resolved database ID. |
| `status.outputs.cosmosdb_account_name` | `string` | The name of the Cosmos DB account, parsed from the resolved database ID -- saves consumers a second reference when they need the full account/database/collection triple. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.mongoDatabaseId` | AzureCosmosdbMongoDatabase | `status.outputs.mongo_database_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
