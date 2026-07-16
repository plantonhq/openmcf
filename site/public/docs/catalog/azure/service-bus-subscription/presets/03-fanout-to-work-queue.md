---
title: "Fan-Out to Work Queue"
description: "This preset creates a routing subscription: it filters the topic's stream and auto-forwards every match into a work queue that a processing fleet drains. Filtering and processing decouple -- the..."
type: "preset"
rank: "03"
presetSlug: "03-fanout-to-work-queue"
componentSlug: "service-bus-subscription"
componentTitle: "Service Bus Subscription"
provider: "azure"
icon: "package"
order: 3
---

# Fan-Out to Work Queue

This preset creates a routing subscription: it filters the topic's
stream and auto-forwards every match into a work queue that a
processing fleet drains. Filtering and processing decouple -- the
routing layer evolves without touching workers.

## When to Use

- A topic fans out by category, and each category feeds a
  competing-consumers worker pool
- Routing topologies where filters change more often than processors

## Key Configuration Choices

- **`forwardTo`** -- references the target queue's `queue_name` output
  (auto-forward addresses entities by NAME within the namespace); the
  queue must exist before the subscription
- **A correlation filter rule** -- cheap equality matching; remove the
  auto-created `$Default` catch-all once (out-of-band) so ONLY matches
  forward into the work queue

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-event-topic` | The AzureServiceBusTopic's Planton resource name | Your messaging composition |
| `priority-route` | ≤50 chars; typically the category being routed | Your naming convention |
| `my-work-queue` | The AzureServiceBusQueue workers drain | Your messaging composition |
| `order-created` | The message Label (Subject) this route admits | Your message contract |
