# AwsS3TableBucket

One S3 table bucket (S3 Tables — managed Apache Iceberg storage) with its full contents: namespaces, tables, resource policies, and replication — seven provider resources under one declarative owner.

## Highlights

- **A data lake as one resource**: the bucket, its logical databases (namespaces), and their Iceberg tables with schemas — declared together, the way charts want analytics storage.
- **AWS maintains the tables**: compaction, snapshot expiry, and unreferenced-file cleanup run continuously per the spec's dials — the ops work that makes self-managed Iceberg painful, gone.
- **Create-only schema, honestly taught**: the Iceberg schema is create-time input the provider never reads back; evolution happens through query engines (ALTER TABLE) — taught on the field and declared write-normalized for imports.
- **The untyped-object caveat carried for you**: the provider's encryption/maintenance arguments are untyped at the pin (plugin protocol v6 pending) — the typed spec is the load-bearing validation, and both modules always send complete, correctly-shaped objects.

## Both Engines

Both modules key namespaces by name and tables by `namespace//table`, and export the same outputs: `table_bucket_arn` (import ID), `owner_account_id`, `table_arns`, `table_warehouse_locations` (both keyed `namespace//table`).

## Chart Wiring

`kms_key_arn` references an AwsKmsKey; replication roles reference AwsIamRole. The `table_bucket_arn` output is what catalog integrations (Athena, Glue, EMR via the analytics-catalog integration) and cross-account policies reference; per-table ARNs come from the `table_arns` map.
