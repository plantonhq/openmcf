---
title: "Presets"
description: "Ready-to-deploy configuration presets for Point-to-Site VPN Gateway"
type: "preset-list"
componentSlug: "point-to-site-vpn-gateway"
componentTitle: "Point-to-Site VPN Gateway"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-remote-users"
    rank: "01"
    title: "Standard Remote Users"
    excerpt: "This preset deploys the classic remote-access shape: one client address pool with split tunneling -- connected devices reach everything the hub reaches over the tunnel, and the internet locally...."
  - slug: "02-forced-tunnel-clients"
    rank: "02"
    title: "Forced-Tunnel Clients"
    excerpt: "This preset routes EVERYTHING from connected clients into the hub -- `internetSecurityEnabled: true` advertises the default route (0.0.0.0/0) into the tunnel, so internet-bound traffic can be..."
---

# Point-to-Site VPN Gateway Presets

Ready-to-deploy configuration presets for Point-to-Site VPN Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
