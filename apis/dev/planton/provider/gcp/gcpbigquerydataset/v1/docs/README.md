# GcpBigQueryDataset -- Research & Design Documentation

## BigQuery Dataset in the GCP Ecosystem

BigQuery is Google Cloud's fully managed, serverless data warehouse designed for
large-scale analytics. A **dataset** is the top-level organizational container
within BigQuery -- it holds tables, views, routines (UDFs), and models. Datasets
are scoped to a project and bound to a geographic location.

The design boundary: the dataset owns **location, defaults, and the ACL**;
tables are their own composable resources. Infrastructure-owned tables
(partitioned fact tables, authorized views, external tables over GCS) are
first-class [GcpBigQueryTable](../../gcpbigquerytable/v1/docs/README.md)
resources that reference the dataset by ID, while application-owned schemas
(dbt models, migration-script tables) can keep managing their own tables
inside the same dataset. Both coexist: the dataset never enumerates its
tables, so IaC-owned and application-owned tables never fight.

## Deployment Landscape

### Method Comparison

| Method | Dataset | Tables | Access | CMEK | Lifecycle |
|--------|---------|--------|--------|------|-----------|
| GCP Console | Yes | Yes | Yes | Yes | Manual |
| `bq` CLI | Yes | Yes | Yes | Yes | Manual |
| Terraform (`google_bigquery_dataset`) | Yes | Separate resource | Yes (authoritative) | Yes | IaC |
| Pulumi (`bigquery.Dataset`) | Yes | Separate resource | Yes | Yes | IaC |
| Planton (this component) | Yes | Separate kind (GcpBigQueryTable) | Yes (authoritative) | Yes | IaC |
| dbt | No | Yes (models) | No | No | Schema |
| Dataform | No | Yes (SQL workflows) | No | No | Schema |

Planton fills the IaC gap for dataset provisioning with cross-resource composability
that Terraform and Pulumi lack natively.

## 90/10 Coverage

| Provider surface | Modeled | Notes |
|---|---|---|
| `dataset_id`, `location`, `project` | Yes | Immutables; ambient-project fallback when project omitted |
| `friendly_name`, `description` | Yes | |
| `labels` | Yes | User labels merged beneath platform labels |
| `resource_tags` | Yes | IAM-condition/org-policy-capable tags |
| `default_table_expiration_ms`, `default_partition_expiration_ms` | Yes | |
| `max_time_travel_hours` | Yes | 48–168; provider types it as string, spec keeps honest int32 |
| `is_case_insensitive`, `default_collation` | Yes | |
| `storage_billing_model` | Yes | LOGICAL / PHYSICAL |
| `delete_contents_on_destroy` | Yes | Destroy-safety guard |
| `default_encryption_configuration.kms_key_name` | Yes | FK -> GcpKmsKey |
| `access` (all 8 arms + condition) | Yes | Authoritative ACL; shape enforced pre-deploy |
| `external_dataset_reference` | Yes | BigQuery Omni / AWS Glue federation |
| `external_catalog_dataset_options` | Yes | Hive Metastore interop |
| `deletion_policy` | No (recorded skip) | Not on the released 6.x provider line |

## Field Analysis

### Immutable Fields (ForceNew)

These fields cannot be changed after creation. Any change destroys and recreates
the dataset:

- `dataset_id` -- the BigQuery dataset identifier
- `location` -- geographic data residency
- `is_case_insensitive` -- case sensitivity behavior
- `external_dataset_reference` -- federated-source binding

### Mutable Fields

- `friendly_name`, `description` -- metadata
- `labels`, `resource_tags` -- attribution and governance
- `default_table_expiration_ms`, `default_partition_expiration_ms` -- lifecycle
- `max_time_travel_hours` -- time travel configuration
- `default_collation` -- collation settings
- `storage_billing_model` -- billing model
- `default_encryption_configuration` -- CMEK key
- `access` -- access control entries
- `external_catalog_options` -- catalog interop metadata
- `delete_contents_on_destroy` -- client-side safety flag

### Labels Support

BigQuery datasets support GCP labels. User labels from `spec.labels` are merged
first, then Planton's platform attribution labels (which win on key conflicts),
identically on both engines:

- `planton-ai_resource: true`
- `planton-ai_name: <dataset_id>`
- `planton-ai_kind: gcpbigquerydataset`
- `planton-ai_organization: <metadata.org>` (if set)
- `planton-ai_environment: <metadata.env>` (if set)
- `planton-ai_id: <metadata.id>` (if set)

## Access Control Model

BigQuery datasets support two access control models:

### 1. Dataset-Level Access (What Planton Models)

Access entries are embedded in the dataset resource itself. This is the
**authoritative** model -- the listed entries become the dataset's complete
ACL. Each entry takes exactly one of two shapes, enforced pre-deploy:

- **Principal grant**: `role` + exactly one of `user_by_email`,
  `group_by_email`, `domain`, `special_group`, `iam_member` -- optionally
  gated by an IAM `condition` (e.g. a time-bounded contractor grant).
- **Resource authorization**: exactly one of `view`, `routine`, or `dataset`,
  with NO role. BigQuery grants the authorized resource implicit read access
  -- the standard pattern for exposing filtered slices of sensitive data
  through authorized views, UDFs, or a whole grantee dataset's views.

### 2. Additive Grants (Not Modeled -- Recorded Skip)

