---
title: "Kinesis Firehose"
description: "Kinesis Firehose deployment documentation"
icon: "package"
order: 100
componentName: "awskinesisfirehose"
---

# AWS Kinesis Firehose

Deploys an Amazon Kinesis Data Firehose delivery stream that captures, transforms, and delivers streaming data to one of eight destinations: S3, OpenSearch, OpenSearch Serverless, HTTP endpoints, Redshift, Splunk, Snowflake, or Apache Iceberg tables. The component supports Direct PUT, Kinesis Data Stream, and Amazon MSK sources; an ordered record-transformation pipeline (Lambda, JQ metadata extraction, decompression, CloudWatch-Logs unwrapping, delimiters, de-aggregation); dynamic partitioning; Parquet/ORC format conversion; and Secrets Manager-sourced destination credentials.

## What Gets Created

When you deploy an AwsKinesisFirehose resource, Planton provisions:

- **Kinesis Firehose Delivery Stream** — the core `aws_kinesis_firehose_delivery_stream` resource configured with the selected destination type
- **Kinesis source configuration** — created only when `kinesisStreamSource` is set, configures Firehose to consume from an existing Kinesis Data Stream with automatic checkpointing and retry
- **MSK source configuration** — created only when `mskSource` is set, configures Firehose as a managed Kafka consumer of a topic on an MSK cluster (IAM auth, private or public connectivity)
- **Server-side encryption** — created only when `sseEnabled` is `true`, encrypts data at rest in the delivery stream buffer using AWS-owned or customer-managed KMS keys (Direct PUT sources only)
- **Extended S3 destination** — primary S3 delivery with optional GZIP/Snappy compression, record processing, dynamic partitioning, and Parquet/ORC format conversion via AWS Glue Data Catalog
- **OpenSearch destination** — direct indexing into an OpenSearch domain with configurable index rotation, document-ID control, VPC delivery, and S3 backup for failed documents
- **OpenSearch Serverless destination** — indexing into an OpenSearch Serverless collection endpoint with S3 backup
- **HTTP endpoint destination** — HTTPS delivery to any endpoint (Datadog, New Relic, Sumo Logic, custom APIs) with S3 backup for failed records
- **Redshift destination** — S3 staging followed by a Redshift COPY command for bulk data warehouse loading
- **Splunk destination** — HEC delivery with indexer acknowledgment and S3 backup for failed events
- **Snowflake destination** — Snowpipe Streaming inserts with key-pair authentication and optional PrivateLink
- **Iceberg destination** — direct snapshot commits into Glue-cataloged Apache Iceberg tables with multi-table routing and unique-key upserts
- **Secrets Manager wiring** — created only when a destination's `secretsManager` block is set, sources the destination credential from AWS Secrets Manager instead of the manifest

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A destination resource** — an S3 bucket, OpenSearch domain or serverless collection, HTTPS endpoint, Redshift cluster, Splunk HEC endpoint, Snowflake table, or Glue-cataloged Iceberg table depending on the chosen destination type
- **An IAM role** with permissions appropriate for the destination (S3 write, OpenSearch index, Redshift COPY, Glue table update, Secrets Manager read, etc.)
- **A Kinesis Data Stream or MSK cluster** if using a stream source instead of Direct PUT
- **An AWS Glue Data Catalog database and table** if enabling Parquet/ORC data format conversion
- **VPC subnets and security groups** if delivering to a VPC-deployed OpenSearch domain

## Quick Start

Create a file `firehose.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKinesisFirehose
metadata:
  name: my-firehose
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsKinesisFirehose.my-firehose
spec:
  region: us-east-1
  extendedS3:
    bucketArn: arn:aws:s3:::my-data-bucket
    roleArn: arn:aws:iam::123456789012:role/firehose-s3-role
```

Deploy:

```shell
planton apply -f firehose.yaml
```

This creates a Direct PUT delivery stream that writes raw records to S3 with no compression or transformation.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the Firehose delivery stream will be created (e.g., `us-west-2`, `eu-west-1`). | Required; non-empty |

Exactly one destination must be configured. The destination type is ForceNew — changing it requires replacing the delivery stream.

