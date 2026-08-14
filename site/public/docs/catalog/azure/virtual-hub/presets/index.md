---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual Hub"
type: "preset-list"
componentSlug: "virtual-hub"
componentTitle: "Virtual Hub"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-hub"
    rank: "01"
    title: "Standard Hub"
    excerpt: "This preset creates the common case: a Standard-tier WAN hub with ARM's defaults stated -- ExpressRoute routing preference, branch-to-branch off, the router at its capacity floor of 2 units, and no..."
  - slug: "02-isolated-spokes-hub"
    rank: "02"
    title: "Isolated Spokes Hub"
    excerpt: "This preset creates a hub prepared for spoke isolation: two custom route tables (\"isolated\" and \"shared-services\", each labeled) that connections associate with and propagate to. The tables..."
  - slug: "03-secured-hub-routing-intent"
    rank: "03"
    title: "Secured Hub (Routing Intent)"
    excerpt: "This preset creates a secured hub: routing intent steers both Internet-bound and private traffic through an Azure Firewall deployed in this hub. Every spoke and branch transiting the hub is inspected..."
---

# Virtual Hub Presets

Ready-to-deploy configuration presets for Virtual Hub. Each preset is a complete manifest you can copy, customize, and deploy.
