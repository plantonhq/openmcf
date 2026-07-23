---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cosmos DB MongoDB Collection"
type: "preset-list"
componentSlug: "cosmos-db-mongodb-collection"
componentTitle: "Cosmos DB MongoDB Collection"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-tenant-sharded"
    rank: "01"
    title: "Tenant-sharded Mongo collection"
    excerpt: "A MongoDB API collection partitioned by `tenantId` with autoscale throughput -- the production-default shape for multi-tenant event or audit streams where each tenant's writes should land on distinct..."
  - slug: "02-shared-throughput"
    rank: "02"
    title: "Shared-throughput Mongo collection"
    excerpt: "A MongoDB API collection that provisions NO throughput of its own -- it draws from the parent database's shared budget (both throughput fields unset). Sharded by `userId` so it still scales out..."
  - slug: "03-ttl-session-store"
    rank: "03"
    title: "TTL session-store Mongo collection"
    excerpt: "A MongoDB API collection where every document expires 24 hours after its last write (`defaultTtlSeconds: 86400` -- Cosmos implements it as an expireAfter index on `_ts`). Sharded by `userId`, with..."
---

# Cosmos DB MongoDB Collection Presets

Ready-to-deploy configuration presets for Cosmos DB MongoDB Collection. Each preset is a complete manifest you can copy, customize, and deploy.
