---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cosmos DB MongoDB Database"
type: "preset-list"
componentSlug: "cosmos-db-mongodb-database"
componentTitle: "Cosmos DB MongoDB Database"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-dedicated-collection-throughput"
    rank: "01"
    title: "Dedicated Collection Throughput"
    excerpt: "This preset creates a database with NO throughput of its own -- a pure namespace. Each AzureCosmosdbMongoCollection inside it provisions its own dedicated RU/s, so workloads are isolated: a hot..."
  - slug: "02-shared-autoscale"
    rank: "02"
    title: "Shared Autoscale Database"
    excerpt: "This preset creates a database with an autoscale throughput budget that every collection inside it SHARES (unless a collection brings its own). Azure scales the budget between 10% and 100% of the..."
  - slug: "03-serverless"
    rank: "03"
    title: "Serverless Database"
    excerpt: "This preset creates a MongoDB API database on a serverless Cosmos DB account -- one whose AzureCosmosdbAccount declares the ENABLE_SERVERLESS capability alongside ENABLE_MONGO. Serverless accounts..."
---

# Cosmos DB MongoDB Database Presets

Ready-to-deploy configuration presets for Cosmos DB MongoDB Database. Each preset is a complete manifest you can copy, customize, and deploy.
