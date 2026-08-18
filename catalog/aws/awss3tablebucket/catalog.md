# AWS S3 Table Bucket

Apache Iceberg without the platform team: a bucket purpose-built for analytics tables, where AWS runs the compaction, snapshot expiry, and orphan cleanup that otherwise consume a data engineer — and Athena, EMR, and any Iceberg engine query the tables in place.

## What Gets Managed

- The table bucket: encryption defaults, unreferenced-file cleanup dials, force-destroy posture, resource policy, and bucket-level replication.
- Its namespaces (logical databases) and their Iceberg tables: create-time schema, table properties, per-table maintenance (compaction, snapshot management), per-table policies and replication.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with S3 Tables permissions.

### AWS Prerequisites

- None for the bucket itself. For KMS encryption: the maintenance service needs a grant on the key (or compaction silently stops). For querying: the analytics-catalog integration (one per account, in the AWS console or Glue).

## After You Deploy

- Point query engines at the bucket via the catalog integration, or manually at each table's `table_warehouse_locations` entry.
- Write data through Iceberg-aware engines (Athena, Spark, PyIceberg) — the bucket stores tables, not loose objects.

## Common Changes

- New dataset: add a table under its namespace — one spec entry.
- Schema change: through the query engine (ALTER TABLE), never the spec — the spec's schema is create-time only.
- Tune maintenance: per-table compaction target size and snapshot retention update in place.
