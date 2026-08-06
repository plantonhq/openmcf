# SNS Topics: The Pub/Sub Hub

## Introduction

An SNS topic is the hub of AWS pub/sub: publishers send to one ARN, and every attached subscription receives (a filtered view of) the stream. The topic owns identity, ordering semantics, encryption, and policy; each consumer relationship is a first-class `AwsSnsSubscription` resource owned by the consumer. This split mirrors how organizations actually operate — one team owns the topic, many teams attach and manage their own subscriptions without touching it.

## Standard vs FIFO

Standard topics maximize throughput with at-least-once, best-effort-ordered delivery across all nine protocols. FIFO topics guarantee strict per-message-group ordering and exactly-once delivery, but only to SQS endpoints; names carry a `.fifo` suffix (appended automatically by the modules); and high-throughput mode (`fifoThroughputScope: MessageGroup`) pairs with SQS FIFO queues configured `perMessageGroupId` for an end-to-end high-throughput ordered pipeline.

## Message Archiving and Replay

FIFO topics can retain published messages (`archivePolicy`) for up to the configured retention window. A subscription added later can replay the archive from a starting point (`replayPolicy` on `AwsSnsSubscription`) before receiving live traffic — the backfill mechanism that makes adding consumers to an established event stream safe. The `beginning_archive_time` output reports the earliest replayable point.

## Policies

- **Access policy** (`policy`): who may Publish/Subscribe. AWS keeps a policy on every topic — deleting one reverts to the owner-only default rather than leaving the topic open or policy-less.
- **Data protection policy** (`dataProtectionPolicy`): standard-topic-only PII/PHI inspection with audit, mask (de-identify), or deny operations. Materialized as its own AWS resource keyed by the topic ARN — a single-per-topic folded satellite.
- **Delivery policy** (`deliveryPolicy`): topic-level HTTP/S retry defaults; subscriptions can override per-edge.

## Delivery-Status Logging

Each delivery protocol (application, firehose, http, lambda, sqs) has an independent opt-in logging block: an IAM role for successes (with a 0-100 sample rate) and a role for failures (always logged when set). SNS assumes these roles to write into CloudWatch Logs. The spec groups the provider's fifteen flat attributes into five per-protocol blocks so a manifest states intent ("log SQS delivery failures") rather than attribute soup.

## Struct-Typed Policy Documents

The access, data-protection, archive, and (on subscriptions) filter/replay policies use `google.protobuf.Struct` for native YAML authoring — users write the JSON structure directly without escaping. The IaC layers serialize to JSON strings for the AWS API. `deliveryPolicy` stays a raw JSON string: it is a rarely-used pass-through document with no composition value in structuring.

## 90/10 Coverage Notes

Deliberately not modeled, with reasons:

- **`aws_sns_platform_application`** — mobile push credentials (APNS/GCM) are a separate product surface with sensitive platform credentials; a candidate kind on demand.
- **`aws_sns_sms_preferences`** — an account-level singleton, not a composable resource.
- **`name_prefix`** — resource names derive from `metadata.name` across the catalog; generated suffixes would break reference-by-name composition.
