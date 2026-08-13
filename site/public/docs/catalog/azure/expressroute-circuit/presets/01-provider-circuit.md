---
title: "Provider Circuit"
description: "This preset creates the common ExpressRoute shape: a circuit bought through a connectivity provider at one of their peering locations. Creation issues the service key; hand it to your provider to..."
type: "preset"
rank: "01"
presetSlug: "01-provider-circuit"
componentSlug: "expressroute-circuit"
componentTitle: "ExpressRoute Circuit"
provider: "azure"
icon: "package"
order: 1
---

# Provider Circuit

This preset creates the common ExpressRoute shape: a circuit bought through a connectivity provider at one of their peering locations. Creation issues the service key; hand it to your provider to complete the physical cross-connect, then configure peerings.

## When to Use

- Datacenter or colocation connectivity to Azure through a carrier (the overwhelmingly common case)
- The anchor of the hybrid chain: circuit → private peering → ExpressRoute gateway → connection

## Key Configuration Choices

- **STANDARD + METERED_DATA is the sane default** -- STANDARD reaches every region in the geopolitical area; metered beats unlimited below roughly two-thirds sustained utilization
- **The trio is fixed at creation** -- provider, peering location, and (downward) bandwidth all force replacement; bandwidth can only GROW in place, so size at the low end of plausible
- **Provider vocabulary is exact** -- `az network express-route list-service-providers` shows the accepted provider names and locations
- **Billing starts at creation** -- the meter runs while the circuit sits `NotProvisioned` waiting on the carrier

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The ARM metadata region (not the physical site) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-connectivity-provider>` | The provider's exact Azure-listed name | `az network express-route list-service-providers` |
| `<your-peering-location>` | The provider's cross-connect site | Your provider's ExpressRoute order |
