---
title: "ExpressRoute Gateway"
description: "ExpressRoute Gateway deployment documentation"
icon: "package"
order: 100
componentName: "azureexpressroutegateway"
---

# Azure ExpressRoute Gateway

Deploys an ExpressRoute Gateway -- the Virtual WAN on-ramp for ExpressRoute circuits -- together with its connections, each joining one circuit peering to the hub. The gateway bills per scale unit from creation and takes roughly 30 minutes to provision; a hub holds at most one. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ExpressRoute Gateway** -- the gateway in the hub, with its scale-unit floor and non-WAN-traffic policy
- **ExpressRoute Connections** (optional) -- one per `connections` entry: the join between a circuit's private peering and the hub, with authorization key, routing block, and weight
- **Azure Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with your `tags`, applied to the gateway

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An Azure Resource Group** where the gateway will be created.
- **An Azure Virtual Hub** the gateway deploys into (must be in the same region).
- **For connections**: an ExpressRoute circuit whose provider side is PROVISIONED -- ARM rejects a connection to an unprovisioned circuit.

## Deploy

### Console

Open the deployment store, find **Azure ExpressRoute Gateway**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Gateway** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureExpressRouteGateway
metadata:
  name: hub-er-gateway
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    value: "network-rg"
  name: hub-er-gateway
  virtualHubId:
    valueFrom:
      kind: AzureVirtualHub
      name: hub-eastus
      fieldPath: status.outputs.virtual_hub_id
  scaleUnits: 1
```

```shell
planton apply -f azure-express-route-gateway.yaml
```

The gateway bills ~$0.42/hr per scale unit from creation, and ARM takes roughly 30 minutes to provision one.

### InfraChart

In a hybrid-connectivity chart, the gateway follows the hub: WAN → hub → **ExpressRoute gateway** → connections referencing circuit peerings.

## Key Configuration

These are the most important decisions when configuring a gateway. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Scale units** -- the MINIMUM (1-10); each unit carries ~2 Gbps of circuit-to-WAN throughput and bills hourly. ARM auto-scales above the floor under load. Updatable in place.

**Connections** -- each joins one circuit PRIVATE PEERING to the hub. A cross-subscription circuit needs the `authorizationKey` its owner generated; a same-subscription circuit needs none.

**Non-Virtual-WAN traffic** -- off by default; on lets classic VNets connected to the same circuit exchange traffic through this gateway.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |
| **AzureVirtualHub** | `virtualHubId` | `status.outputs.virtual_hub_id` |
| **AzureExpressRouteCircuitPeering** (per connection) | `connections[].expressRouteCircuitPeeringId` | `status.outputs.express_route_circuit_peering_id` |
| **AzureVirtualHub** (optional) | `connections[].routing.associatedRouteTableId` | `status.outputs.default_route_table_id` or `status.outputs.route_table_ids.<name>` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `express_route_gateway_id` | ARM ID of the gateway | Operational tooling |
| `express_route_gateway_name` | Name of the gateway | Operational tooling |
| `connection_ids` | Connection ARM IDs, keyed by name | Operational tooling |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard gateway** -- one scale unit in the hub, connections added when circuits are ready. Start from the **Standard Gateway** preset.

**Circuit connection** -- the gateway joined to a provisioned circuit's private peering. Start from the **Circuit Connection** preset.

## Works With

- [**Azure Resource Group**](/cloud-catalog/azure-resource-group) -- provides the resource group the gateway is created in
- [**Azure Virtual Hub**](/cloud-catalog/azure-virtual-hub) -- the hub the gateway deploys into
- [**Azure ExpressRoute Circuit Peering**](/cloud-catalog/azure-express-route-circuit-peering) -- the private peering each connection joins
