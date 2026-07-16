# GCP Backend Bucket

Creates a Compute Engine backend bucket — serves a Cloud Storage bucket's objects through an external HTTP(S) load balancer, optionally cached at Google's edge by Cloud CDN. URL maps route static paths (like `/assets/*`) to a backend bucket while dynamic paths go to backend services.

## What Gets Created

When you deploy a GcpBackendBucket resource, Planton provisions:

- **Backend Bucket** — a `google_compute_backend_bucket` pointing at the origin GCS bucket, with CDN policy, compression, custom response headers, and an optional Cloud Armor edge policy
- **Signed-URL Keys** (optional) — up to 3 named signing keys for serving private content through the CDN

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCS bucket** — referenced via `bucketName`; objects must be publicly readable or served via signed URLs
- **IAM permissions** — any role carrying `compute.backendBuckets.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBackendBucket
metadata:
  name: static-assets
spec:
  projectId:
    value: my-gcp-project-123
  bucketName:
    valueFrom:
      kind: GcpGcsBucket
      name: my-static-assets
      fieldPath: status.outputs.bucket_id
  enableCdn: true
  cdnPolicy:
    cacheMode: CACHE_ALL_STATIC
    defaultTtl: 3600
```

```shell
planton apply -f backend-bucket.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `bucketName` | `StringValueOrRef` | — | Required. The origin GCS bucket. Can reference a GcpGcsBucket. Mutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the backend bucket. Immutable. |
| `backendBucketName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable. |
| `enableCdn` | `bool` | `false` | Cache at Google's edge with Cloud CDN. |
| `cdnPolicy` | object | GCP defaults | Cache mode, TTLs, negative caching, cache keys, stale serving, request coalescing. |
| `compressionMode` | `string` | disabled | `AUTOMATIC` or `DISABLED`. |
| `customResponseHeaders` | `list(string)` | `[]` | Headers the LB adds; values may use `{cdn_cache_status}`. |
| `edgeSecurityPolicy` | `StringValueOrRef` | none | Cloud Armor EDGE policy. Can reference a GcpCloudArmorPolicy. |
| `loadBalancingScheme` | `string` | external | `INTERNAL_MANAGED` for cross-region internal ALBs (incompatible with CDN). Immutable. |
| `signedUrlKeys` | `list(object)` | `[]` | Up to 3 signing keys for CDN signed URLs; `keyValue` is secret material. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value URL maps reference as a path-rule target |
| `backend_bucket_name` | Name of the backend bucket in GCP |
| `bucket_name` | The origin GCS bucket currently being served |

## Related Components

- [GcpGcsBucket](/docs/catalog/gcp/gcpgcsbucket) — the origin bucket whose objects are served
- [GcpCloudArmorPolicy](/docs/catalog/gcp/gcpcloudarmorpolicy) — the edge security policy filtering requests
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the project that owns the backend bucket
