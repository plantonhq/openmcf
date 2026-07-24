---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kafka UI"
type: "preset-list"
componentSlug: "kafka-ui"
componentTitle: "Kafka UI"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-single-cluster-readonly"
    rank: "01"
    title: "Single cluster readonly preset"
    excerpt: "The safe first console: one Kafka cluster wired in with `read_only` on, so the console can browse topics, messages, and consumer lag but cannot create, delete, produce, or edit anything — an app-side..."
  - slug: "02-full-stack-console"
    rank: "02"
    title: "Full stack console preset"
    excerpt: "The whole Kafka family in one pane: a TLS + SCRAM cluster connection, schema browsing through the registry, Connect pipe monitoring, and a login gate on the console itself. This is the shape an infra..."
  - slug: "03-multi-cluster"
    rank: "03"
    title: "Multi cluster preset"
    excerpt: "One console for the whole estate: staging with full console powers (create topics, produce test messages, edit configs) and production locked to observe-only — the per-cluster `read_only` switch is..."
---

# Kafka UI Presets

Ready-to-deploy configuration presets for Kafka UI. Each preset is a complete manifest you can copy, customize, and deploy.
