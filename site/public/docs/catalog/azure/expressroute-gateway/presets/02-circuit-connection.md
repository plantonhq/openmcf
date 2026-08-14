---
title: "Circuit Connection"
description: "This preset creates the gateway joined to a circuit: one connection referencing the circuit's PRIVATE peering, so datacenter routes flow into the hub and every WAN-connected spoke and branch. ARM..."
type: "preset"
rank: "02"
presetSlug: "02-circuit-connection"
componentSlug: "expressroute-gateway"
componentTitle: "ExpressRoute Gateway"
provider: "azure"
icon: "package"
order: 2
---

# Circuit Connection

This preset creates the gateway joined to a circuit: one connection referencing the circuit's PRIVATE peering, so datacenter routes flow into the hub and every WAN-connected spoke and branch. ARM accepts the connection only when the circuit's provider side is provisioned.

## When to Use

- The carrier has provisioned the circuit and its private peering is configured
- Bringing an existing ExpressRoute line into a new Virtual WAN

## Key Configuration Choices

- **The peering, not the circuit** -- the connection references the circuit's private-peering ARM ID (`status.outputs.express_route_circuit_peering_id` on the peering kind)
- **Cross-subscription circuits need a key** -- add `authorizationKey` (the UUID the circuit owner's authorization generated) when the circuit lives in another subscription; leave it out otherwise
- **Routing defaults are ARM's** -- the connection associates with and propagates to the hub's default table; add a `routing` block for isolation or route-map policies

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The gateway's region (must match the hub's) | The hub's configuration |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-virtual-hub-arm-id>` | ARM ID of the hub | `AzureVirtualHub` status outputs (`virtual_hub_id`), or reference it with valueFrom |
| `<your-circuit-private-peering-arm-id>` | ARM ID of the circuit's private peering | `AzureExpressRouteCircuitPeering` status outputs (`express_route_circuit_peering_id`), or reference it with valueFrom |
