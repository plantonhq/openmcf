# AWS S3 Object Set

Deploys one or more objects into an existing AWS S3 bucket, supporting inline text content and base64-encoded binary content. The component manages objects declaratively alongside infrastructure — configuration files, static website assets, seed data, deployment artifacts — with the full per-object surface: HTTP presentation headers, user metadata, storage class, encryption overrides, upload checksums, Object Lock retention, and access control.

## What Gets Created

When you deploy an AwsS3ObjectSet resource, Planton provisions:

- **S3 Object (one per entry)** — an `aws_s3_object` resource for each item in the `objects` list, uploaded to the target bucket with the specified key, content, and per-object settings. Each object's identity is its key: adding, removing, or reordering entries never churns unrelated objects.
- **Merged Tags** — each object receives tags merged from three sources in increasing precedence: resource-identity labels, set-level `tags`, and per-object `tags`

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **An existing S3 bucket** — either a literal bucket name or a deployed AwsS3Bucket resource to reference via `valueFrom`
- **The bucket's AWS region** — must match the region specified in `region`

## Quick Start

Create a file `s3-objects.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3ObjectSet
metadata:
  name: my-objects
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsS3ObjectSet.my-objects
spec:
  region: us-east-1
  bucket:
    value: my-app-bucket
  objects:
    - key: config/app.json
      content: '{"env": "dev", "debug": true}'
      contentType: application/json
```

Deploy:

```shell
planton apply -f s3-objects.yaml
```

This uploads a single JSON configuration file to the `config/app.json` key in the target bucket.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | The AWS region where the target bucket lives (e.g., `us-west-2`). Must match the bucket's own region. | Required; non-empty |
| `bucket` | `StringValueOrRef` | The target S3 bucket. A literal name (`value:`) or a reference to an AwsS3Bucket resource (`valueFrom:`, resolving `status.outputs.bucket_id`). | Required |
| `objects` | `AwsS3Object[]` | The objects to upload. Each object's `key` must be unique within the set. | Minimum 1 item |
| `objects[].key` | `string` | The S3 object key — the object's path within the bucket and its identity in the set. Changing a key replaces that object. | Minimum length 1 |
| `objects[].content` / `objects[].contentBase64` | `string` | The content source: inline UTF-8 text or base64-encoded binary. Exactly one must be set per object. | Exactly one, non-empty |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `tags` | `map<string, string>` | `{}` | Tags applied to every object in the set. Object-level tags merge on top. |
| `objects[].contentType` | `string` | `application/octet-stream` | MIME type stored with the object and returned on every GET. Set it correctly for anything a browser will fetch. |
| `objects[].cacheControl` | `string` | — | Cache-Control header (e.g., `max-age=86400`, `no-cache`, `public, max-age=31536000, immutable`). |
| `objects[].contentEncoding` | `string` | — | Content-Encoding for pre-compressed content (e.g., `gzip`, `br`). |
| `objects[].contentDisposition` | `string` | — | Content-Disposition header (`inline`, `attachment; filename="report.pdf"`). |
| `objects[].contentLanguage` | `string` | — | Content-Language header (e.g., `en-US`). |
| `objects[].metadata` | `map<string, string>` | `{}` | User-defined metadata stored as `x-amz-meta-*` headers. Keys must be lowercase. Immutable in S3 — changes rewrite the object. |
| `objects[].websiteRedirect` | `string` | — | Redirect target (`/other-key` or an absolute URL) served when the bucket has static website hosting. |
| `objects[].storageClass` | `string` | `STANDARD` | Storage class: `STANDARD`, `STANDARD_IA`, `ONEZONE_IA`, `INTELLIGENT_TIERING`, `GLACIER_IR`, `GLACIER`, `DEEP_ARCHIVE`, `REDUCED_REDUNDANCY`, or the special-infrastructure classes (`EXPRESS_ONEZONE`, `OUTPOSTS`, `SNOW`, `FSX_OPENZFS`, `FSX_ONTAP`). |
| `objects[].serverSideEncryption` | `string` | bucket default | Per-object encryption OVERRIDE: `AES256`, `aws:kms`, `aws:kms:dsse` (or `aws:fsx` for FSx-backed access points). Leave unset to inherit the bucket's default encryption — the recommended posture. |
| `objects[].kmsKey` | `StringValueOrRef` | — | KMS key ARN for SSE-KMS, as a literal or a reference to an AwsKmsKey resource. Implies `aws:kms` when `serverSideEncryption` is unset; pairing with `AES256` is rejected. |
| `objects[].bucketKeyEnabled` | `bool` | bucket default | Whether this object uses an S3 Bucket Key for SSE-KMS (batches KMS requests to cut KMS costs). |
| `objects[].checksumAlgorithm` | `string` | — | Upload-integrity checksum stored with the object: `CRC32`, `CRC32C`, `CRC64NVME`, `SHA1`, `SHA256`. |
| `objects[].objectLockMode` | `string` | — | Object Lock retention mode (`GOVERNANCE` or `COMPLIANCE`). Requires an Object Lock-enabled bucket and `objectLockRetainUntilDate`. |
| `objects[].objectLockRetainUntilDate` | `string` | — | RFC 3339 timestamp until which the object version is retained. Required with, and only valid with, `objectLockMode`. |
| `objects[].objectLockLegalHoldStatus` | `string` | — | Legal hold (`ON`/`OFF`) — an indefinite deletion block independent of retention mode. |
| `objects[].acl` | `string` | — | Canned ACL. Leave unset for modern buckets: BucketOwnerEnforced ownership (the AWS default) disables object ACLs and rejects this at apply time. |
| `objects[].forceDestroy` | `bool` | `false` | Bypass GOVERNANCE-mode Object Lock retention and legal holds when deleting this object (requires `s3:BypassGovernanceRetention`). Only valid on Object Lock-enabled buckets — S3 rejects the flag on regular buckets. Ordinary versioned-bucket destroys need no flag: all versions are always removed. |
| `objects[].tags` | `map<string, string>` | `{}` | Tags specific to this object. Merged over set-level tags. |

