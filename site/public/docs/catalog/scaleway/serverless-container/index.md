---
title: "Serverless Container"
description: "Serverless Container deployment documentation"
icon: "package"
order: 100
componentName: "scalewayserverlesscontainer"
---

# Scaleway Serverless Container

Deploys a serverless container on Scaleway as a composite resource that bundles a container namespace, the container deployment, and optional cron triggers into a single declarative unit. Runs pre-built Docker images from any OCI registry with configurable autoscaling (including scale-to-zero), health checks, Private Network connectivity, environment variables, and scheduled invocations. Supports ValueFromRef for Container Registry and Private Network dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Container Namespace** -- a Scaleway namespace that groups the container for lifecycle management and isolation
- **Serverless Container** -- the deployed container instance with the configured image, port, privacy, memory/CPU limits, scaling bounds, health check, and environment variables
- **Cron Triggers** -- created only when `cronTriggers` is populated; each trigger invokes the container on a CRON schedule with JSON arguments
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A container image** in a Scaleway Container Registry or any OCI-compliant registry (Docker Hub, GHCR, ECR). When using Scaleway Container Registry, reference the registry endpoint via ValueFromRef to create an InfraChart dependency edge.
- **A Scaleway Private Network** (optional) in the target region when the container needs to access databases, caches, or other Private Network resources without traversing the public internet.

## Deploy

### Console

Open the deployment store, find **Scaleway Serverless Container**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Public Web Service** preset in the [Presets](#presets) tab for a publicly accessible container that scales from zero.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayServerlessContainer
metadata:
  name: api-service
  org: acme-corp
  env: prod
spec:
  region: fr-par
  image:
    registryEndpoint:
      value: rg.fr-par.scw.cloud/my-namespace
    name: my-api
    tag: v1.0.0
  port: 8080
  privacy: privacy_public
```

```shell
planton apply -f scaleway-serverless-container.yaml
```

This creates a publicly accessible serverless container with default 256 MB memory, scale-to-zero behavior, and a maximum of 20 instances. No Private Network, health check, or cron triggers are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the container to a Container Registry and Private Network deployed in the same InfraPipeline:

```yaml
spec:
  image:
    registryEndpoint:
      valueFrom:
        kind: ScalewayContainerRegistry
        name: app-registry
        fieldPath: status.outputs.endpoint
  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id
```

The InfraPipeline resolves the dependency graph, deploys the Container Registry and Private Network first, then provisions the container with the resolved values.

## Key Configuration

These are the most important decisions when configuring a serverless container. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Privacy** -- The `privacy` field controls endpoint authentication. Use `privacy_public` for internet-accessible APIs and `privacy_private` for internal services that require Scaleway IAM token authentication.

**Scaling** -- Set `minScale` to 0 for scale-to-zero (no cost when idle) or 1+ for always-warm instances (eliminates cold starts). The `maxScale` field caps concurrent instances. Use `scalingOption` to configure autoscaling thresholds based on concurrent requests, CPU usage, or memory usage.

**Memory and CPU** -- The `memoryLimitMb` field sets the memory ceiling per instance (128-4096 MB). The `cpuLimit` field is optional -- when omitted, Scaleway allocates CPU proportional to memory.

**Health check** -- Configure `healthCheck` with a path, failure threshold, and interval for HTTP-based liveness detection. Unhealthy instances are automatically replaced. Without a health check, Scaleway uses TCP port probing.

**Cron triggers** -- Add `cronTriggers` for scheduled container invocations. Each trigger specifies a CRON expression and JSON arguments passed to the container's event object.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayContainerRegistry** | `image.registryEndpoint` | `status.outputs.endpoint` |
| **ScalewayPrivateNetwork** (optional) | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `container_id` | Unique identifier of the deployed container | Scaleway CLI operations, API management, Terraform import |
| `namespace_id` | Unique identifier of the container namespace | External cron triggers, tokens, additional containers |
| `domain_name` | Native Scaleway HTTPS endpoint for the container | ScalewayDnsRecord CNAME targets for custom domain routing |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public web service** -- A publicly accessible container that scales from zero to 20 instances with 256 MB memory and a 5-minute request timeout. The standard configuration for APIs and web applications. Start from the **Public Web Service** preset.

**VPC-connected private service** -- A private container attached to a Private Network with 512 MB memory, a minimum of 1 always-warm instance, HTTP health checking, and a maximum of 10 instances. The standard configuration for backend services accessing databases and caches. Start from the **VPC-Connected** preset.

## Works With

- [**Scaleway Container Registry**](/cloud-catalog/scaleway-container-registry) -- provides the container image registry endpoint referenced by the container image configuration
- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides VPC-internal connectivity for accessing databases, caches, and other private resources