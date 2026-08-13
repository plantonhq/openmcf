---
title: "Microsoft Peering"
description: "This preset configures Microsoft peering -- Microsoft 365 and Azure public services delivered over your circuit instead of the internet. It carries the mandatory advertisement contract (your..."
type: "preset"
rank: "02"
presetSlug: "02-microsoft-peering"
componentSlug: "expressroute-circuit-peering"
componentTitle: "ExpressRoute Circuit Peering"
provider: "azure"
icon: "package"
order: 2
---

# Microsoft Peering

This preset configures Microsoft peering -- Microsoft 365 and Azure public services delivered over your circuit instead of the internet. It carries the mandatory advertisement contract (your registered public prefixes) and the route filter that selects which service communities you receive.

## When to Use

- Microsoft 365 traffic over ExpressRoute (with Microsoft's sign-off -- they gate M365-over-ER scenarios)
- Reaching Azure public service endpoints over the private path

## Key Configuration Choices

- **Public addressing throughout** -- the session /30s use PUBLIC IPs here (unlike private peering), and `advertisedPublicPrefixes` must be prefixes REGISTERED to you; Microsoft validates ownership out-of-band before activating
- **The route filter is not optional in practice** -- without `routeFilterId`, Microsoft peering advertises NOTHING to you
- **`customerAsn` and `routingRegistryName`** cover the on-behalf-of case (prefixes registered to a downstream customer) -- leave defaults otherwise
- **IPv6 rides the `ipv6` block** -- its own /126 pair and its own advertisement contract

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | The circuit's resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-circuit-name>` | The parent circuit's name | `AzureExpressRouteCircuit` status outputs (`express_route_circuit_name`) |
| `<your-primary-slash-30>` / `<your-secondary-slash-30>` | PUBLIC /30s for the BGP sessions | Your provider's handoff + your public allocation |
| `<your-registered-public-prefix>` | A public prefix registered to your ASN | Your RIR records (ARIN/RIPE/APNIC/...) |
| `<your-route-filter-arm-id>` | The route filter selecting service communities | Azure portal → Route filters |
