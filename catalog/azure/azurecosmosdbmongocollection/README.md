# AzureCosmosdbMongoCollection

A MongoDB API collection inside a Cosmos DB database -- the unit of
storage and scale-out where documents actually live. The shard key is
the MongoDB face of the partition key: high cardinality and even request
distribution (tenantId, userId, deviceId) keep the collection scalable;
an unsharded collection is legal but confined to one physical partition.

## When to Use

Use AzureCosmosdbMongoCollection when you need:

- **A document collection with its own scale contract** -- dedicated
  fixed or autoscale RU/s, independent of sibling collections
- **Mongo-style indexes** -- compound keys and unique constraints on
  top of Cosmos DB's Mongo API surface
- **Data lifecycle automation** -- default TTL for expiring documents,
  analytical-store TTL for Synapse Link retention

## Key Configuration

- `mongo_database_id` -- the parent database, referenced from an
  AzureCosmosdbMongoDatabase output; fixed at creation
- `shard_key` -- the document property documents are partitioned by;
  fixed at creation; unset creates an unsharded collection
- `throughput` XOR `autoscale_max_throughput` -- dedicated capacity;
  leave both unset to share the database's throughput (or on serverless
  accounts)
- `indexes` -- Mongo-style index definitions (keys + unique flag);
  Azure requires the `_id` unique index on every collection (the spec
  enforces it up front)
- `default_ttl_seconds` -- document expiry via Cosmos DB's expireAfter
  index on `_ts`; -1 enables TTL without a default expiry

## Composition

```yaml
mongoDatabaseId:
  valueFrom:
    kind: AzureCosmosdbMongoDatabase
    name: app-data
    fieldPath: status.outputs.mongo_database_id
```

Connectivity and keys live on the ACCOUNT (AzureCosmosdbAccount's
endpoint and key outputs); the collection is addressed inside that
connection by database and collection name.

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
