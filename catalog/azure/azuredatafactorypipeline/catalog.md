# Azure Data Factory Pipeline

Deploys one pipeline inside an Azure Data Factory -- an ordered set of activities that executes as a whole when triggered. The activities travel as raw JSON (the Studio's Code view artifact). It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Data Factory pipeline** -- the pipeline definition: activities JSON, run parameters, variables, annotations, concurrency, Studio folder, and the elapsed-time metric

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Data Factory** -- reference an AzureDataFactory's ID output (or provide an existing factory's ARM ID).

### Azure Subscription

- **Azure owns the activity schema** -- invalid JSON or unknown activity shapes fail at deploy time, not at validation time. Author activities in the Data Factory Studio and copy its Code view.
- **Parameters and variables are String-typed on this surface** -- declare other-typed parameters inside the activities JSON if you need them.
- **Linked services and datasets the activities reference must already exist** in the factory, or runs fail at runtime.

## Deploy

### Console

Open the deployment store, find **Azure Data Factory Pipeline**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Scheduled Copy** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f data-factory-pipeline.yaml
```

## After Deploy

The pipeline appears in the factory's Studio under its folder; run it on demand ("Debug" or "Trigger now") or attach a trigger naming it. Monitor runs on the factory's **Monitor** blade -- the `monitor_metrics_after_duration` threshold, when set, fires the elapsed-time metric for runs that exceed it.
