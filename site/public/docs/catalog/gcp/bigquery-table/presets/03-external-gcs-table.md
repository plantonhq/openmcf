---
title: "Preset: External GCS Table"
description: "Use this preset to query data-lake files in GCS without loading them into BigQuery-managed storage: parquet/CSV/JSON exports, hive-partitioned lake layouts, or data shared with Spark and other..."
type: "preset"
rank: "03"
presetSlug: "03-external-gcs-table"
componentSlug: "bigquery-table"
componentTitle: "BigQuery Table"
provider: "gcp"
icon: "package"
order: 3
---

# Preset: External GCS Table

## When to Use

Use this preset to query data-lake files in GCS without loading them into
BigQuery-managed storage: parquet/CSV/JSON exports, hive-partitioned lake
layouts, or data shared with Spark and other engines.

## What It Creates

- An external table over `gs://` parquet files
- Schema autodetection from the files
- Hive partition-key inference from the object paths (`.../dt=2026-01-01/...`)
- Deletion protection OFF (the table holds no data of its own)

## Configuration

| Field | Value | Notes |
|-------|-------|-------|
| sourceFormat | PARQUET | CSV, NEWLINE_DELIMITED_JSON, AVRO, ORC, ICEBERG, ... also supported |
| autodetect | true | Schema inferred from the files |
| hivePartitioningOptions | AUTO | Partition keys and types inferred from paths |

## Upgrading to BigLake

Add a `connectionId` (a BigQuery connection whose service account can read
the bucket) to upgrade the plain external table to a BigLake table:
fine-grained row/column security and metadata caching (`metadataCacheMode`
+ the table's `maxStaleness`) become available, and readers no longer need
direct GCS access.

## How to Use

1. Replace the `datasetId` ref's `name` with your dataset resource name
2. Replace `tableId` and the `sourceUris` bucket path
3. Match `sourceFormat` (and the options block) to your file format
