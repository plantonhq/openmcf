# Preset: Partitioned Analytics Table

## When to Use

Use this preset for the workhorse analytics pattern: an append-only event or
fact table that grows without bound. Day partitioning plus clustering keeps
query cost proportional to the data actually scanned, and the partition
filter requirement stops accidental full-table scans.

## What It Creates

- A native BigQuery table with an explicit JSON schema
- Day partitions on the event timestamp, expiring after 365 days
- Clustering on the two hottest filter columns
- A mandatory partition filter on every query (cost guard)
- Deletion protection ON (the default) — destroy fails until disabled

## Configuration

| Field | Value | Notes |
|-------|-------|-------|
| Partitioning | DAY on `event_time` | Queries filtering on the column scan only matching days |
| Partition expiration | 365 days | Old partitions drop automatically |
| Clustering | `customer_id`, `event_type` | Order matters: lead with the most selective filter |
| requirePartitionFilter | true | Rejects queries without a partition predicate |

## How to Use

1. Replace the `datasetId` ref's `name` with your GcpBigQueryDataset resource name
2. Replace `tableId` and the `schema` JSON with your columns
3. Adjust partitioning granularity (HOUR for very high-volume streams) and
   the clustering columns to your query patterns

## Schema Evolution

Columns are add-only in place: appending a new NULLABLE column updates the
table; removing or retyping a column recreates it (blocked by deletion
protection until you opt out).
