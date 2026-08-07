# Scaleway Serverless Function

Deploys a serverless function on Scaleway as a composite resource that bundles a function namespace, the function itself, and optional cron triggers into a single declarative unit. Supports multiple language runtimes (Node.js, Python, Go, Rust, PHP), scale-to-zero autoscaling, Private Network connectivity for VPC-internal access, environment variables and secrets, optional zip-based code deployment, and scheduled invocations. Supports ValueFromRef for Private Network dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Function Namespace** -- a Scaleway namespace that groups the function for lifecycle management and isolation
- **Serverless Function** -- the deployed function with the configured runtime, handler, privacy, memory limit, scaling bounds, and environment variables
- **Cron Triggers** -- created only when `cronTriggers` is populated; each trigger invokes the function on a CRON schedule with JSON arguments
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **Function source code** deployed via the Scaleway CLI (`scw function deploy`), CI/CD pipeline, or the `zipFile` and `zipHash` spec fields for IaC-managed code uploads.
- **A Scaleway Private Network** (optional) in the target region when the function needs to access databases, caches, or other Private Network resources without traversing the public internet.

## Deploy

### Console

Open the deployment store, find **Scaleway Serverless Function**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **HTTP API** preset in the [Presets](#presets) tab for a public Node.js function that scales from zero.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayServerlessFunction
metadata:
  name: api-handler
  org: acme-corp
  env: prod
spec:
  region: fr-par
  runtime: node20
  handler: handler.handle
  privacy: privacy_public
```

```shell
planton apply -f scaleway-serverless-function.yaml
```

This creates a public Node.js 20 serverless function with default 256 MB memory, scale-to-zero behavior, and a maximum of 20 instances. No Private Network, code upload, or cron triggers are configured. Deploy code separately via the Scaleway CLI.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the function to a Private Network deployed in the same InfraPipeline:

```yaml
spec:
  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id
```

The InfraPipeline resolves the dependency graph, deploys the Private Network first, then provisions the function with the resolved Private Network ID.

## Key Configuration

These are the most important decisions when configuring a serverless function. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Runtime** -- The `runtime` field selects the language runtime (e.g., `node20`, `python311`, `go122`, `rust165`, `php82`). This is a plain string -- Scaleway adds new runtimes frequently. The runtime determines how the `handler` entry point is resolved.

**Privacy** -- The `privacy` field controls endpoint authentication. Use `privacy_public` for internet-accessible APIs and webhooks. Use `privacy_private` for functions invoked only by cron triggers or internal Scaleway services.

**Scaling** -- Set `minScale` to 0 for scale-to-zero (no cost when idle) or 1+ for always-warm instances that eliminate cold starts. The `maxScale` field caps concurrent instances -- set to 1 for scheduled jobs that should not run concurrently.

**Code deployment** -- Set `zipFile` and `zipHash` to upload function code via IaC. When `zipHash` changes, the module re-uploads and redeploys the function. Alternatively, deploy code separately via the Scaleway CLI or CI/CD pipeline and leave these fields empty.

**Cron triggers** -- Add `cronTriggers` for scheduled function invocations. Each trigger specifies a CRON expression and JSON arguments passed to the function's event object. Ideal for background jobs, data cleanup, and periodic synchronization.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayPrivateNetwork** (optional) | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `function_id` | Unique identifier of the deployed function | Scaleway CLI operations, API management, Terraform import |
| `namespace_id` | Unique identifier of the function namespace | External cron triggers, tokens, additional functions |
| `domain_name` | Native Scaleway HTTPS endpoint for the function | ScalewayDnsRecord CNAME targets for custom domain routing |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTP API** -- A public Node.js 20 function that scales from zero to 20 instances with 256 MB memory and a 5-minute timeout. The standard configuration for lightweight REST endpoints and webhooks. Start from the **HTTP API** preset.

**Scheduled job** -- A private Python 3.11 function with a daily 2 AM UTC cron trigger and a maximum of 1 instance. The standard configuration for background tasks like data cleanup, report generation, and cache warming. Start from the **Scheduled Job** preset.

## Works With

- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides VPC-internal connectivity for accessing databases, caches, and other private resources