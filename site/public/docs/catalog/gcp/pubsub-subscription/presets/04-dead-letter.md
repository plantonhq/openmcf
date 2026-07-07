---
title: "Reliable Processing with Dead-Lettering"
description: "The production-money shape: exactly-once delivery, a deep replay window, backoff retries, and a dead-letter queue for poison messages."
type: "preset"
rank: "04"
presetSlug: "04-dead-letter"
componentSlug: "pubsub-subscription"
componentTitle: "Pub/Sub Subscription"
provider: "gcp"
icon: "package"
order: 4
---

# Reliable Processing with Dead-Lettering

The production-money shape: exactly-once delivery, a deep replay window,
backoff retries, and a dead-letter queue for poison messages.

## What this preset creates

A pull subscription hardened for payment-grade processing:

- **Exactly-once delivery** — no redelivery inside the ack deadline.
- **14-day retention with acked messages kept** — seek back and replay
  after a consumer bug.
- **Dead-lettering** — after 10 failed deliveries a message diverts to
  the referenced `payments-dlq` topic instead of blocking the stream.
- **Backoff retries** (30s–300s) — a failing consumer is not hot-looped.
- **Never expires** — quiet periods don't delete the subscription.

## Prerequisites

- `GcpPubSubTopic` resources named `payment-events` and `payments-dlq`.
- The Pub/Sub service agent needs Subscriber on this subscription and
  Publisher on the DLQ topic (grant with `GcpProjectIamMember` or
  topic-level grants).
- Attach a subscription to `payments-dlq` — dead-lettered messages are
  lost if nothing consumes them.

## Remix ideas

- Add a `filter` to shard payment types across parallel subscriptions.
- Lower `maxDeliveryAttempts` (minimum 5) for faster poison detection.
