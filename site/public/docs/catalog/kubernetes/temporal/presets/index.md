---
title: "Presets"
description: "Ready-to-deploy configuration presets for Temporal"
type: "preset-list"
componentSlug: "temporal"
componentTitle: "Temporal"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev"
    rank: "01"
    title: "Dev preset"
    excerpt: "The smallest useful Temporal: all four server services, the Web UI, and one Temporal namespace (`default`) — against a composed KubernetesPostgres named `temporal-db` in the same Kubernetes..."
  - slug: "02-production"
    rank: "02"
    title: "Production preset"
    excerpt: "Temporal sized for real workloads on the composed-PostgreSQL story: replicated frontend/matching/worker, three history replicas (shard ownership redistributes across them — history is where workflow..."
---

# Temporal Presets

Ready-to-deploy configuration presets for Temporal. Each preset is a complete manifest you can copy, customize, and deploy.
