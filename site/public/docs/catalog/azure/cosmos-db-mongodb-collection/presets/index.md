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
---

# Cosmos DB MongoDB Collection Presets

Ready-to-deploy configuration presets for Cosmos DB MongoDB Collection. Each preset is a complete manifest you can copy, customize, and deploy.
