# HTTPS Webhook with Dead-Letter Queue

Push messages to a self-confirming HTTPS endpoint and capture undeliverable messages in an SQS DLQ.

## What this preset gives you

- Webhook delivery with the confirmation handshake handled at deploy time (`endpointAutoConfirms`).
- A five-minute confirmation window before the deployment fails fast.
- A dead-letter queue so retry-exhausted messages are kept, not lost.

## Before you deploy

- The endpoint must respond to SNS's `SubscriptionConfirmation` callback automatically; if a human confirms instead, remove `endpointAutoConfirms` and expect `pending_confirmation: true` until they act.
- The DLQ's queue policy must allow `sns.amazonaws.com` to `SendMessage`.

## Remix ideas

- Add a `deliveryPolicy` JSON override to tune HTTP retry backoff for a slow endpoint.
- Add a `filterPolicy` so the webhook receives only the relevant subset of topic traffic.
- Switch to `protocol: firehose` with a `subscriptionRoleArn` to archive topic traffic to S3 via Kinesis Data Firehose.
