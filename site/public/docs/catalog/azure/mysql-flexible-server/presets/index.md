---
title: "Presets"
description: "Ready-to-deploy configuration presets for MySQL Flexible Server"
type: "preset-list"
componentSlug: "mysql-flexible-server"
componentTitle: "MySQL Flexible Server"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-dev-burstable"
    rank: "01"
    title: "Development Burstable Server"
    excerpt: "This preset creates the smallest practical MySQL Flexible Server: a single Burstable instance on the public endpoint, one application database, and the Azure-services firewall rule. Version (8.0.21),..."
  - slug: "02-production-private-ha"
    rank: "02"
    title: "Production Private Server with Zone-Redundant HA"
    excerpt: "This preset creates a production-grade MySQL Flexible Server: General Purpose compute, 256 GiB storage with auto-grow and elastic IOPS scaling, a zone-redundant standby with automatic failover,..."
  - slug: "03-hardened-cmk-entra"
    rank: "03"
    title: "Hardened Server with CMK Encryption and an Entra Administrator"
    excerpt: "This preset creates a compliance-oriented MySQL Flexible Server: customer-managed-key (CMK) encryption unwrapped through a user-assigned identity, a Microsoft Entra administrator for token-based..."
---

# MySQL Flexible Server Presets

Ready-to-deploy configuration presets for MySQL Flexible Server. Each preset is a complete manifest you can copy, customize, and deploy.
