---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual Network Gateway Connection"
type: "preset-list"
componentSlug: "virtual-network-gateway-connection"
componentTitle: "Virtual Network Gateway Connection"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-site-to-site"
    rank: "01"
    title: "Site-to-Site Tunnel"
    excerpt: "This preset creates the classic site-to-site IPsec connection: the tunnel joining a virtual network gateway to one described on-premises site, with the pre-shared key referenced from a secret...."
  - slug: "02-custom-ipsec-policy"
    rank: "02"
    title: "Custom IPsec Policy Tunnel"
    excerpt: "This preset creates a site-to-site connection with a PINNED IPsec/IKE proposal (DH14 / AES-256 / SHA-256 / PFS2048) -- for on-premises devices whose documentation demands exact algorithms, and for..."
---

# Virtual Network Gateway Connection Presets

Ready-to-deploy configuration presets for Virtual Network Gateway Connection. Each preset is a complete manifest you can copy, customize, and deploy.
