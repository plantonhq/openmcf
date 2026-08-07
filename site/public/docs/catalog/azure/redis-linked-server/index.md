---
title: "Redis Linked Server"
description: "Redis Linked Server deployment documentation"
icon: "package"
order: 100
componentName: "azureredislinkedserver"
---

# Azure Redis Linked Server

Links two PREMIUM Azure Cache for Redis instances into a geo-replication pair: the primary serves reads and writes while continuously replicating to the secondary in another region, which serves as the warm disaster-recovery target. The link is a first-class resource because DELETING it IS the failover operation -- unlinking makes the secondary writable.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Linked Server** -- the geo-replication link, created as a child of the primary cache and named after the secondary

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

Azure enforces three contracts at link time:

- **Both caches must be PREMIUM tier** -- geo-replication is a Premium feature on both ends.
- **The caches must live in DIFFERENT regions** -- that is the point of the pair.
- **The secondary must be the same size or larger** than the primary. Its existing data is flushed when the link is established.

## Deploy

### Console

Open the deployment store, find **Azure Redis Linked Server**, and click **Deploy**. Start from the **Geo DR Link** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRedisLinkedServer
metadata:
  name: prod-cache-geo-link
  org: acme-corp
  env: prod
spec:
  targetRedisCacheId:
    valueFrom:
      kind: AzureRedisCache
      name: primary-cache
      fieldPath: status.outputs.redis_cache_id
  linkedRedisCacheId:
    valueFrom:
      kind: AzureRedisCache
      name: dr-cache
      fieldPath: status.outputs.redis_cache_id
  linkedRedisCacheLocation:
    valueFrom:
      kind: AzureRedisCache
      name: dr-cache
      fieldPath: status.outputs.region
  serverRole: SECONDARY
```

```shell
planton apply -f geo-link.yaml
```

## Key Configuration

**The pair** -- `targetRedisCacheId` is the PRIMARY (the link is created as its child; its resource group and name derive from this ID). `linkedRedisCacheId` is the SECONDARY -- the DR replica that rejects writes while linked.

**The derived location** -- `linkedRedisCacheLocation` references the SAME cache as the secondary (its `region` output), so the location derives from one source of truth instead of being hand-repeated. A literal region is the escape hatch for caches managed outside the manifest set.

**The role** -- `SECONDARY` is the normal shape; `PRIMARY` inverts the pair and is rarely used outside re-linking after a failover.

**Failover** -- Delete the link (the secondary becomes writable), move traffic, and create a new link in the opposite direction once the region recovers. Applications that point at the `geo_replicated_primary_host_name` output instead of either cache's hostname ride through failovers without a connection-string change.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureRedisCache** (primary) | `targetRedisCacheId` | `status.outputs.redis_cache_id` |
| **AzureRedisCache** (secondary) | `linkedRedisCacheId` | `status.outputs.redis_cache_id` |
| **AzureRedisCache** (secondary) | `linkedRedisCacheLocation` | `status.outputs.region` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `linked_server_id` | Azure Resource Manager ID of the link | Audit trails |
| `linked_server_name` | The link's name (equals the secondary cache's name) | Operational tooling |
| `geo_replicated_primary_host_name` | DNS name that always resolves to the CURRENT primary | Application connection strings that survive failovers |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Geo DR link** -- The normal shape: primary in the app's region, secondary in the DR region, role SECONDARY. Start from the **Geo DR Link** preset.

**Re-link after failover** -- Once the failed region recovers, link in the opposite direction. Start from the **Relink After Failover** preset.

**Cross-manifest link** -- Literal ARM IDs and a literal region when the caches are managed outside this manifest set. Start from the **Cross-Manifest Link** preset.

## Works With

- [**Azure Redis Cache**](/cloud-catalog/azure-redis-cache) -- both ends of the pair (Premium tier)
