# Ordered Topic with Duplicate Detection

This preset creates a topic for exactly-once-flavored, ordered
publish-subscribe: publish order is preserved for session-aware
subscriptions, and retried publishes are dropped within the detection
window.

## When to Use

- Event streams where consumers replaying out of order would corrupt
  state (ledgers, inventory movements)
- Publishers that retry on ambiguous failures

## Key Configuration Choices

- **`supportOrdering: true`** -- pair every subscription with
  `requiresSession: true` so each session's events process in order
- **`requiresDuplicateDetection: true`** (ForceNew) + **`PT10M`
  window** -- MessageId-based dedup sized to the publisher's retry
  horizon

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-bus` | The AzureServiceBusNamespace's Planton resource name | Your messaging composition |
| `ledger-entries` | Unique within the namespace | Your naming convention |
