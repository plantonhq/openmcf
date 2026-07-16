---
title: "Per-Application Consumer Group"
description: "This preset creates one consumer group for one consuming application -- the discipline that lets many applications read the same stream independently without their offsets colliding."
type: "preset"
rank: "01"
presetSlug: "01-per-application-group"
componentSlug: "event-hub-consumer-group"
componentTitle: "Event Hub Consumer Group"
provider: "azure"
icon: "package"
order: 1
---

# Per-Application Consumer Group

This preset creates one consumer group for one consuming application --
the discipline that lets many applications read the same stream
independently without their offsets colliding.

## When to Use

- Every real consumer application gets its OWN group; give a new
  processor its own group instead of sharing an existing one
- Never point production applications at `$Default` (Azure's
  service-created catch-all, reserved by the spec) -- shared cursors
  mean applications steal each other's events

## Key Configuration Choices

- **`consumerGroupName`** is the group's identity and is ForceNew --
  renaming replaces the group and resets stored offsets
- **`userMetadata`** -- free-form ownership breadcrumbs operators see
  in the portal; invaluable when auditing who reads a stream

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-telemetry-stream` | The AzureEventHub's Planton resource name | Your streaming composition |
| `analytics` | The consuming application's group name | Your consumer taxonomy |

## Downstream Wiring

Consumers pass the group name to their SDK client:

```yaml
consumerGroup:
  valueFrom:
    kind: AzureEventHubConsumerGroup
    name: my-analytics-consumer
    fieldPath: status.outputs.consumer_group_name
```
