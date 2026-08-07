---
title: "Presets"
description: "Ready-to-deploy configuration presets for Tempo"
type: "preset-list"
componentSlug: "tempo"
componentTitle: "Tempo"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single"
    rank: "01"
    title: "Dev single-node Tempo"
    excerpt: "One monolithic Tempo replica on a persistent volume with OTLP receivers. The smallest honest trace store — traces survive pod restarts (unlike the chart's emptyDir default), but storage is local to..."
  - slug: "02-production-s3"
    rank: "02"
    title: "Production Tempo on object storage"
    excerpt: "Trace blocks in an S3-compatible object store (an in-cluster KubernetesSeaweedFs here; AWS S3, GCS or Azure by swapping the storage block), with the metrics generator deriving service-graph and span..."
  - slug: "03-jaeger-compat"
    rank: "03"
    title: "Jaeger-compatible Tempo"
    excerpt: "Tempo with the four legacy Jaeger receiver protocols opened alongside OTLP and the Jaeger-UI-compatible query sidecar enabled — a migration posture for fleets still emitting Jaeger while they move to..."
---

# Tempo Presets

Ready-to-deploy configuration presets for Tempo. Each preset is a complete manifest you can copy, customize, and deploy.
