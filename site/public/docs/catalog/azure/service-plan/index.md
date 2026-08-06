---
title: "Service Plan"
description: "Service Plan deployment documentation"
icon: "package"
order: 100
componentName: "azureserviceplan"
---

# Azure Service Plan

Deploys an Azure App Service Plan that defines the compute tier, VM size, instance count, and pricing for hosting Azure Web Apps, Function Apps, and Logic Apps Standard workflows. The plan supports Linux, Windows, and Windows Container operating systems, zone-redundant deployments, premium auto-scaling, per-site scaling, elastic worker limits for serverless workloads, and App Service Environment placement for single-tenant compute.

## What Gets Created

When you deploy an AzureServicePlan resource, Planton provisions:

- **App Service Plan** -- an `azurerm_service_plan` / `appservice.ServicePlan` resource in the specified region and resource group, configured with the chosen SKU tier, OS type, instance count, and optional zone balancing
- **Azure Tags** -- platform metadata tags merged with your own `tags` (your tags win on key collision) for tracking, cost allocation, and governance

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the plan will be created (can reference an AzureResourceGroup resource)
- **SKU selection** -- determine the appropriate tier before deployment: `CONSUMPTION_Y1` for pay-per-execution Function Apps, `ELASTIC_PREMIUM_EP1`-`EP3` for pre-warmed Functions, `BASIC_B1`-`B3` for dev/test web apps, `PREMIUM_P1V3`-`P3V3` for production workloads

## Quick Start

Create a file `serviceplan.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServicePlan
metadata:
  name: my-plan
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureServicePlan.my-plan
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  servicePlanName: my-plan
  skuName: BASIC_B1
```

Deploy:

```shell
planton apply -f serviceplan.yaml
```

This creates a Linux Basic B1 App Service Plan with a single worker instance in the `eastus` region.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the plan (e.g., `eastus`, `westeurope`). **ForceNew**: changing this destroys and recreates the plan. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. **ForceNew**: changing this destroys and recreates the plan. | Required |
| `servicePlanName` | `string` | Name of the Service Plan. Unique within the resource group. **ForceNew**: changing this destroys and recreates the plan. | Required, 1-60 characters, pattern `^[0-9a-zA-Z-_]{1,60}$` |
| `skuName` | `enum` | SKU determining pricing tier and compute capacity -- NOT ForceNew, plans re-tier in place. See SKU reference below. | Required, closed vocabulary |

