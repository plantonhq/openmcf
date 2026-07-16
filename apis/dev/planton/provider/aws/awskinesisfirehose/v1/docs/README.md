# AwsKinesisFirehose — Architecture Reference

## Overview

Amazon Kinesis Data Firehose is a fully managed service for loading streaming data into storage and analytics destinations. Unlike Kinesis Data Streams (which requires you to write and operate consumers), Firehose handles the entire delivery pipeline: buffering, transformation, format conversion, compression, and retry — with zero consumer code.

Firehose is the simplest path from streaming data to S3, OpenSearch (domains and serverless collections), HTTP endpoints, Redshift, Splunk, Snowflake, and Apache Iceberg tables.

## How Firehose Works

### Data Flow

```
Source → Buffer → [Process] → [Convert] → [Compress] → Destination
                                                  ↓ (failures)
                                             S3 Backup
```

1. **Ingest** — Data enters the delivery stream from a source (Direct PUT, Kinesis Data Stream, or MSK topic)
2. **Buffer** — Records accumulate in an in-memory buffer until either the size or time threshold is reached
3. **Process** (optional) — An ordered processor pipeline runs: decompression, de-aggregation, CloudWatch-Logs unwrapping, Lambda transformation, JQ metadata extraction, delimiter appending
4. **Convert** (optional) — JSON records are converted to columnar format (Parquet/ORC) using a Glue schema
5. **Compress** (optional) — Data is compressed before writing (GZIP, Snappy, etc.)
6. **Deliver** — The batch is written to the destination
7. **Backup** — Failed records (or all records, if configured) are written to S3

Each delivery stream has exactly **one source** and exactly **one destination**. Both are immutable after creation (ForceNew).

## Source Types

### Direct PUT (Default)

Applications call the `PutRecord` or `PutRecordBatch` API directly on the delivery stream. This is the simplest source — no intermediate streaming infrastructure required.

**Characteristics:**

- Up to 1,000 records/s or 1 MB/s per delivery stream (soft limit, can be increased)
- Records up to 1 MiB each
- Server-side encryption (SSE) can be enabled on the delivery stream buffer
- No ordering guarantees — records may arrive out of order
- No replay — once delivered, records cannot be re-read from the source

**When to use:**

- Application produces data at moderate throughput (<1 MB/s)
- No need for ordered processing or replay
- Want the simplest possible architecture (no Kinesis stream to manage)
- Log forwarding from CloudWatch Logs, IoT Core, or other AWS services

### Kinesis Data Stream Source

Firehose reads from an existing Kinesis Data Stream, acting as a managed consumer with automatic checkpointing and retry.

**Characteristics:**

- Firehose creates an internal consumer and reads all shards
- Inherits the stream's ordering guarantees (per-shard, per-partition-key)
- No SSE on the delivery stream — encryption is handled by the source stream
- Stream must exist before the delivery stream is created
- Source configuration is entirely ForceNew

**When to use:**

- Multiple consumers need the same data (Firehose + Lambda + custom app)
- Need replay capability (Kinesis retains data for 24h–365 days)
- Need ordering guarantees (per partition key)
- High throughput (>1 MB/s) with auto-scaling (ON_DEMAND stream)
- Already have a Kinesis stream in your architecture

### Amazon MSK Source

Firehose reads a Kafka topic on an Amazon MSK cluster (provisioned or serverless), acting as a managed Kafka consumer — no consumer application, no consumer group management.

**Characteristics:**

