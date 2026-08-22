# Worker on Cloudflare

Deploys a Cloudflare Worker — a script that runs on Cloudflare's edge — with its bindings, routing, schedules, and runtime settings. Bindings are grouped by type (the wrangler.toml grain) and each cross-resource binding accepts a literal or a reference to the producing resource, so a Worker composes as a real node in the resource graph.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Workers Script** -- the serverless function, from inline `content` or an R2 `r2Bundle`, with optional static `assets`
- **Bindings** -- typed lists (vars, secrets, KV, R2, D1, Hyperdrive, services, queues, Durable Objects, and the rest of the provider's binding types) flattened into the script's bindings array
- **workers.dev subdomain** -- created only when `workersDev.enabled` is true
- **Custom domains** -- one managed hostname (automatic TLS) per `customDomains` entry
- **Routes** -- one pattern-based route per `routes` entry
- **Cron trigger** -- created only when `schedules` is set
- **Observability** -- Workers Logs and traces when `observability` is set

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection with a Cloudflare API token that has Workers Scripts Write (and Routes / DNS Edit if you attach custom domains or routes).

### Cloudflare Account

- **A script source or assets** -- inline `content`, an R2 `r2Bundle`, a static `assets` directory, or a combination. A Worker with none of these is rejected.
- **Upstream resources** (optional) -- KV, D1, R2, Hyperdrive, another Worker, a Queue, or a Zero Trust tunnel, referenced via `valueFrom`.

## Deploy

### Console

Open the deployment store, find **Worker on Cloudflare**, and click **Deploy**. Start from the **Minimal** preset for an inline hello-world, or **API with Custom Domain** for a production-shaped Worker.

### CLI

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWorker
metadata:
  name: api-gateway
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  workerName: api-gateway
  content: |
    export default { async fetch() { return new Response("ok"); } }
  workersDev:
    enabled: true
```

```shell
planton apply -f cloudflare-worker.yaml
```

### InfraChart

Wire a Worker to a KV namespace deployed in the same pipeline:

```yaml
spec:
  kvNamespaces:
    - name: CONFIG
      namespaceId:
        valueFrom:
          kind: CloudflareKvNamespace
          name: session-cache
          fieldPath: status.outputs.namespace_id
```

## Key Configuration

**Script source** -- `content` for inline ES modules, `r2Bundle` for a CI-built artifact (`bucket` is a CloudflareR2Bucket reference or a literal name), `assets` for a static site or full-stack app.

**Bindings** -- grouped by type. Cross-resource fields take a literal or `valueFrom`. Secrets and secret-key material are sensitive.

**Routing** -- `workersDev` for `*.workers.dev`, `customDomains` for managed hostnames, `routes` for zone patterns.

**Migrations** -- Durable Object class create/rename/transfer/delete. Cloudflare treats the tag as a one-shot; a second apply of the same `newTag` is rejected.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CloudflareKvNamespace** (optional) | `kvNamespaces[].namespaceId` | `status.outputs.namespace_id` |
| **CloudflareD1Database** (optional) | `d1Databases[].databaseId` | `status.outputs.database_id` |
| **CloudflareR2Bucket** (optional) | `r2Buckets[].bucketName`, `r2Bundle.bucket` | `status.outputs.bucket_name` |
| **CloudflareHyperdriveConfig** (optional) | `hyperdriveConfigs[].configId` | `status.outputs.hyperdrive_id` |
| **CloudflareQueue** (optional) | `queues[].queueName` | `status.outputs.queue_name` |
| **CloudflareWorker** (optional) | `services[].service`, `durableObjects[].scriptName`, `tailConsumers[].service` | `status.outputs.script_name` |
| **CloudflareDnsZone** (optional) | `customDomains[].zoneId`, `routes[].zoneId` | `status.outputs.zone_id` |
| **CloudflareZeroTrustTunnel** (optional) | `vpcNetworks[].tunnelId` | `status.outputs.tunnel_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `script_id` | The Cloudflare-assigned identifier of the deployed Worker | Worker management |
| `script_name` | The script name | Service bindings and tail consumers on other Workers |
| `custom_domain_hostnames` | Managed hostnames attached to this Worker | DNS verification |
| `route_patterns` | Route patterns mapped to this Worker | Application endpoint configuration |

## Common Patterns

**API with custom domain** -- CI-built bundle in R2, KV + D1 by reference, managed custom domain. Start from the **API with Custom Domain** preset.

**Minimal worker** -- inline content, workers.dev only. Start from the **Minimal** preset.

**Static site / full-stack** -- `assets` alone, or `assets` plus a script. Start from the **Static Site** or **Full-stack App** presets.

## Works With

- [**KV Namespace on Cloudflare**](/cloud-catalog/cloudflare-kv-namespace)
- [**D1 Database on Cloudflare**](/cloud-catalog/cloudflare-d1-database)
- [**R2 Bucket on Cloudflare**](/cloud-catalog/cloudflare-r2-bucket)
- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone)
