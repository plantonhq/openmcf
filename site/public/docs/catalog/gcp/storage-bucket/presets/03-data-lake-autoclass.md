---
title: "Dual-Region Data Lake with Autoclass"
description: "Analytics storage that manages its own cost: a custom dual-region bucket where Autoclass moves every object to the cheapest justified storage class, hygiene rules reclaim hidden and temporary..."
type: "preset"
rank: "03"
presetSlug: "03-data-lake-autoclass"
componentSlug: "storage-bucket"
componentTitle: "Storage Bucket"
provider: "gcp"
icon: "package"
order: 3
---

# Dual-Region Data Lake with Autoclass

Analytics storage that manages its own cost: a custom dual-region bucket
where Autoclass moves every object to the cheapest justified storage
class, hygiene rules reclaim hidden and temporary storage, and access is
split between a pipeline owner and prefix-scoped readers.

## What this preset creates

A bucket pinned to `US-EAST1` + `US-WEST1` (dual-region durability with
region control), Autoclass enabled with `ARCHIVE` as the terminal class,
a weekly sweep of abandoned multipart uploads, a 30-day TTL on `tmp/`
objects, and two additive grants — the pipeline's `objectAdmin` and a
conditional `objectViewer` limited to the `reports/` prefix.

## Prerequisites

- A `GcpServiceAccount` named `data-pipeline` (the identity that writes
  the lake). Replace the analyst group with your own reader principal.

## Composing analytics

Dataproc clusters take this bucket's `bucket_id` output as
`stagingBucket`/`tempBucket`; BigQuery external tables and load jobs read
`gs://` URIs under it; Pub/Sub Cloud Storage subscriptions can sink
events into it.

## Remix ideas

- Prefer explicit `SetStorageClass` lifecycle transitions instead of
  Autoclass when access patterns are fully predictable (the spec rejects
  configuring both — they fight over storage classes).
- Add `rpo: ASYNC_TURBO` for a 15-minute replication SLA between the two
  regions when the lake backs disaster-recovery commitments.
- Add `enableObjectRetention: true` at creation if regulated datasets
  need per-object WORM retention later.