| Field | Type | Description |
|-------|------|-------------|
| `extendedS3` | `object` | Extended S3 destination for data lake storage with compression, partitioning, and format conversion |
| `opensearch` | `object` | OpenSearch destination for direct indexing with S3 backup |
| `opensearchServerless` | `object` | OpenSearch Serverless destination for indexing into a collection endpoint |
| `httpEndpoint` | `object` | HTTP endpoint destination for HTTPS delivery with S3 backup |
| `redshift` | `object` | Redshift destination for S3 staging + COPY command |
| `splunk` | `object` | Splunk destination for HEC delivery with indexer acknowledgment |
| `snowflake` | `object` | Snowflake destination for Snowpipe Streaming inserts |
| `iceberg` | `object` | Iceberg destination for direct commits into Glue-cataloged tables |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `kinesisStreamSource` | `object` | — | Kinesis Data Stream source configuration. Mutually exclusive with `mskSource`. When both absent, the delivery stream uses Direct PUT. ForceNew. |
| `kinesisStreamSource.streamArn` | `string` | — | ARN of the Kinesis Data Stream to consume from. Required when `kinesisStreamSource` is set. Can reference an AwsKinesisStream resource via `valueFrom`. |
| `kinesisStreamSource.roleArn` | `string` | — | IAM role ARN granting Firehose read access to the Kinesis stream. Required when `kinesisStreamSource` is set. Can reference an AwsIamRole resource via `valueFrom`. |
| `mskSource` | `object` | — | Amazon MSK source configuration. Mutually exclusive with `kinesisStreamSource`. ForceNew. |
| `mskSource.mskClusterArn` | `string` | — | ARN of the MSK cluster (IAM access control required). Required when `mskSource` is set. Can reference an AwsMskCluster resource via `valueFrom`. |
| `mskSource.topicName` | `string` | — | Kafka topic to read from. Required when `mskSource` is set. |
| `mskSource.connectivity` | `string` | — | `PRIVATE` (in-VPC brokers) or `PUBLIC`. Required when `mskSource` is set. |
| `mskSource.roleArn` | `string` | — | IAM role with kafka + kafka-cluster data-plane permissions. Required when `mskSource` is set. Can reference an AwsIamRole resource via `valueFrom`. |
| `mskSource.readFromTimestamp` | `string` | latest offset | RFC 3339 point in time to start reading the topic from. |
| `sseEnabled` | `bool` | `false` | Enables server-side encryption for data at rest in the delivery stream buffer. Only valid for Direct PUT sources. |
| `sseKmsKeyArn` | `string` | — | Customer-managed KMS key ARN for SSE. When absent, uses the AWS-owned CMK. Requires `sseEnabled` to be `true`. Can reference an AwsKmsKey resource via `valueFrom`. |

### Extended S3 Destination Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `extendedS3.bucketArn` | `string` | — | **(Required)** S3 bucket ARN for delivery. Can reference an AwsS3Bucket resource via `valueFrom`. |
| `extendedS3.roleArn` | `string` | — | **(Required)** IAM role ARN granting Firehose write access to S3, KMS, Lambda, and Glue as needed. Can reference an AwsIamRole resource via `valueFrom`. |
| `extendedS3.prefix` | `string` | — | S3 key prefix. Supports Firehose expression syntax (e.g., `year=!{timestamp:yyyy}/`). |
| `extendedS3.errorOutputPrefix` | `string` | — | S3 key prefix for records that fail transformation or delivery. |
| `extendedS3.compressionFormat` | `string` | `UNCOMPRESSED` | Compression applied before writing. Valid: `UNCOMPRESSED`, `GZIP`, `ZIP`, `Snappy`, `HADOOP_SNAPPY`. |
| `extendedS3.kmsKeyArn` | `string` | — | KMS key ARN for S3 server-side encryption (SSE-KMS). Can reference an AwsKmsKey resource via `valueFrom`. |
| `extendedS3.buffering` | `object` | `300s / 5 MiB` | Buffering hints: `intervalInSeconds` (0–900) and `sizeInMbs` (1–128). |
| `extendedS3.customTimeZone` | `string` | `UTC` | IANA time zone for S3 prefix timestamp expressions. |
| `extendedS3.fileExtension` | `string` | — | File extension appended to delivered objects (e.g., `.json`, `.parquet`). Must start with a period. |
| `extendedS3.s3BackupMode` | `string` | `Disabled` | When `Enabled`, a copy of pre-transformation records is written to `s3Backup`. |
| `extendedS3.s3Backup` | `object` | — | S3 configuration for source record backup. Required when `s3BackupMode` is `Enabled`. |
| `extendedS3.processing` | `object` | — | Record-transformation pipeline (see Processing Pipeline below). |
| `extendedS3.logging` | `object` | — | CloudWatch error logging. Set `enabled`, `logGroupName`, and `logStreamName`. |
| `extendedS3.dynamicPartitioning` | `object` | — | Dynamic partitioning by record fields for efficient querying with Athena/Spark. ForceNew. Define partition keys with a `metadataExtraction` processor. |
| `extendedS3.dataFormatConversion` | `object` | — | JSON-to-Parquet/ORC conversion via AWS Glue Data Catalog schema. |

