# AwsLambdaEventSourceMapping

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsLambdaEventSourceMappingSpec defines an event source mapping: the
managed poller that reads records from an event source -- an SQS
queue, a Kinesis stream, a DynamoDB stream, an MSK topic, a
self-managed Kafka cluster, an Amazon MQ broker, or a DocumentDB
change stream -- batches them, and invokes a Lambda function.

The mapping is its own node in the resource graph, not a setting of
the function: it has independent AWS identity (a server-assigned
UUID), a function may have many mappings, a mapping can be repointed
to a different function in place, and it is the edge that wires the
messaging graph (queues, streams, clusters) into compute. Pausing
consumption (disabled), tuning batching, and re-driving failures all
happen here without touching the function or the source.

Batching guidance: the defaults (10 records for SQS, 100 for
streams/Kafka, up to a 0-second batching window) optimize latency.
For throughput and cost, raise batch_size together with
maximum_batching_window_seconds so the poller can actually fill the
larger batches.

## Example

```yaml
# Canonical validated example: an SQS worker mapping with per-record
# failure reporting, both billable CloudWatch metrics, and a per-mapping
# concurrency throttle (the ceiling is the function's own concurrency --
# values above the former 1,000 cap are valid). The Kafka-only arms
# (provisioned pollers with a shared poller group, schema registry) live
# in the msk-shared-pollers preset, which the offline plans also exercise.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsLambdaEventSourceMapping
metadata:
  name: test-sqs-worker
spec:
  region: us-west-2
  functionArn:
    value: arn:aws:lambda:us-west-2:123456789012:function:orders-worker
  eventSourceArn:
    value: arn:aws:sqs:us-west-2:123456789012:orders-queue
  functionResponseTypes:
    - ReportBatchItemFailures
  metrics:
    - EventCount
    - ErrorCount
  scalingMaxConcurrency: 2500
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.functionArn` | `string \| valueFrom` | yes |  | AwsLambda (`status.outputs.function_arn`) |
| `spec.eventSourceArn` | `string \| valueFrom` |  |  | AwsSqsQueue (`status.outputs.queue_arn`) |
| `spec.selfManagedKafka` | `AwsLambdaEventSourceMappingSelfManagedKafka` |  |  |  |
| `spec.selfManagedKafka.bootstrapServers` | `[]string` | yes |  |  |
| `spec.disabled` | `bool` |  |  |  |
| `spec.batchSize` | `int32` |  |  |  |
| `spec.maximumBatchingWindowSeconds` | `int32` |  |  |  |
| `spec.filters` | `[]AwsLambdaEventSourceMappingFilter` |  |  |  |
| `spec.filters[].pattern` | `string` |  |  |  |
| `spec.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.functionResponseTypes` | `[]string` |  |  |  |
| `spec.scalingMaxConcurrency` | `int32` |  |  |  |
| `spec.metrics` | `[]string` |  |  |  |
| `spec.startingPosition` | `string` |  |  |  |
| `spec.startingPositionTimestamp` | `string` |  |  |  |
| `spec.parallelizationFactor` | `int32` |  |  |  |
| `spec.maximumRecordAgeSeconds` | `int32` |  |  |  |
| `spec.maximumRetryAttempts` | `int32` |  |  |  |
| `spec.bisectBatchOnFunctionError` | `bool` |  |  |  |
| `spec.tumblingWindowSeconds` | `int32` |  |  |  |
| `spec.onFailureDestinationArn` | `string \| valueFrom` |  |  | AwsSqsQueue (`status.outputs.queue_arn`) |
| `spec.topics` | `[]string` |  |  |  |
| `spec.kafkaConsumerGroupId` | `string` |  |  |  |
| `spec.sourceAccessConfigurations` | `[]AwsLambdaEventSourceMappingSourceAccess` |  |  |  |
| `spec.sourceAccessConfigurations[].type` | `string` | yes |  |  |
| `spec.sourceAccessConfigurations[].uri` | `string` | yes |  |  |
| `spec.schemaRegistry` | `AwsLambdaEventSourceMappingSchemaRegistry` |  |  |  |
| `spec.schemaRegistry.uri` | `string` | yes |  |  |
| `spec.schemaRegistry.eventRecordFormat` | `string` | yes |  |  |
| `spec.schemaRegistry.validationAttributes` | `[]string` |  |  |  |
| `spec.schemaRegistry.accessConfigurations` | `[]AwsLambdaEventSourceMappingSourceAccess` |  |  |  |
| `spec.schemaRegistry.accessConfigurations[].type` | `string` | yes |  |  |
| `spec.schemaRegistry.accessConfigurations[].uri` | `string` | yes |  |  |
| `spec.provisionedPollers` | `AwsLambdaEventSourceMappingProvisionedPollers` |  |  |  |
| `spec.provisionedPollers.minimumPollers` | `int32` |  |  |  |
| `spec.provisionedPollers.maximumPollers` | `int32` |  |  |  |
| `spec.provisionedPollers.pollerGroupName` | `string` | yes |  |  |
| `spec.mqQueue` | `string` |  |  |  |
| `spec.documentDb` | `AwsLambdaEventSourceMappingDocumentDb` |  |  |  |
| `spec.documentDb.databaseName` | `string` | yes |  |  |
| `spec.documentDb.collectionName` | `string` |  |  |  |
| `spec.documentDb.fullDocument` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the mapping is created in -- must be the region of
both the function and the event source.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.functionArn

