# Overview

The **Azure ExpressRoute Gateway API Resource** provides a consistent and standardized interface for deploying and managing ExpressRoute Gateways -- the on-ramps that bring ExpressRoute circuits into a Virtual WAN hub. This is the Virtual WAN counterpart of the classic VNet-resident ExpressRoute gateway: it lives IN a hub, and each of its connections joins one circuit peering to the hub so branches and spokes across the WAN reach the private line.

## Purpose

We developed this API resource so private-circuit connectivity into the managed WAN is one first-class, versioned object:

- **The circuit's WAN on-ramp**: the gateway plus its circuit-peering connections as one unit
- **Cross-subscription ready**: authorization-key redemption for circuits owned elsewhere, modeled as a validated, sensitive field
- **Chart-ready wiring**: the hub and circuit peering are typed references; connection ARM IDs surface in a name-keyed output map

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Composed Connections**: one ARM child per `connections` entry, each with its own routing block and weight
- **Safe Defaults**: scale from one unit, WAN-only traffic, ARM's routing defaults -- stated rather than implied
- **Validation Rules**: scale bounds (1-10), UUID-shaped authorization keys, unique connection names, and the provider's routing-block contract enforced upfront

## Use Cases

- **Bring a carrier circuit into the WAN**: datacenter traffic reaching every hub-connected spoke and branch
- **Cross-subscription circuits**: redeem an authorization issued by the circuit's owner
- **FastPath**: gateway bypass for on-premises-to-Private-Link traffic on supported circuit SKUs

## Future Enhancements

Future updates will include:

- **Topology views**: the hub with its gateways and circuits as one navigable story

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
