# SQS Queue Worker

An SQS-to-Lambda mapping with partial-batch failure reporting -- the right default for independent record processing.

## What this preset gives you

- A managed poller that reads from an SQS queue and invokes your Lambda function in batches.
- `ReportBatchItemFailures` so only failed messages are retried instead of the whole batch.
- Composed references to an `AwsLambda` function and an `AwsSqsQueue` -- no literal ARNs to wire by hand.

## Before you deploy

- The referenced Lambda function and SQS queue must exist in the same region as the mapping.
- Your function handler must return partial batch failure responses when using `ReportBatchItemFailures`.

## Remix ideas

- Add `batchSize` and `maximumBatchingWindowSeconds` together for higher throughput.
- Set `scalingMaxConcurrency` to throttle how many concurrent invocations this mapping drives.
- Add `filters` with EventBridge-style JSON patterns to discard records before invocation.
