# AzureServicePlan

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureServicePlanSpec** defines the configuration for creating an Azure
App Service Plan (Microsoft.Web/serverfarms), the compute tier that hosts
Azure Web Apps, Function Apps, and Logic Apps Standard workflows.

An App Service Plan defines a set of compute resources for applications to
run on. One or more apps can share the same plan (and therefore the same
VMs); the plan determines the region, VM size, instance count, pricing
tier, and which platform features (staging slots, zone redundancy,
per-site scaling, VNet integration) are available to the apps it hosts.

**OS type**: Each plan is either Linux, Windows, or Windows Container.
The OS type is immutable after creation, and all apps within a plan must
share it. Linux plans set `reserved = true` in the Azure API; Windows
Container plans set `hyperV = true`.

**Choosing a SKU** (the full vocabulary is the `AzureServicePlanSku`
enum; the tier determines both price and capability):
  - Free (F1) / Shared (D1, SHARED): apps share VMs with other customers;
    60/240 CPU-minutes per day, no scale-out, no SLA, no Always-On.
    Shared-tier SKUs are Windows-only.
  - Basic (B1-B3): dedicated VMs, manual scale to 3 instances. Dev/test.
  - Standard (S1-S3): auto-scale to 10 instances, staging slots, backups.
  - Premium (P*v2/P*v3/P*v4, memory-optimized P*m*): faster hardware,
    scale to 30 instances, zone redundancy, VNet features. P*v3/P*v4 offer
    the best price-performance and support the optional premium plan
    auto-scale (`premium_plan_auto_scale_enabled`).
  - Consumption (Y1): Function Apps pay-per-execution; scales to 200
    instances automatically; instances are recycled when idle.
  - Elastic Premium (EP1-EP3): Function Apps with pre-warmed instances
    (no cold start) and event-driven scale to 100 instances; the
    `maximum_elastic_worker_count` field is the cost-control lever.
  - Flex Consumption (FC1): the newest serverless Functions tier
    (per-instance memory selection, always-ready instances).
  - Isolated (I1-I3 for ASEv2, I*v2/I*mv2 for ASEv3): single-tenant
    compute inside an App Service Environment; requires
    `app_service_environment_id`.
  - Workflow (WS1-WS3): Logic Apps Standard hosting.

**Zone redundancy** (`zone_balancing_enabled`): supported on Premium,
Elastic Premium, Consumption, Flex Consumption, Isolated v2, and Workflow
SKUs -- not on Free, Shared, Basic, or Standard. Enabling it on an
existing plan whose `worker_count` is below 2 forces the plan to be
destroyed and recreated; keep the instance count a multiple of the
region's availability-zone count (typically 3) for even spread.

**Per-site scaling** (`per_site_scaling_enabled`): lets individual apps
within the plan scale to fewer instances than the plan itself runs,
instead of every app running on every instance.