### Processing Pipeline (all destinations)

Every destination accepts a `processing` block: `enabled` plus an ordered `processors` list. Each processor entry sets exactly one typed arm:

| Arm | Purpose | Key fields |
|-----|---------|-----------|
| `lambda` | Transform records with a Lambda function | `lambdaArn` (can reference an AwsLambda via `valueFrom`), `bufferSizeInMbs` (1–3), `bufferIntervalInSeconds` (60–900), `numberOfRetries` (0–300) |
| `metadataExtraction` | Extract partition keys from JSON with a JQ expression | `query` (e.g., `{customer_id: .customer_id}`), `jsonParsingEngine` (`JQ-1.6`) |
| `decompression` | Decompress GZIP payloads (CloudWatch Logs subscriptions) | `compressionFormat` (`GZIP`) |
| `cloudwatchLogProcessing` | Unwrap CloudWatch Logs subscription envelopes | `dataMessageExtraction` (bool) |
| `appendDelimiter` | Newline-delimit records for JSON-lines output (extended_s3 only) | `delimiter` (e.g., `"\\n"`) |
| `recordDeaggregation` | Split KPL-aggregated or delimited payloads (extended_s3 only) | `subRecordType` (`JSON`/`DELIMITED`), `delimiter` (base64) |

### OpenSearch Destination Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `opensearch.domainArn` | `string` | — | ARN of the OpenSearch domain. Mutually exclusive with `clusterEndpoint`. Can reference an AwsOpenSearchDomain resource via `valueFrom`. |
| `opensearch.clusterEndpoint` | `string` | — | OpenSearch cluster endpoint URL. Mutually exclusive with `domainArn`. |
| `opensearch.indexName` | `string` | — | **(Required)** Index name (or prefix when rotation is enabled). |
| `opensearch.roleArn` | `string` | — | **(Required)** IAM role ARN with `es:ESHttpPut` and `es:ESHttpGet` permissions. Can reference an AwsIamRole resource via `valueFrom`. |
| `opensearch.s3Config` | `object` | — | **(Required)** S3 configuration for backing up failed or all documents. |
| `opensearch.indexRotationPeriod` | `string` | `OneDay` | Index rotation period. Valid: `NoRotation`, `OneHour`, `OneDay`, `OneWeek`, `OneMonth`. |
| `opensearch.typeName` | `string` | — | Document type name. Only relevant for Elasticsearch 6.x and earlier. |
| `opensearch.defaultDocumentIdFormat` | `string` | `FIREHOSE_DEFAULT` | `FIREHOSE_DEFAULT` (dedupe-safe retries) or `NO_DOCUMENT_ID` (higher indexing throughput). |
| `opensearch.buffering` | `object` | `300s / 5 MiB` | Buffering hints. Max size: 100 MiB for OpenSearch destinations. |
| `opensearch.retryDurationInSeconds` | `int` | `300` | Retry duration for failed index requests. Range: 0–7200. |
| `opensearch.s3BackupMode` | `string` | `FailedDocumentsOnly` | Valid: `FailedDocumentsOnly`, `AllDocuments`. ForceNew. |
| `opensearch.processing` | `object` | — | Record-transformation pipeline before indexing. |
| `opensearch.logging` | `object` | — | CloudWatch error logging for delivery failures. |
| `opensearch.vpcConfig` | `object` | — | VPC configuration for VPC-deployed OpenSearch domains. ForceNew. |

