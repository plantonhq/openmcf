# AzureVirtualNetworkPeering Terraform Module

## Overview

This Terraform module provisions one direction of an Azure virtual network
peering using the `azurerm` provider. It creates a single
`azurerm_virtual_network_peering` -- private connectivity between two
networks over the Microsoft backbone, no gateways or public IPs involved.

One resource is ONE DIRECTION, exactly as ARM models it; traffic only
flows once the reciprocal peering exists on the remote network (the two
directions can deploy concurrently). The peering is an ARM child of its
LOCAL network: the module derives the resource group and network name
from the local network's ARM ID rather than asking for them separately.

The connectivity flags and subnet-name lists update in place. Name, the
two networks, and the complete-vs-subnet-scoped and IPv6-only choices are
the peering's identity -- changing any of them replaces it. Peerings are
not tracked ARM resources, so they carry no tags.

## Resources Created

- `azurerm_virtual_network_peering.main` -- one direction of the peering

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Peering specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Peering name, unique within the local network |
| `virtual_network_id` | yes | ARM ID of the LOCAL network; resource group and network name are parsed from it |
| `remote_virtual_network_id` | yes | ARM ID of the REMOTE network; cross-subscription and global peering work unchanged |
| `allow_virtual_network_access`, `allow_forwarded_traffic`, `allow_gateway_transit`, `use_remote_gateways` | no | The four connectivity flags; defaults mirror Azure's (access on, the rest off) |
| `peer_complete_virtual_networks_enabled` | no | Azure default true; false enables subnet-scoped peering via `local_subnet_names` / `remote_subnet_names` |
| `only_ipv6_peering_enabled` | no | Peer only the IPv6 address space (dual-stack, subnet-scoped) |

## Outputs

| Output | Description |
|--------|-------------|
| `peering_id` | Full ARM ID of the peering |
| `peering_name` | The peering's name within its local network |
| `virtual_network_name` | Local network name, derived from its ARM ID |
| `resource_group_name` | Resource group of the local network, derived from its ARM ID |

## Usage

```hcl
module "peering" {
  source = "./iac/tf"

  metadata = { name = "spoke1-to-hub", org = "mycompany", env = "production" }

  spec = {
    name                      = "spoke1-to-hub"
    virtual_network_id        = "/subscriptions/xxx/resourceGroups/spoke1-rg/providers/Microsoft.Network/virtualNetworks/spoke1"
    remote_virtual_network_id = "/subscriptions/xxx/resourceGroups/hub-rg/providers/Microsoft.Network/virtualNetworks/hub"
    allow_forwarded_traffic   = true
    use_remote_gateways       = true
  }
}
```

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege actions the deploying principal needs on both the local and the remote network.