**ForceNew fields** (changing these destroys and recreates the plan --
and every app on it goes down with it): `service_plan_name`, `os_type`,
`region`, `resource_group`. The SKU is NOT ForceNew: plans scale up,
down, and across tiers in place (moving between Consumption and
dedicated tiers is the exception Azure rejects at apply time).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServicePlan
metadata:
  name: test-plan
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  servicePlanName: planton-hack-plan
  # Explicit OS type exercises the enum mapping; unset deploys LINUX.
  osType: LINUX
  # Premium v3 exercises the premium range gates: auto-scale plus the
  # elastic scale-out ceiling that requires it.
  skuName: PREMIUM_P1V3
  workerCount: 3
  premiumPlanAutoScaleEnabled: true
  maximumElasticWorkerCount: 10
  # Zone balancing exercises the tier gate (premium and above) and the
  # worker-count >= 2 recreate nuance.
  zoneBalancingEnabled: true
  perSiteScalingEnabled: true
  tags:
    cost-center: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.servicePlanName` | `string` | yes |  |  |
| `spec.osType` | `enum` |  |  |  |
| `spec.skuName` | `enum` | yes |  |  |
| `spec.appServiceEnvironmentId` | `string` |  |  |  |
| `spec.workerCount` | `int32` |  |  |  |
| `spec.premiumPlanAutoScaleEnabled` | `bool` |  | `false` |  |
| `spec.maximumElasticWorkerCount` | `int32` |  |  |  |
| `spec.zoneBalancingEnabled` | `bool` |  | `false` |  |
| `spec.perSiteScalingEnabled` | `bool` |  | `false` |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the Service Plan will be created.
Examples: "eastus", "westus2", "westeurope", "southeastasia".
Every app hosted on the plan runs in this region.

**ForceNew**: Changing this destroys and recreates the plan.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure Resource Group where the Service Plan will be created.
Can be a literal string or a reference to an AzureResourceGroup output.

**ForceNew**: Changing this destroys and recreates the plan.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.servicePlanName

`string` · required

The name of the Service Plan.
Allowed characters: alphanumeric, hyphens, and underscores.
Length: 1 to 60 characters.

This name appears in the Azure portal and CLI. It does not need to be
globally unique (uniqueness is scoped to the resource group).

**ForceNew**: Changing this destroys and recreates the plan.

- rule: Service Plan name must contain only alphanumeric characters, hyphens, and underscores
- rule: {"required":true,"string":{"minLen":"1","maxLen":"60"}}

### spec.osType

`enum`

The operating system type for the plan. Unset deploys LINUX -- the
right choice for containers and every modern runtime; the catalog's
app kinds (AzureLinuxWebApp, AzureFunctionApp) are Linux-based.

All apps within a plan must share the plan's OS type. The Shared-tier
SKUs (D1, SHARED) are Windows-only.

**ForceNew**: Changing this destroys and recreates the plan.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_service_plan_os_type_unspecified` -- Not specified -- deploys LINUX.
- `LINUX` -- Linux VMs (`reserved = true` in the Azure API). The right choice for containers and every modern runtime.
- `WINDOWS` -- Windows VMs. Required for .NET Framework apps and the Shared-tier SKUs.
- `WINDOWS_CONTAINER` -- Windows Container VMs (`hyperV = true` in the Azure API) -- hosts Windows containers on Premium v3 and Isolated v2 SKUs.

### spec.skuName

`enum` · required

The SKU that determines the pricing tier and compute capacity -- one
value picks both the tier's capabilities (slots, zone redundancy,
scale-out ceiling) and the VM size. See the message-level comment for
the tier guide and the AzureServicePlanSku enum for the full
vocabulary.

