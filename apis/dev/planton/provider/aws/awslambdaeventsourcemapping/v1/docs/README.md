# Lambda Event Source Mappings: The Edge Between Messaging and Compute

## Introduction

An event source mapping is the managed poller AWS runs on your behalf: it reads records from a queue or stream, batches them, and invokes your Lambda function. It is not a property of the function or the source -- it is its own resource with a server-assigned UUID, independent lifecycle, and a rich tuning surface (batching, filters, failure routing, concurrency throttles).

Modeling it as a first-class Planton node makes the resource graph honest: pausing consumption, repointing to a different function, or tuning batch behavior happens on the mapping without touching the function or the queue.

## Source Families

The spec accepts exactly one event source:

| Family | Spec field | Notes |
|---|---|---|
| SQS | `eventSourceArn` (queue ARN) | Partial-batch failures via `ReportBatchItemFailures`; optional `scalingMaxConcurrency`. |
| Kinesis / DynamoDB Streams | `eventSourceArn` (stream ARN) | Requires `startingPosition`; supports shard parallelism and failure destinations. |
| MSK / Kafka | `eventSourceArn` (cluster ARN) + `topics` | Optional consumer group, schema registry, provisioned pollers. |
| Self-managed Kafka | `selfManagedKafka.bootstrapServers` + `topics` | Requires VPC/auth via `sourceAccessConfigurations`. |
| Amazon MQ | `eventSourceArn` (broker ARN) + `mqQueue` | Single queue name. |
| DocumentDB | `eventSourceArn` (change stream ARN) + `documentDb` | Database/collection scoping. |

## Composition

The mapping composes onto its neighbors by reference:

- **Function** -- `functionArn` references an `AwsLambda` node's `function_arn` output.
- **Source** -- `eventSourceArn` references the appropriate kind (`AwsSqsQueue`, `AwsKinesisStream`, `AwsDynamodb`, `AwsMskCluster`, ...).
- **Encryption / failure routing** -- `kmsKeyArn` and `onFailureDestinationArn` reference `AwsKmsKey` and `AwsSqsQueue` (or SNS) nodes.

## Batching Guidance

AWS defaults optimize for latency (small batches, no batching window). For throughput and cost, raise `batchSize` together with `maximumBatchingWindowSeconds` so the poller can actually fill larger batches.
