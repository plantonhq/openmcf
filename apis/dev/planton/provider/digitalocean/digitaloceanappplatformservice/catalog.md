# Service on DigitalOcean AppPlatform

Deploys a containerized application on DigitalOcean App Platform as a web service, background worker, or one-off job. Supports two deployment sources -- building from a Git repository or running a pre-built container image from DigitalOcean Container Registry (DOCR). Integrates with Planton's Provider Connections for DigitalOcean credential management and ValueFromRef for DNS zone and registry dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **App Platform Application** -- a `digitalocean_app` resource containing a single service component configured according to the chosen `serviceType`
- **Web Service** -- created only when `serviceType` is `web_service`; receives external HTTP traffic with automatic HTTPS, supports CPU-based autoscaling with configurable min/max instance counts
- **Worker** -- created only when `serviceType` is `worker`; runs as a background process without HTTP ingress, supports fixed instance counts
- **Job** -- created only when `serviceType` is `job`; executes as a `PRE_DEPLOY` task for database migrations or one-off scripts
- **Custom Domain** -- created only when `customDomain` is specified; adds a `PRIMARY` domain entry to the app spec
- **Environment Variables** -- injected from the `env` map; scoped to `RUN_AND_BUILD_TIME` for web services and `RUN_TIME` for workers and jobs

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A Git repository** accessible to DigitalOcean App Platform (for `gitSource` deployments), or **a DigitalOcean Container Registry** with a pushed image (for `imageSource` deployments). Provide exactly one source per deployment.
- **A supported region** for App Platform (e.g., `nyc3`, `sfo3`, `fra1`, `sgp1`, `lon1`, `ams3`).

## Deploy

### Console

Open the deployment store, find **Service on DigitalOcean AppPlatform**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Git Source Web Service** preset in the [Presets](#presets) tab to deploy from a GitHub repository with minimal configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: my-web-app
  org: acme-corp
  env: prod
spec:
  serviceName: my-web-app
  region: nyc3
  serviceType: web_service
  gitSource:
    repoUrl: "https://github.com/acme-corp/backend.git"
    branch: main
  instanceSizeSlug: basic_xxs
  instanceCount: 1
```

```shell
planton apply -f app-service.yaml
```

This creates a single-instance web service in NYC3 built from the `main` branch using the smallest instance tier. No autoscaling or custom domain is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the app to a container registry and DNS zone deployed in the same InfraPipeline:

```yaml
spec:
  imageSource:
    registry:
      valueFrom:
        kind: DigitalOceanContainerRegistry
        name: prod-registry
        fieldPath: status.outputs.server_url
    repository: "myapp/backend"
    tag: "v1.0.0"
  customDomain:
    valueFrom:
      kind: DigitalOceanDnsZone
      name: example-zone
      fieldPath: spec.domain_name
```

The InfraPipeline resolves the dependency graph, deploys the container registry and DNS zone first, then provisions the App Platform service with the resolved values.

## Key Configuration

These are the most important decisions when configuring an App Platform service. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Service type** -- The `serviceType` field determines how the application runs. `web_service` receives external HTTP traffic with automatic HTTPS. `worker` runs background processes without an HTTP endpoint. `job` executes one-off or pre-deploy tasks like database migrations.

**Deployment source** -- Choose `gitSource` to let App Platform build from a repository (auto-detects language and framework), or `imageSource` to deploy a pre-built container image from DigitalOcean Container Registry. Only one source can be specified per deployment.

**Instance sizing and scaling** -- The `instanceSizeSlug` field controls CPU and memory per instance, ranging from `basic_xxs` for development to `professional_xl` for heavy production workloads. Enable `enableAutoscale` with `minInstanceCount` and `maxInstanceCount` to let App Platform scale based on CPU utilization (80% threshold). Autoscaling is only available for `web_service` types.

**Custom domain** -- Set `customDomain` to route traffic through your own domain instead of the default `ondigitalocean.app` hostname. Reference a DigitalOceanDnsZone resource via ValueFromRef to coordinate DNS and app deployment.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanDnsZone** (optional) | `customDomain` | `spec.domain_name` |
| **DigitalOceanContainerRegistry** (optional) | `imageSource.registry` | `status.outputs.server_url` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `app_id` | Unique identifier of the App Platform application | DigitalOcean API operations, monitoring dashboards |
| `default_hostname` | Default hostname assigned to the app (ending in `ondigitalocean.app`) | DNS CNAME targets, health check endpoints |
| `live_url` | Publicly accessible URL of the deployed service including protocol | Application integrations, load testing targets |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Git source web service** -- Deploys from a GitHub repository with App Platform's auto-build, using the smallest instance tier and a single instance. Ideal for getting started with source-to-production deployments. Start from the **Git Source Web Service** preset.

**Container image service** -- Deploys from a pre-built image in DigitalOcean Container Registry with professional-tier compute and autoscaling (2-5 instances). Suited for production workloads with external CI/CD pipelines. Start from the **Container Image** preset.

## Works With

- [**DNS Zone on DigitalOcean**](/cloud-catalog/digital-ocean-dns-zone) -- provides the DNS zone referenced by `customDomain` for routing traffic to the app
- [**Container Registry on DigitalOcean**](/cloud-catalog/digital-ocean-container-registry) -- hosts private container images used by `imageSource` deployments