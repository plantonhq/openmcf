---
title: "Pages Project"
description: "Pages Project deployment documentation"
icon: "package"
order: 100
componentName: "cloudflarepagesproject"
---

# Pages Project on Cloudflare

Provisions a Cloudflare Pages project: a managed site host that builds and serves a static site or full-stack app (static assets plus Pages Functions) from Cloudflare's edge. This manages the durable **project** -- its build configuration, optional git connection, per-environment runtime configuration (bindings, env vars, compatibility), and custom domains. Deployments themselves are produced out-of-band (a git push for connected projects, or `wrangler pages deploy` for direct upload). Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Pages Project** -- the managed site project and its `*.pages.dev` subdomain
- **Build configuration** -- how Cloudflare builds the site (for git-connected and wrangler builds)
- **Per-environment runtime config** -- bindings, variables, and secrets for preview and production
- **Custom domains** -- any hostnames you attach (Cloudflare provisions their certificates)

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Pages edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **An account** -- the project is created under an existing Cloudflare account.
- **Git authorization (git-connected projects only)** -- the account must already be authorized with the git provider (the GitHub App install or GitLab OAuth), a one-time browser step in the Cloudflare dashboard.

## Deploy

### Console

Open the deployment store, find **Pages Project on Cloudflare**, and click **Deploy**. The creation wizard captures the account, project name, build configuration, optional git source, the production and preview deployment configs (bindings, variables, secrets), and custom domains.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflarePagesProject
metadata:
  name: marketing-site
  org: acme-corp
  env: prod
spec:
  accountId: a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4
  name: marketing-site
  productionBranch: main
  buildConfig:
    buildCommand: npm run build
    destinationDir: dist
  source:
    type: github
    config:
      owner: acme-corp
      repoName: marketing-site
  deploymentConfigs:
    production:
      kvNamespaces:
        - name: CACHE
          namespaceId:
            valueFrom:
              kind: CloudflareKvNamespace
              name: prod-cache
              fieldPath: status.outputs.namespace_id
  domains:
    - www.example.com
```

```shell
planton apply -f cloudflare-pages-project.yaml
```

This creates a git-connected project that builds on every push to `main`, binds a KV namespace to production Functions, and serves on a custom domain. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Pages project. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Account (`accountId`)** and **Name (`name`)** -- both **immutable**; changing either replaces the project. The name is also the `<name>.pages.dev` subdomain label.

**Build (`buildConfig`)** -- the build command and output directory. Omit for a pre-built direct-upload project.

**Git Source (`source`)** -- connect a GitHub or GitLab repository so Cloudflare builds on every push, or leave it empty and deploy with `wrangler pages deploy`.

**Deployment Configs (`deploymentConfigs.production` / `.preview`)** -- per-environment runtime config: compatibility, environment variables, secrets, and bindings to KV namespaces, D1 databases, R2 buckets, queues, Hyperdrive configs, and other Workers. Bind preview to non-production resources to isolate test traffic.

**Domains (`domains`)** -- custom hostnames in zones on this account; the `*.pages.dev` subdomain is always available.

## Outputs and Dependencies

### What This Component Consumes

Bindings can reference a **CloudflareKvNamespace**, **CloudflareD1Database**, **CloudflareR2Bucket**, **CloudflareQueue**, **CloudflareHyperdriveConfig**, and **CloudflareWorker**. Secrets reference managed secrets resolved just-in-time at deploy.

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `project_name` | The project name (echoed) | Referenced by downstream resources |
| `subdomain` | The `*.pages.dev` subdomain | Accessing the deployed site |
| `domains` | The attached custom domains | DNS and routing |
| `created_on` | Creation timestamp | Auditing |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Git-connected site** -- Cloudflare builds and deploys on every push to the production branch.

**Direct upload** -- deploy pre-built output with `wrangler pages deploy`; no build config or git connection.

## Works With

- [**KV Namespace on Cloudflare**](/cloud-catalog/cloudflare-kv-namespace) -- bound to Functions for edge key-value storage
- [**D1 Database on Cloudflare**](/cloud-catalog/cloudflare-d1-database) -- bound to Functions for serverless SQL
- [**R2 Bucket on Cloudflare**](/cloud-catalog/cloudflare-r2-bucket) -- bound to Functions for object storage
- [**Queue on Cloudflare**](/cloud-catalog/cloudflare-queue) -- bound to Functions as a producer
- [**Hyperdrive Config on Cloudflare**](/cloud-catalog/cloudflare-hyperdrive-config) -- bound to Functions for pooled SQL access
- [**Worker on Cloudflare**](/cloud-catalog/cloudflare-worker) -- bound to Functions for service-to-service calls
