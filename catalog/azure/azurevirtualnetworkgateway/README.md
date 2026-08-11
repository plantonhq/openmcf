# Overview

The **Azure Virtual Network Gateway API Resource** provides a consistent and standardized interface for deploying and managing virtual network gateways -- the managed appliances that terminate hybrid connectivity (site-to-site VPN, point-to-site VPN, VNet-to-VNet, and ExpressRoute) in an Azure virtual network. This resource closes the private-datacenter-connectivity gap: it is the Azure-side anchor of every "connect my network to Azure" story.

## Purpose

We developed this API resource to make hybrid connectivity a first-class, composable building block. A gateway alone is only one third of the story -- it pairs with AzureLocalNetworkGateway (the on-premises side's description) and AzureVirtualNetworkGatewayConnection (the tunnel) -- and this resource keeps each piece independently deployable:

- **Site-to-Site VPN**: terminate IPsec tunnels from datacenters and branch offices
- **Point-to-Site VPN**: give individual clients VPN access with Entra ID, certificate, or RADIUS authentication
- **VNet-to-VNet**: encrypted tunnels between virtual networks across regions
- **ExpressRoute**: terminate a circuit's private peering into the VNet
- **NAT Rules**: translate overlapping address space between on-premises sites and the VNet, composed directly on the gateway

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Composition by Reference**: the gateway lives in a referenced AzureSubnet (the ARM-mandated "GatewaySubnet") and binds referenced AzurePublicIp resources -- addresses stay visible in the resource graph and reusable
- **Full SKU Vocabulary**: BASIC and VPN_GW_1_AZ..5_AZ for VPN (Azure retired new non-AZ VpnGw creates; the AZ tiers deploy in every region), ER_GW SKUs (including autoscaling ER_GW_SCALE) for ExpressRoute, with the type/generation/SKU pairing rules enforced at validation time -- not at minute 30 of a 45-minute deployment
- **Active-Active Support**: two-instance gateways with per-instance public IPs and APIPA BGP peering addresses
- **Composed NAT Rules**: gateway NAT rules are part of the spec, and each rule's ARM id surfaces in the `nat_rule_ids` output for connections to opt into

## Use Cases

- **Datacenter-to-Azure connectivity**: the classic site-to-site IPsec tunnel to an on-premises VPN device
- **Multi-site hub**: one gateway terminating tunnels from many branch offices, each described by its own local network gateway
- **Remote workforce**: point-to-site VPN with Entra ID authentication and per-group address pools
- **Forced tunneling**: route all VNet egress through on-premises inspection via a default-site local network gateway
- **Overlapping address spaces**: NAT rules translating between conflicting RFC 1918 ranges

## Future Enhancements

Future updates will include:

- **ExpressRoute circuit kinds**: first-class circuit and peering resources completing the ExpressRoute story
- **Gateway diagnostics**: built-in tunnel health and BGP route table visibility
- **Virtual WAN integration**: the hub-based alternative for large-scale branch connectivity

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
