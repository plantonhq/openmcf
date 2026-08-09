---
title: "Error Archive to GCS"
description: "The cheapest compliance archive: every ERROR-and-above entry in the project lands as hourly JSON batches in a Cloud Storage bucket."
type: "preset"
rank: "01"
presetSlug: "01-error-archive-to-gcs"
componentSlug: "logging-sink"
componentTitle: "Logging Sink"
provider: "gcp"
icon: "package"
order: 1
---

# Error Archive to GCS

The cheapest compliance archive: every ERROR-and-above entry in the
project lands as hourly JSON batches in a Cloud Storage bucket.

## What it configures

- A project-scoped sink (no `scope` block — the ambient project).
- A `gcsBucket` destination — the module renders the
  `storage.googleapis.com/...` URI from the bucket name.
- `severity>=ERROR` — errors are evidence; INFO is noise at archive
  prices.

## The deploy's second half

Grant the sink's `writer_identity` output
`roles/storage.objectCreator` on the bucket — through the bucket's
`iamMembers` in the same chart. Without the grant the sink reports
success and exports NOTHING.

## Adjust before deploying

- **gcsBucket** — reference a GcpGcsBucket resource via valueFrom in
  charts so the grant flow stays declarative.
- Consider bucket lifecycle rules (the bucket kind's surface) to age
  archives to Coldline.

## When to choose something else

Need to QUERY the logs? The **Audit Logs to BigQuery** preset. Feeding
an external pipeline? The **Log Stream to Pub/Sub** preset.
