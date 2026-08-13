---
title: "Direct Port Circuit"
description: "This preset carves a circuit from your own ExpressRoute Direct port pair -- no third-party provider in the path. It also shows authorization issuance: the named entry generates a key (in the..."
type: "preset"
rank: "03"
presetSlug: "03-direct-port-circuit"
componentSlug: "expressroute-circuit"
componentTitle: "ExpressRoute Circuit"
provider: "azure"
icon: "package"
order: 3
---

# Direct Port Circuit

This preset carves a circuit from your own ExpressRoute Direct port pair -- no third-party provider in the path. It also shows authorization issuance: the named entry generates a key (in the sensitive `authorization_keys` output) that a virtual network gateway in another subscription redeems to connect.

## When to Use

- Estates with their own ExpressRoute Direct ports (10/100 Gbps into Microsoft's edge)
- Multiple circuits carved from one port for different environments or business units
- Hub subscriptions issuing circuit access to spoke subscriptions

## Key Configuration Choices

- **The pair travels together** -- `expressRoutePortId` + `bandwidthInGbps`, mutually exclusive with the provider trio (spec-enforced)
- **Rate limiting is Direct-only and recommended** -- without it, a circuit can burst beyond its configured bandwidth and starve the port's other circuits
- **Name authorizations for their consumers** -- "spoke-subscription", "partner-x": the name keys the generated key in the output, and deleting the entry revokes access
- **No provider handoff** -- Direct circuits skip the NotProvisioned wait; peerings can follow immediately

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The ARM metadata region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-express-route-port-arm-id>` | The Direct port's ARM ID | Azure portal → ExpressRoute Direct → the port → Properties |
