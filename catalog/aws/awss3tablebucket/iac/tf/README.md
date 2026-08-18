# AwsS3TableBucket — Terraform/OpenTofu module

Manages one S3 table bucket and its full contents (`aws_s3tables_table_bucket` + policy + replication + `aws_s3tables_namespace` + `aws_s3tables_table` + table policy + table replication) - seven provider resources under one declarative owner.

Module facts worth knowing before editing:

- **Namespaces key by name; tables key by `namespace.table`** — stable across list reorders; renames replace (both are create-only at AWS).
- **A table's `iceberg_schema`/`properties` are CREATE-ONLY input** — the provider never reads them back (schema evolution happens through query engines via ALTER TABLE), so they cannot drift and never round-trip on import (declared in the import catalog).
- **`format` is module-pinned to ICEBERG** — the provider's OpenTableFormat enum holds exactly that one value (recorded parity exclusion).
- **The object-attribute caveat** — `encryption_configuration` and `maintenance_configuration` are untyped object attributes at the provider (awaiting plugin protocol v6), so no provider-side validation exists for their sub-fields; the typed spec is the load-bearing guard and this module always sends complete, correctly-shaped objects.
- **Policies are JSON-normalized by AWS** (importIgnore at the provider); **replication carries a `version_token`** optimistic-concurrency handshake the provider manages.
- **KMS-encrypted buckets need a maintenance grant** — the S3 Tables maintenance principal must be able to use the key or compaction silently stops (taught on the spec field and in the GUIDE).

Outputs mirror the Pulumi module key-for-key: `table_bucket_arn` (import ID), `owner_account_id`, `table_arns` (keyed `namespace.table`), `table_warehouse_locations`.