NOT ForceNew: plans re-tier in place (apps keep running through most
SKU changes; Azure rejects the few impossible moves at apply time).

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_service_plan_sku_unspecified` -- Not specified -- invalid; pick an explicit tier and size.
- `FREE_F1` -- Free tier: shared VMs, 60 CPU-minutes/day, no scale-out, no SLA, no Always-On. Dev/experiments only.
- `SHARED_D1` -- Shared tier (Windows-only): shared VMs, 240 CPU-minutes/day, custom domains but no SSL bindings, no scale-out, no SLA.
- `SHARED` -- Shared tier's legacy ARM spelling ("SHARED"), functionally D1.
- `BASIC_B1` -- Basic (B-family): dedicated VMs, manual scale to 3 instances, custom domains + SSL. Dev/test and low-traffic production.
- `BASIC_B2`
- `BASIC_B3`
- `STANDARD_S1` -- Standard (S-family): auto-scale to 10 instances, 5 staging slots, daily backups, Traffic Manager integration.
- `STANDARD_S2`
- `STANDARD_S3`
- `PREMIUM_P1V2` -- Premium v2 (P*v2): Dv2-series hardware, scale to 30 instances, 20 slots, zone redundancy.
- `PREMIUM_P2V2`
- `PREMIUM_P3V2`
- `PREMIUM_P0V3` -- Premium v3 (P*v3): the current price-performance sweet spot -- faster hardware, more memory per core, Windows containers, zone redundancy, and premium plan auto-scale.
- `PREMIUM_P1V3`
- `PREMIUM_P2V3`
- `PREMIUM_P3V3`
- `PREMIUM_P1MV3` -- Premium v3 memory-optimized (P*mv3): double the RAM per core of the equivalent P*v3.
- `PREMIUM_P2MV3`
- `PREMIUM_P3MV3`
- `PREMIUM_P4MV3`
- `PREMIUM_P5MV3`
- `PREMIUM_P0V4` -- Premium v4 (P*v4): the newest premium generation.
- `PREMIUM_P1V4`
- `PREMIUM_P2V4`
- `PREMIUM_P3V4`
- `PREMIUM_P1MV4` -- Premium v4 memory-optimized (P*mv4).
- `PREMIUM_P2MV4`
- `PREMIUM_P3MV4`
- `PREMIUM_P4MV4`
- `PREMIUM_P5MV4`
- `CONSUMPTION_Y1` -- Consumption (Y1): Function Apps pay-per-execution -- scales to 200 instances automatically, bills per GB-second, cold starts apply.
- `ELASTIC_PREMIUM_EP1` -- Elastic Premium (EP-family): Function Apps with pre-warmed instances (no cold start), VNet integration, and event-driven scale to 100 instances. Pair with maximum_elastic_worker_count to cap cost.
- `ELASTIC_PREMIUM_EP2`
- `ELASTIC_PREMIUM_EP3`
- `FLEX_CONSUMPTION_FC1` -- Flex Consumption (FC1): the newest serverless Functions tier -- per-instance memory selection, always-ready instances, VNet support. Note: Azure models Flex Consumption function apps as their own resource type, distinct from AzureFunctionApp.
- `ISOLATED_I1` -- Isolated v1 (I-family, ASEv2): single-tenant compute in an App Service Environment v2 (legacy -- prefer Isolated v2 on ASEv3).
- `ISOLATED_I2`
- `ISOLATED_I3`
- `ISOLATED_I1V2` -- Isolated v2 (I*v2, ASEv3): single-tenant, network-isolated compute; scale to 100 instances. Requires app_service_environment_id.
- `ISOLATED_I2V2`
- `ISOLATED_I3V2`
- `ISOLATED_I4V2`
- `ISOLATED_I5V2`
- `ISOLATED_I6V2`
- `ISOLATED_I1MV2` -- Isolated v2 memory-optimized (I*mv2, ASEv3).
- `ISOLATED_I2MV2`
- `ISOLATED_I3MV2`
- `ISOLATED_I4MV2`
- `ISOLATED_I5MV2`
- `WORKFLOW_WS1` -- Workflow (WS-family): Logic Apps Standard hosting with elastic scale-out (maximum_elastic_worker_count applies).
- `WORKFLOW_WS2`
- `WORKFLOW_WS3`

### spec.appServiceEnvironmentId

`string`

The ARM ID of the App Service Environment v3 to place the plan in --
single-tenant, network-isolated App Service compute. Only Isolated
SKUs (I1-I3, I*v2, I*mv2) can be placed in an environment (enforced
here, exactly as Azure enforces it at creation).
Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/
        Microsoft.Web/hostingEnvironments/{name}
Leave empty for regular multi-tenant App Service (the norm).

- rule: app_service_environment_id must be a full App Service Environment ARM ID (/subscriptions/.../providers/Microsoft.Web/hostingEnvironments/{name})

### spec.workerCount

`int32` · optional (explicit presence)

Number of VM instances (workers) allocated to the plan.
Minimum: 1. If not specified, Azure defaults to the SKU's default
capacity (typically 1).

Maximum varies by tier: Basic=3, Standard=10, Premium=30,
Isolated=100. Consumption (Y1), Flex Consumption (FC1), and Elastic
Premium (EP*) manage instance count automatically -- leave this unset
for those tiers.

When `zone_balancing_enabled` is true, set this to a multiple of the
number of availability zones in the region (typically 3) for even
distribution -- and never below 2, or enabling zone balancing later
will force the plan to be recreated.

- rule: {"int32":{"gte":1}}

### spec.premiumPlanAutoScaleEnabled

`bool` · optional (explicit presence)

Enable automatic scaling for Premium (P*v2/P*v3/P*v4) plans -- Azure
adds and removes instances based on HTTP load without an autoscale
rule resource, up to `maximum_elastic_worker_count`. Only valid on
Premium SKUs (enforced here, exactly as the provider enforces it).

Default: false

- default: `false`

### spec.maximumElasticWorkerCount

`int32` · optional (explicit presence)

The upper bound on how many workers the plan can scale to when
handling events -- the primary cost-control lever for serverless
workloads.

Applies to:
- Elastic Premium SKUs (EP1-EP3): default 20, maximum 100.
- Workflow SKUs (WS1-WS3): Logic Apps Standard scale-out.
- Premium SKUs, only when `premium_plan_auto_scale_enabled` is true.

Setting a value above 1 on a Premium plan without premium plan
auto-scale is rejected (enforced here, exactly as the provider
enforces it).

- rule: {"int32":{"gte":0}}

### spec.zoneBalancingEnabled

`bool` · optional (explicit presence)

Distribute the plan's instances across availability zones for higher
resilience. Supported on Premium, Elastic Premium, Consumption, Flex
Consumption, Isolated, and Workflow SKUs -- NOT on Free, Shared,
Basic, or Standard (enforced here, exactly as the provider enforces
it).

Flipping this from false to true on an existing plan with
`worker_count` below 2 destroys and recreates the plan; with 2 or
more workers it applies in place.

Default: false

- default: `false`

### spec.perSiteScalingEnabled

`bool` · optional (explicit presence)

Enable per-site scaling, allowing individual apps within the plan to
scale to fewer instances than the plan itself runs (by default every
app runs on every instance of the plan).

Default: false

- default: `false`

### spec.tags

`map<string, string>`

Free-form Azure resource tags applied to the plan, merged over the
platform's metadata-derived tags (user tags win on key collision) --
the hooks for cost allocation, chargeback reports, and Azure Policy
governance rules that filter or group by them. Updatable in place.

## Validation Rules

- `service_plan_auto_scale_premium_only`: premium_plan_auto_scale_enabled can only be set on Premium SKUs (P*v2, P*v3, P*v4, and their memory-optimized variants)
- `service_plan_elastic_worker_count_sku_gate`: maximum_elastic_worker_count above 1 on a Premium SKU requires premium_plan_auto_scale_enabled: true (Elastic Premium and Workflow SKUs support it natively)
- `service_plan_zone_balancing_sku_gate`: zone_balancing_enabled is not supported on Free, Shared, Basic, or Standard SKUs -- use a Premium, Elastic Premium, Isolated, or Workflow SKU
- `service_plan_ase_requires_isolated_sku`: plans inside an App Service Environment require an Isolated SKU (I1-I3, I1v2-I6v2, or I1mv2-I5mv2)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureServicePlan, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.service_plan_id` | `string` | The Azure Resource Manager ID of the Service Plan. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Web/serverFarms/{name} This is the primary reference output. Referenced by: - AzureFunctionApp (service_plan_id) - AzureLinuxWebApp (service_plan_id) |
| `status.outputs.service_plan_name` | `string` | The name of the Service Plan. Echo of the input name, useful for debugging and audit trails. |
| `status.outputs.os_type` | `string` | The configured operating system type ("Linux", "Windows", or "WindowsContainer"). Informational output for downstream visibility. |
| `status.outputs.sku_name` | `string` | The configured SKU name in Azure's spelling (e.g. "P1v3", "EP1", "Y1"). Informational output for cost tracking and capacity planning. |
| `status.outputs.kind` | `string` | Azure's computed plan kind (e.g. "linux", "elastic", "functionapp") -- the API's own classification of the plan, read back after creation. |
| `status.outputs.reserved` | `bool` | Whether the plan runs Linux workers (`reserved = true` in the Azure API). Read back after creation. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFunctionApp | `spec.servicePlanId` | `status.outputs.service_plan_id` |
| AzureLinuxWebApp | `spec.servicePlanId` | `status.outputs.service_plan_id` |

## See Also

- [Overview](../README.md)
