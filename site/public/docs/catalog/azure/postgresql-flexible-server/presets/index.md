---
title: "Presets"
description: "Ready-to-deploy configuration presets for PostgreSQL Flexible Server"
type: "preset-list"
componentSlug: "postgresql-flexible-server"
componentTitle: "PostgreSQL Flexible Server"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-dev-burstable"
    rank: "01"
    title: "Development Burstable Server"
    excerpt: "This preset creates the smallest practical PostgreSQL Flexible Server: a single Burstable instance on the public endpoint, one application database, and the Azure-services firewall rule. Version,..."
  - slug: "02-production-private-ha"
    rank: "02"
    title: "Production Private Server with Zone-Redundant HA"
    excerpt: "This preset creates the production baseline: a General Purpose server injected into a delegated subnet (no public endpoint), a zone-redundant standby with automatic failover, geo-redundant 35-day..."
  - slug: "03-hardened-cmk-entra"
    rank: "03"
    title: "Hardened Server: Entra-Only Auth + Customer-Managed Key"
    excerpt: "This preset creates the compliance posture: PostgreSQL password authentication is disabled entirely (no admin password exists to leak, rotate, or audit), administration flows through a Microsoft..."
---

# PostgreSQL Flexible Server Presets

Ready-to-deploy configuration presets for PostgreSQL Flexible Server. Each preset is a complete manifest you can copy, customize, and deploy.
