# AWS Lambda Event Source Mapping

Deploys a Lambda event source mapping -- the managed poller that reads an SQS queue, a Kinesis or DynamoDB stream, a Kafka topic, an Amazon MQ queue, or a DocumentDB change stream and invokes a Lambda function with batched records.

## What Gets Created

When you deploy an AwsLambdaEventSourceMapping resource, Planton provisions:

- **Event Source Mapping** -- an `aws_lambda_event_source_mapping` wired to the referenced function and event source, with optional batching, filtering, failure handling, and consumption controls

The event source itself is create-time immutable; batching, filters, and the target function edit in place.

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A Lambda function** in the same region (can be an `AwsLambda` resource)
- **An event source** in the same region (queue, stream, cluster, broker, or change stream)

## Quick Start

Create a file `sqs-worker.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLambdaEventSourceMapping
metadata:
  name: my-worker
spec:
  region: us-west-2
  functionArn:
    value: arn:aws:lambda:us-west-2:123456789012:function:my-function
  eventSourceArn:
    value: arn:aws:sqs:us-west-2:123456789012:my-queue
  functionResponseTypes:
    - ReportBatchItemFailures
```

Deploy:

```shell
planton apply -f sqs-worker.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region for the mapping. | Required; non-empty |
| `functionArn` | `StringValueOrRef` | Lambda function to invoke. | Required |
| `eventSourceArn` or `selfManagedKafka` | `StringValueOrRef` / `object` | Exactly one event source. | CEL-enforced |

### Common Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `disabled` | `bool` | `false` | Pause consumption without deleting the mapping. |
| `batchSize` | `int32` | source default | Records per invocation batch. |
| `maximumBatchingWindowSeconds` | `int32` | `0` | Batching window in seconds (1–300). |
| `functionResponseTypes` | `string[]` | — | `ReportBatchItemFailures` for partial-batch retry. |
| `startingPosition` | `string` | — | Stream/Kafka start position. ForceNew. |

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `uuid` | `string` | Server-assigned mapping UUID |
| `mapping_arn` | `string` | The mapping ARN |
| `function_arn` | `string` | Function ARN as resolved by AWS |
| `state` | `string` | Mapping state at deploy time |

## Related Components

- [AwsLambda](/docs/catalog/aws/awslambda) -- the function this mapping invokes
- [AwsSqsQueue](/docs/catalog/aws/awssqsqueue) -- common SQS event source
- [AwsKinesisStream](/docs/catalog/aws/awskinesisstream) -- Kinesis stream event source
- [AwsKmsKey](/docs/catalog/aws/awskmskey) -- encrypts filter criteria at rest
