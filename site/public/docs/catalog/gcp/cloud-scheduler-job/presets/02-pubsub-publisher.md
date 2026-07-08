---
title: "Pub/Sub Publisher"
description: "Cron-driven event publishing: the job publishes a message to a Pub/Sub topic on schedule, and downstream subscriptions fan the event out to whatever needs to run."
type: "preset"
rank: "02"
presetSlug: "02-pubsub-publisher"
componentSlug: "cloud-scheduler-job"
componentTitle: "Cloud Scheduler Job"
provider: "gcp"
icon: "package"
order: 2
---

# Pub/Sub Publisher

Cron-driven event publishing: the job publishes a message to a Pub/Sub
topic on schedule, and downstream subscriptions fan the event out to
whatever needs to run.

## What this preset creates

A Cloud Scheduler job named `nightly-export-trigger` that publishes a
base64-encoded JSON payload (with routing attributes) to the referenced
`GcpPubSubTopic` at 02:00 UTC nightly, retrying failed publishes up to
3 times.

## When to use

- Triggering daily/hourly batch pipelines through an event, not a direct
  call — consumers scale and retry independently of the scheduler
- Scheduling periodic data exports or ETL kickoffs
- Publishing heartbeat events that multiple systems observe

## Key configuration choices

- `topicName` is a `GcpPubSubTopic` reference resolving to the topic's
  fully qualified path — the composition survives project moves and
  renames of the underlying topic resource.
- `data` must be base64-encoded; `attributes` ride alongside as plain
  metadata (useful for subscription filters).
- A Pub/Sub target ignores `attemptDeadline` (publish is fast); retry
  config still applies to failed publishes.

## Placeholders to replace

- The `export-events` topic reference — the name of your
  `GcpPubSubTopic` resource.
- `data` — base64 of your real payload.

## Related presets

- `03-secure-cloud-run-trigger` — call one service directly instead of
  fanning out through an event.
