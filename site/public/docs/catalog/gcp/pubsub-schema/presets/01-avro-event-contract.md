---
title: "AVRO Event Contract"
description: "The recommended starting point for schema-validated messaging: an Avro record contract for a business event stream."
type: "preset"
rank: "01"
presetSlug: "01-avro-event-contract"
componentSlug: "pubsub-schema"
componentTitle: "Pub/Sub Schema"
provider: "gcp"
icon: "package"
order: 1
---

# AVRO Event Contract

The recommended starting point for schema-validated messaging: an Avro
record contract for a business event stream.

## What this preset creates

A Pub/Sub schema named `order-events` whose definition is an Avro record
with typed fields (including a logical timestamp). Any topic that attaches
this schema (via its `schemaSettings.schema` reference) will reject
published messages that do not conform — moving contract violations from
consumers to publishers, where the producing team can fix them.

## Why AVRO

- Human-readable JSON definitions that review well in pull requests.
- Native fit with Pub/Sub's BigQuery delivery (`useTopicSchema`) and
  Cloud Storage Avro export (`avroConfig.useTopicSchema`) — the schema
  drives table columns and file layout automatically.
- Logical types (timestamps, decimals) carry semantics that plain JSON
  cannot.

## Remix ideas

- Add fields over time by committing new definitions — each update is a
  schema revision (up to 20 retained), not a resource replacement. Keep
  revisions backward-compatible with the encoding your topics use.
- Point `projectId` at a `GcpProject` reference to pin the schema to an
  explicitly managed project instead of the provider default.
