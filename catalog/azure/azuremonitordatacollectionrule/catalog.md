# Azure Monitor Data Collection Rule

Deploys an Azure Monitor data collection rule (DCR) -- the routing table declaring what telemetry the Azure Monitor Agent collects, where it lands, and how the two wire together. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Data collection rule** -- the collection policy with its data sources, destinations, data flows, custom-stream declarations, and (optionally) a managed identity

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **A resource group** -- the rule lives in a referenced resource group.
- **At least one destination** -- typically an AzureLogAnalyticsWorkspace (reference its `workspace_id` output); Event Hub and storage destinations reference their own kinds' outputs.

### Azure Subscription

- **The rule alone collects nothing** -- machines must be attached with `AzureMonitorDataCollectionRuleAssociation` resources, and the Azure Monitor Agent must run on them.
- **Names wire the rule together** -- data flows reference destinations by their rule-local names; a flow naming a missing destination is rejected at deploy time, as are duplicate destination names (one namespace across all arms).
- **Custom streams need a Data Collection Endpoint** -- stream declarations (custom log files, direct ingestion) require `data_collection_endpoint_id`; provide the literal ARM id of a DCE.
- **Platform compatibility is enforced by Azure at deploy time** -- a `Linux` rule cannot carry Windows event logs, a `Windows` rule cannot carry syslog, and the `*_direct` destinations require kind `AgentDirectToStore`.
- **Once `kind` is set, changing it replaces the rule.**
- **The rule object itself is free** -- you pay for the telemetry it lands (workspace ingestion, storage, Event Hub throughput).

## Deploy

### Console

Open the deployment store, find **Azure Monitor Data Collection Rule**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Linux Syslog and Performance to Workspace** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f data-collection-rule.yaml
```

## After Deploy

Attach machines with `AzureMonitorDataCollectionRuleAssociation` resources, then query the destination workspace (`Syslog`, `Perf`, or your custom table) -- first records typically land a few minutes after the agent picks up the association. If nothing arrives, check that the machine runs the Azure Monitor Agent and that the flow's stream names match the data sources' streams.
