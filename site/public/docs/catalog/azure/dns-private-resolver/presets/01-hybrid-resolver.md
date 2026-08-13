---
title: "Hybrid Resolver"
description: "This preset deploys the full hybrid shape: one inbound endpoint (on-premises resolves Azure names through it) and one outbound endpoint (Azure resolves on-premises names through it, steered by a..."
type: "preset"
rank: "01"
presetSlug: "01-hybrid-resolver"
componentSlug: "dns-private-resolver"
componentTitle: "DNS Private Resolver"
provider: "azure"
icon: "package"
order: 1
---

# Hybrid Resolver

This preset deploys the full hybrid shape: one inbound endpoint (on-premises resolves Azure names through it) and one outbound endpoint (Azure resolves on-premises names through it, steered by a forwarding ruleset).

## When to Use

- Hub networks in a hub-and-spoke topology where both directions of hybrid name resolution are needed
- Replacing a pair of IaaS DNS forwarder VMs with one managed service

## Key Configuration Choices

- **One resolver per virtual network** (Azure enforces it) -- anchor it to the hub and let spokes consume it through ruleset links
- **Each endpoint occupies its own delegated subnet** (`Microsoft.Network/dnsResolvers`, /28-/24, nothing else in it) -- carve both subnets when you design the network
- **The inbound IP is dynamically assigned here** -- switch to STATIC allocation if datacenter forwarder configs are expensive to change

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the resolver is created in | The resource-group component's name |
| `<your-hub-vnet>` | The AzureVirtualNetwork the resolver anchors to | The network component's name |
| `<your-dns-inbound-subnet>` | The delegated AzureSubnet for the inbound endpoint | The subnet component's name |
| `<your-dns-outbound-subnet>` | The delegated AzureSubnet for the outbound endpoint | The subnet component's name |

Endpoints bill hourly from provisioning; the resolver object is free. After deploy, point on-premises forwarders at the `inbound_endpoint_ip` output and bind a forwarding ruleset to the `outbound_endpoint_id` output.
