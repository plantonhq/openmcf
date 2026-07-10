---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Bus Topic"
type: "preset-list"
componentSlug: "service-bus-topic"
componentTitle: "Service Bus Topic"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-event-broadcast"
    rank: "01"
    title: "Event Broadcast Topic"
    excerpt: "This preset creates a plain broadcast topic: publishers send once, and every subscription under it receives an independent copy. The right starting point for domain events fanning out to multiple..."
  - slug: "02-ordered-dedup-topic"
    rank: "02"
    title: "Ordered Topic with Duplicate Detection"
    excerpt: "This preset creates a topic for exactly-once-flavored, ordered publish-subscribe: publish order is preserved for session-aware subscriptions, and retried publishes are dropped within the detection..."
---

# Service Bus Topic Presets

Ready-to-deploy configuration presets for Service Bus Topic. Each preset is a complete manifest you can copy, customize, and deploy.
