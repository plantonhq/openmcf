# Azure Container Apps Environment

Serverless microservices without cluster operations: a VNet-injected, zone-redundant Container App Environment running the three workload shapes every product team ends up needing — a public API, a queue-driven worker that scales to zero, and a scheduled batch job — with the messaging, storage, and identity seams between them wired keylessly.

## Who this is for

A product team shipping containers that wants Kubernetes-grade capabilities (Dapr sidecars, KEDA autoscaling, shared volumes, VNet privacy) with none of the cluster to run. Deploying this chart yields a working event-driven skeleton: the API publishes to a topic, the worker wakes from zero to consume it, the job runs on a schedule — replace the three quickstart images with your services and the architecture is already right.

## Architecture

```
                        internet (or VNet-only with internal_only_enabled)
                                        │
                                   ┌────▼────┐  Dapr publish   ┌──────────────────┐
              VNet /16             │   api   │ ───────────────▶│ Service Bus queue │
  ┌────────────────────────────┐   └─────────┘   (keyless)     │  ("orders", DLQ   │
  │  infrastructure subnet /21 │                               │   posture)        │
  │  Microsoft.App delegation  │   ┌─────────┐  Dapr subscribe └────────┬─────────┘
  │  Container App Environment │   │ worker  │ ◀───────────────────────┘
  │  (zone-redundant)          │   └────┬────┘  KEDA scales 0→N on queue depth
  └────────────────────────────┘        │ /data                 (keyless)
                                   ┌────▼────────────┐
                                   │ Azure Files SMB  │   ┌───────────────┐
                                   │ shared volume    │   │  batch job    │
                                   └──────────────────┘   │  (cron, UTC)  │
                                                          └───────────────┘
```

Everything authenticates as one user-assigned identity: the Dapr sidecars publish/consume Service Bus with Entra tokens (Data Sender + Data Receiver grants — never a connection string), and KEDA reads queue depth with the same identity. The Dapr component disables entity management, so the messaging topology stays owned by infrastructure-as-code: the topic's backing queue is a first-class resource with dead-letter posture, TTL, and lock duration you control.

## Resources

| Kind | Name | Purpose |
| --- | --- | --- |
| AzureResourceGroup | `{env}-container-apps` | One container for the whole estate |
| AzureLogAnalyticsWorkspace | `{env}-capps-logs` | Environment console/system logs + broker diagnostics |
| AzureVirtualNetwork | `{env}-capps-vnet` | The environment's own network |
| AzureSubnet | `{env}-capps-infra-subnet` | /21 infrastructure subnet, `Microsoft.App/environments`-delegated |
| AzureUserAssignedIdentity | `{env}-capps-identity` | The shared workload identity everything runs as |
| AzureRoleAssignment | `{env}-capps-bus-send` / `-receive` | Data Sender + Data Receiver on the namespace — the keyless grants |
| AzureServiceBusNamespace | `{env}-capps-bus` | The pub/sub broker (STANDARD) |
| AzureServiceBusQueue | `{env}-{topic}-queue` | The topic's backing queue, DLQ-postured |
| AzureMonitorDiagnosticSetting | `{env}-capps-bus-diag` | Broker telemetry into the workspace |
| AzureContainerAppEnvironment | `{env}-capps` | The VNet-injected, zone-redundant environment |
| AzureContainerAppEnvironmentDaprComponent | `{env}-capps-pubsub` | Keyless Service Bus pub/sub, scoped to api + worker |
| AzureStorageAccount | `{env}-capps-storage` | Backs the shared Azure Files volume |
| AzureStorageShare | `{env}-capps-share` | The SMB share (quota-capped) |
| AzureContainerAppEnvironmentStorage | `{env}-capps-shared-volume` | Registers the share for app volumes (key by reference) |
| AzureContainerApp | `{env}-api` | Public API, external ingress, Dapr publisher |
| AzureContainerApp | `{env}-worker` | Ingress-less consumer, scales 0→N on queue depth, mounts the volume |
| AzureContainerAppJob | `{env}-batch-job` | Cron-triggered run-to-completion work |

## Parameters

| Parameter | Description | Default | Must change |
| --- | --- | --- | --- |
| `region` | Azure region for the estate | `centralus` | |
| `vnet_cidr` | VNet address space | `10.20.0.0/16` | |
| `infra_subnet_cidr` | Infrastructure subnet (/21 or larger) | `10.20.0.0/21` | |
| `servicebus_namespace_name` | Globally unique namespace name | `my-capps-bus` | yes |
| `storage_account_name` | Globally unique storage account name | `mycappsshared` | yes |
| `zone_redundancy_enabled` | Zone-spread the environment (fixed at creation) | `true` | |
| `internal_only_enabled` | VNet-only exposure via internal load balancer | `false` | |
| `api_image` / `worker_image` / `job_image` | The three workload images | quickstart images | eventually |
| `worker_max_replicas` | Worker scale ceiling | `10` | |
| `job_cron_expression` | UTC schedule for the batch job | `0 3 * * *` | |
| `orders_topic_name` | The Dapr topic (= backing queue name) | `orders` | |
| `shared_volume_quota_gb` | Azure Files volume cap (GiB) | `64` | |
| `log_retention_days` | Workspace retention | `30` | |

## After deploying

1. **Point the API at your image** — set `api_image` (and the others) to your registry's images. For a private registry, add a `registries` block to the apps and grant the shared identity `AcrPull`.
2. **Publish and consume through Dapr** — the API publishes with `POST http://localhost:3500/v1.0/publish/pubsub/{orders_topic_name}`; the worker subscribes to the same topic through the Dapr subscription contract. Application code never sees a Service Bus SDK or connection string.
3. **Watch it scale** — enqueue a burst and watch the worker climb from zero (`az containerapp replica list`); the KEDA rule adds a replica per 10 pending messages.
4. **Inspect dead letters** — messages that exhaust 10 delivery attempts or expire past 14 days land in the queue's DLQ; drain it with Service Bus Explorer in the portal.

## Day 2

- **More services**: each new app gets its own Dapr `appId`; add it to the component's `scopes` to grant broker access deliberately.
- **Split identities**: when app permissions diverge, mint an identity per app and narrow each grant (queue-scoped instead of namespace-scoped).
- **Dedicated compute**: declare `workload_profiles` on the environment (D/E families, GPU) and pin heavy apps to them with `workload_profile_name` — note that going from zero profiles to some replaces the environment.
- **Premium messaging**: raise the namespace SKU to PREMIUM for dedicated capacity, CMK, and VNet network rules — an in-place change.
