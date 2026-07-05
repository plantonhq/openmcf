---
title: "Kinesis Stream Consumer"
description: "A Kinesis-to-Lambda mapping that processes the backlog from the oldest available record with partial-batch failure reporting."
type: "preset"
rank: "02"
presetSlug: "02-kinesis-consumer"
componentSlug: "lambda-event-source-mapping"
componentTitle: "Lambda Event Source Mapping"
provider: "aws"
icon: "package"
order: 2
---

# Kinesis Stream Consumer

A Kinesis-to-Lambda mapping that processes the backlog from the oldest available record with partial-batch failure reporting.

## What this preset gives you

- A stream consumer starting at `TRIM_HORIZON` -- processes existing records, not just new ones.
- `ReportBatchItemFailures` for independent per-record retry semantics.
- `parallelizationFactor: 2` to double per-shard throughput while preserving partition-key ordering.
- Composed references to an `AwsLambda` function and an `AwsKinesisStream`.

## Before you deploy

- The function and stream must be in the same region as the mapping.
- Stream sources require `startingPosition`; choose `LATEST` instead if you only want new records.

## Remix ideas

- Set `startingPosition: AT_TIMESTAMP` with `startingPositionTimestamp` to replay from a specific instant.
- Add `onFailureDestinationArn` pointing at an SQS queue for poison-pill batches.
- Raise `batchSize` together with `maximumBatchingWindowSeconds` for cost-efficient batching.
