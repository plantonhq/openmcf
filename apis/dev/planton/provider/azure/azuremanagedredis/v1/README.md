# AzureManagedRedis

Azure Managed Redis -- Azure's current-generation Redis service, built
on Redis Enterprise (the `Microsoft.Cache/redisEnterprise` ARM family).
Azure is retiring classic Azure Cache for Redis; Managed Redis is the
target for NEW Redis deployments, and it carries the capabilities the
classic service never had: customer-managed-key encryption, active
multi-primary geo-replication, Redis modules, and a keyless-by-default
authentication posture.

## When to Use

Use AzureManagedRedis when you need:

- **Any new Redis on Azure** -- the classic service no longer accepts
  new Premium creations in a growing set of regions
- **Secretless caching** -- access keys are OFF by default; clients
  authenticate with Entra tokens under
  AzureManagedRedisAccessPolicyAssignment grants
- **Redis modules** -- RediSearch, RedisJSON, RedisBloom,
  RedisTimeSeries
- **Write-anywhere global caches** -- active geo-replication via
  AzureManagedRedisGeoReplication
- **Bring-your-own-key encryption** -- CMK against an AzureKeyVaultKey

## Key Configuration

- `cluster_name` -- becomes the endpoint
  `{name}.{region}.redis.azure.net`
- `sku_name` -- tier family and memory size in one value (BALANCED /
  COMPUTE_OPTIMIZED / MEMORY_OPTIMIZED / FLASH_OPTIMIZED); many changes
  apply in place
- `high_availability_enabled` -- a replica and the zone-redundant SLA
  (default true; fixed at creation)
- `default_database` -- the Redis process: authentication, clustering,
  eviction, modules, geo-replication membership, persistence
- `customer_managed_key` + `identity` -- BYO encryption key from Key
  Vault, wrapped by a user-assigned identity

## Composition

```yaml
customerManagedKey:
  keyVaultKeyId:
    valueFrom:
      kind: AzureKeyVaultKey
      name: redis-cmk
      fieldPath: status.outputs.key_id
  userAssignedIdentityId:
    valueFrom:
      kind: AzureUserAssignedIdentity
      name: redis-cmk-identity
      fieldPath: status.outputs.identity_id
```

The `managed_redis_id` output is what the geo-replication and grant
kinds reference; `database_id` is the data-plane scope; `hostname` +
`port` are all a keyless client needs.

## Documentation

- [Design research](docs/README.md) -- field mapping, fold verdict, recorded skips
- [Presets](presets/) -- keyless balanced, search+JSON store, geo-replicated member

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
