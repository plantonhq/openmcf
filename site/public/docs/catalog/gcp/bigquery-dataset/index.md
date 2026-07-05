---
title: "BigQuery Dataset"
description: "BigQuery Dataset deployment documentation"
icon: "package"
order: 100
componentName: "gcpbigquerydataset"
---

# GCP BigQuery Dataset

Deploys a GCP BigQuery dataset with configurable data location, table lifecycle defaults, an authoritative access control list, optional CMEK encryption, and BigQuery Omni / open-source catalog interoperability. The dataset is the top-level container for tables, views, and routines in BigQuery; infrastructure-owned tables are separate [GcpBigQueryTable](/docs/catalog/gcp/bigquery-table) resources that reference it.

## What Gets Created

When you deploy a GcpBigQueryDataset resource, Planton provisions:

- **BigQuery API enablement** — the `bigquery.googleapis.com` service is enabled on the target project so a fresh project works first try
- **BigQuery Dataset** — a `google_bigquery_dataset` resource in the specified project and location, labeled with Planton's `planton-ai_*` attribution labels merged over your own labels
- **Access Control Entries** — if the `access` field is provided, an authoritative ACL of principal grants (optionally IAM-condition-gated) and authorized view/routine/dataset entries
- **CMEK Encryption Configuration** — if `kmsKeyName` is provided, all new tables in the dataset default to encryption with the specified Cloud KMS key

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project**; when `projectId` is omitted the provider's default project is used
- **A Cloud KMS key** if enabling customer-managed encryption, with the BigQuery service agent granted `roles/cloudkms.cryptoKeyEncrypterDecrypter`

## Quick Start

Create a file `bigquery-dataset.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBigQueryDataset
metadata:
  name: my-dataset
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpBigQueryDataset.my-dataset
spec:
  datasetId: analytics_events
  location: US
```

Deploy:

```shell
planton apply -f bigquery-dataset.yaml
```

This creates a BigQuery dataset named `analytics_events` in the `US` multi-region in the provider's default project, with default access (project owners = OWNER, project editors = WRITER, project viewers = READER).

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `datasetId` | `string` | Unique identifier for the dataset within the project. Only letters, numbers, and underscores. Immutable after creation. | Required; pattern `^[0-9A-Za-z_]+$`; max 1024 chars |
| `location` | `string` | Geographic location where the dataset resides (e.g., `US`, `EU`, `us-central1`). Immutable after creation. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | Provider default project | GCP project where the dataset is created. Can reference a GcpProject resource via `valueFrom`. |
| `friendlyName` | `string` | — | User-friendly display name for the dataset. |
| `description` | `string` | — | Description of the dataset's contents or purpose. |
| `labels` | `map<string,string>` | — | User labels for cost attribution, merged beneath Planton's platform labels (platform keys win on conflict). |
| `resourceTags` | `map<string,string>` | — | Resource Manager tags (`"123456789012/environment": "production"`); unlike labels, tags participate in IAM conditions and org policy. |
| `defaultTableExpirationMs` | `int64` | `0` (no expiration) | Default lifetime for NEW tables in the dataset, in milliseconds. Minimum 3600000 (1 hour). |
| `defaultPartitionExpirationMs` | `int64` | `0` (no expiration) | Default expiration for partitions in new partitioned tables, in milliseconds. |
| `maxTimeTravelHours` | `int32` | `168` (7 days) | Hours of point-in-time recovery. Range: 48–168. Lower values reduce storage costs on PHYSICAL billing. |
| `isCaseInsensitive` | `bool` | `false` | When `true`, dataset and table names are case-insensitive. Immutable after creation. |
| `defaultCollation` | `string` | — | Default collation for string columns in new tables. Use `und:ci` for case-insensitive collation. |
| `storageBillingModel` | `string` | `LOGICAL` | Billing model: `LOGICAL` (uncompressed bytes) or `PHYSICAL` (compressed bytes, can reduce costs 60–80%; switching allowed once every 14 days). |
| `deleteContentsOnDestroy` | `bool` | `false` | When `true`, all tables are deleted when the dataset is destroyed. When `false`, destroy fails if the dataset contains tables. |
| `kmsKeyName` | `StringValueOrRef` | — | Cloud KMS key for default table encryption (CMEK). Can reference a GcpKmsKey resource via `valueFrom`. |
| `access` | `GcpBigQueryDatasetAccessEntry[]` | Default project access | Authoritative access control entries — the complete ACL; anything not listed is removed. See access entry fields below. |
| `externalDatasetReference` | `object` | — | Makes the dataset a read-only projection of an external source (AWS Glue) via a BigQuery Omni connection. Immutable. Fields: `externalSource`, `connection`. |
| `externalCatalogOptions` | `object` | — | Hive Metastore compatibility metadata for open-source engines. Fields: `defaultStorageLocationUri`, `parameters`. |

### Access Entry Fields

Each entry is either a **principal grant** (`role` + exactly one of the five principal identities, optionally gated by `condition`) or a **resource authorization** (exactly one of `view`, `routine`, `dataset` — no role). The shape is validated pre-deploy.