**SKU reference** -- values by category (deployed in Azure's exact spelling):

- **Free/Shared**: `FREE_F1`, `SHARED_D1` (Windows only)
- **Basic**: `BASIC_B1`-`B3` (manual scale to 3 instances)
- **Standard**: `STANDARD_S1`-`S3` (autoscale to 10, staging slots)
- **Premium**: `PREMIUM_P1V2`-`P3V2`, `PREMIUM_P0V3`-`P3V3`, `PREMIUM_P0V4`-`P3V4`, and memory-optimized `PREMIUM_P1MV3`-`P5MV3` / `PREMIUM_P1MV4`-`P5MV4` (30 instances, zone redundancy, auto-scale)
- **Consumption**: `CONSUMPTION_Y1` (Function Apps pay-per-execution)
- **Elastic Premium**: `ELASTIC_PREMIUM_EP1`-`EP3` (Function Apps pre-warmed)
- **Flex Consumption**: `FLEX_CONSUMPTION_FC1` (the newest serverless Functions tier)
- **Isolated**: `ISOLATED_I1`-`I3` (ASEv2), `ISOLATED_I1V2`-`I6V2` and memory-optimized `ISOLATED_I1MV2`-`I5MV2` (ASEv3; require `appServiceEnvironmentId`)
- **Workflow**: `WORKFLOW_WS1`-`WS3` (Logic Apps Standard)

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `osType` | `enum` | `LINUX` | Operating system type: `LINUX`, `WINDOWS`, or `WINDOWS_CONTAINER`. All apps within a plan share the OS type. Shared-tier SKUs are Windows-only. **ForceNew**: changing this destroys and recreates the plan. |
| `appServiceEnvironmentId` | `string` | - | ARM ID of the App Service Environment v3 to place the plan in (single-tenant compute). Only Isolated SKUs may set this. |
| `workerCount` | `int` | _(SKU default, typically 1)_ | Number of VM instances allocated to the plan. Maximum varies by SKU: Basic=3, Standard=10, Premium=30, Isolated=100. Serverless tiers (Y1/FC1/EP*) manage this automatically. When `zoneBalancingEnabled` is `true`, use a multiple of the zone count (typically 3) and never fewer than 2. |
| `premiumPlanAutoScaleEnabled` | `bool` | `false` | Automatic HTTP-load scaling for Premium plans, up to `maximumElasticWorkerCount`. Premium SKUs only. |
| `maximumElasticWorkerCount` | `int` | _(platform default, typically 20)_ | Scale-out ceiling for Elastic Premium and Workflow plans (and Premium plans with auto-scale enabled). The primary cost-control lever for serverless workloads. Range: 0-100. |
| `zoneBalancingEnabled` | `bool` | `false` | Distribute instances across availability zones. Supported on Premium, Elastic Premium, Consumption, Flex Consumption, Isolated, and Workflow SKUs -- not Free/Shared/Basic/Standard. |
| `perSiteScalingEnabled` | `bool` | `false` | Allow individual apps within the plan to scale independently of the plan's instance count. |
| `tags` | `map` | - | Free-form Azure resource tags, merged over the platform's metadata-derived tags (your tags win on key collision). |

## Examples

### Development Function App Plan

A Consumption plan for pay-per-execution Azure Functions in development:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServicePlan
metadata:
  name: dev-functions
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureServicePlan.dev-functions
spec:
  region: eastus
  resourceGroup:
    value: dev-rg
  servicePlanName: dev-functions
  skuName: CONSUMPTION_Y1
```

### Production Web App Plan with Zone Redundancy

A Premium v3 plan with zone balancing and multiple workers for production web applications:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServicePlan
metadata:
  name: prod-web
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureServicePlan.prod-web
spec:
  region: westeurope
  resourceGroup:
    value: prod-rg
  servicePlanName: prod-web
  skuName: PREMIUM_P1V3
  workerCount: 3
  zoneBalancingEnabled: true
```

### Premium Plan with Automatic HTTP Scaling

A Premium v3 plan that scales itself on HTTP load, capped at 10 instances:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServicePlan
metadata:
  name: autoscale-web
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureServicePlan.autoscale-web
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  servicePlanName: autoscale-web
  skuName: PREMIUM_P1V3
  premiumPlanAutoScaleEnabled: true
  maximumElasticWorkerCount: 10
```

### Elastic Premium with Worker Limits

An Elastic Premium plan for serverless Function Apps with a capped elastic worker count to control scaling costs:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServicePlan
metadata:
  name: events-plan
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureServicePlan.events-plan
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  servicePlanName: events-plan
  skuName: ELASTIC_PREMIUM_EP1
  maximumElasticWorkerCount: 50
```

### Using Foreign Key References

Reference a Planton-managed resource group instead of hardcoding the name:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServicePlan
metadata:
  name: ref-plan
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureServicePlan.ref-plan
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      fieldPath: status.outputs.resource_group_name
  servicePlanName: ref-plan
  skuName: PREMIUM_P1V3
  workerCount: 3
  zoneBalancingEnabled: true
  perSiteScalingEnabled: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `service_plan_id` | `string` | Azure Resource Manager ID of the Service Plan. Referenced by AzureFunctionApp and AzureLinuxWebApp via `servicePlanId`. |
| `service_plan_name` | `string` | Name of the Service Plan |
| `os_type` | `string` | Configured operating system type in Azure's spelling (`Linux`, `Windows`, `WindowsContainer`) |
| `sku_name` | `string` | Configured SKU name in Azure's spelling (e.g., `P1v3`, `EP1`, `Y1`) |
| `kind` | `string` | Azure's computed plan classification (e.g., `linux`, `elastic`, `functionapp`) |
| `reserved` | `bool` | Whether the plan runs Linux workers |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/resource-group) -- provides the resource group for plan placement
- [AzureFunctionApp](/docs/catalog/azure/function-app) -- serverless Function Apps hosted on this plan
- [AzureLinuxWebApp](/docs/catalog/azure/linux-web-app) -- Linux web applications hosted on this plan
