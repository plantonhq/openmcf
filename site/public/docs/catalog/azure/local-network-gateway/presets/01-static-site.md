---
title: "Static Site"
description: "This preset describes a classic on-premises site: the VPN device's static public IP and the address ranges behind it. Azure routes the declared prefixes into whatever tunnel references this..."
type: "preset"
rank: "01"
presetSlug: "01-static-site"
componentSlug: "local-network-gateway"
componentTitle: "Local Network Gateway"
provider: "azure"
icon: "package"
order: 1
---

# Static Site

This preset describes a classic on-premises site: the VPN device's static public IP and the address ranges behind it. Azure routes the declared prefixes into whatever tunnel references this description -- the prefix list IS the routing statement.

## When to Use

- Datacenters and offices with a static public IP on the VPN device
- A handful of stable sites whose prefixes rarely change (static routing is simpler to reason about than BGP at small scale)
- The site half of the three-resource site-to-site story (site + gateway + connection)

## Key Configuration Choices

- **`gatewayAddress`** -- the device's PUBLIC IPv4; use the FQDN form (see the BGP Site preset) only when the address genuinely changes
- **Honest `addressSpaces`** -- declare exactly what lives behind the device. Over-broad prefixes ("10.0.0.0/8 to be safe") blackhole traffic that belongs elsewhere; overlaps with the VNet need gateway NAT rules, not creative declarations
- **Name it after the site** -- "hq-datacenter", "branch-london": one description per site, referenced by that site's connection
- **Free and instant** -- the object is pure ARM metadata; nothing contacts the device at deploy time (a wrong address surfaces later as a tunnel stuck in Connecting)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (by convention, the connecting gateway's) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-device-public-ip>` | The VPN device's static public IPv4 | Your network team / device WAN interface |
| `<your-onprem-cidr>` | The range(s) reachable behind the device | Your on-premises IPAM |
