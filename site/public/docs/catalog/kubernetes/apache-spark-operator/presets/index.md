---
title: "Presets"
description: "Ready-to-deploy configuration presets for Apache Spark Operator"
type: "preset-list"
componentSlug: "apache-spark-operator"
componentTitle: "Apache Spark Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-default"
    rank: "01"
    title: "Default preset"
    excerpt: "The standard operator install: the pinned chart (1.8.0 = operator 1.0.0, the official ASF distribution) into its own `spark-operator` namespace, watching every namespace on the cluster, with the..."
  - slug: "02-fenced-team-namespaces"
    rank: "02"
    title: "Fenced team-namespaces preset"
    excerpt: "The multi-tenant posture: Spark runs ONLY in the listed namespaces. The chart CREATES each namespace (know this before pointing at names you manage elsewhere), plants the `spark` service account and..."
---

# Apache Spark Operator Presets

Ready-to-deploy configuration presets for Apache Spark Operator. Each preset is a complete manifest you can copy, customize, and deploy.
