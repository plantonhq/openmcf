---
title: "Tenant-sharded Mongo collection"
description: "A MongoDB API collection partitioned by `tenantId` with autoscale throughput -- the production-default shape for multi-tenant event or audit streams where each tenant's writes should land on distinct..."
type: "preset"
rank: "01"
presetSlug: "01-tenant-sharded"
componentSlug: "cosmos-db-mongo-collection"
componentTitle: "Cosmos DB Mongo Collection"
provider: "azure"
icon: "package"
order: 1
---

# Tenant-sharded Mongo collection

A MongoDB API collection partitioned by `tenantId` with autoscale
throughput -- the production-default shape for multi-tenant event or
audit streams where each tenant's writes should land on distinct
physical partitions.

Use when sibling collections in the same database bring their own
dedicated throughput (the parent database preset leaves throughput
unset on purpose).

See [`01-tenant-sharded.yaml`](01-tenant-sharded.yaml).
