# AzureRedisLinkedServer

A geo-replication link between two Premium Azure Cache for Redis
instances: the primary serves writes while continuously replicating to
a warm secondary in another region. Deleting the link IS the failover
operation -- it promotes the secondary to a writable cache.

## When to Use

Use AzureRedisLinkedServer when you need:

- **Regional disaster recovery** for caches whose data is expensive or
  slow to rebuild
- **Read locality** -- the linked secondary serves reads in its own
  region
- **A failover-stable endpoint** -- the
  `geo_replicated_primary_host_name` output always resolves to the
  current primary

## Key Configuration

- `target_redis_cache_id` -- the PRIMARY cache (the link's parent; its
  name and resource group derive from this reference)
- `linked_redis_cache_id` -- the SECONDARY cache in another region
- `linked_redis_cache_location` -- referenced from the same secondary
  cache's `region` output, so it can never disagree
- `server_role` -- SECONDARY for the normal DR shape

Azure's requirements: both caches PREMIUM, different regions, and the
secondary at least as large as the primary.

## Composition

```yaml
targetRedisCacheId:
  valueFrom:
    kind: AzureRedisCache
    name: app-cache-east
    fieldPath: status.outputs.redis_cache_id
linkedRedisCacheLocation:
  valueFrom:
    kind: AzureRedisCache
    name: app-cache-west
    fieldPath: status.outputs.region
```

## Documentation

- [Design research](docs/README.md) -- field mapping, the split verdict, failover semantics
- [Presets](presets/) -- DR link, re-link after failover, cross-manifest link
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
