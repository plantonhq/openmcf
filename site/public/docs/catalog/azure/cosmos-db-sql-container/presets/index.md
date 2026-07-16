---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cosmos DB SQL Container"
type: "preset-list"
componentSlug: "cosmos-db-sql-container"
componentTitle: "Cosmos DB SQL Container"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-tenant-partitioned"
    rank: "01"
    title: "Tenant-Partitioned Container"
    excerpt: "This preset creates the workhorse production container: a single high-cardinality partition key, autoscale throughput that follows the traffic curve, and an indexing policy that stops paying write RU..."
  - slug: "02-hierarchical-key"
    rank: "02"
    title: "Hierarchical Partition Key"
    excerpt: "This preset creates a container with a two-level MultiHash partition key (/tenantId, /userId) and a composite index for time-ordered queries -- the shape for tenant data that outgrows a single..."
  - slug: "03-ttl-session-store"
    rank: "03"
    title: "TTL Session Store"
    excerpt: "This preset creates a self-cleaning key-value container: documents expire automatically after 24 hours, indexing is off (point reads only), and throughput is shared from the database -- the cheapest..."
---

# Cosmos DB SQL Container Presets

Ready-to-deploy configuration presets for Cosmos DB SQL Container. Each preset is a complete manifest you can copy, customize, and deploy.
