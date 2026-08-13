# Azure Cognitive Account Project

Deploys an AI Foundry project onto an Azure AI services account -- the workspace a team organizes its AI work in (agents, evaluations, files), isolated from sibling projects on the same account. The parent account must be kind `AIServices` with project management enabled. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **AI Foundry project** -- an ARM child of the account (`.../accounts/{account}/projects/{name}`) with its own managed identity, description, display name and tags

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureCognitiveAccount** of kind `AIServices` with `projectManagementEnabled: true` (which requires the account to carry a managed identity) -- the **AI Foundry Account** preset of that component is exactly this shape.

### Azure Subscription

- Nothing beyond the account: projects are free workspace objects and provision in seconds to a couple of minutes.

## Deploy

### Console

Open the deployment store, find **Azure Cognitive Account Project**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Team Workspace** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureCognitiveAccountProject
metadata:
  name: customer-support
  org: acme-corp
  env: prod
spec:
  cognitiveAccountId:
    valueFrom:
      kind: AzureCognitiveAccount
      name: foundry-prod
      fieldPath: status.outputs.cognitive_account_id
  name: customer-support
  region: eastus
  identity:
    type: SYSTEM_ASSIGNED
  displayName: Customer Support
```

```shell
planton apply -f azure-cognitive-account-project.yaml
```

### InfraChart

In an AI-platform chart the order is: account (AIServices, project management on) → **projects** (one per team), each wiring to the account by reference.

## Key Configuration

These are the most important decisions when configuring the project. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Identity** -- required. `SYSTEM_ASSIGNED` is the zero-management default; bring user-assigned identities when grants must exist before the project does.

**Description and display name** -- set them at creation if you want them at all: ARM cannot UPDATE either to an empty value, so clearing one later replaces the project.

**Naming** -- 2-64 characters, alphanumeric start; the name is the ARM child segment and replaces the project when changed.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureCognitiveAccount** | `cognitiveAccountId` | `status.outputs.cognitive_account_id` |
| **AzureUserAssignedIdentity** (optional) | `identity.identityIds` | `status.outputs.identity_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `project_id` | ARM ID of the project | Operational tooling |
| `project_name` | The project's name | SDK / agent configuration |
| `endpoints` | Data-plane endpoints keyed by service label | What agents and SDKs call |
| `is_default` | Whether this became the account's default project | Automation logic |
| `system_assigned_identity_principal_id` | The project identity's principal ID | Storage / search role grants |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Team workspace** -- a system-identity project per team. Start from the **Team Workspace** preset.

**Pre-granted identity** -- a user-assigned identity whose storage/search grants exist before the project. Start from the **User-Assigned Identity** preset.

## Works With

- [**Azure Cognitive Account**](/cloud-catalog/azure-cognitive-account) -- the parent AIServices account
- [**Azure Cognitive Deployment**](/cloud-catalog/azure-cognitive-deployment) -- the models projects consume
- [**Azure User Assigned Identity**](/cloud-catalog/azure-user-assigned-identity) -- pre-granted project identities
