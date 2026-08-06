# Topic with Message Retention

The replay-ready event stream: topic-level retention so consumers can
seek backwards after a bug or a new subscriber needs history.

## What this preset creates

A topic retaining every published message for 7 days regardless of
acknowledgement state. Any attached subscription can seek to a timestamp
within that window — deploy a fixed consumer, seek back, and reprocess.

## Why topic-level (not subscription-level)

Subscription retention only covers that one subscription's backlog.
Topic retention covers ALL subscriptions — including ones created after
the messages were published. It is the difference between "my consumer
can retry" and "any consumer can replay."

## Remix ideas

- Range is `600s` (10 minutes) to `2678400s` (31 days); billing scales
  with retained bytes.
- Pair with a subscription's `retainAckedMessages` +
  `messageRetentionDuration` for per-consumer replay windows instead.
