# D1 Database on Cloudflare

Deploys a Cloudflare D1 serverless SQLite database with configurable region placement and optional read replication. Integrates with Planton's Provider Connections for Cloudflare credential management and exports the database identifier for binding to Cloudflare Workers.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **D1 Database** -- a serverless SQLite database created in the specified Cloudflare account, with an optional primary location hint to control the region of the primary instance
- **Read Replication** -- created only when `readReplication` is configured; enables D1 Read Replication (Beta) to place read-only replicas across multiple regions for lower global read latency
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has D1 permissions. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A Cloudflare account** with D1 access enabled. The `accountId` field identifies which Cloudflare account owns the database.
- **Schema management** -- D1 tables, indexes, and migrations are managed via the Wrangler CLI, not at the resource level. Deploy the database first, then run migrations against it.

## Deploy

### Console

Open the deployment store, find **D1 Database on Cloudflare**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareD1Database
metadata:
  name: app-cache
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  databaseName: app-cache
```

```shell
planton apply -f cloudflare-d1-database.yaml
```

This creates a D1 database named `app-cache` with Cloudflare selecting the default storage region. No read replication is configured. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a D1 database. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Region placement** -- The `region` field sets the primary location hint for the database. Valid values are `weur` (Western Europe), `eeur` (Eastern Europe), `apac` (Asia Pacific), `oc` (Oceania), `wnam` (Western North America), and `enam` (Eastern North America). Omit to let Cloudflare select based on your account settings. Choose a region close to your Workers for lowest write latency.

**Read replication** -- Set `readReplication.mode` to `auto` to enable automatic read replicas across multiple regions. This reduces global read latency but requires your application code to use the D1 Sessions API for consistency. Set to `disabled` or omit for single-region operation.

**Database naming** -- The `databaseName` must be unique within the account and is limited to 64 characters. Choose a descriptive name since it appears in Wrangler CLI commands and Worker bindings.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `database_id` | The unique identifier (UUID) of the created D1 database | CloudflareWorker D1 bindings, Wrangler CLI configuration |
| `database_name` | The name of the database as confirmed by Cloudflare | Application configuration, monitoring dashboards |
| `connection_string` | Reserved for future use; currently empty as the Pulumi Cloudflare provider does not expose a D1 connection string | Future Worker binding configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard database** -- A D1 database with optional region placement and no read replication. Use for edge databases backing Workers, lightweight relational storage, or key-value workloads with SQL access. Start from the **Standard** preset.

## Works With

This component operates independently and does not reference other components.