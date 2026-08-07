---
title: "Presets"
description: "Ready-to-deploy configuration presets for MySQL Database System"
type: "preset-list"
componentSlug: "mysql-database-system"
componentTitle: "MySQL Database System"
provider: "oci"
icon: "package"
order: 200
presets:
  - slug: "01-high-availability"
    rank: "01"
    title: "High Availability"
    excerpt: "This preset creates a production MySQL HeatWave DB System with High Availability enabled. Three instances are provisioned across fault domains with automatic failover, PITR-enabled backups with..."
  - slug: "02-standalone-development"
    rank: "02"
    title: "Standalone Development"
    excerpt: "This preset creates a single-instance MySQL HeatWave DB System optimized for development and testing. It uses the smallest available shape, minimal storage, short backup retention, and no HA or..."
---

# MySQL Database System Presets

Ready-to-deploy configuration presets for MySQL Database System. Each preset is a complete manifest you can copy, customize, and deploy.
