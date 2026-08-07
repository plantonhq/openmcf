---
title: "Storage Bucket"
description: "Storage Bucket deployment documentation"
icon: "package"
order: 100
componentName: "civobucket"
---

# Storage Bucket on Civo

Deploys an S3-compatible object storage bucket on Civo Cloud with auto-generated access credentials, configurable versioning, and regional placement. Civo Object Storage provides durable, scalable storage accessible via the standard S3 API. Integrates with Planton's Provider Connections for Civo credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo Object Store Credential** -- an access key ID and secret key pair generated automatically for bucket access via the S3-compatible API
- **Civo Object Store Bucket** -- an S3-compatible storage bucket in the specified region, linked to the generated credential
- **Civo Tags** -- metadata tags applied to the bucket for organizational tracking

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A Civo account** with Object Storage enabled in the target region. No additional prerequisites are required -- buckets are standalone resources with no VPC or firewall dependencies.

## Deploy

### Console

Open the deployment store, find **Storage Bucket on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab for a general-purpose storage bucket.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoBucket
metadata:
  name: app-assets
  org: acme-corp
  env: prod
spec:
  bucketName: app-assets
  region: lon1
```

```shell
planton apply -f civo-bucket.yaml
```

This creates an object storage bucket in Civo's London region with versioning disabled and auto-generated access credentials. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Civo storage bucket. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Bucket name** -- The `bucketName` field must be DNS-compatible (3-63 lowercase characters, letters, digits, and hyphens). Choose a name that reflects the content or purpose of the bucket, as renaming requires recreating the resource.

**Versioning** -- Set `versioningEnabled` to `true` to retain all previous versions of every object. Enabled versioning protects against accidental overwrites and deletions but increases storage costs. Leave disabled (default) for general-purpose storage where version history is unnecessary.

**Region** -- The `region` field accepts a Civo region code (e.g., `lon1` for London, `nyc1` for New York). Choose the region closest to your compute workloads to minimize data transfer latency.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `bucket_id` | Unique identifier of the bucket on Civo | Civo API operations, lifecycle management |
| `endpoint_url` | S3-compatible endpoint URL for the bucket | Application storage configuration, SDK/CLI access |
| `access_key_secret_ref` | Reference to the secret storing the access key ID | Application credential injection via environment variables |
| `secret_key_secret_ref` | Reference to the secret storing the secret access key | Application credential injection via environment variables |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard bucket** -- a general-purpose bucket with versioning disabled for application assets, static files, and data exports. Storage costs stay predictable since overwrites replace objects in place. Start from the **Standard** preset.

**Versioned backup bucket** -- a bucket with versioning enabled for database backups, compliance archives, and configuration storage where rollback capability is critical. Start from the **Versioned Backup** preset.

## Works With

This component operates independently and does not reference other deployment components.