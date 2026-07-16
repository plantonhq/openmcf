---
title: "Presets"
description: "Ready-to-deploy configuration presets for SQL Elastic Pool"
type: "preset-list"
componentSlug: "sql-elastic-pool"
componentTitle: "SQL Elastic Pool"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-dtu"
    rank: "01"
    title: "Standard DTU Pool"
    excerpt: "This preset creates a 100-eDTU Standard elastic pool -- the classic shared-compute container for many small databases whose usage peaks do not overlap (the SaaS tenant-per-database pattern)."
  - slug: "02-general-purpose-vcore"
    rank: "02"
    title: "General Purpose vCore Pool with Hybrid Benefit"
    excerpt: "This preset creates a 4-vCore General Purpose elastic pool with a 512 GB shared storage cap, fractional per-database bounds, and Azure Hybrid Benefit licensing (bring your own SQL Server license)."
  - slug: "03-business-critical-zr"
    rank: "03"
    title: "Zone-Redundant Business Critical Pool"
    excerpt: "This preset creates a 4-vCore Business Critical elastic pool with zone-redundant replicas: local-SSD storage for the lowest I/O latency, built-in synchronous replicas for the fastest failover, spread..."
---

# SQL Elastic Pool Presets

Ready-to-deploy configuration presets for SQL Elastic Pool. Each preset is a complete manifest you can copy, customize, and deploy.
