---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual Network Gateway"
type: "preset-list"
componentSlug: "virtual-network-gateway"
componentTitle: "Virtual Network Gateway"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-site-to-site-vpn"
    rank: "01"
    title: "Site-to-Site VPN Gateway"
    excerpt: "This preset creates a route-based VpnGw1 VPN gateway with BGP -- the standard Azure-side anchor for datacenter-to-Azure connectivity. Pair it with an `AzureLocalNetworkGateway` per site and an..."
  - slug: "02-active-active-zone-redundant"
    rank: "02"
    title: "Active-Active Zone-Redundant Gateway"
    excerpt: "This preset creates a Generation2 VpnGw2AZ gateway running as an ACTIVE-ACTIVE pair: two gateway instances, each with its own public IP and APIPA BGP endpoint, both terminating tunnels..."
  - slug: "03-point-to-site-entra-id"
    rank: "03"
    title: "Point-to-Site with Entra ID"
    excerpt: "This preset creates a VPN gateway whose point-to-site clients authenticate with Entra ID (Azure AD) over OpenVPN -- remote-workforce VPN access with your existing identity provider, no certificate..."
---

# Virtual Network Gateway Presets

Ready-to-deploy configuration presets for Virtual Network Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
