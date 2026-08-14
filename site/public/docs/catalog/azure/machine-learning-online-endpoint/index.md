---
title: "Machine Learning Online Endpoint"
description: "Machine Learning Online Endpoint deployment documentation"
icon: "package"
order: 100
componentName: "azuremachinelearningonlineendpoint"
---

# Azure Machine Learning Online Endpoint

Creates a managed online endpoint on an Azure Machine Learning workspace -- the stable HTTPS address applications call to score against deployed models, with authentication and a traffic dial that splits requests across the endpoint's deployments. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Online Endpoint** -- an ARM child of the workspace (`.../workspaces/{ws}/onlineEndpoints/{name}`) with its auth mode, managed identity, traffic maps, and network posture

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureMachineLearningWorkspace** -- the endpoint is created on it.

### Azure Subscription

- **The endpoint name is reserved region-wide** -- per subscription, across all workspaces in the region; it becomes part of the scoring DNS.
- **The endpoint alone serves nothing** -- scoring works once an AzureMachineLearningOnlineDeployment attaches to it and the traffic map routes to it.

## Deploy

### Console

Open the deployment store, find **Azure Machine Learning Online Endpoint**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Key-Auth Endpoint** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMachineLearningOnlineEndpoint
metadata:
  name: fraud-scoring
  org: acme-corp
  env: prod
spec:
  workspaceId:
    valueFrom:
      name: ml-prod
  name: fraud-scoring
  region: eastus
  authMode: Key
  identity:
    type: SYSTEM_ASSIGNED
  traffic:
    blue: 100
```

```shell
planton apply -f azure-machine-learning-online-endpoint.yaml
```

The endpoint creates in a few minutes; deployments then attach to it by reference.

### InfraChart

In an ML-platform chart the order is: workspace → **online endpoint** → online deployment(s), each wiring its parent by reference; the endpoint's traffic map routes to deployments by their names.

## Key Configuration

These are the most important decisions when configuring the endpoint. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Auth mode** -- `Key` (static keys, never expire), `AMLToken` (expiring Azure ML tokens), or `AADToken` (Microsoft Entra tokens -- the keyless posture). Fixed per endpoint contract with callers; changing it in place re-issues credentials.

**Traffic** -- whole-percent shares per deployment name, summing to at most 100 (ARM enforces the sum at apply). A new endpoint usually starts `{blue: 100}`; rollouts move the map in steps. `mirrorTraffic` (at most 50 total) shadow-tests without answering callers.

**Identity** -- required. `SYSTEM_ASSIGNED` covers most cases; grant it AcrPull on the registry and Storage Blob Data Reader on the model store so deployments can pull.

**Initial auth keys** -- Key mode only: bring your own keys from a secret store, or leave unset and let the service mint them (read with `az ml online-endpoint get-credentials`).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureMachineLearningWorkspace** | `workspaceId` | `status.outputs.machine_learning_workspace_id` |
| **AzureUserAssignedIdentity** | `identity.identityIds[]` | `status.outputs.identity_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `online_endpoint_id` | ARM ID of the endpoint | What AzureMachineLearningOnlineDeployment wires to |
| `online_endpoint_name` | The endpoint's name | Traffic-map keys, operational tooling |
| `scoring_uri` | The HTTPS scoring address | Application configuration |
| `swagger_uri` | The OpenAPI document address | Client generation |
| `system_assigned_identity_principal_id` | The system identity's principal ID | Registry / storage role assignments |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Key-auth endpoint** -- the everyday shape: static keys, system identity, all traffic to one deployment. Start from the **Key-Auth Endpoint** preset.

**Entra-auth private endpoint** -- keyless authentication with public network access off. Start from the **Entra-Auth Private Endpoint** preset.

## Works With

- [**Azure Machine Learning Workspace**](/cloud-catalog/azure-machine-learning-workspace) -- the parent workspace
- [**Azure Machine Learning Online Deployment**](/cloud-catalog/azure-machine-learning-online-deployment) -- the model deployments behind the endpoint
- [**Azure User Assigned Identity**](/cloud-catalog/azure-user-assigned-identity) -- bring-your-own identity for registry and storage grants
