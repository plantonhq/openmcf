# AzureCosmosdbSqlContainer

A SQL (NoSQL) API container inside a Cosmos DB database -- the unit of
storage, indexing, and scale-out where documents actually live. Azure
distributes documents across physical partitions by the container's
partition key and indexes them by its indexing policy; throughput,
TTL, unique keys, and conflict resolution are set per container.

## When to Use

Use AzureCosmosdbSqlContainer when you need:

- **A document collection with its own scale contract** -- dedicated
  fixed or autoscale RU/s, independent of sibling containers
- **A tuned indexing policy** -- exclude bulky payload paths to cut
  write RU cost, add composite indexes for multi-property ORDER BY
- **Hierarchical partitioning** -- MULTI_HASH keys (e.g.
  /tenantId + /userId) that route prefix queries efficiently
- **Data lifecycle automation** -- default TTL for expiring documents,
  analytical-store TTL for Synapse Link retention

## Key Configuration

- `sql_database_id` -- the parent database, referenced from an
  AzureCosmosdbSqlDatabase output; fixed at creation
- `partition_key_paths` + `partition_key_kind` + `partition_key_version`
  -- the container's most consequential design decision, fixed at
  creation; MULTI_HASH hierarchical keys require version 2
- `throughput` XOR `autoscale_max_throughput` -- dedicated capacity;
  leave both unset to share the database's throughput (or on
  serverless accounts)
- `indexing_policy` -- updatable in place; a declared policy replaces
  Azure's index-everything default wholesale
- `unique_keys` -- per-logical-partition uniqueness, fixed at creation
- `conflict_resolution_policy` -- last-writer-wins or custom, for
  multi-region-write accounts; fixed at creation

## Composition

```yaml
sqlDatabaseId:
  valueFrom:
    kind: AzureCosmosdbSqlDatabase
    name: app-data
    fieldPath: status.outputs.sql_database_id
```

Connectivity and keys live on the ACCOUNT (AzureCosmosdbAccount's
endpoint and key outputs); the container is addressed inside that
connection by database and container name.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
