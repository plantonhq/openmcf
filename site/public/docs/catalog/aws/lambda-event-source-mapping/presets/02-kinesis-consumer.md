# Kinesis Stream Consumer

A Kinesis-to-Lambda mapping that processes the backlog from the oldest available record with partial-batch failure reporting.

## What this preset gives you

- A stream consumer starting at `TRIM_HORIZON` — processes existing records, not just new ones.
- `ReportBatchItemFailures` for independent per-record retry semantics.
- `parallelization_factor: 2` to double per-shard throughput while preserving partition-key ordering.
- Composed references to an `AwsLambda` function and an `AwsKinesisStream`.

## Before you deploy

- The function and stream must be in the same region as the mapping.
- Stream sources require `starting_position`; choose `LATEST` instead if you only want new records.

## Remix ideas

- Set `starting_position: AT_TIMESTAMP` with `starting_position_timestamp` to replay from a specific instant.
- Add `on_failure_destination_arn` pointing at an SQS queue for poison-pill batches.
- Raise `batch_size` together with `maximum_batching_window_seconds` for cost-efficient batching.
