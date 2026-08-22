# Redis Cache

This preset creates a single-node Redis 7 cluster for caching workloads, with an LRU eviction policy and VPC-private access on the smallest managed-database size.

## When to Use

- Application-level caching (sessions, computed results, hot lookups)
- Workloads that tolerate key eviction under memory pressure
- Private-network caching for apps already running inside a VPC

## Key Configuration Choices

- **Redis 7** (`engine: redis`, `engineVersion: "7"`) -- DigitalOcean treats Redis and Valkey as interchangeable caching engines; switch to `engine: valkey` for the open-source-governed engine without any other change.
- **LRU eviction** (`evictionPolicy: allkeys_lru`) -- evicts the least-recently-used keys when memory fills, the right default for pure caches. Eviction policies apply only to Redis/Valkey clusters.
- **Single node** (`nodeCount: 1`) -- caches usually tolerate a brief failover gap; add standbys only when cache warm-up is expensive.
- **VPC placement** (`vpc.valueFrom`) -- references a `DigitalOceanVpc` resource named `my-vpc`; rename it to your VPC resource, or replace the block with `value: <uuid>` for an unmanaged VPC.

## Related Presets

- **01-postgresql-ha** -- Use for durable relational data
- **02-postgresql-dev** -- Use for dev/test relational databases
