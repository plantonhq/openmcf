---
title: "Presets"
description: "Ready-to-deploy configuration presets for DNS Resolver Virtual Network Link"
type: "preset-list"
componentSlug: "dns-resolver-virtual-network-link"
componentTitle: "DNS Resolver Virtual Network Link"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-spoke-link"
    rank: "01"
    title: "Spoke Link"
    excerpt: "This preset attaches one spoke network to the hub's forwarding ruleset -- the moment it deploys, resources in the spoke resolve the ruleset's domains through the hub resolver's outbound endpoint. No..."
---

# DNS Resolver Virtual Network Link Presets

Ready-to-deploy configuration presets for DNS Resolver Virtual Network Link. Each preset is a complete manifest you can copy, customize, and deploy.
