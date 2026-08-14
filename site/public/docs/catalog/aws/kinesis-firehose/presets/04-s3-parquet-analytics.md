---
title: "S3 Parquet Analytics Preset"
description: "Extended S3 destination with Parquet format conversion for data lake analytics. Consumes from a Kinesis stream, converts JSON to Parquet using a Glue catalog schema, and delivers with dynamic..."
type: "preset"
rank: "04"
presetSlug: "04-s3-parquet-analytics"
componentSlug: "kinesis-firehose"
componentTitle: "Kinesis Firehose"
provider: "aws"
icon: "package"
order: 4
---

# S3 Parquet Analytics Preset

Extended S3 destination with Parquet format conversion for data lake analytics. Consumes from a Kinesis stream, converts JSON to Parquet using a Glue catalog schema, and delivers with dynamic partitioning.

## When to Use

- **Data lake analytics** — Optimized for Athena, Spark, Presto, and Redshift Spectrum
- **Cost and performance** — Parquet provides 60–90% compression and 10–100x faster analytical queries
- **Streaming ingestion** — Kinesis stream as source for high-throughput event pipelines
- **Schema enforcement** — Glue Data Catalog defines and validates record structure

## Key Configuration

- **Kinesis stream source** — `valueFrom` references an AwsKinesisStream; Firehose consumes with automatic checkpointing
- **Data format conversion** — the `openXJson` deserializer arm reads the JSON, the `parquet` serializer arm writes columnar output; exactly one arm on each side selects the format
- **Deserializer tuning** — `columnToJsonKeyMappings` reads the JSON key `timestamp` into the Glue column `ts` (avoiding the Hive reserved word), and dotted JSON keys convert to underscores so they land in valid Glue columns
- **Parquet tuning** — GZIP compression (harder than the SNAPPY default; right when storage/scan cost outweighs write throughput), 128 MiB row groups, dictionary encoding for low-cardinality columns, and writer version V2 (confirm downstream reader support)
- **Glue catalog schema** — `analytics_db.events` table defines the Parquet schema
- **Metadata extraction + dynamic partitioning** — a JQ processor extracts `event_type` from each record; the S3 prefix references it as `!{partitionKeyFromQuery:event_type}` so data lands in query-prunable partitions
- **File extension** — `.parquet` for delivered objects

## Prerequisites

| Resource | Description |
|----------|-------------|
| **Kinesis Data Stream** | Source stream with events in JSON format. Reference via `valueFrom` from AwsKinesisStream. |
| **Glue Data Catalog** | Database and table with schema matching your JSON records. Create table with appropriate column types. |
| **S3 bucket** | Target bucket for Parquet objects. Use lifecycle policies for cost optimization. |
| **IAM roles** | Roles for: (1) Firehose to read from Kinesis, (2) Firehose to write to S3, (3) Firehose to access Glue catalog. |

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `my-events-stream` | Name of your AwsKinesisStream resource |
| `my-analytics-data-lake` | S3 bucket for Parquet output |
| `analytics_db` | Glue database name |
| `events` | Glue table name |
| `123456789012` | Your AWS account ID |
| `firehose-glue-s3-role` | IAM role for S3 and Glue access |
| `firehose-glue-access-role` | IAM role for Glue GetTable/GetTableVersions |
| `firehose-kinesis-consumer-role` | IAM role for Kinesis read (when using valueFrom) |

## Schema Requirements

The Glue table schema must match your JSON record structure. Ensure:

- Column names and types align with incoming JSON
- Use compatible types (e.g., `string`, `bigint`, `double`, `timestamp`, `struct`, `array`)
- Table exists before creating the delivery stream
