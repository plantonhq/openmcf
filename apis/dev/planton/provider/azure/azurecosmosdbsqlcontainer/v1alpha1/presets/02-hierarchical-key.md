# Hierarchical Partition Key

This preset creates a container with a two-level MultiHash partition
key (/tenantId, /userId) and a composite index for time-ordered
queries -- the shape for tenant data that outgrows a single partition
key value.

## When to Use

- Multi-tenant workloads where a single tenant's data can exceed the
  20 GB per-partition-key-value limit of a simple key
- Query patterns that filter by tenant alone AND by tenant+user --
  prefix queries route efficiently at every level of the hierarchy
- Feeds and timelines needing "newest first" per tenant without
  per-document sorts

## Key Configuration Choices

- **`partitionKeyKind: MULTI_HASH` + two paths** -- up to three levels;
  fixed at creation
- **`partitionKeyVersion: 2`** -- required for hierarchical keys
  (large-key support); the spec enforces the pairing
- **`throughput: 2000`** -- fixed RU/s for steady ingest; swap for
  `autoscaleMaxThroughput` on spiky traffic
- **The composite index** matches the query's ORDER BY (tenantId ASC,
  createdAt DESC); composite indexes serve the declared order or its
  exact reverse

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-data` | The AzureCosmosdbSqlDatabase's Planton resource name | Your Cosmos composition |
| `user-events` | The container name | Your data taxonomy |
