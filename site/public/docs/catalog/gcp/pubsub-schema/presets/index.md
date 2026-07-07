---
title: "Presets"
description: "Ready-to-deploy configuration presets for Pub/Sub Schema"
type: "preset-list"
componentSlug: "pubsub-schema"
componentTitle: "Pub/Sub Schema"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-avro-event-contract"
    rank: "01"
    title: "AVRO Event Contract"
    excerpt: "The recommended starting point for schema-validated messaging: an Avro record contract for a business event stream."
  - slug: "02-protobuf-binary-contract"
    rank: "02"
    title: "Protobuf Binary Contract"
    excerpt: "A protobuf-typed schema for high-volume telemetry where compact binary encoding matters."
---

# Pub/Sub Schema Presets

Ready-to-deploy configuration presets for Pub/Sub Schema. Each preset is a complete manifest you can copy, customize, and deploy.
