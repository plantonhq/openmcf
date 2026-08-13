# Azure Event Grid Topic

Deploys an Azure Event Grid custom topic -- the HTTPS endpoint an application publishes its own events to, from which event subscriptions fan events out to handlers. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Event Grid topic** -- the publish endpoint with its pair of access keys, optional managed identity, optional IP firewall, and the input schema every incoming event must match

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name.

### Azure Subscription

- **The topic name is a public DNS hostname** (`{name}.{region}.eventgrid.azure.net`) -- unique across ALL Azure customers in the region; a taken name fails the deploy with a conflict.
- **The input schema is chosen at creation** (Event Grid, CloudEvents 1.0, or custom JSON with envelope mappings) and cannot change afterward -- changing it replaces the topic and its endpoint hostname.
- **Publishing authenticates with an access key or Microsoft Entra ID** -- set `local_auth_enabled: false` to force Entra-only publishing.
- **A topic is free at rest** -- billing is per operation (publishes, deliveries, filtering evaluations).

## Deploy

### Console

Open the deployment store, find **Azure Event Grid Topic**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **CloudEvents Topic** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f eventgrid-topic.yaml
```

## After Deploy

The `endpoint` output is the publish URL and the `primary_access_key` output (a secret) authenticates POSTs via the `aeg-sas-key` header. Send a test event and watch it on the topic's **Metrics** blade (Published Events); until a subscription exists, published events are dropped after evaluation -- create a subscription and events start flowing to its handler.
