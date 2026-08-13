# Overview

The **Azure Virtual Hub Connection API Resource** provides a consistent and standardized interface for deploying and managing Virtual Hub Connections -- the attachments that join spoke virtual networks to a Virtual WAN hub. Once connected, a spoke reaches everything the hub's routing lets it reach: other spokes, branches behind the hub's gateways, and other hubs' networks through the WAN mesh.

## Purpose

We developed this API resource so joining a network to the managed WAN -- and the routing decisions that shape the topology -- is one first-class, versioned object:

- **The spoke's on-ramp**: one connection attaches one VNet to one hub
- **Topology as configuration**: route-table association, label-based propagation, static routes toward NVAs, and internet-security -- the whole isolation/shared-services/service-chaining vocabulary, typed and validated
- **Chart-ready wiring**: the hub and VNet are typed references; the connection's ARM ID surfaces where hub BGP peerings consume it

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Routing Block Contracts**: the provider's at-least-one-of rules enforced upfront -- a configured routing block must actually configure something
- **Safe Defaults**: unset routing attaches any-to-any through the hub's default table (ARM's behavior, stated rather than implied)

## Use Cases

- **Attach application spokes** to a regional hub for managed transit
- **Spoke isolation**: associate with an isolated table, propagate only to shared services
- **Service chaining**: static routes steering prefixes at an NVA inside the spoke
- **Centralized egress**: internet security on, paired with a hub firewall via routing intent

## Future Enhancements

Future updates will include:

- **Topology views**: the hub with its connections as one navigable story

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
