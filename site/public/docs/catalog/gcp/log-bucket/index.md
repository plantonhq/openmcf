---
title: "Log Bucket"
description: "Log Bucket deployment documentation"
icon: "package"
order: 100
componentName: "gcplogbucket"
---

# GCP Log Bucket

Creates a Cloud Logging bucket — the container where log entries are STORED: how long they are retained, whether they are indexed for SQL analytics, who can see which slice (log views), and how the data links into BigQuery. One kind covers all four scopes (project, folder, organization, billing account), and adopting a name that already exists — including GCP's built-in `_Default` bucket — patches it into management instead of failing.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Log bucket** -- exactly ONE of `google_logging_{project|folder|organization|billing_account}_bucket_config`, selected by the spec's `scope`
- **Log views** (optional) -- one `google_logging_log_view` per `logViews[]` entry
- **Linked BigQuery dataset** (optional) -- a `google_logging_linked_dataset` (requires analytics)
- **Scope logging settings** (optional, folder/org only) -- the `google_logging_{folder|organization}_settings` singleton
- **Logging API enablement** -- `logging.googleapis.com` enabled on the project scope (never disabled on destroy)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target scope.
- **Planton Runner** -- required when using Runner-based credential delivery.

### GCP

- **IAM**: `roles/logging.configWriter` on the scope (folder/org scopes need it at that level).
- **CMEK buckets**: grant the Logging service account `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key FIRST, or creation fails.

## Deploy

### Console

Open the deployment store, find **GCP Log Bucket**, and click **Deploy**. Start from the **Compliance Retention** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpLogBucket
metadata:
  name: audit-logs
  org: acme-corp
  env: prod
spec:
  bucketId: audit-logs
  retentionDays: 400
  indexConfigs:
    - fieldPath: jsonPayload.request.status
      type: INDEX_TYPE_INTEGER
```

```shell
planton apply -f log-bucket.yaml
```

A 400-day project bucket with a status-code index — pair it with a GcpLoggingSink routing audit entries in.

### InfraChart

Storage and routing compose in one chart:

```yaml
# The sink routes INTO the bucket this kind creates:
spec:
  destination:
    rawUri:
      valueFrom:
        kind: GcpLogBucket
        name: audit-logs
        fieldPath: status.outputs.bucket_name
```

## Key Configuration

**scope** -- omit for a project bucket in the provider's default project (the common case); `folderId` / `organizationId` / `billingAccount` select the other scopes, which are ADOPT-only (the Logging API creates new custom buckets only under projects).

**retentionDays** -- 1 to 3650; GCP's default 30 is sent explicitly. On a locked bucket, retention can never change again.

**One-way doors** -- `locked` (retention frozen, delete only when empty), `enableAnalytics` (SQL over logs; never disables), `cmekKmsKey` (CMEK never turns off; key rotation allowed). Each is deliberate — read the field docs before arming.

**logViews[]** -- named, filtered slices grantable independently (`roles/logging.viewAccessor` on a view shows a reader only its filter's entries) — how one bucket serves audit and app teams with different eyes.

**linkedBigqueryDataset** -- a read-only BigQuery dataset over the bucket (requires analytics); every field is create-time-only.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `scope.projectId` | `status.outputs.project_id` |
| **GcpKmsKey** (optional) | `cmekKmsKey`, `scopeSettings.kmsKey` | `status.outputs.key_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `bucket_name` | Full bucket resource name | GcpLoggingSink destination (`logging.googleapis.com/{bucket_name}`), bucket-scoped GcpLogMetric |
| `linked_dataset_id` | Linked BigQuery dataset ID | Queries from BigQuery tooling |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Compliance retention** -- long retention + status index, PREVENT posture. Start from the **Compliance Retention** preset.

**Analytics + views** -- Log Analytics on, an errors-only view, and a linked BigQuery dataset. Start from the **Analytics And Views** preset.

**Adopt _Default** -- manage the built-in default bucket's retention declaratively. Start from the **Adopt Default Bucket** preset.

## Works With

- [**GCP Logging Sink**](/cloud-catalog/gcp-logging-sink) -- routes matching entries INTO this bucket
- [**GCP Log Metric**](/cloud-catalog/gcp-log-metric) -- bucket-scoped metrics count entries as they land here
- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- supplies the CMEK key
- [**GCP BigQuery Dataset**](/cloud-catalog/gcp-big-query-dataset) -- the linked dataset appears beside your other datasets