### OpenSearch Serverless Destination Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `opensearchServerless.collectionEndpoint` | `string` | — | **(Required)** Collection endpoint (`https://<id>.<region>.aoss.amazonaws.com`). |
| `opensearchServerless.indexName` | `string` | — | **(Required)** Target index, admitted by the collection's data access policy. |
| `opensearchServerless.roleArn` | `string` | — | **(Required)** IAM role with `aoss:APIAccessAll` and a data-access-policy write grant. Can reference an AwsIamRole resource via `valueFrom`. |
| `opensearchServerless.s3Config` | `object` | — | **(Required)** S3 configuration for backing up failed or all documents. |
| `opensearchServerless.buffering` | `object` | `300s / 5 MiB` | Buffering hints. Max size: 100 MiB. |
| `opensearchServerless.retryDurationInSeconds` | `int` | `300` | Retry duration for failed index requests. Range: 0–7200. |
| `opensearchServerless.s3BackupMode` | `string` | `FailedDocumentsOnly` | Valid: `FailedDocumentsOnly`, `AllDocuments`. ForceNew. |
| `opensearchServerless.processing` | `object` | — | Record-transformation pipeline before indexing. |
| `opensearchServerless.logging` | `object` | — | CloudWatch error logging for delivery failures. |
| `opensearchServerless.vpcConfig` | `object` | — | VPC configuration for collections reached through a VPC endpoint. ForceNew. |

### HTTP Endpoint Destination Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `httpEndpoint.url` | `string` | — | **(Required)** HTTPS URL of the destination endpoint. Must start with `https://`. |
| `httpEndpoint.s3Config` | `object` | — | **(Required)** S3 configuration for backing up failed or all records. |
| `httpEndpoint.name` | `string` | — | Human-readable endpoint name for the AWS Console and CloudWatch metrics. |
| `httpEndpoint.accessKey` | `string` | — | Access key sent in the `X-Amz-Firehose-Access-Key` header. Sensitive. Mutually exclusive with `secretsManager`. |
| `httpEndpoint.secretsManager` | `object` | — | Source the access key from Secrets Manager (`secretArn` + optional `roleArn`; secret shape `{"api_key": "..."}`). Recommended for production. ForceNew. |
| `httpEndpoint.roleArn` | `string` | — | IAM role ARN for delivery and S3 backup. Can reference an AwsIamRole resource via `valueFrom`. |
| `httpEndpoint.buffering` | `object` | `300s / 5 MiB` | Buffering hints for HTTP delivery. Max size: 100 MiB. |
| `httpEndpoint.retryDurationInSeconds` | `int` | `300` | Retry duration for non-2xx responses. Range: 0–7200. |
| `httpEndpoint.s3BackupMode` | `string` | `FailedDataOnly` | Valid: `FailedDataOnly`, `AllData`. |
| `httpEndpoint.processing` | `object` | — | Record-transformation pipeline before HTTP delivery. |
| `httpEndpoint.logging` | `object` | — | CloudWatch error logging for delivery failures. |
| `httpEndpoint.requestConfig` | `object` | — | Request customization: `contentEncoding` (`NONE`, `GZIP`) and `commonAttributes` (key-value headers). |

