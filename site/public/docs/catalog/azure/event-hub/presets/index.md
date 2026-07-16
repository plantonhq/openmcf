---
title: "Presets"
description: "Ready-to-deploy configuration presets for Event Hub"
type: "preset-list"
componentSlug: "event-hub"
componentTitle: "Event Hub"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-telemetry-stream"
    rank: "01"
    title: "Telemetry Stream"
    excerpt: "This preset creates a general-purpose event stream: 8 partitions and a 3-day replay window -- the shape most telemetry, logging, and change-data-capture pipelines start from."
  - slug: "02-captured-archive-stream"
    rank: "02"
    title: "Captured Archive Stream"
    excerpt: "This preset creates a stream with capture enabled: every event is archived to Blob Storage in Avro format automatically -- audit trails, replay beyond the retention window, and batch analytics feed..."
  - slug: "03-compacted-changelog"
    rank: "03"
    title: "Compacted Changelog Stream"
    excerpt: "This preset creates a log-compacted hub: the latest event per partition key is kept forever -- Kafka-style compacted-topic semantics for entity changelogs, materialized views, and cache warming."
---

# Event Hub Presets

Ready-to-deploy configuration presets for Event Hub. Each preset is a complete manifest you can copy, customize, and deploy.
