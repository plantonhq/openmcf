---
title: "Object Bucket"
description: "Object Bucket deployment documentation"
icon: "package"
order: 100
componentName: "scalewayobjectbucket"
---

# Scaleway Object Bucket

Deploys an S3-compatible Object Storage bucket on Scaleway with optional versioning, lifecycle rules for automated tiering and expiration, CORS policies for browser-based access, and Object Lock for WORM compliance. Buckets are regional resources with globally unique names and serve as the primary storage container for application assets, backups, and archival data.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Object Bucket** -- a `scaleway_object_bucket` in the specified region with the configured versioning, lifecycle, and CORS settings
- **Versioning Configuration** -- created only when `versioningEnabled` is true; enables S3-compatible object version tracking with delete markers
- **Lifecycle Rules** -- created only when `lifecycleRules` entries are provided; automates object expiration, storage class transitions (GLACIER, ONEZONE_IA), and incomplete multipart upload cleanup
- **CORS Rules** -- created only when `corsRules` entries are provided; enables cross-origin browser access to the bucket's S3 endpoint
- **Scaleway Tags** -- S3-compatible key-value metadata tags applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair (Access Key + Secret Key). The IaC module authenticates through the Scaleway provider configuration.
- **Choose a region** -- buckets are regional resources. Available regions: `fr-par` (Paris), `nl-ams` (Amsterdam), `pl-waw` (Warsaw). Cannot be changed after creation.
- **Globally unique bucket name** -- the bucket name is derived from `metadata.name` and must be unique across all Scaleway Object Storage. Choose a DNS-compatible name unlikely to collide with other users.

## Deploy

### Console

Open the deployment store, find **Scaleway Object Bucket**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Private Bucket** preset in the [Presets](#presets) tab to create a minimal private bucket with no versioning or lifecycle rules.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayObjectBucket
metadata:
  name: app-media
  org: acme-corp
  env: prod
spec:
  region: fr-par
```

```shell
planton apply -f scaleway-object-bucket.yaml
```

This creates a private bucket in the Paris region with no versioning, no lifecycle rules, and no CORS. Objects are retained indefinitely and the bucket cannot be destroyed while it contains data. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring an object bucket. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Versioning** -- Set `versioningEnabled` to true to retain all versions of every object. Previous versions survive overwrites and deletions (delete markers are inserted instead of removing objects). Once enabled, versioning can only be suspended, not fully disabled. Combine with lifecycle rules to expire old versions and control storage costs.

**Lifecycle rules** -- The `lifecycleRules` array automates object management. Common patterns include expiring logs after N days (`expirationDays`), transitioning infrequently accessed data to `GLACIER` or `ONEZONE_IA` storage classes after a threshold, and aborting incomplete multipart uploads after 7 days. Rules are evaluated daily and take effect within 24 hours.

**Object Lock** -- Set `objectLockEnabled` to true at creation time for WORM (Write Once Read Many) compliance. Requires versioning to be enabled. Cannot be added to an existing bucket or removed once set.

**CORS** -- The `corsRules` array controls which web origins can make cross-origin requests to the bucket's S3 endpoint. Required when browser-based JavaScript uploads or reads objects directly. Without CORS rules, browsers block cross-origin requests.

**Force destroy** -- Set `forceDestroy` to true for dev/test buckets to allow clean teardown even when objects exist. Keep as false (default) for production to prevent accidental data loss.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `bucket_id` | Regional ID of the bucket (`{region}/{bucket-name}`) | Terraform state references, IAM policy targets |
| `endpoint` | FQDN endpoint (`{name}.s3.{region}.scw.cloud`) | S3 client configuration, CDN origin, application config |
| `api_endpoint` | S3 API endpoint (`https://s3.{region}.scw.cloud`) | AWS CLI `--endpoint-url`, SDK configuration |
| `bucket_name` | Name of the bucket in Scaleway Object Storage | S3 client bucket parameter, CI/CD pipeline variables |
| `region` | Region where the bucket is deployed | Co-location decisions for serverless functions and CDN |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private bucket** -- A minimal bucket with no versioning, no lifecycle rules, and no CORS. The fastest path to general-purpose file storage for uploads, documents, and application assets. Start from the **Private Bucket** preset.

**Versioned bucket with lifecycle** -- A bucket with versioning enabled and a 90-day Glacier transition for cost-optimized archival. Incomplete multipart uploads are cleaned up after 7 days. Standard production configuration for data requiring version history and compliance. Start from the **Versioned Lifecycle** preset.

## Works With

This component operates independently and does not reference other deployment components.