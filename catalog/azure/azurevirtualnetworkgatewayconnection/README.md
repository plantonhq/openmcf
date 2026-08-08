# Overview

The **Azure Virtual Network Gateway Connection API Resource** provides a consistent and standardized interface for deploying and managing gateway connections -- the tunnel objects that join a virtual network gateway to its far side: an on-premises VPN device, another gateway, or an ExpressRoute circuit. This is the resource that turns a provisioned gateway into actual connectivity.

## Purpose

We developed this API resource so each tunnel is its own declarative, composable object. The gateway (expensive, long-lived) is deployed once; connections (cheap, per-site) come and go with the network topology:

- **Site-to-Site (IPSEC)**: the classic datacenter-to-Azure tunnel, pointing at an AzureLocalNetworkGateway
- **VNet-to-VNet**: encrypted tunnels between two gateways across regions
- **ExpressRoute**: attach a circuit's private peering to the gateway
- **Custom IPsec/IKE proposals**: pin exact algorithms for on-premises devices that demand them
- **NAT opt-in**: apply the gateway's NAT rules per tunnel via egress/ingress rule id lists

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Composition by Reference**: the connection references its gateway, its site description, and (for VNet-to-VNet) its peer gateway -- the whole tunnel topology stays visible in the resource graph
- **Type-Aware Validation**: each connection type's required far side (site, peer gateway, or circuit) is enforced at validation time, along with the FastPath and policy-based-selector contracts ARM would otherwise reject at apply
- **Secret-Safe Keys**: the IPsec pre-shared key and ExpressRoute authorization key are sensitive reference fields -- resolved at deploy time, never stored in manifests
- **Honest Semantics**: provisioned is not connected -- the docs and outputs distinguish ARM provisioning success from live tunnel establishment

## Use Cases

- **Branch connectivity**: one connection per office to a shared hub gateway
- **Hybrid DR**: VNet-to-VNet tunnels pairing regions for replication traffic
- **Compliance-pinned tunnels**: custom IPsec proposals (DH groups, AES-256, PFS) for regulated environments
- **Overlapping networks**: tunnels that opt into gateway NAT rules to translate conflicting address space
- **ExpressRoute attach**: joining a provisioned circuit to a VNet's ExpressRoute gateway, including cross-subscription circuits via authorization keys

## Future Enhancements

Future updates will include:

- **ExpressRoute circuit kinds**: first-class circuit and peering resources so `expressRouteCircuitId` wires by default reference
- **Connection health surfacing**: live tunnel state (Connected/Connecting) in resource status
- **Shared-key rotation**: coordinated re-key workflows across both tunnel ends

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