- The cluster must have IAM access control enabled; Firehose authenticates with the configured IAM role
- Connectivity is `PRIVATE` (through the cluster's in-VPC brokers — the common case) or `PUBLIC`
- `read_from_timestamp` rewinds the topic to a point in time at creation; otherwise Firehose starts from the latest offset
- No SSE on the delivery stream — the cluster handles encryption
- Source configuration is entirely ForceNew

**When to use:**

- Kafka topics that need to land in S3/OpenSearch/Snowflake/Iceberg without operating a consumer
- Offloading topic data for analytics while Kafka remains the system of record

### Source Comparison

| Feature | Direct PUT | Kinesis Stream Source | MSK Source |
|---------|------------|---------------------|------------|
| Setup complexity | Lowest | Requires existing stream | Requires existing cluster + topic |
| Throughput | 1,000 rec/s or 1 MB/s (soft limit) | Stream capacity (unlimited with ON_DEMAND) | Topic/cluster capacity |
| Ordering | None | Per-shard (partition key) | Per-partition |
| Replay | No | Yes (stream retention) | Point-in-time start (`read_from_timestamp`) |
| Multiple consumers | No | Yes (stream supports many readers) | Yes (Kafka consumer groups) |
| SSE | On delivery stream buffer | On source stream | On source cluster |
| Cost | Firehose per-GB only | Stream cost + Firehose per-GB | Cluster cost + Firehose per-GB |

## Destination Types

### Destination Comparison

| Feature | Extended S3 | OpenSearch | OpenSearch Serverless | HTTP Endpoint | Redshift | Splunk | Snowflake | Iceberg |
|---------|-------------|-----------|----------------------|---------------|----------|--------|-----------|---------|
| **Use case** | Data lake, archive | Log analytics, search | Serverless log analytics | Third-party SaaS, custom APIs | Data warehouse | Splunk observability | Warehouse streaming | Streaming lakehouse |
| **Delivery path** | Direct to S3 | Direct to index | Direct to collection index | HTTPS POST | S3 staging → COPY | HEC POST + ack | Snowpipe Streaming insert | Iceberg snapshot commit |
| **Format conversion** | Parquet/ORC via Glue | No | No | No | No | No | No | Iceberg-native |
| **Dynamic partitioning** | Yes | No | No | No | No | No | No | Table routing via processors |
| **Processor pipeline** | Yes | Yes | Yes | Yes | Yes | Yes | Yes | Yes |
| **Secrets Manager credentials** | N/A (IAM) | N/A (IAM) | N/A (IAM) | Yes (`api_key`) | Yes (username/password) | Yes (HEC token) | Yes (key pair) | N/A (IAM) |
| **S3 backup** | Optional (source records) | Required | Required | Required | Required (staging) + optional | Required | Required | Required |
| **VPC delivery** | N/A | Yes | Yes | No | No (public JDBC) | No | PrivateLink VPCE | No |
| **Buffer caps** | 0–900s / 1–128 MiB | 0–900s / 1–100 MiB | 0–900s / 1–100 MiB | 0–900s / 1–100 MiB | S3 staging | 0–60s / 1–5 MiB | 0–900s / 1–128 MiB (default 0s/1 MiB) | 0–900s / 1–128 MiB |

### Extended S3

The most feature-rich destination and the most common. Data lands in S3 as objects, optionally compressed, converted to columnar format, and dynamically partitioned by record fields.

**Best for:** Data lakes, log archives, analytics pipelines (Athena, Spark, Presto), long-term storage, compliance archives.

### OpenSearch

Indexes records directly into an Amazon OpenSearch Service domain. Supports index rotation (hourly, daily, weekly, monthly), document-ID control, and VPC delivery for private clusters. Failed documents are always backed up to S3.

**Best for:** Log analytics, full-text search, real-time dashboards (OpenSearch Dashboards/Kibana).

### OpenSearch Serverless

Indexes records into an OpenSearch Serverless collection endpoint. No domains, no rotation — the collection scales automatically and the index is a fixed name admitted by the collection's data access policy.

**Best for:** Log analytics without capacity management; teams already standardized on OpenSearch Serverless.

### HTTP Endpoint

Delivers to any HTTPS endpoint that accepts POST requests and returns HTTP 200 on success. The endpoint receives JSON arrays of records. Authentication is via an access key in the `X-Amz-Firehose-Access-Key` header — inline or sourced from Secrets Manager.

**Best for:** Third-party integrations (Datadog, New Relic, Sumo Logic, Honeycomb), custom APIs, webhook-based pipelines.

### Redshift

A two-stage destination: Firehose writes data to an S3 staging bucket, then issues a Redshift `COPY` command to bulk-load the data. This is the standard Redshift ingestion pattern for streaming data. Credentials are inline username/password or a Secrets Manager secret.

**Best for:** Data warehouse loading, business intelligence, reporting pipelines.

### Splunk

Posts events to a Splunk HTTP Event Collector (HEC) endpoint — Splunk Cloud, Splunk Enterprise, or Splunk-managed AWS — and waits for indexer acknowledgment before considering delivery complete. Enforces the tightest buffering of any destination (max 60 seconds / 5 MiB) to keep event latency low. The HEC token is inline or a Secrets Manager secret.

**Best for:** Organizations standardized on Splunk for observability and SIEM.

### Snowflake

Streams records directly into a Snowflake table via Snowpipe Streaming — no intermediate S3 staging and no external Snowpipe configuration. Authenticates with key-pair credentials (inline or Secrets Manager); supports PrivateLink, a dedicated ingestion role, and JSON-to-column or VARIANT loading modes. Defaults to the fastest buffering of any destination (0 seconds / 1 MiB).

**Best for:** Near-real-time Snowflake ingestion replacing batch COPY pipelines.

### Iceberg

Commits records directly into Apache Iceberg tables managed by the AWS Glue Data Catalog — Firehose writes Iceberg snapshots itself, with optional per-record routing across multiple tables and update/delete semantics via unique keys.

**Best for:** Streaming lakehouse ingestion and change-data-capture into Iceberg without a Spark job or custom writer.

## Buffering Model

Firehose buffers incoming records and flushes to the destination when **either** threshold is reached — whichever comes first:

- **Buffer interval** — Time since the last flush (0–900 seconds, default 300)
- **Buffer size** — Accumulated data size (1–128 MiB, default 5 MiB)

### How Buffering Works

```
Records arrive → Buffer fills
                   ├── Size threshold reached? → Flush
                   └── Time threshold reached? → Flush
```

**Example:** With `intervalInSeconds: 300` and `sizeInMbs: 5`:

- If 5 MiB accumulates in 30 seconds → flushes after 30 seconds (size trigger)
- If only 100 KiB arrives in 300 seconds → flushes after 300 seconds (time trigger)

### Tuning Guidelines

| Goal | Interval | Size | Effect |
|------|----------|------|--------|
| Low latency | 60s | 1 MiB | More frequent, smaller deliveries |
| Optimal for analytics | 900s | 128 MiB | Fewer, larger files (better for Athena/Spark) |
| Balanced | 300s | 5 MiB | Default — good for most use cases |
| Parquet conversion | 900s | 128 MiB | Larger batches produce better Parquet files |

**Destination-specific limits:**

- Extended S3, Snowflake, Iceberg: 1–128 MiB
- OpenSearch, OpenSearch Serverless, HTTP Endpoint: 1–100 MiB
- Splunk: 0–60 seconds / 1–5 MiB
- Redshift: Uses S3 staging buffering

## Record-Transformation Pipeline

The processing pipeline is an **ordered list of processors** — each processor's output feeds the next. Six processor types are available:

| Processor | Purpose | Typical position |
|-----------|---------|------------------|
| **Decompression** | Decompress GZIP payloads (CloudWatch Logs subscriptions arrive compressed) | First |
| **CloudWatch log processing** | Unwrap CloudWatch Logs subscription envelopes into individual log events | After decompression |
| **Record de-aggregation** | Split KPL-aggregated or delimited payloads into individual records (Extended S3 only) | Before per-record processing |
| **Lambda** | Arbitrary per-batch transformation (enrich/filter/reshape) | Middle |
| **Metadata extraction** | Extract partition keys from JSON with a JQ expression (drives dynamic partitioning) | After transformation |
| **Append delimiter** | Newline-delimit records (JSON lines) for query engines (Extended S3 only) | Last |

AWS restricts de-aggregation and delimiter appending to S3 delivery — creation fails with them on any other destination.

**Common pipelines:**

- `lambda` — classic transformation
- `metadata_extraction` — dynamic partitioning by record fields
- `decompression → cloudwatch_log_processing → append_delimiter` — CloudWatch Logs to clean JSON-lines in S3
- `record_deaggregation → metadata_extraction` — KPL producers + partitioning

## Lambda Transformation

### How It Works

1. Firehose accumulates records in a processing buffer (1–3 MiB, 60–900s)
2. When the buffer threshold is reached, Firehose invokes the Lambda function with a batch of records
3. Lambda processes each record and returns a result with a **status code** per record
4. Firehose routes records based on their status code

### Lambda Input/Output

**Input event structure:**

```json
{
  "invocationId": "...",
  "deliveryStreamArn": "arn:aws:firehose:...",
  "region": "us-east-1",
  "records": [
    {
      "recordId": "...",
      "approximateArrivalTimestamp": 1234567890,
      "data": "<base64-encoded record>"
    }
  ]
}
```

**Output structure:**

```json
{
  "records": [
    {
      "recordId": "...",
      "result": "Ok",
      "data": "<base64-encoded transformed record>"
    }
  ]
}
```

### Record Status Codes

| Status | Behavior |
|--------|----------|
| `Ok` | Record is delivered to the destination |
| `Dropped` | Record is intentionally discarded (counted as successfully processed) |
| `ProcessingFailed` | Record failed transformation — written to the error output prefix |

### Retry Behavior

- If Lambda returns an error (exception, timeout, throttle), Firehose retries the entire batch
- Retries up to `numberOfRetries` times (default 3, max 300)
- After all retries are exhausted, the batch is written to the error output prefix in S3
- Lambda timeout should be ≤ 5 minutes (Firehose's invocation timeout is 5 minutes)

### Dynamic Partitioning with Lambda

Lambda can extract partition keys from records by adding metadata to the response:

```json
{
  "recordId": "...",
  "result": "Ok",
  "data": "<base64>",
  "metadata": {
    "partitionKeys": {
      "customer_id": "cust-123",
      "region": "us-east-1"
    }
  }
}
```

These keys are referenced in the S3 prefix: `data/customer=!{partitionKeyFromLambda:customer_id}/`

## Data Format Conversion

### JSON → Columnar Conversion

Firehose can convert incoming JSON records to Apache Parquet or Apache ORC format using an AWS Glue Data Catalog schema. This dramatically improves query performance and reduces storage costs for analytics workloads.

### How It Works

1. Firehose receives JSON records
2. Deserializes using the configured input format (OpenX JSON SerDe or Hive JSON SerDe)
3. Maps fields to the Glue table schema
4. Serializes to the output columnar format (Parquet or ORC)
5. Applies columnar-native compression (SNAPPY, GZIP, ZLIB)
6. Writes to S3

### Performance Impact

| Metric | JSON (GZIP) | Parquet (SNAPPY) | Improvement |
|--------|-------------|------------------|-------------|
| Storage size | 1x | 0.1–0.4x | 60–90% smaller |
| Athena scan cost | Full file scan | Columnar pruning | 10–100x cheaper |
| Query latency | Seconds to minutes | Sub-second to seconds | 10–100x faster |

### Parquet vs ORC

| Feature | Parquet | ORC |
|---------|---------|-----|
| Best for | Athena, Spark, Presto, general analytics | Hive workloads |
| Compression | SNAPPY (default), GZIP, UNCOMPRESSED | SNAPPY (default), ZLIB, NONE |
| Predicate pushdown | Yes | Yes |
| ACID support | No | Yes (Hive ACID) |
| Recommendation | **Use Parquet** for most analytics workloads | Use ORC for Hive-centric pipelines |

### Prerequisites

- **Glue Data Catalog** database and table with the schema definition
- Table must define column names, types, and (optionally) partition keys
- Firehose IAM role must have `glue:GetTable` and `glue:GetTableVersions` permissions
- Schema changes require updating the Glue table (Firehose reads the latest version by default)

## Dynamic Partitioning

### How Partition Keys Work

Dynamic partitioning extracts key-value pairs from each record and uses them to construct unique S3 prefixes, creating a partitioned directory layout for efficient querying.

**Two extraction methods:**

1. **From JQ expressions** — a `metadata_extraction` processor applies a JQ query to each JSON record; the resulting keys are referenced as `!{partitionKeyFromQuery:key}` in the prefix
2. **From Lambda metadata** — Partition keys returned by a Lambda transformation function, referenced as `!{partitionKeyFromLambda:key}`

### Prefix Expressions

The S3 prefix uses `!{partitionKeyFromQuery:key}` or `!{partitionKeyFromLambda:key}` syntax:

```
data/region=!{partitionKeyFromLambda:region}/customer=!{partitionKeyFromLambda:customer_id}/year=!{timestamp:yyyy}/month=!{timestamp:MM}/
```

This produces a directory layout like:

```
data/
  region=us-east-1/
    customer=cust-123/
      year=2026/
        month=02/
          file001.parquet
          file002.parquet
    customer=cust-456/
      year=2026/
        month=02/
          file003.parquet
  region=eu-west-1/
    customer=cust-789/
      ...
```

### Important Considerations

- Dynamic partitioning is **ForceNew** — it cannot be enabled or disabled after the delivery stream is created
- Each unique partition key combination creates a separate S3 prefix (and separate buffer)
- High cardinality partition keys (e.g., user_id with millions of values) create many small files — use buffering hints to mitigate
- Retry duration controls how long Firehose retries when a partition write fails (0–7200s, default 300)

## Encryption Model

### Server-Side Encryption (SSE) for Direct PUT

When using Direct PUT as the source, Firehose can encrypt data at rest in its internal buffer:

| Configuration | Encryption | Cost |
|---------------|-----------|------|
| `sseEnabled: false` | No encryption (plaintext in buffer) | No additional cost |
| `sseEnabled: true`, no KMS key | AWS-owned CMK (Firehose manages the key) | No additional cost |
| `sseEnabled: true`, with `sseKmsKeyArn` | Customer-managed CMK | KMS API charges per encryption/decryption |

### Encryption with a Stream Source

When using a Kinesis Data Stream or MSK cluster as the source, **do not enable SSE on the delivery stream**. The source handles encryption:

- The source stream/cluster has its own KMS encryption configuration
- Firehose reads already-encrypted records and the source's KMS key handles decryption
- Enabling SSE on the delivery stream would be redundant and is rejected by the API

### Encryption at S3 Destination

S3 encryption is configured separately in the destination's `kmsKeyArn` field:

- When absent, S3 uses its default encryption (SSE-S3 or bucket default)
- When set, S3 uses SSE-KMS with the specified customer-managed key
- This is independent of the delivery stream's SSE setting

## S3 Backup Behavior

Every non-S3 destination requires S3 backup for error handling. The behavior varies by destination:

| Destination | S3 Role | Backup Modes | Default |
|-------------|---------|-------------|---------|
| **Extended S3** | Primary destination | `Disabled` / `Enabled` (source record backup) | `Disabled` |
| **OpenSearch** | Backup for failed/all documents | `FailedDocumentsOnly` / `AllDocuments` (ForceNew) | `FailedDocumentsOnly` |
| **OpenSearch Serverless** | Backup for failed/all documents | `FailedDocumentsOnly` / `AllDocuments` (ForceNew) | `FailedDocumentsOnly` |
| **HTTP Endpoint** | Backup for failed/all records | `FailedDataOnly` / `AllData` | `FailedDataOnly` |
| **Redshift** | Staging for COPY + optional source backup | `Disabled` / `Enabled` (source record backup) | `Disabled` |
| **Splunk** | Backup for failed/all events | `FailedEventsOnly` / `AllEvents` | `FailedEventsOnly` |
| **Snowflake** | Backup for failed/all records | `FailedDataOnly` / `AllData` | `FailedDataOnly` |
| **Iceberg** | Backup for failed/all records | `FailedDataOnly` / `AllData` | `FailedDataOnly` |

### Extended S3 Backup

For Extended S3 destinations, "backup" means keeping a copy of the **original, pre-transformation** records. This is separate from the primary S3 delivery:

- **Primary delivery** → transformed/converted records (e.g., Parquet with enrichments)
- **Source backup** → original JSON records as received (for auditing/reprocessing)

### Redshift S3 Staging

The Redshift destination uses S3 as an **intermediate staging area** (not backup). Firehose writes data to S3, then issues a `COPY` command. This S3 data is the primary delivery path. Optionally, you can also enable source record backup to a separate S3 location.

## Cost Model

Firehose pricing is straightforward — pay per GB of data ingested. No upfront costs, no provisioning, no idle charges.

### Base Pricing (US East)

| Component | Price |
|-----------|-------|
| Data ingested (first 500 TB/month) | ~$0.029/GB |
| Data ingested (next 1.5 PB/month) | ~$0.025/GB |
| Data ingested (over 2 PB/month) | ~$0.020/GB |
| Format conversion (Parquet/ORC) | ~$0.018/GB |
| Dynamic partitioning | ~$0.020/GB |
| VPC delivery (per hour per AZ) | ~$0.01/hour/AZ |

### What Counts as "Ingested"

- Firehose rounds each record up to the nearest **5 KB** for billing purposes
- A 100-byte record is billed as 5 KB
- A 6 KB record is billed as 10 KB
- Records under 5 KB should be batched client-side for cost efficiency

### Additional Costs

- **Lambda transformation**: Standard Lambda pricing (invocations + duration)
- **S3 storage**: Standard S3 pricing for delivered and backup objects
- **KMS encryption**: Per-API-call KMS charges when using customer-managed keys
- **CloudWatch Logs**: Standard CloudWatch Logs pricing for error logs
- **Glue Data Catalog**: Free for the first million objects; standard pricing after

### Cost Optimization Tips

- Batch records client-side to approach 1 MiB per record (minimize 5 KB rounding overhead)
- Use GZIP compression to reduce delivered data volume (destination storage cost)
- Set appropriate buffer sizes — larger buffers reduce the number of S3 PUT operations
- Use format conversion (Parquet) for analytics workloads — storage savings often outweigh the conversion cost

## Security

### IAM Roles

Firehose requires IAM roles for every interaction with other AWS services. A single role can cover multiple permissions, or you can use separate roles for fine-grained control:

| Permission | Actions | Used By |
|------------|---------|---------|
| S3 write | `s3:PutObject`, `s3:AbortMultipartUpload`, `s3:GetBucketLocation`, `s3:ListBucket` | All destinations |
| Kinesis read | `kinesis:GetRecords`, `kinesis:GetShardIterator`, `kinesis:DescribeStream`, `kinesis:ListShards` | Kinesis source |
| MSK read | `kafka:GetBootstrapBrokers`, `kafka:DescribeCluster(V2)`, `kafka-cluster:Connect`, `kafka-cluster:DescribeTopic`, `kafka-cluster:ReadData`, `kafka-cluster:DescribeGroup` | MSK source |
| Lambda invoke | `lambda:InvokeFunction`, `lambda:GetFunctionConfiguration` | Lambda processing |
| OpenSearch write | `es:ESHttpPut`, `es:ESHttpGet` | OpenSearch destination |
| OpenSearch Serverless write | `aoss:APIAccessAll` + a data-access-policy grant | OpenSearch Serverless destination |
| KMS encrypt/decrypt | `kms:Encrypt`, `kms:Decrypt`, `kms:GenerateDataKey` | SSE, S3 SSE-KMS |
| Glue catalog | `glue:GetTable`, `glue:GetTableVersions` (+ `glue:UpdateTable` for Iceberg) | Format conversion, Iceberg |
| Secrets Manager read | `secretsmanager:GetSecretValue` | Secrets Manager credentials |
| VPC ENI management | `ec2:CreateNetworkInterface`, `ec2:DescribeNetworkInterfaces`, `ec2:DeleteNetworkInterface` | VPC delivery |
| CloudWatch Logs | `logs:PutLogEvents` | Error logging |

### Destination Credentials via Secrets Manager

The Redshift, Splunk, HTTP endpoint, and Snowflake destinations authenticate with a credential (password, HEC token, API key, or key pair). Each supports sourcing that credential from an AWS Secrets Manager secret instead of embedding it in the resource configuration — the recommended production mode:

- The credential never appears in manifests or IaC state
- Rotating the secret in Secrets Manager takes effect without a delivery-stream update
- Enabling/disabling Secrets Manager authentication replaces the delivery stream (create-time decision)
- The expected secret shape is destination-specific: `{"username","password"}` (Redshift), `{"hec_token"}` (Splunk), `{"api_key"}` (HTTP), `{"user","private_key","key_passphrase"}` (Snowflake)

### VPC Delivery

For OpenSearch domains deployed in a VPC, Firehose creates Elastic Network Interfaces (ENIs) in the specified subnets:

- Provide subnets in multiple AZs for high availability
- Security groups must allow outbound HTTPS (port 443) to the OpenSearch domain
- The VPC configuration is **ForceNew** — changing subnets or security groups replaces the delivery stream
- VPC delivery adds ~$0.01/hour per AZ to the cost

### Encryption at Rest

Three layers of encryption are available:

1. **Delivery stream buffer** (SSE) — Encrypts data in Firehose's internal buffer (Direct PUT only)
2. **S3 destination** (SSE-KMS) — Encrypts objects written to S3
3. **Source stream** (KMS) — Encrypts data in the Kinesis Data Stream (Kinesis source only)

All three are independent and can use different KMS keys.

## Limits and Quotas

### Delivery Stream Limits

| Limit | Value | Adjustable |
|-------|-------|-----------|
| Delivery streams per region | 500 | Yes (request increase) |
| Record size | 1 MiB maximum | No |
| `PutRecordBatch` batch size | 500 records or 4 MiB | No |
| `PutRecord` throughput (Direct PUT) | 1,000 records/s or 1 MB/s | Yes |
| `PutRecordBatch` throughput (Direct PUT) | 4,000 records/s or 4 MB/s | Yes |
| Lambda processing buffer | 1–3 MiB | No |
| Lambda processing timeout | 5 minutes | No |
| Lambda retries | 0–300 | No |
| Buffer interval | 0–900 seconds | No |
| Buffer size | 1–128 MiB (destination-dependent) | No |
| Dynamic partitioning active partitions | 500 per delivery stream | Yes |
| Retry duration (delivery) | 0–7200 seconds | No |

### Destination-Specific Limits

| Destination | Max Buffer Size | Notes |
|-------------|----------------|-------|
| Extended S3 | 128 MiB | — |
| OpenSearch / OpenSearch Serverless | 100 MiB | — |
| HTTP Endpoint | 100 MiB | Endpoint must respond within 3 minutes |
| Redshift | N/A | Uses S3 staging; COPY command has separate limits |
| Splunk | 5 MiB (60s max interval) | Tightest buffering — HEC latency budget |
| Snowflake | 128 MiB (default 1 MiB / 0s) | Snowpipe Streaming is designed for near-real-time |
| Iceberg | 128 MiB | Larger buffers produce healthier Iceberg snapshots |

## Firehose vs Kinesis Data Streams vs SQS/SNS

### When to Use Each

| Service | Model | Best For |
|---------|-------|----------|
| **Firehose** | Managed ETL pipeline | Zero-code delivery to the eight supported destinations |
| **Kinesis Data Streams** | Streaming log (pull) | Real-time processing, multiple consumers, replay |
| **SQS** | Message queue (pull) | Task distribution, decoupling, exactly-once processing |
| **SNS** | Pub/sub (push) | Fan-out notifications, multi-subscriber broadcasting |

### Detailed Comparison

| Feature | Firehose | Kinesis Streams | SQS | SNS |
|---------|----------|----------------|-----|-----|
| Consumer code | None | You write it | You write it | You configure subscribers |
| Destinations | S3, OpenSearch (+Serverless), HTTP, Redshift, Splunk, Snowflake, Iceberg | Anything (custom code) | Anything (custom code) | Lambda, SQS, HTTP, email, SMS |
| Ordering | None (or per-shard via Kinesis source) | Per-shard | FIFO or best-effort | None |
| Replay | No | Yes (retention-based) | No (once consumed) | No |
| Throughput | Auto-scaling | Per-shard or ON_DEMAND | Auto-scaling | Auto-scaling |
| Latency | Seconds to minutes (buffer-dependent) | ~200ms (poll) / ~70ms (EFO) | ~1ms | ~1ms |
| Retention | None (delivers immediately) | 24h–365 days | 1min–14 days | None |
| Transformation | Lambda (built-in) | Custom consumer code | Custom consumer code | None |
| Format conversion | Parquet/ORC (built-in) | Custom consumer code | Custom consumer code | None |
| Cost model | Per-GB ingested | Per-shard-hour or per-GB | Per-request | Per-publish + per-delivery |

### Common Architectures

**Simple data lake:**
```
Application → Firehose (Direct PUT) → S3
```

**Real-time + archive:**
```
Application → Kinesis Stream → Firehose → S3 (archive)
                             → Lambda → DynamoDB (real-time)
                             → Custom app (analytics)
```

**Multi-destination fan-out:**
```
Application → Kinesis Stream → Firehose #1 → S3 (data lake)
                             → Firehose #2 → OpenSearch (logs)
                             → Firehose #3 → Datadog (monitoring)
```

**Warehouse pipeline:**
```
Application → Firehose (Direct PUT) → Redshift (via S3 staging)
```

## Operational Best Practices

### Monitoring

Key CloudWatch metrics to monitor:

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `DeliveryToS3.Success` | Percentage of successful S3 deliveries | < 99% |
| `DeliveryToS3.Records` | Number of records delivered | Anomaly detection |
| `IncomingBytes` | Data volume entering the stream | Capacity planning |
| `IncomingRecords` | Record count entering the stream | Throughput analysis |
| `ThrottledRecords` | Records rejected due to throttling | > 0 sustained |
| `FailedConversion.Bytes` | Data that failed format conversion | > 0 |
| `ExecuteProcessing.Duration` | Lambda processing time | Approaching 5min timeout |

### Error Handling

- Always configure CloudWatch logging for production delivery streams
- Monitor the error output prefix in S3 — records here indicate transformation or delivery failures
- For OpenSearch destinations, check the S3 backup bucket for rejected documents (index mapping conflicts, field type mismatches)
- For HTTP endpoints, check the S3 backup for non-2xx responses from the endpoint

### Naming Conventions

- Delivery stream names are unique per AWS account per region
- Names are **ForceNew** — choose carefully at creation time
- 1–64 characters, allowed: letters, digits, hyphens, underscores, periods
- Convention: `{service}-{purpose}-{env}` (e.g., `clickstream-s3-prod`)
