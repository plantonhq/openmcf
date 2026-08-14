# Azure Event Grid Namespace

Deploys an Azure Event Grid namespace -- the capacity-scaled hub of the newer Event Grid that hosts CloudEvents namespace topics and an optional MQTT broker behind one set of regional endpoints. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Event Grid namespace** -- the hub itself: throughput capacity, network posture (public access plus inbound IP rules), optional managed identity, and the optional MQTT broker ("topic spaces") configuration

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name.
- **For MQTT routing (optional)** -- an AzureEventgridTopic in the same region to route MQTT messages into; reference its ID output.

### Azure Subscription

- **The MQTT block is create-only** -- the topic-spaces configuration cannot be added, removed, or changed after create; the namespace is replaced instead. Decide the MQTT posture up front.
- **Capacity is billed per throughput unit** -- a namespace is NOT free at rest (unlike classic topics); one TU is the floor.
- **The name behaves like a DNS label** -- it appears in the namespace's regional hostnames; prefix it with your org.
- **The SKU is always "Standard"** -- Azure defines no other value today, so the platform sends it for you.

## Deploy

### Console

Open the deployment store, find **Azure Event Grid Namespace**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **CloudEvents Hub** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f eventgrid-namespace.yaml
```

## After Deploy

The `namespace_id` output is the wiring edge for streams -- create AzureEventgridNamespaceTopic resources with `namespace_id` referencing it. If the MQTT broker is enabled, clients connect to the namespace's MQTT hostname (the **Overview** blade shows it) with certificate authentication; watch throughput against the provisioned capacity on the **Metrics** blade and raise `capacity` in place when publishers grow.
