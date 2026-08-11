# Overview

The **Azure Virtual Hub API Resource** provides a consistent and standardized interface for deploying and managing Virtual Hubs -- the managed regional routers at the center of Azure's Virtual WAN. Every network that joins the WAN in a region attaches to that region's hub: spoke VNets through hub connections, branches through the hub's gateways, and other hubs through the WAN's automatic mesh.

## Purpose

We developed this API resource so the heart of a managed global network -- and all of its routing customization -- is one first-class, versioned object:

- **The regional anchor of the Virtual WAN stack**: connections, gateways, and firewalls all reference this object's ID
- **Routing as configuration**: custom route tables, route maps, NVA BGP peerings, and routing intent modeled as typed, validated children
- **Chart-ready wiring**: name-keyed output maps (`route_table_ids`, `route_map_ids`) surface every child's ARM ID exactly where connections consume them

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Composed Routing Children**: route tables (with inline routes), route maps, BGP connections, and routing intent deploy with the hub as one unit
- **Safe Defaults**: Standard tier, ExpressRoute routing preference, ARM's router capacity floor -- ARM's defaults, stated rather than implied
- **Validation Rules**: unique child names, explicit enum choices, and the provider's own create-time contracts (non-Drop route-map actions need parameters) enforced upfront

## Use Cases

- **Regional hub of a global WAN**: the router branches, spokes, and other hubs transit through
- **Spoke isolation and shared services**: custom route tables with label-based propagation
- **Secured hub**: routing intent steering Internet/private traffic through an Azure Firewall in the hub
- **NVA integration**: BGP peering the hub router with appliances running in connected spokes

## Future Enhancements

Future updates will include:

- **Topology views**: the hub with its connections, gateways, and route tables as one navigable story

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
