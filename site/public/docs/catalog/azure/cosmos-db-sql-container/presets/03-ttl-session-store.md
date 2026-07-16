---
title: "TTL Session Store"
description: "This preset creates a self-cleaning key-value container: documents expire automatically after 24 hours, indexing is off (point reads only), and throughput is shared from the database -- the cheapest..."
type: "preset"
rank: "03"
presetSlug: "03-ttl-session-store"
componentSlug: "cosmos-db-sql-container"
componentTitle: "Cosmos DB SQL Container"
provider: "azure"
icon: "package"
order: 3
---

# TTL Session Store

This preset creates a self-cleaning key-value container: documents
expire automatically after 24 hours, indexing is off (point reads
only), and throughput is shared from the database -- the cheapest
possible shape for ephemeral data.

## When to Use

- Session state, short-lived tokens, or cache-like documents with a
  natural expiry
- Pure key-value access (reads by id) where queries never filter on
  document properties
- Small auxiliary containers that should ride the database's shared
  throughput instead of billing their own

## Key Configuration Choices

- **`defaultTtl: 86400`** -- documents expire a day after their last
  write; a per-document `ttl` property overrides the default. Set `-1`
  to turn TTL on with per-document expiry only
- **`indexingMode: NONE`** -- only point reads by id work; writes stop
  paying indexing RU entirely. Wrong for anything that queries
- **No throughput fields** -- the container shares the database's
  provisioned throughput (also the required shape on serverless
  accounts)
- **`uniqueKeys`** -- one active session per user within the logical
  partition; fixed at creation

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-data` | The AzureCosmosdbSqlDatabase's Planton resource name | Your Cosmos composition |
| `sessions` | The container name | Your data taxonomy |
