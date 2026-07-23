---
title: "TTL session-store Mongo collection"
description: "A MongoDB API collection where every document expires 24 hours after its last write (`defaultTtlSeconds: 86400` -- Cosmos implements it as an expireAfter index on `_ts`). Sharded by `userId`, with..."
type: "preset"
rank: "03"
presetSlug: "03-ttl-session-store"
componentSlug: "cosmos-db-mongodb-collection"
componentTitle: "Cosmos DB MongoDB Collection"
provider: "azure"
icon: "package"
order: 3
---

# TTL session-store Mongo collection

A MongoDB API collection where every document expires 24 hours after
its last write (`defaultTtlSeconds: 86400` -- Cosmos implements it as
an expireAfter index on `_ts`). Sharded by `userId`, with fixed
dedicated throughput for a predictable steady-state workload.

Use for session stores, short-lived caches, and any data whose value
decays on a clock -- storage stays flat because deletion is automatic.

See [`03-ttl-session-store.yaml`](03-ttl-session-store.yaml).
