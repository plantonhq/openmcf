---
title: "Search + JSON Document Store"
description: "This preset creates a Managed Redis instance with the RediSearch and RedisJSON modules -- a queryable, full-text-searchable JSON document store with Redis latency. Modules exist only on Managed..."
type: "preset"
rank: "02"
presetSlug: "02-search-json-enterprise"
componentSlug: "managed-redis"
componentTitle: "Managed Redis"
provider: "azure"
icon: "package"
order: 2
---

# Search + JSON Document Store

This preset creates a Managed Redis instance with the RediSearch and
RedisJSON modules -- a queryable, full-text-searchable JSON document
store with Redis latency. Modules exist only on Managed Redis; classic
Azure Cache for Redis never had them.

## When to Use

- Product catalogs, user profiles, and session documents that need
  secondary indexes or full-text search
- Replacing a separate search engine for datasets that already live in
  Redis
- JSON-native applications that want document semantics without a
  second database

## Key Configuration Choices

- **MEMORY_OPTIMIZED_M10** -- search indexes live in memory alongside
  the data, so the memory-heavy family fits; scale the size to your
  dataset plus its indexes
- **ENTERPRISE_CLUSTER clustering** -- required by RediSearch: all
  shards proxied behind a single endpoint, so any Redis client works
- **NO_EVICTION** -- required by RediSearch: an index over
  silently-evicted data would lie; writes fail when memory fills, so
  size generously
- **Keyless** (the default) -- grant consumers through
  AzureManagedRedisAccessPolicyAssignment

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<region>` | Azure region, e.g. `eastus` | Your region strategy (check Managed Redis availability) |
| `<resource-group-resource-name>` | The AzureResourceGroup resource to deploy into | Your resource-group manifest |
