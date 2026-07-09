# AzureContainerAppJob

Deploy an Azure Container App Job -- a run-to-completion containerized workload (batch processing, scheduled tasks, queue workers) inside a Container App Environment.

## Overview

Where a Container App runs continuously and serves traffic, a Job starts, does its work, and exits. Each execution runs the job's container template to completion (bounded by `replica_timeout_in_seconds`), with `parallelism` and `replica_completion_count` shaping how many replicas an execution runs and how many must succeed.

Exactly one trigger drives executions: manual (on demand from CLI/SDK/pipelines), schedule (UTC cron), or event (a KEDA scaler fans queue/broker pressure out into executions).

## Key Features

- **Three trigger models**: manual (`az containerapp job start`), cron schedule, and KEDA event triggers with execution fan-out control (min/max executions, polling interval)
- **The Container App container model**: main + init containers, probes with Azure's per-type contracts, env vars, entrypoint overrides, volume mounts
- **Four volume types**: EmptyDir scratch, SMB/NFS Azure Files shares (via `AzureContainerAppEnvironmentStorage`), and Secret mounts
- **Secrets**: plain values or Key Vault references read through a managed identity
- **Registries**: username/password or managed-identity pulls
- **Managed identity** and user tags; workload profile targeting for dedicated compute

## When to Use

- Database migrations gated on deploy pipelines (manual trigger)
- Nightly reports, exports, periodic cleanup (schedule trigger)
- Queue/message-batch processing where each batch is its own isolated execution (event trigger)
- CI runners and other bounded work that should not be modeled as an always-on service

## Spec Highlights

| Field | Notes |
| --- | --- |
| `job_name` | Max 32 lowercase alphanumerics/hyphens (no dots -- unlike app names). ForceNew |
| `replica_timeout_in_seconds` | The hard per-replica deadline; required |
| `manual_trigger` / `schedule_trigger` / `event_trigger` | Exactly one; switching types is ForceNew |
| `event_trigger.scale` | max/min executions, polling interval, KEDA rules (with optional `identity_id`) |
| `containers[]` | cpu/memory required; probes carry per-type threshold ceilings (30/48/240) |

## Outputs

| Output | Purpose |
| --- | --- |
| `job_id` / `job_name` | ARM ID (the handle for starting manual executions) and name |
| `event_stream_endpoint` | Execution event stream for monitoring |
| `outbound_ip_addresses` | Egress allowlisting |
| `identity_principal_id` | RBAC grant target for the system-assigned identity |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppJob
metadata:
  name: nightly-report
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
  replicaTimeoutInSeconds: 3600
  containers:
    - name: report
      image: myregistry.azurecr.io/report:v1.0.0
      cpu: 0.5
      memory: "1Gi"
  scheduleTrigger:
    cronExpression: "0 2 * * *"
```
