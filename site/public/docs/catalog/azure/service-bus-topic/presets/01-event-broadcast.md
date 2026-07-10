---
title: "Event Broadcast Topic"
description: "This preset creates a plain broadcast topic: publishers send once, and every subscription under it receives an independent copy. The right starting point for domain events fanning out to multiple..."
type: "preset"
rank: "01"
presetSlug: "01-event-broadcast"
componentSlug: "service-bus-topic"
componentTitle: "Service Bus Topic"
provider: "azure"
icon: "package"
order: 1
---

# Event Broadcast Topic

This preset creates a plain broadcast topic: publishers send once, and
every subscription under it receives an independent copy. The right
starting point for domain events fanning out to multiple consumers.

## When to Use

- Domain events (order-created, user-registered) consumed by several
  independent applications
- Decoupling producers from an evolving set of consumers -- adding a
  subscription never touches the publisher

## Key Configuration Choices

- **`defaultMessageTtl: P14D`** -- bounded retention for undelivered
  events; subscriptions may shorten it further
- **Consumer dials live on the subscriptions** -- locks, delivery
  counts, sessions, and filters are per-consumer choices, not topic ones

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-bus` | The AzureServiceBusNamespace's Planton resource name | Your messaging composition |
| `order-events` | Unique within the namespace | Your naming convention |

## Downstream Wiring

Each consumer materializes its own view:

```yaml
# On an AzureServiceBusSubscription
topicId:
  valueFrom:
    kind: AzureServiceBusTopic
    name: my-event-topic
    fieldPath: status.outputs.topic_id
subscriptionName: billing-consumer
maxDeliveryCount: 10
```