### Redshift Destination Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `redshift.clusterJdbcurl` | `string` | — | **(Required)** JDBC URL of the Redshift cluster (e.g., `jdbc:redshift://host:5439/db`). |
| `redshift.roleArn` | `string` | — | **(Required)** IAM role ARN for S3 read and Redshift COPY. Can reference an AwsIamRole resource via `valueFrom`. |
| `redshift.dataTableName` | `string` | — | **(Required)** Target Redshift table for the COPY command. |
| `redshift.s3Config` | `object` | — | **(Required)** S3 configuration for intermediate staging. Firehose writes here before issuing COPY. |
| `redshift.dataTableColumns` | `string` | — | Comma-separated column names for the COPY command. When absent, COPY loads all columns in table order. |
| `redshift.copyOptions` | `string` | — | Additional COPY command options (e.g., `JSON 'auto'`, `GZIP`, `DELIMITER ','`). |
| `redshift.username` | `string` | — | Redshift database username. Set together with `password`; exactly one auth mode with `secretsManager`. |
| `redshift.password` | `string` | — | Redshift database password. Sensitive. Prefer `secretsManager`. |
| `redshift.secretsManager` | `object` | — | Source credentials from Secrets Manager (`secretArn` + optional `roleArn`; secret shape `{"username": ..., "password": ...}`). Recommended for production. ForceNew. |
| `redshift.retryDurationInSeconds` | `int` | `3600` | Retry duration for failed COPY commands. Range: 0–7200. |
| `redshift.s3BackupMode` | `string` | `Disabled` | When `Enabled`, original records are backed up to `s3Backup`. |
| `redshift.s3Backup` | `object` | — | S3 configuration for source record backup. Required when `s3BackupMode` is `Enabled`. |
| `redshift.processing` | `object` | — | Record-transformation pipeline before staging. |
| `redshift.logging` | `object` | — | CloudWatch error logging for COPY failures. |

### Splunk Destination Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `splunk.hecEndpoint` | `string` | — | **(Required)** HEC endpoint URL including port (must start with `https://`). |
| `splunk.hecEndpointType` | `string` | `Raw` | `Raw` (preformatted events) or `Event` (Splunk event JSON). |
| `splunk.hecToken` | `string` | — | HEC token. Sensitive. Exactly one auth mode with `secretsManager`. |
| `splunk.secretsManager` | `object` | — | Source the token from Secrets Manager (secret shape `{"hec_token": "..."}`). Recommended for production. ForceNew. |
| `splunk.hecAcknowledgmentTimeoutInSeconds` | `int` | `180` | Indexer acknowledgment wait: 180–600s. |
| `splunk.buffering` | `object` | `60s / 5 MiB` | Buffering hints — Splunk caps: 0–60s / 1–5 MiB. |
| `splunk.retryDurationInSeconds` | `int` | `3600` | Retry duration for failed or unacknowledged deliveries. Range: 0–7200. |
| `splunk.s3BackupMode` | `string` | `FailedEventsOnly` | Valid: `FailedEventsOnly`, `AllEvents`. |
| `splunk.s3Config` | `object` | — | **(Required)** S3 configuration for backing up failed or all events. |
| `splunk.processing` | `object` | — | Record-transformation pipeline before HEC delivery. |
| `splunk.logging` | `object` | — | CloudWatch error logging for delivery failures. |

### Snowflake Destination Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `snowflake.accountUrl` | `string` | — | **(Required)** `https://<account-identifier>.snowflakecomputing.com`. |
| `snowflake.database` / `snowflake.schema` / `snowflake.table` | `string` | — | **(Required)** Target table coordinates. |
| `snowflake.roleArn` | `string` | — | **(Required)** IAM role for S3 backup and Secrets Manager read. Can reference an AwsIamRole resource via `valueFrom`. |
| `snowflake.user` | `string` | — | Snowflake user owning the key pair. Set together with `privateKey`; exactly one auth mode with `secretsManager`. |
| `snowflake.privateKey` | `string` | — | RSA private key (PEM body without header/footer lines). Sensitive. Prefer `secretsManager`. |
| `snowflake.keyPassphrase` | `string` | — | Passphrase for an encrypted private key (7–255 chars). Sensitive. |
| `snowflake.secretsManager` | `object` | — | Source credentials from Secrets Manager (secret shape `{"user": ..., "private_key": ..., "key_passphrase": ...}`). Recommended for production. ForceNew. |
| `snowflake.dataLoadingOption` | `string` | `JSON_MAPPING` | `JSON_MAPPING`, `VARIANT_CONTENT_MAPPING`, or `VARIANT_CONTENT_AND_METADATA_MAPPING`. |
| `snowflake.contentColumnName` | `string` | — | VARIANT column for record content. Required for the VARIANT loading modes. |
| `snowflake.metadataColumnName` | `string` | — | VARIANT column for Firehose metadata. Required for content+metadata mode. |
| `snowflake.snowflakeRole` | `string` | user's default | Snowflake role to assume — a least-privilege ingestion role is recommended. |
| `snowflake.privateLinkVpceId` | `string` | — | PrivateLink VPCE ID for private connectivity. |
| `snowflake.buffering` | `object` | `0s / 1 MiB` | Buffering hints — Snowpipe Streaming defaults to near-real-time. |
| `snowflake.retryDurationInSeconds` | `int` | `60` | Retry duration for failed inserts. Range: 0–7200. |
| `snowflake.s3BackupMode` | `string` | `FailedDataOnly` | Valid: `FailedDataOnly`, `AllData`. |
| `snowflake.s3Config` | `object` | — | **(Required)** S3 configuration for backing up failed or all records. |
| `snowflake.processing` | `object` | — | Record-transformation pipeline before delivery. |
| `snowflake.logging` | `object` | — | CloudWatch error logging for delivery failures. |

