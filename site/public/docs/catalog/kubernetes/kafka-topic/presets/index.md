---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kafka Topic"
type: "preset-list"
componentSlug: "kafka-topic"
componentTitle: "Kafka Topic"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-simple-event-stream"
    rank: "01"
    title: "Simple Event Stream"
    excerpt: "This preset declares the standard application event topic: 6 partitions, replication factor 3, messages retained for 7 days and then deleted. It is the shape most topics should start from — durable..."
  - slug: "02-compacted-changelog"
    rank: "02"
    title: "Compacted Changelog"
    excerpt: "This preset declares a compacted topic: instead of deleting messages by age, Kafka's log cleaner retains the LATEST value for each message key and discards older versions. The topic behaves as a..."
  - slug: "03-high-throughput"
    rank: "03"
    title: "High Throughput"
    excerpt: "This preset declares a wide, short-lived firehose topic: 24 partitions for consumer parallelism, 1 GiB segments so high-volume partitions roll less often, and a 2-day retention window because..."
---

# Kafka Topic Presets

Ready-to-deploy configuration presets for Kafka Topic. Each preset is a complete manifest you can copy, customize, and deploy.
