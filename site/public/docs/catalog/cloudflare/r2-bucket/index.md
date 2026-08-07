---
title: "R2 Bucket"
description: "R2 Bucket deployment documentation"
icon: "package"
order: 100
componentName: "cloudflarer2bucket"
---

# R2 Bucket on Cloudflare

Deploys an R2 object storage bucket on Cloudflare with configurable location hints, public access controls, and optional custom domain serving. R2 provides S3-compatible storage with zero egress fees. Integrates with Planton's Provider Connections for Cloudflare credential management and supports ValueFromRef wiring to DNS zones for custom domain configuration.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **R2 Bucket** -- an object storage bucket in the specified Cloudflare account with the configured location hint (auto, WNAM, ENAM, WEUR, EEUR, APAC, or OC)
- **R2 Custom Domain** -- created only when `customDomain.enabled` is `true`; configures a custom hostname (e.g., media.example.com) for bucket access using the referenced Cloudflare zone
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has R2 permissions. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A Cloudflare account** with R2 enabled. The `accountId` field identifies which account owns the bucket.
- **A Cloudflare DNS zone** (optional) -- required only when using a custom domain. Provide the zone ID directly or reference a CloudflareDnsZone Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **R2 Bucket on Cloudflare**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Private** preset in the [Presets](#presets) tab to pre-populate a secure default configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareR2Bucket
metadata:
  name: app-assets
  org: acme-corp
  env: prod
spec:
  bucketName: app-assets-prod
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  location: auto
  publicAccess: false
```

```shell
planton apply -f r2-bucket.yaml
```

This creates a private R2 bucket with automatic location selection. No public URL or custom domain is configured. Access is limited to Workers, API tokens, or the Cloudflare dashboard. A Stack Job tracks the provisioning in real time.

### InfraChart

When using a custom domain, use ValueFromRef to wire the bucket to a DNS zone deployed in the same InfraPipeline:

```yaml
spec:
  customDomain:
    enabled: true
    zoneId:
      valueFrom:
        kind: CloudflareDnsZone
        name: example-zone
        fieldPath: status.outputs.zone_id
    domain: media.example.com
```

The InfraPipeline resolves the dependency graph, deploys the DNS zone first, then provisions the R2 bucket with the custom domain configured on the resolved zone.

## Key Configuration

These are the most important decisions when configuring an R2 bucket. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Location hint** -- The `location` field specifies where Cloudflare stores bucket data. Use `auto` (default) to let Cloudflare choose the optimal region, or select a specific region (`WNAM`, `ENAM`, `WEUR`, `EEUR`, `APAC`, `OC`) to pin data to a geographic area for latency or compliance requirements. The location is immutable after bucket creation.

**Public access** -- Set `publicAccess: false` (default) to restrict access to Workers, API tokens, and the Cloudflare dashboard. Set `publicAccess: true` to enable the r2.dev public URL for direct HTTP access. Even for public content, consider using a custom domain for branded URLs.

**Custom domain** -- Enable `customDomain` to serve bucket content from a branded hostname (e.g., `media.example.com`). Requires a Cloudflare DNS zone for the domain. The custom domain provides a clean URL and inherits Cloudflare's CDN and caching benefits.

**Bucket naming** -- The `bucketName` must be DNS-compatible: 3-63 characters, lowercase alphanumeric and hyphens, starting and ending with alphanumeric. Choose a name that reflects the environment and purpose (e.g., `app-assets-prod`, `backup-staging`).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CloudflareDnsZone** (optional) | `customDomain.zoneId` | `status.outputs.zone_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `bucket_name` | The name of the R2 bucket | Worker script bundle references, application configuration |
| `bucket_url` | The accessible bucket URL (R2 public endpoint or S3 API URL) | Application asset URLs, SDK endpoint configuration |
| `custom_domain_url` | The custom domain URL if configured (e.g., https://media.example.com) | Frontend asset references, public download URLs |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private bucket** -- No public access, automatic location. Access via Workers or API tokens only. Suitable for backups, logs, CI/CD artifacts, and internal data. Start from the **Private** preset.

**Public CDN bucket** -- Public access enabled with a custom domain for branded content delivery. Combines R2's zero-egress pricing with Cloudflare's CDN for static assets, images, and media files. Start from the **Public CDN** preset.

## Works With

- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone) -- provides the zone ID for custom domain configuration