The provider also offers `google_bigquery_dataset_access` (one additive grant
per resource) and `google_bigquery_dataset_iam_*` (project-IAM-style
bindings). These are deliberately not modeled: mixing additive grants with the
dataset's authoritative inline ACL produces permanent diffs (each write
clobbers the other), and a single authoritative list is the honest,
reviewable source of truth for who can read a dataset.

### Default Access Behavior

When `access` is omitted from the spec, BigQuery applies default access:
- Project owners get `OWNER`
- Project editors get `WRITER`
- Project viewers get `READER`

When `access` is specified, it becomes authoritative -- unlisted entries are
removed, including those defaults. To keep them, list them explicitly.

## CMEK (Customer-Managed Encryption Keys)

BigQuery encrypts all data at rest by default using Google-managed keys. CMEK
provides additional control:

- The KMS key must be in the **same region** as the dataset
- For multi-regional datasets (US, EU), the KMS key must also be multi-regional
- The BigQuery service agent (`bq-<project-number>@bigquery-encryption.iam.gserviceaccount.com`)
  must hold `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key before the
  first table is written
- CMEK is set at the dataset level (affects all new tables)
- Individual tables can override with their own CMEK
- Changing the CMEK key only affects new table data; existing data remains
  encrypted with the previous key

### CMEK and Planton Composability

The `kms_key_name` field uses `StringValueOrRef` with `default_kind = GcpKmsKey`.
In infra charts, this enables:

```yaml
kmsKeyName:
  valueFrom:
    kind: GcpKmsKey
    name: analytics-cmek
    fieldPath: status.outputs.key_id
```

This creates a dependency edge: the KMS key is provisioned before the dataset.

## Storage Billing Model

BigQuery offers two billing models for storage:

| Model | Charges For | Best For |
|-------|------------|----------|
| LOGICAL (default) | Uncompressed data size | Small datasets, predictable billing |
| PHYSICAL | Compressed on-disk size | Large datasets with compressible data |

PHYSICAL billing typically reduces storage costs 60-80% for:
- JSON columns (highly compressible)
- Repeated string values
- Wide tables with many nullable columns

The tradeoff: PHYSICAL billing includes time travel and fail-safe storage in the
billing calculation, while LOGICAL billing does not. Switching to PHYSICAL is
allowed once every 14 days.

## Time Travel

BigQuery's time travel feature allows querying data at any point within the
configured window (48-168 hours). This affects:

- **Cost**: Time travel storage is billed under PHYSICAL model
- **Recovery**: Longer windows provide more recovery options
- **Point-in-time queries**: `SELECT * FROM table FOR SYSTEM_TIME AS OF timestamp`

Reducing `max_time_travel_hours` to 48 (minimum) saves storage costs but limits
the recovery window to 2 days.

## Federation and Catalog Interop

- **`external_dataset_reference`** turns the dataset into a read-only
  projection of an external source (currently AWS Glue databases) through a
  BigQuery Omni connection -- the dataset then contains no BigQuery-managed
  tables of its own. Immutable: converting to/from an external reference
  recreates the dataset.
- **`external_catalog_options`** attaches Hive-Metastore-compatible metadata
  (default storage URI + database parameters) so open-source engines like
  Spark can address the dataset as a Hive database -- the dataset side of the
  BigLake/Iceberg interop story whose table side lives on GcpBigQueryTable.

## Infra-Chart Composability

GcpBigQueryDataset is a **Layer 1** resource in infra chart topology:

```
Layer 0: GcpProject
Layer 0-1: GcpKmsKeyRing -> GcpKmsKey
Layer 1: GcpBigQueryDataset (references Project, optionally KmsKey)
Layer 2: GcpBigQueryTable (references the dataset)
Layer 2+: Application tables, views, dbt models (not IaC)
```

The dataset participates in these infra chart patterns:

- **data-analytics-environment**: BigQuery Dataset + Dataproc + PubSub + GCS + SA
- **ml-notebook-environment**: BigQuery Dataset + Vertex AI Notebook + GCS + SA
- **event-pipeline**: BigQuery Dataset + PubSub + Cloud Function + GCS

### Key Outputs for Composition

| Output | Used By |
|--------|---------|
| `dataset_id` | GcpBigQueryTable, GcpGkeCluster usage export, application code, dbt, SQL |
| `self_link` | API references, audit trails |
| `project` | Cross-project dataset references |
| `location` | Co-locating downstream datasets/tables (joins require one location) |
| `etag` | Optimistic-concurrency API callers |

## Deliberate Exclusions (recorded skips)

| Feature | Reason |
|---------|--------|
| `deletion_policy` | Not on the released 6.x provider line; `delete_contents_on_destroy` covers the destroy-safety story. |
| `google_bigquery_dataset_access` / dataset IAM trio | Additive grants fight the authoritative inline ACL (permanent diffs); one source of truth wins. |
| Routines, row-access policies, jobs | Separate BigQuery resources outside the dataset boundary; revisit on concrete pull. |
| Reservations, connections, data transfer, Analytics Hub, data policies | Separate product families with their own APIs. |

## Best Practices

1. **Use descriptive dataset IDs** -- e.g., `analytics_prod`, `raw_events`
2. **Choose location carefully** -- immutable after creation, affects latency and compliance
3. **Set CMEK for regulated data** -- PCI, HIPAA, SOX workloads
4. **Consider PHYSICAL billing** -- significant savings for compressible data
5. **Be explicit about access** -- document who has access and why; re-state
   the project defaults when you need them alongside custom grants
6. **Use table expiration for staging** -- prevent data sprawl in dev/staging
7. **Keep time travel at 168 hours** -- maximum recovery window unless cost is critical
