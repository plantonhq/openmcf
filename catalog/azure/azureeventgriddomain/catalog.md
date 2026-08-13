# Azure Event Grid Domain

Deploys an Azure Event Grid domain -- one publishing endpoint and one pair of access keys serving many event streams (domain topics), the multi-tenant pattern. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Event Grid domain** -- the shared publish endpoint with its pair of access keys, the domain-topic lifecycle flags, optional managed identity, optional IP firewall, and the input schema every incoming event must match

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name.

### Azure Subscription

- **The domain name is a public DNS hostname** (`{name}.{region}.eventgrid.azure.net`) -- unique across ALL Azure customers in the region; a taken name fails the deploy with a conflict.
- **The input schema is chosen at creation** and applies to every topic in the domain -- changing it replaces the domain and its endpoint hostname.
- **Decide the topic lifecycle up front**: auto-managed (Azure's default -- topics appear with their first subscription and vanish with the last) or pinned (`auto_create`/`auto_delete` false -- topics exist only as declared AzureEventgridDomainTopic resources).
- **A domain is free at rest** -- billing is per operation (publishes, deliveries, filtering evaluations).

## Deploy

### Console

Open the deployment store, find **Azure Event Grid Domain**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Multi-Tenant Domain** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f eventgrid-domain.yaml
```

## After Deploy

The `endpoint` output is the single publish URL for every topic in the domain; events name their topic in the payload's topic field, and the `primary_access_key` output (a secret) authenticates POSTs via the `aeg-sas-key` header. Add topics as AzureEventgridDomainTopic resources referencing this domain's `domain_id` output (or let subscriptions auto-create them, per your lifecycle flags), and watch traffic on the domain's **Metrics** blade.
