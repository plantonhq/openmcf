---
title: "Premium Enterprise Cache"
description: "This preset creates a private, zone-spread, clustered Premium cache with hourly RDB persistence authenticated by managed identity -- the shape a regulated or high-scale workload reaches for."
type: "preset"
rank: "02"
presetSlug: "02-premium-enterprise"
componentSlug: "cache-for-redis"
componentTitle: "Cache for Redis"
provider: "azure"
icon: "package"
order: 2
---

# Premium Enterprise Cache

This preset creates a private, zone-spread, clustered Premium cache with
hourly RDB persistence authenticated by managed identity -- the shape a
regulated or high-scale workload reaches for.

## When to Use

- Data sets beyond a single node's memory (clustering multiplies both
  memory and throughput)
- Workloads that must survive a full cache restart (RDB persistence
  rebuilds the data set)
- Environments where the cache must never answer on a public endpoint

## Key Configuration Choices

- **3 shards at P1** -- 18 GB total; clients must speak the Redis
  Cluster protocol. Clustering cannot be combined with extra
  `replicasPerPrimary` (each shard already has its replica)
- **Private-only** -- public access off; attach an AzurePrivateEndpoint
  for connectivity. VNet injection (`subnetId`) is the legacy
  alternative; Private Link is the recommendation for new designs
- **Managed-identity persistence** -- the system-assigned identity
  writes snapshots, so no storage connection string exists in the spec;
  grant the exported `identity_principal_id` the "Storage Blob Data
  Contributor" role on the persistence account
- **`noeviction`** -- writes fail rather than silently dropping keys;
  watch memory and scale capacity before it fills
- **Geo-DR** -- pair with a second Premium cache in another region
  through AzureRedisLinkedServer when regional failover is required

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<region>` | Azure region, e.g. `eastus` | Your region strategy |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your foundation composition |
| `myorg-prod-cache` | Globally unique DNS name (1-63 letters/digits/hyphens) | Becomes `{cacheName}.redis.cache.windows.net` |
