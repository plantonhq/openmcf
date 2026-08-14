# Azure Monitor Data Collection Rule Association

Attaches one machine (VM, VM scale set, or Arc-enabled server) to an Azure Monitor data collection rule or data collection endpoint -- the resource that actually puts a machine under monitoring. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Data collection rule association** -- an extension resource on the target machine binding it to a rule (or, in the endpoint form, to a Data Collection Endpoint for configuration access)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **A target machine** -- reference an AzureVirtualMachine's `vm_id` output (or a scale set / Arc server ARM id).
- **A data collection rule** -- reference an AzureMonitorDataCollectionRule's `data_collection_rule_id` output (or provide a Data Collection Endpoint ARM id for the endpoint form).

### Azure Subscription

- **Exactly one binding per association** -- a rule OR an endpoint, never both; a machine can carry MANY associations (several rules, plus at most one endpoint association).
- **The name is required for rule bindings** and must be left unset for endpoint bindings (Azure mandates the fixed name `configurationAccessEndpoint` there; the engines apply it automatically).
- **Collection needs the agent** -- the association creates fine on a machine without the Azure Monitor Agent; telemetry starts flowing when the agent runs and picks the association up.
- **The association is free** -- the telemetry the rule collects is billed at its destinations.

## Deploy

### Console

Open the deployment store, find **Azure Monitor Data Collection Rule Association**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Attach VM to Rule** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f dcr-association.yaml
```

## After Deploy

Query the rule's destination workspace for the machine's records (`Syslog`, `Perf`, `Event`, or your custom table) -- first records typically land 3-5 minutes after the agent discovers the association. The portal shows every association on the machine's **Data collection rules** blade and on the rule's **Resources** blade.
