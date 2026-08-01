---
title: "Presets"
description: "Ready-to-deploy configuration presets for Apache Flink Operator"
type: "preset-list"
componentSlug: "apache-flink-operator"
componentTitle: "Apache Flink Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-default"
    rank: "01"
    title: "Default preset"
    excerpt: "The standard operator install: the pinned version (1.15.0, the signed Apache per-version chart channel) into its own `flink-system` namespace, watching every namespace, the admission webhook ON."
  - slug: "02-fenced-ha"
    rank: "02"
    title: "Fenced HA preset"
    excerpt: "The multi-tenant, standby-backed posture: the operator watches ONLY the listed namespaces — the chart scopes its RBAC AND the admission webhook's namespaceSelector to exactly this list, so Flink..."
---

# Apache Flink Operator Presets

Ready-to-deploy configuration presets for Apache Flink Operator. Each preset is a complete manifest you can copy, customize, and deploy.
