---
title: "Queue Sender Credential"
description: "This preset mints a send-only SAS credential scoped to one queue -- the least-privilege shape for a producer service: it can send to its queue and do nothing else in the namespace."
type: "preset"
rank: "01"
presetSlug: "01-queue-sender"
componentSlug: "service-bus-authorization-rule"
componentTitle: "Service Bus Authorization Rule"
provider: "azure"
icon: "package"
order: 1
---

# Queue Sender Credential

This preset mints a send-only SAS credential scoped to one queue -- the
least-privilege shape for a producer service: it can send to its queue
and do nothing else in the namespace.

## When to Use

- A producer service that enqueues work (API frontends, event
  ingesters)
- Any integration that needs a connection string with the narrowest
  possible blast radius

## Key Configuration Choices

- **`queueId` scope** -- rights cover exactly this queue; a compromised
  credential cannot read messages or touch other entities
- **`send: true` only** -- add `listen` for request-reply patterns
  where the producer also drains a response queue (mint a second rule
  instead when the queues differ)

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `orders-api` | The producing application, for the credential's name | Your service inventory |
| `my-work-queue` | The AzureServiceBusQueue's Planton resource name | Your messaging composition |

## Downstream Wiring

The application consumes the connection string output:

```yaml
primary_connection_string  # sensitive output -- wire through managed secrets
```
