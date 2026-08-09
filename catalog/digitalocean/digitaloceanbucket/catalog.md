# Storage Bucket on DigitalOcean

Deploys a DigitalOcean Spaces object storage bucket with configurable access control, optional versioning, and tagging. Spaces buckets provide S3-compatible storage for application assets, backups, static websites, and CDN origins. Integrates with Planton's Provider Connections for DigitalOcean API token management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Spaces Bucket** -- an object storage bucket in the specified DigitalOcean region with the configured access control (private or public-read) and optional versioning
- **Versioning Configuration** -- created only when `versioningEnabled` is true; preserves previous object versions for recovery and audit

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A DigitalOcean account with Spaces enabled** -- Spaces must be activated in your account. Bucket names are globally unique across all Spaces customers.
- **A target region that supports Spaces** -- not all DigitalOcean regions support Spaces. Common Spaces regions include `nyc3`, `sfo3`, `ams3`, `sgp1`, and `fra1`.

## Deploy

### Console

Open the deployment store, find **Storage Bucket on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Private** preset in the [Presets](#presets) tab to create a versioned private bucket.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanBucket
metadata:
  name: app-assets
  org: acme-corp
  env: prod
spec:
  bucketName: acme-app-assets
  region: nyc3
```

```shell
planton apply -f do-bucket.yaml
```

This creates a private Spaces bucket with no versioning. Objects are accessible only via authenticated API calls or signed URLs. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Spaces bucket. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Access control** -- The `accessControl` field sets the bucket ACL: `PRIVATE` (default) restricts access to authenticated API calls and signed URLs, while `PUBLIC_READ` makes objects accessible to anyone with the URL. Use public-read for static websites and CDN origins; use private for backups, logs, and sensitive data.

**Versioning** -- Set `versioningEnabled: true` to preserve previous object versions on overwrite or delete. Versioning protects against accidental data loss and supports compliance requirements. Disabled by default to minimize storage costs.

**Bucket naming** -- The `bucketName` must be DNS-compatible (lowercase, hyphens, 3-63 characters) and globally unique across all DigitalOcean Spaces customers. Choose a name that includes your organization or project prefix to avoid conflicts.

**Region selection** -- The `region` field determines where data is stored. Choose the region closest to your application or users. Spaces endpoints are region-specific (e.g., `https://nyc3.digitaloceanspaces.com`).

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `bucket_id` | Unique bucket identifier (UUID) | API operations, lifecycle policies |
| `endpoint` | Regional endpoint URL for the bucket | Application configuration, CDN origin setup |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private bucket with versioning** -- private access with versioning enabled for backups, logs, and sensitive application data. Objects are accessible only via authenticated API calls. Start from the **Private** preset.

**Public static website bucket** -- public-read access with versioning disabled for hosting static sites, CDN-served assets, and documentation. Suitable for JAMstack deployments and media hosting. Start from the **Public Static Website** preset.

## Works With

This component operates independently and does not reference other components.