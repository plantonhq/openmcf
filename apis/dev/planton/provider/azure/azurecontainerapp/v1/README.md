# AzureContainerApp

Deploy an Azure Container App -- a serverless, continuously running containerized service inside a Container App Environment.

## Overview

A Container App defines the workload: which containers run, how they scale (KEDA rules from HTTP concurrency to any KEDA scaler), how they are exposed (HTTP/TCP ingress with traffic splitting, CORS, mTLS modes, IP restrictions), and what secrets and identities they use. Revisions are the deployment unit -- every template change creates one, and MULTIPLE revision mode splits traffic across them for blue-green and canary rollouts.

For run-to-completion work (batch, scheduled tasks, queue draining), use `AzureContainerAppJob` instead.

## Key Features

- **Full container model**: main + init containers, probes (liveness/readiness/startup with Azure's per-type contracts), env vars (literal or secret-backed), entrypoint overrides, volume mounts
- **Four volume types**: ephemeral EmptyDir, SMB and NFS Azure Files shares (through `AzureContainerAppEnvironmentStorage` registrations), and Secret mounts
- **KEDA scaling**: HTTP/TCP concurrency, Azure Queue depth, and the full KEDA scaler catalog via custom rules -- including scalers executing under a managed identity
- **Ingress depth**: external/internal exposure, auto/HTTP/HTTP2/TCP transports, revision traffic weights, client certificate modes (mTLS), CORS, IP allow/deny lists
- **Secrets**: plain values or Key Vault references read just-in-time through a managed identity
- **Registries**: username/password or managed-identity pulls (ACR keyless)
- **Dapr**: sidecar injection wired to environment-registered components
- **Managed identity** (system- and/or user-assigned) and user tags

## When to Use

- HTTP APIs and web services that should scale to zero when idle
- Event-driven microservices scaling on queue/broker depth
- gRPC services (HTTP2 transport) and raw TCP protocols (TCP transport + exposed port)
- Dapr-based distributed applications

## Spec Highlights

| Field | Notes |
| --- | --- |
| `container_app_name` | Max 32 lowercase alphanumerics/hyphens/dots; part of the FQDN. ForceNew |
| `revision_mode` | SINGLE (default) or MULTIPLE (traffic splitting) |
| `containers[]` | cpu/memory required; probes carry per-type threshold ceilings (30/48/240) |
| `volumes[]` | EMPTY_DIR / AZURE_FILE / NFS_AZURE_FILE / SECRET; share-backed types pair with `storage_name` |
| `secrets[]` | value XOR `key_vault_secret_id` + `identity` |
| `registries[]` | identity XOR username + `password_secret_name` |
| `ingress` | requires `traffic_weight`; `exposed_port` only with TCP transport |
| `workload_profile_name` | Targets an environment workload profile; omit for Consumption |

## Outputs

| Output | Purpose |
| --- | --- |
| `container_app_id` / `container_app_name` | ARM ID and name |
| `ingress_fqdn` | The app's endpoint (empty without ingress) |
| `latest_revision_name` / `latest_revision_fqdn` | CD verification and direct-revision access |
| `outbound_ip_addresses` | Egress allowlisting |
| `custom_domain_verification_id` | TXT-record value for custom-domain binding |
| `identity_principal_id` | RBAC grant target for the system-assigned identity |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerApp
metadata:
  name: my-api
spec:
  resourceGroup:
    value: my-rg
  containerAppName: my-api
  containerAppEnvironmentId:
    valueFrom:
      kind: AzureContainerAppEnvironment
      name: my-env
      fieldPath: status.outputs.environment_id
  containers:
    - name: api
      image: myregistry.azurecr.io/api:v1.0.0
      cpu: 0.5
      memory: "1Gi"
  ingress:
    externalEnabled: true
    targetPort: 8080
    trafficWeight:
      - latestRevision: true
        percentage: 100
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
