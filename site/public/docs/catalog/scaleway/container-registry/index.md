---
title: "Container Registry"
description: "Container Registry deployment documentation"
icon: "package"
order: 100
componentName: "scalewaycontainerregistry"
---

# Scaleway Container Registry

Deploys a fully managed, OCI-compliant Container Registry namespace on Scaleway for storing and distributing container images and Helm charts. Each namespace provides a dedicated Docker endpoint for push and pull operations, with configurable visibility (private or public). The registry integrates with Scaleway Kapsule clusters, Serverless Containers, and CI/CD pipelines for image distribution.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Registry Namespace** -- a `scaleway_registry_namespace` in the specified region with a Docker endpoint at `rg.<region>.scw.cloud/<namespace-name>` and the configured visibility setting

Note: Scaleway Container Registry namespaces do not support tags. Standard metadata tags are not applied to this resource.

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair (Access Key + Secret Key). The IaC module authenticates through the Scaleway provider configuration.
- **Choose a region** -- registry namespaces are regional resources. Available regions: `fr-par` (Paris), `nl-ams` (Amsterdam), `pl-waw` (Warsaw). Cannot be changed after creation.
- **Namespace naming** -- the name is derived from `metadata.name` and becomes part of the Docker endpoint URL. It must be 4-63 characters, lowercase alphanumeric with hyphens, and unique within the Scaleway project.

## Deploy

### Console

Open the deployment store, find **Scaleway Container Registry**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Private Registry** preset in the [Presets](#presets) tab to create a private namespace that requires authentication for image pulls.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayContainerRegistry
metadata:
  name: app-images
  org: acme-corp
  env: prod
spec:
  region: fr-par
```

```shell
planton apply -f scaleway-container-registry.yaml
```

This creates a private registry namespace in the Paris region. After provisioning, authenticate with `docker login rg.fr-par.scw.cloud/app-images -u nologin -p <SCW_SECRET_KEY>` and push images with `docker push rg.fr-par.scw.cloud/app-images/myapp:latest`. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a container registry. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Visibility** -- The `isPublic` field controls whether images require authentication to pull. When false (default), only authenticated users and services can pull images -- the standard choice for proprietary application images. When true, anyone can pull without credentials, suitable for open-source projects and public base images. Pushing always requires authentication regardless of this setting. Can be toggled after creation.

**Region** -- The `region` field determines the Docker endpoint URL and where images are stored. Choose the region closest to your build infrastructure and deployment targets for lowest push/pull latency. Consider data residency requirements for container images.

**Description** -- The `description` field provides a human-readable label for the namespace (e.g., "Production microservices images"). Displayed in the Scaleway Console for identification. No effect on registry behavior.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace_id` | Regional ID of the registry namespace (`{region}/{uuid}`) | Terraform state references, Scaleway API operations |
| `endpoint` | Docker endpoint URL (`rg.<region>.scw.cloud/<name>`) | CI/CD `docker push` target, Kubernetes imagePullSecrets, serverless image source |
| `namespace_name` | Name of the registry namespace in Scaleway | CI/CD pipeline variables, dashboard labels |
| `region` | Region where the namespace is deployed | Co-location decisions for Kapsule clusters and serverless functions |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private registry** -- A private namespace requiring authentication for image pulls. The standard and most common configuration for proprietary application images, internal tools, and CI/CD artifacts. Start from the **Private Registry** preset.

## Works With

This component operates independently and does not reference other components.