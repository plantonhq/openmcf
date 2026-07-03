---
title: "Standard Virtual Network"
description: "This preset creates a general-purpose virtual network with a single /16 address space -- room for dozens of /24 subnets -- and Azure's default DNS resolver, which is what private DNS zone resolution..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "virtual-network"
componentTitle: "Virtual Network"
provider: "azure"
icon: "package"
order: 1
---

# Standard Virtual Network

This preset creates a general-purpose virtual network with a single /16
address space -- room for dozens of /24 subnets -- and Azure's default DNS
resolver, which is what private DNS zone resolution requires. It is the
right starting point for almost every environment.

The network is deliberately just the network. Partition it with
`AzureSubnet` resources (one per workload tier or delegated service),
attach outbound NAT with `AzureNatGateway`, and make private DNS zones
resolvable with `AzurePrivateDnsZoneVirtualNetworkLink` -- each referencing
this network's `virtual_network_id` output.

## When to Use

- The networking foundation for any new environment
- A spoke network in a hub-and-spoke topology
- Anywhere you plan addresses yourself (for centrally-managed allocation,
  see the IPAM alternative in the spec's `ipAddressPools`)

## Key Configuration Choices

- **/16 address space** -- large enough that you will not renumber; add a
  second CIDR block later (in place) if it ever fills
- **Default DNS** -- leave `dnsServers` empty unless integrating
  on-premises DNS; custom servers break direct private-zone resolution
- **Governance tags** -- carried on the network and enforced by Azure
  Policy; user tags win over the Planton-derived ones on collision

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the network in (or use `valueFrom` against an `AzureResourceGroup`) | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
