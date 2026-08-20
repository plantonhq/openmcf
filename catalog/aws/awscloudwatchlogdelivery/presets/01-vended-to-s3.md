# CloudFront Access Logs to S3

The modern access-log pipeline: a CloudFront source delivering Parquet to an S3 archive with Hive-compatible partitioning — Athena queries it with zero ETL. AWS prepends its own `AWSLogs/{account}/CloudFront/` prefix; `suffixPath` adds only your segment beneath it.
