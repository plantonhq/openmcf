# Azure Data Factory

Deploys an Azure Data Factory -- the workspace every other Data Factory resource lives inside: pipelines, data flows, linked services, datasets, triggers, and integration runtimes are all created against it. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Data Factory** -- the workspace itself: managed identity, optional git repository binding, global parameters, network posture, and inline customer-managed-key encryption
- **Credentials** (optional) -- one per named entry, wrapping a user-assigned identity or a service principal whose key lives in Key Vault
- **Managed private endpoints** (optional) -- private egress from the factory's managed virtual network to data stores

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name.
- **For CMK encryption (optional)** -- an AzureKeyVaultKey (reference its versionless ID output) and an AzureUserAssignedIdentity with get/unwrap/wrap permissions on the vault, attached in the identity block.

### Azure Subscription

- **Factory names are globally unique across Azure** -- prefix with your org; a taken name fails at deploy time.
- **The managed virtual network is one-way** -- enabling it updates in place; disabling it REPLACES the factory. Managed private endpoints require it.
- **Private endpoint approval is on the TARGET side** -- each managed private endpoint's connection must be approved on the target resource before traffic flows.
- **The repository binding does not detach by removal** -- removing the block from the spec leaves the repo bound; detach in the Data Factory Studio.

## Deploy

### Console

Open the deployment store, find **Azure Data Factory**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Data Platform Workspace** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f data-factory.yaml
```

## After Deploy

The `data_factory_id` output is the wiring edge for everything inside the factory -- create AzureDataFactoryPipeline resources with `data_factory_id` referencing it. Grant the factory's identity (the `identity_principal_id` output) on the data stores it reads and writes; if managed private endpoints were created, approve their connections on each target's **Networking -> Private endpoint connections** blade before running pipelines through them.
