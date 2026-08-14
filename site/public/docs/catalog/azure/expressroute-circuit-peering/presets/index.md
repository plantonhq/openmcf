---
title: "Presets"
description: "Ready-to-deploy configuration presets for ExpressRoute Circuit Peering"
type: "preset-list"
componentSlug: "expressroute-circuit-peering"
componentTitle: "ExpressRoute Circuit Peering"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-private-peering"
    rank: "01"
    title: "Private Peering"
    excerpt: "This preset configures private peering -- the routing configuration virtually every ExpressRoute deployment needs. Your VNets' private address space flows over the circuit, and an EXPRESS_ROUTE-type..."
  - slug: "02-microsoft-peering"
    rank: "02"
    title: "Microsoft Peering"
    excerpt: "This preset configures Microsoft peering -- Microsoft 365 and Azure public services delivered over your circuit instead of the internet. It carries the mandatory advertisement contract (your..."
---

# ExpressRoute Circuit Peering Presets

Ready-to-deploy configuration presets for ExpressRoute Circuit Peering. Each preset is a complete manifest you can copy, customize, and deploy.
