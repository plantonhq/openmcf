# SQS Fan-out Consumer

Subscribe an SQS queue to a topic with attribute filtering and raw delivery — the standard pub/sub fan-out edge.

## What this preset gives you

- A topic→queue subscription referencing both resources by output, no literal ARNs.
- Attribute-based filtering so the queue receives only the event types its consumer processes.
- Raw message delivery, sparing consumers the SNS JSON envelope.

## Before you deploy

- The queue's `policy` must grant `sqs:SendMessage` to `sns.amazonaws.com` with an `aws:SourceArn` condition matching the topic — without it, deliveries are silently dropped.
- Topic and queue must be in the same region as the subscription.

## Remix ideas

- Add `deadLetterConfig` pointing at a second queue to capture undeliverable messages.
- Switch `filterPolicyScope` to `MessageBody` to filter on payload fields instead of attributes.
- Duplicate this preset per consumer queue to build a full fan-out topology from one topic.