`string | valueFrom` · required

The Lambda function the mapping invokes. Repointing a live mapping
to a different function is an in-place update -- consumption
continues from the tracked position. Reference an AwsLambda
function_arn output or pass a literal function ARN.

- references: AwsLambda (`status.outputs.function_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsLambda, name: <that resource's name>, fieldPath: status.outputs.function_arn}} -- a bare string does not parse

### spec.eventSourceArn

`string | valueFrom`

The ARN of the AWS event source: an SQS queue (queue_arn), a
Kinesis stream (stream_arn), a DynamoDB stream (the table's
stream_arn output), an MSK cluster (cluster_arn -- also set
topics), an Amazon MQ broker (also set queue), or a DocumentDB
change stream. Create-time immutable -- a different source is a
different mapping. The default reference targets AwsSqsQueue;
reference other kinds explicitly (e.g. kind AwsKinesisStream
fieldPath status.outputs.stream_arn, kind AwsDynamodb fieldPath
status.outputs.stream_arn, kind AwsMskCluster fieldPath
status.outputs.cluster_arn) or pass a literal ARN.

- references: AwsSqsQueue (`status.outputs.queue_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSqsQueue, name: <that resource's name>, fieldPath: status.outputs.queue_arn}} -- a bare string does not parse

### spec.selfManagedKafka

`AwsLambdaEventSourceMappingSelfManagedKafka`

A self-managed (non-MSK) Apache Kafka cluster as the event source.
Create-time immutable.

### spec.selfManagedKafka.bootstrapServers

`[]string` · required

Kafka bootstrap servers as "host:port" pairs, e.g.
["kafka1.example.com:9092", "kafka2.example.com:9092"].

- rule: {"repeated":{"minItems":"1"}}

### spec.disabled

`bool`

Pause consumption without deleting the mapping -- the tracked
position is retained, so re-enabling resumes where it stopped.
False (the default) keeps the mapping actively polling.

### spec.batchSize

`int32`

Records per invocation batch. 0 keeps the source's AWS default
(10 for SQS and DocumentDB, 100 for Kinesis/DynamoDB/Kafka/MQ).
Ceilings are per-source (10,000 for streams and SQS FIFO/standard
with long batching; 10 for MQ) -- AWS validates the exact bound.
Batches above 10 records for SQS require a batching window.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.maximumBatchingWindowSeconds

`int32`

How long (seconds, 0-300) the poller gathers records before
invoking, trading latency for fuller batches. 0 keeps the AWS
default (invoke as soon as records are available).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":300,"gte":1}}

### spec.filters

`[]AwsLambdaEventSourceMappingFilter`

Event-pattern filters applied BEFORE invocation -- records matching
no filter are discarded without billing function time. Up to 10
patterns, OR-ed together, each an EventBridge-style JSON pattern
against the record shape of the source.

- rule: {"repeated":{"maxItems":"10"}}

### spec.filters[].pattern

`string`

An EventBridge-style JSON pattern matched against the record,
e.g. {"body":{"type":["order.created"]}} for an SQS source.

- rule: {"string":{"maxLen":"4096"}}

### spec.kmsKeyArn

`string | valueFrom`

The KMS key that encrypts the filter criteria at rest. Empty uses
an AWS-owned key. Reference an AwsKmsKey key_arn output or pass a
literal key ARN.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.functionResponseTypes

`[]string`

Report per-record failures from the function ("ReportBatchItemFailures")
so only the failed records of a batch are retried instead of the
whole batch -- the right setting for almost every SQS and stream
consumer that processes records independently.

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["ReportBatchItemFailures"]}}}}

