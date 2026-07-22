# GCP Cloud Storage Bucket

Deploys a Google Cloud Storage bucket (`google_storage_bucket`) — the durable object store behind static sites, data lakes, build artifacts, backups, and every GCP service that stages data — with additive per-bucket IAM grants.

## Overview

The spec covers the full bucket lifecycle surface:

- **Access model** — uniform bucket-level access (IAM-only, the modern posture), public access prevention, and additive `iamMembers` grants (optionally condition-scoped to prefixes or expiry).
- **Data lifecycle** — object versioning, lifecycle rules (delete, storage-class transitions, multipart-upload cleanup) with explicit-zero semantics, Autoclass (automatic per-object storage classes), soft delete (the 7-day recovery window GCP applies by default), and WORM retention with the irreversible lock.
- **Placement and durability** — regions, multi-regions, predefined and custom dual-regions (pick your two regions), and turbo replication (`rpo: ASYNC_TURBO`).
- **Safety** — `forceDestroy` defaults to false: destroying a non-empty bucket fails instead of silently erasing data.

`bucketName`, `location`, project, custom placement, hierarchical namespace, and `enableObjectRetention` are immutable — changing them replaces the bucket and everything in it.

## When to Use

- **Application object storage** — user uploads, reports, exports
- **Static websites** — public buckets behind the L7 load-balancer family (`GcpBackendBucket`)
- **Data lakes** — Dataproc staging, BigQuery external data, pipeline outputs
- **Artifacts and backups** — build outputs, database exports, disaster-recovery copies

## Prerequisites

- GCP credentials with `roles/storage.admin` on the target project (the Cloud Storage API is enabled automatically)
- For CMEK: a `GcpKmsKey` whose key the GCS service agent can use (`roles/cloudkms.cryptoKeyEncrypterDecrypter`)

## Quick Start

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
```

This creates a private STANDARD bucket that no IAM grant can ever expose publicly.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | StringValueOrRef | No | Owning project; empty uses the provider default. Immutable |
| `bucketName` | string | Yes | Globally unique bucket name (3-63 chars). Immutable |
| `location` | string | Yes | Region, dual-region, or multi-region. Immutable |
| `storageClass` | string | No | `STANDARD` (default), `NEARLINE`, `COLDLINE`, `ARCHIVE` |
| `forceDestroy` | bool | No | Delete all objects on destroy (default false — safe posture) |
| `uniformBucketLevelAccessEnabled` | bool | No | IAM-only access (recommended; GCP default false) |
| `publicAccessPrevention` | string | No | `inherited` (default) or `enforced` |
| `versioningEnabled` | bool | No | Keep noncurrent versions on overwrite/delete |
| `autoclass` | object | No | `enabled` + `terminalStorageClass` (`NEARLINE`/`ARCHIVE`) |
| `lifecycleRules[]` | list | No | Action (`Delete`/`SetStorageClass`/`AbortIncompleteMultipartUpload`) + condition (age, dates, prefixes, versions; explicit 0 expressible) |
| `retentionPolicy` | object | No | WORM: `retentionPeriodSeconds` + `isLocked` (irreversible) |
| `softDeletePolicy` | object | No | Recovery window: 0 (off) or 7–90 days; omitted follows GCP's 7-day default |
| `kmsKeyName` | StringValueOrRef | No | Default CMEK key (reference a `GcpKmsKey`) |
| `requesterPays` | bool | No | Callers pay for access/egress |
| `defaultEventBasedHold` | bool | No | Auto-hold every new object |
| `enableObjectRetention` | bool | No | Per-object retention. Create-time only |
| `website` | object | No | `mainPageSuffix` + `notFoundPage` |
| `corsRules[]` | list | No | Direct browser cross-origin access |
| `logging` | object | No | Access-log delivery to another bucket (referenceable) |
| `customPlacementConfig.dataLocations` | list | No | Exactly two regions (custom dual-region). Immutable |
| `rpo` | string | No | `DEFAULT` or `ASYNC_TURBO` (dual-region turbo replication) |
| `hierarchicalNamespaceEnabled` | bool | No | Real folder semantics. Create-time only |
| `labels` | map | No | User labels, merged beneath platform labels |
| `iamMembers[]` | list | No | Additive grants: `role` + `member` (+ optional IAM condition) |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `bucket_id` | Bucket ID (equals the bucket name) — what consumers reference |
| `bucket_name` | Bucket name |
| `url` | `gs://<name>` |
| `self_link` | API self link |
| `location` | Location as reported by GCS (upper-cased) |
| `project_number` | Numeric owning project |

See the [presets](presets/) for remixable starting points and [docs/README.md](docs/README.md) for the deep dive.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
