---
title: "Worker"
description: "Worker deployment documentation"
icon: "package"
order: 100
componentName: "cloudflareworker"
---

# Worker on Cloudflare

Deploys a Cloudflare Worker with script bundle loading from R2, optional KV namespace bindings, custom domain routing, and environment variable and secret injection. Integrates with Planton's Provider Connections for Cloudflare credential management and supports ValueFromRef wiring to KV namespaces for cross-resource dependency resolution.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Workers Script** -- the serverless function deployed to Cloudflare's edge network, loaded from a pre-built script bundle in an R2 bucket, with module syntax and Node.js compatibility enabled
- **KV Namespace Bindings** -- created only when `kvBindings` entries are provided; binds one or more Workers KV namespaces to the script for edge key-value storage access
- **Plain-Text Variable Bindings** -- created only when `env.variables` entries are provided; injects non-sensitive configuration as Worker bindings
- **Secrets** -- created only when `env.secrets` entries are provided; uploaded via the Cloudflare Secrets API, encrypted at rest, and never logged
- **DNS Record** -- created only when `dns.enabled` is `true`; a proxied A record pointing the custom hostname to Cloudflare's edge
- **Workers Route** -- created only when `dns.enabled` is `true`; attaches the Worker to a URL pattern on the custom domain
- **Observability** -- Workers Logs enabled by default with 100% head sampling rate for full request visibility

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Workers and R2 permissions. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A pre-built script bundle in R2** -- the Worker code must be compiled and uploaded to an R2 bucket before deployment. The `scriptBundle` field specifies the bucket name and object path (e.g., `dist/worker.js`).
- **KV namespaces** (optional) -- if the Worker needs KV storage, create CloudflareKvNamespace resources first. Provide namespace IDs directly or reference them via ValueFromRef.
- **A Cloudflare DNS zone** (optional) -- required only when using custom domain routing (`dns.enabled: true`). Provide the zone ID directly in the `dns.zoneId` field.

## Deploy

### Console

Open the deployment store, find **Worker on Cloudflare**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Minimal** preset in the [Presets](#presets) tab for a bare Worker deployment.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareWorker
metadata:
  name: api-gateway
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  workerName: api-gateway
  scriptBundle:
    bucket: deploy-bundles
    path: api-gateway/v1.0.0/worker.js
```

```shell
planton apply -f cloudflare-worker.yaml
```

This creates a Worker with the script loaded from R2. No KV bindings, DNS routes, or environment variables are configured. The Worker runs but has no route -- attach routes via the Cloudflare dashboard or a subsequent update. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying with KV namespace dependencies, use ValueFromRef to wire the Worker to KV namespaces deployed in the same InfraPipeline:

```yaml
spec:
  kvBindings:
    - valueFrom:
        kind: CloudflareKvNamespace
        name: session-cache
        fieldPath: status.outputs.namespace_id
```

The InfraPipeline resolves the dependency graph, deploys the KV namespace first, then provisions the Worker with the resolved namespace ID bound as a KV binding.

## Key Configuration

These are the most important decisions when configuring a Worker. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Script bundle source** -- The `scriptBundle` field references a pre-built Worker bundle stored in an R2 bucket. Build and upload your Worker code (e.g., via Wrangler or CI/CD) before deploying. The module uses ES module syntax with `index.js` as the main module entry point.

**KV bindings** -- Each entry in `kvBindings` attaches a Workers KV namespace to the script. Use ValueFromRef to reference CloudflareKvNamespace resources, or provide literal namespace IDs. Multiple KV namespaces can be bound to a single Worker for different data stores (cache, sessions, configuration).

**Custom domain routing** -- Set `dns.enabled: true` with a `dns.zoneId` and `dns.hostname` to create a proxied DNS record and Workers Route for the custom domain. The `dns.routePattern` defaults to `hostname/*` if omitted. Omit the `dns` block entirely for Workers that run via Cron Triggers, Queues, or workers.dev subdomain only.

**Environment variables and secrets** -- Use `env.variables` for non-sensitive configuration (plain-text bindings) and `env.secrets` for sensitive values (encrypted at rest via Cloudflare Secrets API). Both support Planton's `$variables-group/` and `$secrets-group/` reference syntax for centralized configuration management.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CloudflareKvNamespace** (optional) | `kvBindings` | `status.outputs.namespace_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `script_id` | The Cloudflare-assigned identifier of the deployed Worker script | Worker management, monitoring dashboards |
| `route_urls` | The route URL patterns where this Worker is active | DNS verification, application endpoint configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**API with custom domain** -- A full-featured Worker with KV bindings, custom domain DNS routing, and environment variables. Use for production REST or GraphQL APIs that need edge storage and a branded hostname. Start from the **API with Custom Domain** preset.

**Minimal worker** -- A bare Worker with only the script bundle. No KV bindings, DNS routes, or environment variables. Use for initial deployments, Cron Trigger workers, or when routes are managed separately. Start from the **Minimal** preset.

## Works With

- [**KV Namespace on Cloudflare**](/cloud-catalog/cloudflare-kv-namespace) -- provides KV namespace IDs for Worker storage bindings