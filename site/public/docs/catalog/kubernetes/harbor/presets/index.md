---
title: "Presets"
description: "Ready-to-deploy configuration presets for Harbor"
type: "preset-list"
componentSlug: "harbor"
componentTitle: "Harbor"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-minimal"
    rank: "01"
    title: "Minimal — evaluation registry, zero dependencies"
    excerpt: "The smallest honest Harbor: the chart's in-cluster PostgreSQL and Redis (single-node, evaluation-grade by upstream's own position), artifact blobs on a 20Gi PersistentVolumeClaim, Trivy scanning on..."
  - slug: "02-production-composed"
    rank: "02"
    title: "Production — composed data plane, object storage, HA components"
    excerpt: "The production posture: every stateful concern leaves the chart's evaluation-grade internals and composes the catalog's own kinds — PostgreSQL from a KubernetesPostgres (the operator-maintained..."
---

# Harbor Presets

Ready-to-deploy configuration presets for Harbor. Each preset is a complete manifest you can copy, customize, and deploy.
