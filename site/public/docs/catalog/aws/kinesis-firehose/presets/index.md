---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kinesis Firehose"
type: "preset-list"
componentSlug: "kinesis-firehose"
componentTitle: "Kinesis Firehose"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-s3-data-lake"
    rank: "01"
    title: "S3 Data Lake Preset"
    excerpt: "Minimal Extended S3 destination for general-purpose log and event storage. Uses Direct PUT source with GZIP compression and time-based partitioning."
  - slug: "02-opensearch-log-analytics"
    rank: "02"
    title: "OpenSearch Log Analytics Preset"
    excerpt: "OpenSearch destination for centralized log indexing with near real-time search. Indexes logs directly into OpenSearch with daily rotation and S3 backup for failed documents."
  - slug: "03-http-endpoint-webhook"
    rank: "03"
    title: "HTTP Endpoint Webhook Preset"
    excerpt: "HTTP endpoint destination for third-party log and metrics integrations. This preset targets Datadog's HTTP intake API with GZIP compression and custom attributes."
  - slug: "04-s3-parquet-analytics"
    rank: "04"
    title: "S3 Parquet Analytics Preset"
    excerpt: "Extended S3 destination with Parquet format conversion for data lake analytics. Consumes from a Kinesis stream, converts JSON to Parquet using a Glue catalog schema, and delivers with dynamic..."
  - slug: "05-snowflake-streaming"
    rank: "05"
    title: "Snowflake Streaming Preset"
    excerpt: "Direct streaming into a Snowflake table via Snowpipe Streaming — no intermediate S3 staging, no external Snowpipe configuration, near-real-time inserts. Credentials are sourced from AWS Secrets..."
  - slug: "06-iceberg-lakehouse"
    rank: "06"
    title: "Iceberg Lakehouse Preset"
    excerpt: "Streams records directly into a Glue-cataloged Apache Iceberg table with upsert semantics — Firehose writes Iceberg snapshots itself, so there is no Spark job, no custom writer, and no staged-file..."
---

# Kinesis Firehose Presets

Ready-to-deploy configuration presets for Kinesis Firehose. Each preset is a complete manifest you can copy, customize, and deploy.
