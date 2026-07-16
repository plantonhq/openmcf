---
title: "Telemetry Stream"
description: "This preset creates a general-purpose event stream: 8 partitions and a 3-day replay window -- the shape most telemetry, logging, and change-data-capture pipelines start from."
type: "preset"
rank: "01"
presetSlug: "01-telemetry-stream"
componentSlug: "event-hub"
componentTitle: "Event Hub"
provider: "azure"
icon: "package"
order: 1
---

# Telemetry Stream

This preset creates a general-purpose event stream: 8 partitions and a
3-day replay window -- the shape most telemetry, logging, and
change-data-capture pipelines start from.

## When to Use

- Application/service telemetry fan-in
- The Kafka-topic equivalent for producers migrating from Kafka (the
  hub's name IS the topic name on the namespace's Kafka endpoint)

## Key Configuration Choices

- **`partitionCount: 8`** -- the unit of parallelism; downstream
  consumers cannot out-scale it. On shared namespaces the count is
  fixed for the hub's life, so size for peak
- **`messageRetention: 3`** days -- the simple retention model; switch
  to `retentionDescription` for hour-granular windows or Kafka-style
  compaction

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-telemetry-hubs` | The AzureEventHubNamespace's Planton resource name | Your streaming composition |
| `telemetry` | The hub/topic name within the namespace | Your stream taxonomy |

## Downstream Wiring

Consumer groups and hub-scoped credentials reference the hub:

```yaml
# On an AzureEventHubConsumerGroup
eventHubId:
  valueFrom:
    kind: AzureEventHub
    name: my-telemetry-stream
    fieldPath: status.outputs.event_hub_id
```
