# GCP Cloud Storage Bucket

Creates a Google Cloud Storage bucket — the durable object store behind static sites, data lakes, build artifacts, and backups — with the full lifecycle surface (versioning, lifecycle rules, Autoclass, soft delete, WORM retention), placement control up to custom dual-regions with turbo replication, and additive per-bucket IAM grants.

## What Gets Created

- The Cloud Storage API is enabled on the project (never disabled on destroy)
- A `google_storage_bucket` carrying your labels merged beneath Planton's attribution labels (`planton-ai_resource`, `planton-ai_name`, `planton-ai_kind`, plus org/env/id when set)
- One `google_storage_bucket_iam_member` per `iamMembers` entry — additive grants that compose safely with grants made elsewhere

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **IAM permissions** — `roles/storage.admin` on the target project
- For CMEK: a `GcpKmsKey` the GCS service agent can use

## Quick Start

Create a file `bucket.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpGcsBucket
metadata:
  name: app-data
spec:
  bucketName: acme-app-data-prod
  location: us-central1
  uniformBucketLevelAccessEnabled: true
  publicAccessPrevention: enforced
  versioningEnabled: true
  lifecycleRules:
    - action:
        type: Delete
      condition:
        withState: ARCHIVED
        numNewerVersions: 3
```

Deploy:

```shell
planton apply -f bucket.yaml
```

This creates a private, versioned bucket whose version history stays bounded — and which no IAM grant can ever expose publicly.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | StringValueOrRef | No | Owning project; empty uses the provider default. Immutable |
| `bucketName` | string | Yes | Globally unique bucket name. Immutable |
| `location` | string | Yes | Region, dual-region, or multi-region. Immutable |
| `storageClass` | string | No | `STANDARD` (default), `NEARLINE`, `COLDLINE`, `ARCHIVE` |
| `forceDestroy` | bool | No | Delete all objects on destroy (default false) |
| `uniformBucketLevelAccessEnabled` | bool | No | IAM-only access control (recommended) |
| `publicAccessPrevention` | string | No | `inherited` (default) or `enforced` |
| `versioningEnabled` | bool | No | Keep noncurrent versions |
| `autoclass` | object | No | Automatic storage classes (`enabled`, `terminalStorageClass`) |
| `lifecycleRules[]` | list | No | Delete / SetStorageClass / AbortIncompleteMultipartUpload over rich conditions |
| `retentionPolicy` | object | No | WORM retention (+ irreversible `isLocked`) |
| `softDeletePolicy.retentionDurationSeconds` | int | No | 0 (off) or 7–90 days; omitted follows GCP's 7-day default |
| `kmsKeyName` | StringValueOrRef | No | Default CMEK key (reference a `GcpKmsKey`) |
| `requesterPays` | bool | No | Callers pay for access/egress |
| `defaultEventBasedHold` | bool | No | Auto-hold every new object |
| `enableObjectRetention` | bool | No | Per-object retention (create-time only) |
| `website` | object | No | `mainPageSuffix` + `notFoundPage` |
| `corsRules[]` | list | No | Browser cross-origin access |
| `logging` | object | No | Access-log delivery to another bucket |
| `customPlacementConfig.dataLocations` | list | No | Exactly two regions (custom dual-region). Immutable |
| `rpo` | string | No | `DEFAULT` or `ASYNC_TURBO` |
| `hierarchicalNamespaceEnabled` | bool | No | Folder semantics (create-time only) |
| `labels` | map | No | User labels (merged beneath platform labels) |
| `iamMembers[]` | list | No | Additive grants: `role` + `member` (+ optional IAM `condition`) |
| `ipFilter` | object | No | Network-layer allowlist: public CIDR ranges and/or VPC networks (reference `GcpVpcNetwork`), evaluated before IAM |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `bucket_id` | Bucket ID (equals the bucket name) — what consumers reference |
| `bucket_name` | Bucket name |
| `url` | `gs://<name>` |
| `self_link` | API self link |
| `location` | Location as reported by GCS |
| `project_number` | Numeric owning project |

## Related Resources

- **GcpBackendBucket** — serve this bucket through the L7 load balancer (CDN, HTTPS)
- **GcpCloudFunction / GcpCloudRun / GcpCloudRunJob** — function sources and GCS volumes
- **GcpDataprocCluster / GcpCloudComposerEnvironment** — staging, temp, and DAG buckets
- **GcpPubSubTopic / GcpPubSubSubscription** — Cloud Storage ingestion and delivery
- **GcpKmsKey** — CMEK protection via `kmsKeyName`
