# FIFO Topic with Message Archive

A strictly-ordered FIFO topic that retains 30 days of messages for consumer backfill.

## What this preset gives you

- Exactly-once, ordered delivery to SQS FIFO subscribers.
- A 30-day message archive — new consumers replay history via their subscription's `replayPolicy` instead of starting blind.
- SHA-256 message signing (`signatureVersion: 2`).

## Before you deploy

- Archiving is FIFO-only; the archive accrues storage cost for the retention window.
- Consumers subscribe through `AwsSnsSubscription` resources referencing this topic's `topic_arn` output; FIFO topics deliver only to SQS endpoints.

## Remix ideas

- Add `kmsKeyId` referencing an `AwsKmsKey` for encryption at rest.
- Raise `MessageRetentionPeriod` (up to 365 days) for longer backfill windows.
- Pair with an `AwsSnsSubscription` carrying `replayPolicy` to bootstrap a new consumer from history.
