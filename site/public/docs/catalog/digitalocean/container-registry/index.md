---
title: "Container Registry"
description: "Container Registry deployment documentation"
icon: "package"
order: 100
componentName: "digitaloceancontainerregistry"
---

# Container Registry on DigitalOcean

Deploys a private, OCI-compliant container registry on DigitalOcean for storing Docker images and Helm charts. Configures the registry name, subscription tier, and region, then exposes the server URL as a stack output for downstream workloads to pull images. Integrates with Planton's Provider Connections for DigitalOcean credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Container Registry** -- a `digitalocean_container_registry` resource with the specified name, subscription tier, and region
- **Docker Credentials** (Terraform only) -- a `digitalocean_container_registry_docker_credentials` resource that generates write-enabled credentials for pushing and pulling images

DigitalOcean restricts each account to a single container registry. Deploying a second DigitalOceanContainerRegistry resource on the same account will fail.

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **No existing container registry** on the target DigitalOcean account. DigitalOcean allows only one registry per account.
- **A supported region** for container registry storage (e.g., `nyc3`, `sfo3`, `fra1`, `sgp1`). Choose the region nearest to your Kubernetes clusters or CI/CD pipelines for lowest pull latency.

## Deploy

### Console

Open the deployment store, find **Container Registry on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Professional Container Registry** preset in the [Presets](#presets) tab for a production-ready configuration with garbage collection.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1
kind: DigitalOceanContainerRegistry
metadata:
  name: prod-registry
  org: acme-corp
  env: prod
spec:
  name: prod-registry
  subscriptionTier: professional
  region: nyc3
  garbageCollectionEnabled: true
```

```shell
planton apply -f registry.yaml
```

This creates a professional-tier container registry in NYC3 with garbage collection enabled. Images are accessible at `registry.digitalocean.com/prod-registry/<repository>:<tag>`.

## Key Configuration

These are the most important decisions when configuring a container registry. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Subscription tier** -- The `subscriptionTier` field controls storage limits and pricing. `starter` is free with limited storage for development. `basic` provides moderate storage for small teams. `professional` offers the highest storage limits and is recommended for production teams pushing many images.

**Region** -- The `region` field determines where registry data is stored. Choose the region closest to your DigitalOcean Kubernetes clusters or build pipelines to minimize image pull latency.

**Garbage collection** -- Set `garbageCollectionEnabled` to `true` to automatically remove untagged images and reduce storage costs. Note that the Pulumi provisioner logs a warning for this field because the upstream DigitalOcean provider does not yet support it natively.

**Registry naming** -- The `name` field must be unique within your DigitalOcean account and is used in image paths (`registry.digitalocean.com/<name>/<image>:<tag>`). Choose a stable name since changing it requires re-tagging and re-pushing all images.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `registry_name` | Name of the created container registry | Image path construction, Kubernetes `imagePullSecret` configuration |
| `server_url` | Full registry URL (e.g., `registry.digitalocean.com/prod-registry`) | Docker login, CI/CD pipeline push targets, App Platform image source |
| `region` | Region slug where the registry is hosted | Verifying data locality alignment with clusters |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Professional registry for production** -- Professional tier with garbage collection enabled, suitable for teams pushing images frequently from CI/CD pipelines. Provides the highest storage limits and automatic cleanup of untagged images. Start from the **Professional Container Registry** preset.

## Works With

This component operates independently and does not reference other components.