### Iceberg Destination Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `iceberg.catalogArn` | `string` | — | **(Required)** Glue Data Catalog ARN (`arn:aws:glue:<region>:<account>:catalog`). ForceNew. |
| `iceberg.roleArn` | `string` | — | **(Required)** IAM role for Glue table read/update and warehouse S3 access. Can reference an AwsIamRole resource via `valueFrom`. |
| `iceberg.destinationTables` | `list` | — | Target Iceberg tables (`databaseName`, `tableName`, optional `s3ErrorOutputPrefix` and `uniqueKeys`). Multiple tables require per-record routing metadata from a processor. ForceNew. |
| `iceberg.appendOnly` | `bool` | `false` | Append-only mode — disables unique-key upserts. ForceNew. |
| `iceberg.buffering` | `object` | `300s / 5 MiB` | Buffering hints. |
| `iceberg.retryDurationInSeconds` | `int` | `300` | Retry duration for failed commits. Range: 0–7200. |
| `iceberg.s3BackupMode` | `string` | `FailedDataOnly` | Valid: `FailedDataOnly`, `AllData`. |
| `iceberg.s3Config` | `object` | — | **(Required)** S3 configuration for backing up failed or all records. |
| `iceberg.processing` | `object` | — | Record-transformation pipeline before the commit. |
| `iceberg.logging` | `object` | — | CloudWatch error logging for delivery failures. |

## Examples

### Extended S3 Data Lake

GZIP-compressed delivery to S3 with timestamp-based prefixes and buffering tuned for throughput:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKinesisFirehose
metadata:
  name: data-lake-firehose
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsKinesisFirehose.data-lake-firehose
spec:
  region: us-east-1
  extendedS3:
    bucketArn: arn:aws:s3:::my-data-lake-bucket
    roleArn: arn:aws:iam::123456789012:role/firehose-s3-delivery-role
    prefix: events/year=!{timestamp:yyyy}/month=!{timestamp:MM}/day=!{timestamp:dd}/
    compressionFormat: GZIP
    fileExtension: .json.gz
    buffering:
      intervalInSeconds: 120
      sizeInMbs: 64
```

### OpenSearch Log Analytics

Indexes application logs into an OpenSearch domain with daily index rotation and S3 backup for failed documents. References the OpenSearch domain via `valueFrom`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKinesisFirehose
metadata:
  name: log-analytics-firehose
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsKinesisFirehose.log-analytics-firehose
spec:
  region: us-east-1
  opensearch:
    domainArn:
      valueFrom:
        kind: AwsOpenSearchDomain
        name: my-log-domain
        fieldPath: status.outputs.domain_arn
    indexName: application-logs
    roleArn: arn:aws:iam::123456789012:role/firehose-opensearch-role
    indexRotationPeriod: OneDay
    s3BackupMode: FailedDocumentsOnly
    buffering:
      intervalInSeconds: 60
      sizeInMbs: 5
    s3Config:
      bucketArn: arn:aws:s3:::my-firehose-backup-bucket
      roleArn: arn:aws:iam::123456789012:role/firehose-s3-backup-role
      prefix: opensearch-backup/failed/
      compressionFormat: GZIP
```

