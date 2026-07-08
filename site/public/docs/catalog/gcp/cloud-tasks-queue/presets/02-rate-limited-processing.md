---
title: "Rate-Limited Processing"
description: "A production queue with explicit throughput ceilings and a bounded exponential-backoff retry policy — the shape that protects a downstream service from both traffic spikes and retry storms."
type: "preset"
rank: "02"
presetSlug: "02-rate-limited-processing"
componentSlug: "cloud-tasks-queue"
componentTitle: "Cloud Tasks Queue"
provider: "gcp"
icon: "package"
order: 2
---

# Rate-Limited Processing

A production queue with explicit throughput ceilings and a bounded
exponential-backoff retry policy — the shape that protects a downstream
service from both traffic spikes and retry storms.

## What this preset creates

A Cloud Tasks queue named `order-processing` dispatching at most 500
tasks/second with 100 concurrent tasks in flight, retrying failures up to
5 times with 1s–3600s exponential backoff, abandoning tasks after 24
hours, and logging 10% of dispatch operations.

## When to use

- Production background processing with a known downstream capacity
- Integrations against rate-limited third-party APIs
- Order/payment pipelines where retry behavior must be predictable

## Key configuration choices

- `maxDispatchesPerSecond: 500` — size to the downstream service's real
  capacity; GCP computes the effective burst size from this rate.
- `maxRetryDuration: "86400s"` — a task older than a day is dropped even
  if attempts remain, keeping the queue from replaying stale work.
- `samplingRatio: 0.1` — raise to `1.0` while debugging dispatch issues.

## Related presets

- `01-basic-queue` — when GCP-managed defaults are enough.
- `03-secure-cloud-run-target` — adds queue-level OIDC auth and routing.
