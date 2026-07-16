---
title: "Cloud Storage Ingestion Topic"
description: "The no-code file-to-stream bridge: Pub/Sub tails a GCS bucket and turns matching objects into messages — no Dataflow job, no custom loader."
type: "preset"
rank: "04"
presetSlug: "04-cloud-storage-ingestion"
componentSlug: "pubsub-topic"
componentTitle: "Pub/Sub Topic"
provider: "gcp"
icon: "package"
order: 4
---

# Cloud Storage Ingestion Topic

The no-code file-to-stream bridge: Pub/Sub tails a GCS bucket and turns
matching objects into messages — no Dataflow job, no custom loader.

## What this preset creates

A topic that ingests newline-delimited JSON objects from a referenced
`GcpGcsBucket` as they land. Each line becomes one message; INFO-level
platform logs surface ingestion progress and errors.

## Prerequisites

- A `GcpGcsBucket` named `data-landing-bucket` (or swap in a literal
  bucket name).
- The Pub/Sub service agent needs `storage.objects.get`/`list` on the
  bucket (`roles/storage.objectViewer` is sufficient).

## Remix ideas

- Swap `textFormat` for `avroFormat` (Avro object payloads) or
  `pubsubAvroFormat` (re-import of Cloud Storage subscription exports).
- Add `minimumObjectCreateTime` to skip historical objects on first
  attach.
- Tighten `matchGlob` per data product (e.g. `orders/**/*.json`).
