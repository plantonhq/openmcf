# Azure Event Grid System Topic

Deploys an Azure Event Grid system topic -- the subscription surface for events Azure itself publishes about one of your resources (a storage account's blob events, a resource group's lifecycle events). It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Event Grid system topic** -- the source binding that makes an Azure service's built-in event stream subscribable, with optional managed identity for secured delivery

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name.
- **The source resource** -- the storage account, Key Vault, or other resource whose events you want; reference its ID output or pass a literal ARM ID.

### Azure Subscription

- **One system topic per source resource per topic type** -- a second create against the same source fails; teams sharing a source share its topic and attach their own subscriptions.
- **The region must match the source's region** -- global sources (subscriptions via `Microsoft.Resources.Subscriptions`, resource groups via `Microsoft.Resources.ResourceGroups`) require `Global`.
- **The source binding is create-only** -- changing the source, type, name, or region replaces the topic and drops every subscription attached to it.
- **A system topic is free at rest** -- billing is per operation (deliveries, filtering evaluations).

## Deploy

### Console

Open the deployment store, find **Azure Event Grid System Topic**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Storage Account Events** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f eventgrid-system-topic.yaml
```

## After Deploy

The `system_topic_id` output is the wiring edge for subscriptions -- create an AzureEventgridEventSubscription with `system_topic_id` referencing it and events start flowing to the handler. Until a subscription exists, the source's events are evaluated and dropped; watch delivery counts on the topic's **Metrics** blade.
