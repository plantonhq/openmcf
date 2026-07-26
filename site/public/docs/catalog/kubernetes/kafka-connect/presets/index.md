---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kafka Connect"
type: "preset-list"
componentSlug: "kafka-connect"
componentTitle: "Kafka Connect"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-minimal-stock-connect"
    rank: "01"
    title: "Minimal stock Connect preset"
    excerpt: "The smallest declarable Connect cluster: one worker on the stock image against a dev Kafka cluster, no TLS, no authentication. The stock image carries only the MirrorMaker 2 connector classes..."
  - slug: "02-debezium-prebuilt-image"
    rank: "02"
    title: "Debezium prebuilt image preset"
    excerpt: "The prebuilt-image arm in production shape: three workers running an image that already carries the Debezium connector plugins (the `image` reference here is an example — point it at the Debezium..."
  - slug: "03-operator-built-image"
    rank: "03"
    title: "Operator-built image preset"
    excerpt: "The build arm end-to-end: the operator builds a Connect image containing the Debezium Postgres connector (pinned by Maven coordinates, resolved at build time) with Kaniko, pushes it to your registry..."
---

# Kafka Connect Presets

Ready-to-deploy configuration presets for Kafka Connect. Each preset is a complete manifest you can copy, customize, and deploy.
