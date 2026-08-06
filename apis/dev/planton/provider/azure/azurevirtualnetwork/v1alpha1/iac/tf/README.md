# AzureVirtualNetwork Terraform Module

## Overview

This Terraform module provisions an Azure Virtual Network using the
`azurerm` provider. It creates a single `azurerm_virtual_network` carrying
the address space and network-wide policy: multi-CIDR (or IPAM-delegated)
addressing, custom DNS, BGP community, DDoS Protection Plan attachment,
VM-to-VM encryption, flow timeout, and edge-zone placement.

Address space, DNS servers, BGP community, DDoS attachment, encryption,
flow timeout, and tags update in place; name, region, resource group, and
edge zone are the network's ARM identity, so changing any of them replaces
the network and everything inside it.

The network is deliberately just the network: subnets (`AzureSubnet`),
outbound NAT (`AzureNatGateway`), and private DNS attachments
(`AzurePrivateDnsZoneVirtualNetworkLink`) are separate composable
resources referencing this network's outputs.

## Resources Created

- `azurerm_virtual_network.main` -- the virtual network

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Virtual network specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region (a regional resource) |
| `resource_group` | yes | Resource group name |
| `name` | yes | Network name, unique within the resource group |
| `address_spaces` | one of | Self-managed CIDR blocks |
| `ip_address_pools` | one of | Network Manager IPAM pool allocations (max 2) |
| `dns_servers` | no | Custom DNS server IPs (empty = Azure's default resolver) |
| `bgp_community` | no | BGP community in `asn:community` notation |
| `ddos_protection_plan` | no | Attach an existing plan (`id` + `enable`) |
| `encryption` | no | `ALLOW_UNENCRYPTED` / `DROP_UNENCRYPTED` (unset = off) |
| `flow_timeout_in_minutes` | no | Connection-tracking timeout, 4-30 |
| `private_endpoint_vnet_policies` | no | `BASIC` (unset = ARM's `Disabled`) |
| `edge_zone` | no | Azure Edge Zone placement |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `virtual_network_id` | Full ARM ID of the network |
| `virtual_network_name` | The network's name as deployed |
| `guid` | ARM's stable network GUID |
| `address_spaces` | The ACTUAL ranges carried (IPAM-provisioned when pools delegate) |

## Usage

```hcl
module "virtual_network" {
  source = "./iac/tf"

  metadata = {
    name = "prod-network"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    region         = "eastus"
    resource_group = "prod-rg"
    name           = "prod-vnet"
    address_spaces = ["10.0.0.0/16"]
  }
}
```

## Required Permissions

The deploying credential needs `Microsoft.Network/virtualNetworks/write`
on the resource group -- held via Network Contributor, Contributor, or
Owner.
