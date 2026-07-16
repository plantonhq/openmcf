# SNS Subscriptions: The Delivery Edge of Pub/Sub

## Introduction

An SNS topic by itself delivers nothing — every consumer relationship is a subscription, and each subscription carries its own delivery contract: the protocol, the endpoint, what subset of messages to receive, what happens when delivery fails, and whether to backfill from the topic's archive. Modeling the subscription as its own resource keeps that contract with the consumer that owns it. A topic may have hundreds of subscriptions owned by different teams; the topic owner controls identity, encryption, and policy, while each subscriber controls its own edge.

## Protocol Landscape

- **sqs** — the workhorse. Durable fan-out into queues that Lambda event source mappings or services consume. Confirms automatically. Pair with `rawMessageDelivery: true` unless consumers want the SNS JSON envelope.
- **lambda** — direct push invocation. Confirms automatically; the function needs an SNS invoke permission.
- **http / https** — webhook delivery with an asynchronous confirmation handshake. The endpoint must acknowledge a SubscriptionConfirmation callback; `endpointAutoConfirms` tells deployment to wait for it.
- **email / email-json** — human-in-the-loop confirmation; deployments cannot wait for the click, so the subscription stays `pending_confirmation` until the recipient acts.
- **firehose** — streams messages into Kinesis Data Firehose for lake/warehouse delivery; requires `subscriptionRoleArn`.
- **sms / application** — phone and mobile-push delivery.

## The Silent-Drop Coupling

SNS-to-SQS delivery is authorized by the QUEUE's resource policy, not the subscription. A subscription to a queue whose policy does not grant `sqs:SendMessage` to `sns.amazonaws.com` (scoped by `aws:SourceArn` to the topic) is created successfully and then drops every message. Always ship the queue policy with the queue (`AwsSqsQueue.policy`); the subscription docs cannot fix this for you at delivery time.

## FIFO Coupling

FIFO topics deliver only to SQS endpoints, and end-to-end exactly-once requires a FIFO queue subscriber. A FIFO topic with `archivePolicy` enables `replayPolicy` on new subscriptions — the backfill mechanism for consumers added after messages were published.

## 90/10 Coverage Notes

- The full provider surface of `aws_sns_topic_subscription` is modeled: filtering (attributes + body scope), raw delivery, redrive, per-subscription delivery policy, replay, Firehose role, and the HTTP/S confirmation knobs.
- Endpoint values are intentionally a `StringValueOrRef` without a default kind — the referenced resource type varies by protocol (queue, function, stream, or a plain URL/email/phone literal).
