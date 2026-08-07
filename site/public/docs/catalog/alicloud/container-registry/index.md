---
title: "Container Registry"
description: "Container Registry deployment documentation"
icon: "package"
order: 100
componentName: "alicloudcontainerregistry"
---

# AliCloud Container Registry

Deploys an Alibaba Cloud Container Registry (ACR) Enterprise Edition instance with bundled namespaces for organizing container images. ACR Enterprise provides three tiers (Basic, Standard, Advanced) with enterprise-grade security, optional multi-region replication, and image scanning. The component integrates with Planton's Provider Connections for AliCloud credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ACR Enterprise Instance** -- an `alicloud_cr_ee_instance` with the selected tier, payment model, and login password
- **Namespaces** -- one `alicloud_cr_ee_namespace` per entry in the `namespaces` list, each with configurable auto-create and default visibility settings

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **An instance name** -- the human-readable identifier for the ACR Enterprise Edition instance within your account.
- **A login password** (optional) -- used for authenticating `docker login`, image push, and image pull operations.

## Deploy

### Console

Open the deployment store, find **AliCloud Container Registry**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic Dev** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudContainerRegistry
metadata:
  name: platform-registry
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  instanceName: acme-platform
  instanceType: Standard
  namespaces:
    - name: backend
      autoCreate: true
      defaultVisibility: PRIVATE
    - name: frontend
      autoCreate: true
      defaultVisibility: PRIVATE
```

```shell
planton apply -f container-registry.yaml
```

This creates a Standard-tier ACR Enterprise instance with Subscription billing and two private namespaces with auto-create enabled. Images are addressed as `{public_endpoint}/{namespace}/{repo}:{tag}`.

## Key Configuration

These are the most important decisions when configuring a container registry. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Instance tier** -- Set `instanceType` to `Basic` (individual developers or small teams), `Standard` (small and medium enterprises), or `Advanced` (large enterprises with geo-replication needs). The tier is immutable after creation and determines image storage limits, concurrent pull capacity, and available features.

**Payment model** -- Set `paymentType` to `Subscription` (default, pre-paid monthly/yearly billing) or `PayAsYouGo` (post-paid usage-based billing). When using Subscription, set `period` to the billing cycle in months (e.g., 1, 3, 6, 12). The payment type is immutable after creation.

**Namespace organization** -- Namespaces are the top-level organizational units for container images. A typical pattern is one namespace per team or application (e.g., `platform`, `frontend`, `backend`). Set `autoCreate: true` to automatically create image repositories when a CI/CD pipeline pushes to a repository name that does not yet exist. Set `defaultVisibility` to `PRIVATE` (recommended) or `PUBLIC` for each namespace.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | ACR Enterprise instance ID | Granting ACK worker nodes pull access via RAM role |
| `instance_name` | Registry instance name | Audit references, organizational lookups |
| `public_endpoint` | Internet-facing registry endpoint domain | CI/CD pipeline `docker login` and image push |
| `vpc_endpoint` | VPC-internal registry endpoint domain | ACK node image pulls (faster, no internet egress cost) |
| `namespace_ids` | Map of namespace names to their IDs | Resource tracking, access control configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Basic for development** -- A Basic-tier instance with PayAsYouGo billing and a single namespace for development and testing. Start from the **Basic Dev** preset.

**Standard for production** -- A Standard-tier instance with Subscription billing and multiple namespaces organized by team or application. Suitable for most production workloads. Start from the **Standard Production** preset.

**Advanced for enterprise** -- An Advanced-tier instance with Subscription billing, geo-replication capabilities, and namespaces for multi-team organizations. Start from the **Advanced Enterprise** preset.

## Works With

This component operates independently and does not reference other deployment components.