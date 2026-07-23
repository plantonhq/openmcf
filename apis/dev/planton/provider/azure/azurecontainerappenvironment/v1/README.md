# AzureContainerAppEnvironment

Deploy an Azure Container Apps Managed Environment -- the hosting platform and secure boundary for Container Apps and Jobs.

## Overview

A Container App Environment is the execution boundary containerized workloads share: one environment provides the networking (external or VNet-injected), the logging pipeline, the Dapr infrastructure, and the compute capacity (serverless Consumption plus optional dedicated/GPU workload profiles) for every app and job inside it.

The environment is the container kind of the Container Apps family. The workloads and platform pieces are their own composable components referencing its `environment_id` output: `AzureContainerApp` (services), `AzureContainerAppJob` (run-to-completion work), `AzureContainerAppEnvironmentStorage` (mountable Azure Files shares), and `AzureContainerAppEnvironmentDaprComponent` (Dapr backends).

## Key Features

- **Three networking modes**: external (Azure-managed, public endpoints), VNet-injected (private connectivity through your `/21`+ subnet), and internal (VNet-injected behind an internal load balancer -- no public exposure)
- **Workload profiles**: serverless Consumption by default; dedicated D/E-family compute and NVIDIA A100/T4 GPU profiles (including serverless GPU) declared per environment and selected by name from apps and jobs
- **Logging destinations**: Log Analytics (KQL, alerting), Azure Monitor diagnostic routing, or streaming-only
- **Zone redundancy** for VNet-injected environments, spreading infrastructure across availability zones
- **Mutual TLS** between apps for encrypted, authenticated east-west traffic
- **Custom DNS suffix**: replace the generated `*.{region}.azurecontainerapps.io` domain with your own wildcard-certificate-backed suffix
- **Managed identity** (system- and/or user-assigned) for the environment's platform operations
- **User tags** merged over platform-derived tags for cost allocation and governance

## When to Use

- Any Azure Container Apps adoption -- every app and job needs an environment first
- Serverless microservices that scale to zero (Consumption)
- Private backend platforms (VNet-injected + internal load balancer)
- GPU inference workloads (Consumption-GPU or dedicated NC profiles)

## Spec Highlights

| Field | Notes |
| --- | --- |
| `environment_name` | 2-60 alphanumerics/hyphens; embedded in every app's default FQDN. ForceNew |
| `logs_destination` + `log_analytics_workspace_id` | LOG_ANALYTICS pairs with a workspace; AZURE_MONITOR forbids one; unset with a workspace deploys log-analytics |
| `infrastructure_subnet_id` | `/21` or larger; unlocks ILB and zone redundancy. ForceNew |
| `workload_profiles` | Dedicated (D4-D32, E4-E32), GPU (NC*-A100), serverless GPU (Consumption-GPU-*); 0-to-some or some-to-0 transitions are ForceNew |
| `custom_domain` | Environment-wide DNS suffix + wildcard PFX certificate (ARM models it as part of the environment) |
| `identity` | SYSTEM_ASSIGNED / USER_ASSIGNED / both |

## Outputs

| Output | Purpose |
| --- | --- |
| `environment_id` | The ARM ID every family kind references |
| `environment_name` | The environment's name |
| `default_domain` | Apps live at `{app-name}.{default_domain}` -- the DNS CNAME seam |
| `static_ip_address` | Public (external) or private (internal) entry IP -- the DNS A-record seam |
| `platform_reserved_cidr` / `platform_reserved_dns_ip_address` / `docker_bridge_cidr` | Network planning for VNet-injected environments |
| `custom_domain_verification_id` | TXT-record value proving domain ownership |
| `identity_principal_id` | RBAC grant target for the system-assigned identity |

## Quick Example

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
    valueFrom:
      kind: AzureLogAnalyticsWorkspace
      name: my-law
      fieldPath: status.outputs.workspace_id
```

## Downstream Usage

```yaml
# AzureContainerApp running inside this environment
spec:
  containerAppEnvironmentId:
    valueFrom:
      kind: AzureContainerAppEnvironment
      name: my-env
      fieldPath: status.outputs.environment_id
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
