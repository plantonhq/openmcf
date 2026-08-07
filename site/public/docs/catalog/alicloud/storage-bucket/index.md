---
title: "Storage Bucket"
description: "Storage Bucket deployment documentation"
icon: "package"
order: 100
componentName: "alicloudstoragebucket"
---

# AliCloud Storage Bucket

Deploys an Alibaba Cloud Object Storage Service (OSS) bucket with configurable access control, versioning, server-side encryption, lifecycle management, CORS rules, and access logging. OSS is a cloud-native, S3-compatible object storage service. Bucket names are globally unique across all Alibaba Cloud accounts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **OSS Bucket** -- an `alicloud_oss_bucket` resource with the specified storage class, redundancy type, ACL, and optional encryption, versioning, lifecycle, CORS, and logging configurations

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **Globally unique bucket name** -- the `bucketName` must be unique across all Alibaba Cloud accounts. 3-63 characters, lowercase letters, digits, and hyphens only.
- **Immutable decisions** -- `storageClass` and `redundancyType` cannot be changed after creation.

## Deploy

### Console

Open the deployment store, find **AliCloud Storage Bucket**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including bucket name, storage class, encryption, and lifecycle rules.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudStorageBucket
metadata:
  name: app-assets
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  bucketName: acme-prod-app-assets
  redundancyType: ZRS
  versioningEnabled: true
  serverSideEncryption:
    sseAlgorithm: AES256
  tags:
    team: platform
```

```shell
planton apply -f alicloud-bucket.yaml
```

This creates a ZRS bucket with versioning and AES-256 encryption. A Stack Job tracks the provisioning in real time.

### InfraChart

OSS buckets are standalone resources with no upstream dependencies. Downstream components reference the bucket name directly in their configuration.

## Key Configuration

These are the most important decisions when configuring an OSS bucket. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Storage class** -- The `storageClass` field is immutable. "Standard" (default) for frequent access. "IA" for infrequent access with retrieval fees. "Archive" / "ColdArchive" / "DeepColdArchive" for long-term storage with increasing restore times.

**Redundancy** -- The `redundancyType` field is immutable. "LRS" (default) stores 3 copies in one zone. "ZRS" stores 3 copies across multiple zones for higher durability.

**Encryption** -- The `serverSideEncryption` block enables encryption at rest. "AES256" uses OSS-managed keys. "KMS" uses customer-managed keys via AliCloudKmsKey.

**Lifecycle rules** -- The `lifecycleRules` field automates object transitions (Standard -> IA -> Archive) and expiration. Combine with `noncurrentVersionExpirationDays` when versioning is enabled.

**Versioning** -- The `versioningEnabled` field preserves all object versions for recovery from accidental overwrites or deletes.

## Outputs and Dependencies

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `bucket_name` | Bucket name (also the resource ID) | Object upload, logging targets, replication |
| `extranet_endpoint` | Public internet endpoint | External client access, CDN origins |
| `intranet_endpoint` | VPC-internal endpoint | Zero-cost access from ECS/containers in the same region |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private standard** -- A minimal private bucket with default settings. Start from the **Private Standard** preset.

**Versioned encrypted** -- A production bucket with ZRS, versioning, and AES-256 encryption. Start from the **Versioned Encrypted** preset.

**Archive with lifecycle** -- A bucket with automatic tiering (Standard -> IA at 30d, Archive at 90d) and 365-day expiration. Start from the **Archive Lifecycle** preset.

## Works With

- [**AliCloud KMS Key**](/cloud-catalog/ali-cloud-kms-key) -- customer-managed encryption key for KMS-based server-side encryption
- [**AliCloud CDN Domain**](/cloud-catalog/ali-cloud-cdn-domain) -- accelerates bucket content delivery globally