## Examples

### Multiple Configuration Files

Upload several configuration files to a shared bucket:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3ObjectSet
metadata:
  name: app-config
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsS3ObjectSet.app-config
spec:
  region: us-east-1
  bucket:
    value: my-app-bucket
  objects:
    - key: config/app.json
      content: '{"env": "dev", "logLevel": "debug"}'
      contentType: application/json
    - key: config/feature-flags.json
      content: '{"darkMode": true, "betaSignup": false}'
      contentType: application/json
    - key: config/robots.txt
      content: |
        User-agent: *
        Disallow: /admin/
      contentType: text/plain
```

### Static Website Assets with Caching

Upload static assets with per-asset cache posture (public access comes from the bucket's policy, not per-object ACLs):

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3ObjectSet
metadata:
  name: website-assets
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsS3ObjectSet.website-assets
spec:
  region: us-west-2
  bucket:
    valueFrom:
      kind: AwsS3Bucket
      name: my-website-bucket
      fieldPath: status.outputs.bucket_id
  tags:
    project: website
  objects:
    - key: index.html
      content: |
        <!DOCTYPE html>
        <html><head><title>My Site</title></head>
        <body><h1>Hello</h1></body></html>
      contentType: text/html
      cacheControl: no-cache
    - key: assets/style.css
      content: "body { font-family: sans-serif; margin: 0; }"
      contentType: text/css
      cacheControl: public, max-age=31536000, immutable
    - key: old/landing.html
      content: "<html><body>moved</body></html>"
      contentType: text/html
      websiteRedirect: /index.html
```

### Encrypted Compliance Drop with Object Lock

Per-object KMS encryption override plus WORM retention (the bucket must have been created with Object Lock enabled):

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3ObjectSet
metadata:
  name: audit-artifacts
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsS3ObjectSet.audit-artifacts
spec:
  region: us-east-1
  bucket:
    value: my-compliance-bucket
  objects:
    - key: audits/2026-q2-report.json
      content: '{"finding": "none", "auditor": "internal"}'
      contentType: application/json
      serverSideEncryption: aws:kms
      kmsKey:
        valueFrom:
          kind: AwsKmsKey
          name: audit-key
          fieldPath: status.outputs.key_arn
      checksumAlgorithm: SHA256
      objectLockMode: COMPLIANCE
      objectLockRetainUntilDate: "2033-01-01T00:00:00Z"
```

### Binary Content in a Cheaper Storage Class

Upload binary assets using base64-encoded content, placed in infrequent access:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3ObjectSet
metadata:
  name: binary-assets
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: staging.AwsS3ObjectSet.binary-assets
spec:
  region: eu-west-1
  bucket:
    value: my-assets-bucket
  objects:
    - key: images/favicon.ico
      contentBase64: AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAA...
      contentType: image/x-icon
      cacheControl: max-age=86400
      storageClass: STANDARD_IA
    - key: data/seed.csv
      content: |
        id,name,email
        1,Alice,alice@example.com
        2,Bob,bob@example.com
      contentType: text/csv
      metadata:
        generated-by: seed-pipeline
      forceDestroy: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `bucket_id` | `string` | The bucket the objects were uploaded to. |
| `object_arns` | `map<string, string>` | Map of object key to ARN (`arn:aws:s3:::bucket/key`), for IAM policy Resource lists. |
| `object_etags` | `map<string, string>` | Map of object key to its ETag (content hash). Changes when content changes — useful for cache invalidation. |
| `object_version_ids` | `map<string, string>` | Map of object key to its version ID. Only populated when the target bucket has versioning enabled. |

## Related Components

- [AwsS3Bucket](/docs/catalog/aws/awss3bucket) — provides the target bucket and owns the uniform security posture (default encryption, versioning, ownership controls); referenced via `valueFrom` in the `bucket` field
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — provides customer-managed keys for per-object SSE-KMS overrides via the `kmsKey` reference
- [AwsCloudFront](/docs/catalog/aws/awscloudfront) — serves objects from S3 via a CDN distribution
- [AwsLambda](/docs/catalog/aws/awslambda) — can be triggered by S3 object events, or deployed from a zip staged by this component
