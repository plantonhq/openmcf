# GCP BigQuery Table

Deploys a BigQuery table — a native table, logical view, materialized view, or external table. All four shapes are arms of one GCP resource; the spec enforces mutual exclusivity pre-deploy.

## What Gets Created

When you deploy a GcpBigQueryTable resource, Planton provisions:

- **BigQuery API enablement** — the `bigquery.googleapis.com` service is enabled on the target project
- **BigQuery Table** — a `google_bigquery_table` in the referenced dataset, labeled with Planton's `planton-ai_*` attribution labels merged over your own labels

## Prerequisites

- **An existing BigQuery dataset** — referenced via `datasetId` (a [GcpBigQueryDataset](/docs/catalog/gcp/gcpbigquerydataset) resource or a literal dataset ID)
- **GCP credentials** with BigQuery data editor permissions on the project
- **A Cloud KMS key** if enabling customer-managed encryption, with the BigQuery service agent granted `roles/cloudkms.cryptoKeyEncrypterDecrypter`

## Quick Start

Create a file `bigquery-table.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBigQueryTable
metadata:
  name: my-events-table
spec:
  datasetId:
    valueFrom:
      kind: GcpBigQueryDataset
      name: my-analytics-dataset
      fieldPath: status.outputs.dataset_id
  tableId: events_raw
  schema: '[{"name":"event_time","type":"TIMESTAMP","mode":"REQUIRED"},{"name":"payload","type":"JSON"}]'
  timePartitioning:
    type: DAY
    field: event_time
```

Deploy:

```shell
planton apply -f bigquery-table.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `datasetId` | `StringValueOrRef` | Parent dataset. Immutable. | Ref → GcpBigQueryDataset `dataset_id` |
| `tableId` | `string` | Table ID within the dataset. Immutable. | Letters, numbers, underscores |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project |
| `schema` | `string` | — | JSON array of field definitions (native tables) |
| `timePartitioning` | object | — | Day/hour/month/year partitions (XOR with range) |
| `rangePartitioning` | object | — | Integer range partitions |
| `clustering` | `string[]` | — | Up to four clustering columns |
| `requirePartitionFilter` | `bool` | `false` | Require partition filter on queries |
| `view` | object | — | Logical view arm |
| `materializedView` | object | — | Materialized view arm |
| `externalDataConfiguration` | object | — | External/BigLake table arm |
| `deletionProtection` | `bool` | `true` | Client-side destroy guard (both engines) |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `table_id` | Short table ID |
| `self_link` | Fully qualified API URI |
| `project` | Resolved GCP project |
| `dataset_id` | Parent dataset |
| `type` | TABLE, VIEW, MATERIALIZED_VIEW, or EXTERNAL |
| `location` | Geographic location |
| `creation_time` | Milliseconds since epoch |

## Examples

See the component presets for partitioned analytics tables, authorized views,
and external GCS tables.

## Related Components

- [GcpBigQueryDataset](/docs/catalog/gcp/gcpbigquerydataset) — parent dataset
- [GcpKmsKey](/docs/catalog/gcp/gcpkmskey) — CMEK key reference
