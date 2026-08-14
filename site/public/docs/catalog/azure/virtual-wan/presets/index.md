---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual WAN"
type: "preset-list"
componentSlug: "virtual-wan"
componentTitle: "Virtual WAN"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-wan"
    rank: "01"
    title: "Standard WAN"
    excerpt: "This preset creates the full-mesh default: a Standard-tier Virtual WAN with ARM's defaults stated -- VPN traffic encrypted, branch-to-branch transit on, no Office 365 local breakout. The WAN object..."
  - slug: "02-isolated-branches-wan"
    rank: "02"
    title: "Isolated Branches WAN"
    excerpt: "This preset creates the hub-and-spoke-only shape: branches (VPN sites) can reach Azure through their hubs but can NOT reach each other through the WAN, and latency-sensitive Office 365 traffic (the..."
---

# Virtual WAN Presets

Ready-to-deploy configuration presets for Virtual WAN. Each preset is a complete manifest you can copy, customize, and deploy.
