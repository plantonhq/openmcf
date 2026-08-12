# Azure Event Grid Event Subscription

Deploys an Azure Event Grid event subscription -- the delivery instruction routing events from a source (custom topic, domain, domain topic, system topic, resource group, or subscription) to a handler, with filtering, retry, and dead-letter behavior. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions exactly one of:

- **Scope-addressed event subscription** -- attaches to any ARM resource that emits Event Grid events (when `scope` is set)
- **System-topic event subscription** -- a child of the referenced system topic (when `system_topic_id` is set)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **The event source** -- an AzureEventgridTopic (the default reference for `scope`), an AzureEventgridSystemTopic (for `system_topic_id`), or any other source's ARM ID.
- **The destination** -- the queue, hub, function, or webhook events are delivered to must exist first; Azure validates it at create time.

### Azure Subscription

- **Webhook destinations answer a validation handshake at create time** -- the endpoint must be live, or the create fails.
- **The advanced filter caps at 25 total values per subscription** (Azure's service limit across all conditions).
- **Delivery properties are ignored on storage-queue destinations** -- queue messages carry no custom properties.
- **Without a dead-letter destination, events that exhaust retries are dropped** -- configure `dead_letter` for at-least-once processing.
- **A subscription is free at rest** -- billing is per operation (deliveries, filtering evaluations).

## Deploy

### Console

Open the deployment store, find **Azure Event Grid Event Subscription**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Queue Fan-Out** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f eventgrid-event-subscription.yaml
```

## After Deploy

Publish (or trigger) an event on the source and watch it arrive at the destination -- the subscription's **Metrics** blade counts matched, delivered, and dead-lettered events. The `event_subscription_id` output records which addressing shape was created (a scoped path or a system-topic child path).
