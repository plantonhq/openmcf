---
title: "Presets"
description: "Ready-to-deploy configuration presets for Strimzi Kafka Operator"
type: "preset-list"
componentSlug: "strimzi-kafka-operator"
componentTitle: "Strimzi Kafka Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-default"
    rank: "01"
    title: "Default"
    excerpt: "This preset installs the Strimzi cluster operator in its standard posture: the pinned `strimzi-kafka-operator` chart, own-namespace watch scope, and the chart's own defaults for everything else..."
  - slug: "02-cluster-wide"
    rank: "02"
    title: "Cluster-Wide"
    excerpt: "This preset installs ONE Strimzi cluster operator that manages Kafka clusters in EVERY namespace — the platform-team posture: the operator lives in a dedicated control-plane namespace..."
  - slug: "03-fenced-teams"
    rank: "03"
    title: "Fenced Teams"
    excerpt: "This preset installs one Strimzi cluster operator with an EXPLICIT namespace fence: it reconciles Kafka clusters only in the listed team namespaces (plus its own installation namespace, which is..."
---

# Strimzi Kafka Operator Presets

Ready-to-deploy configuration presets for Strimzi Kafka Operator. Each preset is a complete manifest you can copy, customize, and deploy.