### spec.scalingMaxConcurrency

`int32`

Maximum concurrent function invocations this SQS mapping may drive
-- a per-mapping throttle below the function's own concurrency.
Minimum 2; the effective ceiling is the function's concurrency
(AWS validates it at deploy time -- it routinely exceeds any fixed
bound, so none is imposed here). 0 leaves scaling to AWS. SQS
sources only.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":2}}

### spec.metrics

`[]string`

Emit the mapping's CloudWatch metrics: "EventCount" (records
delivered to the function), "ErrorCount" (records that failed
processing), and/or "KafkaMetrics" (poller/consumer-lag metrics --
Kafka sources only). Off by default; metrics are billed.

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["EventCount","ErrorCount","KafkaMetrics"]}}}}

### spec.startingPosition

`string`

Where to start reading a stream source, create-time immutable:
"TRIM_HORIZON" (oldest available record -- process the backlog),
"LATEST" (only new records), or "AT_TIMESTAMP" (from
starting_position_timestamp; Kinesis only). Required for Kinesis,
DynamoDB Streams, MSK, self-managed Kafka, and DocumentDB; must
stay empty for SQS and MQ.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["TRIM_HORIZON","LATEST","AT_TIMESTAMP"]}}

### spec.startingPositionTimestamp

`string`

The UTC RFC3339 instant to start reading from (e.g.
"2026-07-04T00:00:00Z"). Required with (and only meaningful for)
starting_position AT_TIMESTAMP.

### spec.parallelizationFactor

`int32`

Concurrent batches processed per shard, 1-10 -- multiplies
per-shard throughput while preserving per-partition-key ordering.
0 keeps the AWS default (1). Kinesis and DynamoDB streams only.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":10,"gte":1}}

### spec.maximumRecordAgeSeconds

`int32`

Discard records older than this (seconds, 60-604800), or -1 (the
AWS default) to never age records out. 0 keeps the AWS default.
Stream sources only.

### spec.maximumRetryAttempts

`int32` · optional (explicit presence)

Retries per failed batch before the batch is discarded (or sent
to on_failure_destination_arn), 0-10000, or -1 (the AWS default)
to retry until the records expire. 0 means no retries. Stream
sources only.

- rule: {"int32":{"lte":10000,"gte":-1}}

### spec.bisectBatchOnFunctionError

`bool`

On a function error, split the failing batch in two and retry the
halves -- isolates a single poison record at the cost of
reprocessing its neighbors (make the function idempotent). Stream
sources only.

### spec.tumblingWindowSeconds

`int32`

Group records into fixed windows (seconds, 0-900) for streaming
aggregations -- the function receives a rolling state across the
window. Stream sources only.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":900,"gte":1}}

### spec.onFailureDestinationArn

`string | valueFrom`

Where discarded batches' metadata is sent after retries are
exhausted: an SQS queue or SNS topic ARN (Kafka sources may also
target an S3 bucket). Reference an AwsSqsQueue queue_arn output or
pass an explicit-kind reference / literal ARN. Stream and Kafka
sources only.

