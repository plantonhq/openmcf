# AzureServicePlan: Research & Design Documentation

## Executive Summary

Azure App Service Plan (`Microsoft.Web/serverfarms`) is the compute abstraction that underpins Azure Web Apps, Function Apps, and Logic Apps Standard. It defines the region, OS type, VM SKU, instance count, and pricing tier for the hosting environment. One or more apps share the same plan's compute resources.

This document captures the research and design rationale behind the `AzureServicePlan` Planton component (enum 442, id_prefix `azsp`). The component models the complete `azurerm_service_plan` surface: the full SKU vocabulary, all three OS types, App Service Environment placement, premium auto-scale, zone balancing, per-site scaling, elastic worker limits, and user tags.

## Azure Deployment Landscape

### Service Plan Architecture

Azure App Service runs on a multi-tenant platform. A Service Plan maps to a **server farm** -- a set of VMs in a specific region, at a specific pricing tier. The key architectural points:

1. **Region-locked**: Plans are created in a specific Azure region. All apps in the plan run in that region.
2. **OS-locked**: Plans are Linux (`reserved = true`), Windows, or Windows Container (`hyperV = true`). This is immutable after creation.
3. **Shared compute**: Multiple apps can run on the same plan, sharing CPU, memory, and instances.
4. **Scale unit**: Scaling the plan (changing `worker_count`) scales all apps in the plan simultaneously, unless per-site scaling is enabled.
5. **Tenancy**: Plans run on multi-tenant App Service by default; Isolated SKUs place the plan inside a single-tenant App Service Environment v3 (referenced by ARM ID).

### SKU Tier Comparison

| Tier | SKUs | Instances | Auto-Scale | Zones | SLA | Use Case |
|------|------|-----------|------------|-------|-----|----------|
| Free | F1 | 1 | No | No | None | Exploration |
| Shared | D1, SHARED | 1 | No | No | None | Low-traffic sites (Windows only) |
| Basic | B1-B3 | 1-3 | No | No | 99.95% | Dev/test |
| Standard | S1-S3 | 1-10 | Rule-based | No | 99.95% | Production entry |
| Premium v2 | P1v2-P3v2 | 1-30 | Yes | Yes | 99.95% | Production |
| Premium v3 | P0v3-P5mv3 | 1-30 | Yes (+ premium auto-scale) | Yes | 99.95% | Production (recommended) |
| Premium v4 | P0v4-P5mv4 | 1-30 | Yes (+ premium auto-scale) | Yes | 99.95% | Newest premium generation |
| Consumption | Y1 | 0-200 | Automatic | Yes | 99.95% | Serverless Functions |
| Elastic Premium | EP1-EP3 | 1-100 | Elastic | Yes | 99.95% | Functions (pre-warmed) |
| Flex Consumption | FC1 | automatic | Automatic | Yes | 99.95% | Newest serverless Functions tier |
| Isolated | I1-I3 (ASEv2), I1v2-I6v2 / I1mv2-I5mv2 (ASEv3) | 1-100 | Yes | Yes | 99.95% | Enterprise (single-tenant) |
| Workflow | WS1-WS3 | elastic | Elastic | Yes | 99.95% | Logic Apps Standard |

### Pricing Model

- **Free/Shared**: CPU minutes per day (shared VMs)
- **Basic-Premium**: Per-instance per-hour (dedicated VMs)
- **Consumption (Y1)**: Per-execution + per-GB-second (no idle cost)
- **Elastic Premium (EP*)**: Per-instance per-hour for pre-warmed + per-execution overage
- **Flex Consumption (FC1)**: Per-execution with per-instance memory selection and always-ready instances

### Premium v3 vs Premium v2

Premium v3 is the recommended tier for production workloads:
- **2x memory**: Compared to equivalent v2 SKUs
- **Better CPU**: Dv3 series VMs
- **Memory-optimized variants**: P1mv3-P5mv3 with doubled RAM per core
- **Same price**: P1v3 costs the same as P1v2 in most regions
- **P0v3**: The smallest Premium (1 vCPU, 4 GB RAM) -- cost-effective for light production

## Design Decisions

### 1. The SKU vocabulary is a closed proto enum

The provider validates SKU names against a closed catalog of known values. The spec models that catalog as the `AzureServicePlanSku` enum, family-prefixed for discoverability (`PREMIUM_P1V3`, `ELASTIC_PREMIUM_EP1`), with both IaC modules mapping enum values row-by-row to Azure's exact wire spellings (`P1v3`, `EP1`). A wrong SKU therefore fails at manifest-validation time instead of apply time, and the enum's family-contiguous numbering lets the SKU-gated rules (below) be expressed as simple range checks. When Azure ships a new SKU generation, the enum grows additively.

### 2. Provider validators are front-loaded as spec rules

The provider enforces four SKU pairings at plan/apply time; the spec enforces the same rules at manifest-validation time:

