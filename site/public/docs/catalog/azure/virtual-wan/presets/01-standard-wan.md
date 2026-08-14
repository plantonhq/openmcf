---
title: "Standard WAN"
description: "This preset creates the full-mesh default: a Standard-tier Virtual WAN with ARM's defaults stated -- VPN traffic encrypted, branch-to-branch transit on, no Office 365 local breakout. The WAN object..."
type: "preset"
rank: "01"
presetSlug: "01-standard-wan"
componentSlug: "virtual-wan"
componentTitle: "Virtual WAN"
provider: "azure"
icon: "package"
order: 1
---

# Standard WAN

This preset creates the full-mesh default: a Standard-tier Virtual WAN with ARM's defaults stated -- VPN traffic encrypted, branch-to-branch transit on, no Office 365 local breakout. The WAN object is free; hubs and gateways created under it carry the cost.

## When to Use

- The root of any managed hub-and-spoke deployment (the overwhelmingly common case)
- The first object of the Virtual WAN chain: WAN → hub(s) → connections and gateways

## Key Configuration Choices

- **Standard is the right tier for almost everyone** -- ExpressRoute, site-to-site and point-to-site VPN, hub-to-hub transit; Basic is a constrained legacy tier that can never be reached back down from Standard
- **The defaults are ARM's** -- encryption on, branch-to-branch on, breakout NONE; override only with a reason
- **Deletion is bottom-up** -- ARM refuses to delete a WAN that still has hubs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The ARM metadata region (hubs choose their own regions) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
