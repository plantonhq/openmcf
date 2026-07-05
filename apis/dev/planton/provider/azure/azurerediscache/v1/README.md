# AzureRedisCache

Azure Cache for Redis: a fully managed, in-memory data store on the
open-source Redis engine -- caching, session state, leaderboards, and
pub/sub with sub-millisecond latency. The spec covers the full surface:
tiers and sizing, Entra (token) authentication and the keyless posture,
VNet injection, clustering, RDB/AOF persistence with managed-identity
storage auth, patch windows, and firewall rules.

## When to Use

Use AzureRedisCache when you need:

- **Application caching / session state** -- STANDARD tier, the
  replicated production default
- **High-scale or regulated workloads** -- PREMIUM: clustering,
  persistence, zone pinning, private networking, geo-replication
- **Secretless Redis access** -- Entra auth + access-policy grants, with
  the shared keys turned off entirely

## Key Configuration

- `cache_name` -- globally unique; becomes
  `{cache_name}.redis.cache.windows.net`
- `sku_name` + `capacity` -- the tier and size (C0-C6 / P1-P5); the
  family letter is derived, never spelled
- `access_keys_authentication_enabled` +
  `redis_configuration.active_directory_authentication_enabled` -- the
  auth posture; keys can only go off once Entra is on
- `redis_configuration` -- eviction policy, memory dials, keyspace
  events, RDB/AOF persistence (Premium)
- `subnet_id` / `public_network_access_enabled` -- network isolation
  (VNet injection is legacy; prefer Private Link)

## Composition

```yaml
resourceGroup:
  valueFrom:
    kind: AzureResourceGroup
    name: app-rg
    fieldPath: status.outputs.resource_group_name
```

The `redis_cache_id` output is what the rest of the Redis family
references: AzureRedisLinkedServer (geo-DR), AzureRedisCacheAccessPolicy
+ AzureRedisCacheAccessPolicyAssignment (Entra grants), and
AzurePrivateEndpoint (private connectivity). The `region` output feeds
the linked server's location so geo-replication composes without
hand-repeated strings.

## Documentation

- [Design research](docs/README.md) -- field mapping, fold/split verdicts, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
