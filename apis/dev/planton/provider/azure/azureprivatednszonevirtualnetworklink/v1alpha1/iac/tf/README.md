# AzurePrivateDnsZoneVirtualNetworkLink Terraform Module

## Overview

This Terraform module provisions a private DNS zone virtual network link
using the `azurerm` provider. It creates a single
`azurerm_private_dns_zone_virtual_network_link`: the attachment that makes
the referenced zone resolvable from the referenced network.

`registration_enabled`, `resolution_policy`, and tags update in place;
name, zone, and network are the link's ARM identity, so changing any of
them replaces the link (a brief resolution gap for the affected network,
nothing else). Azure allows only ONE registration-enabled link per
network.

The zone's name and resource group are derived from the referenced zone's
ARM ID (a `regex()` in locals that fails the plan loudly on a malformed
ID) -- the module never asks for state that is already encoded in the
parent reference.

## Resources Created

- `azurerm_private_dns_zone_virtual_network_link.main` -- the zone-network attachment

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Virtual network link specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | The link's name under the parent zone (1-80 chars) |
| `private_dns_zone_id` | yes | ARM ID of the parent zone (name + resource group derived from it) |
| `virtual_network_id` | yes | ARM ID of the network gaining resolution |
| `registration_enabled` | no | VM A-record auto-registration; default false; one enabled link per network |
| `resolution_policy` | no | `DEFAULT` / `NX_DOMAIN_REDIRECT`; unset = Azure's per-zone-type default |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `link_id` | Full ARM ID of the link |
| `link_name` | The link's name as deployed |

## Usage

```hcl
module "zone_link" {
  source = "./iac/tf"

  metadata = {
    name = "postgres-zone-hub-link"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    name                = "hub-vnet"
    private_dns_zone_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/privateDnsZones/privatelink.postgres.database.azure.com"
    virtual_network_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/network-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet"
  }
}
```

## Required Permissions

The deploying credential needs
`Microsoft.Network/privateDnsZones/virtualNetworkLinks/write` on the zone
and `Microsoft.Network/virtualNetworks/join/action` on the network --
held via Private DNS Zone Contributor + Network Contributor, or
Contributor/Owner.
