---
title: "Container App Job"
description: "Container App Job deployment documentation"
icon: "package"
order: 100
componentName: "azurecontainerappjob"
---

# Azure Container App Job

Deploys a run-to-completion containerized workload -- batch processing, scheduled tasks, queue workers -- inside an Azure Container App Environment, with manual, cron, or KEDA event triggers.

## What Gets Created

When you deploy an AzureContainerAppJob resource, Planton provisions:

- **Container App Job** -- an `azurerm_container_app_job` (`Microsoft.App/jobs`) in the referenced environment, carrying the container template, the trigger configuration, secrets, registries, and identity

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureContainerAppEnvironment** to run in (referenced through `containerAppEnvironmentId`)
- For private images: registry credentials or a managed identity with pull rights
- For share-backed volumes: an `AzureContainerAppEnvironmentStorage` registration

## Quick Start

Create a file `job.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppJob
metadata:
  name: nightly-report
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureContainerAppJob.nightly-report
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  jobName: nightly-report
  containerAppEnvironmentId:
    valueFrom:
      kind: AzureContainerAppEnvironment
      name: my-env
      fieldPath: status.outputs.environment_id
  replicaTimeoutInSeconds: 1800
  containers:
    - name: report
      image: myregistry.azurecr.io/report:v1
      cpu: 0.5
      memory: "1Gi"
  scheduleTrigger:
    cronExpression: "0 2 * * *"
```

Deploy:

```shell
planton apply -f job.yaml
```

Exactly one trigger drives executions: `manualTrigger` (on demand -- `az containerapp job start`), `scheduleTrigger` (UTC cron), or `eventTrigger` (a KEDA scaler fans queue/broker pressure into executions, with `scale` controlling min/max executions and polling). Switching trigger types recreates the job.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `job_id` / `job_name` | The ARM ID (the handle for starting manual executions) and name |
| `event_stream_endpoint` | The execution event stream for monitoring |
| `outbound_ip_addresses` | Egress allowlisting on external services |
| `identity_principal_id` | The RBAC grant target for the system-assigned identity |

## Related Resources

- [Azure Container App Environment](/docs/catalog/azure/container-app-environment) -- the hosting boundary
- [Azure Container App](/docs/catalog/azure/container-app) -- the continuously running sibling
- [Azure Container App Environment Storage](/docs/catalog/azure/container-app-environment-storage) -- share-backed volumes
