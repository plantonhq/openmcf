---
title: "Presets"
description: "Ready-to-deploy configuration presets for ExpressRoute Circuit"
type: "preset-list"
componentSlug: "expressroute-circuit"
componentTitle: "ExpressRoute Circuit"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-provider-circuit"
    rank: "01"
    title: "Provider Circuit"
    excerpt: "This preset creates the common ExpressRoute shape: a circuit bought through a connectivity provider at one of their peering locations. Creation issues the service key; hand it to your provider to..."
  - slug: "02-local-metro-circuit"
    rank: "02"
    title: "Local Metro Circuit"
    excerpt: "This preset creates a LOCAL-tier circuit: connectivity to the Azure regions in the circuit's own metro only, with NO egress fees. When your facility sits near the Azure region you use, this is the..."
  - slug: "03-direct-port-circuit"
    rank: "03"
    title: "Direct Port Circuit"
    excerpt: "This preset carves a circuit from your own ExpressRoute Direct port pair -- no third-party provider in the path. It also shows authorization issuance: the named entry generates a key (in the..."
---

# ExpressRoute Circuit Presets

Ready-to-deploy configuration presets for ExpressRoute Circuit. Each preset is a complete manifest you can copy, customize, and deploy.
