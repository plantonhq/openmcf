---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cosmos DB SQL Database"
type: "preset-list"
componentSlug: "cosmos-db-sql-database"
componentTitle: "Cosmos DB SQL Database"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-dedicated-container-throughput"
    rank: "01"
    title: "Dedicated Container Throughput"
    excerpt: "This preset creates a database with NO throughput of its own -- a pure namespace. Each AzureCosmosdbSqlContainer inside it provisions its own dedicated RU/s, so workloads are isolated: a hot..."
  - slug: "02-shared-autoscale"
    rank: "02"
    title: "Shared Autoscale Database"
    excerpt: "This preset creates a database with an autoscale throughput budget that every container inside it SHARES (unless a container brings its own). Azure scales the budget between 10% and 100% of the..."
  - slug: "03-serverless"
    rank: "03"
    title: "Serverless Database"
    excerpt: "This preset creates a database on a serverless Cosmos DB account -- one whose AzureCosmosdbAccount declares the ENABLE_SERVERLESS capability. Serverless accounts bill per request instead of per..."
---

# Cosmos DB SQL Database Presets

Ready-to-deploy configuration presets for Cosmos DB SQL Database. Each preset is a complete manifest you can copy, customize, and deploy.
