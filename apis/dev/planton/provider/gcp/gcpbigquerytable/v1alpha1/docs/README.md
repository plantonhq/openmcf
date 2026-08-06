# GcpBigQueryTable — Research & Design Documentation

## BigQuery Table in the GCP Ecosystem

A **table** is the unit of storage and query in BigQuery. GCP models native
tables, logical views, materialized views, and external tables as one resource
type (`google_bigquery_table`) with mutually exclusive configuration arms.

The design boundary: tables reference a parent
[GcpBigQueryDataset](../../gcpbigquerydataset/v1alpha1/docs/README.md) for location
and ACL; the dataset never enumerates its tables, so IaC-owned and
application-owned tables coexist in the same dataset.

## Deployment Landscape

| Method | Native | View | Materialized View | External |
|--------|--------|------|-------------------|----------|
| GCP Console | Yes | Yes | Yes | Yes |
| `bq` CLI | Yes | Yes | Yes | Yes |
| Terraform | Yes | Yes | Yes | Yes |
| Planton (this component) | Yes | Yes | Yes | Yes |
| dbt / Dataform | Yes (models) | Yes | Partial | No |

## 90/10 Coverage

| Provider surface | Modeled | Notes |
|---|---|---|
| `table_id`, `dataset_id`, `project` | Yes | Immutables; ambient-project fallback |
| `schema` | Yes | JSON string (provider-authentic) |
| `time_partitioning` / `range_partitioning` | Yes | CEL mutual exclusivity |
| `clustering`, `require_partition_filter` | Yes | |
| `view` / `materialized_view` / `external_data_configuration` | Yes | CEL at-most-one arm |
| `encryption_configuration.kms_key_name` | Yes | FK → GcpKmsKey |
| `deletion_protection` | Yes | Default TRUE both engines |
| `table_constraints`, `biglake_configuration`, etc. | Yes | Released 6.x floor |
| `deletion_policy`, `csv_options.source_column_match` | No | Not on released 6.50.0 |
| Client drift knobs (`ignore_*`, `table_metadata_view`) | No | TF client tuning, not cloud state |

## Field Analysis

### Immutable Fields (ForceNew)

- `table_id`, `dataset_id`, `project`
- Partitioning field names
- `encryption_configuration.kms_key_name`
- `biglake_configuration`, `table_replication_info`
- Materialized view query

### Recorded Skips

- **`deletion_policy`** — not on released 6.x provider line
- **`ignore_auto_generated_schema` / `ignore_schema_changes` / `table_metadata_view`** — client-side drift handling, parity-neutral skip
- **`csv_options.source_column_match`** — not on released 6.50.0

## Composability

- **Inbound**: `datasetId` → GcpBigQueryDataset `outputs.dataset_id`
- **Self-reference**: `table_constraints.foreign_keys.referenced_table_id` → GcpBigQueryTable (same kind)
- **Encryption**: `encryption_configuration.kms_key_name` → GcpKmsKey

## E2E Coverage

Live scenarios (dual-engine on `planton-e2e`):

- Day-partitioned clustered native table on a dataset prerequisite
- Literal-SELECT view (no base-table dependency)

Excluded live (proven offline):

- CMEK — requires org-level BigQuery service agent KMS grant
- External GCS table — requires staged object data
