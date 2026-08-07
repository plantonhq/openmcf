---
title: "Presets"
description: "Ready-to-deploy configuration presets for MSSQL Database"
type: "preset-list"
componentSlug: "mssql-database"
componentTitle: "MSSQL Database"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-general-purpose"
    rank: "01"
    title: "General Purpose Production Database"
    excerpt: "This preset creates a provisioned General Purpose database (2 vCores, 128 GB) with a real backup posture: a 14-day point-in-time restore window plus long-term weekly and monthly backups. Transparent..."
  - slug: "02-serverless-autopause"
    rank: "02"
    title: "Serverless Auto-Pausing Database"
    excerpt: "This preset creates a serverless database that bills compute per second, scales between a warm floor (0.5 vCores) and the SKU ceiling (2 vCores), and pauses entirely after an hour of inactivity --..."
  - slug: "03-hyperscale-replicas"
    rank: "03"
    title: "Hyperscale Database with Readable Replicas"
    excerpt: "This preset creates a Hyperscale database: elastic storage to 100 TB, 4 vCores of independent compute, two readable high-availability replicas spread across availability zones, and zone-redundant..."
---

# MSSQL Database Presets

Ready-to-deploy configuration presets for MSSQL Database. Each preset is a complete manifest you can copy, customize, and deploy.
