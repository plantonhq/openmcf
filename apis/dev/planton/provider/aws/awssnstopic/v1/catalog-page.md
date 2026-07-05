# AWS SNS Topic

Deploys an AWS SNS topic — Standard or FIFO — with optional KMS encryption, IAM access and data-protection policies, message archiving, per-protocol delivery-status logging, and X-Ray tracing. The component handles FIFO naming conventions automatically. Consumers attach through first-class [AwsSnsSubscription](/docs/catalog/aws/awssnssubscription) resources referencing this topic's `topic_arn` output.

## What Gets Created

- **SNS Topic** — one `aws_sns_topic`, Standard or FIFO, named from `metadata.name` (with `.fifo` appended automatically for FIFO topics)
- **Data protection policy** (optional) — one `aws_sns_topic_data_protection_policy` when `dataProtectionPolicy` is set, detecting and auditing/masking/blocking sensitive data in transit

## Prerequisites

- An AWS account connection with permission to manage SNS topics.
- An `AwsKmsKey` (or existing KMS key ARN) if encryption at rest is required.
- IAM roles assumable by `sns.amazonaws.com` if delivery-status logging is configured.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSnsTopic
metadata:
  name: order-events
spec:
  region: us-west-2
  signatureVersion: 2
```

Subscribe consumers with separate `AwsSnsSubscription` resources:

```yaml
apiVersion: aws.planton.dev/v1
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
```

## Configuration Reference

### Required Fields

| Field | Description |
|---|---|
| `region` | AWS region for the topic. |

### Common Optional Fields

| Field | Type | Default | Description |
|---|---|---|---|
| `fifoTopic` | `bool` | `false` | Create a FIFO topic (strict ordering, exactly-once to SQS FIFO). Immutable. |
| `contentBasedDeduplication` | `bool` | `false` | FIFO-only: hash the body as the deduplication ID. |
| `fifoThroughputScope` | `string` | AWS default | FIFO-only: `Topic` or `MessageGroup` (high-throughput mode). |
| `archivePolicy` | `Struct` | — | FIFO-only: retain messages for replay, e.g. `{"MessageRetentionPeriod": 30}`. |
| `displayName` | `string` | — | Human-readable label (SMS "from" name). |
| `kmsKeyId` | `StringValueOrRef` | — | Customer-managed KMS key for encryption at rest. References `AwsKmsKey`. |
| `policy` | `Struct` | AWS owner-only default | Resource-based access policy (who may Publish/Subscribe). |
| `dataProtectionPolicy` | `Struct` | — | Standard-only: PII/PHI detection with audit/mask/deny operations. |
| `deliveryPolicy` | `string` | AWS default | Topic-level HTTP/S retry policy JSON. |
| `deliveryFeedback` | `object` | — | Per-protocol delivery-status logging (application/firehose/http/lambda/sqs), each with success/failure IAM roles and a success sample rate. |
| `tracingConfig` | `string` | `PassThrough` | X-Ray tracing: `Active` or `PassThrough`. |
| `signatureVersion` | `int32` | `1` | Message signing: `1` (SHA1) or `2` (SHA256, recommended). |

## Stack Outputs

| Output | Type | Description |
|---|---|---|
| `topic_arn` | `string` | Topic ARN — the reference target for subscriptions, EventBridge targets, and alarm actions |
| `topic_name` | `string` | Topic name (includes `.fifo` suffix for FIFO topics) |
| `owner` | `string` | AWS account ID owning the topic |
| `beginning_archive_time` | `string` | Start of the replayable archive window (FIFO archive only) |

## Related Components

- [AwsSnsSubscription](/docs/catalog/aws/awssnssubscription) — the delivery edge attaching consumers to this topic
- [AwsSqsQueue](/docs/catalog/aws/awssqsqueue) — the most common fan-out target
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — customer-managed encryption at rest
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — delivery-status logging roles
- [AwsCloudwatchAlarm](/docs/catalog/aws/awscloudwatchalarm) — publishes alarm notifications to topics
