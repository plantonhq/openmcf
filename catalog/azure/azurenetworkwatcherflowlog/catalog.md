# Azure Network Watcher Flow Log

Records network traffic metadata for one virtual network, subnet, or network interface into a storage account, optionally enriched by Traffic Analytics in a Log Analytics workspace. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Network Watcher flow log** -- the recorder, attached to the region's Network Watcher (which Azure auto-creates -- the module never creates a watcher)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **The target** -- an AzureVirtualNetwork, AzureSubnet, or AzureNetworkInterface whose traffic is recorded (referenced by `targetResourceId`).
- **An AzureStorageAccount** flow-log files land in -- one WITHOUT hand-managed lifecycle rules (creating a flow log overwrites them).
- **Optional: an AzureLogAnalyticsWorkspace** for Traffic Analytics.

### Azure Subscription

- **The flow log must live in the target's region** -- flow logging is regional; ARM rejects cross-region pairings.
- **Flow logs are near-free**; the costs that matter are storage (bounded by `retentionPolicy`) and, with Traffic Analytics, workspace ingestion.
- **NSG-targeted flow logs are retired for new creates** (since June 2025) -- target the network scope instead.

## Deploy

### Console

Open the deployment store, find **Azure Network Watcher Flow Log**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Virtual Network Flow Log** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f flow-log.yaml
```

## After Deploy

Flow-log files appear in the storage account under `insights-logs-flowlogflowevent` within minutes of traffic; with Traffic Analytics, processed flows land in the workspace on the configured interval (10 or 60 minutes). The `network_watcher_name` output reports which watcher the flow log attached to.
