# Overview

The **AzureBastionHost** component deploys Azure Bastion -- the managed jump service that opens RDP/SSH sessions to virtual machines over their PRIVATE addresses, from the browser (and, on higher SKUs, from native clients). The machines themselves never expose a public IP, an open port 3389/22, or an inbound NSG rule: Bastion is the one hardened front door.

## Purpose

- **Kill public IPs on VMs**: admin access flows through one audited, Azure-managed entry point instead of per-machine exposure.
- **The SKU picks the shape**: Developer (free, shared infrastructure, dev/test), Basic (dedicated, fixed capacity), Standard (scaling + the feature set), Premium (session recording + private-only deployment).
- **Session-level control**: clipboard, file copy, shareable links, native-client tunneling, Kerberos, and recording are explicit, SKU-gated switches.

## Key Features

- Full azurerm v5 surface: all four SKUs, the complete feature matrix, scale units (2-50), availability zones, and the Developer SKU's virtual-network attachment.
- The provider's whole SKU/feature contract is validated in seconds -- an invalid combination (e.g. tunneling on Basic) never reaches Azure.
- Chart-ready: references the subnet, public IP, and virtual network by typed outputs; publishes the host's ARM ID, DNS name, and private-only state.

## Use Cases

- **Production VM access** without any public exposure -- Standard with tunneling for engineers who live in a terminal.
- **Compliance-driven access** -- Premium with session recording and (optionally) private-only deployment.
- **Dev/test access on a budget** -- the free Developer SKU attached straight to a virtual network.

## Future Enhancements

- The `AzureBastionSubnet` name contract stays documentation until references can be introspected offline.
