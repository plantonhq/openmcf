# Queue Fan-Out

This preset delivers a custom topic's events into a storage queue -- the cheapest at-least-once consumer: workers drain the queue at their own pace, and poison messages stay visible instead of vanishing.

## When to Use

- Worker-pool processing where handlers pull rather than being pushed
- Decoupling bursty publishers from slower consumers (the queue absorbs the burst)

## Key Configuration Choices

- **Storage queue destination** -- pull-based, cheap, and tolerant of consumer downtime; note Azure ignores delivery properties on queue destinations
- **Delivery schema matches the topic** -- Azure rejects a subscription whose delivery schema the source topic cannot map (a CloudEvents-input topic cannot deliver `EventGridSchema`); the preset pins `CloudEventSchemaV1_0` to match the catalog's topic examples -- drop it only for an `EventGridSchema`-input topic
- **Dead-letter configured up front** -- events that exhaust retries land in a blob container instead of being dropped; create the container before deploying
- **No filters** -- every event on the topic reaches the queue; add `subjectFilter` or `advancedFilter` when consumers need less noise

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-eventgrid-topic>` | The Planton name of your `AzureEventgridTopic` resource | Planton console (or replace `valueFrom` with `value:` and any scope's ARM ID) |
| `<your-storage-account>` | The Planton name of your `AzureStorageAccount` resource | Planton console |
| `order-events` | The queue's name (an `AzureStorageQueue` manages it) | Your naming convention |
| `eventgrid-dead-letters` | The dead-letter blob container (an `AzureStorageContainer` manages it) | Your naming convention |
