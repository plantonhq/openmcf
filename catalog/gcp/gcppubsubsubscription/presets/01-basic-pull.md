# Basic Pull Subscription

The default consumer shape: a pull subscription with GCP defaults
everywhere, composed onto a referenced topic.

## What this preset creates

A pull subscription named `order-events-worker` attached to the
`order-events` topic by reference. Consumers pull messages via the
client libraries and acknowledge within the default 10-second deadline;
unacknowledged messages redeliver automatically.

## Remix ideas

- Raise `ackDeadlineSeconds` (10-600) for slow processing.
- Add `enableExactlyOnceDelivery: true` for consumers that cannot
  tolerate redelivery within the deadline window.
- Add a `filter` (immutable) to receive only matching messages —
  non-matching ones are auto-acked and never delivered.
