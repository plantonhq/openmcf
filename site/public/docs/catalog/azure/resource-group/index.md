---
title: "Resource Group"
description: "Resource Group deployment documentation"
icon: "package"
order: 100
componentName: "azureresourcegroup"
---

# Azure Resource Group

Deploys an Azure Resource Group -- the foundational organizational container for all Azure resources. Every Azure resource must belong to a resource group, making this the typical first step in any Azure environment setup. The resource group provides lifecycle management, RBAC boundaries, and cost tracking for all contained resources.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Azure Resource Group** -- a named container in the specified Azure region that holds and organizes related Azure resources for unified lifecycle management, access control, and cost allocation. Planton-derived resource tags (organization, environment, resource kind, resource ID) are applied by the IaC module; the spec itself carries only name and region

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An active Azure subscription** where the resource group will be created. No other prerequisites are required -- resource groups are top-level containers.

## Deploy

### Console

Open the deployment store, find **Azure Resource Group**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureResourceGroup
metadata:
  name: platform-rg
  org: acme-corp
  env: prod
spec:
  name: "acme-prod-rg"
  region: "eastus"
```

```shell
planton apply -f resource-group.yaml
```

This creates a resource group named `acme-prod-rg` in the `eastus` region. Resources within the group can be deployed to different regions -- the resource group region only determines where its metadata is stored.

## Key Configuration

These are the most important decisions when configuring a resource group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Name** -- Must be unique within the Azure subscription. Use a descriptive naming convention like `{project}-{env}-rg` to make resource groups identifiable in the portal and CLI. Maximum 90 characters, supports alphanumeric, underscores, hyphens, periods, and parentheses.

**Region** -- Determines where the resource group metadata is stored. Resources inside the group can be in different regions. Choose a region close to your primary deployment target for consistency, though this has no performance impact on contained resources.

**Lifecycle implications** -- Deleting a resource group deletes all resources within it. Use separate resource groups for resources with different lifecycles (e.g., shared networking in one group, application workloads in another).

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `resource_group_id` | Azure Resource Manager ID of the resource group | Azure Policy assignments, diagnostic settings |
| `resource_group_name` | Name of the resource group | Nearly all Azure Cloud Resources reference this via `resourceGroup` field |
| `region` | Azure region where the resource group was created | Informational -- downstream resources specify their own region |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard resource group** -- Creates a named resource group in a specified region as the foundation for all subsequent Azure deployments. Virtually every Azure InfraChart begins with one or more resource groups. Start from the **Standard** preset.

## Works With

This component operates independently and does not reference other components.