- references: AwsSqsQueue (`status.outputs.queue_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSqsQueue, name: <that resource's name>, fieldPath: status.outputs.queue_arn}} -- a bare string does not parse

### spec.topics

`[]string`

The Kafka topics to consume, 1-249 characters each. Required for
MSK and self-managed Kafka sources; must stay empty otherwise.

- rule: {"repeated":{"unique":true,"items":{"string":{"minLen":"1","maxLen":"249"}}}}

### spec.kafkaConsumerGroupId

`string`

The Kafka consumer group ID to join, create-time immutable. Empty
lets AWS generate one. Setting it lets the mapping resume an
existing group's committed offsets (starting_position is then
ignored for partitions with committed offsets).

- rule: {"string":{"maxLen":"200"}}

### spec.sourceAccessConfigurations

`[]AwsLambdaEventSourceMappingSourceAccess`

Authentication and network access the poller uses to reach the
source: SASL credentials, mTLS client certificates, VPC subnets
and security groups (self-managed Kafka), or broker credentials
(MQ). Each entry pairs a type with the Secrets Manager secret ARN
or VPC resource URI it points at.

- rule: {"repeated":{"maxItems":"22"}}

### spec.sourceAccessConfigurations[].type

`string` · required

The access type: "BASIC_AUTH" (MQ / SASL PLAIN),
"SASL_SCRAM_256_AUTH" / "SASL_SCRAM_512_AUTH" (Kafka SASL),
"CLIENT_CERTIFICATE_TLS_AUTH" (Kafka mTLS),
"SERVER_ROOT_CA_CERTIFICATE" (private CA), "VPC_SUBNET" /
"VPC_SECURITY_GROUP" (self-managed Kafka networking), or
"VIRTUAL_HOST" (RabbitMQ).

- rule: {"required":true,"string":{"in":["BASIC_AUTH","VPC_SUBNET","VPC_SECURITY_GROUP","SASL_SCRAM_512_AUTH","SASL_SCRAM_256_AUTH","VIRTUAL_HOST","CLIENT_CERTIFICATE_TLS_AUTH","SERVER_ROOT_CA_CERTIFICATE"]}}

### spec.sourceAccessConfigurations[].uri

`string` · required

The value for the type: a Secrets Manager secret ARN (auth types),
"subnet:<subnet-id>" / "security_group:<sg-id>" (VPC types), or
the virtual host name (VIRTUAL_HOST).

- rule: {"string":{"minLen":"1"}}

### spec.schemaRegistry

`AwsLambdaEventSourceMappingSchemaRegistry`

Confluent / Glue schema registry integration for Kafka sources:
the poller validates and deserializes records against registered
schemas before invoking the function.

### spec.schemaRegistry.uri

`string` · required

The registry location: an AWS Glue schema registry ARN or a
Confluent registry HTTPS URL.

- rule: {"string":{"minLen":"1"}}

### spec.schemaRegistry.eventRecordFormat

`string` · required

What the function receives: "JSON" (deserialized to JSON) or
"SOURCE" (the original serialized bytes with the schema header
stripped).

- rule: {"required":true,"string":{"in":["JSON","SOURCE"]}}

### spec.schemaRegistry.validationAttributes

`[]string`

Which record parts are validated against the registry: "KEY"
and/or "VALUE".

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["KEY","VALUE"]}}}}

### spec.schemaRegistry.accessConfigurations

`[]AwsLambdaEventSourceMappingSourceAccess`

How the poller authenticates to the registry: type "BASIC_AUTH" /
"CLIENT_CERTIFICATE_TLS_AUTH" / "SERVER_ROOT_CA_CERTIFICATE" with
the Secrets Manager secret ARN in uri.

### spec.schemaRegistry.accessConfigurations[].type

`string` · required

The access type: "BASIC_AUTH" (MQ / SASL PLAIN),
"SASL_SCRAM_256_AUTH" / "SASL_SCRAM_512_AUTH" (Kafka SASL),
"CLIENT_CERTIFICATE_TLS_AUTH" (Kafka mTLS),
"SERVER_ROOT_CA_CERTIFICATE" (private CA), "VPC_SUBNET" /
"VPC_SECURITY_GROUP" (self-managed Kafka networking), or
"VIRTUAL_HOST" (RabbitMQ).

- rule: {"required":true,"string":{"in":["BASIC_AUTH","VPC_SUBNET","VPC_SECURITY_GROUP","SASL_SCRAM_512_AUTH","SASL_SCRAM_256_AUTH","VIRTUAL_HOST","CLIENT_CERTIFICATE_TLS_AUTH","SERVER_ROOT_CA_CERTIFICATE"]}}

### spec.schemaRegistry.accessConfigurations[].uri

`string` · required

The value for the type: a Secrets Manager secret ARN (auth types),
"subnet:<subnet-id>" / "security_group:<sg-id>" (VPC types), or
the virtual host name (VIRTUAL_HOST).

- rule: {"string":{"minLen":"1"}}

### spec.provisionedPollers

`AwsLambdaEventSourceMappingProvisionedPollers`

Dedicated pollers for Kafka sources: pin the fleet between
minimum_pollers and maximum_pollers for predictable throughput
instead of AWS's reactive scaling.

- rule: maximum_pollers must be greater than or equal to minimum_pollers

### spec.provisionedPollers.minimumPollers

`int32`

The floor of always-running pollers, 1-200. 0 keeps the AWS
default (1).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":200,"gte":1}}

### spec.provisionedPollers.maximumPollers

`int32`

The ceiling of pollers AWS may scale to, 1-2000. 0 keeps the AWS
default (200).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":2000,"gte":1}}

### spec.provisionedPollers.pollerGroupName

`string` · required

Share one provisioned poller fleet across mappings by naming a
poller group, 1-128 characters: every mapping naming the same group
draws from (and is jointly capped by) one fleet instead of
provisioning its own -- the cost lever for many low-traffic topics.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"minLen":"1","maxLen":"128"}}

### spec.mqQueue

`string`

The Amazon MQ broker queue to consume (exactly one). Required for
MQ sources; must stay empty otherwise.

- rule: {"string":{"maxLen":"1000"}}

### spec.documentDb

`AwsLambdaEventSourceMappingDocumentDb`

DocumentDB change-stream options. Required for DocumentDB sources.

### spec.documentDb.databaseName

`string` · required

The database whose change stream is consumed.

- rule: {"string":{"minLen":"1","maxLen":"63"}}

### spec.documentDb.collectionName

`string`

Consume a single collection's changes. Empty consumes the whole
database.

- rule: {"string":{"maxLen":"57"}}

### spec.documentDb.fullDocument

`string`

What change events carry: "UpdateLookup" (the full current
document alongside the change) or "Default" (the change delta
only). Empty keeps the AWS default (Default).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["UpdateLookup","Default"]}}

## Validation Rules

- `exactly_one_event_source`: provide the event source as exactly one of event_source_arn (SQS/Kinesis/DynamoDB/MSK/MQ/DocumentDB) or self_managed_kafka
- `timestamp_requires_at_timestamp`: starting_position_timestamp is required with starting_position AT_TIMESTAMP and must stay empty otherwise
- `self_managed_kafka_needs_position_and_topics`: self-managed Kafka sources require starting_position and at least one topic
- `record_age_semantics`: maximum_record_age_seconds must be -1 (never age out) or between 60 and 604800
- `provisioned_pollers_are_kafka_only`: provisioned_pollers applies only to Kafka sources (MSK or self-managed) -- other sources scale automatically
- `schema_registry_is_kafka_only`: schema_registry applies only to Kafka sources (MSK or self-managed)
- `mq_queue_excludes_kafka`: mq_queue (Amazon MQ) and topics (Kafka) are different source families -- set only the one matching the event source

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsLambdaEventSourceMapping, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.uuid` | `string` | The server-assigned mapping UUID -- the identity AWS APIs (GetEventSourceMapping, UpdateEventSourceMapping) key on. |
| `status.outputs.mapping_arn` | `string` | The mapping ARN. |
| `status.outputs.function_arn` | `string` | The ARN of the function the mapping invokes, as resolved by AWS. |
| `status.outputs.state` | `string` | The mapping state as last observed at deploy time (e.g. "Enabled", "Disabled") -- transitional states settle asynchronously. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.functionArn` | AwsLambda | `status.outputs.function_arn` |
| `spec.eventSourceArn` | AwsSqsQueue | `status.outputs.queue_arn` |
| `spec.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.onFailureDestinationArn` | AwsSqsQueue | `status.outputs.queue_arn` |

## See Also

- [Overview](../README.md)
