---
title: "Standard Branch Connection"
description: "This preset creates the simplest tunnel: one link to the site's primary ISP, Azure's default IPsec proposals, an Azure-generated pre-shared key, and ARM's default hub routing (associate with and..."
type: "preset"
rank: "01"
presetSlug: "01-standard-branch-connection"
componentSlug: "vpn-gateway-connection"
componentTitle: "VPN Gateway Connection"
provider: "azure"
icon: "package"
order: 1
---

# Standard Branch Connection

This preset creates the simplest tunnel: one link to the site's primary ISP, Azure's default IPsec proposals, an Azure-generated pre-shared key, and ARM's default hub routing (associate with and propagate to the default table). The connection is free; configure the branch device with the gateway's public IPs and the generated key to bring the tunnel up.

## When to Use

- The first connection of a newly described branch
- Branch devices comfortable with Azure's default IPsec proposal set

## Key Configuration Choices

- **No `sharedKey`** -- Azure generates a strong key; retrieve it for the branch device rather than inventing one in a manifest
- **No routing block** -- unset means ARM's any-to-any default through the hub's default route table
- **One link** -- name it after the site link it connects; add a second `vpnLinks` entry (not an edit) when the branch gains a second ISP

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-vpn-gateway-arm-id>` | ARM ID of the hub's VPN gateway | `AzureVpnGateway` status outputs (`vpn_gateway_id`), or reference it with valueFrom |
| `<your-vpn-site-arm-id>` | ARM ID of the branch's site | `AzureVpnSite` status outputs (`vpn_site_id`), or reference it with valueFrom |
| `<your-vpn-site-link-arm-id>` | ARM ID of the site link being connected | `AzureVpnSite` status outputs (`link_ids.<link-name>`), or reference it with valueFrom |
