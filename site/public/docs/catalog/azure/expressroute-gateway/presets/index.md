---
title: "Presets"
description: "Ready-to-deploy configuration presets for ExpressRoute Gateway"
type: "preset-list"
componentSlug: "expressroute-gateway"
componentTitle: "ExpressRoute Gateway"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-gateway"
    rank: "01"
    title: "Standard Gateway"
    excerpt: "This preset creates the gateway alone: one scale unit (~2 Gbps, the guaranteed floor ARM auto-scales above) in the hub, with no connections yet. Connections join circuit peerings later -- typically..."
  - slug: "02-circuit-connection"
    rank: "02"
    title: "Circuit Connection"
    excerpt: "This preset creates the gateway joined to a circuit: one connection referencing the circuit's PRIVATE peering, so datacenter routes flow into the hub and every WAN-connected spoke and branch. ARM..."
---

# ExpressRoute Gateway Presets

Ready-to-deploy configuration presets for ExpressRoute Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
