---
title: "Work Queue"
description: "This preset creates a competing-consumers work queue: multiple workers drain one queue, failures quarantine to the dead-letter sub-queue, and expired messages stay inspectable. The right starting..."
type: "preset"
rank: "01"
presetSlug: "01-work-queue"
componentSlug: "service-bus-queue"
componentTitle: "Service Bus Queue"
provider: "azure"
icon: "package"
order: 1
---

# Work Queue

This preset creates a competing-consumers work queue: multiple workers
drain one queue, failures quarantine to the dead-letter sub-queue, and
expired messages stay inspectable. The right starting point for job and
command processing.

## When to Use

- Background job processing (orders, exports, notifications)
- Command distribution to a worker fleet

## Key Configuration Choices

- **`maxDeliveryCount: 5`** -- poison messages quarantine after 5
  attempts; raise it if consumers see transient failures worth riding out
- **`deadLetteringOnMessageExpiration: true`** -- expired messages move
  to the dead-letter sub-queue for inspection and replay
- **`defaultMessageTtl: P14D`** -- bounded retention; unprocessed work
  older than two weeks dead-letters

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-bus` | The AzureServiceBusNamespace's Planton resource name | Your messaging composition |
| `orders` | Unique within the namespace | Your naming convention |

## Downstream Wiring

Mint a send-only credential for the producer service:

```yaml
# On an AzureServiceBusAuthorizationRule
ruleName: producer-send
queueId:
  valueFrom:
    kind: AzureServiceBusQueue
    name: my-work-queue
    fieldPath: status.outputs.queue_id
send: true
```
