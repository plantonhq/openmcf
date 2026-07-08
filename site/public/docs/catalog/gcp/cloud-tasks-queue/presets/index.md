---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Tasks Queue"
type: "preset-list"
componentSlug: "cloud-tasks-queue"
componentTitle: "Cloud Tasks Queue"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-basic-queue"
    rank: "01"
    title: "Basic Queue"
    excerpt: "The minimal task queue: a named queue in the ambient project with GCP-managed defaults for dispatch rate and retries."
  - slug: "02-rate-limited-processing"
    rank: "02"
    title: "Rate-Limited Processing"
    excerpt: "A production queue with explicit throughput ceilings and a bounded exponential-backoff retry policy — the shape that protects a downstream service from both traffic spikes and retry storms."
  - slug: "03-secure-cloud-run-target"
    rank: "03"
    title: "Secure Cloud Run Target"
    excerpt: "The recommended modern Cloud Tasks pattern: the queue owns authentication and routing, so producers enqueue bare payloads and every task is dispatched to one Cloud Run service with an automatically..."
---

# Cloud Tasks Queue Presets

Ready-to-deploy configuration presets for Cloud Tasks Queue. Each preset is a complete manifest you can copy, customize, and deploy.
