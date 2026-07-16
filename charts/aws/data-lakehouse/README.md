# AWS Data Lakehouse

A queryable data lake in one deploy: an S3 lake bucket whose storage cost
follows the age of the data, a Glue Data Catalog database that turns objects
into tables, and an Athena workgroup with a per-query spend guard — plus an
optional Firehose delivery stream that lands JSON events GZIP-compressed and
date-partitioned, ready to query the moment they arrive.

This is the storage-and-query foundation for analytics: event history,
product usage, service logs promoted to structured data, ML training sets.
Producers write to S3 (directly or through the ingestion stream), analysts
and jobs query through Athena with SQL, and every engine added later (Spark,
Redshift Spectrum) resolves the same tables through the same Glue catalog.

## Architecture

```
              producers (PutRecord / PutRecordBatch)
                              |
                     [AwsKinesisFirehose]         (ingestion_enabled)
                        Direct PUT, GZIP,
                     date-partitioned prefix
                              |  assumes [AwsIamRole] scoped to the lake
                              v
   [AwsS3Bucket] lake  ── s3://<lake>/data/year=/month=/day=/
     tiering: STANDARD -> STANDARD_IA -> GLACIER_IR (still queryable)
                              ^
                              | LOCATION / locationUri
                  [AwsGlueCatalogDatabase]
                              ^
                              | FROM <database>.<table>
                   [AwsAthenaWorkgroup]
              scan ceiling + enforced result storage
                              |
              results: s3://<lake>-query-results/   (default)
                       AWS-managed, 24h retention   (managed_results_enabled)
```

## Included Cloud Resources

| Resource | Kind | Purpose |
|---|---|---|
| Lake bucket | `AwsS3Bucket` | The data itself: encrypted, private, lifecycle-tiered, never expired by this chart. |
| Query-results bucket | `AwsS3Bucket` | Athena results with a 30-day expiry. Rendered only while `managed_results_enabled` is off. |
| Catalog database | `AwsGlueCatalogDatabase` | The table namespace every query engine resolves; defaults new tables under `s3://<lake>/data/`. |
| Workgroup | `AwsAthenaWorkgroup` | The query surface: enforced configuration, per-query scan ceiling, CloudWatch query metrics. |
| Delivery role | `AwsIamRole` | Least-privilege Firehose delivery, scoped to exactly the lake bucket. Rendered only while `ingestion_enabled` is on. |
| Ingestion stream | `AwsKinesisFirehose` | Direct-PUT delivery: batching, GZIP, date partitions, error prefix. Rendered only while `ingestion_enabled` is on. |

## Parameters

| Name | Description | Default | Required |
|---|---|---|---|
| `aws_region` | Region for every resource; keep the workgroup and buckets together. | `us-east-1` | yes |
| `lake_bucket_name` | Globally unique lake bucket name; also the prefix for companion resources. | `my-org-data-lake` | yes |
| `database_name` | Glue database name — lowercase, digits, underscores only (no hyphens). | `lakehouse` | yes |
| `infrequent_access_after_days` | Days before objects tier to STANDARD_IA. | `30` | yes |
| `glacier_after_days` | Days before objects tier to GLACIER_IR (still Athena-queryable). Must exceed the IA threshold. | `180` | yes |
| `query_scan_limit_gib` | Per-query scan ceiling in GiB; Athena cancels anything larger. `0` removes the limit. | `1024` | yes |
| `managed_results_enabled` | AWS-managed result storage (free, 24h retention) instead of the results bucket. | `false` | no |
| `ingestion_enabled` | Provision the Firehose stream + delivery role. | `true` | no |
| `force_destroy` | Allow teardown to delete non-empty buckets. Keep off for real data. | `false` | no |

## First query (post-deploy)

1. **Land some data.** With ingestion on, send a few records:

   ```bash
   aws firehose put-record \
     --delivery-stream-name <lake_bucket_name>-ingest \
     --record '{"Data":"eyJldmVudCI6ICJzaWdudXAiLCAidXNlciI6ICJhZGEifQo="}'
   ```

   Records appear under `s3://<lake>/data/year=.../` within the buffering
   window (up to 5 minutes / 64 MiB).

2. **Define a table over the landed data.** Run in the chart's workgroup —
   partition projection means never running `MSCK REPAIR` as days roll over:

   ```sql
   CREATE EXTERNAL TABLE events (
     event string,
     user  string
   )
   PARTITIONED BY (year string, month string, day string)
   ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
   LOCATION 's3://<lake_bucket_name>/data/'
   TBLPROPERTIES (
     'projection.enabled'      = 'true',
     'projection.year.type'    = 'integer', 'projection.year.range'  = '2024,2100',
     'projection.month.type'   = 'integer', 'projection.month.range' = '1,12',  'projection.month.digits' = '2',
     'projection.day.type'     = 'integer', 'projection.day.range'   = '1,31',  'projection.day.digits'   = '2',
     'storage.location.template' = 's3://<lake_bucket_name>/data/year=${year}/month=${month}/day=${day}'
   );
   ```

3. **Query with partition pruning** — the WHERE clause on partition columns
   is what keeps scans (and the bill) small:

   ```sql
   SELECT event, count(*) FROM lakehouse.events
   WHERE year = '2026' AND month = '07'
   GROUP BY event;
   ```

## Cost model, honestly

- **S3 storage** is the steady cost; the tiering rules cut it ~45%
  (STANDARD_IA) and ~68% (GLACIER_IR) as data ages, with both tiers still
  millisecond-readable by Athena (small per-GB retrieval fee applies).
- **Athena** bills ~$5 per TB scanned. The workgroup's scan ceiling caps the
  worst single query; partition pruning and GZIP compression (both wired in
  by default) are what keep the routine ones cheap.
- **Firehose** bills per GB ingested (~$0.029/GB Direct PUT) — no hourly
  charge, nothing to pay while idle.
- **Managed query results** are free but vanish after 24 hours; the results
  bucket costs pennies and keeps 30 days.

## Day-2 guidance

- **Deep archive knowingly.** Add a `GLACIER` or `DEEP_ARCHIVE` transition
  to the lake's lifecycle rules only for prefixes you accept restoring
  before querying — Athena cannot read deep-freeze classes in place, which
  is why this chart's default path stops at GLACIER_IR.
- **Convert hot tables to Parquet.** Columnar formats cut scans another
  ~90% for wide tables: add a Glue table for the Parquet layout and either
  a scheduled `CREATE TABLE AS SELECT` compaction, or enable the Firehose
  stream's `dataFormatConversion` (the spec models it) once the Glue table
  exists to convert on ingest.
- **Grow into dynamic partitioning.** Partitioning by record fields
  (customer, event type) uses the stream's `dynamicPartitioning` +
  `metadataExtraction` processor arms; note AWS makes enablement
  create-time-only on a stream, so plan it as a stream replacement.
- **Tighten access by prefix.** Grant analysts `s3:GetObject` on
  `data/*` read paths and reserve `errors/*` for the pipeline owners; the
  bucket policy field on the lake bucket takes a standard IAM document.
- **Watch the workgroup metrics.** `publishCloudwatchMetricsEnabled` emits
  per-query scan volume — alarm on aggregate daily scan bytes to catch
  cost drift before the invoice does.
