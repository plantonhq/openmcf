# Overview

The **Azure Point-to-Site VPN Gateway API Resource** provides a consistent and standardized interface for deploying and managing Point-to-Site VPN Gateways -- the managed receivers inside Virtual WAN hubs that individual devices (laptops, phones) dial into from anywhere. Authentication policy lives on the VPN Server Configuration the gateway points at; the gateway itself owns capacity and the client address pools.

## Purpose

We developed this API resource so connecting people (not branches) to the managed hub network is one first-class, versioned object:

- **Capacity and pools, declared once**: scale units and named client address pools with optional per-pool hub routing
- **Policy by reference**: the VPN Server Configuration is a typed reference -- reusable authentication policy, chart-ready wiring
- **Hub-native**: deploys into a Virtual Hub (one per hub, separate from the site-to-site slot), inheriting the WAN's routing fabric

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Named address pools**: each connection configuration is a client population with its own pool, internet-security switch, and hub routing block
- **Hub routing contracts**: a configured route block requires its associated route table (ARM's rule, validated upfront) and references the hub's route tables and route maps by name
- **Honest lifecycle**: the gateway bills from creation and creates in tens of minutes -- the spec documents every replacing field

## Use Cases

- **Remote workforce access** to everything the hub reaches -- spokes, branches, on-prem via ExpressRoute
- **Forced tunneling**: advertise the default route so client internet traffic rides the tunnel into a hub firewall
- **Segmented user populations**: multiple pools mapped to the server configuration's policy groups (OpenVPN)

## Future Enhancements

Future updates will include:

- **Client connection diagnostics surfacing**: per-pool connection health as Azure exposes it

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
