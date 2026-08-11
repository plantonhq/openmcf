# Overview

The **Azure Virtual WAN API Resource** provides a consistent and standardized interface for deploying and managing Virtual WANs -- the top-level umbrella of Azure's managed hub-and-spoke networking. The WAN object itself is a free, lightweight container of policy; the regional routers (virtual hubs) and the gateways attached to them are separate resources that reference it.

## Purpose

We developed this API resource so the foundation of a managed global network is a first-class, versioned object:

- **The anchor of the Virtual WAN stack**: hubs, hub connections, and every vWAN-attached gateway reference this object's ID
- **Transit policy, typed**: branch-to-branch reachability, VPN encryption, and Office 365 local breakout as explicit, documented choices
- **Standard vs Basic honesty**: the tier choice (and its upgrade-only ratchet) documented where the decision is made

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Chart-Ready Outputs**: the WAN's ARM ID surfaces exactly where hub components consume it
- **Safe Defaults**: Standard tier, encryption on, branch-to-branch transit on -- ARM's defaults, stated rather than implied

## Use Cases

- **Global hub-and-spoke**: the umbrella under which regional hubs interconnect branches, VNets, and ExpressRoute
- **Managed branch connectivity**: SD-WAN/VPN device fleets terminating on hub gateways under one WAN
- **Office 365 breakout policy**: centrally declare which O365 traffic exits at local branches

## Future Enhancements

Future updates will include:

- **Topology views**: the WAN with its hubs and connections as one navigable story

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
