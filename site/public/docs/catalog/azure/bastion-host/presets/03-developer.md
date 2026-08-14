---
title: "Developer Host"
description: "This preset deploys the free Developer tier: Azure-shared infrastructure attached straight to a virtual network -- no `AzureBastionSubnet` to carve, no public IP to allocate, no hourly bill."
type: "preset"
rank: "03"
presetSlug: "03-developer"
componentSlug: "bastion-host"
componentTitle: "Bastion Host"
provider: "azure"
icon: "package"
order: 3
---

# Developer Host

This preset deploys the free Developer tier: Azure-shared infrastructure attached straight to a virtual network -- no `AzureBastionSubnet` to carve, no public IP to allocate, no hourly bill.

## When to Use

- A developer connecting to dev/test VMs in one network
- Trying Bastion before committing to dedicated infrastructure

## Key Configuration Choices

- **One connection per user, no feature knobs** -- shared infrastructure means no scaling, tunneling, file copy, recording, or NSG support on the shared path
- **No peering reach** -- Developer serves only the attached network; hub-spoke topologies need Basic or above
- **Region-limited** -- Developer is offered in a subset of regions; verify availability for yours
- **Upgrading to Basic/Standard/Premium later** requires the `AzureBastionSubnet` and a public IP -- the dedicated-infrastructure ceremony arrives with the upgrade

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the host is created in | The resource-group component's name |
| `<your-virtual-network>` | The AzureVirtualNetwork the shared host attaches to | The network component's name |

Developer is free -- the right answer for one person and one network, and the wrong one the moment either multiplies.
