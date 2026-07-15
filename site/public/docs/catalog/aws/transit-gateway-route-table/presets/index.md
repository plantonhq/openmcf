---
title: "Presets"
description: "Ready-to-deploy configuration presets for Transit Gateway Route Table"
type: "preset-list"
componentSlug: "transit-gateway-route-table"
componentTitle: "Transit Gateway Route Table"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-isolated-domain"
    rank: "01"
    title: "Isolated Domain"
    excerpt: "The canonical segmentation building block: a routing domain whose members see shared services but not the rest of the network, with a blackhole guaranteeing the isolation."
  - slug: "02-inspection-domain"
    rank: "02"
    title: "Inspection Domain"
    excerpt: "The hair-pin pattern's spoke side: every flow leaving a spoke is default-routed through the inspection VPC's attachment, where stateful appliances decide what continues."
---

# Transit Gateway Route Table Presets

Ready-to-deploy configuration presets for Transit Gateway Route Table. Each preset is a complete manifest you can copy, customize, and deploy.
