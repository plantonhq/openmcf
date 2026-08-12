# Azure Data Factory Data Flow

Deploys one data flow inside an Azure Data Factory -- a visually-designed transformation that runs on managed Spark, in either of Azure's two forms: a runnable mapping data flow or a reusable flowlet snippet. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions exactly one of:

- **Mapping data flow** (`flowlet: false`, the default) -- the complete, runnable transformation: script, named sources and sinks, intermediate transformations
- **Flowlet** (`flowlet: true`) -- the reusable snippet other data flows embed; sources and sinks are optional because the embedding flow supplies them

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Data Factory** -- reference an AzureDataFactory's ID output (or provide an existing factory's ARM ID).

### Azure Subscription

- **Author the script in the Data Factory Studio** -- the data flow script language is Azure's own; build the flow visually, open the "Script" view, and ship that artifact. Invalid scripts fail at deploy time.
- **Datasets and linked services the flow references must already exist** in the factory -- they are validated when the flow is saved or run, not by the catalog.
- **Both forms share one name namespace** -- a mapping data flow and a flowlet cannot carry the same name in one factory, and flipping the form replaces the object.

## Deploy

### Console

Open the deployment store, find **Azure Data Factory Data Flow**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Mapping Transformation** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f data-factory-data-flow.yaml
```

## After Deploy

The data flow appears in the factory's Studio under its folder. Run it through a pipeline's Execute Data Flow activity (data flows do not run standalone); debug it interactively in the Studio's data flow debug session -- note debug sessions bill per vCore-hour while active.
