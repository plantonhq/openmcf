# Catch-All Consumer

This preset creates an unfiltered subscription: the consumer receives
every message published to the topic (Azure's auto-created `$Default`
catch-all rule stays in place). The right starting point for a consumer
that processes the whole stream.

## When to Use

- A consumer that genuinely wants every event (audit trails, analytics
  sinks, full-stream processors)
- The first subscription while filtering needs are still unknown

## Key Configuration Choices

- **No `rules`** -- the auto-created catch-all delivers everything;
  add filtered rules later without touching the publisher
- **`maxDeliveryCount: 10`** -- Azure's queue-side convention; required
  here because subscriptions have no server default
- **`deadLetteringOnMessageExpiration: true`** -- expired messages stay
  inspectable

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-event-topic` | The AzureServiceBusTopic's Planton resource name | Your messaging composition |
| `audit-consumer` | ≤50 chars; typically the consuming application's name | Your naming convention |
