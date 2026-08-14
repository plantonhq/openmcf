---
title: "Realtime Serving Features"
description: "This preset is an online-only feature group: a handful of customer features served at low latency, with stale records expiring 30 days after their event time."
type: "preset"
rank: "01"
presetSlug: "01-realtime-serving-features"
componentSlug: "sagemaker-feature-group"
componentTitle: "SageMaker Feature Group"
provider: "aws"
icon: "package"
order: 1
---

# Realtime Serving Features

This preset is an online-only feature group: a handful of customer
features served at low latency, with stale records expiring 30 days
after their event time.

## When to Use

- Real-time inference that looks features up at request time
- Serving-only features where no training dataset lives in this group

## What You Get

- An online store in Standard storage with a 30-day record TTL
  (`ExpiresAt = EventTime + ttl`)
- A five-feature schema keyed on `customer_id` and stamped by
  `event_time`

## Customize

- Tune the `ttl` freely — it is the only online-store setting that
  updates in place; everything else replaces the group
- Add an `offlineStore` at creation if training will ever need this
  data — stores are create-time structure
- Switch `storageType` to `InMemory` (create-time) for the lowest
  latency, and required if you add collection-typed features
