# KV Namespace on Cloudflare

Deploys a Workers KV namespace on Cloudflare for globally replicated, eventually consistent key-value storage at the edge. KV namespaces are bound to Workers for low-latency reads across Cloudflare's network. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Workers KV Namespace** -- a named key-value store in the Cloudflare account, identified by a unique namespace ID that Workers reference as a binding
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Workers KV permissions. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A Cloudflare account** with Workers KV enabled. The KV namespace is created at the account level and can be bound to any Worker in that account.

## Deploy

### Console

Open the deployment store, find **KV Namespace on Cloudflare**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareKvNamespace
metadata:
  name: session-cache
  org: acme-corp
  env: prod
spec:
  namespaceName: session-cache-prod
  ttlSeconds: 3600
  description: "Session cache for production API workers"
```

```shell
planton apply -f cloudflare-kv-namespace.yaml
```

This creates a Workers KV namespace with a 1-hour default TTL for keys. The namespace ID is exported in stack outputs for binding to CloudflareWorker resources. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a KV namespace. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Namespace name** -- The `namespaceName` must be unique within the Cloudflare account, up to 64 characters. Choose a name that reflects the purpose and environment (e.g., `app-cache-prod`, `config-store-staging`). The name is visible in the Cloudflare dashboard.

**Default TTL** -- Set `ttlSeconds: 0` (default) for keys that never expire. Set a positive value (minimum 60 seconds) to auto-expire keys after the specified duration. Individual key writes can override this default. Use expiring TTLs for cache data and session stores; use non-expiring for configuration and feature flags.

**Description** -- The optional `description` field documents the namespace's purpose. Helpful when managing multiple namespaces across environments and services.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace_id` | The unique identifier of the created KV namespace | CloudflareWorker `kvBindings` for binding KV storage to Worker scripts |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard namespace** -- A KV namespace with configurable TTL and description. The only pattern needed -- KV namespaces are simple resources whose complexity lives in the Workers that consume them. Start from the **Standard** preset.

## Works With

This component operates independently and does not reference other components.