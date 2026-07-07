# Basic Topic

The minimal event channel: a named topic in the ambient project with
Google-managed encryption and GCP defaults everywhere.

## What this preset creates

A Pub/Sub topic named `order-events`. Publishers send to it; consumers
attach `GcpPubSubSubscription` resources to receive. Message retention is
left to individual subscriptions (the GCP default).

## Remix ideas

- Add `messageRetentionDuration` (e.g. `"604800s"`) to let any
  subscription seek back through the last 7 days.
- Add `schemaSettings` referencing a `GcpPubSubSchema` to reject
  malformed messages at publish time.
- Add user `labels` for cost attribution across topics.
