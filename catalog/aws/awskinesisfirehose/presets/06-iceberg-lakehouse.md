# Iceberg Lakehouse Preset

Streams records directly into a Glue-cataloged Apache Iceberg table with upsert semantics — Firehose writes Iceberg snapshots itself, so there is no Spark job, no custom writer, and no staged-file compaction pipeline to operate.

## When to Use

- **Streaming lakehouse ingestion** — Land events as queryable Iceberg data (Athena, EMR, Redshift Spectrum) with ACID snapshots
- **Change-data-capture** — `uniqueKeys` turns matching records into UPDATEs, keeping the table a current-state view instead of an append log
- **Replacing Parquet-plus-compaction pipelines** — Iceberg handles small-file compaction and schema evolution natively

## Key Configuration

- **Glue catalog binding** — `catalogArn` is the account's Glue Data Catalog (create-time immutable)
- **Upserts via unique keys** — `order_id` identifies rows; matching records update in place. Set `appendOnly: true` instead for pure append workloads
- **Per-table error prefix** — failed `orders` records land under `errors/orders/`
- **Multi-table routing** — list multiple `destinationTables` and add a `metadata_extraction` processor producing per-record routing metadata
- **S3 backup of failed data** — undeliverable records land in the backup bucket for replay

## Prerequisites

| Resource | Description |
|----------|-------------|
| **Iceberg table** | Created in the Glue Data Catalog (e.g., via Athena `CREATE TABLE ... TBLPROPERTIES ('table_type'='ICEBERG')`) with a schema matching the records. |
| **S3 bucket** | Backup bucket for failed records. |
| **IAM role** | Delivery role with `glue:GetTable`/`glue:UpdateTable` on the catalog + S3 read/write on the table's warehouse location and the backup bucket. |

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `123456789012` | Your AWS account ID (in the catalog ARN and role ARNs) |
| `lakehouse` / `orders` | Glue database / Iceberg table |
| `order_id` | The row-identity column for upserts |
| `my-firehose-backup-bucket` | S3 backup bucket |
