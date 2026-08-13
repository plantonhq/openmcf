# Azure AI Foundry Project

Creates an Azure AI Foundry project -- the workspace one AI team works in, created inside an AzureAiFoundry hub and inheriting its security, storage, and network posture. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **AI Foundry Project** -- ARM-wise an ML workspace of kind "Project" (`Microsoft.MachineLearningServices/workspaces`), linked to its hub and placed in the hub's resource group

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureAiFoundry hub** -- the project is created inside it (the linkage is fixed at creation).

### Azure Subscription

- **No resource group needed** -- the project deploys into the HUB's resource group; the provider derives it from the hub reference.

## Deploy

### Console

Open the deployment store, find **Azure AI Foundry Project**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Team Project** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureAiFoundryProject
metadata:
  name: fraud-detection
  org: acme-corp
  env: prod
spec:
  region: eastus
  name: fraud-detection
  aiServicesHubId:
    valueFrom:
      name: team-hub
  identity:
    type: SYSTEM_ASSIGNED
  friendlyName: Fraud Detection
```

```shell
planton apply -f azure-ai-foundry-project.yaml
```

The project provisions in a few minutes.

### InfraChart

In an AI-platform chart the order is: hub → **projects** (one per team), each wiring the hub by reference.

## Key Configuration

These are the most important decisions when configuring the project. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The hub linkage** -- `aiServicesHubId` is fixed at creation; a project cannot move between hubs. The project's region is its own field but typically matches the hub's.

**Identity** -- `SYSTEM_ASSIGNED` gives the team its own grantable identity. `primaryUserAssignedIdentity` is only legal alongside the identity block (validated at manifest time -- the provider's own pairing).

**Everything else is inherited** -- vault, storage, insights, registry, network isolation, and encryption all come from the hub; a project manifest stays a few lines by design.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureAiFoundry** | `aiServicesHubId` | `status.outputs.ai_foundry_id` |
| **AzureUserAssignedIdentity** | `identity.identityIds[]` | `status.outputs.identity_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `ai_foundry_project_id` | ARM ID of the project | Operational tooling |
| `ai_foundry_project_name` | The project's name | Foundry SDK targeting |
| `project_guid` | The project's immutable GUID | Data-plane calls, diagnostics |
| `system_assigned_identity_principal_id` | The system identity's principal ID | Per-team data grants |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**One project per team** -- the standard shape. Start from the **Team Project** preset.

## Works With

- [**Azure AI Foundry Hub**](/cloud-catalog/azure-ai-foundry) -- the hub this project lives in
- [**Azure Cognitive Account**](/cloud-catalog/azure-cognitive-account) -- the Azure OpenAI models the team's agents call
- [**Azure Search Service**](/cloud-catalog/azure-search-service) -- retrieval for the team's RAG applications
- [**Azure User Assigned Identity**](/cloud-catalog/azure-user-assigned-identity) -- bring-your-own identity for pre-composed grants
