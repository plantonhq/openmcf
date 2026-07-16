---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual Network Peering"
type: "preset-list"
componentSlug: "virtual-network-peering"
componentTitle: "Virtual Network Peering"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-hub-to-spoke"
    rank: "01"
    title: "Hub to Spoke"
    excerpt: "This preset is the **hub-side half** of a hub-and-spoke peering pair. It is written on the hub network and points at one spoke. Connectivity only works once the reciprocal `spoke-to-hub` preset (or..."
  - slug: "02-spoke-to-hub"
    rank: "02"
    title: "Spoke to Hub"
    excerpt: "This preset is the **spoke-side half** of a hub-and-spoke peering pair. It is written on the spoke network and points back at the hub. Deploy it alongside `01-hub-to-spoke` (or an equivalent hub-side..."
  - slug: "03-subnet-scoped-peering"
    rank: "03"
    title: "Subnet-Scoped Peering"
    excerpt: "This preset peers **named subnets only** instead of the networks' complete address spaces. Set `peer_complete_virtual_networks_enabled` to false and list the local and remote subnet names included in..."
---

# Virtual Network Peering Presets

Ready-to-deploy configuration presets for Virtual Network Peering. Each preset is a complete manifest you can copy, customize, and deploy.