| Field | Type | Description |
|-------|------|-------------|
| `access[].role` | `string` | Role for principal grants (e.g., `OWNER`, `WRITER`, `READER`, `roles/bigquery.dataViewer`). Must be omitted for view/routine/dataset authorizations. |
| `access[].userByEmail` | `string` | Email address of a Google Account. |
| `access[].groupByEmail` | `string` | Email address of a Google Group. |
| `access[].domain` | `string` | Domain to grant access to (e.g., `example.com`). |
| `access[].specialGroup` | `string` | Special group: `projectOwners`, `projectReaders`, `projectWriters`, or `allAuthenticatedUsers`. |
| `access[].iamMember` | `string` | IAM member expression (e.g., `serviceAccount:sa@project.iam.gserviceaccount.com`). |
| `access[].view` | `object` | Authorized view: `projectId`, `datasetId`, `tableId`. |
| `access[].routine` | `object` | Authorized routine (UDF / stored procedure): `projectId`, `datasetId`, `routineId`. |
| `access[].dataset` | `object` | Authorized dataset: `projectId`, `datasetId`, `targetTypes` (currently `["VIEWS"]`). |
| `access[].condition` | `object` | IAM condition on a principal grant: `expression` (CEL), `title`, `description`, `location`. |

## Examples

### Dataset with Table Expiration

Automatically delete tables after 90 days, useful for staging or transient data:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBigQueryDataset
metadata:
  name: staging-events
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: staging.GcpBigQueryDataset.staging-events
spec:
  datasetId: staging_events
  location: us-central1
  friendlyName: Staging Events
  description: Transient event data with 90-day auto-expiration
  defaultTableExpirationMs: 7776000000
  maxTimeTravelHours: 48
  deleteContentsOnDestroy: true
```

### Dataset with CMEK Encryption and Physical Billing

Production dataset using customer-managed encryption and physical storage billing for cost optimization:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBigQueryDataset
metadata:
  name: prod-analytics
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpBigQueryDataset.prod-analytics
spec:
  projectId:
    value: my-gcp-project
  datasetId: prod_analytics
  location: US
  friendlyName: Production Analytics
  description: Core analytics dataset with CMEK and physical billing
  storageBillingModel: PHYSICAL
  kmsKeyName:
    value: projects/my-gcp-project/locations/us/keyRings/analytics-ring/cryptoKeys/analytics-key
```

### Dataset with Explicit Access Control

Grant access to specific users, groups, a time-bounded contractor, and an authorized view:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBigQueryDataset
metadata:
  name: finance-data
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpBigQueryDataset.finance-data
spec:
  datasetId: finance_data
  location: EU
  friendlyName: Finance Data
  description: Restricted financial data with explicit access grants
  isCaseInsensitive: true
  defaultCollation: "und:ci"
  access:
    - role: OWNER
      userByEmail: data-owner@example.com
    - role: WRITER
      groupByEmail: data-engineers@example.com
    - role: READER
      groupByEmail: analysts@example.com
    - role: READER
      iamMember: "serviceAccount:etl-pipeline@my-gcp-project.iam.gserviceaccount.com"
    - role: READER
      userByEmail: contractor@example.com
      condition:
        title: expires-2027
        expression: request.time < timestamp("2027-01-01T00:00:00Z")
    - view:
        projectId: my-gcp-project
        datasetId: reporting_views
        tableId: finance_summary
```

### Using Foreign Key References

Reference other Planton-managed resources instead of hardcoding values:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBigQueryDataset
metadata:
  name: ref-dataset
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpBigQueryDataset.ref-dataset
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  datasetId: warehouse
  location: us-central1
  kmsKeyName:
    valueFrom:
      kind: GcpKmsKey
      name: warehouse-key
      fieldPath: status.outputs.key_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `dataset_id` | `string` | The short dataset ID (same as the spec's `datasetId` input), used in BigQuery SQL queries and API calls |
| `self_link` | `string` | Fully qualified URI of the dataset (e.g., `https://bigquery.googleapis.com/bigquery/v2/projects/{project}/datasets/{dataset}`) |
| `project` | `string` | The GCP project that contains this dataset (resolved even under the ambient-project fallback) |
| `creation_time` | `int64` | Creation time of the dataset in milliseconds since epoch |
| `location` | `string` | Geographic location of the dataset; tables joined in one query must share it |
| `etag` | `string` | Entity tag, changing on every metadata modification |

## Related Components

- [GcpBigQueryTable](/docs/catalog/gcp/bigquery-table) — infrastructure-owned tables, views, and external tables inside this dataset
- [GcpProject](/docs/catalog/gcp/project) — provides the GCP project where the dataset is created
- [GcpKmsKeyRing](/docs/catalog/gcp/kms-key-ring) — provides the key ring containing KMS keys for CMEK encryption
- [GcpKmsKey](/docs/catalog/gcp/kms-key) — provides the Cloud KMS encryption key referenced by `kmsKeyName`
- [GcpServiceAccount](/docs/catalog/gcp/service-account) — creates service accounts that can be granted dataset access
