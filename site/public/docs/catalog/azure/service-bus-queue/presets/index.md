---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Bus Queue"
type: "preset-list"
componentSlug: "service-bus-queue"
componentTitle: "Service Bus Queue"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-work-queue"
    rank: "01"
    title: "Work Queue"
    excerpt: "This preset creates a competing-consumers work queue: multiple workers drain one queue, failures quarantine to the dead-letter sub-queue, and expired messages stay inspectable. The right starting..."
  - slug: "02-session-fifo-queue"
    rank: "02"
    title: "Session FIFO Queue"
    excerpt: "This preset creates a session-aware queue with duplicate detection: strict per-session ordering with exclusive consumption, and idempotent producers that can retry sends safely. The right shape for..."
---

# Service Bus Queue Presets

Ready-to-deploy configuration presets for Service Bus Queue. Each preset is a complete manifest you can copy, customize, and deploy.
