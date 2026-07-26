---
title: "Presets"
description: "Ready-to-deploy configuration presets for SigNoz"
type: "preset-list"
componentSlug: "signoz"
componentTitle: "SigNoz"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev"
    rank: "01"
    title: "SigNoz for development"
    excerpt: "The smallest honest SigNoz: the component's defaults against a composed `KubernetesClickHouse` named `telemetry` in the same namespace — the whole platform (UI, API, alerting, the ingestion..."
  - slug: "02-production"
    rank: "02"
    title: "SigNoz for production"
    excerpt: "The destination posture: SigNoz runs the observability product; a `KubernetesClickHouse` (with its `KubernetesAltinityOperator`) runs the database — each with its own lifecycle, sizing and..."
---

# SigNoz Presets

Ready-to-deploy configuration presets for SigNoz. Each preset is a complete manifest you can copy, customize, and deploy.
