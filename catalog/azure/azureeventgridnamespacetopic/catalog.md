# Azure Event Grid Namespace Topic

Deploys one named CloudEvents stream inside an Azure Event Grid namespace. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Event Grid namespace topic** -- one CloudEvents stream inside the referenced namespace, with its delivery retention window

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Event Grid Namespace** -- reference an AzureEventgridNamespace's ID output or provide an existing namespace's ARM ID.

### Azure Subscription

- **The schema and publisher type are fixed** -- CloudEvents v1.0 and "Custom"; Azure defines no alternatives on this resource today.
- **Retention is the only update** -- 1-7 days; changing the name or namespace replaces the topic.
- **Topics share the namespace's capacity** -- a topic adds no cost of its own; throughput bills on the namespace's throughput units.

## Deploy

### Console

Open the deployment store, find **Azure Event Grid Namespace Topic**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Service Stream** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f eventgrid-namespace-topic.yaml
```

## After Deploy

Publishers POST CloudEvents to the namespace's HTTP endpoint addressed at this topic. Subscriptions on namespace topics are managed outside the catalog until the pinned provider ships the resource -- the GUIDE carries the boundary.
