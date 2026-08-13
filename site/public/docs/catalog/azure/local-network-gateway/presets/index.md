---
title: "Presets"
description: "Ready-to-deploy configuration presets for Local Network Gateway"
type: "preset-list"
componentSlug: "local-network-gateway"
componentTitle: "Local Network Gateway"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-static-site"
    rank: "01"
    title: "Static Site"
    excerpt: "This preset describes a classic on-premises site: the VPN device's static public IP and the address ranges behind it. Azure routes the declared prefixes into whatever tunnel references this..."
  - slug: "02-bgp-site"
    rank: "02"
    title: "BGP Site"
    excerpt: "This preset describes a site whose routes arrive dynamically over BGP instead of a static prefix list -- the posture for estates where sites multiply or prefixes churn. The description carries the..."
---

# Local Network Gateway Presets

Ready-to-deploy configuration presets for Local Network Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
