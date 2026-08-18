# AwsS3TableBucket — Pulumi module (Go)

Manages one S3 table bucket and its full contents (`s3tables.TableBucket` + policy + replication + `s3tables.Namespace` + `s3tables.Table` + table policy + table replication) - seven provider resources under one declarative owner.

Module facts worth knowing before editing:

- **Namespaces are named `namespace-{name}`; tables `table-{namespace.table}`** — stable across list reorders; renames replace (both create-only at AWS).
- **A table's iceberg schema/properties are CREATE-ONLY input** — the provider never reads them back (schema evolution happens through query engines via ALTER TABLE), so they cannot drift and never round-trip on import (declared in the import catalog).
- **`Format` is module-pinned to ICEBERG** — the provider's OpenTableFormat enum holds exactly that one value (recorded parity exclusion).
- **Policies are JSON-normalized by AWS** (importIgnore at the provider); **replication carries a `version_token`** optimistic-concurrency handshake the provider manages.
- **KMS-encrypted buckets need a maintenance grant** — the S3 Tables maintenance principal must be able to use the key or compaction silently stops (taught on the spec field and in the GUIDE).

Outputs mirror the Terraform module key-for-key: `table_bucket_arn` (import ID), `owner_account_id`, `table_arns` (keyed `namespace.table`), `table_warehouse_locations`.
