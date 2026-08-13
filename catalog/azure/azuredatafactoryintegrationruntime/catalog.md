# Azure Data Factory Integration Runtime

Deploys one integration runtime inside an Azure Data Factory -- the compute engine the factory's pipelines, data flows, and copy activities run on, in any of three flavors: the managed data-flow compute, the managed SSIS package runtime, or the self-hosted agent registration. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions exactly one integration runtime of the flavor the spec's variant block declares:

- **Azure (data-flow compute)** -- serverless Spark Azure provisions when a mapping data flow runs, sized by compute type and core count, optionally kept warm and joined to the factory's managed virtual network
- **Azure-SSIS** -- a managed cluster of VMs that runs SQL Server Integration Services packages, with an optional SSISDB catalog, node custom setup, virtual network injection, package stores, and an on-premises proxy
- **Self-hosted** -- the registration for the agent you install on your own machines; Azure issues the authorization keys the agent joins with (surfaced as sensitive outputs)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Data Factory** -- the runtime lives in a factory; reference an AzureDataFactory's ID output (or provide a literal).
- **For the azure flavor's virtual network switch** -- the factory must be deployed with its managed virtual network enabled; Azure rejects the runtime otherwise.

### Azure Subscription

- **Billing follows the flavor**: the data-flow compute bills per vCore-hour while clusters run (plus any warm time-to-live); the SSIS runtime is created STOPPED and unbilled -- node-hours bill only after you start it; the self-hosted registration is free on Azure's side.
- **The SSIS runtime does not start itself** -- starting and stopping is an operational action in Data Factory Studio or via the Data Factory API, deliberately outside this definition.
- **A self-hosted registration needs an agent**: creating it issues keys but moves no data until you install the integration runtime agent on your machines and hand it a key.

## Deploy

### Console

Open the deployment store, find **Azure Data Factory Integration Runtime**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Data Flow Compute** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f data-factory-integration-runtime.yaml
```

## After Deploy

The runtime appears in the factory's Studio under Manage -> Integration runtimes. For the self-hosted flavor, copy an authorization key from the component's sensitive outputs and hand it to the agent installer on your machine; the node shows Running in Studio once it joins. For SSIS, start the runtime from Studio when you are ready to run packages -- and stop it when you are not, because node-hours bill while it runs.
