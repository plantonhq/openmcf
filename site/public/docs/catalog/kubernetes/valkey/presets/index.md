---
title: "Presets"
description: "Ready-to-deploy configuration presets for Valkey"
type: "preset-list"
componentSlug: "valkey"
componentTitle: "Valkey"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-single-instance"
    rank: "01"
    title: "Single Instance Valkey"
    excerpt: "This preset deploys one Valkey instance (Redis-compatible) with append-only persistence on a 1Gi volume, a memory ceiling with LRU eviction, and ACL authentication. The most common shape for caching..."
  - slug: "02-persistent-with-replicas"
    rank: "02"
    title: "Production Valkey with Replicas"
    excerpt: "This preset deploys a primary plus two replicas (Redis-compatible) with per-pod append-only persistence, write safety, a dedicated read Service, a PodDisruptionBudget, and ACL authentication. Read..."
---

# Valkey Presets

Ready-to-deploy configuration presets for Valkey. Each preset is a complete manifest you can copy, customize, and deploy.
