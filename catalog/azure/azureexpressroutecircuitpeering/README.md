# Overview

The **Azure ExpressRoute Circuit Peering API Resource** provides a consistent and standardized interface for deploying and managing ExpressRoute circuit peerings -- the BGP routing configurations that make routes flow through a circuit. The circuit is the physical pipe; a peering of each type (private for your VNets, Microsoft for Microsoft public services) is what turns it into connectivity.

## Purpose

We developed this API resource so a circuit's routing surface is explicit, versioned configuration rather than portal state:

- **Private peering**: routes your VNets' private address space over the circuit -- what an EXPRESS_ROUTE virtual network gateway connects to
- **Microsoft peering**: routes Microsoft 365 and Azure public services over the circuit, with the registered-prefix advertisement contract and route-filter selection modeled in full
- **Global Reach**: connections from this circuit's private peering to other circuits' private peerings, linking the on-premises sites behind them across the Microsoft backbone

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Type-Aware Contracts**: route filters and Microsoft configs only on Microsoft peering, no IPv6 on the deprecated public peering, Global Reach only on private peering -- all enforced at validation time
- **Full Dual-Stack Modeling**: IPv4 /30 pairs and the complete `ipv6` block with its own advertisement contract
- **Sensitive BGP Keys**: the MD5 `sharedKey` and Global Reach `authorizationKey` are secret-typed end to end

## Use Cases

- **VNet connectivity**: private peering consumed by an ExpressRoute-type virtual network gateway -- the standard hybrid path
- **Microsoft 365 over ExpressRoute**: Microsoft peering with registered public prefixes and a route filter selecting service communities
- **Site-to-site over the backbone**: Global Reach connections joining datacenters behind different circuits, no VPN needed
- **Dual-stack estates**: IPv6 session pairs alongside IPv4 on private or Microsoft peering

## Future Enhancements

Future updates will include:

- **Session health surfacing**: BGP session state and advertised/received route counts in the console
- **Route filter management**: a first-class route-filter resource for Microsoft-peering community selection

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
