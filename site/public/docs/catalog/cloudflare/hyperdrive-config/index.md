---
title: "Hyperdrive Config"
description: "Hyperdrive Config deployment documentation"
icon: "package"
order: 100
componentName: "cloudflarehyperdriveconfig"
---

# Hyperdrive Config on Cloudflare

Deploys a Cloudflare Hyperdrive -- a connection pooler and global query cache that lets a Worker reach a regional SQL database (PostgreSQL or MySQL) with low latency. A Worker binds to this config and queries the origin without paying the full connection-setup round trip on every request; Hyperdrive reuses warm, pooled connections and can cache read-query results at the edge. Hyperdrive configs are account-scoped and integrate with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Hyperdrive Config** -- an account-scoped configuration pointing at your origin database, with pooling and caching behavior
- **Origin Credentials** -- the database password (and optional Cloudflare Access service-token secret) resolved just-in-time from managed secrets at deploy
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Hyperdrive edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Managed secret for the database password** -- store the origin password as an org secret and reference it; Hyperdrive never accepts a plaintext password.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Cloudflare Account

- **Network reachability** -- the origin must be reachable from Cloudflare's network: a public endpoint, or a private one fronted by Cloudflare Access / a Cloudflare Tunnel (use the Access client ID + secret in that case).

## Deploy

### Console

Open the deployment store, find **Hyperdrive Config on Cloudflare**, and click **Deploy**. The creation wizard captures the owning account and name, the required origin connection (engine, host, port, database, user, and the reference-only password secret), and optional caching, tuning, and mTLS settings.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareHyperdriveConfig
metadata:
  name: prod-postgres
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  name: prod-postgres
  origin:
    scheme: postgres
    host: db.example.com
    database: app_production
    user: hyperdrive
    password:
      value: "$secret/prod-db-password"
  caching:
    maxAge: 60
```

```shell
planton apply -f cloudflare-hyperdrive-config.yaml
```

This creates a Hyperdrive config fronting a PostgreSQL origin with 60-second result caching. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Hyperdrive config. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Database Engine (`origin.scheme`)** -- `postgres`/`postgresql` (default port 5432) or `mysql` (default port 3306). Determines the wire protocol Hyperdrive uses.

**Password (`origin.password`)** -- The database user's password, provided as a managed-secret reference. The backend rejects plaintext; the runner resolves it just-in-time at deploy.

**Cloudflare Access (`origin.accessClientId` / `origin.accessClientSecret`)** -- Set these when the origin host is published behind Cloudflare Access, so Hyperdrive can authenticate through the Access application.

**Caching (`caching`)** -- `disabled`, `maxAge`, and `staleWhileRevalidate` control edge caching of read-query results. Enabled by default (60s / 15s) -- the single biggest latency win for read-heavy workloads.

**Origin Connection Limit (`originConnectionLimit`)** -- Caps pooled connections to the origin (5-100; defaults by plan).

## Outputs and Dependencies

### What This Component Consumes

This component has no Cloud Resource foreign-key dependencies; it points directly at an external origin database. Its credentials are managed-secret references resolved at deploy.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `hyperdrive_id` | The Cloudflare-assigned identifier of the config | Referenced by a CloudflareWorker's `hyperdrive` binding |
| `name` | The config name (echoed) | Verification, dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Serverless Postgres access** -- A Hyperdrive config fronting a managed PostgreSQL database (RDS, Neon, Supabase, Cloud SQL) so Workers query it without per-request connection overhead.

**Private origin via Access** -- A database published behind Cloudflare Access, reached using an Access service token (client ID + secret) rather than a public endpoint.

**Strict TLS** -- An mTLS configuration with `verify-full` for origins that require client certificates and full hostname verification.

## Works With

- [**Worker on Cloudflare**](/cloud-catalog/cloudflare-worker) -- binds this config (a `hyperdrive` binding) to query the origin database with pooled, cached connections
