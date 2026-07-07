---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cosmos DB Account"
type: "preset-list"
componentSlug: "cosmos-db-account"
componentTitle: "Cosmos DB Account"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-sql-api"
    rank: "01"
    title: "Production SQL API Account"
    excerpt: "This preset creates a two-region SQL (NoSQL) API account with automatic failover and continuous backup -- the production baseline for document workloads. Databases and containers are deployed as..."
  - slug: "02-mongodb-api"
    rank: "02"
    title: "MongoDB API Account"
    excerpt: "This preset creates a MongoDB-compatible account: existing MongoDB drivers, tools, and code work unchanged against a fully managed, globally distributable backend. Databases and collections are..."
  - slug: "03-serverless"
    rank: "03"
    title: "Serverless, Entra-Only Account"
    excerpt: "This preset creates a serverless SQL-API account with the fully locked-down access posture: pay-per-request billing, no public endpoint, and no key-based authentication -- every data-plane caller..."
---

# Cosmos DB Account Presets

Ready-to-deploy configuration presets for Cosmos DB Account. Each preset is a complete manifest you can copy, customize, and deploy.
