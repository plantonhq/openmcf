---
title: "Presets"
description: "Ready-to-deploy configuration presets for Container Registry"
type: "preset-list"
componentSlug: "container-registry"
componentTitle: "Container Registry"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard Container Registry"
    excerpt: "This preset creates an Azure Container Registry with Standard SKU and admin user disabled. Standard tier provides 100 GB storage, enhanced throughput for image pulls, and webhook support --..."
  - slug: "02-premium-geo-replicated"
    rank: "02"
    title: "Premium Geo-Replicated Registry"
    excerpt: "This preset creates a Premium registry with a zone-redundant home replica, one geo-replication, and automatic purging of untagged manifests. It is the multi-region production shape: images push once..."
  - slug: "03-premium-network-restricted"
    rank: "03"
    title: "Premium Network-Restricted Registry"
    excerpt: "This preset creates a Premium registry that stays publicly addressable but denies every connection not on an explicit CIDR allowlist, with dedicated data endpoints so downstream egress firewalls can..."
---

# Container Registry Presets

Ready-to-deploy configuration presets for Container Registry. Each preset is a complete manifest you can copy, customize, and deploy.
