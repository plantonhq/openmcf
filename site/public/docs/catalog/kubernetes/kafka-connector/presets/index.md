---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kafka Connector"
type: "preset-list"
componentSlug: "kafka-connector"
componentTitle: "Kafka Connector"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-first-pipe-mirror-source"
    rank: "01"
    title: "First-pipe preset (stock-image MirrorSource)"
    excerpt: "The smallest real pipe that actually runs on the stock image. The Strimzi Connect image carries ONLY the three MirrorMaker 2 connector classes (MirrorSource, MirrorCheckpoint, MirrorHeartbeat —..."
  - slug: "02-debezium-postgres-cdc"
    rank: "02"
    title: "Debezium Postgres CDC preset"
    excerpt: "A production-shaped change-data-capture source: the Debezium Postgres connector streams row-level changes from a Postgres database into Kafka topics (`orders.public.orders`, ...), using Postgres's..."
  - slug: "03-paused-pipe-with-offsets"
    rank: "03"
    title: "Paused pipe with offsets preset"
    excerpt: "The operational-lifecycle preset: a connector declared `paused` with both offset ConfigMap targets wired. Use this shape when a pipe must exist but not move data yet (a cutover waiting on a..."
---

# Kafka Connector Presets

Ready-to-deploy configuration presets for Kafka Connector. Each preset is a complete manifest you can copy, customize, and deploy.