### Production S3 with Kinesis Source and Parquet Conversion

Consumes from an existing Kinesis Data Stream, converts JSON to Parquet via AWS Glue Data Catalog, and writes columnar files to a partitioned S3 data lake:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKinesisFirehose
metadata:
  name: analytics-parquet-firehose
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsKinesisFirehose.analytics-parquet-firehose
spec:
  region: us-east-1
  kinesisStreamSource:
    streamArn:
      valueFrom:
        kind: AwsKinesisStream
        name: my-events-stream
        fieldPath: status.outputs.stream_arn
    roleArn:
      valueFrom:
        kind: AwsIamRole
        name: firehose-kinesis-consumer-role
        fieldPath: status.outputs.role_arn
  extendedS3:
    bucketArn: arn:aws:s3:::my-analytics-data-lake
    roleArn: arn:aws:iam::123456789012:role/firehose-glue-s3-role
    prefix: analytics/year=!{timestamp:yyyy}/month=!{timestamp:MM}/day=!{timestamp:dd}/
    compressionFormat: UNCOMPRESSED
    fileExtension: .parquet
    buffering:
      intervalInSeconds: 60
      sizeInMbs: 64
    processing:
      enabled: true
      processors:
        - metadataExtraction:
            query: '{event_type: .event_type}'
            jsonParsingEngine: JQ-1.6
    dynamicPartitioning:
      enabled: true
      retryDurationInSeconds: 300
    dataFormatConversion:
      enabled: true
      inputFormat: OPENX_JSON
      outputFormat: PARQUET
      parquetCompression: SNAPPY
      schema:
        databaseName: analytics_db
        tableName: events
        roleArn: arn:aws:iam::123456789012:role/firehose-glue-access-role
```

### Snowflake Streaming with Secrets Manager Credentials

Streams records into a Snowflake table via Snowpipe Streaming; the key-pair credential lives in Secrets Manager and never appears in the manifest:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKinesisFirehose
metadata:
  name: snowflake-streaming-firehose
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsKinesisFirehose.snowflake-streaming-firehose
spec:
  region: us-east-1
  snowflake:
    accountUrl: https://myaccount.snowflakecomputing.com
    database: ANALYTICS
    schema: PUBLIC
    table: EVENTS
    roleArn: arn:aws:iam::123456789012:role/firehose-snowflake-role
    secretsManager:
      secretArn: arn:aws:secretsmanager:us-east-1:123456789012:secret:snowflake-firehose-keypair-AbCdEf
    snowflakeRole: FIREHOSE_INGEST
    s3BackupMode: FailedDataOnly
    s3Config:
      bucketArn: arn:aws:s3:::my-firehose-backup-bucket
      roleArn: arn:aws:iam::123456789012:role/firehose-s3-backup-role
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `delivery_stream_arn` | `string` | ARN of the Kinesis Data Firehose delivery stream |
| `delivery_stream_name` | `string` | Name of the delivery stream, unique within the AWS account and region |
| `destination_id` | `string` | Identifier of the destination configuration within the stream |
| `version_id` | `string` | Configuration version, incremented by AWS on every update |

## Related Components

- [AwsKinesisStream](/docs/catalog/aws/kinesis-data-stream) — provides a Kinesis Data Stream as the source for the delivery stream
- [AwsMskCluster](/docs/catalog/aws/msk-cluster) — provides an MSK cluster as the Kafka source for the delivery stream
- [AwsS3Bucket](/docs/catalog/aws/s3-bucket) — serves as the delivery destination, backup target, or Redshift staging area
- [AwsOpenSearchDomain](/docs/catalog/aws/opensearch-domain) — serves as the indexing destination for log and search workloads
- [AwsIamRole](/docs/catalog/aws/iam-role) — provides the permissions Firehose needs for source, destination, and transformation access
- [AwsLambda](/docs/catalog/aws/lambda) — provides the Lambda function for record transformation before delivery
