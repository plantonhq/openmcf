---
title: "Presets"
description: "Ready-to-deploy configuration presets for RabbitMQ Cluster Operator"
type: "preset-list"
componentSlug: "rabbitmq-cluster-operator"
componentTitle: "RabbitMQ Cluster Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard preset"
    excerpt: "The one-per-cluster engine install — and the spec is empty, because the release manifest's own defaults ARE the production-standard posture: the operator watching ALL namespaces (the upstream..."
  - slug: "02-namespace-fenced"
    rank: "02"
    title: "Namespace-fenced preset"
    excerpt: "The multi-tenant posture: the operator watches only the namespaces that are allowed to hold RabbitMQ clusters. The upstream default is the opposite — empty `watch_namespaces` watches EVERYTHING..."
  - slug: "03-air-gapped-mirror"
    rank: "03"
    title: "Air-gapped mirror preset"
    excerpt: "The air-gap posture: every image this install — and the clusters it will create — pulls, re-pointed at a private registry. Three image surfaces travel together:"
---

# RabbitMQ Cluster Operator Presets

Ready-to-deploy configuration presets for RabbitMQ Cluster Operator. Each preset is a complete manifest you can copy, customize, and deploy.
