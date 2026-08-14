# Azure Data Factory Trigger

Deploys one trigger inside an Azure Data Factory -- the instruction that starts pipelines automatically, in any of Azure's four types: schedule, tumbling window, blob event, or custom event. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions exactly one of (the spec's variant block decides):

- **Schedule trigger** -- fires on a wall-clock recurrence, optionally narrowed to specific minutes/hours/days/monthly occurrences
- **Tumbling window trigger** -- fires once per contiguous time window from a fixed start; backfills past windows, retries, rate-limits, and can depend on other windows
- **Blob event trigger** -- fires when blobs are created/deleted in a storage account, filtered by path (wired through Event Grid behind the scenes)
- **Custom event trigger** -- fires on events published to an Event Grid custom topic, filtered by subject and event type

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Data Factory and Pipeline** -- the trigger lives in a factory and starts pipelines by name; reference an AzureDataFactory's ID output and AzureDataFactoryPipeline name outputs (or provide literals).
- **Blob event**: an AzureStorageAccount reference. **Custom event**: an AzureEventgridTopic reference.

### Azure Subscription

- **Blob event triggers require the Microsoft.EventGrid resource provider registered** -- Azure creates the Event Grid subscription on the storage account behind the scenes.
- **A started trigger RUNS PIPELINES** -- deploying with `activated: true` (the default) and a due schedule starts real, billed pipeline runs immediately. Set the start time deliberately.
- **Tumbling window `start_time` in the past backfills** -- every window between then and now runs (rate-limited by `max_concurrency`). That is a feature pointed at history and an incident pointed at the wrong date.

## Deploy

### Console

Open the deployment store, find **Azure Data Factory Trigger**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Daily Schedule** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f data-factory-trigger.yaml
```

## After Deploy

The trigger appears in the factory's Studio under Manage -> Triggers, showing its started/stopped state. Fired runs land on the Monitor blade's Trigger runs view. Flip `activated: false` and re-apply to pause the trigger without deleting it (in-flight pipeline runs finish; new firings stop).
