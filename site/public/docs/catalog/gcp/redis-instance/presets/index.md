---
title: "Presets"
description: "Ready-to-deploy configuration presets for Redis Instance"
type: "preset-list"
componentSlug: "redis-instance"
componentTitle: "Redis Instance"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-basic-cache"
    rank: "01"
    title: "Basic Cache"
    excerpt: "This preset deploys a single-node BASIC tier Redis instance — the smallest, cheapest Memorystore shape — peered to your VPC over direct peering."
  - slug: "02-ha-production"
    rank: "02"
    title: "Production HA Cache"
    excerpt: "This preset deploys a STANDARD_HA Redis instance — a primary with an automatic-failover replica in a second zone (99.9% SLA) — hardened with AUTH, TLS, RDB persistence, and a pinned maintenance..."
  - slug: "03-private-services-access"
    rank: "03"
    title: "Private Services Access with Read Replicas"
    excerpt: "This preset deploys a STANDARD_HA Redis instance over the VPC's private services access connection — the connectivity mode Shared VPC requires — with a read endpoint for scale-out, CMEK encryption at..."
---

# Redis Instance Presets

Ready-to-deploy configuration presets for Redis Instance. Each preset is a complete manifest you can copy, customize, and deploy.
