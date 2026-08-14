---
title: "Single-Link Branch"
description: "This preset describes the classic branch: one ISP with a static public IP, and a static list of the prefixes reachable behind the branch. The site is free and deploys nothing at the branch -- it is..."
type: "preset"
rank: "01"
presetSlug: "01-single-link-branch"
componentSlug: "vpn-site"
componentTitle: "VPN Site"
provider: "azure"
icon: "package"
order: 1
---

# Single-Link Branch

This preset describes the classic branch: one ISP with a static public IP, and a static list of the prefixes reachable behind the branch. The site is free and deploys nothing at the branch -- it is the address-book entry a VPN Gateway Connection points at.

## When to Use

- Branches with one internet uplink and a static public IP
- Datacenters whose reachable prefixes are known and stable (no BGP)

## Key Configuration Choices

- **Static routing** -- `addressCidrs` carries the branch's prefixes; Azure routes them into the tunnel (no BGP to operate)
- **One link, named for its meaning** -- the link's name ("primary-isp") is the key the connection and the `link_ids` output use
- **Speed is informational** -- `speedInMbps` aids portal display and partner automation; it does not rate-limit

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the site object | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-virtual-wan-arm-id>` | ARM ID of the Virtual WAN | `AzureVirtualWan` status outputs (`virtual_wan_id`), or reference it with valueFrom |

Replace the example endpoint (`203.0.113.10`) and prefixes (`192.168.10.0/24`) with the branch's real values.
