---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Queue"
type: "preset-list"
componentSlug: "storage-queue"
componentTitle: "Storage Queue"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-work-queue"
    rank: "01"
    title: "Work Queue"
    excerpt: "This preset creates a plain work queue -- the decoupling primitive: producers enqueue up to 64 KB messages, a worker pool polls and deletes them at its own pace."
  - slug: "02-poison-companion"
    rank: "02"
    title: "Poison-Queue Companion"
    excerpt: "This preset declares the dead-letter companion of a work queue. Azure Functions moves messages that exhaust their retries to `{queue}-poison` by NAMING CONVENTION -- if nothing declares that queue,..."
  - slug: "03-ingest-buffer"
    rank: "03"
    title: "Ingest Buffer Queue"
    excerpt: "This preset creates a queue absorbing a bursty external producer -- webhooks, device telemetry, upload notifications -- so downstream processing drains at its own pace instead of scaling to the burst."
---

# Storage Queue Presets

Ready-to-deploy configuration presets for Storage Queue. Each preset is a complete manifest you can copy, customize, and deploy.
