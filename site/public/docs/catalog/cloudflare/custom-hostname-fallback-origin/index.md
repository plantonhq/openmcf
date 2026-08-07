---
title: "Custom Hostname Fallback Origin"
description: "Custom Hostname Fallback Origin deployment documentation"
icon: "package"
order: 100
componentName: "cloudflarecustomhostnamefallbackorigin"
---

# Custom Hostname Fallback Origin on Cloudflare

Sets the default origin for a Cloudflare-for-SaaS zone: the backend that all of the zone's custom hostnames route to unless an individual hostname overrides it. It is a zone-level singleton (one per zone) and a prerequisite for serving traffic to any custom hostname in the zone. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Fallback Origin** -- the zone's default backend for custom-hostname traffic

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has SSL and Certificates edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A SaaS zone** -- the fallback origin is set on an existing `CloudflareDnsZone` configured for Cloudflare for SaaS.

## Deploy

### Console

Open the deployment store, find **Custom Hostname Fallback Origin on Cloudflare**, and click **Deploy**. The creation wizard captures the zone and the default backend origin.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareCustomHostnameFallbackOrigin
metadata:
  name: saas-fallback
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: saas-zone
      fieldPath: status.outputs.zone_id
  origin:
    value: origin.helpdesk.io
```

```shell
planton apply -f cloudflare-custom-hostname-fallback-origin.yaml
```

This points every custom hostname in the SaaS zone at `origin.helpdesk.io` by default. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a fallback origin. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Zone (`zoneId`)** -- The SaaS zone this fallback origin belongs to. Immutable -- changing it replaces the resource. Reference a `CloudflareDnsZone` to keep the dependency in the graph.

**Origin (`origin`)** -- The default backend hostname custom hostnames route to. A literal hostname within the SaaS zone, or a reference to a backend endpoint output. Editable.

## Outputs and Dependencies

### What This Component Consumes

The fallback origin references a **CloudflareDnsZone** (via `zoneId`).

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `status` | The deployment status of the fallback origin | Monitoring activation |
| `created_at` | RFC3339 creation timestamp | Auditing |
| `updated_at` | RFC3339 last-updated timestamp | Auditing |
| `errors` | Any errors reported while deploying | Troubleshooting |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single shared backend** -- route all customer hostnames to one origin, overriding per-hostname only when needed.

## Works With

- [**Custom Hostname on Cloudflare**](/cloud-catalog/cloudflare-custom-hostname) -- the per-customer hostnames that route to this fallback origin
- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone) -- the SaaS zone this origin serves
