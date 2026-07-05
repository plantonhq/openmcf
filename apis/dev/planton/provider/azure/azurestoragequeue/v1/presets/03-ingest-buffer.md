# Ingest Buffer Queue

This preset creates a queue absorbing a bursty external producer --
webhooks, device telemetry, upload notifications -- so downstream
processing drains at its own pace instead of scaling to the burst.

## When to Use

- Webhook receivers that must ACK fast and process later
- Telemetry or event ingestion where bursts dwarf steady-state
- Buffering in front of a rate-limited downstream (an external API, a
  database)

## When NOT to Use (reach for the richer sibling)

Storage queues deliver at-least-once with no ordering, no sessions, no
topics, and 7-day maximum message TTL. If the integration needs
pub/sub fan-out, FIFO sessions, or first-class dead-lettering, use
Service Bus -- Storage queues win when the requirement is "cheap,
huge, and simple."

## Key Configuration Choices

- **One queue per producer system** -- isolates backpressure and lets
  grants stay least-privilege (the producer gets Message Sender on just
  its queue)
- **`metadata.producer`** -- records the source system for operators

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `<queue-name>` | 3-63 lowercase letters/digits/hyphens | Your naming convention |
| `<producer-system>` | What enqueues into this buffer | Your integration inventory |
