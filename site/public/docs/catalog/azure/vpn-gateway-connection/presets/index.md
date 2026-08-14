---
title: "Presets"
description: "Ready-to-deploy configuration presets for VPN Gateway Connection"
type: "preset-list"
componentSlug: "vpn-gateway-connection"
componentTitle: "VPN Gateway Connection"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-branch-connection"
    rank: "01"
    title: "Standard Branch Connection"
    excerpt: "This preset creates the simplest tunnel: one link to the site's primary ISP, Azure's default IPsec proposals, an Azure-generated pre-shared key, and ARM's default hub routing (associate with and..."
  - slug: "02-pinned-ipsec-connection"
    rank: "02"
    title: "Pinned IPsec Connection"
    excerpt: "This preset pins the tunnel to an exact proposal (AES-256/SHA-256, DH group 14, PFS 2048 -- a widely supported compliance baseline) and carries a pre-agreed key. With a pinned proposal there is NO..."
---

# VPN Gateway Connection Presets

Ready-to-deploy configuration presets for VPN Gateway Connection. Each preset is a complete manifest you can copy, customize, and deploy.
