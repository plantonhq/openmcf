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
  - slug: "01-dev-appliance"
    rank: "01"
    title: "Dev appliance"
    excerpt: "The smallest honest SigNoz: the whole platform — UI, ingestion collector, schema migrator and the bundled ClickHouse — from a four-line manifest. The one thing that is NOT upstream-default here is..."
  - slug: "02-production-bundled"
    rank: "02"
    title: "Production on the bundled ClickHouse"
    excerpt: "A production-shaped single-cluster observability platform: replicated ClickHouse (one shard, two replicas over a 3-node ZooKeeper quorum), autoscaled ingestion, alert emails over authenticated SMTP,..."
  - slug: "03-external-clickhouse"
    rank: "03"
    title: "SigNoz on your own ClickHouse"
    excerpt: "The composition posture: SigNoz runs the observability product, a KubernetesClickHouse (with its KubernetesAltinityOperator) runs the database — each with its own lifecycle, sizing and operations...."
---

# SigNoz Presets

Ready-to-deploy configuration presets for SigNoz. Each preset is a complete manifest you can copy, customize, and deploy.
