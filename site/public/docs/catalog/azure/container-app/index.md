---
title: "Container App"
description: "Container App deployment documentation"
icon: "package"
order: 100
componentName: "azurecontainerapp"
---

# Azure Container App

Deploy a serverless, continuously running containerized service inside an Azure Container App Environment -- with KEDA scaling, revision-based deployments, and full ingress control.

## What Gets Created

- An Azure Container App (`Microsoft.App/containerApps`) with its first revision

## Prerequisites

- An `AzureContainerAppEnvironment` (the hosting boundary)
- For private images: registry credentials or a managed identity with pull rights
- For share-backed volumes: an `AzureContainerAppEnvironmentStorage` registration

## Quick Start

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerApp
metadata:
  name: my-web-service
spec:
  resourceGroup:
    value: my-rg
  containerAppName: my-web-service
  containerAppEnvironmentId:
    value: /subscriptions/.../managedEnvironments/my-env
  containers:
    - name: web
      image: mcr.microsoft.com/k8se/quickstart:latest
      cpu: 0.25
      memory: "0.5Gi"
  ingress:
    externalEnabled: true
    targetPort: 8080
    trafficWeight:
      - latestRevision: true
        percentage: 100
```

## Configuration Reference

### Required Fields

| Field | Description |
| --- | --- |
| `resourceGroup` | Resource group name or `AzureResourceGroup` reference (ForceNew) |
| `containerAppName` | Max 32 lowercase alphanumerics/hyphens/dots; part of the FQDN (ForceNew) |
| `containerAppEnvironmentId` | The hosting environment's ARM ID (ForceNew) |
| `containers` | At least one container with name, image, cpu, and memory |

### Optional Fields

| Field | Description |
| --- | --- |
| `revisionMode` | `SINGLE` (default) or `MULTIPLE` (traffic splitting) |
| `workloadProfileName` | Environment workload profile; omit for serverless Consumption |
| `maxInactiveRevisions` | Retention for inactive revisions (0-100) |
| `initContainers` | Run to completion before main containers; no probes |
| `volumes` | `EMPTY_DIR`, `AZURE_FILE`, `NFS_AZURE_FILE` (with `storageName`), or `SECRET` |
| `minReplicas` / `maxReplicas` | Scale bounds (0-300 / 1-300); 0 enables scale-to-zero |
| `cooldownPeriodInSeconds` / `pollingIntervalInSeconds` | KEDA scaler dials (300 / 30) |
| `revisionSuffix` / `terminationGracePeriodSeconds` | Revision pinning and shutdown grace |
| `httpScaleRules` / `tcpScaleRules` / `azureQueueScaleRules` / `customScaleRules` | The KEDA scaling surface |
| `secrets` | Plain values or Key Vault references (with `identity`) |
| `registries` | Managed identity or username + password-secret pulls |
| `ingress` | Exposure, transports, traffic weights, mTLS modes, CORS, IP restrictions |
| `dapr` | Sidecar injection (`appId`, `appPort`, `appProtocol`) |
| `identity` | `SYSTEM_ASSIGNED`, `USER_ASSIGNED` (+ identity references), or both |
| `tags` | Free-form Azure tags merged over platform tags |

## Examples

### Canary Rollout with Traffic Splitting

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerApp
metadata:
  name: my-api
spec:
  resourceGroup:
    value: my-rg
  containerAppName: my-api
  containerAppEnvironmentId:
    value: /subscriptions/.../managedEnvironments/my-env
  revisionMode: MULTIPLE
  revisionSuffix: v2
  containers:
    - name: api
      image: myregistry.azurecr.io/api:v2.0.0
      cpu: 0.5
      memory: "1Gi"
  ingress:
    externalEnabled: true
    targetPort: 8080
    trafficWeight:
      - revisionSuffix: v1
        percentage: 90
      - revisionSuffix: v2
        percentage: 10
        label: canary
```

### Queue Worker Scaling to Zero

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerApp
metadata:
  name: my-worker
spec:
  resourceGroup:
    value: my-rg
  containerAppName: my-worker
  containerAppEnvironmentId:
    value: /subscriptions/.../managedEnvironments/my-env
  containers:
    - name: worker
      image: myregistry.azurecr.io/worker:v1.0.0
      cpu: 0.25
      memory: "0.5Gi"
      env:
        - name: QUEUE_CONNECTION
          secretName: queue-conn
  minReplicas: 0
  maxReplicas: 10
  azureQueueScaleRules:
    - name: queue-depth
      queueName: work
      queueLength: 5
      authentication:
        - secretName: queue-conn
          triggerParameter: connection
  secrets:
    - name: queue-conn
      value: DefaultEndpointsProtocol=https;AccountName=...
```

### Persistent Share Mount + Key Vault Secret

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerApp
metadata:
  name: my-cms
spec:
  resourceGroup:
    value: my-rg
  containerAppName: my-cms
  containerAppEnvironmentId:
    value: /subscriptions/.../managedEnvironments/my-env
  containers:
    - name: cms
      image: myregistry.azurecr.io/cms:v1.0.0
      cpu: 1.0
      memory: "2Gi"
      env:
        - name: DB_PASSWORD
          secretName: db-password
      volumeMounts:
        - name: content
          path: /var/content
  volumes:
    - name: content
      storageType: AZURE_FILE
      storageName:
        valueFrom:
          kind: AzureContainerAppEnvironmentStorage
          name: cms-content
          fieldPath: status.outputs.storage_name
  secrets:
    - name: db-password
      keyVaultSecretId: https://my-vault.vault.azure.net/secrets/db-password
      identity: System
  identity:
    type: SYSTEM_ASSIGNED
  ingress:
    externalEnabled: true
    targetPort: 8080
    trafficWeight:
      - latestRevision: true
        percentage: 100
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `container_app_id` / `container_app_name` | ARM ID and name |
| `ingress_fqdn` | The app's main endpoint (empty without ingress) |
| `latest_revision_name` / `latest_revision_fqdn` | Latest-revision handle and direct FQDN |
| `outbound_ip_addresses` | Egress IPs for external allowlists |
| `custom_domain_verification_id` | TXT-record value for custom-domain binding |
| `identity_principal_id` | System-assigned identity's RBAC principal |

## Related Components

- [Azure Container App Environment](/docs/catalog/azure/container-app-environment) -- the hosting boundary
- [Azure Container App Job](/docs/catalog/azure/container-app-job) -- run-to-completion workloads
- [Azure Container App Environment Storage](/docs/catalog/azure/container-app-environment-storage) -- share-backed volumes
- [Azure Container App Environment Dapr Component](/docs/catalog/azure/container-app-environment-dapr-component) -- Dapr backends
- [Azure User Assigned Identity](/docs/catalog/azure/user-assigned-identity) -- keyless auth for registries, Key Vault, and scalers