- `premium_plan_auto_scale_enabled` is Premium-only.
- `maximum_elastic_worker_count` above 1 on a Premium SKU requires premium auto-scale (Elastic Premium and Workflow support it natively).
- `zone_balancing_enabled` is rejected on Free/Shared/Basic/Standard.
- `app_service_environment_id` requires an Isolated SKU.

One provider behavior cannot be front-loaded and is documented on the fields instead: flipping `zone_balancing_enabled` from false to true with `worker_count < 2` forces the plan to be recreated (a state transition, not a static property of the manifest).

### 3. OS type is a three-value enum defaulting to Linux

Azure's API accepts Linux, Windows, and WindowsContainer; the spec models all three. An unset value deploys LINUX -- the catalog's app kinds (AzureLinuxWebApp, AzureFunctionApp) are Linux-based, and Linux is the right default for containers and every modern runtime.

### 4. App Service Environment placement is a raw ARM ID

`app_service_environment_id` accepts the ARM ID of an App Service Environment v3. The environment itself (`azurerm_app_service_environment_v3`) is not currently a catalog kind -- it is a heavyweight, subnet-anchored, hours-to-provision enterprise resource; the field composes with it by ID today and becomes a foreign-key reference if/when an ASE kind is added.

### 5. `maximum_elastic_worker_count` is the serverless cost lever

Elastic Premium plans auto-scale to 20 workers by default and can spike costs without a cap. The field applies to Elastic Premium and Workflow SKUs natively, and to Premium SKUs when premium auto-scale is enabled.

### 6. SKU changes never recreate the plan

`sku_name` is deliberately not ForceNew: plans scale up, down, and across tiers in place, and apps keep running through most SKU changes. Azure rejects the few impossible moves (e.g. Consumption <-> dedicated) at apply time.

## Terraform Provider Analysis

### Source Files

- `internal/services/appservice/service_plan_resource.go` -- Resource implementation
- `internal/services/appservice/helpers/service_plan.go` -- SKU catalog + tier-classification helpers
- `internal/services/appservice/validate/service_plan_name.go` -- Name validation

### Key Behaviors

1. **Name validation**: `^[0-9a-zA-Z-_]{1,60}$` (alphanumeric + hyphens + underscores, max 60)
2. **SKU case-insensitivity**: `DiffSuppressFunc` handles case-insensitive SKU comparison (the modules always send Azure's canonical spelling)
3. **worker_count computed**: If not specified, Azure sets it to the SKU's default capacity
4. **ForceNew fields**: `name`, `resource_group_name`, `location`, `os_type`
5. **Zone balancing ForceNew**: Enabling zone balancing with `worker_count < 2` forces recreation
6. **CustomizeDiff**: Validates the SKU pairings listed under Design Decisions at plan time; the spec front-loads all of them
7. **ASE create check**: Plans with a hosting environment require an `I`-prefixed (Isolated) SKU

### API Version

- Azure API: `Microsoft.Web` version `2023-12-01`
- Resource ID: `/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Web/serverFarms/{name}`

## Pulumi Provider Analysis

### Package

- `github.com/pulumi/pulumi-azure/sdk/v6/go/azure/appservice`
- Resource: `appservice.NewServicePlan`
- The provider is built through the shared `pulumiazureprovider` builder (static client secret, keyless web identity, or ambient chain)

### Field Mapping

| Spec Field | Pulumi Property |
|------------|----------------|
| `service_plan_name` | `Name` |
| `region` | `Location` |
| `resource_group` | `ResourceGroupName` |
| `os_type` (enum -> wire string) | `OsType` |
| `sku_name` (enum -> wire string) | `SkuName` |
| `app_service_environment_id` | `AppServiceEnvironmentId` |
| `worker_count` | `WorkerCount` |
| `premium_plan_auto_scale_enabled` | `PremiumPlanAutoScaleEnabled` |
| `maximum_elastic_worker_count` | `MaximumElasticWorkerCount` |
| `zone_balancing_enabled` | `ZoneBalancingEnabled` |
| `per_site_scaling_enabled` | `PerSiteScalingEnabled` |
| `tags` | `Tags` |

## Downstream Dependencies

### Resources that reference AzureServicePlan

| Resource | Field | Reference Path |
|----------|-------|---------------|
| AzureFunctionApp | `service_plan_id` | `status.outputs.service_plan_id` |
| AzureLinuxWebApp | `service_plan_id` | `status.outputs.service_plan_id` |

### Infra Charts

| Chart | Role |
|-------|------|
| `function-app-environment` | Foundation resource (Service Plan -> Function App) |
| `web-app-environment` | Foundation resource (Service Plan -> Linux Web App) |

## References

- [Azure App Service Plan documentation](https://learn.microsoft.com/en-us/azure/app-service/overview-hosting-plans)
- [Terraform azurerm_service_plan](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/service_plan)
- [Azure App Service pricing](https://azure.microsoft.com/en-us/pricing/details/app-service/linux/)
- [App Service Plan SKU comparison](https://learn.microsoft.com/en-us/azure/app-service/overview-compare)
