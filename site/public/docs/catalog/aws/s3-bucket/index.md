---
title: "S3 Bucket"
description: "S3 Bucket deployment documentation"
icon: "package"
order: 100
componentName: "awss3bucket"
---

# AWS S3 Bucket

Deploys an Amazon S3 bucket with its complete behavioral surface: versioning, default encryption, public-access posture, bucket policy, lifecycle management, replication, static website hosting, server access logging, CORS, event notifications, Object Lock, transfer acceleration, requester pays, and Intelligent-Tiering archive configurations.

## What Gets Created

When you deploy an AwsS3Bucket resource, Planton provisions:

- **S3 Bucket** — the bucket itself, named from `metadata.name` (immutable, globally unique), with identity tags
- **Public Access Block** — always configured; ALL four guards enabled unless `publicAccessBlock` explicitly relaxes them
- **Ownership Controls** — `BucketOwnerEnforced` (ACLs disabled) unless `objectOwnership` overrides it, plus an optional canned ACL when ACLs are re-enabled
- **Bucket Policy** — when `policy` is provided
- **Versioning** — when `versioningStatus` is set
- **Server-Side Encryption** — when `encryption` is set (otherwise AWS's own SSE-S3 default applies); supports SSE-S3, SSE-KMS, DSSE-KMS, and the S3 Bucket Key
- **Lifecycle Configuration** — transitions, expiration, noncurrent-version handling, and multipart-upload cleanup, when `lifecycleRules` are provided
- **Replication Configuration** — cross-region or same-region rules with RTC, replica KMS, and cross-account ownership translation, when `replication` is provided
- **Website Configuration** — website or redirect mode with routing rules, when `website` is provided
- **Access Logging** — when `logging` is provided, with optional partitioned log-object keys
- **CORS Configuration** — when `corsRules` are provided
- **Event Notifications** — EventBridge and/or Lambda/SQS/SNS targets, when `notification` is provided
- **Object Lock Default Retention** — when `objectLockDefaultRetention` is provided (requires `objectLockEnabled`)
- **Transfer Acceleration / Requester Pays / Intelligent-Tiering configurations** — when their fields are set

Both engines (Terraform/OpenTofu and Pulumi) implement the same contract with identical stack outputs.

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A KMS key** (or `AwsKmsKey` resource) if using SSE-KMS/DSSE-KMS
- **An IAM role and a versioned destination bucket** if configuring replication
- **Delivery permission on the target** (queue policy / topic policy / Lambda invoke permission) before configuring SQS/SNS/Lambda notifications — the EventBridge arm needs no such grant

## Quick Start

Create a file `bucket.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3Bucket
metadata:
  name: my-bucket
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsS3Bucket.my-bucket
spec:
  region: us-east-1
```

Deploy:

```shell
planton apply -f bucket.yaml
```

This creates a fully private bucket in `us-east-1`: all four public-access guards on, ACLs disabled, AWS-default SSE-S3 encryption.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the bucket will be created (e.g., `us-west-2`). | Required; non-empty |

### Bucket Root

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `forceDestroy` | `bool` | `false` | Delete all objects (including versions) when destroying the bucket. |
| `objectLockEnabled` | `bool` | `false` | Enable Object Lock (WORM). Immutable after creation; requires `versioningStatus: Enabled`. |

### Versioning and Encryption

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `versioningStatus` | `string` | unversioned | `Enabled` or `Suspended`. Enabled buckets can never return to unversioned. |
| `encryption.sseAlgorithm` | `string` | AWS default (SSE-S3) | `AES256`, `aws:kms`, or `aws:kms:dsse`. |
| `encryption.kmsKeyId` | `string \| valueFrom` | AWS-managed `aws/s3` key | Customer-managed KMS key for the KMS algorithms. References `AwsKmsKey`. |
| `encryption.bucketKeyEnabled` | `bool` | `false` | S3 Bucket Key — cuts KMS request costs by up to 99%. |

### Public Access, Ownership, and Policy

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `publicAccessBlock` | `object` | all four guards ON | Four independent guards: `blockPublicAcls`, `blockPublicPolicy`, `ignorePublicAcls`, `restrictPublicBuckets`. Absence = fully private. |
| `objectOwnership` | `string` | `BucketOwnerEnforced` | `BucketOwnerPreferred` / `ObjectWriter` re-enable ACLs for legacy patterns. |
| `acl` | `string` | — | Canned ACL; only valid when ownership re-enables ACLs (validated). |
| `policy` | `object` | — | Bucket resource policy as a native YAML structure — the primary access-control surface. |

### Lifecycle

| Field | Type | Description |
|-------|------|-------------|
| `transitionDefaultMinimumObjectSize` | `string` | `all_storage_classes_128K` (AWS default) or `varies_by_storage_class`. |
| `lifecycleRules[].id` | `string` | Unique rule identifier. |
| `lifecycleRules[].status` | `string` | `Enabled` (default) or `Disabled`. |
| `lifecycleRules[].filter` | `object` | Scope: `prefix`, `tags`, `objectSizeGreaterThan`, `objectSizeLessThan` (AND-combined). Absent = whole bucket. |
| `lifecycleRules[].transitions[]` | `object[]` | `{days XOR date, storageClass}` — STANDARD_IA, ONEZONE_IA, INTELLIGENT_TIERING, GLACIER_IR, GLACIER, DEEP_ARCHIVE. |
| `lifecycleRules[].expiration` | `object` | Exactly one of `days`, `date`, or `expiredObjectDeleteMarker`. |
| `lifecycleRules[].noncurrentVersionTransitions[]` | `object[]` | `{noncurrentDays, newerNoncurrentVersions, storageClass}`. |
| `lifecycleRules[].noncurrentVersionExpiration` | `object` | `{noncurrentDays, newerNoncurrentVersions}` — the essential cost control for versioned buckets. |
| `lifecycleRules[].abortIncompleteMultipartUploadDays` | `int` | Days to abort failed multipart uploads (7 is a common choice). |

### Replication

| Field | Type | Description |
|-------|------|-------------|
| `replication.roleArn` | `string \| valueFrom` | IAM role S3 assumes. References `AwsIamRole`. Required. |
| `replication.rules[].destination.bucketArn` | `string \| valueFrom` | Destination bucket (must be versioned). References `AwsS3Bucket`. |
| `replication.rules[].destination.account` | `string` | Destination account for cross-account replication. |
| `replication.rules[].destination.storageClass` | `string` | Storage class for replicas (empty keeps source class). |
| `replication.rules[].destination.replicaKmsKeyId` | `string \| valueFrom` | Destination KMS key; required when replicating SSE-KMS objects. |
| `replication.rules[].destination.metricsEnabled` / `replicationTimeControlEnabled` | `bool` | Replication metrics / the 15-minute RTC SLA (RTC requires metrics — validated). |
| `replication.rules[].filter` | `object` | `prefix` and/or `tags` scope. |
| `replication.rules[].deleteMarkerReplication` | `bool` | Keep deletions in sync. |
| `replication.rules[].existingObjectReplication` | `bool` | Replicate pre-existing objects (may require an AWS entitlement). |
| `replication.rules[].replicateSseKmsEncryptedObjects` | `bool` | Include SSE-KMS objects (requires `replicaKmsKeyId` — validated). |

### Website, Logging, CORS, Notifications

| Field | Type | Description |
|-------|------|-------------|
| `website.indexDocumentSuffix` / `errorDocumentKey` | `string` | Website mode (e.g., `index.html` / `error.html`). |
| `website.redirectAllRequestsTo` | `object` | Redirect-bucket mode — mutually exclusive with website mode (validated). |
| `website.routingRules[]` | `object[]` | Conditional redirects (`condition` + `redirect`). |
| `logging.targetBucket` | `string \| valueFrom` | Log destination bucket (same region). References `AwsS3Bucket`. |
| `logging.targetPrefix` | `string` | Log object key prefix. |
| `logging.partitionedPrefixDateSource` | `string` | `EventTime` or `DeliveryTime` — Athena-friendly partitioned log keys. |
| `corsRules[]` | `object[]` | `allowedMethods`, `allowedOrigins` (required), headers, `maxAgeSeconds`. |
| `notification.eventbridge` | `bool` | Deliver all events to EventBridge (no permission setup needed). |
| `notification.lambdaFunctions[]` / `queues[]` / `topics[]` | `object[]` | Targets by ARN or reference (`AwsLambda`, `AwsSqsQueue`, `AwsSnsTopic`) with `events` + key prefix/suffix filters. |

### Object Lock, Acceleration, Requester Pays, Intelligent-Tiering

| Field | Type | Description |
|-------|------|-------------|
| `objectLockDefaultRetention` | `object` | `{mode: GOVERNANCE\|COMPLIANCE, days XOR years}`. COMPLIANCE is absolute until expiry. |
| `accelerationStatus` | `string` | `Enabled` or `Suspended` — Transfer Acceleration. |
| `requestPayer` | `string` | `BucketOwner` (default) or `Requester`. |
| `intelligentTieringConfigurations[]` | `object[]` | Named archive configs: `{name, status, filterPrefix, filterTags, tiers: [{accessTier, days}]}` — ARCHIVE_ACCESS ≥ 90 days, DEEP_ARCHIVE_ACCESS ≥ 180 (validated). |

## Examples

### Private Versioned Bucket with KMS and Version Pruning

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3Bucket
metadata:
  name: app-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsS3Bucket.app-data
spec:
  region: us-east-1
  versioningStatus: Enabled
  encryption:
    sseAlgorithm: aws:kms
    kmsKeyId:
      valueFrom:
        kind: AwsKmsKey
        name: app-data-key
        fieldPath: status.outputs.key_arn
    bucketKeyEnabled: true
  lifecycleRules:
    - id: prune-noncurrent
      noncurrentVersionExpiration:
        noncurrentDays: 30
      abortIncompleteMultipartUploadDays: 7
```

### Log Bucket with Tiering and Expiration

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3Bucket
metadata:
  name: app-logs
  annotations:
    planton.dev/provisioner: tofu
    terraform.planton.dev/stack.name: prod.AwsS3Bucket.app-logs
spec:
  region: us-west-2
  lifecycleRules:
    - id: tier-and-expire
      filter:
        prefix: "logs/"
      transitions:
        - days: 30
          storageClass: STANDARD_IA
        - days: 90
          storageClass: GLACIER
      expiration:
        days: 365
      abortIncompleteMultipartUploadDays: 7
```

### Cross-Region Replication with RTC

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3Bucket
metadata:
  name: dr-source
spec:
  region: us-east-1
  versioningStatus: Enabled
  replication:
    roleArn:
      valueFrom:
        kind: AwsIamRole
        name: s3-replication-role
        fieldPath: status.outputs.role_arn
    rules:
      - id: everything-to-dr
        destination:
          bucketArn:
            valueFrom:
              kind: AwsS3Bucket
              name: dr-replica
              fieldPath: status.outputs.bucket_arn
          metricsEnabled: true
          replicationTimeControlEnabled: true
        deleteMarkerReplication: true
```

### Event Notifications into EventBridge

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsS3Bucket
metadata:
  name: media-uploads
spec:
  region: us-east-1
  notification:
    eventbridge: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `bucket_id` | `string` | Name/ID of the S3 bucket |
| `bucket_arn` | `string` | ARN of the bucket (e.g., `arn:aws:s3:::my-bucket`) |
| `region` | `string` | AWS region where the bucket was created |
| `bucket_regional_domain_name` | `string` | Regional domain name (e.g., `my-bucket.s3.us-east-1.amazonaws.com`) — also the right CloudFront origin domain |
| `bucket_domain_name` | `string` | Legacy global domain name (`my-bucket.s3.amazonaws.com`) |
| `hosted_zone_id` | `string` | Route53 hosted zone ID for the bucket's region, used for alias records |
| `website_endpoint` | `string` | Website endpoint (only when website hosting is configured) |
| `website_domain` | `string` | Website service domain for Route53 aliases (only when website hosting is configured) |

## Related Components

- [AwsKmsKey](/docs/catalog/aws/kms-key) — customer-managed key for SSE-KMS encryption
- [AwsS3ObjectSet](/docs/catalog/aws/s3-object-set) — manage objects inside a bucket
- [AwsCloudFront](/docs/catalog/aws/cloudfront) — serve a private bucket publicly with TLS and caching via Origin Access Control
- [AwsIamRole](/docs/catalog/aws/iam-role) — replication role
- [AwsSqsQueue](/docs/catalog/aws/sqs-queue) / [AwsSnsTopic](/docs/catalog/aws/sns-topic) / [AwsLambda](/docs/catalog/aws/lambda) — event notification targets
