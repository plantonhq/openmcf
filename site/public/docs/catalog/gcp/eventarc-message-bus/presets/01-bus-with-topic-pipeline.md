---
title: "Bus with Topic Pipeline"
description: "The smallest useful hub: everything published to the bus lands on one Pub/Sub topic, where downstream consumers subscribe — Eventarc Advanced as a managed firehose."
type: "preset"
rank: "01"
presetSlug: "01-bus-with-topic-pipeline"
componentSlug: "eventarc-message-bus"
componentTitle: "Eventarc Message Bus"
provider: "gcp"
icon: "package"
order: 1
---

# Bus with Topic Pipeline

The smallest useful hub: everything published to the bus lands on one
Pub/Sub topic, where downstream consumers subscribe — Eventarc Advanced
as a managed firehose.

## What it configures

- A message bus with INFO platform logging (the onboarding posture).
- One pipeline delivering to a Pub/Sub topic, one enrollment routing
  everything (`celMatch: "true"`).

## Adjust before deploying

- **destination.topic** — reference your GcpPubSubTopic's `topic_id`
  output.
- **location** — Eventarc Advanced serves a subset of regions; the API
  rejects unsupported ones at create time.
- **celMatch** — tighten from `"true"` as consumers specialize (route
  on `message.type` / `message.source`).

## After deploying

Publish a CloudEvent to the bus (`gcloud eventarc message-buses
publish`); it appears on the topic within seconds. Add API sources to
feed Google-service events in without touching publishers.

## When to choose something else

For CEL-split routing to multiple consumers, start from the **Audit
Fan-Out** preset. For a single point-to-point route, use
GcpEventarcTrigger (Eventarc Standard) instead of a bus.
