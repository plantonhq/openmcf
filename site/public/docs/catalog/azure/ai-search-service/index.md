---
title: "AI Search Service"
description: "AI Search Service deployment documentation"
icon: "package"
order: 100
componentName: "azuresearchservice"
---

# Azure AI Search Service

Creates an Azure AI Search service -- the managed search-and-retrieval engine AI applications index and query their own data with (keyword, vector, and semantic search; the standard retrieval companion to Azure OpenAI). It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Search Service** -- the service itself (`Microsoft.Search/searchServices`) with its SKU, capacity, auth posture, network controls, and identity
- **Shared Private Links** (optional) -- one ARM child per `sharedPrivateLinkServices` entry (`.../sharedPrivateLinkResources/{name}`), giving indexers private reach to data sources behind private endpoints

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureResourceGroup** -- the service lives in one.

### Azure Subscription

- **The name is GLOBALLY unique** -- it forms the endpoint `{name}.search.windows.net`; a name-taken error means a genuine global collision.
- **Higher tiers need quota approval** -- `standard2`, `standard3`, and both storage-optimized SKUs require a quota increase request to Microsoft before ARM accepts them.
- **One free service per subscription** -- the `free` SKU is a shared-cluster evaluation tier.

## Deploy

### Console

Open the deployment store, find **Azure AI Search Service**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production Search** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureSearchService
metadata:
  name: acme-search
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: ai-platform-rg
  name: acme-search-prod
  sku: standard
  replicaCount: 3
  identity:
    type: SYSTEM_ASSIGNED
```

```shell
planton apply -f azure-search-service.yaml
```

The service provisions in a few minutes.

### InfraChart

In an AI-platform chart the order is: resource group → **search service** alongside the Azure OpenAI account; applications consume `endpoint` plus a key (or Entra RBAC), and indexers reach data sources via the service identity or shared private links.

## Key Configuration

These are the most important decisions when configuring the service. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Capacity** -- SKU × partitions × replicas. Partitions scale storage/indexing, replicas scale queries and availability (3+ replicas = 99.9% read-write SLA). Counts resize in place; the SKU upgrades in place ONLY along basic → standard → standard2 → standard3 (anything else replaces the service).

**Auth posture** -- default is API keys. Set `authenticationFailureMode` to add Entra RBAC alongside keys; set `localAuthenticationEnabled: false` for RBAC-only (keys stop working, and the failure mode must stay unset -- validated at manifest time).

**Semantic ranking** -- `semanticSearchSku: free` for evaluation (1000 requests/month), `standard` for production RAG; not available on the free service SKU.

**Shared private links** -- each entry sits "Pending" until the target resource's owner approves the connection; creation itself never needs the target side's consent.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |
| **AzureUserAssignedIdentity** | `identity.identityIds[]` | `status.outputs.identity_id` |
| *(any private-linkable kind)* | `sharedPrivateLinkServices[].targetResourceId` | named explicitly in `valueFrom` (storage, SQL, Key Vault, ...) |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `search_service_id` | ARM ID of the service | Diagnostic settings, RBAC scopes |
| `search_service_name` | The service's name | Operational tooling |
| `endpoint` | `https://{name}.search.windows.net` | What applications and SDKs call |
| `primary_key` / `secondary_key` | Admin API keys (sensitive) | Application config via managed secrets |
| `default_query_key` | The built-in read-only query key (sensitive) | Client-side query access |
| `customer_managed_key_encryption_compliance_status` | CMK compliance readback | Governance checks |
| `system_assigned_identity_principal_id` | The system identity's principal ID | Data-source grants for indexers |
| `shared_private_link_service_ids` | Per-link ARM IDs, keyed by name | Operational tooling |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production search** -- standard SKU, SLA replicas, RBAC alongside keys. Start from the **Production Search** preset.

**Development** -- the cheap basic tier. Start from the **Dev Basic Search** preset.

**RAG retrieval** -- semantic ranking for Azure OpenAI grounding. Start from the **Semantic RAG Search** preset.

## Works With

- [**Azure Cognitive Account**](/cloud-catalog/azure-cognitive-account) -- the Azure OpenAI models RAG applications pair with this service
- [**Azure AI Foundry Hub**](/cloud-catalog/azure-ai-foundry) -- Foundry projects consume search for agent retrieval
- [**Azure Storage Account**](/cloud-catalog/azure-storage-account) -- the classic indexer data source (and shared-private-link target)
- [**Azure Resource Group**](/cloud-catalog/azure-resource-group) -- where the service lives
