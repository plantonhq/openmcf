---
title: "Container App Environment"
description: "Container App Environment deployment documentation"
icon: "package"
order: 100
componentName: "azurecontainerappenvironment"
---

# Azure Container App Environment

Deploy the hosting platform for Azure Container Apps: the secure boundary providing networking, logging, Dapr infrastructure, and compute capacity that every Container App and Job inside it shares.

## What Gets Created

- An Azure Container Apps Managed Environment (`Microsoft.App/managedEnvironments`)
- Optionally, the environment's custom DNS suffix configuration (when `customDomain` is set)

## Prerequisites

- An Azure Resource Group (`AzureResourceGroup`)
- For VNet injection: an `AzureSubnet` with a `/21` or larger address space
- For persistent logs: an `AzureLogAnalyticsWorkspace`

## Quick Start

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironment
metadata:
  name: my-env
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  environmentName: my-env
  logsDestination: LOG_ANALYTICS
  logAnalyticsWorkspaceId:
    value: /subscriptions/.../workspaces/my-law
```

## Configuration Reference

### Required Fields

| Field | Description |
| --- | --- |
| `region` | Azure region (ForceNew) |
| `resourceGroup` | Resource group name or `AzureResourceGroup` reference (ForceNew) |
| `environmentName` | 2-60 alphanumerics/hyphens; no leading/trailing hyphen (ForceNew) |

### Optional Fields

| Field | Description |
| --- | --- |
| `logsDestination` | `LOG_ANALYTICS` (requires the workspace) or `AZURE_MONITOR` (forbids it); unset with a workspace deploys log-analytics; unset without one is streaming-only |
| `logAnalyticsWorkspaceId` | The workspace persisting application logs |
| `daprApplicationInsightsConnectionString` | Dapr service-to-service telemetry export (secret; ForceNew) |
| `infrastructureSubnetId` | VNet injection; `/21`+ subnet (ForceNew) |
| `infrastructureResourceGroupName` | Name for the platform-managed infra RG; requires workload profiles (ForceNew) |
| `internalLoadBalancerEnabled` | VNet-only access; requires the subnet (ForceNew) |
| `zoneRedundancyEnabled` | Spread infrastructure across zones; requires the subnet (ForceNew) |
| `publicNetworkAccess` | `ENABLED` / `DISABLED`; unset lets Azure derive it |
| `mutualTlsEnabled` | Encrypt + authenticate all app-to-app traffic |
| `workloadProfiles` | Dedicated (D4-D32, E4-E32), GPU (NC24/48/96-A100), serverless GPU (CONSUMPTION_GPU_*) pools |
| `customDomain` | Environment-wide DNS suffix + wildcard PFX certificate |
| `identity` | `SYSTEM_ASSIGNED`, `USER_ASSIGNED` (+ identity references), or both |
| `tags` | Free-form Azure tags merged over platform tags |

## Examples

### VNet-Injected Internal Environment with Dedicated Compute

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironment
metadata:
  name: prod-env
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  environmentName: prod-env
  infrastructureSubnetId:
    valueFrom:
      kind: AzureSubnet
      name: apps-subnet
      fieldPath: status.outputs.subnet_id
  internalLoadBalancerEnabled: true
  zoneRedundancyEnabled: true
  mutualTlsEnabled: true
  logsDestination: LOG_ANALYTICS
  logAnalyticsWorkspaceId:
    valueFrom:
      kind: AzureLogAnalyticsWorkspace
      name: prod-law
      fieldPath: status.outputs.workspace_id
  workloadProfiles:
    - name: dedicated
      workloadProfileType: D4
      minimumCount: 1
      maximumCount: 5
```

### Serverless GPU Environment

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironment
metadata:
  name: inference-env
spec:
  region: westus3
  resourceGroup:
    value: ml-rg
  environmentName: inference-env
  workloadProfiles:
    - name: gpu-serverless
      workloadProfileType: CONSUMPTION_GPU_NC8AS_T4
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `environment_id` | ARM ID referenced by every family kind |
| `environment_name` | The environment's name |
| `default_domain` | Apps live at `{app}.{default_domain}` |
| `static_ip_address` | Public or private entry IP |
| `platform_reserved_cidr`, `platform_reserved_dns_ip_address`, `docker_bridge_cidr` | VNet planning values |
| `custom_domain_verification_id` | TXT-record value for domain ownership |
| `identity_principal_id` | System-assigned identity's RBAC principal |

## Related Components

- [Azure Container App](/docs/catalog/azure/container-app) -- continuously running services in this environment
- [Azure Container App Job](/docs/catalog/azure/container-app-job) -- run-to-completion workloads
- [Azure Container App Environment Storage](/docs/catalog/azure/container-app-environment-storage) -- mountable Azure Files shares
- [Azure Container App Environment Dapr Component](/docs/catalog/azure/container-app-environment-dapr-component) -- Dapr building-block backends
- [Azure Subnet](/docs/catalog/azure/subnet) -- the VNet-injection seam
- [Azure Log Analytics Workspace](/docs/catalog/azure/log-analytics-workspace) -- the logging seam
