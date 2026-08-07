# AWS SNS Subscription

Deliver messages from an SNS topic to an endpoint — SQS, Lambda, HTTP/S, email, SMS, Firehose, or a mobile push application. Each subscription is its own composable node with its own filtering, dead-lettering, and replay lifecycle.

## What Gets Created

- An SNS subscription binding one topic to one endpoint over one protocol.
- Optional message filtering evaluated against message attributes or the message body.
- Optional delivery dead-letter wiring to an SQS queue for undeliverable messages.
- Optional replay of a FIFO topic's archived messages to backfill the new consumer.

## Prerequisites

- An `AwsSnsTopic` (or the ARN of an external/cross-account topic that granted you `sns:Subscribe`).
- The endpoint resource for your protocol — e.g. an `AwsSqsQueue` whose policy allows `sns.amazonaws.com` to `SendMessage`, or an `AwsLambda` with an SNS invoke permission.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSnsSubscription
metadata:
  name: fulfillment-events
spec:
  region: us-west-2
  topicArn:
    valueFrom:
      kind: AwsSnsTopic
      name: order-events
      fieldPath: status.outputs.topic_arn
  protocol: sqs
  endpoint:
    valueFrom:
      kind: AwsSqsQueue
      name: fulfillment-queue
      fieldPath: status.outputs.queue_arn
  rawMessageDelivery: true
```

## Configuration Reference

### Required Fields

| Field | Description |
|---|---|
| `region` | AWS region; must match the topic's region. |
| `topicArn` | The topic to subscribe to (reference or literal ARN). Replaces the subscription when changed. |
| `protocol` | `sqs`, `lambda`, `http`, `https`, `email`, `email-json`, `sms`, `firehose`, or `application`. Replaces when changed. |
| `endpoint` | The delivery target; format depends on the protocol. Replaces when changed. |

### Common Optional Fields

| Field | Description |
|---|---|
| `filterPolicy` / `filterPolicyScope` | Deliver only matching messages; evaluate against attributes (default) or the body. |
| `rawMessageDelivery` | Skip the JSON envelope for SQS/HTTP/S/Firehose targets. |
| `deadLetterConfig.deadLetterTargetArn` | SQS queue that receives messages whose delivery ultimately failed. |
| `replayPolicy` | Backfill this consumer from a FIFO topic's message archive. |
| `subscriptionRoleArn` | Required for `firehose`: the role SNS assumes to write to the stream. |
| `endpointAutoConfirms` / `confirmationTimeoutMinutes` | Tune the HTTP/S confirmation handshake. |

## Stack Outputs

| Output | Description |
|---|---|
| `subscription_arn` | The subscription's ARN. |
| `owner_id` | AWS account that owns the subscription. |
| `pending_confirmation` | True while the endpoint has not confirmed. |
| `confirmation_was_authenticated` | True when the confirmation was signed. |

## Related Components

- [AWS SNS Topic](/docs/catalog/aws/awssnstopic) — the topic this subscription delivers from.
- [AWS SQS Queue](/docs/catalog/aws/awssqsqueue) — the most common delivery target and the DLQ target.
- [AWS Lambda](/docs/catalog/aws/awslambda) — serverless consumers of topic messages.
- [AWS Lambda Event Source Mapping](/docs/catalog/aws/awslambdaeventsourcemapping) — the SQS-to-Lambda edge that pairs with topic→queue fan-out.
