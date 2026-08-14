---
title: "Standard Port"
description: "This preset creates the common ExpressRoute Direct shape: a 10 Gbps Dot1Q port pair with both links administratively enabled and metered billing. The per-link outputs (router, interface, patch panel,..."
type: "preset"
rank: "01"
presetSlug: "01-standard-port"
componentSlug: "expressroute-port"
componentTitle: "ExpressRoute Port"
provider: "azure"
icon: "package"
order: 1
---

# Standard Port

This preset creates the common ExpressRoute Direct shape: a 10 Gbps Dot1Q port pair with both links administratively enabled and metered billing. The per-link outputs (router, interface, patch panel, rack) are the letter-of-authorization facts your colocation facility needs to order the two physical cross-connects.

## When to Use

- Dedicated 10 Gbps+ capacity into Azure for a large hybrid estate
- The anchor of the ExpressRoute Direct chain: port → Direct-mode circuit → private peering → gateway connection

## Key Configuration Choices

- **The physical trio is fixed at creation** -- peering location, bandwidth, and encapsulation all force replacement (of the port AND every circuit on it)
- **DOT1Q is the common encapsulation** -- one VLAN tag per circuit; choose QINQ only when overlapping customer VLAN ranges must share the port
- **Billing starts at creation** -- the port's full monthly rate runs from the moment ARM creates it, cross-connects or not
- **Links start enabled here** -- flip `adminEnabled` to false if you prefer to enable them only after the facility completes the cross-connects

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The ARM metadata region (not the physical site) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-peering-location>` | The ExpressRoute Direct facility | `az network express-route port location list` |
