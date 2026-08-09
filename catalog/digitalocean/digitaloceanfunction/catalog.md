# Function on DigitalOcean

Deploys a serverless function on DigitalOcean via App Platform with configurable runtime, memory, timeout, environment variables, and optional cron scheduling. Functions can serve as HTTP endpoints or run as scheduled background jobs. Integrates with Planton's Provider Connections for DigitalOcean API token management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **App Platform Application** -- a managed application container that hosts the serverless function in the specified DigitalOcean region
- **Function Component** -- the serverless function within the App Platform app, configured with the specified runtime, memory, and timeout
- **Environment Variables** -- regular and encrypted secret environment variables injected into the function runtime

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A GitHub repository** containing the function source code. The `githubSource` configuration specifies the repository, branch, and optional auto-deploy on push.
- **A source directory** within the repository that contains the function code and `project.yml` file required by DigitalOcean Functions.
- **A supported runtime** -- Node.js 18/20, Python 3.9/3.10/3.11, Go 1.20/1.21, or PHP 8.2.

## Deploy

### Console

Open the deployment store, find **Function on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web API** preset in the [Presets](#presets) tab to deploy an HTTP-accessible function from GitHub.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanFunction
metadata:
  name: api-handler
  org: acme-corp
  env: prod
spec:
  functionName: api-handler
  region: nyc1
  runtime: nodejs_20
  githubSource:
    repo: "acme-corp/functions"
    branch: main
    deployOnPush: true
  sourceDirectory: /functions/api
  memoryMb: 256
  timeoutMs: 3000
  isWeb: true
```

```shell
planton apply -f do-function.yaml
```

This creates a Node.js 20 serverless function exposed as an HTTP endpoint with 256 MB memory, a 3-second timeout, and automatic redeployment on push. No environment variables or cron schedule are configured. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a serverless function. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Runtime selection** -- The `runtime` field sets the language and version (e.g., `nodejs_20`, `python_311`, `go_121`). Choose the runtime that matches your function code. Each runtime has its own entry point conventions and dependency management.

**Memory and timeout** -- The `memoryMb` field sets the memory allocation (128, 256, 512, 1024, or 2048 MB). The `timeoutMs` field sets the maximum execution time (up to 300,000 ms / 5 minutes). Default is 256 MB and 3,000 ms, sufficient for most API handlers. Increase both for data-processing workloads.

**Web vs. scheduled execution** -- Set `isWeb: true` (default) to expose the function as an HTTP endpoint. Set `isWeb: false` with a `cronSchedule` (standard cron expression) to run the function on a timer. A function cannot be both web-accessible and cron-triggered simultaneously.

**Auto-deploy** -- Set `githubSource.deployOnPush: true` to automatically redeploy when changes are pushed to the configured branch. Disable for production environments where you want manual deployment control.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `function_id` | Unique identifier of the deployed function | API operations, monitoring dashboards |
| `https_endpoint` | Public HTTPS URL for invoking the function | Webhook targets, API gateway routing, application integration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web API function** -- a Node.js 20 HTTP endpoint deployed from GitHub with auto-deploy, 256 MB memory, and a 3-second timeout. Suitable for API handlers, webhooks, and lightweight HTTP services. Start from the **Web API** preset.

**Scheduled background job** -- a Python 3.11 function triggered hourly by cron, not exposed as an HTTP endpoint, with 512 MB memory and a 60-second timeout. Suitable for ETL pipelines, data synchronization, and periodic report generation. Start from the **Scheduled Job** preset.

## Works With

This component operates independently and does not reference other components.