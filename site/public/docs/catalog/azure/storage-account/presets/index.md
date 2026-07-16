---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Account"
type: "preset-list"
componentSlug: "storage-account"
componentTitle: "Storage Account"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-general-purpose-v2"
    rank: "01"
    title: "General-Purpose v2 Account"
    excerpt: "This preset creates a StorageV2 account on the standard tier with local redundancy and blob data protection -- the baseline for application assets, uploads, and scratch data. Containers are added as..."
  - slug: "02-production-locked-down"
    rank: "02"
    title: "Production Locked-Down Account"
    excerpt: "This preset creates a geo-zone-redundant account with the production security posture: a DENY firewall admitting only declared subnets and trusted Microsoft services, anonymous access made..."
  - slug: "03-data-lake-gen2"
    rank: "03"
    title: "Data Lake Storage Gen2 Account"
    excerpt: "This preset creates a hierarchical-namespace (ADLS Gen2) account: real directories, POSIX ACLs, and the dfs endpoint that analytics engines -- Spark, Databricks, Synapse -- address. SFTP ingestion is..."
---

# Storage Account Presets

Ready-to-deploy configuration presets for Storage Account. Each preset is a complete manifest you can copy, customize, and deploy.
